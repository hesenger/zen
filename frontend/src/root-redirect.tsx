import { useEffect, useState } from 'react'
import { Navigate } from 'react-router-dom'

export default function RootRedirect() {
  const [status, setStatus] = useState<'loading' | 'ready' | 'pending-setup'>('loading')

  useEffect(() => {
    fetch('/api/check')
      .then(res => res.json())
      .then(data => setStatus(data.status))
      .catch(() => setStatus('pending-setup'))
  }, [])

  if (status === 'loading') {
    return null
  }

  return <Navigate to={status === 'ready' ? '/login' : '/setup'} replace />
}
