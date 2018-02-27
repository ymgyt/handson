package awsctr

import (
	"io"
	"io/ioutil"
	"os"

	"github.com/juju/errors"
	yaml "gopkg.in/yaml.v2"
)

func ReadCfgFromYAMLFile(path string, cfg interface{}) error {
	f, err := os.Open(path)
	if err != nil {
		return errors.Trace(err)
	}
	return ReadCfgYAML(f, cfg)
}

func ReadCfgYAML(r io.Reader, cfg interface{}) error {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return errors.Trace(err)
	}
	return yaml.Unmarshal(data, cfg)
}
