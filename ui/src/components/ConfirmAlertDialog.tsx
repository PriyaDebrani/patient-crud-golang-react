import {
  AlertDialogBody,
  AlertDialogCloseButton,
  AlertDialogContent,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogOverlay,
  Button,
  AlertDialog,
} from "@chakra-ui/react";
import { FunctionComponent } from "react";

export interface ConfirmAlertDialogProps {
  isOpen: boolean;
  onClose: () => void;
  onConfirm: () => void;
  leastDestructiveRef: React.RefObject<HTMLButtonElement>;
  message: string;
}

const ConfirmAlertDialog: FunctionComponent<ConfirmAlertDialogProps> = ({
  isOpen,
  onClose,
  onConfirm,
  leastDestructiveRef,
  message,
}) => (
  <AlertDialog
    motionPreset="slideInBottom"
    isOpen={isOpen}
    leastDestructiveRef={leastDestructiveRef}
    onClose={onClose}
  >
    <AlertDialogOverlay />

    <AlertDialogContent>
      <AlertDialogHeader>Confirm Action?</AlertDialogHeader>
      <AlertDialogCloseButton />
      <AlertDialogBody>{message}</AlertDialogBody>
      <AlertDialogFooter>
        <Button ref={leastDestructiveRef} onClick={onClose}>
          No
        </Button>
        <Button colorScheme="red" ml={3} onClick={onConfirm}>
          Yes
        </Button>
      </AlertDialogFooter>
    </AlertDialogContent>
  </AlertDialog>
);

export default ConfirmAlertDialog;
