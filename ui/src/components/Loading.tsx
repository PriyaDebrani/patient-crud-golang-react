import { Spinner } from "@chakra-ui/react";
import React, { FunctionComponent } from "react";

const Loading: FunctionComponent = () => (
  <Spinner
    thickness="4px"
    speed="0.65s"
    emptyColor="gray.200"
    color="blue.500"
    size="xl"
  />
);

export default Loading;
