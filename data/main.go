package main

import (
	"fmt"
	"sync"
)

// 数据越界
// func main() {
// 	var a uint8 = 255
// 	fmt.Printf("%d\n", a+1)
// 	fmt.Printf("%d\n", a+2)
// 	fmt.Printf("%d\n", a+3)
// }

// func main() {
// 	go func() { // goroutine协程
// 		defer func() {
// 			if e := recover(); e != nil { // 捕获该协程的panic 0000000
// 				fmt.Println("recover ", e)
// 			}
// 		}()
// 		panic("0000000")
// 		fmt.Println("11111111") // 无法打印，panic之后就终止了不再执行后续的逻辑
// 	}()

// 	defer func() {
// 		if e := recover(); e != nil { // 捕获main协程的panic 222222
// 			fmt.Println("recover ", e)
// 		}
// 	}()
// 	panic("222222 ")            // mian 协程的panic
// 	fmt.Println("33333")        // panic之后就终止了不再执行后续的逻辑
// 	time.Sleep(2 * time.Second) // 保证运行gofunc协程
// }

func callGoroutine(handlers ...func() error) (err error) {
	var wg sync.WaitGroup
	for _, f := range handlers {
		wg.Add(1)
		// 每个函数启动一个协程
		go func(handler func() error) {
			defer func() {
				// 每个协程内部使用recover捕获可能在调用逻辑中发生的panic
				if e := recover(); e != nil {
					// 日志记录失败信息，捕获panic，不影响其他协程跟主协程运行
					fmt.Println("recover ", e)
				}
				defer wg.Done()
			}()
			// 取第一个报错的handler调用逻辑，并最终向外返回
			e := handler()
			if err == nil && e != nil {
				err = e
			}
		}(f)
	}

	wg.Wait()

	return
}

func main() {
	userRpc := func() error {
		panic("userRpc fail ")
		return nil
	}

	// 调用逻辑2
	orderRpc := func() error {
		panic("orderRpc fail")
		return nil
	}

	err := callGoroutine(userRpc, orderRpc)
	if err != nil {
		fmt.Println(err)
	}
}
