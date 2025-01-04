package futures

type FuturePollFunc func(state any) (Result, error)

type FuturePoll struct {
	state   any
	poll_fn FuturePollFunc
}

func (p *FuturePoll) Poll() (Result, error) {
	return p.poll_fn(p.state)
}

func (p *FuturePoll) Then(then_fn FutureThenFunc) Future {
	return &FutureThen{
		Left:    p,
		Then_fn: then_fn,
	}
}

func (p *FuturePoll) Catch(catch_fn FutureCatchFunc) Future {
	return &FutureCatch{
		Left:     p,
		Catch_fn: catch_fn,
	}
}

func Poll[S any](poll_fn FuturePollFunc, state S) Future {
	return &FuturePoll{
		state:   state,
		poll_fn: poll_fn,
	}
}

