package avatarMenager

import (
	"bytes"
	"fmt"
	"github.com/chai2010/webp"
	"github.com/google/uuid"
	"github.com/nfnt/resize"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"server/internal/lib/CustomErrors" // Импортируйте ваш пакет кастомных ошибок
)

const avatarDir = "./avatars"

func SaveAvatar(file *multipart.File) (string, error) {
	buffer := new(bytes.Buffer)
	if _, err := io.Copy(buffer, *file); err != nil {
		return "", CustomErrors.ErrFileReadError
	}

	var img image.Image
	var err error
	contentType := http.DetectContentType(buffer.Bytes())

	switch contentType {
	case "image/png":
		img, err = png.Decode(buffer)
	case "image/jpeg":
		img, err = jpeg.Decode(buffer)
	case "image/gif":
		isNonAnimated, err := isNonAnimatedGIF(bytes.NewReader(buffer.Bytes()))
		if err != nil || !isNonAnimated {
			return "", CustomErrors.ErrAnimatedGIFNotSupported
		}
		img, err = gif.Decode(buffer)
	case "image/webp":
		img, err = webp.Decode(buffer)
	case "image/vnd.microsoft.icon":
		return "", CustomErrors.ErrUnsupportedImageFormat
	default:
		return "", CustomErrors.ErrUnsupportedImageFormat
	}

	if err != nil {
		return "", CustomErrors.ErrImageDecodingError
	}

	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	if width != height {
		return "", CustomErrors.ErrImageMustBeSquare
	}

	resizedImg := resize.Resize(512, 512, img, resize.Lanczos3)

	uniqueFileName := fmt.Sprintf("%s.webp", uuid.New().String())
	fileName := filepath.Join(avatarDir, uniqueFileName)

	out, err := os.Create(fileName)
	if err != nil {
		return "", CustomErrors.ErrFileSaveError
	}
	defer out.Close()

	if err := webp.Encode(out, resizedImg, &webp.Options{Quality: 80}); err != nil {
		return "", CustomErrors.ErrWebPEncodingError
	}

	return "http://localhost:8080/user/avatar?filename=" + uniqueFileName, nil
}

func isNonAnimatedGIF(reader io.Reader) (bool, error) {
	img, err := gif.DecodeAll(reader)
	if err != nil {
		return false, err
	}
	return len(img.Image) == 1, nil
}
