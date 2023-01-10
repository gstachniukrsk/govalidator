package go_validator_test

import (
	"context"
	"github.com/gstachniukrsk/go_validator"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSimpleErrorPresenter(t *testing.T) {
	type args struct {
		ctx  context.Context
		path []string
		err  error
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "required error",
			args: args{
				ctx:  context.Background(),
				path: []string{},
				err:  go_validator.RequiredError{},
			},
			want: "required",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := go_validator.SimpleErrorPresenter()

			out := p(tt.args.ctx, tt.args.path, tt.args.err)
			assert.Equal(t, tt.want, out)
		})
	}
}
