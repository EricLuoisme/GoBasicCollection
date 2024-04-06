package main

import (
	"fmt"
	"net/http"
)

type DB interface {
	Store(string) error
}

type storage struct {
}

func (s *storage) Store(value string) error {
	fmt.Println("storing into db", value)
	return nil
}

// myExecuteFunc 作为展示, 在使用第三方library时, 我们需要嵌入一些额外的内容
// 这里Decorator就是这个func, 我们返回的type是一个func, 但同时
func myExecuteFunc(db DB) ExecuteFn {
	return func(s string) {
		// access to DB?
		fmt.Println("extra edit from my function", s)
		db.Store(s)
	}
}

func makeHttpFunc(db DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		db.Store("some http stuff")
	}
}

func main() {
	// 模拟传入一个db
	s := &storage{}
	http.HandleFunc("/", makeHttpFunc(s))
	Execute(myExecuteFunc(s))
}

// ExecuteFn coming from a third party lib
type ExecuteFn func(string)

// Execute 接收一个function作为入参, 内部填入"Foo bar baz"来调用入参
// 也就是我们的decorator
func Execute(fn ExecuteFn) {
	fn("Foo bar baz")
}
