package encoding

import (
	"encoding/base64"
	"io"
)

type Decoder struct {
	r io.Reader
}

func (d *Decoder) Read(dest []byte) (n int, err error) {
	buflen := base64.RawURLEncoding.EncodedLen(len(dest))
	buf := make([]byte, buflen)

	n, err = d.r.Read(buf)

	if err != nil {
		return
	}

	n, err = base64.RawURLEncoding.Decode(dest, buf[:n])

	return
}

func NewDecoderFromReader(r io.Reader) (*Decoder) {
	decoder := Decoder{}
	decoder.r = r

	return &decoder
}

func Decode(data []byte) ([]byte, error) {
	max_len := base64.RawURLEncoding.DecodedLen(len(data))
	decoded_buf := make([]byte, max_len)

	n, err := base64.RawURLEncoding.Decode(decoded_buf, data)

	return decoded_buf[:n], err
}
