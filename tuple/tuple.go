package tuple

type Pair[A, B any] struct {
	First  A
	Second B
}

type Triple[A, B, C any] struct {
	First  A
	Second B
	Third  C
}
