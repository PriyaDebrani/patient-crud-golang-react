import {
  Alert,
  AlertIcon,
  AlertTitle,
  Box,
  Button,
  FormControl,
  FormLabel,
  Heading,
  Input,
  ListItem,
  Textarea,
  UnorderedList,
  VStack,
} from "@chakra-ui/react";
import { ChangeEvent, FormEvent, FunctionComponent, useState } from "react";
import { Patient } from "../types/patient";

function validate(patient: Patient): [boolean, { [key: string]: string }] {
  let valid = true;
  const newValidationErrors: { [key: string]: string } = {};
  if (patient.id == 0 || isNaN(patient.id)) {
    newValidationErrors["id"] = "ID is required";
    valid = false;
  }
  if (patient.name == "") {
    newValidationErrors["name"] = "Name is required";
    valid = false;
  }
  if (patient.phone == 0 || isNaN(patient.phone)) {
    newValidationErrors["phone"] = "Contact No must be a valid 10-digit number";
    valid = false;
  }
  if (patient.disease == "") {
    newValidationErrors["disease"] = "Disease is required";
    valid = false;
  }
  if (patient.year == 0 || isNaN(patient.year)) {
    newValidationErrors["year"] = "Year must be a valid number";
    valid = false;
  }
  if (patient.month == 0 || isNaN(patient.month)) {
    newValidationErrors["month"] = "Month must be a valid number";
    valid = false;
  }
  if (patient.date == 0 || isNaN(patient.date)) {
    newValidationErrors["date"] = "Date must be a valid number";
    valid = false;
  }
  if (patient.address == "") {
    newValidationErrors["address"] = "Address is required";
    valid = false;
  }
  return [valid, newValidationErrors];
}

export interface PatientFormProps {
  initialPatient: Patient;
  handleSubmit: (patient: Patient) => Promise<void>;
}

const PatientForm: FunctionComponent<PatientFormProps> = ({
  initialPatient,
  handleSubmit,
}) => {
  const [patient, setPatient] = useState<Patient>(initialPatient);
  const [errorResponses, setrErrorResponses] = useState<
    string[] | string | undefined
  >(undefined);
  const [validationErrors, setValidationErrors] = useState<{
    [key: string]: string;
  }>({});

  const handleChange = (event: ChangeEvent<HTMLInputElement>) => {
    const { name, value, valueAsNumber } = event.target;
    if (
      name === "id" ||
      name === "phone" ||
      name === "year" ||
      name === "month" ||
      name === "date"
    ) {
      setPatient({
        ...patient,
        [name]: valueAsNumber,
      });
    } else {
      setPatient({
        ...patient,
        [name]: value,
      });
    }
    setValidationErrors((prevErrors) => ({
      ...prevErrors,
      [name]: "",
    }));
  };

  const handleAddressChange = (event: ChangeEvent<HTMLTextAreaElement>) => {
    const { name, value } = event.target;
    setPatient({
      ...patient,
      [name]: value,
    });
    setValidationErrors((prevErrors) => ({
      ...prevErrors,
      [name]: "",
    }));
  };

  const onSubmit = (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    const [isValid, validationErrors] = validate(patient);
    setValidationErrors(validationErrors);

    if (isValid) {
      handleSubmit(patient).catch((error) => {
        if (error.response != undefined) {
          setrErrorResponses(error.response.data.messages);
        } else {
          setrErrorResponses(error.message);
        }
      });
    }
  };

  return (
    <Box p={5} maxW="full" mx="auto">
      <VStack spacing={4} align="stretch" maxW="container.md" mx="auto">
        <Heading as="h1" size="lg" textAlign="left">
          Patient Form
        </Heading>
        {Array.isArray(errorResponses) && errorResponses.length > 0 ? (
          <Alert status="error">
            <UnorderedList>
              {errorResponses.map((err, idx) => (
                <ListItem key={idx}>{err}</ListItem>
              ))}
            </UnorderedList>
          </Alert>
        ) : (
          errorResponses != undefined && (
            <Alert status="error">
              <AlertIcon />
              <AlertTitle>{errorResponses}</AlertTitle>
            </Alert>
          )
        )}

        <form method="post" onSubmit={onSubmit}>
          <FormControl id="id">
            <FormLabel>ID</FormLabel>
            <Input
              type="number"
              placeholder="ID"
              name="id"
              value={patient.id}
              onChange={handleChange}
            />
            {validationErrors["id"] && (
              <div style={{ color: "red" }}>{validationErrors["id"]}</div>
            )}
          </FormControl>

          <FormControl id="name">
            <FormLabel>Name</FormLabel>
            <Input
              type="text"
              placeholder="Your Name"
              name="name"
              value={patient.name}
              onChange={handleChange}
            />
            {validationErrors["name"] && (
              <div style={{ color: "red" }}>{validationErrors["name"]}</div>
            )}
          </FormControl>

          <FormControl id="phone">
            <FormLabel>Contact No</FormLabel>
            <Input
              type="number"
              placeholder="Contact Number"
              name="phone"
              value={patient.phone}
              onChange={handleChange}
            />
            {validationErrors["phone"] && (
              <div style={{ color: "red" }}>{validationErrors["phone"]}</div>
            )}
          </FormControl>

          <FormControl id="disease">
            <FormLabel>Disease</FormLabel>
            <Input
              type="text"
              placeholder="Disease"
              name="disease"
              value={patient.disease}
              onChange={handleChange}
            />
            {validationErrors["disease"] && (
              <div style={{ color: "red" }}>{validationErrors["disease"]}</div>
            )}
          </FormControl>

          <FormControl id="date">
            <FormLabel>Date</FormLabel>
            <Input
              type="number"
              placeholder="Date"
              name="date"
              value={patient.date}
              onChange={handleChange}
            />
            {validationErrors["date"] && (
              <div style={{ color: "red" }}>{validationErrors["date"]}</div>
            )}
          </FormControl>

          <FormControl id="month">
            <FormLabel>Month</FormLabel>
            <Input
              type="number"
              placeholder="Month"
              name="month"
              value={patient.month}
              onChange={handleChange}
            />
            {validationErrors["month"] && (
              <div style={{ color: "red" }}>{validationErrors["month"]}</div>
            )}
          </FormControl>

          <FormControl id="year">
            <FormLabel>Year</FormLabel>
            <Input
              type="number"
              placeholder="Year"
              name="year"
              value={patient.year}
              onChange={handleChange}
            />
            {validationErrors["year"] && (
              <div style={{ color: "red" }}>{validationErrors["year"]}</div>
            )}
          </FormControl>

          <FormControl id="address" mt={4}>
            <FormLabel>Address</FormLabel>
            <Textarea
              placeholder="Address"
              name="address"
              value={patient.address}
              onChange={handleAddressChange}
            />
            {validationErrors["address"] && (
              <div style={{ color: "red" }}>{validationErrors["address"]}</div>
            )}
          </FormControl>

          <Button type="submit" colorScheme="blue" mt={4}>
            Submit
          </Button>
        </form>
      </VStack>
    </Box>
  );
};

export default PatientForm;
