package conf

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
)

// ParseFunc is a type for the parsing function
type ParseFunc func(ctx context.Context, r io.Reader) (any, error)

// Parser is an extension for the Reader interface.
type Parser interface {
	Reader

	WithParser(parser ParseFunc) Parser
	WithPrefix(prefix string) Parser
}

type parser struct {
	stream io.Reader
	parser ParseFunc
	prefix string
}

func (p *parser) Prefix() string {
	return p.prefix
}

var (
	// ErrNoParser is an error returned if the parser was not given
	ErrNoParser = errors.New("no parser")
	// ErrNoStream is an error returned if the stream was not given
	ErrNoStream = errors.New("no data stream")
)

func (p *parser) Read(ctx context.Context) (any, error) {
	if p.stream == nil {
		return nil, ErrNoStream
	}
	if p.parser == nil {
		return nil, ErrNoParser
	}

	data, err := p.parser(ctx, p.stream)
	if err != nil {
		return nil, fmt.Errorf("parser %T failed: %w", p.parser, err)
	}

	if v, ok := p.stream.(io.Closer); ok {
		if err := v.Close(); err != nil {
			return nil, fmt.Errorf("failed to close stream: %w", err)
		}
	}

	return data, nil
}

func (p *parser) WithPrefix(prefix string) Parser {
	p.prefix = prefix
	return p
}

func (p *parser) WithParser(parser ParseFunc) Parser {
	p.parser = parser
	return p
}

// NewStreamParser creates an instance of the Parser to read from a given stream
func NewStreamParser(stream io.Reader) Parser {
	return &parser{
		stream: stream,
	}
}

// NewFileParser creates an instance of the Parser and opens the given file
func NewFileParser(filename string) (Parser, error) {
	f, err := os.Open(filename) //nolint:gosec
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}

	return NewStreamParser(f), nil
}
