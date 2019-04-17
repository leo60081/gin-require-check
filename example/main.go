package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/leo60081/gin-require-check"
)

type parameter struct {
	PersonID string `parameter:"person-id" binding:"required"` //must match with gin relative path design.
}
type header struct {
	Sing      string `header:"user-sing" binding:"required"`
	TimeStamp int    `header:"timestamp" binding:"required"`
}
type body struct {
	ID      int       `json:"uerid" binding:"required"`
	Name    string    `json:"user_name" binding:"required"`
	Phone   string    `json:"phone" binding:"required"`
	Habbits *[]string `json:"habbits" binding:"exists"`
	Weight  float32   `json:"weight"`
}
type query struct {
	Token string `form:"t" binding:"required"`
}

func main() {
	app := gin.New()

	r := middleware.NewRequestCheck()
	r.BindBody(&body{})
	r.BindHeader(&header{})
	r.BindParameter(&parameter{})
	r.BindQuery(&query{})

	//r.SetBadrequestHandler(b)

	app.POST("person/:person-id", r.CheckBind, echo)
	app.Run(":8080")
}

/*
func b(c *gin.Context, err error) {
	log.Println("oops~ ", err.Error())
	c.Abort()
}*/

func echo(c *gin.Context) {
	var b body
	var h header
	var q query
	var p parameter

	middleware.BodyUnmarshal(c, &b)
	middleware.ParameterUnmarshal(c, &p)
	middleware.QueryUnmarshal(c, &q)
	middleware.HeaderUnmarshal(c, &h)

	m := make(map[string]interface{})
	m["header"] = h
	m["parameter"] = p
	m["query"] = q
	m["body"] = b

	c.JSON(200, m)

}
