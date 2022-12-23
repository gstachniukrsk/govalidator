package main_test

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"validator"
)

func TestNewPathToErrCollector(t *testing.T) {
	t.Run("new path to err collector", func(t *testing.T) {
		out := main.NewPathToErrCollector(main.PathPresenter("."), main.SimpleErrorPresenter())
		assert.NotNil(t, out)
	})
}

func Test_pathToErrCollector_Collect(t *testing.T) {
	t.Run("collect, it's not a set", func(t *testing.T) {
		c := main.NewPathToErrCollector(main.PathPresenter("."), main.SimpleErrorPresenter())
		c.Collect(context.Background(), []string{}, main.RequiredError{})
		c.Collect(context.Background(), []string{}, main.RequiredError{})
		assert.NotNil(t, c)

		out := c.GetErrors()

		assert.Equal(t, 1, len(out))

		root := out[""]

		assert.Equal(t, "required", root[0])
		assert.Equal(t, "required", root[1])
	})

	t.Run("two fields", func(t *testing.T) {
		c := main.NewPathToErrCollector(main.PathPresenter("."), main.SimpleErrorPresenter())
		c.Collect(context.Background(), []string{"a"}, main.RequiredError{})
		c.Collect(context.Background(), []string{"b"}, main.RequiredError{})
		assert.NotNil(t, c)

		out := c.GetErrors()

		assert.Equal(t, 2, len(out))

		a := out["a"]
		b := out["b"]

		assert.Equal(t, "required", a[0])
		assert.Equal(t, "required", b[0])
	})

	t.Run("nothing collected", func(t *testing.T) {
		c := main.NewPathToErrCollector(main.PathPresenter("."), main.SimpleErrorPresenter())
		assert.NotNil(t, c)

		out := c.GetErrors()

		require.NotNil(t, out)
		assert.Equal(t, 0, len(out))
	})
}
