package structures

import (
	"github.com/google/uuid"
)

type User struct {
	Id       uuid.UUID `db:"id"`
	Email    string    `db:"email"`
	Password string    `db:"password"`
	Type     string    `db:"type"`
}
