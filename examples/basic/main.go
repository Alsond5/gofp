package main

import (
	"fmt"
	"io"
	"net/http"

	"github.com/Alsond5/gofp"
	"github.com/Alsond5/gofp/result"
)

func httpGet(url string) gofp.Result[*http.Response] {
	return gofp.Of(http.Get(url))
}

func readBody(resp *http.Response) gofp.Result[string] {
	defer resp.Body.Close()
	return result.Map(
		gofp.Of(io.ReadAll(resp.Body)),
		func(b []byte) string { return string(b) },
	)
}

func main() {
	result.FlatMap(
		httpGet("https://jsonplaceholder.typicode.com/todos/1"),
		readBody,
	).
		IfOk(func(body string) { fmt.Println("Result:", body) }).
		IfErr(func(err error) { fmt.Println("Error:", err) })
}
