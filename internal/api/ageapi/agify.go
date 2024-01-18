package ageapi

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"go.uber.org/zap"
)

type agifyResponse struct {
	Count int    `json:"count"`
	Name  string `json:"name"`
	Age   int    `json:"age"`
}

type Config struct {
	URI string
}

type Dependencies struct {
	Client *http.Client
	Log    *zap.Logger
}

type agify struct {
	cfg  Config
	deps Dependencies

	log *zap.Logger
}

func NewAgify(cfg Config, deps Dependencies) AgeAPI {
	return &agify{
		cfg:  cfg,
		deps: deps,
		log:  deps.Log.Named("agify"),
	}
}

func (ag *agify) Get(ctx context.Context, name string) (int, error) {
	req, err := http.NewRequest(http.MethodGet, ag.cfg.URI, nil)
	if err != nil {
		ag.log.Error("failed to create http request", zap.Error(err))
		return 0, err
	}

	q := req.URL.Query()
	q.Add("name", name)
	req.URL.RawQuery = q.Encode()

	req = req.WithContext(ctx)
	resp, err := ag.deps.Client.Do(req)
	if err != nil {
		ag.log.Error("failed to do http request", zap.Error(err))
		return 0, err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		ag.log.Error("failed to read response body", zap.Error(err))
		return 0, err
	}

	res := agifyResponse{}
	if err := json.Unmarshal(body, &res); err != nil {
		ag.log.Error("failed to unmarshal response", zap.Error(err))
		return 0, err
	}

	ag.log.Debug("success request", zap.Any("resp", res))
	return res.Age, nil
}
