package format

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"os"

	"gopkg.in/yaml.v2"
)

func WriteOrPrint(data interface{}, format string, pretty bool, output string) error {
	if output != "" {
		f, err := os.OpenFile(output, os.O_WRONLY|os.O_CREATE, 0644)
		if err != nil {
			return err
		}
		defer f.Close()
		return Write(data, format, Indent, f)
	} else {
		return Print(data, format, pretty)
	}
}

func Write(data interface{}, format string, indent string, writer io.Writer) error {
	switch format {
	case "yaml", "":
		return WriteYaml(data, writer)
	case "json":
		return WriteJson(data, writer, Indent)
	case "xml":
		return WriteXml(data, writer, Indent)
	default:
		return fmt.Errorf("unsupported format: %s", format)
	}
}

func WriteYaml(data interface{}, writer io.Writer) error {
	encoder := yaml.NewEncoder(writer)
	if slice, ok := data.([]interface{}); ok {
		// YAML separates each entry with "---"
		// (In JSON the slice would be written as an array)
		for _, d := range slice {
			err := encoder.Encode(d)
			if err != nil {
				return err
			}
		}
		return nil
	} else {
		return encoder.Encode(data)
	}
}

func WriteJson(data interface{}, writer io.Writer, indent string) error {
	encoder := json.NewEncoder(writer)
	encoder.SetIndent("", indent)
	return encoder.Encode(data)
}

func WriteXml(data interface{}, writer io.Writer, indent string) error {
	data, err := Normalize(data)
	if err != nil {
		return err
	}

	data = EnsureXml(data)

	_, err = io.WriteString(writer, xml.Header)
	if err != nil {
		return err
	}
	encoder := xml.NewEncoder(writer)
	encoder.Indent("", indent)
	return encoder.Encode(data)
}
