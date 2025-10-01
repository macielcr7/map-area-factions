'use client'

import { useMemo, useState } from 'react'
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
  TableRow,
} from '@/components/ui/table'
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import {
  Activity,
  Search,
  Calendar,
  User,
  Database,
  Eye,
  Loader2,
} from 'lucide-react'
import { useAuditLogs } from '@/lib/queries'
import { useDebounce } from '@/lib/use-debounce'
import { format } from 'date-fns'
import ptBR from 'date-fns/locale/pt-BR'

interface AuditLog {
  id: string
  user: {
    id?: string
    name: string
    email: string
  } | null
  entity: string
  entity_id: string | null
  action: string
  changes: string | null
  ip_address: string | null
  user_agent: string | null
  created_at: string
}

const ACTION_LABELS: Record<string, string> = {
  create: 'Criação',
  update: 'Atualização',
  delete: 'Exclusão',
  read: 'Leitura',
}

const ENTITY_LABELS: Record<string, string> = {
  user: 'Usuário',
  faction: 'Facção',
  geometry: 'Geometria',
  report: 'Relatório',
}

const PAGE_SIZE = 20

export default function AuditPage() {
  const [page, setPage] = useState(1)
  const [searchTerm, setSearchTerm] = useState('')
  const debouncedSearch = useDebounce(searchTerm, 300)
  const [entityFilter, setEntityFilter] = useState('')
  const [actionFilter, setActionFilter] = useState('')
  const [selectedLog, setSelectedLog] = useState<AuditLog | null>(null)
  const [showDetails, setShowDetails] = useState(false)

  const { data, isLoading, isFetching } = useAuditLogs({
    page,
    limit: PAGE_SIZE,
    entity: entityFilter || undefined,
    action: actionFilter || undefined,
  })

  const logs = useMemo(() => (data?.logs ?? []) as AuditLog[], [data?.logs])
  const pagination = data?.pagination as
    | {
        page: number
        limit: number
        total: number
        pages: number
        has_next: boolean
        has_prev: boolean
      }
    | undefined

  const filteredLogs = useMemo(() => {
    if (!debouncedSearch) {
      return logs
    }

    const term = debouncedSearch.toLowerCase()

    return logs.filter((log) => {
      const userName = log.user?.name?.toLowerCase() ?? 'anônimo'
      const userEmail = log.user?.email?.toLowerCase() ?? ''
      const entity = log.entity.toLowerCase()
      const action = log.action.toLowerCase()
      return (
        userName.includes(term) ||
        userEmail.includes(term) ||
        entity.includes(term) ||
        action.includes(term)
      )
    })
  }, [logs, debouncedSearch])

  const stats = useMemo(() => {
    const totalLogs = logs.length

    const byAction = logs.reduce<Record<string, number>>((acc, log) => {
      const key = log.action
      acc[key] = (acc[key] ?? 0) + 1
      return acc
    }, {})

    const todayLogs = logs.filter((log) => {
      const today = new Date()
      const logDate = new Date(log.created_at)
      return (
        logDate.getDate() === today.getDate() &&
        logDate.getMonth() === today.getMonth() &&
        logDate.getFullYear() === today.getFullYear()
      )
    }).length

    return {
      total_logs: totalLogs,
      today_logs: todayLogs,
      by_action: {
        create: byAction.create ?? 0,
        update: byAction.update ?? 0,
        delete: byAction.delete ?? 0,
      },
    }
  }, [logs])

  const isLogsLoading = isLoading || isFetching

  const viewLogDetails = (log: AuditLog) => {
    setSelectedLog(log)
    setShowDetails(true)
  }

  const formatChanges = (changesString: string | null) => {
    if (!changesString) {
      return 'Sem detalhes registrados.'
    }

    try {
      const changes = JSON.parse(changesString) as Record<string, unknown>
      return Object.entries(changes).map(([key, value]) => (
        <div key={key} className="mb-1">
          <span className="font-medium">{key}:</span> {JSON.stringify(value)}
        </div>
      ))
    } catch (error) {
      return changesString
    }
  }

  const resetFilters = () => {
    setSearchTerm('')
    setEntityFilter('')
    setActionFilter('')
    setPage(1)
  }

  return (
    <div className="space-y-6">
      <div className="flex justify-between items-center">
        <h1 className="text-3xl font-bold">Log de Auditoria</h1>
        <Button variant="outline" onClick={resetFilters}>
          <Calendar className="mr-2 h-4 w-4" />
          Resetar Filtros
        </Button>
      </div>

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
                <p className="text-2xl font-bold text-green-600">{stats.today_logs}</p>
              </div>
              <User className="h-8 w-8 text-green-600" />
            </div>
          </CardContent>
        </Card>
        <Card>
          <CardContent className="pt-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-gray-600">Criações</p>
                <p className="text-2xl font-bold text-blue-600">{stats.by_action.create ?? 0}</p>
              </div>
              <Database className="h-8 w-8 text-blue-600" />
            </div>
          </CardContent>
        </Card>
        <Card>
          <CardContent className="pt-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-gray-600">Exclusões</p>
                <p className="text-2xl font-bold text-red-600">{stats.by_action.delete ?? 0}</p>
              </div>
              <Database className="h-8 w-8 text-red-600" />
            </div>
          </CardContent>
        </Card>
      </div>

      <Card>
        <CardContent className="pt-6">
          <div className="flex flex-col gap-4 md:flex-row">
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
            <div className="w-full md:w-48">
              <select
                value={entityFilter}
                onChange={(e) => {
                  setEntityFilter(e.target.value)
                  setPage(1)
                }}
                className="w-full p-2 border rounded-md"
              >
                <option value="">Todas as entidades</option>
                {Object.entries(ENTITY_LABELS).map(([value, label]) => (
                  <option key={value} value={value}>
                    {label}
                  </option>
                ))}
              </select>
            </div>
            <div className="w-full md:w-48">
              <select
                value={actionFilter}
                onChange={(e) => {
                  setActionFilter(e.target.value)
                  setPage(1)
                }}
                className="w-full p-2 border rounded-md"
              >
                <option value="">Todas as ações</option>
                {Object.entries(ACTION_LABELS).map(([value, label]) => (
                  <option key={value} value={value}>
                    {label}
                  </option>
                ))}
              </select>
            </div>
          </div>
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <CardTitle className="flex items-center">
            <Activity className="mr-2 h-5 w-5" />
            Registros de Auditoria ({pagination?.total ?? filteredLogs.length})
          </CardTitle>
        </CardHeader>
        <CardContent>
          <div className="rounded-md border">
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
                {isLogsLoading ? (
                  <TableRow>
                    <TableCell colSpan={6} className="py-8 text-center">
                      <div className="flex items-center justify-center gap-2 text-sm text-muted-foreground">
                        <Loader2 className="h-4 w-4 animate-spin" />
                        Carregando logs...
                      </div>
                    </TableCell>
                  </TableRow>
                ) : filteredLogs.length === 0 ? (
                  <TableRow>
                    <TableCell colSpan={6} className="py-8 text-center text-sm text-muted-foreground">
                      Nenhum registro encontrado.
                    </TableCell>
                  </TableRow>
                ) : (
                  filteredLogs.map((log) => (
                    <TableRow key={log.id}>
                      <TableCell className="whitespace-nowrap">
                        {format(new Date(log.created_at), 'dd/MM/yyyy HH:mm', { locale: ptBR })}
                      </TableCell>
                      <TableCell>
                        {log.user ? (
                          <div>
                            <p className="font-medium">{log.user.name}</p>
                            <p className="text-sm text-gray-500">{log.user.email}</p>
                          </div>
                        ) : (
                          <Badge variant="outline">Sistema</Badge>
                        )}
                      </TableCell>
                      <TableCell>
                        {ENTITY_LABELS[log.entity] ?? log.entity}
                      </TableCell>
                      <TableCell>
                        <Badge variant="outline">
                          {ACTION_LABELS[log.action] ?? log.action}
                        </Badge>
                      </TableCell>
                      <TableCell className="text-sm text-gray-500">
                        {log.ip_address ?? '-'}
                      </TableCell>
                      <TableCell>
                        <Button variant="outline" size="sm" onClick={() => viewLogDetails(log)}>
                          <Eye className="h-4 w-4" />
                        </Button>
                      </TableCell>
                    </TableRow>
                  ))
                )}
              </TableBody>
            </Table>
          </div>

          {pagination && pagination.pages > 1 && (
            <div className="mt-4 flex items-center justify-between text-sm text-muted-foreground">
              <div>
                Página {pagination.page} de {pagination.pages} — {pagination.total} registros
              </div>
              <div className="flex items-center gap-2">
                <Button
                  variant="outline"
                  size="sm"
                  onClick={() => setPage((prev) => Math.max(1, prev - 1))}
                  disabled={!pagination.has_prev || isLogsLoading}
                >
                  Anterior
                </Button>
                <Button
                  variant="outline"
                  size="sm"
                  onClick={() => setPage((prev) => prev + 1)}
                  disabled={!pagination.has_next || isLogsLoading}
                >
                  Próxima
                </Button>
              </div>
            </div>
          )}
        </CardContent>
      </Card>

      <Dialog open={showDetails} onOpenChange={setShowDetails}>
        <DialogContent className="max-w-xl">
          <DialogHeader>
            <DialogTitle>Detalhes do Log</DialogTitle>
          </DialogHeader>
          {selectedLog && (
            <div className="space-y-4">
              <div>
                <p className="font-medium text-gray-600">Usuário</p>
                <p>{selectedLog.user?.name || 'Sistema'}</p>
                {selectedLog.user?.email && (
                  <p className="text-sm text-gray-500">{selectedLog.user.email}</p>
                )}
              </div>
              <div className="grid grid-cols-2 gap-4 text-sm">
                <div>
                  <p className="font-medium text-gray-600">Entidade</p>
                  <p>{ENTITY_LABELS[selectedLog.entity] ?? selectedLog.entity}</p>
                </div>
                <div>
                  <p className="font-medium text-gray-600">Ação</p>
                  <p>{ACTION_LABELS[selectedLog.action] ?? selectedLog.action}</p>
                </div>
                <div>
                  <p className="font-medium text-gray-600">Data</p>
                  <p>{format(new Date(selectedLog.created_at), 'dd/MM/yyyy HH:mm:ss', { locale: ptBR })}</p>
                </div>
                <div>
                  <p className="font-medium text-gray-600">IP</p>
                  <p>{selectedLog.ip_address ?? '-'}</p>
                </div>
              </div>

              <div>
                <p className="font-medium text-gray-600 mb-2">Alterações</p>
                <div className="rounded border bg-muted p-3 text-sm">
                  {formatChanges(selectedLog.changes)}
                </div>
              </div>

              {selectedLog.user_agent && (
                <div>
                  <p className="font-medium text-gray-600 mb-2">User Agent</p>
                  <p className="text-sm text-gray-500 break-all">{selectedLog.user_agent}</p>
                </div>
              )}
            </div>
          )}
        </DialogContent>
      </Dialog>
    </div>
  )
}
