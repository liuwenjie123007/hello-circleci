package db

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type DBContext struct {
	db *sqlx.DB
}

type Config struct {
	User     string
	Password string
	URI      string // "(domain:port)/database"
}

func NewDB(config *Config) (*DBContext, error) {
	db, err := sqlx.Connect("mysql", fmt.Sprintf("%s:%s@%s?parseTime=true&loc=Asia%%2FTokyo", config.User, config.Password, config.URI))
	d := &DBContext{
		db: db,
	}

	return d, err
}

func NewDBWithDB(db *sqlx.DB) *DBContext {
	return &DBContext{
		db: db,
	}
}

func (dc *DBContext) Begin(ctx context.Context) (context.Context, error) {
	if _, ok := getTx(ctx); ok {
		return ctx, nil
	}
	tx, err := dc.db.Beginx()
	if err != nil {
		return nil, fmt.Errorf("begin transaction error: %w", err)
	}

	return withTx(ctx, tx), nil
}

func (dc *DBContext) Commit(ctx context.Context) error {
	tx, ok := getTx(ctx)
	if !ok {
		return nil
	}
	return tx.Commit()
}

func (dc *DBContext) Rollback(ctx context.Context) error {
	tx, ok := getTx(ctx)
	if !ok {
		return nil
	}
	return tx.Rollback()
}

// Select is get []interface{}
func (dc *DBContext) Select(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	tx, ok := getTx(ctx)
	if ok {
		return tx.SelectContext(ctx, dest, query, args...)
	}
	return dc.db.SelectContext(ctx, dest, query, args...)
}

// Get is get interface{}
func (dc *DBContext) Get(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	tx, ok := getTx(ctx)
	if ok {
		return tx.GetContext(ctx, dest, query, args...)
	}
	return dc.db.GetContext(ctx, dest, query, args...)
}

// Exec wrapper
func (dc *DBContext) Exec(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	tx, ok := getTx(ctx)
	if ok {
		return tx.ExecContext(ctx, query, args...)
	}
	return dc.db.ExecContext(ctx, query, args...)
}

// NamedExec wrapper
func (dc *DBContext) NamedExec(ctx context.Context, query string, arg interface{}) (sql.Result, error) {
	tx, ok := getTx(ctx)
	if ok {
		return tx.NamedExecContext(ctx, query, arg)
	}
	return dc.db.NamedExecContext(ctx, query, arg)
}

// QueryContext wrapper
func (dc *DBContext) QueryContext(ctx context.Context, query string, args ...interface{}) (*sqlx.Rows, error) {
	tx, ok := getTx(ctx)
	if ok {
		return tx.QueryxContext(ctx, query, args...)
	}
	return dc.db.QueryxContext(ctx, query, args...)
}

var txKey = struct{}{}

func getTx(ctx context.Context) (*sqlx.Tx, bool) {
	t := ctx.Value(txKey)
	if t == nil {
		return nil, false
	}
	tx, ok := t.(*sqlx.Tx)
	return tx, ok
}

func withTx(ctx context.Context, tx *sqlx.Tx) context.Context {
	return context.WithValue(ctx, txKey, tx)
}
