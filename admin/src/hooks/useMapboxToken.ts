'use client'

import { useState, useEffect } from 'react'

export function useMapboxToken() {
  const [token, setToken] = useState<string | null>(null)
  const [isLoading, setIsLoading] = useState(true)

  useEffect(() => {
    // Get token from environment variables (inlined at build time)
    const envToken = process.env.NEXT_PUBLIC_MAPBOX_TOKEN
    
    if (envToken && envToken !== 'pk.your-mapbox-token' && envToken !== 'pk.your-mapbox-token-here') {
      setToken(envToken)
    } else {
      setToken(null)
    }
    
    setIsLoading(false)
  }, [])

  return { token, isLoading }
}
