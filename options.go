package msgfmt

type Option interface {
	apply(*parser)
}

// type fnOption func(*parser)

// func (fn fnOption) apply(p *parser) {
// 	fn(p)
// }
