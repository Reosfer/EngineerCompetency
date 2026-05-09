package helper

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"os"
	"strings"

	"github.com/google/uuid"
	"github.com/nfnt/resize"
)

func Base64ToImage(data string, outputPath string, outputFilename string) (string, error) {
	idx := strings.Index(data, ";base64,")

	if idx < 0 {
		newErr := NewError("Invalid Image")
		return "", newErr
	}

	imageType := data[11:idx]
	unbased, err := base64.StdEncoding.DecodeString(data[idx+8:])

	if err != nil {
		return "", err
	}

	r := bytes.NewReader(unbased)
	newFilename := fmt.Sprintf("%s/%s.%s", outputPath, outputFilename, imageType)

	switch imageType {
	case "png":
		im, err := png.Decode(r)

		if err != nil {
			return "", err
		}

		f, err := os.OpenFile(newFilename, os.O_WRONLY|os.O_CREATE, 0777)
		if err != nil {
			return "", err
		}

		png.Encode(f, im)

		width, height := GetImageDimension(newFilename)
		var maxWidth, maxHeight int

		if width > height {
			maxWidth = 1024
			maxHeight = 550
		} else {
			maxHeight = 1024
			maxWidth = 550
		}

		ratioWidth := float64(maxWidth) / float64(width)
		ratioHeight := float64(maxHeight) / float64(height)

		finalWidth := ratioWidth * float64(width)
		finalHeight := ratioHeight * float64(height)

		file, err := os.Open(newFilename)
		if err != nil {
			return "", err
		}

		img, err := jpeg.Decode(file)
		if err != nil {
			return "", err
		}

		file.Close()

		m := resize.Resize(uint(finalWidth), uint(finalHeight), img, resize.Lanczos3)
		out, err := os.Create(newFilename)
		if err != nil {
			return "", err
		}

		png.Encode(out, m)

		break
	case "jpeg":
		im, err := jpeg.Decode(r)

		if err != nil {
			return "", err
		}

		f, err := os.OpenFile(newFilename, os.O_WRONLY|os.O_CREATE, 0777)

		if err != nil {
			return "", err
		}

		jpeg.Encode(f, im, nil)

		width, height := GetImageDimension(newFilename)
		var maxWidth, maxHeight int

		if width > height {
			maxWidth = 1024
			maxHeight = 550
		} else {
			maxHeight = 1024
			maxWidth = 550
		}

		ratioWidth := float64(maxWidth) / float64(width)
		ratioHeight := float64(maxHeight) / float64(height)

		finalWidth := ratioWidth * float64(width)
		finalHeight := ratioHeight * float64(height)

		file, err := os.Open(newFilename)
		if err != nil {
			return "", err
		}

		img, err := jpeg.Decode(file)
		if err != nil {
			return "", err
		}

		file.Close()

		m := resize.Resize(uint(finalWidth), uint(finalHeight), img, resize.Lanczos3)
		out, err := os.Create(newFilename)
		if err != nil {
			return "", err
		}

		defer out.Close()

		jpeg.Encode(out, m, nil)
		break
	case "gif":
		im, err := gif.Decode(r)

		if err != nil {
			return "", err
		}

		f, err := os.OpenFile(newFilename, os.O_WRONLY|os.O_CREATE, 0777)

		if err != nil {
			return "", err
		}

		gif.Encode(f, im, nil)

		width, height := GetImageDimension(newFilename)
		var maxWidth, maxHeight int

		if width > height {
			maxWidth = 1024
			maxHeight = 550
		} else {
			maxHeight = 1024
			maxWidth = 550
		}

		ratioWidth := float64(maxWidth) / float64(width)
		ratioHeight := float64(maxHeight) / float64(height)

		finalWidth := ratioWidth * float64(width)
		finalHeight := ratioHeight * float64(height)

		file, err := os.Open(newFilename)
		if err != nil {
			return "", err
		}

		img, err := jpeg.Decode(file)
		if err != nil {
			return "", err
		}

		file.Close()

		m := resize.Resize(uint(finalWidth), uint(finalHeight), img, resize.Lanczos3)

		out, err := os.Create(newFilename)
		if err != nil {
			return "", err
		}

		gif.Encode(out, m, nil)

		break
	}

	return newFilename, nil
}

func GetMime(fileExtension string) string {
	var mime string

	if fileExtension == "jpg" || fileExtension == "jpeg" || fileExtension == "png" || fileExtension == "gif" {
		mime = "image"
	} else {
		mime = "application"
	}

	return mime
}

func GetUniqueImageName() string {
	uniqueID := uuid.New()
	outputFilename := fmt.Sprintf("%s", uniqueID.String())
	return outputFilename
}

func GetExtensionImage(imageName string) string {
	fmt.Println(imageName)
	imageExplode := strings.Split(imageName, ".")
	return imageExplode[1]
}

func GetImageDimension(imagePath string) (int, int) {
	file, err := os.Open(imagePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
	}

	image, _, err := image.DecodeConfig(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %v\n", imagePath, err)
	}
	return image.Width, image.Height
}
