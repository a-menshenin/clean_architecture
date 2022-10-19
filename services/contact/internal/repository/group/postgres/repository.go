package postgres

import (
	"fmt"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/spf13/viper"

	"github.com/pressly/goose"

	"architecture_go_2/services/contact/internal/repository/contact"
)

func init() {
	viper.SetDefault("MIGRATIONS_DIR", "./services/contact/internal/repository/storage/postgres/migrations")
}

type Repository struct {
	db     *pgxpool.Pool
	genSQL squirrel.StatementBuilderType

	repoContact contact.Contact

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

func New(db *pgxpool.Pool, repoContact contact.Contact, setters ...Option) (*Repository, error) {
	if err := migrations(db); err != nil {
		return nil, err
	}

	var r = &Repository{
		genSQL:      squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
		repoContact: repoContact,
		db:          db,
	}

	r.SetOptions(setters...)
	return r, nil
}

func migrations(pool *pgxpool.Pool) (err error) {
	db, err := goose.OpenDBWithDriver("postgres", pool.Config().ConnConfig.ConnString())
	if err != nil {
		return err
	}
	defer func() {
		if errClose := db.Close(); errClose != nil {
			err = errClose
			return
		}
	}()

	dir := viper.GetString("MIGRATIONS_DIR")
	goose.SetTableName("contact_version")
	if err = goose.Run("up", db, dir); err != nil {
		return fmt.Errorf("goose %s error : %w", "up", err)
	}
	return
}
