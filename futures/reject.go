package futures

type FutureReject struct {
	err error
}

func (r *FutureReject) Poll() (Result, error) {
	return Result{}, r.err
}

func (r *FutureReject) Then(then_fn FutureThenFunc) Future {
	return &FutureThen{
		Left:    r,
		Then_fn: then_fn,
	}
}

func (r *FutureReject) Catch(catch_fn FutureCatchFunc) Future {
	return &FutureCatch{
		Left:     r,
		Catch_fn: catch_fn,
	}
}

func Reject(error error) *FutureReject {
	return &FutureReject{err: error}
}
