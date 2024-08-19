package storage

import (
	"context"
	"log/slog"
	"time"

	"github.com/dugtriol/backend-bootcamp-assignment-2024/internal/datasource/storage/structures"
	"github.com/dugtriol/backend-bootcamp-assignment-2024/pkg/db"
	"github.com/google/uuid"
)

type Storage struct {
	db  *db.Database
	log *slog.Logger
}

func New(database *db.Database, log *slog.Logger) *Storage {
	return &Storage{db: database, log: log}
}

func (r *Storage) SaveUser(ctx context.Context, email, password, userType string) (uuid.UUID, error) {
	id := uuid.New()

	if _, err := r.db.Exec(
		ctx,
		`INSERT INTO users(id, email, password, type) VALUES($1, $2, $3, $4);`,
		id,
		email,
		password,
		userType,
	); err != nil {
		r.log.Error("database: failed to save user")
		return uuid.Nil, err
	}

	return id, nil
}

func (r *Storage) GetUserById(ctx context.Context, id uuid.UUID) (*structures.User, error) {
	var a structures.User

	err := r.db.Get(ctx, &a, "SELECT id,email,password,type FROM users WHERE id=$1", id)
	if err != nil {
		r.log.Error("database: failed to get user by id")
		return nil, err
	}
	return &a, nil
}

func (r *Storage) SaveHouse(ctx context.Context, address, developer string, year int) (*structures.House, error) {
	var house structures.House
	err := r.db.ExecQueryRow(
		ctx,
		`INSERT INTO houses(address, developer, year) VALUES($1, $2, $3) RETURNING *`,
		address,
		developer,
		year,
	).Scan(&house.Id, &house.Address, &house.Year, &house.Developer, &house.CreatedAt, &house.UpdateAt)
	if err != nil {
		r.log.Error("database: failed to save house")
		return nil, err
	}
	return &house, err
}

func (r *Storage) GetHouse(ctx context.Context, id int) (*structures.House, error) {
	var house structures.House
	err := r.db.Get(
		ctx,
		&house,
		"SELECT id,address,year,developer,created_at,update_at FROM houses WHERE id=$1",
		id,
	)
	if err != nil {
		r.log.Error("database: failed to egt house")
		return nil, err
	}

	return &house, nil
}

func (r *Storage) SaveFlat(ctx context.Context, houseId, price, rooms int) (*structures.Flat, error) {
	var flat structures.Flat
	err := r.db.ExecQueryRow(
		ctx,
		`INSERT INTO flats(house_id, price, rooms) VALUES($1, $2, $3) RETURNING *`,
		houseId,
		price,
		rooms,
	).Scan(&flat.Id, &flat.HouseId, &flat.Price, &flat.Rooms, &flat.Status)

	if err != nil {
		r.log.Error("database: failed to save flat")
		return nil, err
	}
	return &flat, nil
}

func (r *Storage) GetFlat(ctx context.Context, id int) (*structures.Flat, error) {
	var flat structures.Flat
	err := r.db.Get(
		ctx,
		&flat,
		"SELECT id,house_id,price,rooms,status FROM flats WHERE id=$1", id,
	)
	if err != nil {
		r.log.Error("database: failed to get flat")
		return nil, err
	}
	return &flat, nil
}

func (r *Storage) UpdateDate(ctx context.Context, time time.Time, id int) error {
	_, err := r.db.Exec(
		ctx,
		"UPDATE houses SET update_at = $1 WHERE id = $2",
		time,
		id,
	)
	if err != nil {
		r.log.Error("database: failed to update date")
		return err
	}
	return nil
}

func (r *Storage) UpdateStatus(ctx context.Context, id int, status string) error {
	_, err := r.db.Exec(
		ctx,
		"UPDATE flats SET status = $1 WHERE id = $2",
		status,
		id,
	)
	if err != nil {
		r.log.Error("database: failed to update status")
		return err
	}
	return nil
}

func (r *Storage) GetListByClient(ctx context.Context, id int) (*[]structures.Flat, error) {
	status := "approved"
	rows, err := r.db.Query(
		ctx,
		"SELECT id,house_id,price,rooms,status FROM flats WHERE house_id=$1 AND status=$2", id, status,
	)
	if err != nil {
		r.log.Error("database: failed to get list by client", err)
		return nil, err
	}
	defer rows.Close()

	var flats []structures.Flat
	for rows.Next() {
		var flat structures.Flat
		if err := rows.Scan(&flat.Id, &flat.HouseId, &flat.Price, &flat.Rooms, &flat.Status); err != nil {
			return &flats, err
		}
		flats = append(flats, flat)
	}
	if err = rows.Err(); err != nil {
		r.log.Error("database: failed to get list by client", err)
		return &flats, err
	}
	return &flats, nil
}

func (r *Storage) GetListByModerator(ctx context.Context, id int) (*[]structures.Flat, error) {
	r.log.Info("database start")
	rows, err := r.db.Query(
		ctx,
		"SELECT id,house_id,price,rooms,status FROM flats WHERE house_id=$1", id,
	)
	if err != nil {
		r.log.Error("database: failed to get list by client", err)
		return nil, err
	}
	defer rows.Close()

	var flats []structures.Flat
	for rows.Next() {
		var flat structures.Flat
		if err := rows.Scan(&flat.Id, &flat.HouseId, &flat.Price, &flat.Rooms, &flat.Status); err != nil {
			r.log.Error("database: failed to get list by client", err)
			return &flats, err
		}
		flats = append(flats, flat)
	}
	if err = rows.Err(); err != nil {
		r.log.Error("database: failed to get list by client", err)
		return &flats, err
	}
	r.log.Info("database end")
	return &flats, nil
}
