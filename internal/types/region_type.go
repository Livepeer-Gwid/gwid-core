package types

import "github.com/google/uuid"

type RegionRes struct {
	ID         uuid.UUID `json:"id"`
	RegionName string    `json:"region_name"`
	Status     string    `json:"status"`
	Endpoint   string    `json:"endpoint"`
}
