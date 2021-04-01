package dasrpc

import (
	"fmt"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/rs/cors"
	"net"
	"net/http"
)

/**
 * Copyright (C), 2019-2020
 * FileName: rpc
 * Author:   LinGuanHong
 * Date:     2020/12/21 11:54 上午
 * Description:
 */

type JsonrpcOptions struct {
	Port string
}

type RpcServiceDelegate struct {
	Name    string
	Element interface{}
}

type JsonrpcServiceImpl struct {
	port       string
	httpServer *http.Server
	delegates  []*RpcServiceDelegate
}

func NewJsonrpcService(port string, delegate ...*RpcServiceDelegate) *JsonrpcServiceImpl {
	l := &JsonrpcServiceImpl{}
	l.port = port
	l.delegates = append(l.delegates, delegate...)
	return l
}

func (*JsonrpcServiceImpl) Ping(val string, val2 int) (res string, err error) {
	return
}

func (j *JsonrpcServiceImpl) registerHandler(delegate ...*RpcServiceDelegate) (*rpc.Server, error) {
	handler := rpc.NewServer()
	size := len(delegate)
	for i := 0; i < size; i++ {
		j.delegates = append(j.delegates, delegate[i])
	}
	size2 := len(j.delegates)
	for k := 0; k < size2; k++ {
		if err := handler.RegisterName(j.delegates[k].Name, j.delegates[k].Element); err != nil {
			return nil, fmt.Errorf("rpc RegisterName err: %s", err.Error())
		}
	}
	return handler, nil
}

func (j *JsonrpcServiceImpl) Start(beforeServeFunc BeforeServeFunc) error {
	if j.httpServer != nil {
		return nil
	}
	var (
		listener net.Listener
		err      error
	)
	handler, err := j.registerHandler()
	if err != nil {
		return fmt.Errorf("jsonrpc register handler err: %s", err.Error())
	}
	if listener, err = net.Listen("tcp", ":"+j.port); err != nil {
		panic(err.Error())
	}
	j.httpServer = &http.Server{Handler: newCorsHandler(handler, []string{"*"}, beforeServeFunc)}
	if err = j.httpServer.Serve(listener); err != nil {
		return fmt.Errorf("jsonrpc serve err: %s", err.Error())
	}
	return nil
}

func (j *JsonrpcServiceImpl) Stop() {
	if j.httpServer != nil {
		_ = j.httpServer.Close()
	}
}

func newCorsHandler(srv *rpc.Server, allowedOrigins []string, bf BeforeServeFunc) http.Handler {
	if len(allowedOrigins) == 0 {
		return srv
	}
	c := cors.New(cors.Options{
		AllowedOrigins:   allowedOrigins,
		AllowedMethods:   []string{"POST", "GET", "OPTIONS"},
		MaxAge:           600,
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
		Debug:            false,
	})
	return NewRPCHandler(srv, c, bf)
}
