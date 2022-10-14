package jsonschema

import (
	"fmt"
	"io"
	"io/fs"

	"github.com/xeipuuv/gojsonschema"
)

type FS struct {
	raw fs.FS
}

// NewFS will typically be passed an embed.FS to easily load separate .json schema files
func NewFS(raw fs.FS) *FS {
	return &FS{raw: raw}
}

// Loads the given filename suffixing with ".json"
func (f *FS) LoadJSONExt(filepath string) gojsonschema.JSONLoader {
	return f.Load(fmt.Sprintf("%s.json", filepath))
}

// Loads the given filename
func (f *FS) Load(filepath string) gojsonschema.JSONLoader {
	file, err := f.raw.Open(filepath)
	if err != nil {
		panic(fmt.Errorf("open json schema file: %w", err))
	}

	b, err := io.ReadAll(file)
	if err != nil {
		panic(fmt.Errorf("read json schema file: %s", err))
	}

	return gojsonschema.NewBytesLoader(b)
}
