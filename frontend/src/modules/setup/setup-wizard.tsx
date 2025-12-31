import {
  Container,
  Stepper,
  Box,
  Stack,
  TextInput,
  PasswordInput,
  Group,
  Button,
} from "@mantine/core";
import { useForm } from "@mantine/form";
import { useState } from "react";
import { useSetup } from "./setup-context";

export function SetupWizard() {
  const [active, setActive] = useState(0);
  const { updateSetupData } = useSetup();

  const form = useForm({
    initialValues: {
      username: "",
      password: "",
      confirmPassword: "",
    },
    validate: {
      username: (value) => {
        if (!value) return "Username is required";
        if (value.length < 3) return "Username must be at least 3 characters";
        return null;
      },
      password: (value) => {
        if (!value) return "Password is required";
        if (value.length < 3) return "Password must be at least 8 characters";
        return null;
      },
      confirmPassword: (value, values) => {
        if (!value) return "Please confirm your password";
        if (value !== values.password) return "Passwords do not match";
        return null;
      },
    },
  });

  const nextStep = () => {
    if (active === 0) {
      const validation = form.validate();
      if (!validation.hasErrors) {
        updateSetupData({
          username: form.values.username,
          password: form.values.password,
        });
        setActive((current) => current + 1);
      }
    } else {
      setActive((current) => (current < 2 ? current + 1 : current));
    }
  };

  const prevStep = () =>
    setActive((current) => (current > 0 ? current - 1 : current));

  return (
    <Container size="md" style={{ marginTop: 50 }}>
      <Stepper active={active} onStepClick={setActive}>
        <Stepper.Step label="Step 1" description="Account">
          <Box mt="xl">
            <Stack>
              <TextInput
                label="Username"
                placeholder="Enter username"
                {...form.getInputProps("username")}
              />
              <PasswordInput
                label="Password"
                placeholder="Enter password"
                {...form.getInputProps("password")}
              />
              <PasswordInput
                label="Confirm Password"
                placeholder="Confirm password"
                {...form.getInputProps("confirmPassword")}
              />
            </Stack>
          </Box>
        </Stepper.Step>
        <Stepper.Step label="Step 2" description="">
          <Box mt="xl"></Box>
        </Stepper.Step>
        <Stepper.Step label="Step 3" description="">
          <Box mt="xl"></Box>
        </Stepper.Step>
      </Stepper>

      <Group justify="center" mt="xl">
        <Button variant="default" onClick={prevStep} disabled={active === 0}>
          Back
        </Button>
        <Button onClick={nextStep} disabled={active === 2}>
          Next
        </Button>
      </Group>
    </Container>
  );
}
