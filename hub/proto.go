//go:generate colf -b .. go proto.colf
//go:generate mv Colfer.go proto_auto.go
package hub

import (
	"encoding"
	"io"
	"math"

	"github.com/pkg/errors"
)

type (
	// RequestType defines all possible requests
	// version is not neeeded because a new version of request
	// is a new request
	RequestType     uint8
	colferMarshaler interface {
		MarshalTo(buf []byte) int
		MarshalLen() (int, error)
	}
	colferUnmarshaler encoding.BinaryUnmarshaler
)

const (
	QUIT RequestType = iota
	PUT
	GET
)

func (o *Request) WriteTo(w io.Writer) (n int64, err error)   { return writeTo256(w, o) }
func (o *Response) WriteTo(w io.Writer) (n int64, err error)  { return writeTo256(w, o) }
func (o *Request) ReadFrom(r io.Reader) (n int64, err error)  { return readFrom256(r, o) }
func (o *Response) ReadFrom(r io.Reader) (n int64, err error) { return readFrom256(r, o) }

func writeTo256(w io.Writer, o colferMarshaler) (int64, error) {
	l, err := o.MarshalLen()
	if err != nil {
		return 0, err
	}
	if l > math.MaxUint8 {
		return 0, errors.New("length out of range")
	}
	data := make([]byte, l+1)
	data[0] = uint8(l)
	o.MarshalTo(data[1:])
	nn, err := w.Write(data)
	return int64(nn), errors.Wrap(err, "fail to write data frame")
}

func readFrom256(r io.Reader, o colferUnmarshaler) (n int64, err error) {
	var l [1]byte
	nn, err := r.Read(l[:])
	n += int64(nn)
	if err != nil {
		return n, errors.Wrap(err, "fail to read data size")
	}
	buf := make([]byte, int(l[0]))
	nn, err = io.ReadAtLeast(r, buf, int(l[0]))
	n += int64(nn)
	if err != nil {
		return n, errors.Wrap(err, "fail to read data bytes")
	}
	return n, o.UnmarshalBinary(buf)
}
