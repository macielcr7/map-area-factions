import axios, { AxiosInstance, AxiosError } from 'axios'
import { getSession, signOut } from 'next-auth/react'

class ApiClient {
  private client: AxiosInstance

  constructor() {
    const isBrowser = typeof window !== 'undefined'
    const baseURL = isBrowser
      ? process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api/v1'
      : process.env.NEXT_INTERNAL_API_URL || process.env.NEXT_PUBLIC_API_URL || 'http://backend:8080/api/v1'

    this.client = axios.create({
      baseURL,
      timeout: 10000,
    })

    this.setupInterceptors()
  }

  private setupInterceptors() {
    // Request interceptor
    this.client.interceptors.request.use(
      async (config) => {
        const session = await getSession()
        if (session?.accessToken) {
          config.headers.Authorization = `Bearer ${session.accessToken}`
        }
        return config
      },
      (error) => {
        return Promise.reject(error)
      }
    )

    // Response interceptor
    this.client.interceptors.response.use(
      (response) => response,
      async (error: AxiosError) => {
        if (error.response?.status === 401) {
          // Token expirado, fazer logout
          await signOut()
        }
        return Promise.reject(error)
      }
    )
  }

  // Auth methods
  async login(email: string, password: string) {
    const response = await this.client.post('/auth/login', { email, password })
    return response.data
  }

  async getMe() {
    const response = await this.client.get('/auth/me')
    return response.data
  }

  // User methods
  async getUsers(params?: { page?: number; limit?: number; role?: string }) {
    const response = await this.client.get('/users', { params })
    return response.data
  }

  async createUser(userData: any) {
    const response = await this.client.post('/users', userData)
    return response.data
  }

  async updateUser(id: string, userData: any) {
    const response = await this.client.put(`/users/${id}`, userData)
    return response.data
  }

  async deleteUser(id: string) {
    const response = await this.client.delete(`/users/${id}`)
    return response.data
  }

  // Faction methods
  async getFactions(params?: { active?: boolean }) {
    const response = await this.client.get('/factions', { params })
    return response.data
  }

  async createFaction(factionData: any) {
    const response = await this.client.post('/factions', factionData)
    return response.data
  }

  async updateFaction(id: string, factionData: any) {
    const response = await this.client.put(`/factions/${id}`, factionData)
    return response.data
  }

  async deleteFaction(id: string) {
    const response = await this.client.delete(`/factions/${id}`)
    return response.data
  }

  // Geometry methods
  async getGeometries(params?: { 
    region_id?: string; 
    faction_id?: string; 
    bbox?: string 
  }) {
    const response = await this.client.get('/geometries', { params })
    return response.data
  }

  async createGeometry(geometryData: any) {
    const response = await this.client.post('/geometries', geometryData)
    return response.data
  }

  async updateGeometry(id: string, geometryData: any) {
    const response = await this.client.put(`/geometries/${id}`, geometryData)
    return response.data
  }

  async deleteGeometry(id: string) {
    const response = await this.client.delete(`/geometries/${id}`)
    return response.data
  }

  // Region methods
  async getRegions(params?: { state?: string; city?: string }) {
    const response = await this.client.get('/regions', { params })
    return response.data
  }

  // Health check
  async healthCheck() {
    const response = await this.client.get('/health')
    return response.data
  }
}

export const apiClient = new ApiClient()
