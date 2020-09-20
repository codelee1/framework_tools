package utils

import (
	"github.com/afex/hystrix-go/hystrix"
	"sync"
)

// 主要是RequestVolumeThreshold、SleepWindow和ErrorPercentThreshold的配置
var config = hystrix.CommandConfig{
	Timeout:                3000, // 执行command的超时时间，毫秒
	MaxConcurrentRequests:  8,    // command最大并发量
	RequestVolumeThreshold: 5,    // 请求阈值（统计10秒内的请求数量），熔断器是否开启首先要满足这个条件：这里设置5，表示只是要有5个请求，才开始计算ErrorPercentThreshold（错误百分比）
	SleepWindow:            1000, // 过多长时间，熔断器尝试再次检测是否开启，毫秒
	ErrorPercentThreshold:  10,   // 服务错误率，百分比，如果错误率大于这个值，则启动熔断器
}

type Hystrix struct {
	fallback fallbackFunc
	loadMap  *sync.Map
}

type runFunc func() error
type fallbackFunc func(error) error

func NewHystrix(fallback fallbackFunc) *Hystrix {
	return &Hystrix{
		fallback: fallback,
		loadMap:  new(sync.Map),
	}
}

func (h *Hystrix) Run(name string, run runFunc) error {
	if _, ok := h.loadMap.Load(name); !ok {
		hystrix.ConfigureCommand(name, config)
		h.loadMap.Store(name, config)
	}
	err := hystrix.Do(name, func() error {
		return run()
	}, func(err error) error {
		return h.fallback(err)
	})
	return err
}
