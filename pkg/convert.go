package convert

import (
	"fmt"
	"image"
	_ "image/jpeg" // For JPEG format
	_ "image/png"  // For PNG format
	"mime/multipart"
	"os"

	"github.com/chai2010/webp"
)

func ToWebP(file *multipart.File, uploadFile *os.File) error{
	// Decode the input image
    img, _, err := image.Decode(*file)
    if err != nil {
        return fmt.Errorf("unable to decode the image file: %s", err.Error())
    }

	// Encode the image to WebP format
    err = webp.Encode(uploadFile, img, nil)
    if err != nil { 
        return fmt.Errorf("unable to encode the image to WebP format: %s", err.Error())
    }

	return err
} 