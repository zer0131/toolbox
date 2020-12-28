package common

import (
	"bytes"
	"compress/flate"
	"compress/gzip"
	"github.com/pkg/errors"
	"io/ioutil"
)

func GzipMarshal(raw []byte) ([]byte, error) {
	var b bytes.Buffer
	gz, err := gzip.NewWriterLevel(&b, flate.DefaultCompression)
	if err != nil {
		return nil, errors.Wrap(err, "")
	}
	if _, err := gz.Write(raw); err != nil {
		return nil, errors.Wrap(err, "")
	}
	if err := gz.Close(); err != nil {
		return nil, errors.Wrap(err, "")
	}
	return b.Bytes(), nil
}

func GzipUnmarshal(raw []byte) ([]byte, error) {
	rdata := bytes.NewReader(raw)
	r, err := gzip.NewReader(rdata)
	if err != nil {
		return nil, errors.Wrap(err, "")
	}

	rr, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, errors.Wrap(err, "")
	}
	return rr, nil
}
