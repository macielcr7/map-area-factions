import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { apiClient } from './api'
import { toast } from 'sonner'

// Users
export function useUsers(params?: { page?: number; limit?: number; role?: string }) {
  return useQuery({
    queryKey: ['users', params],
    queryFn: () => apiClient.getUsers(params),
  })
}

export function useCreateUser() {
  const queryClient = useQueryClient()
  
  return useMutation({
    mutationFn: apiClient.createUser,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['users'] })
      toast.success('Usuário criado com sucesso')
    },
    onError: (error: any) => {
      toast.error(error.response?.data?.message || 'Erro ao criar usuário')
    },
  })
}

export function useUpdateUser() {
  const queryClient = useQueryClient()
  
  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: any }) => 
      apiClient.updateUser(id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['users'] })
      toast.success('Usuário atualizado com sucesso')
    },
    onError: (error: any) => {
      toast.error(error.response?.data?.message || 'Erro ao atualizar usuário')
    },
  })
}

export function useDeleteUser() {
  const queryClient = useQueryClient()
  
  return useMutation({
    mutationFn: apiClient.deleteUser,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['users'] })
      toast.success('Usuário removido com sucesso')
    },
    onError: (error: any) => {
      toast.error(error.response?.data?.message || 'Erro ao remover usuário')
    },
  })
}

// Factions
export function useFactions(params?: { active?: boolean }) {
  return useQuery({
    queryKey: ['factions', params],
    queryFn: () => apiClient.getFactions(params),
  })
}

export function useCreateFaction() {
  const queryClient = useQueryClient()
  
  return useMutation({
    mutationFn: apiClient.createFaction,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['factions'] })
      toast.success('Facção criada com sucesso')
    },
    onError: (error: any) => {
      toast.error(error.response?.data?.message || 'Erro ao criar facção')
    },
  })
}

export function useUpdateFaction() {
  const queryClient = useQueryClient()
  
  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: any }) => 
      apiClient.updateFaction(id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['factions'] })
      toast.success('Facção atualizada com sucesso')
    },
    onError: (error: any) => {
      toast.error(error.response?.data?.message || 'Erro ao atualizar facção')
    },
  })
}

export function useDeleteFaction() {
  const queryClient = useQueryClient()
  
  return useMutation({
    mutationFn: apiClient.deleteFaction,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['factions'] })
      toast.success('Facção removida com sucesso')
    },
    onError: (error: any) => {
      toast.error(error.response?.data?.message || 'Erro ao remover facção')
    },
  })
}

// Geometries
export function useGeometries(params?: { 
  region_id?: string; 
  faction_id?: string; 
  bbox?: string 
}) {
  return useQuery({
    queryKey: ['geometries', params],
    queryFn: () => apiClient.getGeometries(params),
  })
}

export function useCreateGeometry() {
  const queryClient = useQueryClient()
  
  return useMutation({
    mutationFn: apiClient.createGeometry,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['geometries'] })
      toast.success('Geometria criada com sucesso')
    },
    onError: (error: any) => {
      toast.error(error.response?.data?.message || 'Erro ao criar geometria')
    },
  })
}

export function useUpdateGeometry() {
  const queryClient = useQueryClient()
  
  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: any }) => 
      apiClient.updateGeometry(id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['geometries'] })
      toast.success('Geometria atualizada com sucesso')
    },
    onError: (error: any) => {
      toast.error(error.response?.data?.message || 'Erro ao atualizar geometria')
    },
  })
}

export function useDeleteGeometry() {
  const queryClient = useQueryClient()
  
  return useMutation({
    mutationFn: apiClient.deleteGeometry,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['geometries'] })
      toast.success('Geometria removida com sucesso')
    },
    onError: (error: any) => {
      toast.error(error.response?.data?.message || 'Erro ao remover geometria')
    },
  })
}

// Regions
export function useRegions(params?: { state?: string; city?: string }) {
  return useQuery({
    queryKey: ['regions', params],
    queryFn: () => apiClient.getRegions(params),
  })
}