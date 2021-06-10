package dasrpc

import (
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/rs/cors"
	"net/http"
)

/**
 * Copyright (C), 2019-2019
 * FileName: context
 * Author:   LinGuanHong
 * Date:     2019-11-28 15:12
 * Description:
 */

type BeforeServeFunc func(w http.ResponseWriter, r *http.Request)

type RPCHandler struct {
	srv         *rpc.Server
	cors        *cors.Cors
	BeforeServe BeforeServeFunc
}

func NewRPCHandler(srv *rpc.Server, cors *cors.Cors, cb BeforeServeFunc) RPCHandler {
	return RPCHandler{srv: srv, cors: cors, BeforeServe: cb}
}

func (rpc RPCHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if rpc.BeforeServe != nil {
		rpc.BeforeServe(w, r)
	}
	if rpc.cors != nil {
		rpc.cors.ServeHTTP(w, r, rpc.srv.ServeHTTP)
	}
}
