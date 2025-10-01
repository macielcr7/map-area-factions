import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { apiClient } from './api'
import { toast } from 'sonner'

// Users
export function useUsers(params?: { page?: number; limit?: number; role?: string; search?: string }) {
  return useQuery({
    queryKey: ['users', params],
    queryFn: () => apiClient.getUsers(params),
  })
}

export function useCreateUser() {
  const queryClient = useQueryClient()
  
  return useMutation({
    mutationFn: (userData: any) => apiClient.createUser(userData),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['users'] })
      toast.success('Usuário criado com sucesso')
    },
    onError: (error: any) => {
      console.log(error)
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
    mutationFn: (id: string) => apiClient.deleteUser(id),
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
    mutationFn: (data: any) => apiClient.createFaction(data),
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
    mutationFn: (id: string) => apiClient.deleteFaction(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['factions'] })
      toast.success('Facção removida com sucesso')
    },
    onError: (error: any) => {
      console.log(error)
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
    mutationFn: (data: any) => apiClient.createGeometry(data),
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
    mutationFn: (id: string) => apiClient.deleteGeometry(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['geometries'] })
      toast.success('Geometria removida com sucesso')
    },
    onError: (error: any) => {
      toast.error(error.response?.data?.message || 'Erro ao remover geometria')
    },
  })
}

// Reports
export function useReports(params?: { 
  status?: string; 
  type?: string; 
  page?: number; 
  limit?: number 
}) {
  return useQuery({
    queryKey: ['reports', params],
    queryFn: () => apiClient.getReports(params),
  })
}

export function useReport(id: string) {
  return useQuery({
    queryKey: ['report', id],
    queryFn: () => apiClient.getReport(id),
    enabled: !!id,
  })
}

export function useCreateReport() {
  const queryClient = useQueryClient()
  
  return useMutation({
    mutationFn: (data: any) => apiClient.createReport(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['reports'] })
      toast.success('Relatório criado com sucesso')
    },
    onError: (error: any) => {
      toast.error(error.response?.data?.message || 'Erro ao criar relatório')
    },
  })
}

export function useUpdateReport() {
  const queryClient = useQueryClient()
  
  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: any }) => 
      apiClient.updateReport(id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['reports'] })
      toast.success('Relatório atualizado com sucesso')
    },
    onError: (error: any) => {
      toast.error(error.response?.data?.message || 'Erro ao atualizar relatório')
    },
  })
}

export function useDeleteReport() {
  const queryClient = useQueryClient()
  
  return useMutation({
    mutationFn: (id: string) => apiClient.deleteReport(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['reports'] })
      toast.success('Relatório removido com sucesso')
    },
    onError: (error: any) => {
      toast.error(error.response?.data?.message || 'Erro ao remover relatório')
    },
  })
}

// Audit
export function useAuditLogs(params?: { 
  entity?: string; 
  action?: string; 
  user_id?: string; 
  page?: number; 
  limit?: number 
}) {
  return useQuery({
    queryKey: ['audit', params],
    queryFn: () => apiClient.getAuditLogs(params),
  })
}

export function useAuditLog(id: string) {
  return useQuery({
    queryKey: ['audit-log', id],
    queryFn: () => apiClient.getAuditLog(id),
    enabled: !!id,
  })
}

// Regions
export function useRegions(params?: { state?: string; city?: string }) {
  return useQuery({
    queryKey: ['regions', params],
    queryFn: () => apiClient.getRegions(params),
  })
}

// Settings
export function useSettings() {
  return useQuery({
    queryKey: ['settings'],
    queryFn: () => apiClient.getSettings(),
  })
}

export function useUpdateSettings() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (payload: any) => apiClient.updateSettings(payload),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['settings'] })
      toast.success('Configurações salvas com sucesso')
    },
    onError: (error: any) => {
      toast.error(error.response?.data?.message || 'Erro ao atualizar configurações')
    },
  })
}
