package example

//
// import (
// 	"errors"
// 	"fmt"
//
// 	"github.com/practice/futures"
// )
//
// func Test() {
// 	allFutures := []futures.Future{}
// 	i := 0
//
// 	exFuture := futures.Poll(func(state any) futures.Result {
// 		i, ok := state.(*int)
// 		if !ok {
// 			panic("Some other type")
// 		}
// 		fmt.Println("count", *i)
// 		if *i < 10 {
// 			*i += 1
// 			return futures.Pending()
// 		} else {
// 			return futures.Finished("Done counting")
// 		}
// 	}, &i).Then(func(reslut any) futures.Future {
// 		fmt.Println(reslut)
// 		return futures.Resolve(nil)
// 	})
//
// 	exFuture2 := futures.Resolve(20).Then(func(reslut any) futures.Future {
// 		fmt.Println(reslut)
// 		return futures.Resolve(30)
// 	}).Then(func(reslut any) futures.Future {
// 		fmt.Println(reslut)
// 		return futures.Resolve(40)
// 	}).Then(func(reslut any) futures.Future {
// 		fmt.Println(reslut)
// 		return futures.Reject(errors.New("throwing some random error"))
// 	}).Catch(func(reslut error) futures.Future {
// 		fmt.Println("Recovering from error", reslut)
// 		return futures.Resolve(nil)
// 	})
//
// 	allFutures = append(allFutures, exFuture, exFuture2)
//
// 	Run(allFutures)
// }
//
// func Run(futures []futures.Future) {
// 	quit := false
// 	for !quit {
// 		quit = true
// 		for _, f := range futures {
// 			res, _ := f.Poll()
// 			if res.Value != nil {
// 				fmt.Println(res.Value)
// 			}
// 			if !res.Finished {
// 				quit = false
// 			}
// 		}
// 	}
// }
