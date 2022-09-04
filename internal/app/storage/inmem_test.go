package storage

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"sync"
	"testing"
)

func TestInMemRepo_Add(t *testing.T) {
	type fields struct {
		urls map[string]string
		*sync.RWMutex
	}
	type args struct {
		longURL string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name:   "Id in repo",
			fields: fields{urls: make(map[string]string)},
			args:   args{longURL: "https://github.com/"},
			want:   "1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &InMemRepo{
				urls: tt.fields.urls,
			}
			assert.Equal(t, tt.want, r.Add(tt.args.longURL))
			assert.Contains(t, r.urls, tt.want)
		})
	}
}

func TestInMemRepo_Get(t *testing.T) {
	type fields struct {
		urls map[string]string
		*sync.RWMutex
	}
	type args struct {
		id string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "Id in repo",
			fields: fields{
				urls: map[string]string{
					"1": "https://github.com/",
					"2": "https://github.com/",
				},
			},
			args:    args{id: "2"},
			want:    "https://github.com/",
			wantErr: false,
		},
		{
			name: "Id not in repo",
			fields: fields{
				urls: map[string]string{
					"1": "https://github.com/",
					"2": "https://github.com/",
				},
			},
			args:    args{id: "3"},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &InMemRepo{
				urls: tt.fields.urls,
			}
			got, err := r.Get(tt.args.id)
			if !tt.wantErr {
				require.NoError(t, err)
				assert.Equal(t, tt.want, got)
				return
			}
			assert.Error(t, err)
		})
	}
}

func TestInMemRepo_NewID(t *testing.T) {
	type fields struct {
		urls map[string]string
		*sync.RWMutex
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		{
			name: "Test #1",
			fields: fields{
				urls: map[string]string{
					"1": "https://github.com/",
					"2": "https://github.com/",
				},
			},
			want: 3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &InMemRepo{
				urls: tt.fields.urls,
			}
			assert.Equal(t, tt.want, r.NewID())
		})
	}
}
