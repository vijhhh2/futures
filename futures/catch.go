package futures

import "github.com/practice/utils"

type FutureCatchFunc func(reslut error) Future

type FutureCatch struct {
	Left     Future
	Right    Future
	Catch_fn FutureCatchFunc
}

func (t *FutureCatch) Poll() (Result, error) {
	if t.Left != nil {
		result, err := t.Left.Poll()
		if err != nil {
			t.Right = t.Catch_fn(err)
			t.Left = nil
			return Pending(), nil
		}
		return result, nil
	} else {
		utils.Assert(t.Right, "Right should be a future")
		return t.Right.Poll()
	}
}

func (t *FutureCatch) Then(then_fn FutureThenFunc) Future {
	return &FutureThen{
		Left:    t,
		Then_fn: then_fn,
	}
}

func (t *FutureCatch) Catch(catch_fn FutureCatchFunc) Future {
	return &FutureCatch{
		Left:     t,
		Catch_fn: catch_fn,
	}
}