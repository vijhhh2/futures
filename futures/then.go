package futures

import "github.com/practice/utils"

type FutureThenFunc func(reslut any) Future

type FutureThen struct {
	Left    Future
	Right   Future
	Then_fn FutureThenFunc
}

func (t *FutureThen) Poll() (Result, error) {
	if t.Left != nil {
		result, err := t.Left.Poll()
		if err != nil {
			return Result{}, err
		}
		if result.Finished {
			t.Right = t.Then_fn(result.Value)
			t.Left = nil
		}
		return Pending(), nil
	} else {
		utils.Assert(t.Right, "Right should be a future")
		return t.Right.Poll()
	}
}

func (t *FutureThen) Then(then_fn FutureThenFunc) Future {
	return &FutureThen{
		Left:    t,
		Then_fn: then_fn,
	}
}

func (t *FutureThen) Catch(catch_fn FutureCatchFunc) Future {
	return &FutureCatch{
		Left:     t,
		Catch_fn: catch_fn,
	}
}