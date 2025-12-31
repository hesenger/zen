import { Box, Stack, Loader, Text, Center, Button, Alert } from "@mantine/core";
import { useEffect, useState } from "react";
import { useSetup } from "./setup-context";

export function CompletionStep() {
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const { setupData } = useSetup();

  useEffect(() => {
    const submitSetup = async () => {
      try {
        const response = await fetch("/api/setup", {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
          },
          body: JSON.stringify(setupData),
        });

        if (!response.ok) {
          const errorData = await response.json();
          setError(errorData.error || errorData.message || "Failed to complete setup");
          setIsLoading(false);
          return;
        }

        setIsLoading(false);
      } catch (err) {
        setError(err instanceof Error ? err.message : "Failed to complete setup");
        setIsLoading(false);
      }
    };

    submitSetup();
  }, [setupData]);

  return (
    <Box mt="xl">
      <Center>
        <Stack align="center" gap="lg">
          {isLoading ? (
            <>
              <Loader size="xl" />
              <Text size="lg">Setting up your workspace...</Text>
            </>
          ) : error ? (
            <>
              <Alert color="red" title="Setup Failed" w="100%" maw={400}>
                {error}
              </Alert>
            </>
          ) : (
            <>
              <Text size="xl" fw={700} c="green">
                Setup Complete!
              </Text>
              <Text c="dimmed">Your workspace is ready to use</Text>
              <Button
                component="a"
                href="/login"
                size="lg"
                mt="md"
              >
                Go to Login
              </Button>
            </>
          )}
        </Stack>
      </Center>
    </Box>
  );
}
