'use client'

import { useState } from 'react'
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
  TableRow 
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
  XCircle,
  Clock,
  Eye,
  Edit
} from 'lucide-react'

// Mock reports data - replace with real API calls
const mockReports = [
  {
    id: '1',
    user: { name: 'Ana Silva', email: 'ana@citizen.com' },
    geometry_id: 'uuid-geo-1',
    type: 'inaccuracy',
    description: 'A área marcada como PCC não corresponde à realidade atual do bairro.',
    attachments: ['evidence1.jpg'],
    status: 'pending',
    reporter_lat: -3.7319,
    reporter_lng: -38.5267,
    created_at: '2024-01-20T10:30:00Z',
    assigned_to: null,
    reviewed_by: null,
    review_notes: null
  },
  {
    id: '2',
    user: { name: 'Carlos Santos', email: 'carlos@citizen.com' },
    geometry_id: 'uuid-geo-2',
    type: 'inappropriate',
    description: 'Conteúdo ofensivo nas descrições da área.',
    attachments: [],
    status: 'reviewing',
    reporter_lat: -3.7420,
    reporter_lng: -38.5367,
    created_at: '2024-01-19T14:15:00Z',
    assigned_to: { name: 'Maria Moderadora', email: 'maria@moderator.com' },
    reviewed_by: null,
    review_notes: null
  },
  {
    id: '3',
    user: { name: 'João Oliveira', email: 'joao@citizen.com' },
    geometry_id: 'uuid-geo-3',
    type: 'outdated',
    description: 'Informações desatualizadas sobre domínio territorial.',
    attachments: ['photo1.jpg', 'photo2.jpg'],
    status: 'resolved',
    reporter_lat: -3.7220,
    reporter_lng: -38.5167,
    created_at: '2024-01-18T09:45:00Z',
    assigned_to: { name: 'Pedro Moderador', email: 'pedro@moderator.com' },
    reviewed_by: { name: 'Pedro Moderador', email: 'pedro@moderator.com' },
    review_notes: 'Informações atualizadas conforme evidências fornecidas.'
  },
  {
    id: '4',
    user: null, // Anonymous report
    geometry_id: 'uuid-geo-4',
    type: 'spam',
    description: 'Conteúdo duplicado e sem relevância.',
    attachments: [],
    status: 'rejected',
    reporter_lat: -3.7519,
    reporter_lng: -38.5467,
    created_at: '2024-01-17T16:20:00Z',
    assigned_to: { name: 'Maria Moderadora', email: 'maria@moderator.com' },
    reviewed_by: { name: 'Maria Moderadora', email: 'maria@moderator.com' },
    review_notes: 'Relatório classificado como spam e rejeitado.'
  }
]

const statusColors = {
  pending: 'bg-yellow-500',
  reviewing: 'bg-blue-500',
  resolved: 'bg-green-500',
  rejected: 'bg-red-500'
}

const statusLabels = {
  pending: 'Pendente',
  reviewing: 'Em Análise',
  resolved: 'Resolvido',
  rejected: 'Rejeitado'
}

const typeLabels = {
  inaccuracy: 'Imprecisão',
  inappropriate: 'Conteúdo Inadequado',
  spam: 'Spam',
  outdated: 'Desatualizado',
  other: 'Outros'
}

const typeColors = {
  inaccuracy: 'bg-orange-500',
  inappropriate: 'bg-red-500',
  spam: 'bg-gray-500',
  outdated: 'bg-blue-500',
  other: 'bg-purple-500'
}

export default function ReportsPage() {
  const [reports, setReports] = useState(mockReports)
  const [searchTerm, setSearchTerm] = useState('')
  const [statusFilter, setStatusFilter] = useState('')
  const [typeFilter, setTypeFilter] = useState('')
  const [selectedReport, setSelectedReport] = useState(null)
  const [showDetails, setShowDetails] = useState(false)
  const [showReview, setShowReview] = useState(false)
  const [reviewData, setReviewData] = useState({
    status: '',
    review_notes: '',
    assigned_to_id: ''
  })

  const filteredReports = reports.filter(report => {
    const matchesSearch = (report.user?.name || 'Anônimo').toLowerCase().includes(searchTerm.toLowerCase()) ||
                         report.description.toLowerCase().includes(searchTerm.toLowerCase()) ||
                         report.type.toLowerCase().includes(searchTerm.toLowerCase())
    const matchesStatus = !statusFilter || report.status === statusFilter
    const matchesType = !typeFilter || report.type === typeFilter
    return matchesSearch && matchesStatus && matchesType
  })

  const viewReportDetails = (report) => {
    setSelectedReport(report)
    setShowDetails(true)
  }

  const openReviewDialog = (report) => {
    setSelectedReport(report)
    setReviewData({
      status: report.status,
      review_notes: report.review_notes || '',
      assigned_to_id: report.assigned_to?.email || ''
    })
    setShowReview(true)
  }

  const handleReviewSubmit = () => {
    setReports(reports.map(report => 
      report.id === selectedReport.id 
        ? { 
            ...report, 
            status: reviewData.status,
            review_notes: reviewData.review_notes,
            reviewed_by: { name: 'Moderador Atual', email: 'moderator@admin.com' }
          }
        : report
    ))
    setShowReview(false)
    setSelectedReport(null)
  }

  // Mock statistics
  const stats = {
    total_reports: reports.length,
    pending_reports: reports.filter(r => r.status === 'pending').length,
    resolved_reports: reports.filter(r => r.status === 'resolved').length,
    by_type: {
      inaccuracy: reports.filter(r => r.type === 'inaccuracy').length,
      inappropriate: reports.filter(r => r.type === 'inappropriate').length,
      spam: reports.filter(r => r.type === 'spam').length,
      outdated: reports.filter(r => r.type === 'outdated').length
    }
  }

  return (
    <div className="space-y-6">
      <div className="flex justify-between items-center">
        <h1 className="text-3xl font-bold">Gestão de Relatórios</h1>
        <Button variant="outline">
          <FileText className="mr-2 h-4 w-4" />
          Exportar Dados
        </Button>
      </div>

      {/* Statistics Cards */}
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

      {/* Filters */}
      <Card>
        <CardContent className="pt-6">
          <div className="flex gap-4">
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
            <div className="w-48">
              <select
                value={statusFilter}
                onChange={(e) => setStatusFilter(e.target.value)}
                className="w-full p-2 border rounded-md"
              >
                <option value="">Todos os status</option>
                <option value="pending">Pendente</option>
                <option value="reviewing">Em Análise</option>
                <option value="resolved">Resolvido</option>
                <option value="rejected">Rejeitado</option>
              </select>
            </div>
            <div className="w-48">
              <select
                value={typeFilter}
                onChange={(e) => setTypeFilter(e.target.value)}
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

      {/* Reports Table */}
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center">
            <FileText className="mr-2 h-5 w-5" />
            Relatórios ({filteredReports.length})
          </CardTitle>
        </CardHeader>
        <CardContent>
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
              {filteredReports.map((report) => (
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
                        onClick={() => viewReportDetails(report)}
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
              ))}
            </TableBody>
          </Table>
        </CardContent>
      </Card>

      {/* Report Details Modal */}
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
                  <div className="flex gap-2">
                    {selectedReport.attachments.map((attachment, index) => (
                      <Badge key={index} variant="outline">{attachment}</Badge>
                    ))}
                  </div>
                </div>
              )}

              {selectedReport.reporter_lat && selectedReport.reporter_lng && (
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

      {/* Review Modal */}
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
                onChange={(e) => setReviewData({...reviewData, status: e.target.value})}
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
                onChange={(e) => setReviewData({...reviewData, review_notes: e.target.value})}
                placeholder="Adicione suas observações sobre a revisão..."
                rows={4}
              />
            </div>
          </div>
          <DialogFooter>
            <Button variant="outline" onClick={() => setShowReview(false)}>
              Cancelar
            </Button>
            <Button onClick={handleReviewSubmit}>
              Salvar Revisão
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  )
}