import { Box, Stack, TextInput, PasswordInput, Button, Paper, Title, Container } from "@mantine/core";
import { useForm } from "@mantine/form";

export default function Login() {
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

  const handleSubmit = form.onSubmit((values) => {
    console.log(values);
  });

  return (
    <Container size={420} my={40}>
      <Title ta="center" mb="xl">
        Login
      </Title>
      <Paper withBorder shadow="md" p={30} radius="md">
        <form onSubmit={handleSubmit}>
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
            <Button type="submit" fullWidth mt="md">
              Login
            </Button>
          </Stack>
        </form>
      </Paper>
    </Container>
  );
}
