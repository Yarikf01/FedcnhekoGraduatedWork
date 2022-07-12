package api

import (
	"context"
	"errors"
)

type Manager interface {
	Upload(ctx context.Context) error
	Download(ctx context.Context) error
}

type Config struct {
}

func NewManager(cfg Config) Manager {
	return &manager{
	}
}

// impl

type manager struct {
}

func (m *manager) Upload(ctx context.Context) error {
	return errors.New("implement me")
}

func (m *manager) Download(ctx context.Context) error {
	return errors.New("implement me")
}
