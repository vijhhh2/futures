package utils

func Swap[T any](a T, b T) {
	b, a = a, b
}
