package db

import (
	"context"
	"dataservice/internal/pgxprovider"
	"dataservice/internal/schema"
	"dataservice/internal/userdb"
	"errors"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

type Config struct {
	QueryTimeout time.Duration
}

type Dependencies struct {
	Log *zap.Logger
	PGX *pgxprovider.PGXProvider
}

type Postgres struct {
	cfg  Config
	deps Dependencies
}

func New(cfg Config, deps Dependencies) userdb.DB {
	return &Postgres{
		cfg:  cfg,
		deps: deps,
	}
}

func (p *Postgres) AddPersonInfo(ctx context.Context, personInfo schema.PersonInfo) error {
	_, err := p.deps.PGX.Exec(ctx, `INSERT INTO userDB (user_name, surname, age, gender, country)
									     VALUES ($1, $2, $3, $4, $5)`, personInfo.Name, personInfo.Surname,
		personInfo.Age, personInfo.Gender, personInfo.Country)
	if err != nil {
		p.deps.Log.Error("failed to insert", zap.Error(err))
		return err
	}
	p.deps.Log.Info("adding personal info to database")
	return nil
}

func (p *Postgres) buildGetQuery(request schema.GetRequest) (string, []interface{}, error) {
	pred := squirrel.Eq{}
	if request.Age != 0 {
		pred["age"] = request.Age
	}
	if request.Country != "" {
		pred["country"] = request.Country
	}
	if request.ID != 0 {
		pred["user_id"] = request.ID
	}
	if request.Name != "" {
		pred["user_name"] = request.Name
	}
	if request.Surname != "" {
		pred["surname"] = request.Surname
	}
	if request.Gender != "" {
		pred["gender"] = request.Gender
	}

	b := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	q := b.Select("user_id, user_name, surname, age, gender, country").
		From("userDB").Where(pred)
	if request.Count != 0 {
		q = q.Limit(uint64(request.Count))
	}
	if request.Offset != 0 {
		q = q.Offset(uint64(request.Offset))
	}

	return q.ToSql()
}

func (p *Postgres) GetPersonInfo(ctx context.Context, request schema.GetRequest) ([]schema.PersonInfo, error) {
	sql, args, err := p.buildGetQuery(request)
	if err != nil {
		p.deps.Log.Error("failed to build query", zap.Error(err))
		return nil, err
	}

	p.deps.Log.Debug("select query", zap.String("query", sql), zap.Any("args", args))

	res, err := p.deps.PGX.Query(ctx, sql, args...)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, userdb.ErrNotFound
	} else if err != nil {
		p.deps.Log.Error("failed to select", zap.Error(err))
		return nil, err
	}

	defer res.Close()
	ret := make([]schema.PersonInfo, 0)
	for res.Next() {
		cur := schema.PersonInfo{}
		err := res.Scan(&cur.ID, &cur.Name, &cur.Surname, &cur.Age, &cur.Gender, &cur.Country)
		if err != nil {
			p.deps.Log.Error("failed to scan rows", zap.Error(err))
			return nil, err
		}

		ret = append(ret, cur)
	}

	return ret, nil
}

func (p *Postgres) DeletePersonInfo(ctx context.Context, id int) error {
	_, err := p.deps.PGX.Exec(ctx, `DELETE FROM userDB WHERE user_id = $1`, id)
	if errors.Is(err, pgx.ErrNoRows) {
		return userdb.ErrNotFound
	} else if err != nil {
		p.deps.Log.Error("failed to delete", zap.Error(err))
		return err
	}
	p.deps.Log.Info("deleting personal info from the database")
	return nil
}

func (p *Postgres) UpdatePersonInfo(ctx context.Context, info schema.PersonInfo) error {
	_, err := p.deps.PGX.Exec(ctx, `UPDATE userDB 
									SET user_name = $1, surname = $2, age = $3, gender = $4, country = $5
									WHERE user_id = $6`,
		info.Name, info.Surname, info.Age, info.Gender, info.Country, info.ID)
	if errors.Is(err, pgx.ErrNoRows) {
		return userdb.ErrNotFound
	} else if err != nil {
		p.deps.Log.Error("failed to update", zap.Error(err))
		return err
	}
	p.deps.Log.Info("updating personal info in the database")
	return nil
}
