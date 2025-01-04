package futures

type FutureForeverFunc func(state any) Future

type FutureForever struct {
	State            any
	Future_forver_fn FutureForeverFunc
	Future           Future
}

func (f *FutureForever) Poll() (Result, error) {
	if f.Future == nil {
		f.Future = f.Future_forver_fn(f.State)
	}

	res, err := f.Future.Poll()
	if err != nil {
		f.Future = nil
		return Result{}, err
	}

	if res.Finished {
		f.Future = nil
	}

	return Pending(), nil
}

func (t *FutureForever) Then(then_fn FutureThenFunc) Future {
	return &FutureThen{
		Left:    t,
		Then_fn: then_fn,
	}
}

func (t *FutureForever) Catch(catch_fn FutureCatchFunc) Future {
	return &FutureCatch{
		Left:     t,
		Catch_fn: catch_fn,
	}
}

func Forever(forever_fn FutureForeverFunc, state any) Future {
	return &FutureForever{
		State:            state,
		Future_forver_fn: forever_fn,
	}
}
