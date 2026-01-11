package output

import (
	"io"

	"gopkg.in/yaml.v3"
)

func WriteYAML(w io.Writer, payload any) error {
	data, err := yaml.Marshal(payload)
	if err != nil {
		return err
	}
	_, err = w.Write(data)
	return err
}
