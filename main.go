package main

import (
	"bytes"
	"github.com/nfnt/resize"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

func resizeFile(file *bytes.Buffer, width uint, height uint) (image.Image, string) {

	img, format, err := image.Decode(file)

	if err != nil {
		log.Fatal(err)
	}

	m := resize.Resize(width, height, img, resize.NearestNeighbor)

	return m, format
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

	result, format := resizeFile(bytes.NewBuffer(image), uint(width), uint(height))

	if format == "jpeg" {
		jpeg.Encode(w, result, &jpeg.Options{Quality: 90})
	} else if format == "gif" {
		gif.Encode(w, result, nil)
	} else if format == "png" {
		png.Encode(w, result)
	}
}

func main() {

	log.Println("Starting server")

	http.HandleFunc("/resize", resizeHandler)
	http.ListenAndServe(":4500", nil)

	log.Println("Server stopped")

}
