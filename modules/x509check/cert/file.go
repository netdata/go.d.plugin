package cert

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
)

func NewFile(path string) *File {
	return &File{path}
}

type File struct {
	Path string
}

func (f File) Gather() ([]*x509.Certificate, error) {
	content, err := ioutil.ReadFile(f.Path)
	if err != nil {
		return nil, fmt.Errorf("error on reading '%s' : %v", f.Path, err)
	}

	block, _ := pem.Decode(content)
	if block == nil {
		return nil, fmt.Errorf("error on decoding '%s' : %v", f.Path, err)
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("error on parsing certigicate '%s' : %v", f.Path, err)
	}

	return []*x509.Certificate{cert}, nil
}
