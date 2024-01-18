package api

import (
	"dataservice/internal/api/ageapi"
	"dataservice/internal/api/genderapi"
	"dataservice/internal/api/nationalizeapi"
)

type API interface {
	AgeAPI() ageapi.AgeAPI
	GenderAPI() genderapi.GenderAPI
	NationalizeAPI() nationalizeapi.NationalizeAPI
}

type Dependencies struct {
	Age ageapi.AgeAPI
	Gender genderapi.GenderAPI
	Nationalize nationalizeapi.NationalizeAPI
}

type api struct {
	deps Dependencies
}

func NewAPI(deps Dependencies) API {
	return &api{
		deps: deps,
	}
}

func (api *api) AgeAPI() ageapi.AgeAPI {
	return api.deps.Age
}

func (api *api) GenderAPI() genderapi.GenderAPI {
	return api.deps.Gender
}

func (api *api) NationalizeAPI() nationalizeapi.NationalizeAPI {
	return api.deps.Nationalize
}
