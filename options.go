package g3server

import (
	"net/http"

	gwRuntime "github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

type g3SOption struct {
	grpcDialOptList   []grpc.DialOption
	gwServeMuxOptList []gwRuntime.ServeMuxOption
	httpServer        *http.Server
}

var defaultg3SOption = g3SOption{
	grpcDialOptList:   nil,
	gwServeMuxOptList: nil,
	httpServer:        nil,
}

type Option interface {
	apply(*g3SOption)
}

type funcG3SOption struct {
	f func(*g3SOption)
}

func (fdo *funcG3SOption) apply(do *g3SOption) {
	fdo.f(do)
}

func newG3SFuncOption(f func(*g3SOption)) *funcG3SOption {
	return &funcG3SOption{
		f: f,
	}
}

func WithGRPCDialOption(optList []grpc.DialOption) Option {
	return newG3SFuncOption(func(o *g3SOption) {
		o.grpcDialOptList = optList
	})
}

func WithGWServeMuxOption(optList []gwRuntime.ServeMuxOption) Option {
	return newG3SFuncOption(func(o *g3SOption) {
		o.gwServeMuxOptList = optList
	})
}

func WithHTTPServerOption(s *http.Server) Option {
	return newG3SFuncOption(func(o *g3SOption) {
		o.httpServer = s
	})
}
