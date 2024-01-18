package manager

import (
	"context"
	"dataservice/internal/api"
	"dataservice/internal/schema"
	"dataservice/internal/userdb"
	"dataservice/internal/utils"
	"time"

	"go.uber.org/zap"
)

type Config struct {
	Timeout time.Duration
}

type Dependencies struct {
	API api.API
	DB  userdb.DB

	Log *zap.Logger
}

type Manager struct {
	cfg  Config
	deps Dependencies
}

func New(cfg Config, deps Dependencies) *Manager {
	return &Manager{
		cfg:  cfg,
		deps: deps,
	}
}

func (m *Manager) enrichMessage(ctx context.Context, req schema.PutRequest) (schema.PersonInfo, error) {
	var (
		age         int
		gender      string
		nationalize string
	)

	err := utils.ParallelRequest(ctx, m.cfg.Timeout,
		func(ctx context.Context) error {
			var err error
			age, err = m.deps.API.AgeAPI().Get(ctx, req.Name)
			return err
		},
		func(ctx context.Context) error {
			var err error
			gender, err = m.deps.API.GenderAPI().Get(ctx, req.Name)
			return err
		},
		func(ctx context.Context) error {
			var err error
			nationalize, err = m.deps.API.NationalizeAPI().Get(ctx, req.Name)
			return err
		},
	)
	if err != nil {
		m.deps.Log.Error("failed to API reqeusts", zap.Error(err))
		return schema.PersonInfo{}, err
	}

	ret := schema.PersonInfo{
		Name:    req.Name,
		Surname: req.Surname,
		Age:     age,
		Gender:  gender,
		Country: nationalize,
	}

	return ret, nil
}

func (m *Manager) AddPersonInfo(ctx context.Context, req schema.PutRequest) error {
	info, err := m.enrichMessage(ctx, req)
	if err != nil {
		return err
	}

	if err = m.deps.DB.AddPersonInfo(ctx, info); err != nil {
		m.deps.Log.Error("error adding to database:", zap.Error(err))
		return err
	}
	return nil
}

func (m *Manager) GetPersonInfo(ctx context.Context, req schema.GetRequest) ([]schema.PersonInfo, error) {
	ret, err := m.deps.DB.GetPersonInfo(ctx, req)
	if err != nil {
		m.deps.Log.Error("error getting from database", zap.Error(err))
		return nil, err
	}
	return ret, nil
}

func (m *Manager) DeletePersonInfo(ctx context.Context, id int) error {
	if err := m.deps.DB.DeletePersonInfo(ctx, id); err != nil {
		m.deps.Log.Error("error deleting from database", zap.Error(err))
		return err
	}
	return nil
}

func (m *Manager) UpdatePersonInfo(ctx context.Context, info schema.PersonInfo) error {
	if err := m.deps.DB.UpdatePersonInfo(ctx, info); err != nil {
		m.deps.Log.Error("error updating database information", zap.Error(err))
		return err
	}
	return nil
}
