'use client'

import { Card, CardHeader, CardTitle, CardContent } from '@/components/ui/card'
import { FileText } from 'lucide-react'

export default function ReportsPage() {
  return (
    <div className="space-y-6">
      <h1 className="text-3xl font-bold">Relatórios</h1>

      <Card>
        <CardHeader>
          <CardTitle className="flex items-center">
            <FileText className="mr-2 h-5 w-5" />
            Sistema de Relatórios
          </CardTitle>
        </CardHeader>
        <CardContent>
          <div className="bg-gray-100 rounded-lg h-64 flex items-center justify-center">
            <div className="text-center">
              <FileText className="mx-auto h-12 w-12 text-gray-400 mb-4" />
              <p className="text-gray-600">
                Relatórios em desenvolvimento
              </p>
            </div>
          </div>
        </CardContent>
      </Card>
    </div>
  )
}