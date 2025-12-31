import { BrowserRouter, Routes, Route } from 'react-router-dom'
import Setup from './modules/setup/setup'
import Login from './modules/login/login'
import Dashboard from './modules/dashboard/dashboard'
import RootRedirect from './root-redirect'
import ProtectedRoute from './protected-route'
import { AuthProvider } from './auth-context'

function App() {
  return (
    <BrowserRouter>
      <AuthProvider>
        <Routes>
          <Route path="/" element={<RootRedirect />} />
          <Route path="/login" element={<Login />} />
          <Route path="/setup" element={<Setup />} />
          <Route
            path="/dashboard"
            element={
              <ProtectedRoute>
                <Dashboard />
              </ProtectedRoute>
            }
          />
        </Routes>
      </AuthProvider>
    </BrowserRouter>
  )
}

export default App
