'use client'

import { useEffect, useRef, useState } from 'react'
import mapboxgl from 'mapbox-gl'
import MapboxDraw from '@mapbox/mapbox-gl-draw'
import { toast } from 'sonner'
import { useMapboxToken } from '@/hooks/useMapboxToken'
import 'mapbox-gl/dist/mapbox-gl.css'
import '@mapbox/mapbox-gl-draw/dist/mapbox-gl-draw.css'

interface MapEditorProps {
  center?: [number, number]
  zoom?: number
  onGeometryCreate?: (geometry: any) => void
  onGeometryUpdate?: (geometry: any) => void
  onGeometryDelete?: (id: string) => void
  selectedTool?: 'select' | 'polygon' | 'line' | 'edit'
  factions?: any[]
  geometries?: any[]
  layerVisibility?: { [key: string]: boolean }
}

export function MapEditor({
  center = [-38.5267, -3.7172], // Fortaleza, CE
  zoom = 10,
  onGeometryCreate,
  onGeometryUpdate,
  onGeometryDelete,
  selectedTool = 'select',
  factions = [],
  geometries = [],
  layerVisibility = {}
}: MapEditorProps) {
  const mapContainer = useRef<HTMLDivElement>(null)
  const map = useRef<mapboxgl.Map | null>(null)
  const draw = useRef<MapboxDraw | null>(null)
  const [mapLoaded, setMapLoaded] = useState(false)
  const { token: mapboxToken, isLoading: tokenLoading } = useMapboxToken()

  // Initialize map
  useEffect(() => {
    if (!mapContainer.current || map.current || tokenLoading) return

    if (!mapboxToken) {
      console.warn('Mapbox token not configured')
      return
    }

    console.log('Mapbox token configured:', mapboxToken.substring(0, 20) + '...')
    mapboxgl.accessToken = mapboxToken

    try {
      const mapInstance = new mapboxgl.Map({
        container: mapContainer.current,
        style: 'mapbox://styles/mapbox/streets-v12',
        center: center,
        zoom: zoom,
        attributionControl: true
      })

      // Add navigation controls
      mapInstance.addControl(new mapboxgl.NavigationControl(), 'top-right')

      // Add scale control
      mapInstance.addControl(new mapboxgl.ScaleControl(), 'bottom-left')

      // Initialize MapboxDraw
      const drawInstance = new MapboxDraw({
        displayControlsDefault: false,
        controls: {
          polygon: false,
          line_string: false,
          point: false,
          trash: false
        },
        defaultMode: 'simple_select'
      })

      mapInstance.addControl(drawInstance, 'top-left')

      // Map loaded event
      mapInstance.on('load', () => {
        setMapLoaded(true)
        toast.success('Mapa carregado com sucesso')
      })

      // Draw events
      mapInstance.on('draw.create', (e: any) => {
        if (onGeometryCreate && e.features && e.features.length > 0) {
          const feature = e.features[0]
          onGeometryCreate(feature)
        }
      })

      mapInstance.on('draw.update', (e: any) => {
        if (onGeometryUpdate && e.features && e.features.length > 0) {
          const feature = e.features[0]
          onGeometryUpdate(feature)
        }
      })

      mapInstance.on('draw.delete', (e: any) => {
        if (onGeometryDelete && e.features && e.features.length > 0) {
          const feature = e.features[0]
          onGeometryDelete(feature.id as string)
        }
      })

      map.current = mapInstance
      draw.current = drawInstance

    } catch (error) {
      console.error('Error initializing map:', error)
      toast.error('Erro ao carregar o mapa')
    }

    // Cleanup
    return () => {
      if (map.current) {
        map.current.remove()
        map.current = null
        draw.current = null
      }
    }
  }, [mapboxToken, tokenLoading]) // Run when token changes

  // Handle tool changes
  useEffect(() => {
    if (!draw.current || !mapLoaded) return

    try {
      switch (selectedTool) {
        case 'polygon':
          (draw.current as any).changeMode('draw_polygon')
          break
        case 'line':
          (draw.current as any).changeMode('draw_line_string')
          break
        case 'edit':
          (draw.current as any).changeMode('direct_select')
          break
        case 'select':
        default:
          (draw.current as any).changeMode('simple_select')
          break
      }
    } catch (error) {
      console.error('Error changing draw mode:', error)
    }
  }, [selectedTool, mapLoaded])

  // Load existing geometries
  useEffect(() => {
    if (!map.current || !mapLoaded) return

    // Remove existing sources and layers
    factions.forEach((faction: any) => {
      const layerId = `faction-${faction.id}`
      const sourceId = `faction-${faction.id}-source`

      if (map.current?.getLayer(layerId)) {
        map.current.removeLayer(layerId)
      }
      if (map.current?.getSource(sourceId)) {
        map.current.removeSource(sourceId)
      }
    })

    // Add layers for each faction
    factions.forEach((faction: any) => {
      const factionGeometries = geometries.filter(
        (geo: any) => geo.faction_id === faction.id
      )

      if (factionGeometries.length === 0) return

      const layerId = `faction-${faction.id}`
      const sourceId = `faction-${faction.id}-source`

      const features = factionGeometries
        .filter((geo: any) => geo.geometry)
        .map((geo: any) => ({
          type: 'Feature' as const,
          id: geo.id,
          geometry: typeof geo.geometry === 'string' 
            ? JSON.parse(geo.geometry) 
            : geo.geometry,
          properties: {
            id: geo.id,
            faction_id: faction.id,
            faction_name: faction.name,
            faction_color: faction.color_hex,
            risk_level: geo.risk_level || 1,
            region_name: geo.region_name || 'Indefinida'
          }
        }))

      if (features.length === 0) return

      const geojsonData = {
        type: 'FeatureCollection' as const,
        features: features
      }

      try {
        map.current?.addSource(sourceId, {
          type: 'geojson',
          data: geojsonData
        })

        // Add fill layer
        map.current?.addLayer({
          id: layerId,
          type: 'fill',
          source: sourceId,
          layout: {
            visibility: layerVisibility[faction.id] === false ? 'none' : 'visible'
          },
          paint: {
            'fill-color': faction.color_hex || '#888888',
            'fill-opacity': 0.5
          }
        })

        // Add outline layer
        map.current?.addLayer({
          id: `${layerId}-outline`,
          type: 'line',
          source: sourceId,
          layout: {
            visibility: layerVisibility[faction.id] === false ? 'none' : 'visible'
          },
          paint: {
            'line-color': faction.color_hex || '#888888',
            'line-width': 2
          }
        })

        // Add click popup
        map.current?.on('click', layerId, (e: any) => {
          if (!e.features || e.features.length === 0) return

          const feature = e.features[0]
          const props = feature.properties

          new mapboxgl.Popup()
            .setLngLat(e.lngLat)
            .setHTML(`
              <div style="padding: 8px;">
                <h3 style="margin: 0 0 4px 0; font-weight: bold;">${props.faction_name}</h3>
                <p style="margin: 2px 0;"><strong>Região:</strong> ${props.region_name}</p>
                <p style="margin: 2px 0;"><strong>Nível de Risco:</strong> ${props.risk_level}/5</p>
              </div>
            `)
            .addTo(map.current!)
        })

        // Change cursor on hover
        map.current?.on('mouseenter', layerId, () => {
          if (map.current) {
            map.current.getCanvas().style.cursor = 'pointer'
          }
        })

        map.current?.on('mouseleave', layerId, () => {
          if (map.current) {
            map.current.getCanvas().style.cursor = ''
          }
        })
      } catch (error) {
        console.error(`Error adding layer for faction ${faction.id}:`, error)
      }
    })
  }, [geometries, factions, mapLoaded, layerVisibility])

  // Update layer visibility
  useEffect(() => {
    if (!map.current || !mapLoaded) return

    factions.forEach((faction: any) => {
      const layerId = `faction-${faction.id}`
      const outlineLayerId = `${layerId}-outline`

      if (map.current?.getLayer(layerId)) {
        map.current.setLayoutProperty(
          layerId,
          'visibility',
          layerVisibility[faction.id] === false ? 'none' : 'visible'
        )
      }

      if (map.current?.getLayer(outlineLayerId)) {
        map.current.setLayoutProperty(
          outlineLayerId,
          'visibility',
          layerVisibility[faction.id] === false ? 'none' : 'visible'
        )
      }
    })
  }, [layerVisibility, factions, mapLoaded])

  if (tokenLoading) {
    return (
      <div className="w-full h-[500px] rounded-lg flex items-center justify-center bg-gray-100">
        <div className="text-center">
          <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600 mx-auto mb-2"></div>
          <p className="text-sm text-gray-600">Carregando mapa...</p>
        </div>
      </div>
    )
  }

  if (!mapboxToken) {
    return (
      <div className="w-full h-[500px] rounded-lg flex items-center justify-center bg-red-50">
        <div className="text-center">
          <div className="text-red-500 mb-2">⚠️</div>
          <p className="text-sm text-red-600">Token do Mapbox não configurado</p>
        </div>
      </div>
    )
  }

  return (
    <div 
      ref={mapContainer} 
      className="w-full h-[500px] rounded-lg"
      style={{ minHeight: '500px' }}
    />
  )
}
