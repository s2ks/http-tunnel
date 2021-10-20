package encoding

import "encoding/base64"

func Decode(data []byte) ([]byte, error) {
	max_len := base64.RawURLEncoding.DecodedLen(len(data))
	decoded_buf := make([]byte, max_len)

	n, err := base64.RawURLEncoding.Decode(decoded_buf, data)

	return decoded_buf[0:n], err
}
