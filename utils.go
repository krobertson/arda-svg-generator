package main

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	go_qr "github.com/piglig/go-qr"
)

var qrcfg = go_qr.NewQrCodeImgConfig(10, 4, go_qr.WithSVGXMLHeader(true), go_qr.WithOptimalSVG())

// generateQrCode is used to generate an SVG version of the specified URL.
func generateQrCode(s string) (string, error) {
	errCorLvl := go_qr.Low
	qr, err := go_qr.EncodeText(s, errCorLvl)
	if err != nil {
		return "", err
	}
	return qr.SVGString(qrcfg, "#FFFFFF", "#000000")
}

// cacheImage is used to fetch the specified image and store it locally in the
// "cache" folder in data URI format. This is used so we can cache item iamges
// locally without redownloading them on each run. It returns the string of the
// local filename, or error on a failure. Note, it is possible to return a
// non-blank string on an error while writing. It will return "", nil when it
// gets a non-2xx status code. That isn't cached and should be treated as no
// image.
func cacheImage(url string) (string, error) {
	h := md5.Sum([]byte(url))
	urlfile := filepath.Join("cache", hex.EncodeToString(h[:])+".cache")

	if _, err := os.Stat(urlfile); err == nil {
		return urlfile, nil
	}

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		fmt.Printf("%q returned %d, not storing\n", url, resp.StatusCode)
		return "", nil
	}

	f, err := os.Create(urlfile)
	if err != nil {
		return "", err
	}
	defer f.Close()

	f.WriteString(fmt.Sprintf("data:%s;base64,", resp.Header.Get("Content-Type")))

	b := base64.NewEncoder(base64.StdEncoding, f)
	defer b.Close()

	_, err = io.Copy(b, resp.Body)
	return urlfile, err
}
