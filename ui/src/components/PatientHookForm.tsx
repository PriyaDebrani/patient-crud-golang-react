import {
  Alert,
  Box,
  Button,
  FormControl,
  FormLabel,
  Heading,
  Input,
  ListItem,
  UnorderedList,
  VStack,
} from "@chakra-ui/react";
import { superstructResolver } from "@hookform/resolvers/superstruct";
import { FunctionComponent, useState } from "react";
import { useForm } from "react-hook-form";
import { Patient, patientSchema } from "../types/patient";

export interface PatientFormProps {
  initialPatient: Patient;
  onSubmit: (patient: Patient) => Promise<void>;
  isUpdate?: boolean;
}

const PatientHookForm: FunctionComponent<PatientFormProps> = ({
  initialPatient,
  onSubmit,
  isUpdate = false,
}) => {
  const [responseErrs, setResponseErrs] = useState<string[] | undefined>(
    undefined
  );
  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<Patient>({
    defaultValues: initialPatient,
    resolver: superstructResolver(patientSchema),
  });

  const handleFormSubmit = (data: Patient) => {
    onSubmit(data).catch((err) => {
      if (err.response !== undefined || err.response.data !== undefined) {
        setResponseErrs(err.response.data.messages);
        return;
      }
      setResponseErrs(err.message);
    });
  };

  return (
    <Box p={5} maxW="full" mx="auto">
      <VStack spacing={4} align="stretch" maxW="container.md" mx="auto">
        <Heading as="h1" size="lg" textAlign="left">
          Patient Form
        </Heading>

        {Array.isArray(responseErrs) && responseErrs.length > 0 && (
          <Alert status="error">
            <UnorderedList>
              {responseErrs.map((err, idx) => (
                <ListItem key={idx}>{err}</ListItem>
              ))}
            </UnorderedList>
          </Alert>
        )}

        <form onSubmit={handleSubmit(handleFormSubmit)}>
          <FormControl id="id">
            <FormLabel>ID</FormLabel>
            <Input
              type="number"
              placeholder="ID"
              {...register("id", {
                valueAsNumber: true,
              })}
              isDisabled={isUpdate}
            />
            <span style={{ color: "red", marginTop: "8px" }}>
              {errors.id && errors.id.message}
            </span>
          </FormControl>

          <FormControl id="name">
            <FormLabel>Name</FormLabel>
            <Input type="text" placeholder="Your Name" {...register("name")} />
            <span style={{ color: "red", marginTop: "8px" }}>
              {errors.name && errors.name.message}
            </span>
          </FormControl>

          <FormControl id="disease">
            <FormLabel>Disease</FormLabel>
            <Input type="text" placeholder="Disease" {...register("disease")} />
            <span style={{ color: "red", marginTop: "8px" }}>
              {errors.disease && errors.disease.message}
            </span>
          </FormControl>

          <FormControl id="phone">
            <FormLabel>Phone</FormLabel>
            <Input
              type="number"
              placeholder="Your Phone No."
              {...register("phone", { valueAsNumber: true })}
            />
            <span style={{ color: "red", marginTop: "8px" }}>
              {errors.phone && errors.phone.message}
            </span>
          </FormControl>

          <FormControl id="year">
            <FormLabel>Year</FormLabel>
            <Input
              type="number"
              placeholder="Year"
              {...register("year", { valueAsNumber: true })}
            />
            <span style={{ color: "red", marginTop: "8px" }}>
              {errors.year && errors.year.message}
            </span>
          </FormControl>

          <FormControl id="month">
            <FormLabel>Month</FormLabel>
            <Input
              type="number"
              placeholder="Month"
              {...register("month", { valueAsNumber: true })}
            />
            <span style={{ color: "red", marginTop: "8px" }}>
              {errors.month && errors.month.message}
            </span>
          </FormControl>

          <FormControl id="date">
            <FormLabel>Date</FormLabel>
            <Input
              type="number"
              placeholder="Date"
              {...register("date", { valueAsNumber: true })}
            />
            <span style={{ color: "red", marginTop: "8px" }}>
              {errors.date && errors.date.message}
            </span>
          </FormControl>

          <FormControl id="address">
            <FormLabel>Address</FormLabel>
            <textarea
              placeholder="Your Address"
              {...register("address")}
              cols={76}
            />
            <span style={{ color: "red", marginTop: "8px" }}>
              {errors.address && errors.address.message}
            </span>
          </FormControl>

          <Button type="submit" colorScheme="blue" mt={4}>
            Submit
          </Button>
        </form>
      </VStack>
    </Box>
  );
};

export default PatientHookForm;
