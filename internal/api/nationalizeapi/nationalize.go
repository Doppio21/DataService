package nationalizeapi

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"sort"

	"go.uber.org/zap"
)

type nationalizeResponse struct {
	Count   int    `json:"count"`
	Name    string `json:"name"`
	Country []country
}

type country struct {
	CountryID   string  `json:"country_id"`
	Probability float32 `json:"probability"`
}

type Config struct {
	URI string
}

type Dependencies struct {
	Client *http.Client
	Log    *zap.Logger
}

type nationalize struct {
	cfg  Config
	deps Dependencies

	log *zap.Logger
}

func NewNationalize(cfg Config, deps Dependencies) NationalizeAPI {
	return &nationalize{
		cfg:  cfg,
		deps: deps,
		log:  deps.Log.Named("nationalize"),
	}
}

func (n *nationalize) Get(ctx context.Context, name string) (string, error) {
	req, err := http.NewRequest(http.MethodGet, n.cfg.URI, nil)
	if err != nil {
		n.log.Error("failed to create http request", zap.Error(err))
		return "", err
	}

	q := url.Values{}
	q.Add("name", name)
	req.URL.RawQuery = q.Encode()

	req = req.WithContext(ctx)
	resp, err := n.deps.Client.Do(req)
	if err != nil {
		n.log.Error("failed to do http request", zap.Error(err))
		return "", err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		n.log.Error("failed to read response body", zap.Error(err))
		return "", err
	}

	res := nationalizeResponse{}
	if err := json.Unmarshal(body, &res); err != nil {
		n.log.Error("failed to unmarshal response", zap.Error(err))
		return "", err
	}

	n.log.Debug("success request", zap.Any("resp", res))
	if len(res.Country) == 0 {
		return "unknown", nil
	}

	sort.Slice(res.Country, func(i, j int) bool {
		return res.Country[i].Probability < res.Country[j].Probability
	})

	return res.Country[0].CountryID, nil
}
