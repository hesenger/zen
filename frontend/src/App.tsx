import { BrowserRouter, Routes, Route } from 'react-router-dom'
import Setup from './modules/setup/setup'
import Login from './modules/login/login'
import RootRedirect from './root-redirect'

function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<RootRedirect />} />
        <Route path="/login" element={<Login />} />
        <Route path="/setup" element={<Setup />} />
      </Routes>
    </BrowserRouter>
  )
}

export default App
