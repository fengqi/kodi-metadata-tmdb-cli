package utils

import (
	"encoding/xml"
	"os"
)

func SaveNfo(file string, v interface{}) error {
	if file == "" {
		return nil
	}

	bytes, err := xml.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}

	f, err := os.OpenFile(file, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			Logger.WarningF("save nfo close file err: %v", err)
		}
	}(f)

	_, err = f.Write([]byte(xml.Header))
	if err != nil {
		return err
	}

	_, err = f.Write(bytes)
	if err != nil {
		return err
	}

	return nil
}
