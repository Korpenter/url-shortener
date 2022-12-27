package encoders

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHMACString(t *testing.T) {
	tests := []struct {
		name   string
		userID string
		want   string
	}{
		{name: "Test HMAC #1",
			userID: "user1",
			want:   "60e8d0babc58e796ac223a64b5e68b998de7d3b203bc8a859bc0ec15ee66f5f9",
		},
		{name: "Test HMAC #2",
			userID: "user2",
			want:   "bfe70caa6f0a26dbc64e5cd31121cb3d5d13075f60b0663b4328375bc3f47456",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, HMACString(tt.userID, []byte("defaultKeyUrlSHoRtenEr")))
		})
	}
}
