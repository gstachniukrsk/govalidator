package validator

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPathPresenter(t *testing.T) {
	type args struct {
		glue string
	}
	tests := []struct {
		name  string
		args  args
		input []string
		want  string
	}{
		{
			name: "empty",
			args: args{
				glue: ".",
			},
			input: []string{},
			want:  "",
		},
		{
			name: "one",
			args: args{
				glue: ".",
			},
			input: []string{"one"},
			want:  "one",
		},
		{
			name: "two",
			args: args{
				glue: ".",
			},
			input: []string{"one", "two"},
			want:  "one.two",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := PathPresenter(tt.args.glue)
			got := p(nil, tt.input, nil)
			assert.Equal(t, tt.want, got)
		})
	}
}
