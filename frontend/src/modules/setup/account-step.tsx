import {
  Box,
  Stack,
  TextInput,
  PasswordInput,
} from "@mantine/core";
import { useForm } from "@mantine/form";
import { forwardRef, useImperativeHandle } from "react";

interface AccountStepProps {
  onNext: (data: { username: string; password: string }) => void;
}

export interface AccountStepRef {
  validate: () => boolean;
}

export const AccountStep = forwardRef<AccountStepRef, AccountStepProps>(
  ({ onNext }, ref) => {
    const form = useForm({
      initialValues: {
        username: "",
        password: "",
        confirmPassword: "",
      },
      validate: {
        username: (value) => {
          if (!value) return "Username is required";
          if (value.length < 3)
            return "Username must be at least 3 characters";
          return null;
        },
        password: (value) => {
          if (!value) return "Password is required";
          if (value.length < 3)
            return "Password must be at least 8 characters";
          return null;
        },
        confirmPassword: (value, values) => {
          if (!value) return "Please confirm your password";
          if (value !== values.password) return "Passwords do not match";
          return null;
        },
      },
    });

    useImperativeHandle(ref, () => ({
      validate: () => {
        const validation = form.validate();
        if (!validation.hasErrors) {
          onNext({
            username: form.values.username,
            password: form.values.password,
          });
          return true;
        }
        return false;
      },
    }));

    return (
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
    );
  }
);
