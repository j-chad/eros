package utils

import (
	"encoding/json"
	"io"
)

const maxJSONSize = 1 << 20 // 1 MB

func SafeJSONDecode(data io.Reader, v any) error {
	safeReader := io.LimitReader(data, maxJSONSize)
	dec := json.NewDecoder(safeReader)
	return dec.Decode(v)
}
