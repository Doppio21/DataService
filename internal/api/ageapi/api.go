package ageapi

import "context"

//go:generate mockgen -package ageapi -destination api_mock.go . AgeAPI
type AgeAPI interface {
	Get(_ context.Context, name string) (age int, _ error)
}
