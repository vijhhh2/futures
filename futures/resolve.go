package futures

type FutureResolve struct {
	result any
}

func (r *FutureResolve) Poll() (Result, error) {
	return Finished(r.result), nil
}

func (r *FutureResolve) Then(then_fn FutureThenFunc) Future {
	return &FutureThen{
		Left:    r,
		Then_fn: then_fn,
	}
}

func (r *FutureResolve) Catch(catch_fn FutureCatchFunc) Future {
	return &FutureCatch{
		Left:     r,
		Catch_fn: catch_fn,
	}
}

func Resolve(result any) *FutureResolve {
	return &FutureResolve{result: result}
}
