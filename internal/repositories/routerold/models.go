package routerold

import (
	"time"

	"github.com/uptrace/bun"
)

// RouterModel Database model for router
type RouterModel struct {
	bun.BaseModel `bun:"table:routers"`

	ID        int    `bun:"id,pk"`
	TenantID  int    `bun:"id_tenant"`
	IMEI      string `bun:"imei"`
	TotemID   string `bun:"idTotem"`
	TotemKey  string `bun:"totemKey"`
	IpVPN     string `bun:"ip_vpn"`
	Latitude  string `bun:"lat"`
	Longitude string `bun:"lng"`
	//Altitude  float32    `bun:"alt"` ////Altitude
	Active    int        `bun:"active"`
	NotifyC5  int        `bun:"notify_c5_cdmx"`
	NotifyC5J int        `bun:"notify_c5_jal"`
	CreatedAt *time.Time `bun:"created"`
	UpdatedAt *time.Time `bun:"updated"`
}
