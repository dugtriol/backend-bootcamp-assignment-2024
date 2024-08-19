package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/allegro/bigcache/v3"
	"github.com/dugtriol/backend-bootcamp-assignment-2024/internal/datasource"
	"github.com/dugtriol/backend-bootcamp-assignment-2024/internal/datasource/storage/structures"
	"github.com/dugtriol/backend-bootcamp-assignment-2024/internal/handlers/house"
	"github.com/google/uuid"
)

const (
	moderatorAll = "moderator:all"
	clientAll    = "client:all"
)

type Client struct {
	source datasource.Datasource
	conn   *bigcache.BigCache
	log    *slog.Logger
}

func NewClient(log *slog.Logger, conn *bigcache.BigCache, source datasource.Datasource) *Client {
	return &Client{source: source, conn: conn, log: log}
}

func (c Client) SaveUser(ctx context.Context, email, password, userType string) (uuid.UUID, error) {
	var err error
	flat, err := c.source.SaveUser(ctx, email, password, userType)
	if err != nil {
		return uuid.Nil, err
	}

	return flat, nil
}

func (c Client) GetUserById(ctx context.Context, id uuid.UUID) (*structures.User, error) {
	result, err := c.source.GetUserById(ctx, id)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c Client) SaveHouse(ctx context.Context, address, developer string, year int) (*structures.House, error) {
	var err error
	result, err := c.source.SaveHouse(ctx, address, developer, year)
	if err != nil {
		return nil, err
	}

	if err = c.conn.Delete(fmt.Sprintf("%s:%d", clientAll, result.Id)); err != nil {
		c.log.Error("failed to delete list of flats from cache (client)", err)
	}

	if err = c.conn.Delete(fmt.Sprintf("%s:%d", moderatorAll, result.Id)); err != nil {
		c.log.Error("failed to delete list of flats from cache (moderator)", err)
	}

	return result, nil
}

func (c Client) GetHouse(ctx context.Context, id int) (*structures.House, error) {
	result, err := c.source.GetHouse(ctx, id)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c Client) UpdateDate(ctx context.Context, time time.Time, id int) error {
	var err error
	err = c.source.UpdateDate(ctx, time, id)
	if err != nil {
		return err
	}

	if err = c.conn.Delete(fmt.Sprintf("%s:%d", clientAll, id)); err != nil {
		c.log.Error("failed to delete list of flats from cache (client)", err)
	}

	if err = c.conn.Delete(fmt.Sprintf("%s:%d", moderatorAll, id)); err != nil {
		c.log.Error("failed to delete list of flats from cache (moderator)", err)
	}
	return nil
}

func (c Client) SaveFlat(ctx context.Context, houseId, price, rooms int) (*structures.Flat, error) {
	var err error
	flat, err := c.source.SaveFlat(ctx, houseId, price, rooms)
	if err != nil {
		return nil, err
	}

	if err = c.conn.Delete(fmt.Sprintf("%s:%d", clientAll, houseId)); err != nil {
		c.log.Error("failed to delete list of flats from cache (client)", err)
	}

	if err = c.conn.Delete(fmt.Sprintf("%s:%d", moderatorAll, houseId)); err != nil {
		c.log.Error("failed to delete list of flats from cache (moderator)", err)
	}

	return flat, nil
}

func (c Client) GetFlat(ctx context.Context, id int) (*structures.Flat, error) {
	flat, err := c.source.GetFlat(ctx, id)
	if err != nil {
		return nil, err
	}
	return flat, nil
}

func (c Client) UpdateStatus(ctx context.Context, id int, status string) error {
	var err error
	err = c.source.UpdateStatus(ctx, id, status)
	if err != nil {
		return err
	}

	if err = c.conn.Delete(fmt.Sprintf("%s:%d", clientAll, id)); err != nil {
		c.log.Error("failed to delete list of flats from cache (client)", err)
	}

	if err = c.conn.Delete(fmt.Sprintf("%s:%d", moderatorAll, id)); err != nil {
		c.log.Error("failed to delete list of flats from cache (moderator)", err)
	}
	return nil
}

func (c Client) GetListByClient(ctx context.Context, id int) (*[]structures.Flat, error) {
	c.log.Info("cache start")
	var err error
	data, err := c.conn.Get(fmt.Sprintf("%s:%d", clientAll, id))
	if errors.Is(err, bigcache.ErrEntryNotFound) {
		var err error
		list, err := c.source.GetListByClient(ctx, id)
		if err != nil {
			c.log.Error("failed to cache list flats client")
			return nil, err
		}

		flats := house.GetListResponse{Flats: list}

		resp, err := json.Marshal(flats)

		if err != nil {
			c.log.Error("failed to cache list flats client")
			return nil, err
		}

		if err := c.conn.Set(fmt.Sprintf("%s:%d", clientAll, id), resp); err != nil {
			c.log.Error("failed to cache list flats client")
			return nil, err
		}
		c.log.Info("cache end")
		return list, nil
	}
	if err != nil {
		c.log.Error("failed to cache list flats client")
		return nil, err
	}

	var result house.GetListResponse
	err = json.Unmarshal(data, &result)
	if err != nil {
		c.log.Error("failed to cache list flats client")
		return nil, err
	}
	c.log.Info("cache end")

	return result.Flats, nil
}

func (c Client) GetListByModerator(ctx context.Context, id int) (*[]structures.Flat, error) {
	c.log.Info("cache start")
	var err error
	data, err := c.conn.Get(fmt.Sprintf("%s:%d", moderatorAll, id))
	if errors.Is(err, bigcache.ErrEntryNotFound) {
		var err error
		list, err := c.source.GetListByModerator(ctx, id)
		if err != nil {
			c.log.Error("failed to cache list flats client")
			return nil, err
		}

		flats := house.GetListResponse{Flats: list}

		resp, err := json.Marshal(flats)
		if err != nil {
			c.log.Error("failed to cache list flats client")
			return nil, err
		}

		if err = c.conn.Set(fmt.Sprintf("%s:%d", moderatorAll, id), resp); err != nil {
			c.log.Error("failed to cache list flats client")
			return nil, err
		}
		c.log.Info("cache end")
		return list, nil
	}
	if err != nil {
		c.log.Error("failed to cache list flats client")
		return nil, err
	}

	var result house.GetListResponse
	err = json.Unmarshal(data, &result)
	if err != nil {
		c.log.Error("failed to cache list flats client")
		return nil, err
	}
	c.log.Info("cache end")

	return result.Flats, nil
}
