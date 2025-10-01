package models

type CreateGeometryRequest struct {
	RegionID      *string `json:"region_id" validate:"required,uuid4"`
	GeometryType  string  `json:"geometry_type" validate:"required,oneof=polygon polyline"`
	GeoJSON       GeoJSON `json:"geojson" validate:"required"`
	FactionID     *string `json:"faction_id,omitempty" validate:"omitempty,uuid4"`
	RiskLevel     *string `json:"risk_level,omitempty" validate:"omitempty,oneof=low medium high"`
	ValidityStart *string `json:"validity_start,omitempty"`
	ValidityEnd   *string `json:"validity_end,omitempty"`
	Source        string  `json:"source,omitempty" validate:"omitempty,max=255"`
	Notes         string  `json:"notes,omitempty" validate:"omitempty"`
	Status        *string `json:"status,omitempty" validate:"omitempty,oneof=draft review published"`
}

type UpdateGeometryRequest struct {
	RegionID      *string  `json:"region_id,omitempty" validate:"omitempty,uuid4"`
	GeometryType  *string  `json:"geometry_type,omitempty" validate:"omitempty,oneof=polygon polyline"`
	GeoJSON       *GeoJSON `json:"geojson,omitempty"`
	FactionID     *string  `json:"faction_id,omitempty" validate:"omitempty,uuid4"`
	RiskLevel     *string  `json:"risk_level,omitempty" validate:"omitempty,oneof=low medium high"`
	ValidityStart *string  `json:"validity_start,omitempty"`
	ValidityEnd   *string  `json:"validity_end,omitempty"`
	Source        *string  `json:"source,omitempty" validate:"omitempty,max=255"`
	Notes         *string  `json:"notes,omitempty"`
	Status        *string  `json:"status,omitempty" validate:"omitempty,oneof=draft review published"`
}
