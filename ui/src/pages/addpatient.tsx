import PatientForm from "../components/PatientForm";
import { Patient } from "@/types/patient";
import axios from "axios";
import router from "next/router";
import { FunctionComponent } from "react";

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

const AddPatient: FunctionComponent = () => {
  const handleSubmit = async (patient: Patient) => {
    await axios.post("/api/patients", patient);
    await router.push("/");
  };

  return (
    <PatientForm handleSubmit={handleSubmit} initialPatient={initialPatient} />
  );
};

export default AddPatient;
