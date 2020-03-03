package http

import (
	"io"
	"log"
	"net/http"
	"sync"

	"git.xtools.tv/tv/udf-tests/tc-proxy/tc"
)

type session struct {
	mu       sync.Mutex
	rule     tc.Rule
	id       string
	doneChan chan bool
}

type sessionParams struct {
	rule     tc.Rule
	doneChan chan bool
	id       string
}

func newSession(params sessionParams) *session {
	return &session{
		id:       params.id,
		rule:     params.rule,
		doneChan: params.doneChan,
	}
}

func (s *session) run(in io.Reader, out http.ResponseWriter) {
	fl, ok := out.(http.Flusher)
	if !ok {
		http.Error(out, "flushing not supported", http.StatusNotImplemented)
		return
	}

	for {
		s.mu.Lock()
		err := s.rule.Process(in, out)
		s.mu.Unlock()

		if err != nil {
			log.Printf("stop session id=%s: %v", s.id, err)
			break
		}
		fl.Flush()
	}
	s.doneChan <- true
}

func (s *session) setRule(rule tc.Rule) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.rule = rule
}
