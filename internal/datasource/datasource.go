package datasource

import (
	"context"
	"time"

	"github.com/dugtriol/backend-bootcamp-assignment-2024/internal/datasource/storage/structures"
	"github.com/google/uuid"
)

type Datasource interface {
	User
	House
	Flat
	GetList
}

type User interface {
	SaveUser(ctx context.Context, email, password, userType string) (uuid.UUID, error)
	GetUserById(ctx context.Context, id uuid.UUID) (*structures.User, error)
}

type House interface {
	SaveHouse(ctx context.Context, address, developer string, year int) (*structures.House, error)
	GetHouse(ctx context.Context, id int) (*structures.House, error)
	UpdateDate(ctx context.Context, time time.Time, id int) error
}

type Flat interface {
	SaveFlat(ctx context.Context, houseId, price, rooms int) (*structures.Flat, error)
	GetFlat(ctx context.Context, id int) (*structures.Flat, error)
	UpdateStatus(ctx context.Context, id int, status string) error
}

type GetList interface {
	GetListByClient(ctx context.Context, id int) (*[]structures.Flat, error)
	GetListByModerator(ctx context.Context, id int) (*[]structures.Flat, error)
}
