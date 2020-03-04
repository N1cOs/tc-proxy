package http

import (
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"

	"git.xtools.tv/tv/udf-tests/tc-proxy/tc"
)

var (
	errEmptyHost = errors.New("host isn't specified")
	errEmptyPort = errors.New("port isn't specified")
)

type Proxy struct {
	store    *store
	listener net.Listener
	client   *http.Client
}

type ProxyParams struct {
	Host string
	Port int
}

func NewProxy(params ProxyParams) (*Proxy, error) {
	if params.Host == "" {
		return nil, errEmptyHost
	}

	if params.Port == 0 {
		return nil, errEmptyPort
	}

	addr := fmt.Sprintf("%s:%d", params.Host, params.Port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}
	log.Printf("http proxy started on %s", addr)

	proxy := &Proxy{
		store:    newStore(),
		listener: listener,
		client:   http.DefaultClient,
	}
	return proxy, nil
}

func (p *Proxy) Serve() error {
	if err := http.Serve(p.listener, p); err != nil {
		return err
	}
	return nil
}

func (p *Proxy) ServeHTTP(dest http.ResponseWriter, target *http.Request) {
	destAddr, _, err := net.SplitHostPort(target.RemoteAddr)
	if err != nil {
		http.Error(dest, "can't parse client address", http.StatusBadRequest)
		return
	}

	srcs, err := net.LookupHost(target.URL.Hostname())
	if err != nil {
		fmt.Print(err)
		http.Error(dest, "can't resolve destination host", http.StatusBadRequest)
		return
	}
	srcAddr := srcs[0]

	src, err := p.openStream(target)
	if err != nil {
		http.Error(dest, err.Error(), http.StatusBadRequest)
		return
	}
	defer src.Body.Close()

	done := make(chan bool)
	reqParams := requestParams{
		src:      srcAddr,
		dest:     destAddr,
		doneChan: done,
	}

	session := p.store.new(reqParams)
	go session.run(src.Body, dest)
	<-done

	err = p.store.remove(session.id)
	if err != nil {
		http.Error(dest, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (p *Proxy) openStream(target *http.Request) (*http.Response, error) {
	req := &http.Request{
		Method: target.Method,
		URL:    target.URL,
		Header: target.Header,
		Body:   target.Body,
	}
	return p.client.Do(req)
}

func (p *Proxy) SetRule(params tc.KeyParams, rule tc.Rule) {
	key := key{
		src:  params.Src,
		dest: params.Dest,
	}
	p.store.setRule(key, rule)
}
