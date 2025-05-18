import axios from "axios";
import { Patient } from "@/types/patient";

export const getPatientById = async (id: string): Promise<Patient> => {
  const response = await axios.get(`/api/patients/${id}`);
  return response.data;
};

export const updatePatient = async (
  id: string,
  patient: Patient
): Promise<void> => {
  await axios.put(`/api/patients/${id}`, patient);
};
