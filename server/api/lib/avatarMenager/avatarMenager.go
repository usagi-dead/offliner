package avatarManager

import (
	"bytes"
	"errors"
	"github.com/chai2010/webp"
	"github.com/nfnt/resize"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"mime/multipart"
	"net/http"
	"sync"
)

func ParsingAvatarImage(file *multipart.File) ([]byte, []byte, error) {
	buffer := new(bytes.Buffer)
	if _, err := io.Copy(buffer, *file); err != nil {
		return nil, nil, errors.New("")
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
			return nil, nil, errors.New("")
		}
		img, err = gif.Decode(buffer)
	case "image/webp":
		img, err = webp.Decode(buffer)
	default:
		return nil, nil, errors.New("")
	}

	if err != nil {
		return nil, nil, errors.New("")
	}

	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	if width != height {
		return nil, nil, errors.New("")
	}

	var wg sync.WaitGroup
	var buf512, buf52 []byte
	var err512, err52 error

	// Обработка 512x512
	wg.Add(1)
	go func() {
		defer wg.Done()
		resized := resize.Resize(512, 512, img, resize.Lanczos3)
		buffer := new(bytes.Buffer)
		if err := webp.Encode(buffer, resized, &webp.Options{Quality: 80}); err != nil {
			err512 = errors.New("")
			return
		}
		buf512 = buffer.Bytes()
	}()

	// Обработка 52x52
	wg.Add(1)
	go func() {
		defer wg.Done()
		resized := resize.Resize(52, 52, img, resize.Lanczos3)
		buffer := new(bytes.Buffer)
		if err := webp.Encode(buffer, resized, &webp.Options{Quality: 80}); err != nil {
			err52 = errors.New("")
			return
		}
		buf52 = buffer.Bytes()
	}()

	// Ожидание завершения всех горутин
	wg.Wait()

	// Проверка на ошибки
	if err512 != nil {
		return nil, nil, err512
	}
	if err52 != nil {
		return nil, nil, err52
	}

	return buf512, buf52, nil
}

func isNonAnimatedGIF(reader io.Reader) (bool, error) {
	img, err := gif.DecodeAll(reader)
	if err != nil {
		return false, err
	}
	return len(img.Image) == 1, nil
}
