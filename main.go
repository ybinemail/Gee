package main

import (
	"database/sql"
	"fmt"
	"geego/cache"
	"geego/gee"
	"log"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
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

var r = gee.New()

func main() {
	/*
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

		r.Run(":9999")
	*/

	/*
		var port int
		var api bool
		flag.IntVar(&port, "port", 8001, "Geecache server port")
		flag.BoolVar(&api, "api", false, "Start a api server?")
		flag.Parse()

		apiAddr := "http://localhost:9999"
		addrMap := map[int]string{
			8001: "http://localhost:8001",
			8002: "http://localhost:8002",
			8003: "http://localhost:8003",
		}

		var addrs []string
		for _, v := range addrMap {
			addrs = append(addrs, v)
		}

		gee := createGroup()
		if api {
			go startAPIServer(apiAddr, gee)
		}
		startCacheServer(addrMap[port], []string(addrs), gee)
	*/

	demoGeeORM()

}

func demoGeeORM() {
	db, _ := sql.Open("sqlite3", "gee.db")

	defer func() {
		_ = db.Close()
	}()

	_, _ = db.Exec("drop table if exists user;")
	_, _ = db.Exec("CREATE TABLE User(Name text);")

	result, err := db.Exec("insert into user('name') values(?),(?)", "Tom", "Sam")

	if err == nil {
		affected, _ := result.RowsAffected()
		log.Println(affected)
	}

	row := db.QueryRow("select name from user limit 1")
	var name string
	if er := row.Scan(&name); er == nil {
		log.Println(name)
	}

}

func createGroup() *cache.Group {
	return cache.NewGroup("socers", 2<<10, cache.GetterFunc(
		func(key string) ([]byte, error) {
			log.Println("[SlowDB] search key", key)
			if v, ok := db[key]; ok {
				return []byte(v), nil
			}
			return nil, fmt.Errorf("%s not exist", key)
		}))
}

func startCacheServer(addr string, addrs []string, gee *cache.Group) {
	peers := cache.NewHTTPPool(addr)
	peers.Set(addrs...)
	gee.RegisterPeers(peers)
	log.Println("geecache is running at", addr)
	log.Fatal(http.ListenAndServe(addr[7:], peers))
}

func startAPIServer(apiAddr string, gee *cache.Group) {
	http.Handle("/api", http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			key := r.URL.Query().Get("key")
			view, err := gee.Get(key)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/octet-stream")
			w.Write(view.ByteSlice())

		}))
	log.Println("fontend server is running at", apiAddr)
	log.Fatal(http.ListenAndServe(apiAddr[7:], nil))
}
