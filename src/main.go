package main

import (
	"flag"

	"git.xtools.tv/tv/udf-tests/tc-proxy/grpc"
	"git.xtools.tv/tv/udf-tests/tc-proxy/http"
)

const (
	defaultHost     = "0.0.0.0"
	defaultHTTPPort = 80
	defaultGRPCPort = 530
)

func main() {
	var httpHost, grpcHost string
	var httpPort, grpcPort int

	flag.StringVar(&httpHost, "http-host", defaultHost, "--http-host HOST")
	flag.StringVar(&grpcHost, "grpc-host", defaultHost, "--grpc-host HOST")
	flag.IntVar(&httpPort, "http-port", defaultHTTPPort, "--http-port PORT")
	flag.IntVar(&grpcPort, "grpc-port", defaultGRPCPort, "--grpc-port PORT")
	flag.Parse()

	params := http.ProxyParams{
		Host: httpHost,
		Port: httpPort,
	}
	proxy, err := http.NewProxy(params)
	if err != nil {
		panic(err)
	}

	grpcParams := grpc.ServerParams{
		Host:  grpcHost,
		Port:  grpcPort,
		Proxy: proxy,
	}
	grpcServer := grpc.NewServer(grpcParams)
	go grpcServer.Serve()

	if err = proxy.Serve(); err != nil {
		panic(err)
	}
}
