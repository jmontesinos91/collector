package unitsold

import (
	"time"

	"github.com/uptrace/bun"
)

// UnitsModel Database model for units
type UnitsModel struct {
	bun.BaseModel `bun:"table:units"`

	ID                   int    `bun:"id,pk"`
	RouterID             int    `bun:"id_router"`
	Description          string `bun:"description"`
	Address1             string `bun:"address_1"`
	Address2             string `bun:"address_2"`
	ExteriorNumber       string `bun:"exterior_number"`
	Neighborhood         string `bun:"suburb"`
	Number               string `bun:"number"`
	PostalCode           string `bun:"postal_code"`
	CountryID            int    `bun:"id_country"`
	RegionID             int    `bun:"id_region"`
	IsActive             bool   `bun:"active"`
	Route                string `bun:"route"`
	PlateNumber          string `bun:"plate_number"`
	VIN                  string `bun:"vin"`
	Driver               string `bun:"driver"`
	Provider             string `bun:"provider"`
	Owner                string `bun:"owner"`
	Company              string `bun:"company"`
	IsVehicle            bool   `bun:"vehicle"`
	IsVirtualButtonAlarm bool   `bun:"allow_virtual_button_alarm"`
	PatrolReaction       string `bun:"reaction_patrol"`
	AlarmLevelsID        int    `bun:"id_alarms_levels"`
	IsDigitalOutPut      bool   `bun:"active_digital_output"`
	MunicipalityID       int    `bun:"id_municipality"`
}

// UnitsCamerasModel Database model for units_cameras
type UnitsCamerasModel struct {
	bun.BaseModel `bun:"table:units_cameras"`

	UnitID    int        `bun:"id_unit"`
	CameraID  int        `bun:"id_camera"`
	CreatedAt *time.Time `bun:"created_at"`
	UpdatedAt *time.Time `bun:"updated_at"`
}
