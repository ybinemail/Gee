package gee

import (
	"log"
	"time"
)

func Logger() HandlerFunc {
	return func(c *Context) {
		startTime := time.Now()

		//Process next
		c.Next()
		// calclate resolution time
		log.Printf("[%d]  %s  in  %v", c.StatusCode, c.Req.RequestURI, time.Since(startTime))
	}
}
