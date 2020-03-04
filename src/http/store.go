package http

import (
	"fmt"
	"log"
	"strconv"
	"sync"

	"git.xtools.tv/tv/udf-tests/tc-proxy/tc"
)

type key struct {
	src  string
	dest string
}

type store struct {
	mu       sync.Mutex
	sessions map[key]map[string]*session
	rules    map[key]tc.Rule
	lastID   uint64
}

func newStore() *store {
	return &store{
		sessions: make(map[key]map[string]*session),
		rules:    make(map[key]tc.Rule),
	}
}

type requestParams struct {
	src      string
	dest     string
	doneChan chan bool
}

func (s *store) new(req requestParams) *session {
	s.mu.Lock()
	defer s.mu.Unlock()

	key := key{
		src:  req.src,
		dest: req.dest,
	}

	rule, ok := s.rules[key]
	if !ok {
		rule = tc.NewAcceptRule(tc.AcceptParams{})
	}

	id := s.nextID()
	params := sessionParams{
		id:       id,
		rule:     rule,
		doneChan: req.doneChan,
	}
	sess := newSession(params)

	_, ok = s.sessions[key]
	if !ok {
		s.sessions[key] = make(map[string]*session)
	}

	s.sessions[key][id] = sess
	log.Printf("added session %s -> %s: id=%s", key.src, key.dest, id)

	return sess
}

func (s *store) remove(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, sessions := range s.sessions {
		_, ok := sessions[id]
		if ok {
			delete(sessions, id)
			log.Printf("remove session: id=%s", id)
			return nil
		}
	}
	return fmt.Errorf("unknown session id %s", id)
}

func (s *store) nextID() string {
	s.lastID++
	return strconv.FormatUint(s.lastID, 10)
}

func (s *store) setRule(key key, rule tc.Rule) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.rules[key] = rule
	for _, session := range s.sessions[key] {
		session.setRule(rule)
	}
	log.Printf("set %s: src=%s, dest=%s", rule, key.src, key.dest)
}
