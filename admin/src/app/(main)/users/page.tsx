'use client'

import { Card, CardHeader, CardTitle, CardContent } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Users, Plus, UserCheck } from 'lucide-react'

export default function UsersPage() {
  return (
    <div className="space-y-6">
      <div className="flex justify-between items-center">
        <h1 className="text-3xl font-bold">Usuários</h1>
        <Button>
          <Plus className="mr-2 h-4 w-4" />
          Novo Usuário
        </Button>
      </div>

      <Card>
        <CardHeader>
          <CardTitle className="flex items-center">
            <Users className="mr-2 h-5 w-5" />
            Gestão de Usuários
          </CardTitle>
        </CardHeader>
        <CardContent>
          <div className="bg-gray-100 rounded-lg h-64 flex items-center justify-center">
            <div className="text-center">
              <UserCheck className="mx-auto h-12 w-12 text-gray-400 mb-4" />
              <p className="text-gray-600">
                Interface de usuários em desenvolvimento
              </p>
              <p className="text-sm text-gray-500 mt-2">
                CRUD completo de usuários será implementado aqui
              </p>
            </div>
          </div>
        </CardContent>
      </Card>
    </div>
  )
}