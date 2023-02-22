package repository_test

import (
	"github.com/stretchr/testify/assert"
	"hello-circleci/internal/modle"
	"hello-circleci/internal/repository"
	"hello-circleci/internal/testutil"
	"testing"
)

func TestGet(t *testing.T) {
	db, f := testutil.TearUp(t)
	defer f()

	ur := repository.NewUserRepository(db)

	tests := []struct {
		name    string
		input   string
		want    *modle.User
		wantErr error
	}{
		{
			name:  "add",
			input: "hoge",
			want: &modle.User{
				Id:   1,
				Name: "hoge",
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		user, err := ur.Get(tt.input)
		if tt.wantErr == nil {
			if assert.NoError(t, err) {
				assert.Equal(t, tt.want, user)
			}
		} else {
			assert.Error(t, err)
		}
	}
}
