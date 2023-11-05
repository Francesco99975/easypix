package convert

import (
	"errors"
	"image"
	"mime/multipart"
	"os"

	"github.com/chai2010/webp"
)

func ToWebP(file *multipart.File) (*os.File, error) {
	outputPath := "output.webp"
    outputFile, err := os.CreateTemp(outputPath, "*")
	if err != nil {
        return &os.File{}, errors.New("unable to create the output file")
    }

	// Decode the input image
    img, _, err := image.Decode(*file)
    if err != nil {
        return &os.File{}, errors.New("unable to decode the image file")
    }

	// Encode the image to WebP format
    err = webp.Encode(outputFile, img, nil)
    if err != nil { 
        return &os.File{}, errors.New("unable to encode the image to WebP format")
    }

	return outputFile, err
} 