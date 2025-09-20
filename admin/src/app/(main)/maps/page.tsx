'use client'

import { Card, CardHeader, CardTitle, CardContent } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Map, Plus, Globe } from 'lucide-react'

export default function MapsPage() {
  return (
    <div className="space-y-6">
      <div className="flex justify-between items-center">
        <h1 className="text-3xl font-bold">Editor de Mapas</h1>
        <Button>
          <Plus className="mr-2 h-4 w-4" />
          Nova Área
        </Button>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        {/* Map Editor Preview */}
        <div className="lg:col-span-2">
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center">
                <Map className="mr-2 h-5 w-5" />
                Editor Interativo
              </CardTitle>
            </CardHeader>
            <CardContent>
              <div className="bg-gray-100 rounded-lg h-96 flex items-center justify-center">
                <div className="text-center">
                  <Globe className="mx-auto h-12 w-12 text-gray-400 mb-4" />
                  <p className="text-gray-600">
                    Editor de mapas em desenvolvimento
                  </p>
                  <p className="text-sm text-gray-500 mt-2">
                    Integração com Mapbox GL JS será implementada aqui
                  </p>
                </div>
              </div>
            </CardContent>
          </Card>
        </div>

        {/* Controls Panel */}
        <div className="space-y-4">
          <Card>
            <CardHeader>
              <CardTitle>Ferramentas de Desenho</CardTitle>
            </CardHeader>
            <CardContent className="space-y-3">
              <Button variant="outline" className="w-full justify-start">
                <div className="w-4 h-4 bg-blue-500 rounded mr-2" />
                Desenhar Polígono
              </Button>
              <Button variant="outline" className="w-full justify-start">
                <div className="w-4 h-4 bg-green-500 rounded mr-2" />
                Desenhar Linha
              </Button>
              <Button variant="outline" className="w-full justify-start">
                <div className="w-4 h-4 bg-red-500 rounded mr-2" />
                Editar Geometria
              </Button>
            </CardContent>
          </Card>

          <Card>
            <CardHeader>
              <CardTitle>Camadas</CardTitle>
            </CardHeader>
            <CardContent>
              <div className="space-y-2">
                <div className="flex items-center justify-between">
                  <div className="flex items-center space-x-2">
                    <div className="w-3 h-3 bg-red-500 rounded-full" />
                    <span className="text-sm">PCC</span>
                  </div>
                  <input type="checkbox" defaultChecked />
                </div>
                <div className="flex items-center justify-between">
                  <div className="flex items-center space-x-2">
                    <div className="w-3 h-3 bg-blue-500 rounded-full" />
                    <span className="text-sm">CV</span>
                  </div>
                  <input type="checkbox" defaultChecked />
                </div>
                <div className="flex items-center justify-between">
                  <div className="flex items-center space-x-2">
                    <div className="w-3 h-3 bg-green-500 rounded-full" />
                    <span className="text-sm">GDE</span>
                  </div>
                  <input type="checkbox" defaultChecked />
                </div>
              </div>
            </CardContent>
          </Card>

          <Card>
            <CardHeader>
              <CardTitle>Propriedades</CardTitle>
            </CardHeader>
            <CardContent>
              <div className="space-y-2 text-sm">
                <p className="text-gray-600">
                  Selecione uma geometria para editar suas propriedades
                </p>
                <div className="border rounded p-2 bg-gray-50">
                  <p><strong>Facção:</strong> Não selecionada</p>
                  <p><strong>Risco:</strong> Não definido</p>
                  <p><strong>Região:</strong> Não definida</p>
                </div>
              </div>
            </CardContent>
          </Card>
        </div>
      </div>
    </div>
  )
}