package download

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

// Verbose if set to true package will print more information
var Verbose bool

// FromURL -  download file from provided url returns string - path where file was downloaded
func FromURL(saveTo string, url string, verbose bool) (string, error) {
	output, err := os.Create(saveTo)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	if Verbose {
		log.Printf("Downloading %v => %v\n", url, saveTo)
	}

	zipResp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	defer zipResp.Body.Close()

	// output, err := os.Create(filename)
	n, err := io.Copy(output, zipResp.Body)
	if err != nil {
		fmt.Println("error while downloading", url, "-", err)
		return string(""), err
	}
	if Verbose {
		log.Printf("Downloaded %v bytes", n)
	}

	return saveTo, nil
}

// VerifyChecksum - TODO
func VerifyChecksum(fileToCheck string, checksum string) (bool, error) {

	f, err := os.Open(fileToCheck)
	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()
	hasher := sha256.New()

	if _, err := io.Copy(hasher, f); err != nil {
		log.Fatal(err)
	}

	downloadedFileChecksum := hex.EncodeToString(hasher.Sum(nil))

	if downloadedFileChecksum != checksum {
		if Verbose {
			log.Printf("Checksum mismatch! %v != %v\n", checksum, downloadedFileChecksum)
		}
		return false, errors.New("Checksum mismatched")
	}

	return true, nil
}
