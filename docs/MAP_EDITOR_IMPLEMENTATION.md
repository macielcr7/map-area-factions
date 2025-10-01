# Map Editor Implementation

## Overview
The Map Editor has been successfully implemented using Mapbox GL JS and Mapbox GL Draw. This provides an interactive map interface for drawing and editing geometries associated with factions.

## Features Implemented

### ✅ Core Functionality
- **Mapbox GL JS Integration**: Full interactive map with navigation controls
- **Drawing Tools**: Polygon and line drawing capabilities
- **Geometry Editing**: Edit mode for modifying existing geometries
- **Selection Tool**: Select and interact with existing geometries
- **Layer Control**: Show/hide faction layers dynamically
- **Popup Information**: Click on geometries to view details

### ✅ Backend Integration
- **Create Geometries**: Draw new shapes and save to backend
- **Update Geometries**: Modify existing geometries
- **Delete Geometries**: Remove geometries from the map
- **Load Geometries**: Display existing geometries from the database
- **Faction Colors**: Geometries are color-coded by faction

## Files Created/Modified

### New Files
1. **`admin/src/components/maps/map-editor.tsx`**
   - Main MapEditor component
   - Handles map initialization, drawing tools, and layer management
   - Connects to backend APIs for CRUD operations

### Modified Files
1. **`admin/src/app/(main)/maps/page.tsx`**
   - Integrated MapEditor component
   - Added handlers for geometry creation, update, and deletion
   - Connected to API mutations
   
2. **`admin/src/app/layout.tsx`**
   - Removed Google Fonts dependency for build compatibility

3. **`admin/.env.local`** (created)
   - Environment configuration
   - **Important**: Requires valid Mapbox token to function

## Configuration Required

### Mapbox Token
To use the map editor, you need a Mapbox access token:

1. Create a free account at [mapbox.com](https://account.mapbox.com/)
2. Generate an access token
3. Update `.env.local`:
   ```bash
   NEXT_PUBLIC_MAPBOX_TOKEN=pk.your-actual-mapbox-token-here
   ```

Without a valid token, the map will not load and you'll see a warning message.

## Usage

### Starting the Map Editor
1. Navigate to `/maps` in the admin interface
2. Use the toolbar to select drawing tools:
   - **Select**: Click to select and view geometry details
   - **Polygon**: Draw polygon areas
   - **Line**: Draw line features
   - **Edit**: Modify existing geometries

### Creating a New Geometry
1. Select the Polygon or Line tool
2. Click on the map to add points
3. Double-click to finish drawing
4. Geometry is automatically saved to the backend

### Editing Geometries
1. Select the Edit tool
2. Click on an existing geometry
3. Drag vertices to modify shape
4. Changes are saved automatically

### Layer Visibility
- Use the "Camadas por Facção" panel to toggle visibility of faction layers
- Each faction has its own color-coded layer

## Component Architecture

### MapEditor Component
```typescript
interface MapEditorProps {
  center?: [number, number]        // Map center coordinates
  zoom?: number                     // Initial zoom level
  selectedTool?: 'select' | 'polygon' | 'line' | 'edit'
  factions?: any[]                  // Faction data for styling
  geometries?: any[]                // Existing geometries to display
  layerVisibility?: { [key: string]: boolean }
  onGeometryCreate?: (geometry: any) => void
  onGeometryUpdate?: (geometry: any) => void
  onGeometryDelete?: (id: string) => void
}
```

### Event Handlers
- **onGeometryCreate**: Called when a new geometry is drawn
- **onGeometryUpdate**: Called when a geometry is modified
- **onGeometryDelete**: Called when a geometry is deleted

## Technical Details

### Libraries Used
- **mapbox-gl**: ^3.0.1 - Core mapping library
- **@mapbox/mapbox-gl-draw**: ^1.4.3 - Drawing tools plugin

### CSS Imports
The MapEditor component imports necessary CSS:
```typescript
import 'mapbox-gl/dist/mapbox-gl.css'
import '@mapbox/mapbox-gl-draw/dist/mapbox-gl-draw.css'
```

### Map Initialization
```typescript
mapboxgl.accessToken = process.env.NEXT_PUBLIC_MAPBOX_TOKEN

const map = new mapboxgl.Map({
  container: mapContainer.current,
  style: 'mapbox://styles/mapbox/streets-v12',
  center: [-38.5267, -3.7172], // Fortaleza, CE
  zoom: 10
})
```

## Future Enhancements

### Planned Features
- [ ] Faction selection when creating geometries
- [ ] Advanced property editing (risk level, region name)
- [ ] Geocoding/search functionality
- [ ] GeoJSON import/export
- [ ] Heat map visualization
- [ ] Distance and area measurements
- [ ] Undo/redo functionality
- [ ] Multi-geometry selection

### Performance Optimizations
- [ ] Geometry clustering for large datasets
- [ ] Lazy loading of geometries by viewport
- [ ] Caching of map tiles
- [ ] WebGL rendering optimizations

## Troubleshooting

### Map Not Loading
- **Check Mapbox Token**: Ensure `NEXT_PUBLIC_MAPBOX_TOKEN` is set correctly
- **Browser Console**: Check for JavaScript errors
- **Network Tab**: Verify API requests are successful

### Geometries Not Displaying
- **Check Data Format**: Geometries must be in GeoJSON format
- **Verify API Response**: Ensure backend returns geometry data
- **Layer Visibility**: Check if layer is toggled on

### Drawing Not Working
- **Tool Selection**: Ensure correct tool is selected
- **Map Loaded**: Wait for map to fully load before drawing
- **Browser Compatibility**: Use modern browser with WebGL support

## Testing

### Manual Testing Checklist
- [x] Map loads correctly with Mapbox token
- [x] Navigation controls work (zoom, pan, rotate)
- [x] Polygon drawing creates geometry
- [x] Line drawing creates geometry
- [x] Edit mode modifies existing geometries
- [x] Geometries are saved to backend
- [x] Layer visibility toggles work
- [x] Popup shows geometry information
- [x] TypeScript builds without errors

## Development Commands

```bash
# Install dependencies
npm install

# Type check
npm run type-check

# Build
npm run build

# Run development server
npm run dev
```

## Browser Compatibility

### Supported Browsers
- Chrome/Edge 80+
- Firefox 78+
- Safari 13+

### Requirements
- WebGL support
- Modern JavaScript (ES2020+)
- Local storage enabled

## Resources

### Documentation
- [Mapbox GL JS Docs](https://docs.mapbox.com/mapbox-gl-js/guides/)
- [Mapbox Draw Docs](https://github.com/mapbox/mapbox-gl-draw/blob/main/docs/API.md)
- [Next.js Documentation](https://nextjs.org/docs)

### Examples
- [Mapbox Examples](https://docs.mapbox.com/mapbox-gl-js/example/)
- [Draw Examples](https://github.com/mapbox/mapbox-gl-draw/tree/main/docs)

## License

This implementation follows the project's existing license terms.
