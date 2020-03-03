package grpc

import (
	"context"
	"fmt"
	"log"
	"net"

	it "git.xtools.tv/tv/udf-tests/tc-proxy/grpc/internal"
	"git.xtools.tv/tv/udf-tests/tc-proxy/tc"
	"google.golang.org/grpc"
)

type Server struct {
	addr  string
	proxy tc.Proxy
}

type ServerParams struct {
	Host  string
	Port  int
	Proxy tc.Proxy
}

func NewServer(params ServerParams) *Server {
	addr := fmt.Sprintf("%s:%d", params.Host, params.Port)
	return &Server{
		addr:  addr,
		proxy: params.Proxy,
	}
}

func (s *Server) Serve() error {
	listener, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}
	log.Printf("grpc server started on %s", s.addr)

	server := grpc.NewServer()
	it.RegisterProxyServer(server, s)
	return server.Serve(listener)
}

func (s *Server) SetAcceptRule(ctx context.Context, req *it.AcceptRuleRequest) (*it.Empty, error) {
	rule := tc.NewAcceptRule(tc.AcceptParams{})
	params, err := s.keyParams(req.Key)
	if err != nil {
		return nil, err
	}

	s.proxy.SetRule(*params, rule)
	return &it.Empty{}, nil
}

func (s *Server) SetDropRule(ctx context.Context, req *it.DropRuleRequest) (*it.Empty, error) {
	ruleParams := tc.DropParams{
		MsgPattern: req.Pattern,
	}
	rule, err := tc.NewDropRule(ruleParams)
	if err != nil {
		return nil, err
	}

	params, err := s.keyParams(req.Key)
	if err != nil {
		return nil, err
	}

	s.proxy.SetRule(*params, rule)
	return &it.Empty{}, nil
}

func (s *Server) keyParams(key *it.Key) (*tc.KeyParams, error) {
	srcHosts, err := net.LookupHost(key.Src)
	if err != nil {
		return nil, err
	}
	srcHost := srcHosts[0]

	destHosts, err := net.LookupHost(key.Dest)
	if err != nil {
		return nil, err
	}
	destHost := destHosts[0]

	params := &tc.KeyParams{
		Src:  srcHost,
		Dest: destHost,
	}
	return params, nil
}
