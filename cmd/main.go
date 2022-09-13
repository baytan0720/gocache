package main

import (
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"gocache"

	"github.com/labstack/echo/v4"
)

var C *gocache.Cache

func main() {
	if len(os.Args) < 2 {
		log.Fatalln("Usage: ./xxx 8080")
	}
	port := ":" + os.Args[1]
	C = gocache.New()
	s := echo.New()

	s.GET("/set", set)
	s.GET("/setwithtimeout", setWithTimeout)
	s.GET("/setifnotexist", setIfNotExist)
	s.GET("/setcap", setCap)
	s.GET("/get", get)
	s.GET("/getorset", getOrSet)
	s.GET("/gettimeout", getTimeout)
	s.GET("/del", del)
	s.GET("/size", size)
	s.GET("/keys", keys)
	s.GET("/vals", vals)
	s.GET("/entrys", entrys)

	s.Start(port)
}

func set(c echo.Context) error {
	key := c.Request().FormValue("key")
	val := c.Request().FormValue("val")
	if key == "" {
		return c.String(200, "Invalid argument: key cannot be null")
	}
	C.Set(key, val)
	return c.String(200, "OK")
}

func setWithTimeout(c echo.Context) error {
	key := c.Request().FormValue("key")
	val := c.Request().FormValue("val")
	timeout, err := time.ParseDuration(c.Request().FormValue("timeout"))
	if err != nil {
		return c.String(200, "Invalid argument: timeout unavailable")
	}
	if key == "" {
		return c.String(200, "Invalid argument: key cannot be null")
	}
	C.SetWithTimeout(key, val, timeout)
	return c.String(200, "OK")
}

func setIfNotExist(c echo.Context) error {
	key := c.Request().FormValue("key")
	val := c.Request().FormValue("val")
	if key == "" {
		return c.String(200, "Invalid argument: key cannot be null")
	}
	if ok := C.SetIfNotExist(key, val); !ok {
		return c.String(200, "Exist")
	}

	return c.String(200, "OK")
}

func setCap(c echo.Context) error {
	cap, err := strconv.Atoi(c.Request().FormValue("cap"))
	if err != nil {
		return c.String(200, "Invalid argument: cap unavailable")
	}
	C.SetCap(cap)
	return c.String(200, "OK")
}

func get(c echo.Context) error {
	key := c.Request().FormValue("key")
	if key == "" {
		return c.String(200, "Invalid argument: key cannot be null")
	}
	if v, ok := C.Get(key); ok {
		return c.JSON(200, v)
	}
	return c.JSON(200, nil)
}

func getOrSet(c echo.Context) error {
	key := c.Request().FormValue("key")
	val := c.Request().FormValue("val")
	if key == "" {
		return c.String(200, "Invalid argument: key cannot be null")
	}
	v := C.GetOrSet(key, val)
	return c.JSON(200, v)
}

func getTimeout(c echo.Context) error {
	key := c.Request().FormValue("key")
	if key == "" {
		return c.String(200, "Invalid argument: key cannot be null")
	}
	timeout, ok := C.GetTimeOut(key)
	if !ok {
		return c.JSON(200, nil)
	}
	if timeout == -1 {
		c.JSON(200, -1)
	}
	return c.JSON(200, timeout)
}

func del(c echo.Context) error {
	key := c.Request().FormValue("key")
	keys := strings.Split(c.Request().FormValue("keys"), ",")
	if key == "" && len(keys) == 0 {
		return c.String(200, "Invalid argument: key(s) cannot be null")
	}
	Keys := make([]interface{}, 0, len(keys)+1)
	for _, v := range keys {
		if v != "" {
			keys = append(keys, v)
		}
	}
	if key != "" {
		Keys = append(Keys, key)
	}
	C.Del(Keys...)
	return c.String(200, "OK")
}

func size(c echo.Context) error {
	return c.JSON(200, C.Size())
}

func keys(c echo.Context) error {
	return c.JSON(200, C.Keys())
}

func vals(c echo.Context) error {
	return c.JSON(200, C.Vals())
}

func entrys(c echo.Context) error {
	return c.JSON(200, C.Entrys())
}
