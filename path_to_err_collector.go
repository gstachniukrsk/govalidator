package go_validator

import "context"
type pathToErrCollector struct {
	errs          map[string][]string
	pathPresenter PresenterFunc
	errPresenter  PresenterFunc
}

type PathToErrCollector interface {
	Collect(ctx context.Context, path []string, err error)
	GetErrors() map[string][]string
}

func NewPathToErrCollector(pathPresenter PresenterFunc, errPresenter PresenterFunc) PathToErrCollector {
	return &pathToErrCollector{
		errs:          make(map[string][]string, 0),
		pathPresenter: pathPresenter,
		errPresenter:  errPresenter,
	}
}

func (c *pathToErrCollector) Collect(ctx context.Context, path []string, err error) {
	key := c.pathPresenter(ctx, path, err)
	c.errs[key] = append(c.errs[key], c.errPresenter(ctx, path, err))
}

func (c *pathToErrCollector) GetErrors() map[string][]string {
	return c.errs
}
