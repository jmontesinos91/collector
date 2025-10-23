package traffic

import (
	"github.com/jmontesinos91/collector/domains/pagination"
	"time"

	"github.com/uptrace/bun"
)

// Model Database model for traffic
type Model struct {
	bun.BaseModel `bun:"table:traffic"`

	ID         string    `bun:"id,pk"`
	Request    string    `bun:"request,pk"`
	IMEI       string    `bun:"imei"`
	Ip         string    `bun:"ip"`
	IsAlarm    bool      `bun:"isAlarm"`
	IsNotified bool      `bun:"isnotified"`
	Counter    int       `bun:"counter"`
	CreatedAt  time.Time `bun:"created_at"`
	UpdatedAt  time.Time `bun:"updated_at"`
}

// Metadata struct filter for repository layer
type Metadata struct {
	Qparam  string
	ID      string
	Request string
	IMEI    string
	Ip      string
	IsAlarm *bool
	Counter *int
	Filter  pagination.Filter
}
