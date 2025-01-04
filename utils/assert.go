package utils

func Assert(value any, message string) {
	if value == nil {
		panic(message)
	}
}
