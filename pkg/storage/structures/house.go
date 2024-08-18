package structures

import "time"

type House struct {
	Id        int       `db:"id" json:"id"`
	Address   string    `db:"address" json:"address"`
	Year      int       `db:"year" json:"year"`
	Developer string    `db:"developer" json:"developer"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdateAt  time.Time `db:"update_at" json:"update_at"`
}
