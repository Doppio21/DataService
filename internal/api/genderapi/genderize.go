package genderapi

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"go.uber.org/zap"
)

type genderizeResponse struct {
	Count       int     `json:"count"`
	Name        string  `json:"name"`
	Gender      string  `json:"gender"`
	Probability float32 `json:"probability"`
}

type Config struct {
	URI string
}

type Dependencies struct {
	Client *http.Client
	Log    *zap.Logger
}

type genderize struct {
	cfg  Config
	deps Dependencies

	log *zap.Logger
}

func NewGenderize(cfg Config, deps Dependencies) GenderAPI {
	return &genderize{
		cfg:  cfg,
		deps: deps,
		log:  deps.Log.Named("genderize"),
	}
}

func (g *genderize) Get(ctx context.Context, name string) (string, error) {
	req, err := http.NewRequest(http.MethodGet, g.cfg.URI, nil)
	if err != nil {
		g.log.Error("failed to create http request", zap.Error(err))
		return "", err
	}

	q := req.URL.Query()
	q.Add("name", name)
	req.URL.RawQuery = q.Encode()

	req = req.WithContext(ctx)
	resp, err := g.deps.Client.Do(req)
	if err != nil {
		g.log.Error("failed to do http request", zap.Error(err))
		return "", err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		g.log.Error("failed to read response body", zap.Error(err))
		return "", err
	}

	res := genderizeResponse{}
	if err := json.Unmarshal(body, &res); err != nil {
		g.log.Error("failed to unmarshal response", zap.Error(err))
		return "", err
	}

	g.log.Debug("success request", zap.Any("resp", res))
	return res.Gender, nil
}
