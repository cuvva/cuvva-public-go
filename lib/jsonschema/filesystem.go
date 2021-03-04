package jsonschema

import (
	"fmt"
	"io/fs"
	"io/ioutil"

	"github.com/xeipuuv/gojsonschema"
)

type FS struct {
	raw fs.FS
}

func NewFS(raw fs.FS) *FS {
	return &FS{raw: raw}
}

func (f *FS) Load(filepath string) gojsonschema.JSONLoader {
	file, err := f.raw.Open(fmt.Sprintf("%s.json", filepath))
	if err != nil {
		panic(fmt.Errorf("open json schema file: %w", err))
	}

	b, err := ioutil.ReadAll(file)
	if err != nil {
		panic(fmt.Errorf("read json schema file: %s", err))
	}

	return gojsonschema.NewBytesLoader(b)
}
