'use client'

import { useEffect, useMemo, useState } from 'react'
import { Card, CardHeader, CardTitle, CardContent } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Badge } from '@/components/ui/badge'
import { Textarea } from '@/components/ui/textarea'
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
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import {
  FileText,
  Search,
  AlertTriangle,
  CheckCircle,
  Clock,
  Eye,
  Edit,
  Loader2,
  RefreshCw,
} from 'lucide-react'
import {
  useReports,
  useUpdateReport,
} from '@/lib/queries'
import { toast } from 'sonner'

const PAGE_SIZE = 20

const statusLabels: Record<ReportStatus, string> = {
  pending: 'Pendente',
  reviewing: 'Em Análise',
  resolved: 'Resolvido',
  rejected: 'Rejeitado',
}

const statusColors: Record<ReportStatus, string> = {
  pending: 'bg-yellow-500',
  reviewing: 'bg-blue-500',
  resolved: 'bg-green-500',
  rejected: 'bg-red-500',
}

const typeLabels: Record<ReportType, string> = {
  inaccuracy: 'Imprecisão',
  inappropriate: 'Conteúdo Inadequado',
  spam: 'Spam',
  outdated: 'Desatualizado',
  other: 'Outros',
}

const typeColors: Record<ReportType, string> = {
  inaccuracy: 'bg-orange-500',
  inappropriate: 'bg-red-500',
  spam: 'bg-gray-500',
  outdated: 'bg-blue-500',
  other: 'bg-purple-500',
}

type ReportStatus = 'pending' | 'reviewing' | 'resolved' | 'rejected'
type ReportType = 'inaccuracy' | 'inappropriate' | 'spam' | 'outdated' | 'other'

type ReportUser = {
  id?: string
  name: string
  email: string
}

type AdminReport = {
  id: string
  user: ReportUser | null
  geometry_id: string | null
  type: ReportType
  description: string
  attachments: string[]
  status: ReportStatus
  reporter_lat: number | null
  reporter_lng: number | null
  created_at: string
  assigned_to: ReportUser | null
  reviewed_by: ReportUser | null
  review_notes: string | null
}

type ReviewFormData = {
  status: ReportStatus
  review_notes: string
  assigned_to_id: string
}

export default function ReportsPage() {
  const [page, setPage] = useState(1)
  const [searchTerm, setSearchTerm] = useState('')
  const [statusFilter, setStatusFilter] = useState<ReportStatus | ''>('')
  const [typeFilter, setTypeFilter] = useState<ReportType | ''>('')
  const [selectedReport, setSelectedReport] = useState<AdminReport | null>(null)
  const [showDetails, setShowDetails] = useState(false)
  const [showReview, setShowReview] = useState(false)
  const [reviewData, setReviewData] = useState<ReviewFormData>({
    status: 'pending',
    review_notes: '',
    assigned_to_id: '',
  })

  const { data, isLoading, isFetching, isError, error } = useReports({
    page,
    limit: PAGE_SIZE,
    status: statusFilter || undefined,
    type: typeFilter || undefined,
  })

  const updateReportMutation = useUpdateReport()

  useEffect(() => {
    if (!isError || !error) {
      return
    }
    toast.error('Erro ao carregar relatórios')
    // eslint-disable-next-line no-console
    console.error(error)
  }, [isError, error])

  const reports = useMemo(() => (data?.reports ?? []) as AdminReport[], [data?.reports])
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

  const filteredReports = useMemo(() => {
    if (!searchTerm) {
      return reports
    }

    const normalizedSearch = searchTerm.toLowerCase()

    return reports.filter((report) => {
      const reporterName = report.user?.name?.toLowerCase() ?? 'anônimo'
      const reporterEmail = report.user?.email?.toLowerCase() ?? ''
      return (
        reporterName.includes(normalizedSearch) ||
        reporterEmail.includes(normalizedSearch) ||
        report.description.toLowerCase().includes(normalizedSearch) ||
        typeLabels[report.type].toLowerCase().includes(normalizedSearch)
      )
    })
  }, [reports, searchTerm])

  const stats = useMemo(() => {
    const totalReports = reports.length
    const pendingReports = reports.filter((report) => report.status === 'pending').length
    const resolvedReports = reports.filter((report) => report.status === 'resolved').length
    const byType = reports.reduce<Record<ReportType, number>>((acc, report) => {
      acc[report.type] = (acc[report.type] ?? 0) + 1
      return acc
    }, {
      inaccuracy: 0,
      inappropriate: 0,
      spam: 0,
      outdated: 0,
      other: 0,
    })

    return {
      total_reports: totalReports,
      pending_reports: pendingReports,
      resolved_reports: resolvedReports,
      by_type: byType,
    }
  }, [reports])

  const isReportsLoading = isLoading || isFetching

  const openReviewDialog = (report: AdminReport) => {
    setSelectedReport(report)
    setReviewData({
      status: report.status,
      review_notes: report.review_notes ?? '',
      assigned_to_id: report.assigned_to?.id ?? '',
    })
    setShowReview(true)
  }

  const handleReviewSubmit = async () => {
    if (!selectedReport) {
      return
    }

    try {
      await updateReportMutation.mutateAsync({
        id: selectedReport.id,
        data: {
          status: reviewData.status,
          review_notes: reviewData.review_notes,
          assigned_to_id: reviewData.assigned_to_id || undefined,
        },
      })
      setShowReview(false)
      setSelectedReport(null)
    } catch (mutationError) {
      // eslint-disable-next-line no-console
      console.error(mutationError)
    }
  }

  const resetFilters = () => {
    setSearchTerm('')
    setStatusFilter('')
    setTypeFilter('')
    setPage(1)
  }

  return (
    <div className="space-y-6">
      <div className="flex justify-between items-center">
        <h1 className="text-3xl font-bold">Gestão de Relatórios</h1>
        <Button variant="outline" onClick={resetFilters}>
          <RefreshCw className="mr-2 h-4 w-4" />
          Resetar Filtros
        </Button>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
        <Card>
          <CardContent className="pt-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-gray-600">Total</p>
                <p className="text-2xl font-bold">{stats.total_reports}</p>
              </div>
              <FileText className="h-8 w-8 text-blue-600" />
            </div>
          </CardContent>
        </Card>
        <Card>
          <CardContent className="pt-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-gray-600">Pendentes</p>
                <p className="text-2xl font-bold text-yellow-600">{stats.pending_reports}</p>
              </div>
              <Clock className="h-8 w-8 text-yellow-600" />
            </div>
          </CardContent>
        </Card>
        <Card>
          <CardContent className="pt-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-gray-600">Resolvidos</p>
                <p className="text-2xl font-bold text-green-600">{stats.resolved_reports}</p>
              </div>
              <CheckCircle className="h-8 w-8 text-green-600" />
            </div>
          </CardContent>
        </Card>
        <Card>
          <CardContent className="pt-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-gray-600">Imprecisões</p>
                <p className="text-2xl font-bold text-orange-600">{stats.by_type.inaccuracy}</p>
              </div>
              <AlertTriangle className="h-8 w-8 text-orange-600" />
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
                  placeholder="Buscar por usuário, descrição ou tipo..."
                  value={searchTerm}
                  onChange={(e) => setSearchTerm(e.target.value)}
                  className="pl-10"
                />
              </div>
            </div>
            <div className="w-full md:w-48">
              <select
                value={statusFilter}
                onChange={(e) => {
                  setStatusFilter(e.target.value as ReportStatus | '')
                  setPage(1)
                }}
                className="w-full p-2 border rounded-md"
              >
                <option value="">Todos os status</option>
                <option value="pending">Pendente</option>
                <option value="reviewing">Em Análise</option>
                <option value="resolved">Resolvido</option>
                <option value="rejected">Rejeitado</option>
              </select>
            </div>
            <div className="w-full md:w-48">
              <select
                value={typeFilter}
                onChange={(e) => {
                  setTypeFilter(e.target.value as ReportType | '')
                  setPage(1)
                }}
                className="w-full p-2 border rounded-md"
              >
                <option value="">Todos os tipos</option>
                <option value="inaccuracy">Imprecisão</option>
                <option value="inappropriate">Inadequado</option>
                <option value="spam">Spam</option>
                <option value="outdated">Desatualizado</option>
                <option value="other">Outros</option>
              </select>
            </div>
          </div>
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <CardTitle className="flex items-center">
            <FileText className="mr-2 h-5 w-5" />
            Relatórios ({pagination?.total ?? filteredReports.length})
          </CardTitle>
        </CardHeader>
        <CardContent>
          <div className="rounded-md border">
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>Data</TableHead>
                  <TableHead>Usuário</TableHead>
                  <TableHead>Tipo</TableHead>
                  <TableHead>Status</TableHead>
                  <TableHead>Descrição</TableHead>
                  <TableHead>Responsável</TableHead>
                  <TableHead>Ações</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {isReportsLoading ? (
                  <TableRow>
                    <TableCell colSpan={7} className="py-8 text-center">
                      <div className="flex items-center justify-center gap-2 text-sm text-muted-foreground">
                        <Loader2 className="h-4 w-4 animate-spin" />
                        Carregando relatórios...
                      </div>
                    </TableCell>
                  </TableRow>
                ) : filteredReports.length === 0 ? (
                  <TableRow>
                    <TableCell colSpan={7} className="py-8 text-center text-sm text-muted-foreground">
                      Nenhum relatório encontrado.
                    </TableCell>
                  </TableRow>
                ) : (
                  filteredReports.map((report) => (
                    <TableRow key={report.id}>
                      <TableCell>
                        {new Date(report.created_at).toLocaleDateString('pt-BR')}
                      </TableCell>
                      <TableCell>
                        {report.user ? (
                          <div>
                            <p className="font-medium">{report.user.name}</p>
                            <p className="text-sm text-gray-500">{report.user.email}</p>
                          </div>
                        ) : (
                          <Badge variant="outline">Anônimo</Badge>
                        )}
                      </TableCell>
                      <TableCell>
                        <Badge className={`${typeColors[report.type]} text-white`}>
                          {typeLabels[report.type]}
                        </Badge>
                      </TableCell>
                      <TableCell>
                        <Badge className={`${statusColors[report.status]} text-white`}>
                          {statusLabels[report.status]}
                        </Badge>
                      </TableCell>
                      <TableCell className="max-w-xs">
                        <p className="truncate">{report.description}</p>
                      </TableCell>
                      <TableCell>
                        {report.assigned_to ? (
                          <p className="text-sm">{report.assigned_to.name}</p>
                        ) : (
                          <Badge variant="outline">Não atribuído</Badge>
                        )}
                      </TableCell>
                      <TableCell>
                        <div className="flex gap-2">
                          <Button
                            variant="outline"
                            size="sm"
                            onClick={() => {
                              setSelectedReport(report)
                              setShowDetails(true)
                            }}
                          >
                            <Eye className="h-4 w-4" />
                          </Button>
                          <Button
                            variant="outline"
                            size="sm"
                            onClick={() => openReviewDialog(report)}
                          >
                            <Edit className="h-4 w-4" />
                          </Button>
                        </div>
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
                Página {pagination.page} de {pagination.pages} — {pagination.total} relatórios
              </div>
              <div className="flex items-center gap-2">
                <Button
                  variant="outline"
                  size="sm"
                  onClick={() => setPage((prev) => Math.max(1, prev - 1))}
                  disabled={!pagination.has_prev || isReportsLoading}
                >
                  Anterior
                </Button>
                <Button
                  variant="outline"
                  size="sm"
                  onClick={() => setPage((prev) => prev + 1)}
                  disabled={!pagination.has_next || isReportsLoading}
                >
                  Próxima
                </Button>
              </div>
            </div>
          )}
        </CardContent>
      </Card>

      <Dialog open={showDetails} onOpenChange={setShowDetails}>
        <DialogContent className="max-w-2xl">
          <DialogHeader>
            <DialogTitle>Detalhes do Relatório</DialogTitle>
          </DialogHeader>
          {selectedReport && (
            <div className="space-y-4">
              <div className="grid grid-cols-2 gap-4">
                <div>
                  <p className="font-medium text-gray-600">Reportado por:</p>
                  <p>{selectedReport.user?.name || 'Usuário Anônimo'}</p>
                  {selectedReport.user && (
                    <p className="text-sm text-gray-500">{selectedReport.user.email}</p>
                  )}
                </div>
                <div>
                  <p className="font-medium text-gray-600">Data:</p>
                  <p>{new Date(selectedReport.created_at).toLocaleString('pt-BR')}</p>
                </div>
                <div>
                  <p className="font-medium text-gray-600">Tipo:</p>
                  <Badge className={`${typeColors[selectedReport.type]} text-white`}>
                    {typeLabels[selectedReport.type]}
                  </Badge>
                </div>
                <div>
                  <p className="font-medium text-gray-600">Status:</p>
                  <Badge className={`${statusColors[selectedReport.status]} text-white`}>
                    {statusLabels[selectedReport.status]}
                  </Badge>
                </div>
              </div>

              <div>
                <p className="font-medium text-gray-600 mb-2">Descrição:</p>
                <p className="bg-gray-100 p-3 rounded">{selectedReport.description}</p>
              </div>

              {selectedReport.attachments && selectedReport.attachments.length > 0 && (
                <div>
                  <p className="font-medium text-gray-600 mb-2">Anexos:</p>
                  <div className="flex gap-2 flex-wrap">
                    {selectedReport.attachments.map((attachment, index) => (
                      <Badge key={index} variant="outline">{attachment}</Badge>
                    ))}
                  </div>
                </div>
              )}

              {selectedReport.reporter_lat !== null && selectedReport.reporter_lng !== null && (
                <div>
                  <p className="font-medium text-gray-600 mb-2">Localização do Relato:</p>
                  <p className="font-mono text-sm">
                    {selectedReport.reporter_lat.toFixed(6)}, {selectedReport.reporter_lng.toFixed(6)}
                  </p>
                </div>
              )}

              {selectedReport.review_notes && (
                <div>
                  <p className="font-medium text-gray-600 mb-2">Notas da Revisão:</p>
                  <p className="bg-blue-50 p-3 rounded border border-blue-200">
                    {selectedReport.review_notes}
                  </p>
                </div>
              )}
            </div>
          )}
        </DialogContent>
      </Dialog>

      <Dialog open={showReview} onOpenChange={setShowReview}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Revisar Relatório</DialogTitle>
            <DialogDescription>
              Atualize o status e adicione notas de revisão para este relatório.
            </DialogDescription>
          </DialogHeader>
          <div className="space-y-4">
            <div>
              <label className="block text-sm font-medium mb-2">Status</label>
              <select
                value={reviewData.status}
                onChange={(e) =>
                  setReviewData((prev) => ({ ...prev, status: e.target.value as ReportStatus }))
                }
                className="w-full p-2 border rounded-md"
              >
                <option value="pending">Pendente</option>
                <option value="reviewing">Em Análise</option>
                <option value="resolved">Resolvido</option>
                <option value="rejected">Rejeitado</option>
              </select>
            </div>

            <div>
              <label className="block text-sm font-medium mb-2">Notas da Revisão</label>
              <Textarea
                value={reviewData.review_notes}
                onChange={(e) => setReviewData((prev) => ({ ...prev, review_notes: e.target.value }))}
                placeholder="Adicione suas observações sobre a revisão..."
                rows={4}
              />
            </div>
          </div>
          <DialogFooter>
            <Button variant="outline" onClick={() => setShowReview(false)}>
              Cancelar
            </Button>
            <Button onClick={handleReviewSubmit} disabled={updateReportMutation.isPending}>
              {updateReportMutation.isPending ? (
                <>
                  <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                  Salvando...
                </>
              ) : (
                'Salvar Revisão'
              )}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  )
}
