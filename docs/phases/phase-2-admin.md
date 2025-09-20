# Fase 2: Admin Interface - Guia de Implementação

## 🎯 Objetivo da Fase
Desenvolver a interface administrativa completa em Next.js 14, incluindo editor de mapas interativo, gestão de usuários, CRUD de facções e dashboard com métricas.

## 📋 Checklist Detalhado

### Day 1-3: Setup e Estrutura Base

#### 1. Setup Next.js 14 e Dependencies
```bash
# Instalar dependências core
cd admin
npm install next@14 react@18 react-dom@18 typescript
npm install @types/node @types/react @types/react-dom

# UI Framework
npm install tailwindcss postcss autoprefixer
npm install @radix-ui/react-* lucide-react class-variance-authority
npm install @shadcn/ui

# Forms e Validation
npm install react-hook-form @hookform/resolvers zod

# HTTP Client e State
npm install axios @tanstack/react-query zustand

# Maps
npm install mapbox-gl @mapbox/mapbox-gl-draw
npm install @types/mapbox-gl @types/mapbox__mapbox-gl-draw

# Auth
npm install next-auth
```

#### 2. Configuração Tailwind e Shadcn/ui
```typescript
// admin/tailwind.config.js
import type { Config } from 'tailwindcss'

const config: Config = {
  content: [
    './src/pages/**/*.{js,ts,jsx,tsx,mdx}',
    './src/components/**/*.{js,ts,jsx,tsx,mdx}',
    './src/app/**/*.{js,ts,jsx,tsx,mdx}',
  ],
  theme: {
    extend: {
      colors: {
        border: "hsl(var(--border))",
        input: "hsl(var(--input))",
        ring: "hsl(var(--ring))",
        background: "hsl(var(--background))",
        foreground: "hsl(var(--foreground))",
        primary: {
          DEFAULT: "hsl(var(--primary))",
          foreground: "hsl(var(--primary-foreground))",
        },
        // Cores específicas para facções
        faction: {
          tdn: "#ff0000",
          cv: "#cc0000", 
          gde: "#0000ff",
          neutral: "#808080",
        }
      },
    },
  },
  plugins: [require("tailwindcss-animate")],
}
export default config
```

#### 3. Configuração TypeScript e ESLint
```json
// admin/tsconfig.json
{
  "compilerOptions": {
    "target": "es5",
    "lib": ["dom", "dom.iterable", "es6"],
    "allowJs": true,
    "skipLibCheck": true,
    "strict": true,
    "noEmit": true,
    "esModuleInterop": true,
    "module": "esnext",
    "moduleResolution": "bundler",
    "resolveJsonModule": true,
    "isolatedModules": true,
    "jsx": "preserve",
    "incremental": true,
    "plugins": [
      {
        "name": "next"
      }
    ],
    "baseUrl": ".",
    "paths": {
      "@/*": ["./src/*"],
      "@/components/*": ["./src/components/*"],
      "@/lib/*": ["./src/lib/*"],
      "@/hooks/*": ["./src/hooks/*"],
      "@/types/*": ["./src/types/*"]
    }
  },
  "include": ["next-env.d.ts", "**/*.ts", "**/*.tsx", ".next/types/**/*.ts"],
  "exclude": ["node_modules"]
}
```

**Tarefas:**
- [ ] Configurar Next.js 14 com App Router
- [ ] Setup Tailwind CSS + Shadcn/ui
- [ ] Configurar TypeScript paths
- [ ] Setup ESLint e Prettier
- [ ] Configurar estrutura de pastas

### Day 4-6: Authentication e Layout

#### 4. Configuração NextAuth
```typescript
// admin/src/app/api/auth/[...nextauth]/route.ts
import NextAuth from 'next-auth'
import CredentialsProvider from 'next-auth/providers/credentials'
import { AuthOptions } from 'next-auth'

const authOptions: AuthOptions = {
  providers: [
    CredentialsProvider({
      name: 'credentials',
      credentials: {
        email: { label: 'Email', type: 'email' },
        password: { label: 'Password', type: 'password' }
      },
      async authorize(credentials) {
        if (!credentials?.email || !credentials?.password) {
          return null
        }

        try {
          const response = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/auth/login`, {
            method: 'POST',
            headers: {
              'Content-Type': 'application/json',
            },
            body: JSON.stringify({
              email: credentials.email,
              password: credentials.password,
            }),
          })

          const data = await response.json()

          if (!response.ok) {
            return null
          }

          return {
            id: data.user.id,
            email: data.user.email,
            name: data.user.name,
            role: data.user.role,
            accessToken: data.access_token,
            refreshToken: data.refresh_token,
          }
        } catch (error) {
          return null
        }
      }
    })
  ],
  callbacks: {
    async jwt({ token, user }) {
      if (user) {
        token.accessToken = user.accessToken
        token.refreshToken = user.refreshToken
        token.role = user.role
      }
      return token
    },
    async session({ session, token }) {
      session.accessToken = token.accessToken
      session.refreshToken = token.refreshToken
      session.user.role = token.role
      return session
    }
  },
  pages: {
    signIn: '/auth/login',
  },
  session: {
    strategy: 'jwt',
  },
}

const handler = NextAuth(authOptions)
export { handler as GET, handler as POST }
```

#### 5. Páginas de Autenticação
```typescript
// admin/src/app/(auth)/login/page.tsx
'use client'

import { useState } from 'react'
import { signIn, getSession } from 'next-auth/react'
import { useRouter } from 'next/navigation'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Card, CardHeader, CardTitle, CardContent } from '@/components/ui/card'
import { Alert, AlertDescription } from '@/components/ui/alert'
import { Loader2 } from 'lucide-react'

export default function LoginPage() {
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [error, setError] = useState('')
  const [isLoading, setIsLoading] = useState(false)
  const router = useRouter()

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setIsLoading(true)
    setError('')

    try {
      const result = await signIn('credentials', {
        email,
        password,
        redirect: false,
      })

      if (result?.error) {
        setError('Credenciais inválidas')
        return
      }

      // Verificar se login foi bem-sucedido
      const session = await getSession()
      if (session) {
        router.push('/dashboard')
      }
    } catch (err) {
      setError('Erro ao fazer login')
    } finally {
      setIsLoading(false)
    }
  }

  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-50">
      <Card className="w-full max-w-md">
        <CardHeader>
          <CardTitle className="text-2xl text-center">
            Map Area Factions
          </CardTitle>
          <p className="text-center text-gray-600">
            Painel Administrativo
          </p>
        </CardHeader>
        <CardContent>
          <form onSubmit={handleSubmit} className="space-y-4">
            {error && (
              <Alert variant="destructive">
                <AlertDescription>{error}</AlertDescription>
              </Alert>
            )}
            
            <div className="space-y-2">
              <label htmlFor="email" className="text-sm font-medium">
                Email
              </label>
              <Input
                id="email"
                type="email"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                required
                disabled={isLoading}
              />
            </div>
            
            <div className="space-y-2">
              <label htmlFor="password" className="text-sm font-medium">
                Senha
              </label>
              <Input
                id="password"
                type="password"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                required
                disabled={isLoading}
              />
            </div>
            
            <Button 
              type="submit" 
              className="w-full"
              disabled={isLoading}
            >
              {isLoading && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
              Entrar
            </Button>
          </form>
        </CardContent>
      </Card>
    </div>
  )
}
```

#### 6. Layout Principal
```typescript
// admin/src/components/layout/sidebar.tsx
'use client'

import Link from 'next/link'
import { usePathname } from 'next/navigation'
import { cn } from '@/lib/utils'
import { 
  LayoutDashboard, 
  Map, 
  Users, 
  Flag, 
  FileText, 
  Settings,
  Activity
} from 'lucide-react'

const sidebarItems = [
  {
    title: 'Dashboard',
    href: '/dashboard',
    icon: LayoutDashboard,
  },
  {
    title: 'Mapas',
    href: '/maps',
    icon: Map,
  },
  {
    title: 'Usuários',
    href: '/users',
    icon: Users,
  },
  {
    title: 'Facções',
    href: '/factions',
    icon: Flag,
  },
  {
    title: 'Auditoria',
    href: '/audit',
    icon: Activity,
  },
  {
    title: 'Relatórios',
    href: '/reports',
    icon: FileText,
  },
  {
    title: 'Configurações',
    href: '/settings',
    icon: Settings,
  },
]

export function Sidebar() {
  const pathname = usePathname()

  return (
    <div className="pb-12 w-64">
      <div className="space-y-4 py-4">
        <div className="px-3 py-2">
          <h2 className="mb-2 px-4 text-lg font-semibold">
            Map Factions Admin
          </h2>
          <div className="space-y-1">
            {sidebarItems.map((item) => (
              <Link
                key={item.href}
                href={item.href}
                className={cn(
                  'flex items-center rounded-lg px-3 py-2 text-sm font-medium hover:bg-accent hover:text-accent-foreground',
                  pathname === item.href ? 'bg-accent text-accent-foreground' : 'transparent'
                )}
              >
                <item.icon className="mr-2 h-4 w-4" />
                {item.title}
              </Link>
            ))}
          </div>
        </div>
      </div>
    </div>
  )
}
```

**Tarefas:**
- [ ] Configurar NextAuth com backend
- [ ] Implementar página de login responsiva
- [ ] Criar layout principal com sidebar
- [ ] Implementar header com user menu
- [ ] Configurar middleware de proteção de rotas
- [ ] Implementar logout functionality

### Day 7-9: API Integration e State Management

#### 7. API Client Configuration
```typescript
// admin/src/lib/api.ts
import axios, { AxiosInstance, AxiosError } from 'axios'
import { getSession, signOut } from 'next-auth/react'

class ApiClient {
  private client: AxiosInstance

  constructor() {
    this.client = axios.create({
      baseURL: process.env.NEXT_PUBLIC_API_URL,
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
}

export const apiClient = new ApiClient()
```

#### 8. React Query Setup
```typescript
// admin/src/lib/queries.ts
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
```

**Tarefas:**
- [ ] Configurar Axios client com interceptors
- [ ] Implementar refresh token automático
- [ ] Setup React Query para cache
- [ ] Criar hooks para todas as entidades
- [ ] Implementar error handling global
- [ ] Configurar TypeScript types para API

### Day 10-12: CRUD Interfaces

#### 9. User Management Interface
```typescript
// admin/src/app/users/page.tsx
'use client'

import { useState } from 'react'
import { useUsers, useCreateUser, useUpdateUser, useDeleteUser } from '@/lib/queries'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Card, CardHeader, CardTitle, CardContent } from '@/components/ui/card'
import { 
  Table, 
  TableBody, 
  TableCell, 
  TableHead, 
  TableHeader, 
  TableRow 
} from '@/components/ui/table'
import { Badge } from '@/components/ui/badge'
import { UserDialog } from '@/components/users/user-dialog'
import { MoreHorizontal, Plus, Search } from 'lucide-react'

export default function UsersPage() {
  const [searchTerm, setSearchTerm] = useState('')
  const [selectedRole, setSelectedRole] = useState<string>('')
  const [isDialogOpen, setIsDialogOpen] = useState(false)
  const [editingUser, setEditingUser] = useState<any>(null)

  const { data: usersData, isLoading } = useUsers({
    role: selectedRole || undefined,
  })

  const createUserMutation = useCreateUser()
  const updateUserMutation = useUpdateUser()

  const filteredUsers = usersData?.data?.filter((user: any) =>
    user.name.toLowerCase().includes(searchTerm.toLowerCase()) ||
    user.email.toLowerCase().includes(searchTerm.toLowerCase())
  ) || []

  const handleCreateUser = (userData: any) => {
    createUserMutation.mutate(userData, {
      onSuccess: () => {
        setIsDialogOpen(false)
      }
    })
  }

  const handleUpdateUser = (userData: any) => {
    if (editingUser) {
      updateUserMutation.mutate(
        { id: editingUser.id, data: userData },
        {
          onSuccess: () => {
            setIsDialogOpen(false)
            setEditingUser(null)
          }
        }
      )
    }
  }

  const getRoleColor = (role: string) => {
    switch (role) {
      case 'admin': return 'destructive'
      case 'moderator': return 'default'
      case 'collaborator': return 'secondary'
      default: return 'outline'
    }
  }

  return (
    <div className="space-y-6">
      <div className="flex justify-between items-center">
        <h1 className="text-3xl font-bold">Usuários</h1>
        <Button onClick={() => setIsDialogOpen(true)}>
          <Plus className="mr-2 h-4 w-4" />
          Novo Usuário
        </Button>
      </div>

      <Card>
        <CardHeader>
          <div className="flex gap-4">
            <div className="relative flex-1">
              <Search className="absolute left-3 top-3 h-4 w-4 text-muted-foreground" />
              <Input
                placeholder="Buscar usuários..."
                value={searchTerm}
                onChange={(e) => setSearchTerm(e.target.value)}
                className="pl-9"
              />
            </div>
            <select
              value={selectedRole}
              onChange={(e) => setSelectedRole(e.target.value)}
              className="px-3 py-2 border rounded-md"
            >
              <option value="">Todas as roles</option>
              <option value="admin">Admin</option>
              <option value="moderator">Moderador</option>
              <option value="collaborator">Colaborador</option>
              <option value="citizen">Cidadão</option>
            </select>
          </div>
        </CardHeader>
        <CardContent>
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>Nome</TableHead>
                <TableHead>Email</TableHead>
                <TableHead>Role</TableHead>
                <TableHead>Status</TableHead>
                <TableHead>Criado em</TableHead>
                <TableHead className="w-[100px]">Ações</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {isLoading ? (
                <TableRow>
                  <TableCell colSpan={6} className="text-center">
                    Carregando...
                  </TableCell>
                </TableRow>
              ) : filteredUsers.length === 0 ? (
                <TableRow>
                  <TableCell colSpan={6} className="text-center">
                    Nenhum usuário encontrado
                  </TableCell>
                </TableRow>
              ) : (
                filteredUsers.map((user: any) => (
                  <TableRow key={user.id}>
                    <TableCell className="font-medium">{user.name}</TableCell>
                    <TableCell>{user.email}</TableCell>
                    <TableCell>
                      <Badge variant={getRoleColor(user.role)}>
                        {user.role}
                      </Badge>
                    </TableCell>
                    <TableCell>
                      <Badge variant={user.status === 'active' ? 'default' : 'secondary'}>
                        {user.status}
                      </Badge>
                    </TableCell>
                    <TableCell>
                      {new Date(user.created_at).toLocaleDateString('pt-BR')}
                    </TableCell>
                    <TableCell>
                      <Button
                        variant="ghost"
                        size="sm"
                        onClick={() => {
                          setEditingUser(user)
                          setIsDialogOpen(true)
                        }}
                      >
                        <MoreHorizontal className="h-4 w-4" />
                      </Button>
                    </TableCell>
                  </TableRow>
                ))
              )}
            </TableBody>
          </Table>
        </CardContent>
      </Card>

      <UserDialog
        open={isDialogOpen}
        onOpenChange={setIsDialogOpen}
        user={editingUser}
        onSubmit={editingUser ? handleUpdateUser : handleCreateUser}
        isLoading={createUserMutation.isPending || updateUserMutation.isPending}
      />
    </div>
  )
}
```

#### 10. Faction Management Interface
```typescript
// admin/src/app/factions/page.tsx
'use client'

import { useState } from 'react'
import { useFactions, useCreateFaction, useUpdateFaction } from '@/lib/queries'
import { Button } from '@/components/ui/button'
import { Card, CardHeader, CardTitle, CardContent } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { FactionDialog } from '@/components/factions/faction-dialog'
import { Plus, Palette } from 'lucide-react'

export default function FactionsPage() {
  const [isDialogOpen, setIsDialogOpen] = useState(false)
  const [editingFaction, setEditingFaction] = useState<any>(null)

  const { data: factionsData, isLoading } = useFactions()
  const createFactionMutation = useCreateFaction()
  const updateFactionMutation = useUpdateFaction()

  const handleCreateFaction = (factionData: any) => {
    createFactionMutation.mutate(factionData, {
      onSuccess: () => {
        setIsDialogOpen(false)
      }
    })
  }

  const handleUpdateFaction = (factionData: any) => {
    if (editingFaction) {
      updateFactionMutation.mutate(
        { id: editingFaction.id, data: factionData },
        {
          onSuccess: () => {
            setIsDialogOpen(false)
            setEditingFaction(null)
          }
        }
      )
    }
  }

  return (
    <div className="space-y-6">
      <div className="flex justify-between items-center">
        <h1 className="text-3xl font-bold">Facções</h1>
        <Button onClick={() => setIsDialogOpen(true)}>
          <Plus className="mr-2 h-4 w-4" />
          Nova Facção
        </Button>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        {isLoading ? (
          <div className="col-span-full text-center">Carregando...</div>
        ) : (
          factionsData?.data?.map((faction: any) => (
            <Card key={faction.id} className="cursor-pointer hover:shadow-lg transition-shadow">
              <CardHeader>
                <div className="flex items-center justify-between">
                  <div className="flex items-center gap-3">
                    <div 
                      className="w-6 h-6 rounded-full"
                      style={{ backgroundColor: faction.color_hex }}
                    />
                    <div>
                      <CardTitle className="text-lg">{faction.name}</CardTitle>
                      <p className="text-sm text-muted-foreground">{faction.acronym}</p>
                    </div>
                  </div>
                  <Badge variant={faction.active ? 'default' : 'secondary'}>
                    {faction.active ? 'Ativa' : 'Inativa'}
                  </Badge>
                </div>
              </CardHeader>
              <CardContent>
                <div className="space-y-2">
                  <div className="flex justify-between text-sm">
                    <span>Prioridade:</span>
                    <span>{faction.display_priority}</span>
                  </div>
                  <div className="flex justify-between text-sm">
                    <span>Cor:</span>
                    <div className="flex items-center gap-2">
                      <div 
                        className="w-4 h-4 rounded border"
                        style={{ backgroundColor: faction.color_hex }}
                      />
                      <span className="font-mono text-xs">{faction.color_hex}</span>
                    </div>
                  </div>
                </div>
                <Button
                  variant="outline"
                  size="sm"
                  className="w-full mt-4"
                  onClick={() => {
                    setEditingFaction(faction)
                    setIsDialogOpen(true)
                  }}
                >
                  <Palette className="mr-2 h-4 w-4" />
                  Editar
                </Button>
              </CardContent>
            </Card>
          ))
        )}
      </div>

      <FactionDialog
        open={isDialogOpen}
        onOpenChange={setIsDialogOpen}
        faction={editingFaction}
        onSubmit={editingFaction ? handleUpdateFaction : handleCreateFaction}
        isLoading={createFactionMutation.isPending || updateFactionMutation.isPending}
      />
    </div>
  )
}
```

**Tarefas:**
- [ ] Implementar interface de usuários com CRUD
- [ ] Implementar interface de facções com color picker
- [ ] Criar componentes de tabela reutilizáveis
- [ ] Implementar filtros e busca
- [ ] Adicionar paginação
- [ ] Implementar modais de confirmação

### Day 13-17: Map Editor

#### 11. Map Editor Component
```typescript
// admin/src/components/maps/map-editor.tsx
'use client'

import { useEffect, useRef, useState } from 'react'
import mapboxgl from 'mapbox-gl'
import MapboxDraw from '@mapbox/mapbox-gl-draw'
import { useGeometries, useCreateGeometry, useUpdateGeometry } from '@/lib/queries'
import { useFactions } from '@/lib/queries'
import { Button } from '@/components/ui/button'
import { Card, CardContent } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { Trash2, Edit, Save, X } from 'lucide-react'

interface MapEditorProps {
  regionId?: string
}

export function MapEditor({ regionId }: MapEditorProps) {
  const mapContainer = useRef<HTMLDivElement>(null)
  const map = useRef<mapboxgl.Map | null>(null)
  const draw = useRef<MapboxDraw | null>(null)
  
  const [selectedGeometry, setSelectedGeometry] = useState<any>(null)
  const [editingGeometry, setEditingGeometry] = useState<any>(null)
  const [isDrawing, setIsDrawing] = useState(false)

  const { data: geometriesData } = useGeometries({ region_id: regionId })
  const { data: factionsData } = useFactions({ active: true })
  const createGeometryMutation = useCreateGeometry()
  const updateGeometryMutation = useUpdateGeometry()

  // Initialize map
  useEffect(() => {
    if (!mapContainer.current) return

    mapboxgl.accessToken = process.env.NEXT_PUBLIC_MAPBOX_TOKEN!

    map.current = new mapboxgl.Map({
      container: mapContainer.current,
      style: 'mapbox://styles/mapbox/streets-v12',
      center: [-38.5267, -3.7172], // Fortaleza
      zoom: 10,
    })

    // Initialize drawing tools
    draw.current = new MapboxDraw({
      displayControlsDefault: false,
      controls: {
        polygon: true,
        line_string: true,
        trash: true,
      },
      defaultMode: 'simple_select',
    })

    map.current.addControl(draw.current)

    // Event listeners
    map.current.on('draw.create', handleGeometryCreate)
    map.current.on('draw.update', handleGeometryUpdate)
    map.current.on('draw.delete', handleGeometryDelete)
    map.current.on('draw.selectionchange', handleSelectionChange)

    return () => {
      if (map.current) {
        map.current.remove()
      }
    }
  }, [])

  // Load geometries on map
  useEffect(() => {
    if (!map.current || !geometriesData?.features) return

    // Remove existing geometries
    if (draw.current) {
      draw.current.deleteAll()
    }

    // Add geometries to map
    geometriesData.features.forEach((feature: any) => {
      if (draw.current) {
        draw.current.add(feature)
      }
    })
  }, [geometriesData])

  const handleGeometryCreate = (e: any) => {
    const feature = e.features[0]
    setSelectedGeometry(feature)
    setIsDrawing(false)
  }

  const handleGeometryUpdate = (e: any) => {
    const feature = e.features[0]
    setEditingGeometry(feature)
  }

  const handleGeometryDelete = (e: any) => {
    // Handle deletion if needed
  }

  const handleSelectionChange = (e: any) => {
    const features = e.features
    setSelectedGeometry(features.length > 0 ? features[0] : null)
  }

  const saveNewGeometry = async (factionId: string, riskLevel: string) => {
    if (!selectedGeometry) return

    const geometryData = {
      geojson: {
        type: 'Feature',
        geometry: selectedGeometry.geometry,
        properties: {
          faction_id: factionId,
          risk_level: riskLevel,
        },
      },
      region_id: regionId,
      geometry_type: selectedGeometry.geometry.type.toLowerCase().includes('polygon') 
        ? 'polygon' 
        : 'polyline',
      faction_id: factionId,
      risk_level: riskLevel,
      status: 'draft',
    }

    createGeometryMutation.mutate(geometryData, {
      onSuccess: () => {
        setSelectedGeometry(null)
        if (draw.current) {
          draw.current.changeMode('simple_select')
        }
      }
    })
  }

  const updateExistingGeometry = async (geometryId: string, updates: any) => {
    updateGeometryMutation.mutate(
      { id: geometryId, data: updates },
      {
        onSuccess: () => {
          setEditingGeometry(null)
        }
      }
    )
  }

  const getFactionColor = (factionId: string) => {
    const faction = factionsData?.data?.find((f: any) => f.id === factionId)
    return faction?.color_hex || '#808080'
  }

  return (
    <div className="flex h-[600px] gap-4">
      {/* Map Container */}
      <div className="flex-1 relative">
        <div ref={mapContainer} className="w-full h-full rounded-lg" />
        
        {/* Drawing Controls */}
        <div className="absolute top-4 left-4 space-y-2">
          <Button
            size="sm"
            variant={isDrawing ? 'default' : 'outline'}
            onClick={() => {
              if (draw.current) {
                draw.current.changeMode('draw_polygon')
                setIsDrawing(true)
              }
            }}
          >
            Desenhar Polígono
          </Button>
          <Button
            size="sm"
            variant="outline"
            onClick={() => {
              if (draw.current) {
                draw.current.changeMode('draw_line_string')
                setIsDrawing(true)
              }
            }}
          >
            Desenhar Linha
          </Button>
        </div>
      </div>

      {/* Sidebar */}
      <div className="w-80 space-y-4">
        {/* Selected Geometry Panel */}
        {selectedGeometry && (
          <Card>
            <CardContent className="p-4">
              <h3 className="font-medium mb-3">Nova Geometria</h3>
              <div className="space-y-3">
                <div>
                  <label className="text-sm font-medium">Facção</label>
                  <select className="w-full mt-1 px-3 py-2 border rounded-md">
                    <option value="">Selecione uma facção</option>
                    {factionsData?.data?.map((faction: any) => (
                      <option key={faction.id} value={faction.id}>
                        {faction.name} ({faction.acronym})
                      </option>
                    ))}
                  </select>
                </div>
                
                <div>
                  <label className="text-sm font-medium">Nível de Risco</label>
                  <select className="w-full mt-1 px-3 py-2 border rounded-md">
                    <option value="low">Baixo</option>
                    <option value="medium">Médio</option>
                    <option value="high">Alto</option>
                  </select>
                </div>

                <div className="flex gap-2">
                  <Button
                    size="sm"
                    onClick={() => saveNewGeometry('faction-id', 'medium')}
                    disabled={createGeometryMutation.isPending}
                  >
                    <Save className="mr-2 h-4 w-4" />
                    Salvar
                  </Button>
                  <Button
                    size="sm"
                    variant="outline"
                    onClick={() => {
                      setSelectedGeometry(null)
                      if (draw.current) {
                        draw.current.trash()
                      }
                    }}
                  >
                    <X className="mr-2 h-4 w-4" />
                    Cancelar
                  </Button>
                </div>
              </div>
            </CardContent>
          </Card>
        )}

        {/* Existing Geometries List */}
        <Card>
          <CardContent className="p-4">
            <h3 className="font-medium mb-3">Geometrias na Região</h3>
            <div className="space-y-2 max-h-64 overflow-y-auto">
              {geometriesData?.features?.map((feature: any) => (
                <div
                  key={feature.id}
                  className="flex items-center justify-between p-2 border rounded"
                >
                  <div className="flex items-center gap-2">
                    <div
                      className="w-3 h-3 rounded"
                      style={{ 
                        backgroundColor: getFactionColor(feature.properties.faction_id) 
                      }}
                    />
                    <span className="text-sm">
                      {feature.properties.faction_name || 'Sem nome'}
                    </span>
                    <Badge variant="outline" className="text-xs">
                      {feature.properties.risk_level}
                    </Badge>
                  </div>
                  <div className="flex gap-1">
                    <Button size="sm" variant="ghost">
                      <Edit className="h-3 w-3" />
                    </Button>
                    <Button size="sm" variant="ghost">
                      <Trash2 className="h-3 w-3" />
                    </Button>
                  </div>
                </div>
              ))}
            </div>
          </CardContent>
        </Card>

        {/* Layer Control */}
        <Card>
          <CardContent className="p-4">
            <h3 className="font-medium mb-3">Camadas</h3>
            <div className="space-y-2">
              {factionsData?.data?.map((faction: any) => (
                <label key={faction.id} className="flex items-center gap-2">
                  <input
                    type="checkbox"
                    defaultChecked
                    className="rounded"
                  />
                  <div
                    className="w-3 h-3 rounded"
                    style={{ backgroundColor: faction.color_hex }}
                  />
                  <span className="text-sm">{faction.name}</span>
                </label>
              ))}
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  )
}
```

**Tarefas:**
- [ ] Implementar MapEditor com Mapbox GL JS
- [ ] Integrar Mapbox GL Draw para desenho
- [ ] Implementar controles de desenho
- [ ] Criar sidebar com propriedades
- [ ] Implementar salvamento de geometrias
- [ ] Adicionar controle de camadas por facção
- [ ] Implementar edição de geometrias existentes

### Day 18-20: Dashboard e Finalização

#### 12. Dashboard com Métricas
```typescript
// admin/src/app/dashboard/page.tsx
'use client'

import { Card, CardHeader, CardTitle, CardContent } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { BarChart, Bar, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer } from 'recharts'
import { Users, Map, Flag, Activity } from 'lucide-react'

export default function DashboardPage() {
  // Mock data - substitute with real API calls
  const stats = {
    totalUsers: 1245,
    totalGeometries: 89,
    totalFactions: 8,
    totalReports: 23,
  }

  const activityData = [
    { name: 'Jan', geometries: 12, users: 45 },
    { name: 'Fev', geometries: 19, users: 52 },
    { name: 'Mar', geometries: 15, users: 61 },
    { name: 'Abr', geometries: 23, users: 75 },
    { name: 'Mai', geometries: 18, users: 89 },
    { name: 'Jun', geometries: 31, users: 95 },
  ]

  const recentActivities = [
    {
      id: 1,
      user: 'João Silva',
      action: 'Criou geometria',
      target: 'Messejana, Fortaleza',
      time: '2 horas atrás',
      type: 'create',
    },
    {
      id: 2,
      user: 'Maria Santos',
      action: 'Atualizou facção',
      target: 'TDN',
      time: '4 horas atrás',
      type: 'update',
    },
    {
      id: 3,
      user: 'Carlos Lima',
      action: 'Aprovou geometria',
      target: 'Centro, Fortaleza',
      time: '6 horas atrás',
      type: 'approve',
    },
  ]

  return (
    <div className="space-y-6">
      <h1 className="text-3xl font-bold">Dashboard</h1>

      {/* Stats Cards */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Total de Usuários</CardTitle>
            <Users className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{stats.totalUsers}</div>
            <p className="text-xs text-muted-foreground">
              +12% em relação ao mês passado
            </p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Geometrias</CardTitle>
            <Map className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{stats.totalGeometries}</div>
            <p className="text-xs text-muted-foreground">
              +8% em relação ao mês passado
            </p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Facções Ativas</CardTitle>
            <Flag className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{stats.totalFactions}</div>
            <p className="text-xs text-muted-foreground">
              Estável
            </p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Relatórios Pendentes</CardTitle>
            <Activity className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{stats.totalReports}</div>
            <p className="text-xs text-muted-foreground">
              -3% em relação ao mês passado
            </p>
          </CardContent>
        </Card>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        {/* Activity Chart */}
        <Card>
          <CardHeader>
            <CardTitle>Atividade dos Últimos Meses</CardTitle>
          </CardHeader>
          <CardContent>
            <ResponsiveContainer width="100%" height={300}>
              <BarChart data={activityData}>
                <CartesianGrid strokeDasharray="3 3" />
                <XAxis dataKey="name" />
                <YAxis />
                <Tooltip />
                <Bar dataKey="geometries" fill="#8884d8" name="Geometrias" />
                <Bar dataKey="users" fill="#82ca9d" name="Usuários" />
              </BarChart>
            </ResponsiveContainer>
          </CardContent>
        </Card>

        {/* Recent Activity */}
        <Card>
          <CardHeader>
            <CardTitle>Atividade Recente</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="space-y-4">
              {recentActivities.map((activity) => (
                <div key={activity.id} className="flex items-start space-x-3">
                  <div className="flex-shrink-0">
                    <div className={`w-2 h-2 rounded-full mt-2 ${
                      activity.type === 'create' ? 'bg-green-500' :
                      activity.type === 'update' ? 'bg-blue-500' :
                      'bg-yellow-500'
                    }`} />
                  </div>
                  <div className="flex-1 min-w-0">
                    <div className="text-sm">
                      <span className="font-medium">{activity.user}</span>
                      {' '}{activity.action}{' '}
                      <span className="font-medium">{activity.target}</span>
                    </div>
                    <div className="text-xs text-muted-foreground">
                      {activity.time}
                    </div>
                  </div>
                </div>
              ))}
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  )
}
```

**Tarefas:**
- [ ] Implementar dashboard com cards de estatísticas
- [ ] Integrar gráficos com Recharts
- [ ] Criar feed de atividades recentes
- [ ] Implementar métricas em tempo real
- [ ] Adicionar filtros de período
- [ ] Criar alertas para métricas críticas

## 🎯 Critérios de Aceite Fase 2

### Interface Administrativa
- [ ] ✅ Login funcional integrado com backend
- [ ] ✅ Layout responsivo com sidebar
- [ ] ✅ CRUD completo de usuários
- [ ] ✅ CRUD completo de facções com color picker
- [ ] ✅ Editor de mapas com Mapbox GL Draw
- [ ] ✅ Dashboard com métricas e gráficos

### Funcionalidades do Mapa
- [ ] ✅ Desenho de polígonos e polilinhas
- [ ] ✅ Edição de geometrias existentes
- [ ] ✅ Controle de camadas por facção
- [ ] ✅ Atribuição de facção e risco
- [ ] ✅ Salvamento no backend
- [ ] ✅ Visualização por cores

### Qualidade e UX
- [ ] ✅ Design responsivo (desktop/tablet)
- [ ] ✅ Loading states e error handling
- [ ] ✅ Validação de formulários
- [ ] ✅ Feedback visual para ações
- [ ] ✅ Navegação intuitiva
- [ ] ✅ Performance otimizada

### Integração
- [ ] ✅ API client configurado
- [ ] ✅ Authentication com NextAuth
- [ ] ✅ State management com React Query
- [ ] ✅ Type safety com TypeScript
- [ ] ✅ Error boundary implementado

## 🚀 Comandos para Testar

```bash
# Desenvolvimento
cd admin
npm run dev

# Build
npm run build

# Testes
npm test
npm run test:e2e

# Linting
npm run lint
npm run type-check

# Acessar aplicação
open http://localhost:3000
```

## 📋 Próximos Passos

Após completar a Fase 2:
1. **Demo com stakeholders** usando editor de mapas
2. **User testing** com usuários finais
3. **Performance optimization**
4. **Mobile responsiveness** final
5. **Preparação para Fase 3** (Mobile App)

**🎯 Com a Fase 2 completa, você terá uma interface administrativa poderosa para gerenciar todo o sistema de mapas!**