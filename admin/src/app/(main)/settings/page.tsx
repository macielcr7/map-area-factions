'use client'

import { useEffect, useMemo, useState } from 'react'
import { Card, CardHeader, CardTitle, CardContent } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Switch } from '@/components/ui/switch'
import { Textarea } from '@/components/ui/textarea'
import {
  Settings,
  Globe,
  Shield,
  Bell,
  Palette,
  Database,
  MapPin,
  Save,
  RefreshCw,
  Loader2,
} from 'lucide-react'
import { toast } from 'sonner'
import { useSettings, useUpdateSettings } from '@/lib/queries'

interface SystemSettings {
  general: {
    site_name: string
    site_description: string
    maintenance_mode: boolean
    timezone: string
    default_language: string
  }
  map: {
    default_style: string
    default_zoom: number
    default_center_lat: number
    default_center_lng: number
    max_zoom: number
    min_zoom: number
  }
  security: {
    session_timeout_minutes: number
    require_email_verification: boolean
    enable_two_factor: boolean
    password_min_length: number
    enable_audit_logging: boolean
  }
  notifications: {
    email_enabled: boolean
    email_smtp_host: string
    email_smtp_port: number
    push_enabled: boolean
    webhook_url: string
    notify_new_reports: boolean
  }
}

const DEFAULT_SETTINGS: SystemSettings = {
  general: {
    site_name: '',
    site_description: '',
    maintenance_mode: false,
    timezone: 'America/Sao_Paulo',
    default_language: 'pt-BR',
  },
  map: {
    default_style: '',
    default_zoom: 10,
    default_center_lat: -3.7319,
    default_center_lng: -38.5267,
    max_zoom: 18,
    min_zoom: 8,
  },
  security: {
    session_timeout_minutes: 480,
    require_email_verification: true,
    enable_two_factor: false,
    password_min_length: 8,
    enable_audit_logging: true,
  },
  notifications: {
    email_enabled: true,
    email_smtp_host: '',
    email_smtp_port: 587,
    push_enabled: false,
    webhook_url: '',
    notify_new_reports: true,
  },
}

export default function SettingsPage() {
  const { data, isLoading, isFetching, isError, error } = useSettings()
  const updateSettingsMutation = useUpdateSettings()

  const [settings, setSettings] = useState<SystemSettings | null>(null)

  useEffect(() => {
    if (data?.settings) {
      setSettings(data.settings as SystemSettings)
    }
  }, [data?.settings])

  useEffect(() => {
    if (!isError || !error) {
      return
    }
    toast.error('Erro ao carregar configurações')
    // eslint-disable-next-line no-console
    console.error(error)
  }, [isError, error])

  const isSaving = updateSettingsMutation.isPending
  const isBusy = isLoading || isFetching || !settings

  const dirtyStateHash = useMemo(() => JSON.stringify(settings), [settings])
  const originalStateHash = useMemo(() => JSON.stringify(data?.settings), [data?.settings])
  const isDirty = dirtyStateHash !== originalStateHash

  const handleSave = async () => {
    if (!settings) {
      return
    }

    try {
      await updateSettingsMutation.mutateAsync(settings)
    } catch (mutationError) {
      // eslint-disable-next-line no-console
      console.error(mutationError)
    }
  }

  const handleReset = () => {
    if (!data?.settings) {
      setSettings(DEFAULT_SETTINGS)
      return
    }

    setSettings(data.settings as SystemSettings)
    toast.success('Configurações restauradas')
  }

  const updateSetting = (section: keyof SystemSettings, key: string, value: any) => {
    setSettings((prev) => {
      if (!prev) {
        return prev
      }
      return {
        ...prev,
        [section]: {
          ...prev[section],
          [key]: value,
        },
      }
    })
  }

  if (isBusy) {
    return (
      <div className="flex h-full items-center justify-center">
        <div className="flex items-center gap-2 text-muted-foreground">
          <Loader2 className="h-5 w-5 animate-spin" /> Carregando configurações...
        </div>
      </div>
    )
  }

  return (
    <div className="space-y-6">
      <div className="flex justify-between items-center">
        <h1 className="text-3xl font-bold">Configurações do Sistema</h1>
        <div className="flex gap-2">
          <Button variant="outline" onClick={handleReset} disabled={isSaving}>
            <RefreshCw className="mr-2 h-4 w-4" />
            Resetar
          </Button>
          <Button onClick={handleSave} disabled={isSaving || !isDirty}>
            {isSaving ? (
              <>
                <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                Salvando...
              </>
            ) : (
              <>
                <Save className="mr-2 h-4 w-4" />
                Salvar
              </>
            )}
          </Button>
        </div>
      </div>

      {/* General Settings */}
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center">
            <Settings className="mr-2 h-5 w-5" />
            Configurações Gerais
          </CardTitle>
        </CardHeader>
        <CardContent className="space-y-4">
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div className="space-y-2">
              <Label htmlFor="site_name">Nome do Site</Label>
              <Input
                id="site_name"
                value={settings?.general.site_name ?? ''}
                onChange={(e) => updateSetting('general', 'site_name', e.target.value)}
              />
            </div>
            <div className="space-y-2">
              <Label htmlFor="timezone">Fuso Horário</Label>
              <select
                id="timezone"
                value={settings?.general.timezone ?? ''}
                onChange={(e) => updateSetting('general', 'timezone', e.target.value)}
                className="w-full p-2 border rounded-md"
              >
                <option value="America/Sao_Paulo">America/São Paulo</option>
                <option value="America/New_York">America/New York</option>
                <option value="Europe/London">Europe/London</option>
                <option value="UTC">UTC</option>
              </select>
            </div>
          </div>
          <div className="space-y-2">
            <Label htmlFor="site_description">Descrição</Label>
            <Textarea
              id="site_description"
              value={settings?.general.site_description ?? ''}
              onChange={(e) => updateSetting('general', 'site_description', e.target.value)}
              placeholder="Descrição exibida na interface administrativa"
            />
          </div>
          <div className="flex items-center justify-between rounded border p-3">
            <div>
              <p className="font-medium">Modo de Manutenção</p>
              <p className="text-sm text-muted-foreground">
                Controla o acesso público à aplicação.
              </p>
            </div>
            <Switch
              checked={settings?.general.maintenance_mode ?? false}
              onCheckedChange={(checked) => updateSetting('general', 'maintenance_mode', checked)}
            />
          </div>
        </CardContent>
      </Card>

      {/* Map Settings */}
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center">
            <MapPin className="mr-2 h-5 w-5" />
            Configurações de Mapa
          </CardTitle>
        </CardHeader>
        <CardContent className="space-y-4">
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div className="space-y-2">
              <Label htmlFor="map_style">Estilo do Mapbox</Label>
              <Input
                id="map_style"
                value={settings?.map.default_style ?? ''}
                onChange={(e) => updateSetting('map', 'default_style', e.target.value)}
              />
            </div>
            <div className="space-y-2">
              <Label htmlFor="default_zoom">Zoom Padrão</Label>
              <Input
                id="default_zoom"
                type="number"
                min={0}
                max={22}
                step={0.1}
                value={settings?.map.default_zoom ?? 0}
                onChange={(e) => updateSetting('map', 'default_zoom', Number(e.target.value))}
              />
            </div>
            <div className="space-y-2">
              <Label htmlFor="center_lat">Latitude Padrão</Label>
              <Input
                id="center_lat"
                type="number"
                step={0.000001}
                value={settings?.map.default_center_lat ?? 0}
                onChange={(e) => updateSetting('map', 'default_center_lat', Number(e.target.value))}
              />
            </div>
            <div className="space-y-2">
              <Label htmlFor="center_lng">Longitude Padrão</Label>
              <Input
                id="center_lng"
                type="number"
                step={0.000001}
                value={settings?.map.default_center_lng ?? 0}
                onChange={(e) => updateSetting('map', 'default_center_lng', Number(e.target.value))}
              />
            </div>
          </div>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div className="space-y-2">
              <Label htmlFor="min_zoom">Zoom Mínimo</Label>
              <Input
                id="min_zoom"
                type="number"
                min={0}
                max={22}
                step={0.1}
                value={settings?.map.min_zoom ?? 0}
                onChange={(e) => updateSetting('map', 'min_zoom', Number(e.target.value))}
              />
            </div>
            <div className="space-y-2">
              <Label htmlFor="max_zoom">Zoom Máximo</Label>
              <Input
                id="max_zoom"
                type="number"
                min={0}
                max={22}
                step={0.1}
                value={settings?.map.max_zoom ?? 0}
                onChange={(e) => updateSetting('map', 'max_zoom', Number(e.target.value))}
              />
            </div>
          </div>
        </CardContent>
      </Card>

      {/* Security Settings */}
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center">
            <Shield className="mr-2 h-5 w-5" />
            Segurança
          </CardTitle>
        </CardHeader>
        <CardContent className="space-y-4">
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div className="space-y-2">
              <Label htmlFor="session_timeout">Expiração de Sessão (min)</Label>
              <Input
                id="session_timeout"
                type="number"
                min={5}
                max={1440}
                value={settings?.security.session_timeout_minutes ?? 0}
                onChange={(e) =>
                  updateSetting('security', 'session_timeout_minutes', Number(e.target.value))
                }
              />
            </div>
            <div className="space-y-2">
              <Label htmlFor="password_length">Tamanho mínimo da senha</Label>
              <Input
                id="password_length"
                type="number"
                min={6}
                max={64}
                value={settings?.security.password_min_length ?? 0}
                onChange={(e) =>
                  updateSetting('security', 'password_min_length', Number(e.target.value))
                }
              />
            </div>
          </div>
          <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
            <div className="flex items-center justify-between rounded border p-3">
              <div>
                <p className="font-medium">Verificação por email</p>
                <p className="text-sm text-muted-foreground">Obrigar confirmação de email no cadastro.</p>
              </div>
              <Switch
                checked={settings?.security.require_email_verification ?? false}
                onCheckedChange={(checked) =>
                  updateSetting('security', 'require_email_verification', checked)
                }
              />
            </div>
            <div className="flex items-center justify-between rounded border p-3">
              <div>
                <p className="font-medium">Autenticação em duas etapas</p>
                <p className="text-sm text-muted-foreground">Exigir 2FA para usuários privilegiados.</p>
              </div>
              <Switch
                checked={settings?.security.enable_two_factor ?? false}
                onCheckedChange={(checked) => updateSetting('security', 'enable_two_factor', checked)}
              />
            </div>
            <div className="flex items-center justify-between rounded border p-3">
              <div>
                <p className="font-medium">Registro de auditoria</p>
                <p className="text-sm text-muted-foreground">Salvar operações críticas no log.</p>
              </div>
              <Switch
                checked={settings?.security.enable_audit_logging ?? false}
                onCheckedChange={(checked) => updateSetting('security', 'enable_audit_logging', checked)}
              />
            </div>
          </div>
        </CardContent>
      </Card>

      {/* Notifications */}
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center">
            <Bell className="mr-2 h-5 w-5" />
            Notificações
          </CardTitle>
        </CardHeader>
        <CardContent className="space-y-4">
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div className="flex items-center justify-between rounded border p-3">
              <div>
                <p className="font-medium">Notificações por Email</p>
                <p className="text-sm text-muted-foreground">Enviar alertas por email para a equipe.</p>
              </div>
              <Switch
                checked={settings?.notifications.email_enabled ?? false}
                onCheckedChange={(checked) => updateSetting('notifications', 'email_enabled', checked)}
              />
            </div>
            <div className="flex items-center justify-between rounded border p-3">
              <div>
                <p className="font-medium">Notificações Push</p>
                <p className="text-sm text-muted-foreground">Habilitar notificações push no app.</p>
              </div>
              <Switch
                checked={settings?.notifications.push_enabled ?? false}
                onCheckedChange={(checked) => updateSetting('notifications', 'push_enabled', checked)}
              />
            </div>
          </div>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div className="space-y-2">
              <Label htmlFor="smtp_host">Servidor SMTP</Label>
              <Input
                id="smtp_host"
                value={settings?.notifications.email_smtp_host ?? ''}
                onChange={(e) => updateSetting('notifications', 'email_smtp_host', e.target.value)}
                placeholder="smtp.exemplo.com"
              />
            </div>
            <div className="space-y-2">
              <Label htmlFor="smtp_port">Porta SMTP</Label>
              <Input
                id="smtp_port"
                type="number"
                min={1}
                max={65535}
                value={settings?.notifications.email_smtp_port ?? 0}
                onChange={(e) =>
                  updateSetting('notifications', 'email_smtp_port', Number(e.target.value))
                }
              />
            </div>
          </div>
          <div className="space-y-2">
            <Label htmlFor="webhook_url">Webhook</Label>
            <Input
              id="webhook_url"
              value={settings?.notifications.webhook_url ?? ''}
              onChange={(e) => updateSetting('notifications', 'webhook_url', e.target.value)}
              placeholder="https://hooks.exemplo.com"
            />
          </div>
          <div className="flex items-center justify-between rounded border p-3">
            <div>
              <p className="font-medium">Avisar novo relatório</p>
              <p className="text-sm text-muted-foreground">Enviar alerta quando um novo relatório for criado.</p>
            </div>
            <Switch
              checked={settings?.notifications.notify_new_reports ?? false}
              onCheckedChange={(checked) =>
                updateSetting('notifications', 'notify_new_reports', checked)
              }
            />
          </div>
        </CardContent>
      </Card>
    </div>
  )
}
