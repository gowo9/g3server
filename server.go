package g3server

import (
	"context"
	"net"
	"net/http"
	"strings"
	"time"

	gwRuntime "github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const readHeaderTimeout = 3 * time.Second

// Server gRPC 和 gRPC-Gateway 共用 port 的 server
type Server struct {
	pipe *PipeListener

	gRPCServer *grpc.Server

	gwProxyMux *gwRuntime.ServeMux

	httpServer *http.Server
}

type GatewayRegisterFunc func(ctx context.Context, mux *gwRuntime.ServeMux, conn *grpc.ClientConn) error

// New 建立一個新的 ggserver
func New(ctx context.Context, gs *grpc.Server, gwRegFunc GatewayRegisterFunc, optArgs ...Option) (s *Server, err error) {
	opt := defaultg3SOption
	for _, o := range optArgs {
		o.apply(&opt)
	}

	pipe := ListenPipe()

	// setting gRPC Dial

	if opt.grpcDialOptList == nil {
		opt.grpcDialOptList = make([]grpc.DialOption, 0, 2)
	}
	opt.grpcDialOptList = append(opt.grpcDialOptList,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithContextDialer(func(c context.Context, s string) (net.Conn, error) {
			return pipe.DialContext(c, `pipe`, s)
		}),
	)
	clientPipeConn, err := grpc.Dial(`pipe`, opt.grpcDialOptList...)
	if err != nil {
		return
	}

	// setting gRPC Gateway

	proxyMux := gwRuntime.NewServeMux(opt.gwServeMuxOptList...)
	if ctx == nil {
		ctx = context.Background()
	}
	err = gwRegFunc(ctx, proxyMux, clientPipeConn)
	if err != nil {
		return
	}

	if opt.httpServer == nil {
		opt.httpServer = &http.Server{
			ReadHeaderTimeout: readHeaderTimeout,
		}
	}

	s = &Server{
		pipe:       pipe,
		gRPCServer: gs,
		gwProxyMux: proxyMux,
		httpServer: opt.httpServer,
	}
	s.httpServer.Handler = s

	return s, nil
}

// ServeHTTP 實作 http.Handler
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
		s.gRPCServer.ServeHTTP(w, r) // application/grpc 直接給 gRPC 處理
	} else {
		s.gwProxyMux.ServeHTTP(w, r) // 非 gRPC 交給 gRPC-Gateway 處理
	}
}

func (s *Server) startGRPCServer() {
	go func() {
		_ = s.gRPCServer.Serve(s.pipe)
	}()
}

func (s *Server) Serve(l net.Listener) (err error) {
	s.startGRPCServer()

	// 配置 h2c
	// gPRC 必須使用 http2 連線，但是內建的 http 並不支援非 TLS 的 HTTP2(h2c)
	// 所以這裡得要配置 h2c，讓 gRPC 客戶端可以正常連線

	var http2Server http2.Server
	err = http2.ConfigureServer(s.httpServer, &http2Server)
	if err != nil {
		return
	}
	s.httpServer.Handler = h2c.NewHandler(s, &http2Server)

	return s.httpServer.Serve(l)
}

func (s *Server) ServeTLS(l net.Listener, certFile, keyFile string) (err error) {
	s.startGRPCServer()

	return s.httpServer.ServeTLS(l, certFile, keyFile)
}
