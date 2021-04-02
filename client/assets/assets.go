package assets

import (
	"io/ioutil"
)

func Asset(certPath string) ([]byte, error) {
	return ioutil.ReadFile(certPath)
}
