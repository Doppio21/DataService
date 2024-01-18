package genderapi

import "context"

//go:generate mockgen -package genderapi -destination api_mock.go . GenderAPI
type GenderAPI interface {
	Get(_ context.Context, name string) (gender string, _ error)
}
