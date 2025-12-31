import { Box, Stack, TextInput } from "@mantine/core";
import { useForm } from "@mantine/form";
import { forwardRef, useImperativeHandle } from "react";

interface TokensStepProps {
  onNext: (data: { githubToken: string }) => void;
}

export interface TokensStepRef {
  validate: () => boolean;
}

export const TokensStep = forwardRef<TokensStepRef, TokensStepProps>(
  ({ onNext }, ref) => {
    const form = useForm({
      initialValues: {
        githubToken: "",
      },
      validate: {
        githubToken: (value) => {
          if (
            value &&
            !value.startsWith("ghp_") &&
            !value.startsWith("github_pat_")
          )
            return "Invalid GitHub token format";
          return null;
        },
      },
    });

    useImperativeHandle(ref, () => ({
      validate: () => {
        const validation = form.validate();
        if (!validation.hasErrors) {
          onNext({
            githubToken: form.values.githubToken,
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
            label="GitHub Personal Access Token"
            placeholder="ghp_xxxxxxxxxxxx"
            description="Required for accessing GitHub repositories"
            {...form.getInputProps("githubToken")}
          />
        </Stack>
      </Box>
    );
  }
);
