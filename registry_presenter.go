package go_validator

import (
"context"
"reflect"
)
type registryPresenter struct {
	registry map[string]PresenterFunc
	fallback PresenterFunc
}

type RegistryPresenter interface {
	Register(err error, presenter PresenterFunc)
	Present(ctx context.Context, path []string, err error) string
}

func NewRegistryPresenter(
	fallback PresenterFunc,
	registry map[error]PresenterFunc,
) RegistryPresenter {
	rp := &registryPresenter{
		fallback: fallback,
	}

	for err, presenter := range registry {
		rp.Register(err, presenter)
	}

	return rp
}

func (rp *registryPresenter) Register(err error, presenter PresenterFunc) {
	if rp.registry == nil {
		rp.registry = make(map[string]PresenterFunc)
	}
	key := reflect.ValueOf(err).Type().String()
	rp.registry[key] = presenter
}

func (rp *registryPresenter) Present(ctx context.Context, path []string, err error) string {
	key := reflect.ValueOf(err).Type().String()
	if fn, ok := rp.registry[key]; ok {
		return fn(ctx, path, err)
	}

	return rp.fallback(ctx, path, err)
}
