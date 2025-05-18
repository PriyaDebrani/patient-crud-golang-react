package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

type httpTransport struct {
	service Service
}

func newHttpTransport(service Service) *httpTransport {
	return &httpTransport{service: service}
}

type webSocketSubscriber struct {
	conn *websocket.Conn
	name string
}

func (ws *webSocketSubscriber) update(notification Notification) {
	err := ws.conn.WriteJSON(notification)
	if err != nil {
		log.Println("error sending message to websocket:", err)
	}
}

func (ws *webSocketSubscriber) getName() string {
	return ws.name
}

type errResponse struct {
	Messages []string `json:"messages"`
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (t *httpTransport) ConnectionHandler(w http.ResponseWriter, req *http.Request) {
	conn, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Println("error in connecting websocket:", err)
		return
	}
	fmt.Println("websocket connected")

	remoteAddr := conn.RemoteAddr().String()
	wsSubscriber := &webSocketSubscriber{
		conn: conn,
		name: remoteAddr,
	}

	t.service.addSubscriber(wsSubscriber)

	defer func() {
		t.service.removeSubscriber(wsSubscriber)
		conn.Close()
	}()

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			log.Println("connection closed", err)
			if websocket.IsCloseError(err) {
				return
			}
			log.Println("error reading message:", err)
			return
		}
	}
}

func writeErrResponse(w http.ResponseWriter, statusCode int, res errResponse) {
	w.WriteHeader(statusCode)
	if jsonErr := json.NewEncoder(w).Encode(res); jsonErr != nil {
		log.Println("error sending json response:", jsonErr)
	}
}

func (t *httpTransport) createPatientHandler(w http.ResponseWriter, req *http.Request) {
	var patient Patient

	if err := json.NewDecoder(req.Body).Decode(&patient); err != nil {
		writeErrResponse(w, http.StatusBadRequest, errResponse{Messages: []string{"error while decoding json"}})
		return
	}

	if err := t.service.createPatient(patient); err != nil {
		if errors.Is(err, errDuplicateId) {
			writeErrResponse(w, http.StatusConflict, errResponse{Messages: []string{errDuplicateId.Error()}})
			return
		}
		var validation *ValidationError
		if errors.As(err, &validation) {
			writeErrResponse(w, http.StatusBadRequest, errResponse{Messages: validation.Mistakes})
			return
		}

		writeErrResponse(w, http.StatusInternalServerError, errResponse{Messages: []string{err.Error()}})
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (t *httpTransport) getPatientHandler(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	id := vars["id"]
	idint, err := strconv.Atoi(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	patient, err := t.service.getPatient(idint)
	if err != nil {
		if errors.Is(err, errPatientNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(patient); err != nil {
		log.Println("error sending response:", err)
	}
}

func (t *httpTransport) getPatientsHandler(w http.ResponseWriter, req *http.Request) {
	patients, err := t.service.getPatients()
	if err != nil {
		http.Error(w, "error fetching patient", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(patients); err != nil {
		log.Println("error writing response:", err)
	}
}

func (t *httpTransport) updatePatientHandler(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	id := vars["id"]
	idint, err := strconv.Atoi(id)
	if err != nil {
		writeErrResponse(w, http.StatusBadRequest, errResponse{Messages: []string{err.Error()}})
		return
	}

	var updatedPatient Patient
	err = json.NewDecoder(req.Body).Decode(&updatedPatient)
	if err != nil {
		writeErrResponse(w, http.StatusBadRequest, errResponse{Messages: []string{err.Error()}})
		return
	}

	updatedPatient.Id = idint

	if err := t.service.updatePatient(updatedPatient); err != nil {
		if errors.Is(err, errPatientNotFound) {
			writeErrResponse(w, http.StatusNotFound, errResponse{Messages: []string{"patient not found"}})
			return
		}

		var validation *ValidationError
		if errors.As(err, &validation) {
			writeErrResponse(w, http.StatusBadRequest, errResponse{Messages: validation.Mistakes})
			return
		}

		writeErrResponse(w, http.StatusInternalServerError, errResponse{Messages: []string{err.Error()}})
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (t *httpTransport) deletePatientHandler(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	id := vars["id"]
	idint, err := strconv.Atoi(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = t.service.deletePatient(idint)
	if err != nil {
		if errors.Is(err, errPatientNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func buildRoutes(t *httpTransport) *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/api/patients", t.createPatientHandler).Methods("POST")
	router.HandleFunc("/api/patients", t.getPatientsHandler).Methods("GET")
	router.HandleFunc("/api/patients/{id}", t.getPatientHandler).Methods("GET")
	router.HandleFunc("/api/patients/{id}", t.updatePatientHandler).Methods("PUT")
	router.HandleFunc("/api/patients/{id}", t.deletePatientHandler).Methods("DELETE")
	router.HandleFunc("/websocket", t.ConnectionHandler)

	return router
}
