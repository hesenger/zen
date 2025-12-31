import { Stack, TextInput, PasswordInput, Button, Paper, Title, Container, Text } from "@mantine/core";
import { useForm } from "@mantine/form";
import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { useAuth } from "../../auth-context";

export default function Login() {
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);
  const navigate = useNavigate();
  const { login } = useAuth();

  const form = useForm({
    initialValues: {
      username: "",
      password: "",
    },
    validate: {
      username: (value) => {
        if (!value) return "Username is required";
        return null;
      },
      password: (value) => {
        if (!value) return "Password is required";
        return null;
      },
    },
  });

  const handleSubmit = form.onSubmit(async (values) => {
    setError("");
    setLoading(true);

    try {
      await login(values.username, values.password);
      navigate("/dashboard");
    } catch (err) {
      setError(err instanceof Error ? err.message : "Login failed");
    } finally {
      setLoading(false);
    }
  });

  return (
    <Container size={420} my={40}>
      <Title ta="center" mb="xl">
        Login
      </Title>
      <Paper withBorder shadow="md" p={30} radius="md">
        <form onSubmit={handleSubmit}>
          <Stack>
            {error && (
              <Text c="red" size="sm">
                {error}
              </Text>
            )}
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
            <Button type="submit" fullWidth mt="md" loading={loading}>
              Login
            </Button>
          </Stack>
        </form>
      </Paper>
    </Container>
  );
}
