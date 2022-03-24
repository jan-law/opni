package auth

import (
	"errors"
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/rancher/opni-monitoring/pkg/logger"
)

type Middleware interface {
	Description() string
	Handle(*fiber.Ctx) error
}

type NamedMiddleware interface {
	Middleware
	Name() string
}

type namedMiddlewareImpl struct {
	Middleware
	name string
}

func (nm *namedMiddlewareImpl) Name() string {
	return nm.name
}

func namedMiddleware(name string, mw Middleware) NamedMiddleware {
	return &namedMiddlewareImpl{
		Middleware: mw,
		name:       name,
	}
}

var (
	authMiddlewares = make(map[string]Middleware)

	ErrInvalidMiddlewareName    = errors.New("invalid or empty auth middleware name")
	ErrMiddlewareAlreadyExists  = errors.New("auth middleware already exists")
	ErrNilMiddleware            = errors.New("auth middleware is nil")
	ErrInvalidConfig            = errors.New("middleware config must be a non-nil pointer to a struct")
	ErrMiddlewareNotFound       = errors.New("auth middleware not found")
	ErrMiddlewareConfigNotFound = errors.New("auth middleware config not found")

	authLogger = logger.New().Named("auth")
)

func RegisterMiddleware(name string, m Middleware) error {
	name = strings.TrimSpace(name)
	if len(name) == 0 {
		return ErrInvalidMiddlewareName
	}
	if _, ok := authMiddlewares[name]; ok {
		return fmt.Errorf("%w: %s", ErrMiddlewareAlreadyExists, name)
	}
	if m == nil {
		return ErrNilMiddleware
	}
	authMiddlewares[name] = m
	return nil
}

func GetMiddleware(name string) (NamedMiddleware, error) {
	if m, ok := authMiddlewares[name]; ok {
		return namedMiddleware(name, m), nil
	}
	return nil, fmt.Errorf("%w: %s", ErrMiddlewareNotFound, name)
}

func ResetMiddlewares() {
	authMiddlewares = make(map[string]Middleware)
}

func NamedMiddlewareAs[T Middleware](nmw NamedMiddleware) T {
	return nmw.(*namedMiddlewareImpl).Middleware.(T)
}
