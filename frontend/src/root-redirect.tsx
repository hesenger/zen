import { useEffect, useState } from 'react'
import { Navigate } from 'react-router-dom'
import { useAuth } from './auth-context'

export default function RootRedirect() {
  const [status, setStatus] = useState<'loading' | 'ready' | 'authenticated' | 'pending-setup'>('loading')
  const { loading } = useAuth()

  useEffect(() => {
    fetch('/api/check')
      .then(res => res.json())
      .then(data => setStatus(data.status))
      .catch(() => setStatus('pending-setup'))
  }, [])

  if (status === 'loading' || loading) {
    return null
  }

  if (status === 'pending-setup') {
    return <Navigate to="/setup" replace />
  }

  if (status === 'authenticated') {
    return <Navigate to="/dashboard" replace />
  }

  return <Navigate to="/login" replace />
}
