package process

import (
	"bytes"
	"io"
	"strings"
)

type output struct {
	stdout io.Writer
	buffer strings.Builder
}

func (o *output) Write(p []byte) (n int, err error) {
	if n, err = o.stdout.Write(p); err != nil {
		return n, err
	}

	p = bytes.ReplaceAll(p, []byte{'\x00'}, []byte{})
	p = bytes.ReplaceAll(p, []byte{'\r'}, []byte{})

	return o.buffer.Write(p)
}

func (o *output) Reset() {
	o.buffer.Reset()
}

func (o *output) String() string {
	return o.buffer.String()
}
