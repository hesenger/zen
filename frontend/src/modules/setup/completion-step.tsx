import { Box, Stack, Loader, Text, Center, Button } from "@mantine/core";
import { useEffect, useState } from "react";

export function CompletionStep() {
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    const timer = setTimeout(() => {
      setIsLoading(false);
    }, 2000);

    return () => clearTimeout(timer);
  }, []);

  return (
    <Box mt="xl">
      <Center>
        <Stack align="center" gap="lg">
          {isLoading ? (
            <>
              <Loader size="xl" />
              <Text size="lg">Setting up your workspace...</Text>
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
