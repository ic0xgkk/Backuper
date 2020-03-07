package compress

import (
	"archive/tar"
	"compress/gzip"
	"errors"
	"io"
	"os"
	"path/filepath"
)

type Compress struct {
	// Compress file (or encrypted file)
	file             *os.File

	// Inner stream
	gzipWriter       *gzip.Writer
	tarWriter        *tar.Writer

	// Global output writer
	writer           *io.WriteCloser

	fileInfoHook     func(path string, fileInfo os.FileInfo)
}

func (c *Compress) Initialize(writer *io.WriteCloser, fileInfoHook func(path string, fileInfo os.FileInfo)) (err error) {
	// Binary stream: tar -> gz -> writer

	c.writer = writer
	c.fileInfoHook = fileInfoHook

	// Initialize gzip compress handle
	c.gzipWriter, err = gzip.NewWriterLevel(*c.writer, gzip.BestCompression)
	if err != nil {
		return errors.New("Initialize gzip writer failed: " + err.Error())
	}

	// Initialize tar compress handle
	c.tarWriter = tar.NewWriter(c.gzipWriter)

	return nil
}

func (c *Compress) walkFunction(path string, fileInfo os.FileInfo, inputErr error) (err error) {
	// Check walk if be with error
	if inputErr != nil {
		return errors.New("Parent walk function reported an error: " + inputErr.Error())
	}

	c.fileInfoHook(path, fileInfo)

	var header *tar.Header

	if fileInfo.Mode().IsRegular() || fileInfo.Mode().IsDir() || (fileInfo.Mode() & os.ModeSymlink != 0) {
		if header, err = tar.FileInfoHeader(fileInfo, path); err != nil {
			return errors.New("Get header failed: " + err.Error())
		}

		header.Name = filepath.ToSlash(path)
		if err = c.tarWriter.WriteHeader(header); err != nil {
			return errors.New("Write header failed: " + err.Error())
		}
	}
	if fileInfo.Mode().IsRegular() {
		var data *os.File

		data, err = os.Open(path)
		if err != nil {
			return errors.New("Open file failed: " + err.Error())
		}
		defer data.Close()

		if _, err = io.Copy(c.tarWriter, data); err != nil {
			return errors.New("Copy data failed: " + err.Error())
		}
	}

	return nil
}

func (c *Compress) Compress(inputPath string) (err error) {
	err = filepath.Walk(inputPath, c.walkFunction)
	if err != nil {
		return errors.New("Compress failed: " + err.Error())
	}
	return nil
}

func (c *Compress) Close() {
	c.tarWriter.Close()
	c.gzipWriter.Close()
}

