# gocache
### 描述   
golang实现的高性能缓存   
#### 简要介绍
1.gocache使用interface作为key和value，这意味着可以以任意类型的键缓存任意类型的值，通常建议key为string或[]byte，gocache底层均使用标准库提供的数据类型。
2.gocache可以以两种方式运行，一种是无容量限制，意味可以存任意数量的缓存，这取决于内存大小；另一种是基于LRU淘汰策略的容量限制，当超过容量限制时会淘汰最久未使用的缓存，可以在运行中调整容量大小，但不得小于当前的缓存数量，避免出现大面积缓存失效导致性能下降。
3.gocache支持超时控制，给缓存一个持续时间，到期后自动删除。（超时缓存同样适用LRU淘汰策略，这可能会出现还没到期的缓存被删除的情况）
### 安装
```bash
$ go get github.com/baytan0720/gocache@latest
```
### 使用
#### 基本使用
```go
package main

import (
	"fmt"

	"github.com/baytan0720/gocache"
)

func main() {
	// 创建一个缓存对象
	c := gocache.New()

	// 写入缓存
	c.Set(1, 1)
	c.Set(2, 2)
	c.Set(3, 3)

	// 获取缓存
	v1, _ := c.Get(1)
	fmt.Println(v1)

	// 修改缓存容量大小
	c.SetCap(4)

	c.Set(4, 4)
	c.Set(5, 5)

	// 查询被淘汰的缓存是否存在
	b := c.Contains(2)
	fmt.Println(b)

	//删除一个或多个缓存
	c.Del(4, 5)

	// 获取缓存大小
	size := c.Size()
	fmt.Println(size)
}

```
#### 超时控制
