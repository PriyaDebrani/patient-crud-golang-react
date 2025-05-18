package main

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
)

func TestTransport_createPatient(t *testing.T) {
	tests := []struct {
		name           string
		url            string
		requestBody    string
		existingBody   []Patient
		wantResponse   string
		wantErr        error
		wantStatusCode int
	}{
		{
			name: "invalid json syntax for name :NEG",
			url:  "/api/patients",
			requestBody: `
			{
				"id": 98,
				"name": abc,
				"address": "SRT",
				"phone": 123,
				"disease": "fever",
				"year": 2012,
				"month": 10,
				"date" :12,
			}
			`,
			wantResponse: `{
				"messages": ["error while decoding json"]
				}`,
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name: "duplicate id :NEG",
			url:  "/api/patients",
			requestBody: `
			{
				"id": 1,
				"name": "abc",
				"address": "SRT",
				"disease": "fever",
				"phone": 123,
				"year": 2012,
				"month": 10,
				"date" :12
			}
			`,
			existingBody: []Patient{
				{
					Id:      1,
					Name:    "abc",
					Address: "SRT",
					Disease: "fever",
					Phone:   123,
					Year:    2024,
					Month:   10,
					Date:    12,
				},
			},
			wantErr:        errDuplicateId,
			wantStatusCode: http.StatusConflict,
		},
		{
			name: "empty name :NEG",
			url:  "/api/patients",
			requestBody: `
			{
				"id": 1,
				"name": "",
				"address": "surat",
				"phone": 12345,
				"disease": "fever",
				"year": 2024,
				"month": 2,
				"date" :12
			}
			`,
			wantResponse: `{
				"messages": ["name cannot be empty"]
				}`,

			wantStatusCode: http.StatusBadRequest,
		},
		{
			name: "empty address :NEG",
			url:  "/api/patients",
			requestBody: `
			{
				"id": 1,
				"name": "dfrf",
				"address": "",
				"phone": 12345,
				"disease": "fever",
				"year": 2024,
				"month": 2,
				"date" :12
			}
			`,
			wantResponse: `{
				"messages": ["address cannot be empty"]
				}`,
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name: "empty phone :NEG",
			url:  "/api/patients",
			requestBody: `
			{
				"id": 1,
				"name": "fwer",
				"address": "ewre",
				"phone": 0,
				"disease": "fever",
				"year": 2024,
				"month": 2,
				"date" :12
			}
			`,
			wantResponse: `{
				"messages": ["contact Number should be postive"]
				}`,
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name: "multiple validation errors :NEG",
			url:  "/api/patients",
			requestBody: `
			{
				"id": 1,
				"name": "",
				"address": "ewre",
				"phone": 0,
				"disease": "",
				"year": 0,
				"month": 2,
				"date" :0
			}
			`,
			wantResponse: `{
				"messages": ["name cannot be empty",
							"disease cannot be empty",
							"contact Number should be postive",
							"year should be positive or negative",
							"date should be positive or less than 32"]
				}`,
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name: "patient created :POS",
			url:  "/api/patients",
			requestBody: `
			{
				"id": 1,
				"name": "priya",
				"address": "surat",
				"phone": 12345,
				"disease": "fever",
				"year": 2024,
				"month": 2,
				"date" :12
			}
			`,
			wantStatusCode: http.StatusCreated,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newInMemoryRepository()
			service := newPatientsService(repo)
			transport := newHttpTransport(service)
			repo.patients = tt.existingBody

			router := buildRoutes(transport)
			body := strings.NewReader(tt.requestBody)

			res := httptest.NewRecorder()
			req := httptest.NewRequest("POST", tt.url, body)
			router.ServeHTTP(res, req)

			gotResponse := res.Body.String()
			assert.Equal(t, tt.wantStatusCode, res.Code, "expect status code to match")
			if tt.wantResponse != "" {
				assert.JSONEq(t, tt.wantResponse, gotResponse, "expect response body to match")
			}
		})
	}
}

func TestTransport_getPatients(t *testing.T) {
	testTime := time.Date(2024, time.August, 20, 17, 0, 0, 0, time.UTC)
	tests := []struct {
		name             string
		url              string
		existingPatients []Patient
		wantResponse     string
		wantStatusCode   int
	}{
		{
			name:             "empty list :POS",
			url:              "/api/patients",
			existingPatients: []Patient{},
			wantResponse:     "[]",
			wantStatusCode:   http.StatusOK,
		},
		{
			name: "patients list :POS",
			url:  "/api/patients",
			existingPatients: []Patient{
				{
					Id:        1,
					Name:      "priya",
					Address:   "surat",
					Disease:   "fever",
					Phone:     12345,
					Year:      2024,
					Month:     2,
					Date:      12,
					CreatedAt: testTime,
					UpdatedAt: testTime,
				},
				{
					Id:        2,
					Name:      "priya",
					Address:   "surat",
					Disease:   "fever",
					Phone:     12345,
					Year:      2024,
					Month:     2,
					Date:      12,
					CreatedAt: testTime,
					UpdatedAt: testTime,
				},
			},
			wantResponse: `
			[
				{
					"id": 1,
					"name": "priya",
					"address": "surat",
					"phone": 12345,
					"disease": "fever",
					"year": 2024,
					"month": 2,
					"date": 12,
					"createdAt" : "2024-08-20T17:00:00Z",
					"updatedAt" : "2024-08-20T17:00:00Z"
				},
				{
					"id": 2,
					"name": "priya",
					"address": "surat",
					"phone": 12345,
					"disease": "fever",
					"year": 2024,
					"month": 2,
					"date": 12,
					"createdAt" : "2024-08-20T17:00:00Z",
					"updatedAt" : "2024-08-20T17:00:00Z"
				}
			]
			`,
			wantStatusCode: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newInMemoryRepository()
			repo.patients = tt.existingPatients
			service := newPatientsService(repo)
			transport := newHttpTransport(service)

			router := buildRoutes(transport)

			res := httptest.NewRecorder()
			req := httptest.NewRequest("GET", tt.url, nil)
			router.ServeHTTP(res, req)

			gotResponse := res.Body.String()

			assert.Equal(t, tt.wantStatusCode, res.Code, "expect status code to match")
			assert.JSONEq(t, tt.wantResponse, gotResponse, "expect patients to match")
		})
	}
}

func TestTransport_getPatient(t *testing.T) {
	testTime := time.Date(2024, time.August, 20, 17, 0, 0, 0, time.UTC)
	tests := []struct {
		name               string
		url                string
		existingPatients   []Patient
		expectedResponse   string
		expectedStatusCode int
	}{
		{
			name:               "patient does not exist :NEG",
			url:                "/api/patients/1",
			existingPatients:   []Patient{},
			expectedStatusCode: http.StatusNotFound,
		},
		{
			name: "patient exists :POS",
			url:  "/api/patients/1",
			existingPatients: []Patient{
				{
					Id:        1,
					Name:      "priya",
					Address:   "surat",
					Disease:   "fever",
					Phone:     12345,
					Year:      2024,
					Month:     2,
					Date:      12,
					CreatedAt: testTime,
					UpdatedAt: testTime,
				},
			},
			expectedResponse: `{
				"id": 1,
				"name": "priya",
				"address": "surat",
				"phone": 12345,
				"disease": "fever",
				"year": 2024,
				"month": 2,
				"date" :12,
				"createdAt" : "2024-08-20T17:00:00Z",
				"updatedAt" : "2024-08-20T17:00:00Z"
			}`,
			expectedStatusCode: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newInMemoryRepository()
			repo.patients = tt.existingPatients
			service := newPatientsService(repo)
			transport := newHttpTransport(service)

			router := buildRoutes(transport)

			res := httptest.NewRecorder()
			req := httptest.NewRequest("GET", tt.url, nil)
			router.ServeHTTP(res, req)

			gotResponse := res.Body.String()
			assert.Equal(t, tt.expectedStatusCode, res.Code, "status code mismatched")
			if tt.expectedStatusCode == http.StatusOK {
				assert.JSONEq(t, tt.expectedResponse, gotResponse, "expected response mismatched")
			}
		})
	}
}

func TestTransport_updatePatient(t *testing.T) {
	tests := []struct {
		name             string
		url              string
		requestBody      string
		existingPatients []Patient
		wantResponse     string
		wantStatusCode   int
	}{
		{
			name: "patient not found :NEG",
			url:  "/api/patients/1",
			requestBody: `
			{
				"id": 1,
				"name": "abc",
				"address": "surat",
				"phone": 12345,
				"disease": "fever",
				"year": 2024,
				"month": 2,
				"date" :12
			}
			`,
			existingPatients: []Patient{},
			wantResponse: `{
				"messages": ["patient not found"]
				}`,
			wantStatusCode: http.StatusNotFound,
		},
		{
			name: "empty name :NEG",
			url:  "/api/patients/1",
			requestBody: `
			{
				"id": 1,
				"name": "",
				"address": "surat",
				"phone": 12345,
				"disease": "fever",
				"year": 2024,
				"month": 2,
				"date" :12
			}
			`,
			existingPatients: []Patient{
				{
					Id:      1,
					Name:    "priya",
					Address: "surat",
					Disease: "fever",
					Phone:   12345,
					Year:    2024,
					Month:   2,
					Date:    12,
				},
			},
			wantResponse: `{
				"messages": ["name cannot be empty"]
				}`,
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name: "empty address :NEG",
			url:  "/api/patients/1",
			requestBody: `
			{
				"id": 1,
				"name": "gdfgh",
				"address": "",
				"phone": 12345,
				"disease": "fever",
				"year": 2024,
				"month": 2,
				"date" :12
			}
			`,
			existingPatients: []Patient{
				{
					Id:      1,
					Name:    "priya",
					Address: "surat",
					Disease: "fever",
					Phone:   12345,
					Year:    2024,
					Month:   2,
					Date:    12,
				},
			},
			wantResponse: `{
				"messages": ["address cannot be empty"]
				}`,
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name: "empty phone :NEG",
			url:  "/api/patients/1",
			requestBody: `
			{
				"id": 1,
				"name": "fwer",
				"address": "ewre",
				"phone": 0,
				"disease": "fever",
				"year": 2024,
				"month": 2,
				"date" :12
			}
			`,
			existingPatients: []Patient{
				{
					Id:      1,
					Name:    "priya",
					Address: "surat",
					Disease: "fever",
					Phone:   12345,
					Year:    2024,
					Month:   2,
					Date:    12,
				},
			},
			wantResponse: `{
				"messages": ["contact Number should be postive"]
				}`,
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name: "patient found :POS",
			url:  "/api/patients/1",
			requestBody: `
			{
				"id": 1,
				"name": "priya",
				"address": "SRT",
				"phone": 12345,
				"disease": "fever",
				"year": 2024,
				"month": 2,
				"date" :12
			}
			`,
			existingPatients: []Patient{
				{
					Id:      1,
					Name:    "priya",
					Address: "surat",
					Disease: "fever",
					Phone:   12345,
					Year:    2024,
					Month:   2,
					Date:    12,
				},
			},
			wantStatusCode: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newInMemoryRepository()
			repo.patients = tt.existingPatients
			service := newPatientsService(repo)
			transport := newHttpTransport(service)

			router := buildRoutes(transport)
			body := strings.NewReader(tt.requestBody)

			res := httptest.NewRecorder()
			req := httptest.NewRequest("PUT", tt.url, body)
			router.ServeHTTP(res, req)

			gotResponse := res.Body.String()

			assert.Equal(t, tt.wantStatusCode, res.Code, "expect status code to match")
			if tt.wantResponse != "" {
				assert.JSONEq(t, tt.wantResponse, gotResponse, "expect response body to match")
			}
		})
	}
}

func TestTransport_deletePatient(t *testing.T) {
	tests := []struct {
		name             string
		url              string
		existingPatients []Patient
		wantPatients     []Patient
		wantStatusCode   int
	}{

		{
			name: "patient not found :NEG",
			url:  "/api/patients/2",
			existingPatients: []Patient{
				{
					Id:      1,
					Name:    "abc",
					Address: "srt",
					Phone:   12345,
					Disease: "Cold",
					Year:    2024,
					Month:   2,
					Date:    12,
				},
			},
			wantPatients: []Patient{
				{
					Id:      1,
					Name:    "abc",
					Address: "srt",
					Phone:   12345,
					Disease: "Cold",
					Year:    2024,
					Month:   2,
					Date:    12,
				},
			},
			wantStatusCode: http.StatusNotFound,
		},
		{
			name: "patient deleted :POS",
			url:  "/api/patients/2",
			existingPatients: []Patient{
				{
					Id:      1,
					Name:    "abc",
					Address: "srt",
					Phone:   12345,
					Disease: "cold",
					Year:    2024,
					Month:   2,
					Date:    12,
				},
				{
					Id:      2,
					Name:    "abc",
					Address: "srt",
					Phone:   12345,
					Disease: "cold",
					Year:    2024,
					Month:   2,
					Date:    12,
				},
			},
			wantPatients: []Patient{
				{
					Id:      1,
					Name:    "abc",
					Address: "srt",
					Phone:   12345,
					Disease: "cold",
					Year:    2024,
					Month:   2,
					Date:    12,
				},
			},
			wantStatusCode: http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newInMemoryRepository()
			repo.patients = tt.existingPatients
			service := newPatientsService(repo)
			transport := newHttpTransport(service)

			router := buildRoutes(transport)

			res := httptest.NewRecorder()
			req := httptest.NewRequest("DELETE", tt.url, nil)
			router.ServeHTTP(res, req)

			assert.Equal(t, tt.wantStatusCode, res.Code, "expect status code to match")
			assert.Equal(t, tt.wantPatients, repo.patients, "expect patients to match")
		})
	}
}

func TestWebSocket_createPatient(t *testing.T) {
	repo := newInMemoryRepository()
	service := newPatientsService(repo)
	transport := newHttpTransport(service)
	router := buildRoutes(transport)
	ts := httptest.NewServer(router)
	defer ts.Close()

	wsURL := "ws" + ts.URL[len("http"):] + "/websocket"
	conn, res, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Errorf("Failed to connect to WebSocket: %v", err)
	}

	assert.Equal(t, http.StatusSwitchingProtocols, res.StatusCode, "status code mismatched")

	defer conn.Close()
	log.Println("WebSocket connection established successfully")

	requestBody := `
		{
			"id": 1,
			"name": "abc",
			"address": "surat",
			"disease": "fever",
			"phone": 123,
			"year": 2012,
			"month": 10,
			"date" :12
		}
	`

	resp, err := http.Post(ts.URL+"/api/patients", "application/json", bytes.NewBufferString(requestBody))

	if err != nil {
		t.Errorf("Failed to send HTTP request: %v", err)
	}
	defer resp.Body.Close()

	wantNotification := Notification{
		Message: "New patient added with id: 1",
		NewPatients: []Patient{
			{
				Id:      1,
				Name:    "abc",
				Address: "surat",
				Disease: "fever",
				Phone:   123,
				Year:    2012,
				Month:   10,
				Date:    12,
			},
		},
	}

	var notification Notification
	if err := conn.ReadJSON(&notification); err != nil {
		t.Errorf("Failed to read message from WebSocket: %v", err)
	}

	assert.Equal(t, wantNotification.Message, notification.Message, "expect message to match")
	for i := range notification.NewPatients {
		assertPatientEqual(t, wantNotification.NewPatients[i], notification.NewPatients[i])
	}
}

func TestWebSocket_deletePatient(t *testing.T) {
	wantNotification := `
		{
			"message": "Patient removed with id: 1",
			"newPatients": []
		}
	`

	repo := newInMemoryRepository()
	repo.patients = []Patient{
		{
			Id:      1,
			Name:    "abc",
			Address: "srt",
			Phone:   12345,
			Disease: "cold",
			Year:    2024,
			Month:   2,
			Date:    12,
		},
	}
	service := newPatientsService(repo)
	transport := newHttpTransport(service)
	router := buildRoutes(transport)
	ts := httptest.NewServer(router)
	defer ts.Close()

	wsURL := "ws" + ts.URL[len("http"):] + "/websocket"
	conn, res, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Errorf("Failed to connect to WebSocket: %v", err)
	}
	defer conn.Close()
	log.Println("WebSocket connection established successfully")

	assert.Equal(t, http.StatusSwitchingProtocols, res.StatusCode, "status code mismatched")

	req, err := http.NewRequest("DELETE", ts.URL+"/api/patients/1", nil)
	if err != nil {
		t.Fatalf("Failed to create HTTP request: %v", err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("failed to send HTTP request: %v", err)
	}
	defer resp.Body.Close()

	messageType, msg, err := conn.ReadMessage()
	if err != nil {
		t.Errorf("failed to read message from websocket: %v", err)
	}

	assert.Equal(t, websocket.TextMessage, messageType, "message type should be text")
	notification := string(msg)
	assert.JSONEq(t, wantNotification, notification, "expect WebSocket notification to match")
}

func TestWebSocket_updatePatient(t *testing.T) {
	repo := newInMemoryRepository()
	repo.patients = []Patient{
		{
			Id:      1,
			Name:    "abc",
			Address: "srt",
			Disease: "cold",
			Phone:   12345,
			Year:    2024,
			Month:   2,
			Date:    12,
		},
	}
	service := newPatientsService(repo)
	transport := newHttpTransport(service)
	router := buildRoutes(transport)

	ts := httptest.NewServer(router)
	defer ts.Close()

	wsURL := "ws" + ts.URL[len("http"):] + "/websocket"
	conn, res, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Errorf("Failed to connect to WebSocket: %v", err)
	}
	defer conn.Close()
	log.Println("WebSocket connection established successfully")

	assert.Equal(t, http.StatusSwitchingProtocols, res.StatusCode, "status code mismatched")

	requestBody := `
	{
		"id": 1,
		"name": "priya",
		"address": "surat",
		"phone": 54321,
		"disease": "cold",
		"year": 2025,
		"month": 3,
		"date": 15
	}
	`
	wantNotification := Notification{
		Message: "Patient updated with id: 1",
		NewPatients: []Patient{
			{
				Id:      1,
				Name:    "priya",
				Address: "surat",
				Disease: "cold",
				Phone:   54321,
				Year:    2025,
				Month:   3,
				Date:    15,
			},
		},
	}

	req, err := http.NewRequest("PUT", ts.URL+"/api/patients/1", bytes.NewBufferString(requestBody))
	if err != nil {
		t.Fatalf("Failed to create HTTP request: %v", err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Failed to send HTTP request: %v", err)
	}
	defer resp.Body.Close()

	var notification Notification
	if err := conn.ReadJSON(&notification); err != nil {
		t.Errorf("Failed to read message from WebSocket: %v", err)
	}

	assert.Equal(t, notification.Message, wantNotification.Message, "message type should be text")
	for i := range notification.NewPatients {
		assertPatientEqual(t, wantNotification.NewPatients[i], notification.NewPatients[i])
	}
}
