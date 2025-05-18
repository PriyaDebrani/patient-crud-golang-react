import Loading from "@/components/Loading";
import PatientHookForm from "@/components/PatientHookForm";
import { Patient } from "@/types/patient";
import { Alert, AlertIcon, AlertTitle, Box, Text } from "@chakra-ui/react";
import { useRouter } from "next/router";
import { Fragment, FunctionComponent, useEffect, useState } from "react";
import { getPatientById, updatePatient } from "../../api";

const initialPatient: Patient = {
  id: 0,
  name: "",
  address: "",
  disease: "",
  phone: 0,
  year: 0,
  month: 0,
  date: 0,
};

interface PatientState {
  patient: Patient;
  isLoading: boolean;
  error: string | undefined;
}

const getIdFromParameters = (id: string | string[] | undefined): string => {
  if (id == undefined) {
    return "";
  }
  if (Array.isArray(id)) {
    return id[0];
  }
  return id;
};

const UpdatePatient: FunctionComponent = () => {
  const router = useRouter();
  const { id } = router.query;
  const idStr = getIdFromParameters(id);

  const [patientsState, setPatientsState] = useState<PatientState>({
    patient: initialPatient,
    isLoading: true,
    error: undefined,
  });

  const getPatient = () => {
    getPatientById(idStr)
      .then((patient) =>
        setPatientsState({
          ...patientsState,
          patient: patient,
          isLoading: false,
        })
      )
      .catch((error) => {
        if (error.response !== undefined || error.response.data !== undefined) {
          setPatientsState({
            ...patientsState,
            error: error.response.data.messages,
            isLoading: false,
          });
        } else {
          setPatientsState({
            ...patientsState,
            error: error.messages,
            isLoading: false,
          });
        }
      });
  };

  const handleUpdate = async (updatedPatient: Patient) => {
    updatePatient(idStr, updatedPatient).then(() =>
      router.push("/").catch((error) => {
        if (error.response !== undefined || error.response.data !== undefined) {
          setPatientsState({
            ...patientsState,
            error: error.response.data.messages,
            isLoading: false,
          });
        } else if (error.response.data == null) {
          setPatientsState({
            ...patientsState,
            error: error.messages,
            isLoading: false,
          });
        }
      })
    );
  };

  useEffect(() => {
    if (router.isReady && id !== undefined) {
      getPatient();
    }
  }, [router.isReady]);

  if (!router.isReady || patientsState.isLoading) {
    return <Loading />;
  }

  return (
    <Fragment>
      {patientsState.error && (
        <Alert status="error">
          <AlertIcon />
          <AlertTitle>{patientsState.error}</AlertTitle>
        </Alert>
      )}
      <PatientHookForm
        initialPatient={patientsState.patient}
        onSubmit={handleUpdate}
        isUpdate={true}
      />
    </Fragment>
  );
};

export default UpdatePatient;
