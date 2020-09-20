package utils

import (
	"errors"
	"fmt"
	"github.com/afex/hystrix-go/hystrix"
	"math/rand"
	"testing"
	"time"
)

func TestHystrix(t *testing.T) {
	funcName := "testSth"
	fallback := func(err error) error {
		//fmt.Println("todo some  after error was happen ")
		return nil
	}
	hy := NewHystrix(fallback)

	c, _, _ := hystrix.GetCircuit(funcName)

	rand.Seed(time.Now().Unix())
	for i := 0; i < 100; i++ {
		time.Sleep(time.Millisecond * 100)
		err := hy.Run(funcName, func() error {
			// 模拟调用服务失败
			r := rand.Int31n(10)
			if r > 4 {
				return nil
			} else {
				return errors.New("err")
			}
		})
		fmt.Println(i, " 熔断器是否开启：", c.IsOpen(), ". 请求是否允许 :", c.AllowRequest())
		if err != nil {
			fmt.Println("request ", i, " is err:", err)
		}
	}
}
