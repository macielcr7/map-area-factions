'use client'

import { useState } from 'react'
import { useFactions, useCreateFaction, useUpdateFaction } from '@/lib/queries'
import { Button } from '@/components/ui/button'
import { Card, CardHeader, CardTitle, CardContent } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { Plus, Palette, Edit2, Trash2 } from 'lucide-react'
import { FactionDialog } from '@/components/factions/faction-dialog'

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

  const mockFactions = [
    {
      id: '1',
      name: 'Primeiro Comando da Capital',
      acronym: 'PCC',
      color_hex: '#ff0000',
      active: true,
      display_priority: 1,
      description: 'Facção criminosa originária do sistema prisional paulista'
    },
    {
      id: '2',
      name: 'Comando Vermelho',
      acronym: 'CV',
      color_hex: '#cc0000',
      active: true,
      display_priority: 2,
      description: 'Organização criminosa do Rio de Janeiro'
    },
    {
      id: '3',
      name: 'Terceiro Comando Puro',
      acronym: 'TCP',
      color_hex: '#800080',
      active: true,
      display_priority: 3,
      description: 'Facção dissidente do Comando Vermelho'
    },
    {
      id: '4',
      name: 'Guardiões do Estado',
      acronym: 'GDE',
      color_hex: '#0000ff',
      active: true,
      display_priority: 4,
      description: 'Facção local do Ceará'
    },
    {
      id: '5',
      name: 'Família do Norte',
      acronym: 'FDN',
      color_hex: '#008000',
      active: false,
      display_priority: 5,
      description: 'Facção do Amazonas'
    }
  ]

  const factions = factionsData?.data || mockFactions

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
          factions?.map((faction: any) => (
            <Card key={faction.id} className="cursor-pointer hover:shadow-lg transition-shadow">
              <CardHeader>
                <div className="flex items-center justify-between">
                  <div className="flex items-center gap-3">
                    <div 
                      className="w-6 h-6 rounded-full border-2 border-gray-200"
                      style={{ backgroundColor: faction.color_hex }}
                    />
                    <div>
                      <CardTitle className="text-lg">{faction.name}</CardTitle>
                      <p className="text-sm text-muted-foreground font-mono">
                        {faction.acronym}
                      </p>
                    </div>
                  </div>
                  <Badge variant={faction.active ? 'default' : 'secondary'}>
                    {faction.active ? 'Ativa' : 'Inativa'}
                  </Badge>
                </div>
              </CardHeader>
              <CardContent>
                <div className="space-y-3">
                  <p className="text-sm text-gray-600">
                    {faction.description || 'Sem descrição'}
                  </p>
                  
                  <div className="flex justify-between text-sm">
                    <span className="text-muted-foreground">Prioridade:</span>
                    <span className="font-medium">#{faction.display_priority}</span>
                  </div>
                  
                  <div className="flex justify-between text-sm">
                    <span className="text-muted-foreground">Cor:</span>
                    <div className="flex items-center gap-2">
                      <div 
                        className="w-4 h-4 rounded border"
                        style={{ backgroundColor: faction.color_hex }}
                      />
                      <span className="font-mono text-xs">{faction.color_hex}</span>
                    </div>
                  </div>

                  <div className="flex gap-2 mt-4">
                    <Button
                      variant="outline"
                      size="sm"
                      className="flex-1"
                      onClick={() => {
                        setEditingFaction(faction)
                        setIsDialogOpen(true)
                      }}
                    >
                      <Edit2 className="mr-2 h-4 w-4" />
                      Editar
                    </Button>
                    <Button
                      variant="outline"
                      size="sm"
                      disabled
                    >
                      <Trash2 className="h-4 w-4" />
                    </Button>
                  </div>
                </div>
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