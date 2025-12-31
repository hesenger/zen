import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom'
import Setup from './modules/setup/setup'

function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<Navigate to="/setup" replace />} />
        <Route path="/setup" element={<Setup />} />
      </Routes>
    </BrowserRouter>
  )
}

export default App
