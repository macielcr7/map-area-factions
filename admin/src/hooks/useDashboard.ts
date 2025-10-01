import { useState, useEffect } from 'react'
import { useSession } from 'next-auth/react'

interface DashboardStats {
  totalUsers: number
  totalGeometries: number
  totalFactions: number
  totalReports: number
  totalIncidents: number
  totalSubscriptions: number
}

interface ActivityData {
  name: string
  geometries: number
  users: number
  incidents: number
}

interface RecentActivity {
  id: string
  user_id: string
  user_name: string
  entity: string
  action: string
  target: string
  created_at: string
  type?: 'create' | 'update' | 'approve' | 'delete'
}

interface Alert {
  id: string
  type: 'warning' | 'info' | 'error' | 'success'
  message: string
  count?: number
  time?: string
}

export function useDashboard() {
  const { data: session } = useSession()
  console.log('🎯 useDashboard hook called', { session: !!session, accessToken: !!session?.accessToken })
  
  const [stats, setStats] = useState<DashboardStats>({
    totalUsers: 0,
    totalGeometries: 0,
    totalFactions: 0,
    totalReports: 0,
    totalIncidents: 0,
    totalSubscriptions: 0,
  })
  const [activityData, setActivityData] = useState<ActivityData[]>([])
  const [recentActivities, setRecentActivities] = useState<RecentActivity[]>([])
  const [alerts, setAlerts] = useState<Alert[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  const fetchDashboardData = async () => {
    console.log('🔍 Fetching dashboard data...', { session: !!session, accessToken: !!session?.accessToken })
    
    if (!session?.accessToken) {
      console.log('❌ No session or access token')
      return
    }

    try {
      setLoading(true)
      setError(null)

      // Fetch stats
      console.log('📊 Fetching stats from: /api/dashboard/stats')
      const statsResponse = await fetch('/api/dashboard/stats', {
        headers: {
          'Content-Type': 'application/json',
        },
      })

      console.log('📊 Stats response:', statsResponse.status, statsResponse.ok)
      if (statsResponse.ok) {
        const statsData = await statsResponse.json()
        console.log('📊 Stats data:', statsData)
        setStats(statsData)
      } else {
        console.error('❌ Stats request failed:', statsResponse.status, await statsResponse.text())
      }

      // Fetch activity data
      const activityResponse = await fetch('/api/dashboard/activity', {
        headers: {
          'Content-Type': 'application/json',
        },
      })

      if (activityResponse.ok) {
        const activityData = await activityResponse.json()
        setActivityData(activityData)
      }

      // Fetch recent activities
      const activitiesResponse = await fetch('/api/dashboard/activities', {
        headers: {
          'Content-Type': 'application/json',
        },
      })

      if (activitiesResponse.ok) {
        const activitiesData = await activitiesResponse.json()
        setRecentActivities(activitiesData)
      }

      // Fetch alerts
      const alertsResponse = await fetch('/api/dashboard/alerts', {
        headers: {
          'Content-Type': 'application/json',
        },
      })

      if (alertsResponse.ok) {
        const alertsData = await alertsResponse.json()
        setAlerts(alertsData)
      }

    } catch (err) {
      console.error('Error fetching dashboard data:', err)
      setError('Erro ao carregar dados do dashboard')
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    if (session?.accessToken) {
      fetchDashboardData()
    }
  }, [session?.accessToken])

  return {
    stats,
    activityData,
    recentActivities,
    alerts,
    loading,
    error,
    refetch: fetchDashboardData,
  }
}
