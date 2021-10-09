package encoding

import "encoding/base64"

func Encode(data []byte) string {
	return base64.RawURLEncoding.EncodeToString(data)
}
