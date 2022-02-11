package ent

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/vmihailenco/msgpack/v5"
)

func (c Cursor) String() string {
	buf := &bytes.Buffer{}
	wc := base64.NewEncoder(base64.RawStdEncoding, buf)
	defer wc.Close()
	_ = msgpack.NewEncoder(wc).Encode(c)

	return buf.String()
}

func ParseCursorFromString(s *string) (*Cursor, error) {
	if s == nil {
		return nil, nil
	}

	var c Cursor
	if err := msgpack.NewDecoder(
		base64.NewDecoder(
			base64.RawStdEncoding,
			strings.NewReader(*s),
		),
	).Decode(&c); err != nil {
		return nil, fmt.Errorf("cannot decode cursor: %w", err)
	}
	return &c, nil
}
