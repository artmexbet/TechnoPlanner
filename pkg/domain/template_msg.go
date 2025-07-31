package domain

type Envelope[P any] struct {
	ID      string
	Payload P
}
