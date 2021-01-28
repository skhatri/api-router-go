package model

import "strings"

type Container struct {
	Errors   []MessageItem          `json:"errors,omitempty"`
	Data     interface{}            `json:"data,omitempty"`
	Warnings []MessageItem          `json:"warnings,omitempty"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
	status   int
	headers  map[string]string
	decorate *bool
}

func (c *Container) GetStatus() int {
	if c.status == 0 {
		return 200
	}
	return c.status
}

func (c *Container) AddHeader(key string, value string) *Container {
	if c.headers == nil {
		c.headers = make(map[string]string)
	}
	c.headers[strings.ToLower(key)] = value
	return c
}
func (c *Container) GetHeaders() map[string]string {
	return c.headers
}
func (c *Container) AsJson() *Container {
	c.AddHeader("Content-Type", "application/json;charset=utf-8")
	return c
}


func (c *Container) IsDecorated() bool {
	return c.decorate == nil || *c.decorate
}


type MessageItem struct {
	Code    string                 `json:"code"`
	Message string                 `json:"message"`
	Details map[string]interface{} `json:"details,omitempty"`
}

func ErrorResponse(mi MessageItem, statusCode int) *Container {
	c := &Container{
		Errors: []MessageItem{
			mi,
		},
		status: statusCode,
	}
	return c.AsJson()
}

func WithErrorOnly(mi MessageItem, statusCode int) *Container {
	decorateFlag := false
	c := &Container{
		Errors: []MessageItem{
			mi,
		},
		status:   statusCode,
		decorate: &decorateFlag,
	}
	return c.AsJson()
}

func ListResponse(data []interface{}) *Container {
	c := &Container{
		Data: data,
	}
	return c.AsJson()
}

func Response(data interface{}) *Container {
	c := &Container{
		Data: data,
	}
	return c.AsJson()
}

func ResponseWithStatusCode(data interface{}, statusCode int) *Container {
	c := &Container{
		Data:   data,
		status: statusCode,
	}
	return c.AsJson()
}

func WithDataOnly(data interface{}) *Container {
	decorateFlag := false
	c := &Container{
		Data:     data,
		decorate: &decorateFlag,
	}
	return c.AsJson()
}
