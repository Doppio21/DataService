package nationalizeapi

import "context"

//go:generate mockgen -package nationalizeapi -destination api_mock.go . NationalizeAPI
type NationalizeAPI interface {
	Get(_ context.Context, name string) (nationalize string, _ error)
}
