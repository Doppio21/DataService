package api

import (
	"dataservice/internal/api/ageapi"
	"dataservice/internal/api/genderapi"
	"dataservice/internal/api/nationalizeapi"

	"go.uber.org/mock/gomock"
)

type APIMock struct {
	Age         *ageapi.MockAgeAPI
	Gender      *genderapi.MockGenderAPI
	Nationalize *nationalizeapi.MockNationalizeAPI
}

func NewAPIMock(ctrl *gomock.Controller) *APIMock {
	return &APIMock{
		Age:         ageapi.NewMockAgeAPI(ctrl),
		Gender:      genderapi.NewMockGenderAPI(ctrl),
		Nationalize: nationalizeapi.NewMockNationalizeAPI(ctrl),
	}
}

func (m *APIMock) AgeAPI() ageapi.AgeAPI {
	return m.Age
}

func (m *APIMock) GenderAPI() genderapi.GenderAPI {
	return m.Gender
}

func (m *APIMock) NationalizeAPI() nationalizeapi.NationalizeAPI {
	return m.Nationalize
}
