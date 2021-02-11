package msgfmt

import (
	"errors"
	"fmt"
	"io"
)

type parser struct{}

func (p *parser) Parse(s string) (Message, error) {
	root, err := p.readRoot(newReader(s))
	if errors.Is(err, io.EOF) {
		err = nil
	}

	return root, err
}

func (p *parser) readRoot(r *reader) (root rootNode, err error) {
	var node baseNode

	for {
		var next rune
		node, next, err = p.readText(r)
		root = append(root, node)

		if err != nil {
			return
		}

		switch next {
		case '{':
			node, err = p.readControl(r)
		case '\'':
			node, err = p.readEscape(r)
		}

		root = append(root, node)

		if err != nil {
			return
		}
	}
}

func (p *parser) readText(r *reader) (textNode, rune, error) {
	condition, result := anyOf('{', '\'')
	text, err := r.ReadUntil(condition)

	return textNode(text), result(), err
}

func (p *parser) readEscape(r *reader) (rootNode, error) {
	var root rootNode

	for {
		text, err := r.ReadUntil(equal('\''))
		if err != nil {
			return append(root, textNode(text)), err
		}

		if text != "" {
			root = append(root, textNode(text))
		} else {
			root = append(root, quoteNode{})
		}

		next, _, err := r.ReadRune()
		if next != '\'' || err != nil {
			_ = r.UnreadRune()
			return root, err
		}

		root = append(root, quoteNode{})
	}
}

func (p *parser) readControl(r *reader) (baseNode, error) {
	if err := p.skipWhitespaces(r); err != nil {
		return rootNode{}, err
	}

	key, err := r.ReadWhile(keyword())
	if err != nil {
		return rootNode{}, err
	}

	if key == "" {
		return rootNode{}, ErrEmptyKey
	}

	_ = r.UnreadByte()

	if err := p.skipWhitespaces(r); err != nil {
		return rootNode{}, err
	}

	switch next, _, _ := r.ReadRune(); next {
	case '}':
		return variableNode(key), nil
	case ',':
		return nil, fmt.Errorf("%w: complex control", ErrNotImplemented)
	default:
		return rootNode{}, ErrInvalidTemplate
	}
}

func (*parser) skipWhitespaces(r *reader) error {
	_, err := r.ReadWhile(whitespace())
	if err != nil {
		return err
	}

	return r.UnreadRune()
}
