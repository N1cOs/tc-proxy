package tc

import (
	"bytes"
	"io"
	"regexp"
	"strings"
)

const (
	defaultBufSize = 512
)

type Rule interface {
	Process(io.Reader, io.Writer) error
}

type AcceptRule struct {
	buf []byte
}

type AcceptParams struct {
	BufSize int
}

func NewAcceptRule(params AcceptParams) *AcceptRule {
	bufSize := params.BufSize
	if bufSize == 0 {
		bufSize = defaultBufSize
	}

	buf := make([]byte, bufSize)
	return &AcceptRule{buf}
}

func (rule *AcceptRule) Process(in io.Reader, out io.Writer) error {
	read, err := in.Read(rule.buf)
	if read > 0 {
		_, err := out.Write(rule.buf[:read])
		if err != nil {
			return err
		}
	}
	return err
}

func (rule *AcceptRule) String() string {
	return "accept rule"
}

type DropRule struct {
	buf     []byte
	restMsg []byte
	regexp  *regexp.Regexp
}

type DropParams struct {
	MsgPattern string
	BufSize    int
}

func NewDropRule(params DropParams) (*DropRule, error) {
	regex, err := regexp.Compile(params.MsgPattern)
	if err != nil {
		return nil, err
	}

	bufSize := params.BufSize
	if bufSize == 0 {
		bufSize = defaultBufSize
	}
	buf := make([]byte, bufSize)

	rule := &DropRule{
		buf:     buf,
		restMsg: make([]byte, 0),
		regexp:  regex,
	}
	return rule, nil
}

func (rule *DropRule) Process(in io.Reader, out io.Writer) error {
	read, err := in.Read(rule.buf)
	if read > 0 {
		var data []byte
		data = append(data, rule.restMsg...)
		data = append(data, rule.buf[:read]...)

		msgs := bytes.NewBuffer(data)
		for {
			msg, _ := msgs.ReadString('\n')
			if !strings.HasSuffix(msg, "\n") {
				rule.restMsg = []byte(msg)
				break
			}

			if !rule.regexp.MatchString(msg) {
				_, err := out.Write([]byte(msg))
				if err != nil {
					return err
				}
			}
		}
	}
	return err
}

func (rule *DropRule) String() string {
	return "drop rule"
}
