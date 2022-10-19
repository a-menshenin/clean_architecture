package postgres

import (
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Repository struct {
	db      *pgxpool.Pool
	genSQL  squirrel.StatementBuilderType
	options options
}

type options struct {
	Timeout       time.Duration
	DefaultLimit  uint64
	DefaultOffset uint64
}

type Option func(*options)

func WithTimeout(timeout time.Duration) Option {
	return func(args *options) {
		args.Timeout = timeout
	}
}

func WithDefaultLimit(limit uint64) Option {
	return func(args *options) {
		args.DefaultLimit = limit
	}
}

func WithDefaultOffset(offset uint64) Option {
	return func(args *options) {
		args.DefaultOffset = offset
	}
}

func (r *Repository) SetOptions(setters ...Option) {
	args := &options{
		Timeout:       time.Second * 30,
		DefaultLimit:  10,
		DefaultOffset: 0,
	}

	for _, setter := range setters {
		setter(args)
	}

	r.options = *args
}

func New(db *pgxpool.Pool, setters ...Option) *Repository {
	var r = &Repository{
		genSQL: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
		db:     db,
	}

	r.SetOptions(setters...)
	return r
}
