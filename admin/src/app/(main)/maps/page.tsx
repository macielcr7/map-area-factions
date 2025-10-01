'use client'

import { useState } from 'react'
import { useFactions, useGeometries, useCreateGeometry, useUpdateGeometry, useDeleteGeometry } from '@/lib/queries'
import { Card, CardHeader, CardTitle, CardContent } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Badge } from '@/components/ui/badge'
import { Switch } from '@/components/ui/switch'
import { MapEditor } from '@/components/maps/map-editor'
import { 
  Map, 
  Plus, 
  Globe, 
  Search, 
  Layers, 
  Edit3, 
  Square, 
  Minus,
  Eye,
  EyeOff,
  Filter,
  MapPin,
  BarChart3
} from 'lucide-react'
import { toast } from 'sonner'

interface LayerVisibility {
  [key: string]: boolean
}

interface SelectedGeometry {
  id: string
  faction_id: string
  risk_level: number
  region_name: string
  area: number
}

export default function MapsPage() {
  const [searchTerm, setSearchTerm] = useState('')
  const [selectedTool, setSelectedTool] = useState<'select' | 'polygon' | 'line' | 'edit'>('select')
  const [layerVisibility, setLayerVisibility] = useState<LayerVisibility>({})
  const [selectedGeometry, setSelectedGeometry] = useState<SelectedGeometry | null>(null)
  const [mapCenter] = useState<[number, number]>([-38.5267, -3.7172])
  const [mapZoom] = useState(10)

  // API hooks
  const { data: factionsData } = useFactions({ active: true })
  const { data: geometriesData } = useGeometries()
  const createGeometryMutation = useCreateGeometry()
  const updateGeometryMutation = useUpdateGeometry()
  const deleteGeometryMutation = useDeleteGeometry()

  const factions = factionsData?.data || []
  const geometries = geometriesData?.data || []

  // Initialize layer visibility
  const initializeLayerVisibility = () => {
    const visibility: LayerVisibility = {}
    factions.forEach((faction: any) => {
      if (!(faction.id in layerVisibility)) {
        visibility[faction.id] = true
      }
    })
    setLayerVisibility(prev => ({ ...prev, ...visibility }))
  }

  // Calculate statistics
  const stats = {
    total_areas: geometries.length,
    total_factions: factions.length,
    coverage_km2: geometries.reduce((sum: number, geo: any) => sum + (geo.area_km2 || 0), 0),
    by_faction: factions.map((faction: any) => ({
      id: faction.id,
      name: faction.name,
      color: faction.color_hex,
      count: geometries.filter((geo: any) => geo.faction_id === faction.id).length
    }))
  }

  const toggleLayerVisibility = (factionId: string) => {
    setLayerVisibility(prev => ({
      ...prev,
      [factionId]: !prev[factionId]
    }))
  }

  const handleToolSelect = (tool: 'select' | 'polygon' | 'line' | 'edit') => {
    setSelectedTool(tool)
    const toolNames = {
      select: 'Seleção',
      polygon: 'Polígono',
      line: 'Linha',
      edit: 'Edição'
    }
    toast.info(`Ferramenta selecionada: ${toolNames[tool]}`)
  }

  const handleGeometryCreate = async (feature: any) => {
    try {
      // Prepare geometry data for API
      const geometryData = {
        geometry: JSON.stringify(feature.geometry),
        faction_id: factions[0]?.id || '', // TODO: Allow selecting faction
        region_name: 'Nova Área',
        risk_level: 1,
        active: true
      }

      await createGeometryMutation.mutateAsync(geometryData)
      toast.success('Geometria criada com sucesso!')
    } catch (error) {
      console.error('Error creating geometry:', error)
      toast.error('Erro ao criar geometria')
    }
  }

  const handleGeometryUpdate = async (feature: any) => {
    try {
      if (!feature.id) return
      
      const geometryData = {
        geometry: JSON.stringify(feature.geometry)
      }

      await updateGeometryMutation.mutateAsync({
        id: feature.id as string,
        data: geometryData
      })
      toast.success('Geometria atualizada com sucesso!')
    } catch (error) {
      console.error('Error updating geometry:', error)
      toast.error('Erro ao atualizar geometria')
    }
  }

  const handleGeometryDelete = async (id: string) => {
    try {
      await deleteGeometryMutation.mutateAsync(id)
      toast.success('Geometria removida com sucesso!')
    } catch (error) {
      console.error('Error deleting geometry:', error)
      toast.error('Erro ao remover geometria')
    }
  }

  const handleGeometrySelect = (geometry: any) => {
    setSelectedGeometry({
      id: geometry.id,
      faction_id: geometry.faction_id,
      risk_level: geometry.risk_level || 1,
      region_name: geometry.region_name || 'Indefinida',
      area: geometry.area_km2 || 0
    })
  }

  return (
    <div className="space-y-6">
      <div className="flex justify-between items-center">
        <h1 className="text-3xl font-bold">Editor de Mapas</h1>
        <div className="flex gap-2">
          <Button variant="outline">
            <Filter className="mr-2 h-4 w-4" />
            Filtros
          </Button>
        </div>
      </div>

      {/* Statistics Cards */}
      <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
        <Card>
          <CardContent className="pt-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-gray-600">Total de Áreas</p>
                <p className="text-2xl font-bold">{stats.total_areas}</p>
              </div>
              <MapPin className="h-8 w-8 text-blue-600" />
            </div>
          </CardContent>
        </Card>
        
        <Card>
          <CardContent className="pt-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-gray-600">Facções Ativas</p>
                <p className="text-2xl font-bold">{stats.total_factions}</p>
              </div>
              <Layers className="h-8 w-8 text-green-600" />
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardContent className="pt-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-gray-600">Cobertura</p>
                <p className="text-2xl font-bold">{stats.coverage_km2.toFixed(1)} km²</p>
              </div>
              <Globe className="h-8 w-8 text-purple-600" />
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardContent className="pt-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-gray-600">Última Atualização</p>
                <p className="text-lg font-bold">Agora</p>
              </div>
              <BarChart3 className="h-8 w-8 text-orange-600" />
            </div>
          </CardContent>
        </Card>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-4 gap-6">
        {/* Map Editor */}
        <div className="lg:col-span-3">
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center justify-between">
                <div className="flex items-center">
                  <Map className="mr-2 h-5 w-5" />
                  Editor Interativo - Fortaleza, CE
                </div>
                <div className="flex items-center gap-2 text-sm text-gray-500">
                  <span>Zoom: {mapZoom}</span>
                  <span>|</span>
                  <span>{mapCenter[1].toFixed(4)}, {mapCenter[0].toFixed(4)}</span>
                </div>
              </CardTitle>
            </CardHeader>
            <CardContent>
              {/* Toolbar */}
              <div className="flex items-center gap-2 mb-4 p-2 bg-gray-50 rounded-lg">
                <Button 
                  size="sm" 
                  variant={selectedTool === 'select' ? 'default' : 'outline'}
                  onClick={() => handleToolSelect('select')}
                >
                  <MapPin className="h-4 w-4" />
                </Button>
                <Button 
                  size="sm" 
                  variant={selectedTool === 'polygon' ? 'default' : 'outline'}
                  onClick={() => handleToolSelect('polygon')}
                >
                  <Square className="h-4 w-4" />
                </Button>
                <Button 
                  size="sm" 
                  variant={selectedTool === 'line' ? 'default' : 'outline'}
                  onClick={() => handleToolSelect('line')}
                >
                  <Minus className="h-4 w-4" />
                </Button>
                <Button 
                  size="sm" 
                  variant={selectedTool === 'edit' ? 'default' : 'outline'}
                  onClick={() => handleToolSelect('edit')}
                >
                  <Edit3 className="h-4 w-4" />
                </Button>
                <div className="flex-1" />
                <div className="flex items-center gap-2">
                  <Search className="h-4 w-4 text-gray-400" />
                  <Input
                    placeholder="Buscar local..."
                    value={searchTerm}
                    onChange={(e) => setSearchTerm(e.target.value)}
                    className="w-64"
                  />
                </div>
              </div>

              {/* Map Container */}
              <MapEditor
                center={mapCenter}
                zoom={mapZoom}
                selectedTool={selectedTool}
                factions={factions}
                geometries={geometries}
                layerVisibility={layerVisibility}
                onGeometryCreate={handleGeometryCreate}
                onGeometryUpdate={handleGeometryUpdate}
                onGeometryDelete={handleGeometryDelete}
              />
            </CardContent>
          </Card>
        </div>

        {/* Controls Panel */}
        <div className="space-y-4">
          {/* Drawing Tools */}
          <Card>
            <CardHeader>
              <CardTitle className="text-lg">Ferramentas</CardTitle>
            </CardHeader>
            <CardContent className="space-y-3">
              <Button 
                variant={selectedTool === 'select' ? 'default' : 'outline'} 
                className="w-full justify-start"
                onClick={() => handleToolSelect('select')}
              >
                <MapPin className="mr-2 h-4 w-4" />
                Selecionar
              </Button>
              <Button 
                variant={selectedTool === 'polygon' ? 'default' : 'outline'} 
                className="w-full justify-start"
                onClick={() => handleToolSelect('polygon')}
              >
                <Square className="mr-2 h-4 w-4" />
                Desenhar Polígono
              </Button>
              <Button 
                variant={selectedTool === 'line' ? 'default' : 'outline'} 
                className="w-full justify-start"
                onClick={() => handleToolSelect('line')}
              >
                <Minus className="mr-2 h-4 w-4" />
                Desenhar Linha
              </Button>
              <Button 
                variant={selectedTool === 'edit' ? 'default' : 'outline'} 
                className="w-full justify-start"
                onClick={() => handleToolSelect('edit')}
              >
                <Edit3 className="mr-2 h-4 w-4" />
                Editar Geometria
              </Button>
            </CardContent>
          </Card>

          {/* Layer Controls */}
          <Card>
            <CardHeader>
              <CardTitle className="text-lg">Camadas por Facção</CardTitle>
            </CardHeader>
            <CardContent>
              <div className="space-y-3">
                {factions.map((faction: any) => (
                  <div key={faction.id} className="flex items-center justify-between">
                    <div className="flex items-center space-x-2">
                      <div 
                        className="w-3 h-3 rounded-full" 
                        style={{ backgroundColor: faction.color_hex }}
                      />
                      <span className="text-sm font-medium">{faction.name}</span>
                      <Badge variant="outline" className="text-xs">
                        {stats.by_faction.find((f: { id: any }) => f.id === faction.id)?.count || 0}
                      </Badge>
                    </div>
                    <Switch 
                      checked={layerVisibility[faction.id] !== false}
                      onCheckedChange={() => toggleLayerVisibility(faction.id)}
                    />
                  </div>
                ))}
              </div>
            </CardContent>
          </Card>

          {/* Properties Panel */}
          <Card>
            <CardHeader>
              <CardTitle className="text-lg">Propriedades</CardTitle>
            </CardHeader>
            <CardContent>
              {selectedGeometry ? (
                <div className="space-y-3">
                  <div>
                    <Label className="text-sm font-medium">Facção</Label>
                    <div className="flex items-center gap-2 mt-1">
                      {(() => {
                        const faction = factions.find((f: any) => f.id === selectedGeometry.faction_id)
                        return faction ? (
                          <>
                            <div 
                              className="w-3 h-3 rounded-full" 
                              style={{ backgroundColor: faction.color_hex }}
                            />
                            <span className="text-sm">{faction.name}</span>
                          </>
                        ) : (
                          <span className="text-sm text-gray-500">Não definida</span>
                        )
                      })()}
                    </div>
                  </div>
                  
                  <div>
                    <Label className="text-sm font-medium">Nível de Risco</Label>
                    <div className="flex items-center gap-1 mt-1">
                      {[1, 2, 3, 4, 5].map(level => (
                        <div 
                          key={level}
                          className={`w-3 h-3 rounded-full ${
                            level <= selectedGeometry.risk_level 
                              ? 'bg-red-500' 
                              : 'bg-gray-200'
                          }`}
                        />
                      ))}
                      <span className="text-sm ml-2">{selectedGeometry.risk_level}/5</span>
                    </div>
                  </div>
                  
                  <div>
                    <Label className="text-sm font-medium">Região</Label>
                    <p className="text-sm mt-1">{selectedGeometry.region_name}</p>
                  </div>
                  
                  <div>
                    <Label className="text-sm font-medium">Área</Label>
                    <p className="text-sm mt-1">{selectedGeometry.area.toFixed(2)} km²</p>
                  </div>

                  <Button size="sm" className="w-full mt-4">
                    <Edit3 className="mr-2 h-4 w-4" />
                    Editar Propriedades
                  </Button>
                </div>
              ) : (
                <div className="text-center py-4">
                  <MapPin className="mx-auto h-8 w-8 text-gray-400 mb-2" />
                  <p className="text-sm text-gray-600">
                    Selecione uma geometria no mapa para editar suas propriedades
                  </p>
                </div>
              )}
            </CardContent>
          </Card>

          {/* Quick Stats */}
          <Card>
            <CardHeader>
              <CardTitle className="text-lg">Estatísticas</CardTitle>
            </CardHeader>
            <CardContent>
              <div className="space-y-2 text-sm">
                <div className="flex justify-between">
                  <span>Áreas mapeadas:</span>
                  <span className="font-medium">{stats.total_areas}</span>
                </div>
                <div className="flex justify-between">
                  <span>Cobertura total:</span>
                  <span className="font-medium">{stats.coverage_km2.toFixed(1)} km²</span>
                </div>
                <div className="flex justify-between">
                  <span>Facções ativas:</span>
                  <span className="font-medium">{stats.total_factions}</span>
                </div>
                <div className="flex justify-between">
                  <span>Última edição:</span>
                  <span className="font-medium">Agora</span>
                </div>
              </div>
            </CardContent>
          </Card>
        </div>
      </div>
    </div>
  )
}