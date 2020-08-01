package configfile

import (
	"io/ioutil"

	"github.com/go-yaml/yaml"
	"github.com/hanjunlee/awscred/core"
)

// YamlHandler handle the yaml-format configuration file.
type YamlHandler struct {
	filepath string
}

// NewYamlHandler create a new yaml handler.
func NewYamlHandler()

// Read read the configuration file.
func (h *YamlHandler) Read() (map[string]core.Config, error) {
	data, err := ioutil.ReadFile(h.filepath)
	if err != nil {
		return nil, err
	}

	c := make(map[string]core.Config)
	err = yaml.Unmarshal(data, c)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func (h *YamlHandler) Write(c map[string]core.Config) error {
	out, err := yaml.Marshal(c)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(h.filepath, out, 0644)
}
