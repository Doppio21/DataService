package manager

import (
	"context"
	"dataservice/internal/api"
	"dataservice/internal/schema"
	"dataservice/internal/userdb"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestAddPersonInfo(t *testing.T) {
	const (
		name        = "Dmitry"
		surname     = "Federov"
		age         = 22
		gender      = "male"
		nationalize = "RU"
	)

	ctrl := gomock.NewController(t)
	api := api.NewAPIMock(ctrl)
	db := userdb.NewMockDB(ctrl)

	api.Age.EXPECT().Get(gomock.Any(), name).Return(age, nil)
	api.Gender.EXPECT().Get(gomock.Any(), name).Return(gender, nil)
	api.Nationalize.EXPECT().Get(gomock.Any(), name).Return(nationalize, nil)

	db.EXPECT().AddPersonInfo(gomock.Any(), schema.PersonInfo{
		ID:      0,
		Name:    name,
		Surname: surname,
		Age:     age,
		Gender:  gender,
		Country: nationalize,
	}).Return(nil)

	mgr := New(Config{Timeout: time.Second}, Dependencies{
		API: api,
		DB:  db,
	})

	err := mgr.AddPersonInfo(context.Background(), schema.PutRequest{
		Name:    name,
		Surname: surname,
	})
	require.NoError(t, err)
}

func TestGetPersonInfo(t *testing.T) {
	exp := []schema.PersonInfo{
		{
			ID:      12,
			Name:    "Dmitry",
			Surname: "Federov",
			Age:     22,
			Gender:  "male",
			Country: "RU",
		},
	}

	ctrl := gomock.NewController(t)
	db := userdb.NewMockDB(ctrl)

	db.EXPECT().GetPersonInfo(gomock.Any(), schema.GetRequest{
		ID: exp[0].ID,
	}).Return(exp, nil)

	mgr := New(Config{Timeout: time.Second}, Dependencies{
		DB: db,
	})

	res, err := mgr.GetPersonInfo(context.Background(), schema.GetRequest{
		ID: exp[0].ID,
	})

	require.Equal(t, exp, res)
	require.NoError(t, err)
}

func TestDeletePersonInfo(t *testing.T) {
	const id = 12

	ctrl := gomock.NewController(t)
	db := userdb.NewMockDB(ctrl)

	db.EXPECT().DeletePersonInfo(gomock.Any(), id).Return(nil)
	mgr := New(Config{Timeout: time.Second}, Dependencies{
		DB: db,
	})

	err := mgr.DeletePersonInfo(context.Background(), id)
	require.NoError(t, err)
}

func TestUpdatePersonInfo(t *testing.T) {
	const (
		id          = 12
		name        = "Dmitry"
		surname     = "Federov"
		age         = 22
		gender      = "male"
		nationalize = "RU"
	)

	ctrl := gomock.NewController(t)
	db := userdb.NewMockDB(ctrl)

	mgr := New(Config{Timeout: time.Second}, Dependencies{
		DB: db,
	})

	db.EXPECT().UpdatePersonInfo(gomock.Any(), schema.PersonInfo{
		ID:      id,
		Name:    name,
		Surname: surname,
		Age:     age,
		Gender:  gender,
		Country: nationalize,
	}).Return(nil)

	err := mgr.UpdatePersonInfo(context.Background(), schema.PersonInfo{
		ID:      id,
		Name:    name,
		Surname: surname,
		Age:     age,
		Gender:  gender,
		Country: nationalize,
	})
	require.NoError(t, err)
}
