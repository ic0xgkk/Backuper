package file

import (
	"errors"
	"os"
)

type File struct {
	fileWriter *os.File
}

func (f *File) Initialize(filePath string, mode string) (err error) {
	switch mode {
	case "ro":
		f.fileWriter, err = os.OpenFile(filePath, os.O_RDONLY, 0644)
		break
	case "rw":
		f.fileWriter, err = os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
		break
	case "wa":
		f.fileWriter, err = os.OpenFile(filePath, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
		break
	default:
		panic("Error file mode")
	}
	if err != nil {
		return errors.New("Open file failed: " + err.Error())
	}

	return nil
}

func (f *File) GetWriter() *os.File {
	return f.fileWriter
}

func (f *File) Close() () {
	f.fileWriter.Close()
}
