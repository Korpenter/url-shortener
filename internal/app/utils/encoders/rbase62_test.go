package encoders

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestToBase62(t *testing.T) {
	tests := []struct {
		name string
		num  int
		want string
	}{
		{name: "Test #1",
			num:  1243,
			want: "3K",
		},
		{name: "Test #1",
			num:  53467,
			want: "NuD",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, ToRBase62(tt.num))
		})
	}
}
