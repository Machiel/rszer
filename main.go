package main

import (
	"bytes"
	"github.com/gographics/imagick/imagick"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

func resizeFile(file []byte, width uint, height uint) (image.Image, string) {

	mw := imagick.NewMagickWand()
	defer mw.Destroy()

	err := mw.ReadImageBlob(file)

	if err != nil {
		log.Fatal(err)
	}

	err = mw.ResizeImage(width, height, imagick.FILTER_LANCZOS, 1.0)

	if err != nil {
		log.Fatal(err)
	}

	err = mw.SetImageCompressionQuality(95)

	blob := mw.GetImageBlob()

	image, format, err := image.Decode(bytes.NewBuffer(blob))

	if err != nil {
		log.Fatal(err)
	}

	return image, format
}

func resizeHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Requesting resize")

	image, err := ioutil.ReadAll(r.Body)

	if err != nil {
		log.Fatal(err)
		return
	}

	width, err := strconv.Atoi(r.URL.Query().Get("width"))

	if err != nil {
		log.Fatal(err)
		return
	}

	height, err := strconv.Atoi(r.URL.Query().Get("height"))

	if err != nil {
		log.Fatal(err)
		return
	}

	log.Printf("Requested resize, width %d, height %d", width, height)

	result, format := resizeFile(image, uint(width), uint(height))

	if format == "jpeg" {
		jpeg.Encode(w, result, &jpeg.Options{Quality: 90})
	} else if format == "gif" {
		gif.Encode(w, result, nil)
	} else if format == "png" {
		png.Encode(w, result)
	}
}

func main() {

	imagick.Initialize()
	defer imagick.Terminate()

	log.Println("Starting server")

	http.HandleFunc("/resize", resizeHandler)
	http.ListenAndServe(":4500", nil)

	log.Println("Server stopped")

}
