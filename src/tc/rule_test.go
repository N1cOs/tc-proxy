package tc

import (
	"bytes"
	"io"
	"testing"
)

type testCase struct {
	Input    string
	Expected string
}

func TestAcceptRule(t *testing.T) {
	params := AcceptParams{
		BufSize: 5,
	}
	rule := NewAcceptRule(params)

	for _, tCase := range []testCase{
		{
			"",
			"",
		},
		{
			"h,42,126\r\n",
			"h,42,126\r\n",
		},
		{
			"h,42,126\r\nc,42,127\r\nh,42,128\r\n",
			"h,42,126\r\nc,42,127\r\nh,42,128\r\n",
		},
	} {
		reader := newReader(tCase.Input)
		writer := newWriter()

		var err error
		for err != io.EOF {
			err = rule.Process(reader, writer)
		}

		if res := writer.String(); res != tCase.Expected {
			t.Errorf("actual %q, expected %q", res, tCase.Expected)
		}
	}
}

func TestDropRule(t *testing.T) {
	params := DropParams{
		MsgPattern: "^c,42,",
		BufSize:    5,
	}
	rule, err := NewDropRule(params)
	if err != nil {
		t.Error(err)
	}

	for _, tCase := range []testCase{
		{
			"c,42,123\r\n",
			"",
		},
		{
			"h,42,126\r\n",
			"h,42,126\r\n",
		},
		{
			"c,43,125\r\nc,42,124\r\n",
			"c,43,125\r\n",
		},
		{
			"h,42,126\r\nc,42,127\r\nh,42,128\r\n",
			"h,42,126\r\nh,42,128\r\n",
		},
	} {
		reader := newReader(tCase.Input)
		writer := newWriter()

		var err error
		for err != io.EOF {
			err = rule.Process(reader, writer)
		}

		if res := writer.String(); res != tCase.Expected {
			t.Errorf("actual %q, expected %q", res, tCase.Expected)
		}
	}
}

func newReader(msgs string) io.Reader {
	return bytes.NewBufferString(msgs)
}

func newWriter() *bytes.Buffer {
	return &bytes.Buffer{}
}
