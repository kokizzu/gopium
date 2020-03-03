package fmts

import (
	"1pkg/gopium"
	"bytes"
	"encoding/json"
)

// StructFormat defines abstraction for
// formatting gopium.Struct to byte slice
type StructFormat func(gopium.Struct) ([]byte, error)

// PrettyJson defines json.Marshal
// with json.Indent TypeFormat implementation
func PrettyJson(st gopium.Struct) ([]byte, error) {
	r, err := json.Marshal(st)
	if err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	err = json.Indent(&buf, r, "", "\t")
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}