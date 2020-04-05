package model

import (
	"github.com/skhatri/api-router-go/test"
	"testing"
)

func TestContainerProperties(t *testing.T) {
	var data = make(map[string]string)
	container := Response(data)
	test.EqualTo(t, 200, container.GetStatus())

	container.AddHeader("content-type", "application/json;charset=UTF-8")
	test.EqualTo(t, "application/json;charset=UTF-8", container.GetHeaders()["content-type"])
}

func TestErrorResponse(t *testing.T) {
	messageItem := MessageItem{
		Code:    "not-found",
		Message: "Not Found",
		Details: nil,
	}
	container := ErrorResponse(messageItem, 400)
	test.EqualTo(t, 400, container.GetStatus())
	test.EqualTo(t, 1, len(container.Errors))
	test.EqualTo(t, messageItem.Code, container.Errors[0].Code)
}

func TestListResponse(t *testing.T) {
	var items = make([]interface{}, 0)
	for i := 0; i < 5; i++ {
		items = append(items, string(i))
	}
	container := ListResponse(items)
	typeCheck := container.Data.([]interface{})
	test.EqualTo(t, 5, len(typeCheck))
}

func TestResponse(t *testing.T) {
	data := Response(MessageItem{Code: "Test"}).Data.(MessageItem)
	test.EqualTo(t, "Test", data.Code)
}
