package msgfmt_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/vporoshok/msgfmt"
)

type MsgfmtSuite struct {
	suite.Suite
}

func TestMsgfmt(t *testing.T) {
	suite.Run(t, new(MsgfmtSuite))
}

func (s *MsgfmtSuite) TestPlainText() {
	p := msgfmt.New()
	msg, err := p.Parse("some text here")
	s.Require().NoError(err)

	ctx := context.Background()
	s.Equal("some text here", msg.Format(ctx, nil))
}

func (s *MsgfmtSuite) TestBlockEscape() {
	p := msgfmt.New()
	msg, err := p.Parse("some '{text}' here")
	s.Require().NoError(err)

	ctx := context.Background()
	s.Equal("some {text} here", msg.Format(ctx, nil))
}

func (s *MsgfmtSuite) TestEscapeUntilEnd() {
	p := msgfmt.New()
	msg, err := p.Parse("some '{text} here")
	s.Require().NoError(err)

	ctx := context.Background()
	s.Equal("some {text} here", msg.Format(ctx, nil))
}

func (s *MsgfmtSuite) TestQuote() {
	p := msgfmt.New()
	msg, err := p.Parse("some text '' here")
	s.Require().NoError(err)

	ctx := context.Background()
	s.Equal("some text ' here", msg.Format(ctx, nil))
}

func (s *MsgfmtSuite) TestQuoteAtEnd() {
	p := msgfmt.New()
	msg, err := p.Parse("some text here''")
	s.Require().NoError(err)

	ctx := context.Background()
	s.Equal("some text here'", msg.Format(ctx, nil))
}

func (s *MsgfmtSuite) TestQuoteInEscape() {
	p := msgfmt.New()
	msg, err := p.Parse("some 'text '' here'")
	s.Require().NoError(err)

	ctx := context.Background()
	s.Equal("some text ' here", msg.Format(ctx, nil))
}

func (s *MsgfmtSuite) TestVariable() {
	p := msgfmt.New()
	msg, err := p.Parse("some {foo} here")
	s.Require().NoError(err)

	ctx := context.Background()
	s.Equal("some text here", msg.Format(ctx, msgfmt.D{"foo": "text"}))
}

func (s *MsgfmtSuite) TestBreakControl() {
	p := msgfmt.New()
	msg, err := p.Parse("some {foo")
	s.Require().NoError(err)

	ctx := context.Background()
	s.Equal("some ", msg.Format(ctx, msgfmt.D{"foo": "text"}))
}

func (s *MsgfmtSuite) TestEmptyKey() {
	p := msgfmt.New()
	msg, err := p.Parse("some { , } bar")
	s.Error(err)

	ctx := context.Background()
	s.Equal("some ", msg.Format(ctx, msgfmt.D{"foo": "text"}))
}

func (s *MsgfmtSuite) TestBadControl() {
	p := msgfmt.New()
	msg, err := p.Parse("some {foo bar}")
	s.Error(err)

	ctx := context.Background()
	s.Equal("some ", msg.Format(ctx, msgfmt.D{"foo": "text"}))
}
