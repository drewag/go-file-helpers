package filehelpers

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"errors"
	"io"
	"math"
	"os"
)

// Constants

const sFileChunk = 256

// Public Functions

// Generate a hash for the given file
func HashForFileAtPath(path string) (string, error) {
	file, err := os.Open(path)

	if err != nil {
		return "", err
	}

	defer file.Close()

	// calculate the file size
	info, _ := file.Stat()

	filesize := info.Size()

	blocks := uint64(math.Ceil(float64(filesize) / float64(sFileChunk)))

	hash := md5.New()

	for i := uint64(0); i < blocks; i++ {
		blocksize := int(math.Min(sFileChunk, float64(filesize-int64(i*sFileChunk))))
		buf := make([]byte, blocksize)

		file.Read(buf)
		io.WriteString(hash, string(buf)) // append into the hash
	}

	return base64.URLEncoding.EncodeToString(hash.Sum(nil)), nil
}

// Copy the file at src to dst
func CopyFile(src, dst string) error {
	s, err := os.Open(src)
	if err != nil {
		return err
	}
	defer s.Close()
	d, err := os.Create(dst)
	if err != nil {
		return err
	}
	if _, err := io.Copy(d, s); err != nil {
		d.Close()
		return err
	}
	return d.Close()
}

// Compare files
func AreFilesIdentical(lhs, rhs string) (bool, error) {
	lFile, err := os.Open(lhs)
	if err != nil {
		return false, err
	}
	rFile, err := os.Open(rhs)
	if err != nil {
		return false, err
	}
	lFileInfo, err := lFile.Stat()
	if err != nil {
		return false, err
	}
	rFileInfo, err := rFile.Stat()
	if err != nil {
		return false, err
	}

	if lFileInfo.Size() != rFileInfo.Size() {
		return false, nil
	}
	if lFileInfo.IsDir() || rFileInfo.IsDir() {
		return false, errors.New("Cannot compare directories")
	}

	const readSize = 1024
	lBytes := make([]byte, readSize)
	rBytes := make([]byte, readSize)
	for err == nil {
		_, err = lFile.Read(lBytes)
		if err != nil {
			break
		}
		_, err = rFile.Read(rBytes)
		if err != nil {
			break
		}

		if bytes.Compare(lBytes, rBytes) != 0 {
			return false, nil
		}
	}

	if err != io.EOF {
		return false, err
	}

	return true, nil
}
