package msgfmt

import "context"

type D map[string]interface{}

type Message interface {
	Format(context.Context, D) string
}

type Parser interface {
	Parse(string) (Message, error)
}

func New(opts ...Option) Parser {
	p := &parser{}

	for _, opt := range opts {
		opt.apply(p)
	}

	return p
}
