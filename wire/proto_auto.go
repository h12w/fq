package wire

// Code generated by colf(1); DO NOT EDIT.
// The compiler used schema file proto.colf.

import (
	"encoding/binary"
	"fmt"
	"io"
)

var intconv = binary.BigEndian

// Colfer configuration attributes
var (
	// ColferSizeMax is the upper limit for serial byte sizes.
	ColferSizeMax = 16 * 1024 * 1024
)

// ColferMax signals an upper limit breach.
type ColferMax string

// Error honors the error interface.
func (m ColferMax) Error() string { return string(m) }

// ColferError signals a data mismatch as as a byte index.
type ColferError int

// Error honors the error interface.
func (i ColferError) Error() string {
	return fmt.Sprintf("colfer: unknown header at byte %d", i)
}

// ColferTail signals data continuation as a byte index.
type ColferTail int

// Error honors the error interface.
func (i ColferTail) Error() string {
	return fmt.Sprintf("colfer: data continuation at byte %d", i)
}

type Request struct {
	ID uint64

	Type uint8

	ClientID string

	JournalDir string

	Offset uint64
}

// MarshalTo encodes o as Colfer into buf and returns the number of bytes written.
// If the buffer is too small, MarshalTo will panic.
func (o *Request) MarshalTo(buf []byte) int {
	var i int

	if x := o.ID; x >= 1<<49 {
		buf[i] = 0 | 0x80
		intconv.PutUint64(buf[i+1:], x)
		i += 9
	} else if x != 0 {
		buf[i] = 0
		i++
		for x >= 0x80 {
			buf[i] = byte(x | 0x80)
			x >>= 7
			i++
		}
		buf[i] = byte(x)
		i++
	}

	if x := o.Type; x != 0 {
		buf[i] = 1
		i++
		buf[i] = x
		i++
	}

	if l := len(o.ClientID); l != 0 {
		buf[i] = 2
		i++
		x := uint(l)
		for x >= 0x80 {
			buf[i] = byte(x | 0x80)
			x >>= 7
			i++
		}
		buf[i] = byte(x)
		i++
		i += copy(buf[i:], o.ClientID)
	}

	if l := len(o.JournalDir); l != 0 {
		buf[i] = 3
		i++
		x := uint(l)
		for x >= 0x80 {
			buf[i] = byte(x | 0x80)
			x >>= 7
			i++
		}
		buf[i] = byte(x)
		i++
		i += copy(buf[i:], o.JournalDir)
	}

	if x := o.Offset; x >= 1<<49 {
		buf[i] = 4 | 0x80
		intconv.PutUint64(buf[i+1:], x)
		i += 9
	} else if x != 0 {
		buf[i] = 4
		i++
		for x >= 0x80 {
			buf[i] = byte(x | 0x80)
			x >>= 7
			i++
		}
		buf[i] = byte(x)
		i++
	}

	buf[i] = 0x7f
	i++
	return i
}

// MarshalLen returns the Colfer serial byte size.
// The error return option is wire.ColferMax.
func (o *Request) MarshalLen() (int, error) {
	l := 1

	if x := o.ID; x >= 1<<49 {
		l += 9
	} else if x != 0 {
		for l += 2; x >= 0x80; l++ {
			x >>= 7
		}
	}

	if x := o.Type; x != 0 {
		l += 2
	}

	if x := len(o.ClientID); x != 0 {
		if x > ColferSizeMax {
			return 0, ColferMax(fmt.Sprintf("colfer: field wire.Request.ClientID exceeds %d bytes", ColferSizeMax))
		}
		for l += x + 2; x >= 0x80; l++ {
			x >>= 7
		}
	}

	if x := len(o.JournalDir); x != 0 {
		if x > ColferSizeMax {
			return 0, ColferMax(fmt.Sprintf("colfer: field wire.Request.JournalDir exceeds %d bytes", ColferSizeMax))
		}
		for l += x + 2; x >= 0x80; l++ {
			x >>= 7
		}
	}

	if x := o.Offset; x >= 1<<49 {
		l += 9
	} else if x != 0 {
		for l += 2; x >= 0x80; l++ {
			x >>= 7
		}
	}

	if l > ColferSizeMax {
		return l, ColferMax(fmt.Sprintf("colfer: struct wire.Request exceeds %d bytes", ColferSizeMax))
	}
	return l, nil
}

// MarshalBinary encodes o as Colfer conform encoding.BinaryMarshaler.
// The error return option is wire.ColferMax.
func (o *Request) MarshalBinary() (data []byte, err error) {
	l, err := o.MarshalLen()
	if err != nil {
		return nil, err
	}
	data = make([]byte, l)
	o.MarshalTo(data)
	return data, nil
}

// Unmarshal decodes data as Colfer and returns the number of bytes read.
// The error return options are io.EOF, wire.ColferError and wire.ColferMax.
func (o *Request) Unmarshal(data []byte) (int, error) {
	if len(data) == 0 {
		return 0, io.EOF
	}
	header := data[0]
	i := 1

	if header == 0 {
		start := i
		i++
		if i >= len(data) {
			goto eof
		}
		x := uint64(data[start])

		if x >= 0x80 {
			x &= 0x7f
			for shift := uint(7); ; shift += 7 {
				b := uint64(data[i])
				i++
				if i >= len(data) {
					goto eof
				}

				if b < 0x80 || shift == 56 {
					x |= b << shift
					break
				}
				x |= (b & 0x7f) << shift
			}
		}
		o.ID = x

		header = data[i]
		i++
	} else if header == 0|0x80 {
		start := i
		i += 8
		if i >= len(data) {
			goto eof
		}
		o.ID = intconv.Uint64(data[start:])
		header = data[i]
		i++
	}

	if header == 1 {
		start := i
		i++
		if i >= len(data) {
			goto eof
		}
		o.Type = data[start]
		header = data[i]
		i++
	}

	if header == 2 {
		if i >= len(data) {
			goto eof
		}
		x := uint(data[i])
		i++

		if x >= 0x80 {
			x &= 0x7f
			for shift := uint(7); ; shift += 7 {
				if i >= len(data) {
					goto eof
				}
				b := uint(data[i])
				i++

				if b < 0x80 {
					x |= b << shift
					break
				}
				x |= (b & 0x7f) << shift
			}
		}

		if x > uint(ColferSizeMax) {
			return 0, ColferMax(fmt.Sprintf("colfer: wire.Request.ClientID size %d exceeds %d bytes", x, ColferSizeMax))
		}

		start := i
		i += int(x)
		if i >= len(data) {
			goto eof
		}
		o.ClientID = string(data[start:i])

		header = data[i]
		i++
	}

	if header == 3 {
		if i >= len(data) {
			goto eof
		}
		x := uint(data[i])
		i++

		if x >= 0x80 {
			x &= 0x7f
			for shift := uint(7); ; shift += 7 {
				if i >= len(data) {
					goto eof
				}
				b := uint(data[i])
				i++

				if b < 0x80 {
					x |= b << shift
					break
				}
				x |= (b & 0x7f) << shift
			}
		}

		if x > uint(ColferSizeMax) {
			return 0, ColferMax(fmt.Sprintf("colfer: wire.Request.JournalDir size %d exceeds %d bytes", x, ColferSizeMax))
		}

		start := i
		i += int(x)
		if i >= len(data) {
			goto eof
		}
		o.JournalDir = string(data[start:i])

		header = data[i]
		i++
	}

	if header == 4 {
		start := i
		i++
		if i >= len(data) {
			goto eof
		}
		x := uint64(data[start])

		if x >= 0x80 {
			x &= 0x7f
			for shift := uint(7); ; shift += 7 {
				b := uint64(data[i])
				i++
				if i >= len(data) {
					goto eof
				}

				if b < 0x80 || shift == 56 {
					x |= b << shift
					break
				}
				x |= (b & 0x7f) << shift
			}
		}
		o.Offset = x

		header = data[i]
		i++
	} else if header == 4|0x80 {
		start := i
		i += 8
		if i >= len(data) {
			goto eof
		}
		o.Offset = intconv.Uint64(data[start:])
		header = data[i]
		i++
	}

	if header != 0x7f {
		return 0, ColferError(i - 1)
	}
	if i < ColferSizeMax {
		return i, nil
	}
eof:
	if i >= ColferSizeMax {
		return 0, ColferMax(fmt.Sprintf("colfer: struct wire.Request size exceeds %d bytes", ColferSizeMax))
	}
	return 0, io.EOF
}

// UnmarshalBinary decodes data as Colfer conform encoding.BinaryUnmarshaler.
// The error return options are io.EOF, wire.ColferError, wire.ColferTail and wire.ColferMax.
func (o *Request) UnmarshalBinary(data []byte) error {
	i, err := o.Unmarshal(data)
	if i < len(data) && err == nil {
		return ColferTail(i)
	}
	return err
}

type Response struct {
	RequestID uint64

	Err string
}

// MarshalTo encodes o as Colfer into buf and returns the number of bytes written.
// If the buffer is too small, MarshalTo will panic.
func (o *Response) MarshalTo(buf []byte) int {
	var i int

	if x := o.RequestID; x >= 1<<49 {
		buf[i] = 0 | 0x80
		intconv.PutUint64(buf[i+1:], x)
		i += 9
	} else if x != 0 {
		buf[i] = 0
		i++
		for x >= 0x80 {
			buf[i] = byte(x | 0x80)
			x >>= 7
			i++
		}
		buf[i] = byte(x)
		i++
	}

	if l := len(o.Err); l != 0 {
		buf[i] = 1
		i++
		x := uint(l)
		for x >= 0x80 {
			buf[i] = byte(x | 0x80)
			x >>= 7
			i++
		}
		buf[i] = byte(x)
		i++
		i += copy(buf[i:], o.Err)
	}

	buf[i] = 0x7f
	i++
	return i
}

// MarshalLen returns the Colfer serial byte size.
// The error return option is wire.ColferMax.
func (o *Response) MarshalLen() (int, error) {
	l := 1

	if x := o.RequestID; x >= 1<<49 {
		l += 9
	} else if x != 0 {
		for l += 2; x >= 0x80; l++ {
			x >>= 7
		}
	}

	if x := len(o.Err); x != 0 {
		if x > ColferSizeMax {
			return 0, ColferMax(fmt.Sprintf("colfer: field wire.Response.Err exceeds %d bytes", ColferSizeMax))
		}
		for l += x + 2; x >= 0x80; l++ {
			x >>= 7
		}
	}

	if l > ColferSizeMax {
		return l, ColferMax(fmt.Sprintf("colfer: struct wire.Response exceeds %d bytes", ColferSizeMax))
	}
	return l, nil
}

// MarshalBinary encodes o as Colfer conform encoding.BinaryMarshaler.
// The error return option is wire.ColferMax.
func (o *Response) MarshalBinary() (data []byte, err error) {
	l, err := o.MarshalLen()
	if err != nil {
		return nil, err
	}
	data = make([]byte, l)
	o.MarshalTo(data)
	return data, nil
}

// Unmarshal decodes data as Colfer and returns the number of bytes read.
// The error return options are io.EOF, wire.ColferError and wire.ColferMax.
func (o *Response) Unmarshal(data []byte) (int, error) {
	if len(data) == 0 {
		return 0, io.EOF
	}
	header := data[0]
	i := 1

	if header == 0 {
		start := i
		i++
		if i >= len(data) {
			goto eof
		}
		x := uint64(data[start])

		if x >= 0x80 {
			x &= 0x7f
			for shift := uint(7); ; shift += 7 {
				b := uint64(data[i])
				i++
				if i >= len(data) {
					goto eof
				}

				if b < 0x80 || shift == 56 {
					x |= b << shift
					break
				}
				x |= (b & 0x7f) << shift
			}
		}
		o.RequestID = x

		header = data[i]
		i++
	} else if header == 0|0x80 {
		start := i
		i += 8
		if i >= len(data) {
			goto eof
		}
		o.RequestID = intconv.Uint64(data[start:])
		header = data[i]
		i++
	}

	if header == 1 {
		if i >= len(data) {
			goto eof
		}
		x := uint(data[i])
		i++

		if x >= 0x80 {
			x &= 0x7f
			for shift := uint(7); ; shift += 7 {
				if i >= len(data) {
					goto eof
				}
				b := uint(data[i])
				i++

				if b < 0x80 {
					x |= b << shift
					break
				}
				x |= (b & 0x7f) << shift
			}
		}

		if x > uint(ColferSizeMax) {
			return 0, ColferMax(fmt.Sprintf("colfer: wire.Response.Err size %d exceeds %d bytes", x, ColferSizeMax))
		}

		start := i
		i += int(x)
		if i >= len(data) {
			goto eof
		}
		o.Err = string(data[start:i])

		header = data[i]
		i++
	}

	if header != 0x7f {
		return 0, ColferError(i - 1)
	}
	if i < ColferSizeMax {
		return i, nil
	}
eof:
	if i >= ColferSizeMax {
		return 0, ColferMax(fmt.Sprintf("colfer: struct wire.Response size exceeds %d bytes", ColferSizeMax))
	}
	return 0, io.EOF
}

// UnmarshalBinary decodes data as Colfer conform encoding.BinaryUnmarshaler.
// The error return options are io.EOF, wire.ColferError, wire.ColferTail and wire.ColferMax.
func (o *Response) UnmarshalBinary(data []byte) error {
	i, err := o.Unmarshal(data)
	if i < len(data) && err == nil {
		return ColferTail(i)
	}
	return err
}
