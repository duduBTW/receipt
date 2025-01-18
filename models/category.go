package models

import "time"

type Category struct {
	ID             int64     `db:"id"`
	Name           string    `db:"name"`
	LucideIconName string    `db:"lucide_icon_name"`
	Hue            int       `db:"hue"`
	Saturation     int       `db:"saturation"`
	Lightness      int       `db:"lightness"`
	CreatedAt      time.Time `db:"created_at"`
}
