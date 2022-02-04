package infra

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
)

type FileIO struct {
	PrettyJSON bool
	HomePath   string
	DataFolder string
}

func (fio FileIO) Read(fname string, dst interface{}) error {
	fpath := path.Join(fio.HomePath, fio.DataFolder, fname)

	data, err := ioutil.ReadFile(fpath)
	if err != nil {
		switch {
		case os.IsNotExist(err):
			return nil
		default:
			return err
		}
	}

	if err := json.Unmarshal(data, dst); err != nil {
		return err
	}

	return nil
}

func (fio FileIO) Write(fname string, src interface{}) error {
	var data []byte
	var err error
	if fio.PrettyJSON {
		data, err = json.MarshalIndent(src, "", "  ")
	} else {
		data, err = json.Marshal(src)
	}
	if err != nil {
		return nil
	}

	fpath := path.Join(fio.HomePath, fio.DataFolder, fname)

	if _, err := os.Stat(fpath); os.IsNotExist(err) {
		err := os.MkdirAll(path.Dir(fpath), 0o700)
		if err != nil {
			return err
		}
	}

	if err := ioutil.WriteFile(fpath, data, 0o700); err != nil {
		return err
	}

	return nil
}

func (fio FileIO) Delete(fname string) error {
	fpath := path.Join(fio.HomePath, fio.DataFolder, fname)
	if err := os.Remove(fpath); !os.IsNotExist(err) {
		return err
	}

	return nil
}
