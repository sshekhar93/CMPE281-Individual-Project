package hash

import (
	"github.com/catinello/base62"
)

type hashBase62 struct{}

func NewHashBase62() HashALgorithm {
	return &hashBase62{}
}

func (*hashBase62) Encode(num int) string {
	return base62.Encode(num)
}
func (*hashBase62) Decode(s string) (int, error) {
	return base62.Decode(s)
}
