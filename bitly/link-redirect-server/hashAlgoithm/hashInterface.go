package hash

type HashALgorithm interface {
	Encode(num int) string
	Decode(s string) (int, error)
}
