package server

import (
	"context"
	"dataservice/internal/manager"
	"dataservice/internal/schema"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

const (
	shutdownTimeout = 5 * time.Second
)

type Config struct {
	Address string
}

type Dependencies struct {
	Manager manager.Manager
	Log     *zap.Logger
}

type Server struct {
	cfg  Config
	deps Dependencies
}

func New(cfg Config, deps Dependencies) *Server {
	return &Server{
		cfg:  cfg,
		deps: deps,
	}
}

func (s *Server) Run(ctx context.Context) error {
	router := gin.New()

	router.Use(LoggerMiddleware(s.deps.Log))
	router.PUT("/", s.addHandler)
	router.GET("/", s.getHandler)
	router.DELETE("/:id", s.deleteHandler)
	router.POST("/:id", s.updateHandler)

	srv := &http.Server{
		Addr:    s.cfg.Address,
		Handler: router,
	}

	serverClosed := make(chan struct{})
	go func() {
		s.deps.Log.Info("server started")
		defer close(serverClosed)
		if err := srv.ListenAndServe(); err == nil && err != http.ErrServerClosed {
			s.deps.Log.Fatal("listen and serve:", zap.Error(err))
		}
	}()

	select {
	case <-ctx.Done():
		s.deps.Log.Info("shutting down server gracefully")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()

		if err := srv.Shutdown(shutdownCtx); err != nil {
			return fmt.Errorf("shutdown: %w", err)
		}
	case <-serverClosed:
	}

	s.deps.Log.Info("server finished")
	return nil
}

func (s *Server) addHandler(c *gin.Context) {
	req := schema.PutRequest{}
	data, err := io.ReadAll(c.Request.Body)
	if s.replyError(c, err) {
		s.deps.Log.Error("failed to read body:", zap.Error(err))
		return
	}

	err = json.Unmarshal(data, &req)
	if s.replyError(c, err) {
		s.deps.Log.Error("failed to unmarshal request:", zap.Error(err))
		return
	}

	err = s.deps.Manager.AddPersonInfo(c, req)
	if s.replyError(c, err) {
		s.deps.Log.Error("failed add person:", zap.Error(err))
		return
	}

	c.Status(http.StatusOK)
}

func (s *Server) getQuery(query url.Values) (schema.GetRequest, error) {
	ret := schema.GetRequest{}
	for key, value := range query {
		v := value[len(value)-1]
		var err error

		switch key {
		case "id":
			ret.ID, err = strconv.Atoi(v)
		case "name":
			ret.Name = v
		case "surname":
			ret.Surname = v
		case "age":
			ret.Age, err = strconv.Atoi(v)
		case "gender":
			ret.Gender = v
		case "country":
			ret.Country = v
		case "count":
			ret.Count, err = strconv.Atoi(v)
		case "offset":
			ret.Offset, err = strconv.Atoi(v)
		default:
		}

		if err != nil {
			s.deps.Log.Error("failed to read query",
				zap.String("field", key), zap.String("value", v))
			return schema.GetRequest{},
				errors.WithMessagef(err, "failed to read query: %s=%s", key, v)
		}
	}

	return ret, nil
}

func (s *Server) getHandler(c *gin.Context) {
	req, err := s.getQuery(c.Request.URL.Query())
	if s.replyError(c, err) {
		return
	}

	res, err := s.deps.Manager.GetPersonInfo(c, req)
	if s.replyError(c, err) {
		return
	}

	c.JSON(http.StatusOK, res)
}

func (s *Server) deleteHandler(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if s.replyError(c, err) {
		s.deps.Log.Error("incorrect ID:", zap.Error(err))
		return
	}

	err = s.deps.Manager.DeletePersonInfo(c, id)
	if s.replyError(c, err) {
		return
	}

	c.Status(http.StatusOK)
}

func (s *Server) updateHandler(c *gin.Context) {
	value := c.Param("id")
	id, err := strconv.Atoi(value)
	if s.replyError(c, err) {
		s.deps.Log.Error("incorrect ID:", zap.Error(err))
		return
	}

	info := schema.PersonInfo{}
	data, err := io.ReadAll(c.Request.Body)
	if s.replyError(c, err) {
		s.deps.Log.Error("failed to read body:", zap.Error(err))
		return
	}

	err = json.Unmarshal(data, &info)
	if s.replyError(c, err) {
		s.deps.Log.Error("failed to unmarshal request:", zap.Error(err))
		return
	}

	info.ID = id
	err = s.deps.Manager.UpdatePersonInfo(c, info)
	if s.replyError(c, err) {
		return
	}

	c.Status(http.StatusOK)
}

type errorResponse struct {
	Message string
}

func (s *Server) replyError(c *gin.Context, err error) bool {
	if err == nil {
		return false
	}

	resp := errorResponse{Message: err.Error()}
	c.JSON(http.StatusBadRequest, &resp)
	return true
}
