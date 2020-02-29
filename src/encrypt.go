package main

import (
	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/armor"
	"golang.org/x/crypto/openpgp/packet"
	"io"
	"os"
)

type GnuPG struct {
	WriteCloser io.WriteCloser
	file   *os.File
}

func (g *GnuPG) InitializeGnuPG(pubKey string, finalFilePath string) (err error) {
	var file *os.File
	file, err = os.Open(pubKey)
	if err != nil {
		return
	}
	defer func() { _ = file.Close() }()

	var block *armor.Block
	block, err = armor.Decode(file)
	if err != nil { return }

	var recipient *openpgp.Entity
	recipient, err = openpgp.ReadEntity(packet.NewReader(block.Body))
	if err != nil { return }

	if g.file, err = os.OpenFile(finalFilePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644); err != nil { return }

	g.WriteCloser, err = openpgp.Encrypt(g.file, []*openpgp.Entity{recipient},
	nil, &openpgp.FileHints{IsBinary: true}, nil)
	if err != nil { return err }

	return nil
}

func (g *GnuPG) Close()  {
	_ = g.WriteCloser.Close()
	_ = g.file.Close()
}

