package functions

import (
	"github.com/skhatri/api-router-go/router"
	"github.com/skhatri/api-router-go/test"
	"testing"
)

func makeWebRequest() *router.WebRequest {
	return &router.WebRequest{
		Headers:     make(map[string]string),
		QueryParams: make(map[string]string),
		Body:        nil,
		Uri:         "/test",
		QueryString: "",
	}
}

func TestStatusFunc(t *testing.T) {
	result := StatusFunc(makeWebRequest())
	test.EqualTo(t, 200, result.GetStatus())
}

func TestEchoFunc(t *testing.T) {
	webReq := makeWebRequest()
	webReq.Headers["x-client"] = "client1"
	test.EqualTo(t, "client1", EchoFunc(webReq).GetHeaders()["x-client"])
}

func TestIndexFunc(t *testing.T) {
	webReq := makeWebRequest()
	webReq.Headers["x-client"] = "client1"
	test.EqualTo(t, 200, IndexFunc(webReq).GetStatus())
}
