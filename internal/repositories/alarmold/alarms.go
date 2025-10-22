package alarmold

import (
	"time"

	"github.com/uptrace/bun"
)

// AlarmModel Database model for old alarms
type AlarmModel struct {
	bun.BaseModel `bun:"table:alarms"`

	ID                string    `bun:"id,pk"`
	RouterID          string    `bun:"id_router"`
	Latitude          float64   `bun:"lat"`
	Longitude         float64   `bun:"lng"`
	CanceledComments  string    `bun:"cancel_comments"`
	SolvedComments    string    `bun:"solved_comments"`
	Waiting           int       `bun:"waiting"`
	Attending         int       `bun:"attending"`
	Attended          int       `bun:"attended"`
	Canceled          int       `bun:"canceled"`
	Request           string    `bun:"request"`
	NumerCad          string    `bun:"number_cad"`
	StatusCad         string    `bun:"status_cad"`
	ClosedCad         int       `bun:"closed_cad"`
	UserID            int       `bun:"id_user"`
	AlarmType         int       `bun:"type"`
	SmartDetection    time.Time `bun:"smart_detection_timestamp"`
	SmartDetectionID  int       `bun:"smart_detection_id"`
	AttendingComments string    `bun:"attending_comments"`
	AlarmLevelID      int       `bun:"id_alarm_level"`
	Protocol          string    `bun:"protocol"`
	Source            string    `bun:"source"`
	CreatedAt         time.Time `bun:"created"`
	UpdatedAt         time.Time `bun:"updated"`
}
