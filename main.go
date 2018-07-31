package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"

	"strings"

	"github.com/dlsteuer/slack-snake/commands"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.Use(RequestLogger())
	r.GET("/", func(c *gin.Context) {
		c.String(200, "Hello World!")
	})

	r.POST("/", func(c *gin.Context) {
		fmt.Println(c.PostForm("token"))
		text := c.PostForm("text")
		parts := strings.Split(text, " ")
		f, ok := commands.Mapping[parts[0]]
		if !ok {
			c.String(400, "invalid command: %s", parts[0])
		}
		msg, err := f(parts[1:]...)
		if err != nil {
			fmt.Println(err)
			c.String(400, "error while running command: %s, %v", parts[0], err)
			return
		}
		c.String(200, msg)
	})
	err := r.Run(":7000")
	if err != nil {
		fmt.Println(err)
	}
}

func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		buf, _ := ioutil.ReadAll(c.Request.Body)
		rdr1 := ioutil.NopCloser(bytes.NewBuffer(buf))
		rdr2 := ioutil.NopCloser(bytes.NewBuffer(buf)) //We have to create a new Buffer, because rdr1 will be read.

		fmt.Println(readBody(rdr1)) // Print request body

		c.Request.Body = rdr2
		c.Next()
	}
}

func readBody(reader io.Reader) string {
	buf := new(bytes.Buffer)
	buf.ReadFrom(reader)

	s := buf.String()
	return s
}
