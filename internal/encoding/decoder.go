package encoding

import (
	"encoding/base64"
	"io"
)

type Decoder struct {
	r io.Reader
}

func (d *Decoder) Read(dest []byte) (int, error) {
	return d.r.Read(dest)
}

func NewDecoderFromReader(r io.Reader) (*Decoder) {
	decoder := Decoder{}
	decoder.r =  base64.NewDecoder(base64.RawURLEncoding, r)

	return &decoder
}

func Decode(data []byte) ([]byte, error) {
	max_len := base64.RawURLEncoding.DecodedLen(len(data))
	decoded_buf := make([]byte, max_len)

	n, err := base64.RawURLEncoding.Decode(decoded_buf, data)

	return decoded_buf[:n], err
}
