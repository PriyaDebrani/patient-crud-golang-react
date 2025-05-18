import { Patient } from "@/types/patient";
import { Alert, AlertIcon, AlertTitle, Button } from "@chakra-ui/react";
import axios from "axios";
import { Fragment, FunctionComponent, useEffect, useState } from "react";
import Loading from "../components/Loading";
import PatientsTable from "../components/PatientsTable";
import router from "next/router";

interface PatientsState {
  patients: Patient[];
  isLoading: boolean;
  error: string | undefined;
}

const PatientsPage: FunctionComponent = () => {
  const [patientsState, setPatientsState] = useState<PatientsState>({
    patients: [],
    isLoading: false,
    error: undefined,
  });
  const [messages, setMessages] = useState<string[]>([]);

  useEffect(() => {
    const socket = new WebSocket("ws://localhost:8000/websocket");

    socket.onopen = () => {
      console.log("WebSocket connection opened");
    };

    socket.onmessage = (event) => {
      const notification = JSON.parse(event.data);
      console.log("WebSocket message received:", notification);
      setMessages((prev) => [...prev, notification]);
      setPatientsState({
        ...patientsState,
        patients: notification.newPatients,
      });
    };

    socket.onerror = (error) => {
      console.error("WebSocket error:", error);
    };

    socket.onclose = () => {
      console.log("WebSocket connection closed");
    };

    return () => {
      socket.close();
    };
  }, []);

  useEffect(() => {
    getPatients();
  }, []);

  const getPatients = () => {
    setPatientsState({
      ...patientsState,
      isLoading: true,
    });

    axios
      .get("/api/patients")
      .then((response) => {
        setPatientsState({
          patients: response.data,
          isLoading: false,
          error: undefined,
        });
      })
      .catch((error) => {
        console.log(error);
        if (error.response?.data) {
          setPatientsState({
            ...patientsState,
            error: error.response.data,
            isLoading: false,
          });
        } else {
          setPatientsState({
            ...patientsState,
            error: error.message,
            isLoading: false,
          });
        }
      });
  };

  const deletePatient = (id: number) => {
    axios
      .delete(`/api/patients/${id}`)
      .then(() => {
        getPatients();
      })
      .catch((error) => {
        console.log(error);
        if (error.response?.data) {
          setPatientsState({
            ...patientsState,
            error: error.response.data,
          });
        } else {
          setPatientsState({
            ...patientsState,
            error: error.message,
          });
        }
      });
  };

  useEffect(() => {
    getPatients();
  }, []);

  if (patientsState.isLoading) {
    return <Loading />;
  }

  return (
    <Fragment>
      <Button
        onClick={() =>
          router.push("/patients/add").catch((error) =>
            setPatientsState({
              ...patientsState,
              error: error,
            })
          )
        }
      >
        Add Patient
      </Button>
      {patientsState.error && (
        <Alert status="error">
          <AlertIcon />
          <AlertTitle>{patientsState.error}</AlertTitle>
        </Alert>
      )}
      <PatientsTable
        patients={patientsState.patients}
        deletePatient={deletePatient}
      />
    </Fragment>
  );
};

export default PatientsPage;
