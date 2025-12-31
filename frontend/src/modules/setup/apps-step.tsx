import {
  Box,
  Stack,
  TextInput,
  Textarea,
  Select,
  Button,
  Group,
  Paper,
  Text,
} from "@mantine/core";
import { useForm } from "@mantine/form";
import { forwardRef, useImperativeHandle } from "react";
import type { App } from "./setup-context";

interface AppsStepProps {
  onNext: (data: { apps: App[] }) => void;
}

export interface AppsStepRef {
  validate: () => boolean;
}

export const AppsStep = forwardRef<AppsStepRef, AppsStepProps>(
  ({ onNext }, ref) => {
    const form = useForm({
      initialValues: {
        apps: [] as App[],
      },
      validate: {
        apps: {
          key: (value) => {
            if (!value) return "Key is required";
            if (!value.includes("/")) return "Key must be in user/repo format";
            return null;
          },
          command: (value) => {
            if (!value) return "Command is required";
            return null;
          },
        },
      },
    });

    useImperativeHandle(ref, () => ({
      validate: () => {
        const validation = form.validate();
        if (!validation.hasErrors) {
          onNext({
            apps: form.values.apps,
          });
          return true;
        }
        return false;
      },
    }));

    const addApp = () => {
      form.insertListItem("apps", {
        provider: "github",
        key: "",
        command: "",
      });
    };

    const removeApp = (index: number) => {
      form.removeListItem("apps", index);
    };

    return (
      <Box mt="xl">
        <Stack>
          {form.values.apps.map((_, index) => (
            <Paper key={index} p="md" withBorder>
              <Stack>
                <Select
                  label="Provider"
                  data={[{ value: "github", label: "GitHub" }]}
                  {...form.getInputProps(`apps.${index}.provider`)}
                />
                <TextInput
                  label="Key"
                  placeholder="user/repo"
                  description="Repository in user/repo format"
                  {...form.getInputProps(`apps.${index}.key`)}
                />
                <Textarea
                  label="Command"
                  placeholder="Enter command"
                  minRows={3}
                  {...form.getInputProps(`apps.${index}.command`)}
                />
                <Group justify="flex-end">
                  <Button
                    variant="subtle"
                    color="red"
                    onClick={() => removeApp(index)}
                  >
                    Remove
                  </Button>
                </Group>
              </Stack>
            </Paper>
          ))}

          {form.values.apps.length === 0 && (
            <Text c="dimmed" ta="center">
              No apps added yet
            </Text>
          )}

          <Button onClick={addApp} variant="light">
            Add App
          </Button>
        </Stack>
      </Box>
    );
  }
);
