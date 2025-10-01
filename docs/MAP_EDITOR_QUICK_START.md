# Map Editor Quick Start Guide

## Overview
The Map Editor is now fully integrated into the admin interface at `/maps`. This guide will help you get started quickly.

## Prerequisites

### 1. Install Dependencies
```bash
cd admin
npm install
```

### 2. Configure Mapbox Token
Create or update `admin/.env.local`:
```bash
NEXT_PUBLIC_MAPBOX_TOKEN=pk.your-actual-mapbox-token
```

Get your token from: https://account.mapbox.com/access-tokens/

## Using the Map Editor

### Navigation
Go to: **http://localhost:3000/maps**

### Interface Layout

```
┌─────────────────────────────────────────────────────────────┐
│ Editor de Mapas                             [Filtros]       │
├─────────────────────────────────────────────────────────────┤
│ Statistics Cards: [Total Áreas] [Facções] [Cobertura]      │
├───────────────────────────────────┬─────────────────────────┤
│                                   │  Ferramentas            │
│                                   │  ┌────────────────────┐ │
│                                   │  │ [📍] Selecionar    │ │
│      INTERACTIVE MAP              │  │ [⬜] Polígono       │ │
│      (Mapbox GL JS)               │  │ [─] Linha          │ │
│                                   │  │ [✏️] Editar         │ │
│      - Navigation Controls        │  └────────────────────┘ │
│      - Zoom/Pan/Rotate            │                         │
│      - Faction Layers             │  Camadas por Facção     │
│      - Drawing Tools              │  ┌────────────────────┐ │
│                                   │  │ 🔴 TDN      [ON]   │ │
│                                   │  │ 🔵 CV       [ON]   │ │
│                                   │  │ 🟡 PCC      [ON]   │ │
│                                   │  └────────────────────┘ │
│                                   │                         │
│                                   │  Propriedades           │
│                                   │  [Geometry Details]     │
│                                   │                         │
│                                   │  Estatísticas           │
│                                   │  [Quick Stats]          │
└───────────────────────────────────┴─────────────────────────┘
```

### Drawing Tools

#### 1. Select Tool (📍)
- **Purpose**: Select and view existing geometries
- **Usage**: Click on any colored area on the map
- **Result**: Properties panel shows geometry details

#### 2. Polygon Tool (⬜)
- **Purpose**: Draw new polygon areas
- **Usage**:
  1. Click tool to activate
  2. Click on map to add vertices
  3. Double-click to complete polygon
- **Result**: Geometry saved to backend automatically

#### 3. Line Tool (─)
- **Purpose**: Draw linear features (boundaries, routes)
- **Usage**:
  1. Click tool to activate
  2. Click on map to add points
  3. Double-click to complete line
- **Result**: Line saved to backend

#### 4. Edit Tool (✏️)
- **Purpose**: Modify existing geometries
- **Usage**:
  1. Click tool to activate
  2. Click on a geometry to select
  3. Drag vertices to reshape
- **Result**: Changes saved automatically

### Layer Controls

#### Toggle Faction Visibility
- Each faction has a colored layer
- Use switches to show/hide specific factions
- Useful for focusing on specific areas

Example:
```
🔴 TDN  [ON]  → Click to hide TDN geometries
🔵 CV   [OFF] → Click to show CV geometries
```

### Viewing Geometry Details

When you click on a geometry, a popup shows:
- **Faction Name**: Which faction controls this area
- **Region Name**: Name of the region/neighborhood
- **Risk Level**: Security risk rating (1-5)

### Workflow Example

#### Creating a New Area
1. Select the **Polygon Tool** from the toolbar
2. Click on the map to outline the area
3. Double-click to finish
4. Geometry is automatically saved
5. Toast notification confirms success

#### Editing an Existing Area
1. Select the **Edit Tool** from the toolbar
2. Click on the geometry you want to modify
3. Drag the vertices to reshape the area
4. Changes save automatically
5. Toast notification confirms update

## Features

### ✅ Currently Working
- Interactive map with Mapbox GL JS
- Draw polygons and lines
- Edit existing geometries
- Delete geometries (via edit mode + delete key)
- Layer visibility controls
- Faction-based color coding
- Real-time updates via API
- Popup information on click
- Navigation controls (zoom, pan, rotate)
- Scale indicator

### 🚧 Coming Soon
- Faction selection when creating
- Advanced property editing
- Search/geocoding
- GeoJSON import/export
- Measurement tools
- Undo/redo

## Keyboard Shortcuts

- **Escape**: Cancel current drawing
- **Delete**: Remove selected geometry (in edit mode)
- **+/-**: Zoom in/out
- **Arrow Keys**: Pan the map

## Map Controls

### Navigation Controls (Top Right)
- **+** button: Zoom in
- **-** button: Zoom out
- **⟳** button: Reset north
- **3D** button: Toggle tilt

### Scale Control (Bottom Left)
- Shows map scale in kilometers

## Troubleshooting

### Map Not Loading
**Problem**: Gray box instead of map
**Solution**: Check that `NEXT_PUBLIC_MAPBOX_TOKEN` is set correctly in `.env.local`

### Can't Draw
**Problem**: Nothing happens when clicking
**Solution**: 
1. Ensure correct tool is selected
2. Wait for map to fully load (check for "Mapa carregado" toast)
3. Check browser console for errors

### Geometries Not Showing
**Problem**: No colored areas visible
**Solution**:
1. Check layer visibility switches (ensure they're ON)
2. Verify backend API is running
3. Check that geometries exist in database

## API Integration

The Map Editor automatically handles:
- **Create**: POST /api/v1/geometries
- **Update**: PUT /api/v1/geometries/:id
- **Delete**: DELETE /api/v1/geometries/:id
- **Load**: GET /api/v1/geometries

All operations include:
- Loading states
- Success/error notifications
- Automatic data refresh

## Data Format

Geometries are stored in GeoJSON format:
```json
{
  "geometry": {
    "type": "Polygon",
    "coordinates": [[[lng, lat], [lng, lat], ...]]
  },
  "faction_id": "uuid",
  "region_name": "Centro",
  "risk_level": 3,
  "active": true
}
```

## Performance Tips

1. **Toggle Off Unused Layers**: Hide factions you're not working with
2. **Zoom In**: Better performance at higher zoom levels
3. **Limit Vertices**: Keep polygons simple for better performance
4. **Regular Saves**: Geometries auto-save, but check for success toasts

## Browser Compatibility

### Recommended Browsers
- ✅ Chrome 80+
- ✅ Firefox 78+
- ✅ Edge 80+
- ✅ Safari 13+

### Requirements
- WebGL support
- JavaScript enabled
- Modern browser (2020+)

## Development Mode

To run the admin interface locally:

```bash
cd admin
npm run dev
```

Access at: http://localhost:3000

## Production Build

To build for production:

```bash
cd admin
npm run build
npm start
```

## Next Steps

1. **Get a Mapbox Token**: https://account.mapbox.com/
2. **Configure .env.local**: Add your token
3. **Start the App**: `npm run dev`
4. **Navigate to /maps**: Start drawing!

## Support

For issues or questions:
- Check documentation in `docs/MAP_EDITOR_IMPLEMENTATION.md`
- Review Mapbox docs: https://docs.mapbox.com/
- Check console for error messages

## Summary

The Map Editor is fully functional and ready to use. The main requirement is a valid Mapbox token. Once configured, you can:
- ✅ Draw new areas
- ✅ Edit existing areas
- ✅ Manage faction layers
- ✅ View area details
- ✅ All changes save automatically

Happy mapping! 🗺️
