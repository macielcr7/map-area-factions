'use client'

import { useState, useEffect } from 'react'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import * as z from 'zod'
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogFooter,
} from '@/components/ui/dialog'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Textarea } from '@/components/ui/textarea'
import { Switch } from '@/components/ui/switch'
import { Loader2, Palette } from 'lucide-react'

const factionSchema = z.object({
  name: z.string().min(1, 'Nome é obrigatório'),
  acronym: z.string().min(1, 'Sigla é obrigatória').max(10, 'Sigla muito longa'),
  color_hex: z.string().regex(/^#[0-9A-F]{6}$/i, 'Cor deve ser um hex válido'),
  description: z.string().optional(),
  active: z.boolean(),
  display_priority: z.number().min(1).max(100),
})

type FactionFormData = z.infer<typeof factionSchema>

interface FactionDialogProps {
  open: boolean
  onOpenChange: (open: boolean) => void
  faction?: any
  onSubmit: (data: FactionFormData) => void
  isLoading: boolean
}

const predefinedColors = [
  '#ff0000', '#cc0000', '#800000', // Vermelhos
  '#0000ff', '#0066cc', '#003d99', // Azuis
  '#008000', '#00cc00', '#009900', // Verdes
  '#800080', '#cc00cc', '#990099', // Roxos
  '#ff8000', '#ff9900', '#cc6600', // Laranjas
  '#808080', '#666666', '#404040', // Cinzas
]

export function FactionDialog({ 
  open, 
  onOpenChange, 
  faction, 
  onSubmit, 
  isLoading 
}: FactionDialogProps) {
  const [selectedColor, setSelectedColor] = useState('#ff0000')

  const {
    register,
    handleSubmit,
    reset,
    setValue,
    watch,
    formState: { errors },
  } = useForm<FactionFormData>({
    resolver: zodResolver(factionSchema),
    defaultValues: {
      name: '',
      acronym: '',
      color_hex: '#ff0000',
      description: '',
      active: true,
      display_priority: 1,
    },
  })

  const watchedColor = watch('color_hex')

  useEffect(() => {
    setSelectedColor(watchedColor)
  }, [watchedColor])

  useEffect(() => {
    if (faction) {
      reset({
        name: faction.name || '',
        acronym: faction.acronym || '',
        color_hex: faction.color_hex || '#ff0000',
        description: faction.description || '',
        active: faction.active ?? true,
        display_priority: faction.display_priority || 1,
      })
      setSelectedColor(faction.color_hex || '#ff0000')
    } else {
      reset({
        name: '',
        acronym: '',
        color_hex: '#ff0000',
        description: '',
        active: true,
        display_priority: 1,
      })
      setSelectedColor('#ff0000')
    }
  }, [faction, reset])

  const handleColorSelect = (color: string) => {
    setSelectedColor(color)
    setValue('color_hex', color)
  }

  const onSubmitForm = (data: FactionFormData) => {
    onSubmit(data)
  }

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-[500px]">
        <DialogHeader>
          <DialogTitle>
            {faction ? 'Editar Facção' : 'Nova Facção'}
          </DialogTitle>
        </DialogHeader>

        <form onSubmit={handleSubmit(onSubmitForm)} className="space-y-4">
          <div className="grid grid-cols-2 gap-4">
            <div className="space-y-2">
              <Label htmlFor="name">Nome</Label>
              <Input
                id="name"
                {...register('name')}
                placeholder="Ex: Primeiro Comando da Capital"
              />
              {errors.name && (
                <p className="text-sm text-red-600">{errors.name.message}</p>
              )}
            </div>

            <div className="space-y-2">
              <Label htmlFor="acronym">Sigla</Label>
              <Input
                id="acronym"
                {...register('acronym')}
                placeholder="Ex: PCC"
                maxLength={10}
              />
              {errors.acronym && (
                <p className="text-sm text-red-600">{errors.acronym.message}</p>
              )}
            </div>
          </div>

          <div className="space-y-2">
            <Label htmlFor="description">Descrição</Label>
            <Textarea
              id="description"
              {...register('description')}
              placeholder="Descrição opcional da facção"
              rows={3}
            />
          </div>

          <div className="space-y-2">
            <Label>Cor da Facção</Label>
            <div className="flex items-center gap-3">
              <div
                className="w-12 h-12 rounded-lg border-2 border-gray-200 flex items-center justify-center"
                style={{ backgroundColor: selectedColor }}
              >
                <Palette className="w-5 h-5 text-white mix-blend-difference" />
              </div>
              <Input
                {...register('color_hex')}
                placeholder="#ff0000"
                className="font-mono"
                onChange={(e) => setSelectedColor(e.target.value)}
              />
            </div>
            {errors.color_hex && (
              <p className="text-sm text-red-600">{errors.color_hex.message}</p>
            )}
            
            <div className="grid grid-cols-6 gap-2 mt-2">
              {predefinedColors.map((color) => (
                <button
                  key={color}
                  type="button"
                  className={`w-8 h-8 rounded-md border-2 ${
                    selectedColor === color ? 'border-gray-800' : 'border-gray-200'
                  }`}
                  style={{ backgroundColor: color }}
                  onClick={() => handleColorSelect(color)}
                />
              ))}
            </div>
          </div>

          <div className="grid grid-cols-2 gap-4">
            <div className="space-y-2">
              <Label htmlFor="display_priority">Prioridade de Exibição</Label>
              <Input
                id="display_priority"
                type="number"
                {...register('display_priority', { valueAsNumber: true })}
                min={1}
                max={100}
              />
              {errors.display_priority && (
                <p className="text-sm text-red-600">{errors.display_priority.message}</p>
              )}
            </div>

            <div className="flex items-center space-x-2">
              <Switch
                id="active"
                {...register('active')}
                defaultChecked={true}
              />
              <Label htmlFor="active">Facção ativa</Label>
            </div>
          </div>

          <DialogFooter>
            <Button
              type="button"
              variant="outline"
              onClick={() => onOpenChange(false)}
              disabled={isLoading}
            >
              Cancelar
            </Button>
            <Button type="submit" disabled={isLoading}>
              {isLoading && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
              {faction ? 'Atualizar' : 'Criar'}
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  )
}