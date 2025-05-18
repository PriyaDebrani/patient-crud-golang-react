import PatientHookForm from "@/components/PatientHookForm";
import { Patient } from "@/types/patient";
import axios from "axios";
import router from "next/router";
import { FunctionComponent } from "react";

const defaultPatient: Patient = {
  id: 0,
  name: "",
  address: "",
  disease: "",
  phone: 0,
  year: 0,
  month: 0,
  date: 0,
};

const Add: FunctionComponent = () => {
  const handleSubmit = async (patient: Patient) => {
    await axios.post("/api/patients", patient);
    await router.push("/");
  };

  return (
    <PatientHookForm onSubmit={handleSubmit} initialPatient={defaultPatient} />
  );
};

export default Add;
