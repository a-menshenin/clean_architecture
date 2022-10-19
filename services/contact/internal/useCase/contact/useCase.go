package contact

import (
	"time"

	"architecture_go_2/services/contact/internal/useCase/adapters/storage"
)

type UseCase struct {
	adapterStorage storage.Contact
	options        options
}

type options struct {
	Timeout time.Duration
}

type Option func(*options)

func WithTimeout(timeout time.Duration) Option {
	return func(args *options) {
		args.Timeout = timeout
	}
}

func (uc *UseCase) SetOptions(setters ...Option) {
	args := &options{
		Timeout: time.Second * 30,
	}

	for _, setter := range setters {
		setter(args)
	}

	uc.options = *args
}

func New(storage storage.Contact, setters ...Option) *UseCase {
	var uc = &UseCase{
		adapterStorage: storage,
	}
	uc.SetOptions(setters...)
	return uc
}
