package main

import (
	"fmt"
	"time"

	"gocache"
)

func main() {
	// 创建一个缓存对象
	c := gocache.New()

	// 写入缓存并设置超时
	c.SetWithTimeout(1, 1, 1*time.Second)
	c.GetOrSet(2, 2)

	// 获取缓存以及剩余时间
	v1, _ := c.Get(1)
	fmt.Println(v1)
	timeout, _ := c.GetTimeOut(1)
	fmt.Println(timeout)

	c.GetOrSet(2, 2)

	// 等待超时
	time.Sleep(time.Second)

	//打印所有key
	keys := c.Keys()
	fmt.Println(keys)
}
