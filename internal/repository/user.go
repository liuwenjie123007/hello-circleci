package repository

import (
	"context"
	"hello-circleci/internal/modle"
	"hello-circleci/pkg/db"
)

type UserRepository interface {
	Get(name string) (*modle.User, error)
	Add(name string) error
	Delete(name string) error
}

func NewUserRepository(db *db.DBContext) UserRepository {
	return userRepository{
		db: db,
	}
}

type userRepository struct {
	db *db.DBContext
}

func (u userRepository) Get(name string) (*modle.User, error) {
	var user modle.User
	if err := u.db.Get(context.Background(), &user, `
		SELECT id, name
		FROM t_user
		WHERE name = ?`, name); err != nil {
		return nil, err
	}

	return &user, nil
}

func (u userRepository) Add(name string) error {
	if _, err := u.db.Exec(context.Background(), `
		INSERT INTO t_user(name) VALUES(?)`, name); err != nil {
		return err
	}

	return nil
}

func (u userRepository) Delete(name string) error {
	if _, err := u.db.Exec(context.Background(), `
		DELETE FROM t_user WHERE name = ?`, name); err != nil {
		return err
	}

	return nil
}
