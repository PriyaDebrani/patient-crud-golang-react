import { DeleteIcon, EditIcon } from "@chakra-ui/icons";
import {
  HStack,
  IconButton,
  Link,
  Table,
  TableContainer,
  Tbody,
  Td,
  Th,
  Thead,
  Tr,
} from "@chakra-ui/react";
import { Fragment, FunctionComponent, useRef, useState } from "react";
import { Patient } from "../types/patient";
import ConfirmAlertDialog from "./ConfirmAlertDialog";

export interface PatientsTableProps {
  patients: Patient[];
  deletePatient: (id: number) => void;
}

const PatientsTable: FunctionComponent<PatientsTableProps> = ({
  patients,
  deletePatient,
}) => {
  const cancelRef = useRef<HTMLButtonElement>(null);
  const [isOpen, setIsOpen] = useState<boolean>(false);
  const [selectedPatientId, setSelectedPatientId] = useState<number | null>(
    null
  );

  const onOpen = (id: number) => {
    setIsOpen(true);
    setSelectedPatientId(id);
  };

  const onClose = () => {
    setSelectedPatientId(null);
    setIsOpen(false);
  };

  const handleDelete = () => {
    if (selectedPatientId != null) {
      deletePatient(selectedPatientId);
      onClose();
    }
  };

  return (
    <Fragment>
      <TableContainer>
        <Table>
          <Thead>
            <Tr>
              <Th>Id</Th>
              <Th>Name</Th>
              <Th>Address</Th>
              <Th>Disease</Th>
              <Th>Phone</Th>
              <Th>Date</Th>
              <Th>Action</Th>
            </Tr>
          </Thead>
          <Tbody>
            {patients.map((patient) => (
              <Tr key={patient.id}>
                <Td>{patient.id}</Td>
                <Td>{patient.name}</Td>
                <Td>{patient.address}</Td>
                <Td>{patient.disease}</Td>
                <Td>{patient.phone}</Td>
                <Td>
                  {patient.date}/{patient.month}/{patient.year}
                </Td>
                <Td>
                  <HStack>
                    <IconButton
                      aria-label="Delete patient"
                      icon={<DeleteIcon />}
                      colorScheme="red"
                      onClick={() => onOpen(patient.id)}
                    />
                    <Link href={`/patients/${patient.id}`}>
                      <IconButton
                        aria-label="Edit patient"
                        icon={<EditIcon />}
                      />
                    </Link>
                  </HStack>
                </Td>
              </Tr>
            ))}
          </Tbody>
        </Table>
      </TableContainer>
      <ConfirmAlertDialog
        isOpen={isOpen}
        onClose={onClose}
        onConfirm={handleDelete}
        leastDestructiveRef={cancelRef}
        message={`Are you sure you want to delete patient id ${selectedPatientId}?`}
      />
    </Fragment>
  );
};

export default PatientsTable;
