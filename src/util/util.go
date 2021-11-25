package util

import (
	"io/ioutil"
)

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

func ReadROM(path string) []byte {
	dat, err := ioutil.ReadFile(path)
	checkError(err)
	return dat
}
