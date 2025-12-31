import { Container, Title, Button, Paper } from '@mantine/core'
import { useAuth } from '../../auth-context'
import { useNavigate } from 'react-router-dom'

export default function Dashboard() {
  const { logout } = useAuth()
  const navigate = useNavigate()

  const handleLogout = async () => {
    await logout()
    navigate('/login')
  }

  return (
    <Container size="lg" my={40}>
      <Paper withBorder shadow="md" p={30} radius="md">
        <Title mb="xl">Dashboard</Title>
        <Button onClick={handleLogout}>Logout</Button>
      </Paper>
    </Container>
  )
}
