package download

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
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
		return "", err
	}

	if Verbose {
		log.Printf("Downloading %v => %v\n", url, saveTo)
	}

	zipResp, err := http.Get(url)

	if err != nil {
		return "", err
	}

	if zipResp.StatusCode != 200 {
		return "", errors.New(zipResp.Status) //&httpError{"wtf", zipResp.StatusCode}
	}

	defer zipResp.Body.Close()

	// output, err := os.Create(filename)
	n, err := io.Copy(output, zipResp.Body)
	if err != nil {
		log.Println("error while downloading", url, "-", err)
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
		log.Printf("error: copy failed %v", err)
		os.Exit(3)
	}

	downloadedFileChecksum := hex.EncodeToString(hasher.Sum(nil))

	if downloadedFileChecksum != checksum {
		if Verbose {
			log.Printf("checksum mismatch! %v != %v\n", checksum, downloadedFileChecksum)
		}
		return false, errors.New("checksum mismatched")
	}

	return true, nil
}
