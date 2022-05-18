package main

import (
	"fmt"
	"geego/cache"
	"geego/gee"
	"log"
	"net/http"
	"time"
)

func onlyForV2() gee.HandlerFunc {
	return func(c *gee.Context) {
		// Start timer
		t := time.Now()
		// if a server error occurred
		c.Fail(500, "Internal Server Error")
		// Calculate resolution time
		log.Printf("[%d] %s in %v for group v2", c.StatusCode, c.Req.RequestURI, time.Since(t))
	}
}

var db = map[string]string{
	"Tom":  "630",
	"Jack": "589",
	"Sam":  "567",
}

func main() {
	r := gee.New()

	r.Use(gee.Logger())

	r.GET("/hello", func(c *gee.Context) {
		c.HTML(http.StatusOK, "<h1>Hello Gee</h1>")
	})

	v1 := r.NewGroup("/v1")
	{
		v1.GET("/", func(c *gee.Context) {
			c.HTML(http.StatusOK, "<h1>Hello Gee</h1>")
		})

		v1.GET("/hello", func(c *gee.Context) {
			// expect /hello?name=geektutu
			c.String(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path)
		})
	}
	v2 := r.NewGroup("/v2")
	v2.Use(onlyForV2())
	{
		v2.GET("/hello/:name", func(c *gee.Context) {
			// expect /hello/geektutu
			c.String(http.StatusOK, "hello %s, you're at %s\n", c.Param("name"), c.Path)
		})
		v2.POST("/login", func(c *gee.Context) {
			c.JSON(http.StatusOK, gee.H{
				"username": c.PostForm("username"),
				"password": c.PostForm("password"),
			})
		})

	}

	r.GET("/panic", func(c *gee.Context) {
		names := []string{"geektutu"}
		c.String(http.StatusOK, names[100])
	})
	r.Use(gee.Recovery())
	fmt.Println("Welcome to Gee-go!")

	r.GET("/cache/:group/:key", func(c *gee.Context) {
		groupName := c.Params["group"]
		key := c.Params["key"]

		group := cache.GetGroup(groupName)

		if group == nil {
			c.Fail(http.StatusNotFound, "no such group"+groupName)
		}

		view, err := group.Get(key)
		if err != nil {
			c.Fail(http.StatusInternalServerError, err.Error())
			return
		}

		c.SetHeader("Content-Type", "application/octet-stream")
		c.Data(http.StatusOK, view.ByteSlice())

	})

	cache.NewGroup("socers", 2<<10, cache.GetterFunc(
		func(key string) ([]byte, error) {
			log.Println("[SlowDB] search key", key)
			if v, ok := db[key]; ok {
				return []byte(v), nil
			}
			return nil, fmt.Errorf("%s not exist", key)
		}))

	r.Run(":9999")

}
