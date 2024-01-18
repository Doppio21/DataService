package userdb

import (
	"context"
	"dataservice/internal/schema"
	"errors"
)

var ErrNotFound = errors.New("not found")

//go:generate mockgen -package userdb -destination db_mock.go . DB
type DB interface {
	AddPersonInfo(ctx context.Context, info schema.PersonInfo) error
	GetPersonInfo(ctx context.Context, req schema.GetRequest) ([]schema.PersonInfo, error)
	DeletePersonInfo(ctx context.Context, id int) error
	UpdatePersonInfo(ctx context.Context, info schema.PersonInfo) error
}
