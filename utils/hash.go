package utils

import (
	"bufio"
	"crypto/sha1"
	"fmt"
	"io"
	"os"
)

func GetFileSHA1(filepath string) (string, error) {

	f, err := os.Open(filepath)
	if err != nil {
		return "", fmt.Errorf("con't open file")
	}
	defer f.Close()
	r := bufio.NewReader(f)

	h := sha1.New()

	_, err = io.Copy(h, r)
	if err != nil {
		return "", fmt.Errorf("read file fail")
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
