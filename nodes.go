package msgfmt

import (
	"context"
	"fmt"
	"io"
	"strings"
)

type baseNode interface {
	String() string
	format(context.Context, io.Writer, D)
}

type rootNode []baseNode

func (n rootNode) String() string {
	return fmt.Sprintf("%v", []baseNode(n))
}

func (n rootNode) Format(ctx context.Context, d D) string {
	res := new(strings.Builder)
	n.format(ctx, res, d)

	return res.String()
}

func (n rootNode) format(ctx context.Context, w io.Writer, d D) {
	for i := range n {
		n[i].format(ctx, w, d)
	}
}

type textNode string

const textNodeWrap = 10

func (n textNode) String() string {
	if len(n) > textNodeWrap {
		return fmt.Sprintf("%.10s...", string(n))
	}

	return string(n)
}

func (n textNode) format(_ context.Context, w io.Writer, _ D) {
	_, _ = w.Write([]byte(n))
}

type quoteNode struct{}

func (quoteNode) String() string {
	return "'"
}

func (quoteNode) format(_ context.Context, w io.Writer, _ D) {
	_, _ = w.Write([]byte{'\''})
}

type variableNode string

func (n variableNode) String() string {
	return string(n)
}

func (n variableNode) format(_ context.Context, w io.Writer, d D) {
	_, _ = fmt.Fprint(w, d[string(n)])
}
