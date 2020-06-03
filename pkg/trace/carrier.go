package trace

import (
	"strings"
)

// MetadataReaderWriter ...
type MetadataReaderWriter struct {
	MD map[string][]string
}

// Set ...
func (w MetadataReaderWriter) Set(key, val string) {
	key = strings.ToLower(key)
	w.MD[key] = append(w.MD[key], val)
}

// ForeachKey ...
func (w MetadataReaderWriter) ForeachKey(handler func(key, val string) error) error {
	for k, vals := range w.MD {
		for _, v := range vals {
			if err := handler(k, v); err != nil {
				return err
			}
		}
	}

	return nil
}
