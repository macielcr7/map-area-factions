'use client'

import { useState } from 'react'
import { Card, CardHeader, CardTitle, CardContent } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Badge } from '@/components/ui/badge'
import { 
  Table, 
  TableBody, 
  TableCell, 
  TableHead, 
  TableHeader, 
  TableRow 
} from '@/components/ui/table'
import { 
  Activity, 
  Search, 
  Calendar,
  User,
  Database,
  Eye
} from 'lucide-react'

// Mock audit data - replace with real API calls
type AuditLogAction = 'create' | 'update' | 'delete' | 'read'
type AuditLogEntity = 'user' | 'faction' | 'geometry' | 'report'

interface AuditLog {
  id: string
  user: { name: string; email: string }
  entity: AuditLogEntity
  entity_id: string
  action: AuditLogAction
  changes: string
  ip_address: string
  user_agent: string
  created_at: string
}

const mockAuditLogs: AuditLog[] = [
  {
    id: '1',
    user: { name: 'João Silva', email: 'joao@admin.com' },
    entity: 'faction',
    entity_id: 'uuid-1', 
    action: 'create',
    changes: JSON.stringify({ name: 'PCC', color: '#ff0000' }),
    ip_address: '192.168.1.100',
    user_agent: 'Mozilla/5.0...',
    created_at: '2024-01-20T10:30:00Z'
  },
  {
    id: '2',
    user: { name: 'Maria Santos', email: 'maria@moderator.com' },
    entity: 'user',
    entity_id: 'uuid-2',
    action: 'update', 
    changes: JSON.stringify({ role: 'moderator', active: true }),
    ip_address: '192.168.1.101',
    user_agent: 'Mozilla/5.0...',
    created_at: '2024-01-20T11:15:00Z'
  },
  {
    id: '3',
    user: { name: 'Pedro Oliveira', email: 'pedro@collaborator.com' },
    entity: 'geometry',
    entity_id: 'uuid-3',
    action: 'delete',
    changes: JSON.stringify({ geometry_type: 'polygon', faction_id: 'uuid-1' }),
    ip_address: '192.168.1.102', 
    user_agent: 'Mozilla/5.0...',
    created_at: '2024-01-20T12:45:00Z'
  },
  {
    id: '4',
    user: { name: 'Ana Costa', email: 'ana@citizen.com' },
    entity: 'report',
    entity_id: 'uuid-4',
    action: 'create',
    changes: JSON.stringify({ type: 'inaccuracy', status: 'pending' }),
    ip_address: '192.168.1.103',
    user_agent: 'Mozilla/5.0...',
    created_at: '2024-01-20T14:20:00Z'
  }
]

const actionColors: Record<AuditLogAction, string> = {
  create: 'bg-green-500',
  update: 'bg-blue-500',
  delete: 'bg-red-500',
  read: 'bg-gray-500'
}

const actionLabels: Record<AuditLogAction, string> = {
  create: 'Criação',
  update: 'Atualização', 
  delete: 'Exclusão',
  read: 'Leitura'
}

const entityLabels: Record<AuditLogEntity, string> = {
  user: 'Usuário',
  faction: 'Facção',
  geometry: 'Geometria',
  report: 'Relatório'
}

export default function AuditPage() {
  const [logs, setLogs] = useState<AuditLog[]>(mockAuditLogs)
  const [searchTerm, setSearchTerm] = useState('')
  const [entityFilter, setEntityFilter] = useState<AuditLogEntity | ''>('')
  const [actionFilter, setActionFilter] = useState<AuditLogAction | ''>('')
  const [selectedLog, setSelectedLog] = useState<AuditLog | null>(null)
  const [showDetails, setShowDetails] = useState(false)

  const filteredLogs = logs.filter(log => {
    const matchesSearch = log.user.name.toLowerCase().includes(searchTerm.toLowerCase()) ||
                         log.user.email.toLowerCase().includes(searchTerm.toLowerCase()) ||
                         log.entity.toLowerCase().includes(searchTerm.toLowerCase())
    const matchesEntity = !entityFilter || log.entity === entityFilter
    const matchesAction = !actionFilter || log.action === actionFilter
    return matchesSearch && matchesEntity && matchesAction
  })

  const viewLogDetails = (log: AuditLog) => {
    setSelectedLog(log)
    setShowDetails(true)
  }

  const formatChanges = (changesString: string) => {
    try {
      const changes = JSON.parse(changesString)
      return Object.entries(changes).map(([key, value]) => (
        <div key={key} className="mb-1">
          <span className="font-medium">{key}:</span> {JSON.stringify(value)}
        </div>
      ))
    } catch {
      return changesString
    }
  }

  // Mock statistics
  const stats = {
    total_logs: logs.length,
    today_logs: logs.filter(log => {
      const today = new Date().toDateString()
      return new Date(log.created_at).toDateString() === today
    }).length,
    by_action: {
      create: logs.filter(log => log.action === 'create').length,
      update: logs.filter(log => log.action === 'update').length,
      delete: logs.filter(log => log.action === 'delete').length
    }
  }

  return (
    <div className="space-y-6">
      <div className="flex justify-between items-center">
        <h1 className="text-3xl font-bold">Log de Auditoria</h1>
        <Button variant="outline">
          <Calendar className="mr-2 h-4 w-4" />
          Exportar Relatório
        </Button>
      </div>

      {/* Statistics Cards */}
      <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
        <Card>
          <CardContent className="pt-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-gray-600">Total de Logs</p>
                <p className="text-2xl font-bold">{stats.total_logs}</p>
              </div>
              <Activity className="h-8 w-8 text-blue-600" />
            </div>
          </CardContent>
        </Card>
        
        <Card>
          <CardContent className="pt-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-gray-600">Hoje</p>
                <p className="text-2xl font-bold">{stats.today_logs}</p>
              </div>
              <Calendar className="h-8 w-8 text-green-600" />
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardContent className="pt-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-gray-600">Criações</p>
                <p className="text-2xl font-bold">{stats.by_action.create}</p>
              </div>
              <div className="h-8 w-8 bg-green-500 rounded" />
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardContent className="pt-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-gray-600">Exclusões</p>
                <p className="text-2xl font-bold text-red-600">{stats.by_action.delete}</p>
              </div>
              <div className="h-8 w-8 bg-red-500 rounded" />
            </div>
          </CardContent>
        </Card>
      </div>

      {/* Filters */}
      <Card>
        <CardContent className="pt-6">
          <div className="flex gap-4">
            <div className="flex-1">
              <div className="relative">
                <Search className="absolute left-3 top-3 h-4 w-4 text-gray-400" />
                <Input
                  placeholder="Buscar por usuário, email ou entidade..."
                  value={searchTerm}
                  onChange={(e) => setSearchTerm(e.target.value)}
                  className="pl-10"
                />
              </div>
            </div>
            <div className="w-48">
              <select
                value={entityFilter}
                onChange={(e) => setEntityFilter(e.target.value as AuditLogEntity | '')}
                className="w-full p-2 border rounded-md"
              >
                <option value="">Todas as entidades</option>
                <option value="user">Usuário</option>
                <option value="faction">Facção</option>
                <option value="geometry">Geometria</option>
                <option value="report">Relatório</option>
              </select>
            </div>
            <div className="w-48">
              <select
                value={actionFilter}
                onChange={(e) => setActionFilter(e.target.value as AuditLogAction | '')}
                className="w-full p-2 border rounded-md"
              >
                <option value="">Todas as ações</option>
                <option value="create">Criação</option>
                <option value="update">Atualização</option>
                <option value="delete">Exclusão</option>
              </select>
            </div>
          </div>
        </CardContent>
      </Card>

      {/* Audit Logs Table */}
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center">
            <Activity className="mr-2 h-5 w-5" />
            Registros de Auditoria ({filteredLogs.length})
          </CardTitle>
        </CardHeader>
        <CardContent>
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>Data/Hora</TableHead>
                <TableHead>Usuário</TableHead>
                <TableHead>Entidade</TableHead>
                <TableHead>Ação</TableHead>
                <TableHead>IP</TableHead>
                <TableHead>Detalhes</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {filteredLogs.map((log) => (
                <TableRow key={log.id}>
                  <TableCell>
                    {new Date(log.created_at).toLocaleString('pt-BR')}
                  </TableCell>
                  <TableCell>
                    <div>
                      <p className="font-medium">{log.user.name}</p>
                      <p className="text-sm text-gray-500">{log.user.email}</p>
                    </div>
                  </TableCell>
                  <TableCell>
                    <Badge variant="outline">
                      {entityLabels[log.entity]}
                    </Badge>
                  </TableCell>
                  <TableCell>
                    <Badge className={`${actionColors[log.action]} text-white`}>
                      {actionLabels[log.action]}
                    </Badge>
                  </TableCell>
                  <TableCell className="font-mono text-sm">
                    {log.ip_address}
                  </TableCell>
                  <TableCell>
                    <Button
                      variant="outline"
                      size="sm"
                      onClick={() => viewLogDetails(log)}
                    >
                      <Eye className="h-4 w-4" />
                    </Button>
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </CardContent>
      </Card>

      {/* Log Details Modal */}
      {showDetails && selectedLog && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
          <div className="bg-white rounded-lg p-6 max-w-2xl w-full mx-4">
            <div className="flex justify-between items-center mb-4">
              <h3 className="text-lg font-semibold">Detalhes do Log de Auditoria</h3>
              <Button variant="outline" onClick={() => setShowDetails(false)}>
                ✕
              </Button>
            </div>
            
            <div className="space-y-4">
              <div className="grid grid-cols-2 gap-4">
                <div>
                  <p className="font-medium text-gray-600">Usuário:</p>
                  <p>{selectedLog.user.name} ({selectedLog.user.email})</p>
                </div>
                <div>
                  <p className="font-medium text-gray-600">Data/Hora:</p>
                  <p>{new Date(selectedLog.created_at).toLocaleString('pt-BR')}</p>
                </div>
                <div>
                  <p className="font-medium text-gray-600">Entidade:</p>
                  <p>{entityLabels[selectedLog.entity]} (ID: {selectedLog.entity_id})</p>
                </div>
                <div>
                  <p className="font-medium text-gray-600">Ação:</p>
                  <Badge className={`${actionColors[selectedLog.action]} text-white`}>
                    {actionLabels[selectedLog.action]}
                  </Badge>
                </div>
                <div>
                  <p className="font-medium text-gray-600">Endereço IP:</p>
                  <p className="font-mono">{selectedLog.ip_address}</p>
                </div>
              </div>
              
              <div>
                <p className="font-medium text-gray-600 mb-2">Alterações:</p>
                <div className="bg-gray-100 p-3 rounded border font-mono text-sm">
                  {formatChanges(selectedLog.changes)}
                </div>
              </div>

              <div>
                <p className="font-medium text-gray-600 mb-2">User Agent:</p>
                <p className="text-sm bg-gray-100 p-2 rounded break-all">
                  {selectedLog.user_agent}
                </p>
              </div>
            </div>
          </div>
        </div>
      )}
    </div>
  )
}
