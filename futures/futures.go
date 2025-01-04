package futures

type Result struct {
	Finished bool
	Value    any
}

func Pending() Result {
	return Result{}
}

func Finished(value any) Result {
	return Result{
		Finished: true,
		Value:    value,
	}
}

// Resolve a future similar to javascript Resolve promise
type Future interface {
	Poll() (Result, error)
	Then(then_fn FutureThenFunc) Future
	Catch(catch_fn FutureCatchFunc) Future
}