package funcOps

import (
	"fmt"
	"testing"
)

type Opts struct {
	maxConn int
	id      string
	tls     bool
}

type OptsFunc func(*Opts)

// defaultOps 更好的作为Library提供一个默认
func defaultOps() Opts {
	return Opts{
		maxConn: 10,
		id:      "default",
		tls:     false,
	}
}

// withTLS 单独处理某个值
func withTLS() OptsFunc {
	return func(opts *Opts) {
		opts.tls = true
	}
}

// withMaxConn 传递返回方法作为结果
func withMaxConn(n int) OptsFunc {
	return func(o *Opts) {
		o.maxConn = n
	}
}

func withID(id string) OptsFunc {
	return func(o *Opts) {
		o.id = id
	}
}

type Server struct {
	Opts
}

// newServer 的任意参数输入, 方便调用方选择合适的配置
func newServer(opts ...OptsFunc) *Server {
	o := defaultOps()
	for _, fn := range opts {
		fn(&o)
	}
	return &Server{
		Opts: o,
	}

}

func Test_FuncStruct(t *testing.T) {
	server := newServer(withID("234"), withMaxConn(101), withTLS())
	defaultServer := newServer()
	fmt.Printf("Server: %+v\n", server)
	fmt.Printf("Default: %+v\n", defaultServer)
}
