package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type require struct {
	header            interface{}
	params            interface{}
	query             interface{}
	body              interface{}
	badRequestHandler func(c *gin.Context, err error)
}

var context_key string

var defaultBadrequestHandler func(c *gin.Context, err error)

func init() {
	context_key = uuid.New().String()
	defaultBadrequestHandler = abortWithBadrequest
}

func (c *require) SetBadrequestHandler(f func(c *gin.Context, err error)) {
	if f == nil {
		c.badRequestHandler = defaultBadrequestHandler
		return
	}
	c.badRequestHandler = f
	return
}

func abortWithBadrequest(c *gin.Context, err error) {
	c.JSON(http.StatusBadRequest, gin.H{
		"code":    http.StatusBadRequest,
		"message": err.Error(),
	})
	c.Abort()
}

func NewRequestCheck() (r *require) {
	return &require{
		badRequestHandler: defaultBadrequestHandler,
	}
}

//declare header require check
//using tag 'header' to naming field
func (c *require) BindHeader(t interface{}) (err error) {
	if t != nil {
		v := reflect.Indirect(reflect.ValueOf(t))
		count := v.NumField()
		for i := 0; i < count; i++ {
			f := v.Field(i)
			switch f.Kind() {
			case reflect.String:
			case reflect.Int, reflect.Int32, reflect.Int64:
			case reflect.Float32, reflect.Float64:
			default:
				return errors.New("type error: " + f.Kind().String())
			}
		}
		c.header = t
	}
	return nil
}

//declare query require check
//using gin BindQuery
func (c *require) BindQuery(t interface{}) (err error) {
	if t != nil {
		c.query = t
	}
	return errors.New("binding fail.")
}

//declare parameter require check
//using tag 'parameter' to naming field
func (c *require) BindParameter(t interface{}) error {
	if t != nil {
		v := reflect.Indirect(reflect.ValueOf(t))
		count := v.NumField()
		for i := 0; i < count; i++ {
			f := v.Field(i)
			switch f.Kind() {
			case reflect.String:
			case reflect.Int, reflect.Int32, reflect.Int64:
			case reflect.Float32, reflect.Float64:
			default:
				return errors.New("type error")
			}
		}
		c.params = t
	}
	return nil
}

//declare body require check
//using gin Bind
func (c *require) BindBody(t interface{}) (err error) {
	if t != nil {
		c.body = t
		return
	}
	return errors.New("binding fail.")
}

//gin middlware
func (a *require) CheckBind(c *gin.Context) {
	out := ""
	if a.header != nil {
		_a := newStruct(a.header)
		err := checkAndBuildSturct("header", c.GetHeader, _a)
		if err != nil {
			logger.Println(err.Error())
			a.badRequestHandler(c, err)
			return
		}
		c.Set("header-"+context_key, _a)
		if out != "" {
			out += "\t"
		}
		out += fmt.Sprintf("Header= %+v", _a)
	}
	if a.params != nil {
		_a := newStruct(a.params)
		err := checkAndBuildSturct("parameter", c.Param, _a)
		if err != nil {
			logger.Println(err.Error())
			a.badRequestHandler(c, err)
			return
		}
		c.Set("params-"+context_key, _a)

		if out != "" {
			out += "\t"
		}
		out += fmt.Sprintf("Params= %+v", _a)
	}

	if a.query != nil {
		_a := newStruct(a.query)
		err := c.BindQuery(_a)
		if err != nil {
			logger.Println(err.Error())
			a.badRequestHandler(c, err)
			return
		}
		c.Set("query-"+context_key, _a)

		if out != "" {
			out += "\t"
		}
		out += fmt.Sprintf("Query= %+v", _a)
	}

	if a.body != nil {
		body := newStruct(a.body)
		err := c.Bind(body)
		if err != nil {
			logger.Println(err.Error())
			a.badRequestHandler(c, err)
			return
		}

		c.Set("body-"+context_key, body)
		if out != "" {
			out += "\t"
		}
		out += fmt.Sprintf("Body= %+v", body)
	}
	if out != "" {
		logger.Println(out)
	}
}

func BodyUnmarshal(c *gin.Context, t interface{}) error {
	s, _ := c.Get("body-" + context_key)
	return unmarshal(s, t)
}

func HeaderUnmarshal(c *gin.Context, t interface{}) error {
	s, _ := c.Get("header-" + context_key)
	return unmarshal(s, t)
}

func ParameterUnmarshal(c *gin.Context, t interface{}) error {
	s, _ := c.Get("params-" + context_key)
	return unmarshal(s, t)
}

func QueryUnmarshal(c *gin.Context, t interface{}) error {
	s, _ := c.Get("query-" + context_key)
	return unmarshal(s, t)
}

func unmarshal(s, d interface{}) error {
	_s := reflect.Indirect(reflect.ValueOf(s))
	if _s.Type() == reflect.Indirect(reflect.ValueOf(d)).Type() {
		reflect.Indirect(reflect.ValueOf(d)).Set(_s)
	} else {
		d = nil
		logger.Println("can not unmarshal to", reflect.TypeOf(d).Name())
		return errors.New("can not unmarshal")
	}
	return nil
}
