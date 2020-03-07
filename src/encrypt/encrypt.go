package encrypt

import (
	"errors"
	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/armor"
	"golang.org/x/crypto/openpgp/packet"
	"io"
	"os"
)

type Encrypt struct {
	writeCloser io.WriteCloser
	file        *os.File
}

func (g *Encrypt) Initialize(file *os.File, pubKeyPath string) (err error) {
	g.file = file

	file, err = os.Open(pubKeyPath)
	if err != nil {
		return errors.New("Open pub key failed: " + err.Error())
	}
	defer file.Close()

	var block *armor.Block
	block, err = armor.Decode(file)
	if err != nil {
		return errors.New("Decode pub key failed: " + err.Error())
	}

	var recipient *openpgp.Entity
	recipient, err = openpgp.ReadEntity(packet.NewReader(block.Body))
	if err != nil {
		return errors.New("Read entity failed: " + err.Error())
	}

	g.writeCloser, err = openpgp.Encrypt(g.file, []*openpgp.Entity{recipient},
	nil, &openpgp.FileHints{IsBinary: true}, nil)
	if err != nil {
		return errors.New("Open encrypt writer failed: " + err.Error())
	}

	return nil
}

func (g *Encrypt) GetWriter() *io.WriteCloser {
	return &g.writeCloser
}

func (g *Encrypt) Close()  {
	g.writeCloser.Close()
}

