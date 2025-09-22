'use client'

import { useState } from 'react'
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
  RefreshCw
} from 'lucide-react'
import { toast } from 'sonner'

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

export default function SettingsPage() {
  const [isLoading, setIsLoading] = useState(false)
  const [settings, setSettings] = useState<SystemSettings>({
    general: {
      site_name: 'Map Area Factions',
      site_description: 'Sistema de mapeamento de áreas por facção',
      maintenance_mode: false,
      timezone: 'America/Sao_Paulo',
      default_language: 'pt-BR'
    },
    map: {
      default_style: 'mapbox://styles/mapbox/streets-v11',
      default_zoom: 10,
      default_center_lat: -3.7319,
      default_center_lng: -38.5267,
      max_zoom: 18,
      min_zoom: 8
    },
    security: {
      session_timeout_minutes: 480,
      require_email_verification: true,
      enable_two_factor: false,
      password_min_length: 8,
      enable_audit_logging: true
    },
    notifications: {
      email_enabled: true,
      email_smtp_host: 'smtp.gmail.com',
      email_smtp_port: 587,
      push_enabled: false,
      webhook_url: '',
      notify_new_reports: true
    }
  })

  const handleSave = async () => {
    setIsLoading(true)
    try {
      // Simulate API call
      await new Promise(resolve => setTimeout(resolve, 1000))
      toast.success('Configurações salvas com sucesso')
    } catch (error) {
      toast.error('Erro ao salvar configurações')
    } finally {
      setIsLoading(false)
    }
  }

  const handleReset = () => {
    if (confirm('Tem certeza que deseja resetar todas as configurações?')) {
      // Reset to default values
      setSettings({
        general: {
          site_name: 'Map Area Factions',
          site_description: 'Sistema de mapeamento de áreas por facção',
          maintenance_mode: false,
          timezone: 'America/Sao_Paulo',
          default_language: 'pt-BR'
        },
        map: {
          default_style: 'mapbox://styles/mapbox/streets-v11',
          default_zoom: 10,
          default_center_lat: -3.7319,
          default_center_lng: -38.5267,
          max_zoom: 18,
          min_zoom: 8
        },
        security: {
          session_timeout_minutes: 480,
          require_email_verification: true,
          enable_two_factor: false,
          password_min_length: 8,
          enable_audit_logging: true
        },
        notifications: {
          email_enabled: true,
          email_smtp_host: 'smtp.gmail.com',
          email_smtp_port: 587,
          push_enabled: false,
          webhook_url: '',
          notify_new_reports: true
        }
      })
      toast.success('Configurações resetadas')
    }
  }

  const updateSetting = (section: keyof SystemSettings, key: string, value: any) => {
    setSettings(prev => ({
      ...prev,
      [section]: {
        ...prev[section],
        [key]: value
      }
    }))
  }

  return (
    <div className="space-y-6">
      <div className="flex justify-between items-center">
        <h1 className="text-3xl font-bold">Configurações do Sistema</h1>
        <div className="flex gap-2">
          <Button variant="outline" onClick={handleReset}>
            <RefreshCw className="mr-2 h-4 w-4" />
            Resetar
          </Button>
          <Button onClick={handleSave} disabled={isLoading}>
            <Save className="mr-2 h-4 w-4" />
            {isLoading ? 'Salvando...' : 'Salvar'}
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
                value={settings.general.site_name}
                onChange={(e) => updateSetting('general', 'site_name', e.target.value)}
              />
            </div>
            <div className="space-y-2">
              <Label htmlFor="timezone">Fuso Horário</Label>
              <select
                id="timezone"
                value={settings.general.timezone}
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
            <Label htmlFor="site_description">Descrição do Site</Label>
            <Textarea
              id="site_description"
              value={settings.general.site_description}
              onChange={(e) => updateSetting('general', 'site_description', e.target.value)}
              rows={3}
            />
          </div>

          <div className="flex items-center space-x-2">
            <Switch
              id="maintenance_mode"
              checked={settings.general.maintenance_mode}
              onCheckedChange={(checked) => updateSetting('general', 'maintenance_mode', checked)}
            />
            <Label htmlFor="maintenance_mode">Modo de Manutenção</Label>
          </div>
        </CardContent>
      </Card>

      {/* Map Settings */}
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center">
            <MapPin className="mr-2 h-5 w-5" />
            Configurações do Mapa
          </CardTitle>
        </CardHeader>
        <CardContent className="space-y-4">
          <div className="space-y-2">
            <Label htmlFor="map_style">Estilo Padrão do Mapa</Label>
            <select
              id="map_style"
              value={settings.map.default_style}
              onChange={(e) => updateSetting('map', 'default_style', e.target.value)}
              className="w-full p-2 border rounded-md"
            >
              <option value="mapbox://styles/mapbox/streets-v11">Streets</option>
              <option value="mapbox://styles/mapbox/satellite-v9">Satellite</option>
              <option value="mapbox://styles/mapbox/light-v10">Light</option>
              <option value="mapbox://styles/mapbox/dark-v10">Dark</option>
            </select>
          </div>

          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div className="space-y-2">
              <Label htmlFor="default_center_lat">Latitude Central Padrão</Label>
              <Input
                id="default_center_lat"
                type="number"
                step="0.000001"
                value={settings.map.default_center_lat}
                onChange={(e) => updateSetting('map', 'default_center_lat', parseFloat(e.target.value))}
              />
            </div>
            <div className="space-y-2">
              <Label htmlFor="default_center_lng">Longitude Central Padrão</Label>
              <Input
                id="default_center_lng"
                type="number"
                step="0.000001"
                value={settings.map.default_center_lng}
                onChange={(e) => updateSetting('map', 'default_center_lng', parseFloat(e.target.value))}
              />
            </div>
          </div>

          <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
            <div className="space-y-2">
              <Label htmlFor="default_zoom">Zoom Padrão</Label>
              <Input
                id="default_zoom"
                type="number"
                min="1"
                max="20"
                value={settings.map.default_zoom}
                onChange={(e) => updateSetting('map', 'default_zoom', parseInt(e.target.value))}
              />
            </div>
            <div className="space-y-2">
              <Label htmlFor="min_zoom">Zoom Mínimo</Label>
              <Input
                id="min_zoom"
                type="number"
                min="1"
                max="20"
                value={settings.map.min_zoom}
                onChange={(e) => updateSetting('map', 'min_zoom', parseInt(e.target.value))}
              />
            </div>
            <div className="space-y-2">
              <Label htmlFor="max_zoom">Zoom Máximo</Label>
              <Input
                id="max_zoom"
                type="number"
                min="1"
                max="20"
                value={settings.map.max_zoom}
                onChange={(e) => updateSetting('map', 'max_zoom', parseInt(e.target.value))}
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
            Configurações de Segurança
          </CardTitle>
        </CardHeader>
        <CardContent className="space-y-4">
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div className="space-y-2">
              <Label htmlFor="session_timeout">Timeout da Sessão (minutos)</Label>
              <Input
                id="session_timeout"
                type="number"
                min="30"
                max="1440"
                value={settings.security.session_timeout_minutes}
                onChange={(e) => updateSetting('security', 'session_timeout_minutes', parseInt(e.target.value))}
              />
            </div>
            <div className="space-y-2">
              <Label htmlFor="password_min_length">Tamanho Mínimo da Senha</Label>
              <Input
                id="password_min_length"
                type="number"
                min="6"
                max="20"
                value={settings.security.password_min_length}
                onChange={(e) => updateSetting('security', 'password_min_length', parseInt(e.target.value))}
              />
            </div>
          </div>

          <div className="space-y-3">
            <div className="flex items-center space-x-2">
              <Switch
                id="email_verification"
                checked={settings.security.require_email_verification}
                onCheckedChange={(checked) => updateSetting('security', 'require_email_verification', checked)}
              />
              <Label htmlFor="email_verification">Exigir Verificação de Email</Label>
            </div>
            
            <div className="flex items-center space-x-2">
              <Switch
                id="two_factor"
                checked={settings.security.enable_two_factor}
                onCheckedChange={(checked) => updateSetting('security', 'enable_two_factor', checked)}
              />
              <Label htmlFor="two_factor">Habilitar Autenticação de Dois Fatores</Label>
            </div>
            
            <div className="flex items-center space-x-2">
              <Switch
                id="audit_logging"
                checked={settings.security.enable_audit_logging}
                onCheckedChange={(checked) => updateSetting('security', 'enable_audit_logging', checked)}
              />
              <Label htmlFor="audit_logging">Habilitar Log de Auditoria</Label>
            </div>
          </div>
        </CardContent>
      </Card>

      {/* Notification Settings */}
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center">
            <Bell className="mr-2 h-5 w-5" />
            Configurações de Notificação
          </CardTitle>
        </CardHeader>
        <CardContent className="space-y-4">
          <div className="flex items-center space-x-2 mb-4">
            <Switch
              id="email_enabled"
              checked={settings.notifications.email_enabled}
              onCheckedChange={(checked) => updateSetting('notifications', 'email_enabled', checked)}
            />
            <Label htmlFor="email_enabled">Habilitar Notificações por Email</Label>
          </div>

          {settings.notifications.email_enabled && (
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div className="space-y-2">
                <Label htmlFor="smtp_host">Servidor SMTP</Label>
                <Input
                  id="smtp_host"
                  value={settings.notifications.email_smtp_host}
                  onChange={(e) => updateSetting('notifications', 'email_smtp_host', e.target.value)}
                />
              </div>
              <div className="space-y-2">
                <Label htmlFor="smtp_port">Porta SMTP</Label>
                <Input
                  id="smtp_port"
                  type="number"
                  value={settings.notifications.email_smtp_port}
                  onChange={(e) => updateSetting('notifications', 'email_smtp_port', parseInt(e.target.value))}
                />
              </div>
            </div>
          )}

          <div className="space-y-3">
            <div className="flex items-center space-x-2">
              <Switch
                id="push_enabled"
                checked={settings.notifications.push_enabled}
                onCheckedChange={(checked) => updateSetting('notifications', 'push_enabled', checked)}
              />
              <Label htmlFor="push_enabled">Habilitar Notificações Push</Label>
            </div>
            
            <div className="flex items-center space-x-2">
              <Switch
                id="notify_new_reports"
                checked={settings.notifications.notify_new_reports}
                onCheckedChange={(checked) => updateSetting('notifications', 'notify_new_reports', checked)}
              />
              <Label htmlFor="notify_new_reports">Notificar Novos Relatórios</Label>
            </div>
          </div>

          <div className="space-y-2">
            <Label htmlFor="webhook_url">URL do Webhook (opcional)</Label>
            <Input
              id="webhook_url"
              type="url"
              placeholder="https://api.exemplo.com/webhook"
              value={settings.notifications.webhook_url}
              onChange={(e) => updateSetting('notifications', 'webhook_url', e.target.value)}
            />
          </div>
        </CardContent>
      </Card>

      {/* System Status */}
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center">
            <Database className="mr-2 h-5 w-5" />
            Status do Sistema
          </CardTitle>
        </CardHeader>
        <CardContent>
          <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
            <div className="bg-green-50 p-4 rounded-lg">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-sm font-medium text-green-600">Banco de Dados</p>
                  <p className="text-2xl font-bold text-green-800">Online</p>
                </div>
                <div className="w-3 h-3 bg-green-500 rounded-full"></div>
              </div>
            </div>
            
            <div className="bg-green-50 p-4 rounded-lg">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-sm font-medium text-green-600">Redis Cache</p>
                  <p className="text-2xl font-bold text-green-800">Online</p>
                </div>
                <div className="w-3 h-3 bg-green-500 rounded-full"></div>
              </div>
            </div>
            
            <div className="bg-blue-50 p-4 rounded-lg">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-sm font-medium text-blue-600">WebSocket</p>
                  <p className="text-2xl font-bold text-blue-800">Ativo</p>
                </div>
                <div className="w-3 h-3 bg-blue-500 rounded-full"></div>
              </div>
            </div>
          </div>
        </CardContent>
      </Card>
    </div>
  )
}