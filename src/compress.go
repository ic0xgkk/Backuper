package main

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"os"
	"path/filepath"
)

type Compress struct {
	file             *os.File
	gzipWriter       *gzip.Writer
	tarWriter        *tar.Writer
	Logger           *Logger
	GnuPG            *GnuPG
}


func (c *Compress) InitializeCompress(logger *Logger, gnupg *GnuPG) (err error) {
	// Binary stream: tar -> gz -> GnuPG -> file

	c.Logger = logger
	c.GnuPG = gnupg

	// Initialize gzip compress handle
	if c.gzipWriter, err = gzip.NewWriterLevel(c.GnuPG.WriteCloser, gzip.BestCompression); err != nil {
		return
	}

	// Initialize tar compress handle
	c.tarWriter = tar.NewWriter(c.gzipWriter)

	return
}

func (c *Compress) walkFunction(path string, fileInfo os.FileInfo, inputErr error) (err error) {
	// Check walk if be with error
	if inputErr != nil {
		c.Logger.Error("Walk file failed: ", inputErr.Error())
		return
	}

	var header *tar.Header

	// Regular file
	if fileInfo.Mode().IsRegular() {
		c.Logger.Info("Compress: ", path, " [Regular File]")
		if header, err = tar.FileInfoHeader(fileInfo, path); err != nil {
			c.Logger.Error(err.Error())
			return
		}

		header.Name = filepath.ToSlash(path)
		if err = c.tarWriter.WriteHeader(header); err != nil {
			c.Logger.Error(err.Error())
			return
		}

		var data *os.File
		data, err = os.Open(path)
		if err != nil {
			c.Logger.Error(err.Error())
			return
		}
		if _, err = io.Copy(c.tarWriter, data); err != nil {
			c.Logger.Error(err.Error())
			return
		}
		_ = data.Close()

		return nil
	}

	// Dir
	if fileInfo.Mode().IsDir() {
		c.Logger.Info("Compress: ", path, " [Directory]")
		if header, err = tar.FileInfoHeader(fileInfo, path); err != nil {
			c.Logger.Error(err.Error())
			return
		}
		header.Name = filepath.ToSlash(path)
		if err = c.tarWriter.WriteHeader(header); err != nil {
			c.Logger.Error(err.Error())
			return
		}

		return nil
	}

	// Symlink
	if fileInfo.Mode() & os.ModeSymlink != 0 {
		c.Logger.Info("Compress: ", path, " [Symlink]")
		if header, err = tar.FileInfoHeader(fileInfo, path); err != nil {
			c.Logger.Error(err.Error())
			return
		}
		header.Name = filepath.ToSlash(path)
		if err = c.tarWriter.WriteHeader(header); err != nil {
			c.Logger.Error(err.Error())
			return
		}

		return nil
	}

	c.Logger.Error("Unsupported file type: ", path)

	return nil
}

func (c *Compress) Compress(inputPath string) (err error) {
	err = filepath.Walk(inputPath, c.walkFunction)
	return
}

func (c *Compress) Close() {
	_ = c.tarWriter.Close()
	_ = c.gzipWriter.Close()
	_ = c.file.Close()
}


