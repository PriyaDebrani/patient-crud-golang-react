import { min, number, object, refine, size, string } from "superstruct";

export interface Patient {
  id: number;
  name: string;
  address: string;
  disease: string;
  phone: number;
  year: number;
  month: number;
  date: number;
}

const phoneSchema = refine(number(), "phone", (value) => {
  if (/^\d{10}$/.test(String(value))) {
    return true;
  }
  return "phone number length should be 10 digits";
});

export const patientSchema = object({
  id: min(number(), 1),
  name: string(),
  address: string(),
  disease: string(),
  phone: phoneSchema,
  year: min(number(), 1000),
  month: size(number(), 1, 12),
  date: size(number(), 1, 31),
});
