package config

import (
	"bytes"
	"io"
	"os"

	"github.com/fgrimme/ankiPDF/layout"
	yaml "gopkg.in/yaml.v2"
)

type Side struct {
	Fields []string `yaml:"fields"`
	Layout Layout   `yaml:"layout"`
}

type Config struct {
	CardSize     layout.DIN        `yaml:"cards"`
	FieldLayouts map[string]Layout `yaml:"field_layouts"`
	Front        Side              `yaml:"front"`
	Back         Side              `yaml:"back"`
	Empty        map[string]string `yaml:"empty"`
}

type Layout struct {
	Font   string  `yaml:"font"`
	Size   float64 `yaml:"size"`
	Height float64 `yaml:"height"`
}

// FromFile loads a configuration from file.
func FromFile(cfgpath string) (*Config, error) {
	f, err := os.Open(cfgpath)
	if err != nil {
		return nil, err
	}
	return load(f)
}

// load loads configuration from an io.Reader.
// Note, the configuration does not get validated or sanitized.
func load(in io.Reader) (*Config, error) {
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(in)
	if err != nil {
		return nil, err
	}

	c := &Config{}
	err = yaml.Unmarshal(buf.Bytes(), c)
	if err != nil {
		return nil, err
	}
	return c, nil
}
