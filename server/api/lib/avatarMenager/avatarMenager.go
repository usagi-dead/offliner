package avatarManager

import (
	"bytes"
	"github.com/chai2010/webp"
	"github.com/nfnt/resize"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"mime/multipart"
	"net/http"
	u "server/internal/features/user"
	"sync"
)

func ParsingAvatarImage(file *multipart.File) ([]byte, []byte, error) {
	buffer := new(bytes.Buffer)
	if _, err := io.Copy(buffer, *file); err != nil {
		return nil, nil, u.ErrInternal
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
			return nil, nil, u.ErrInvalidTypeAvatar
		}
		img, err = gif.Decode(buffer)
	case "image/webp":
		img, err = webp.Decode(buffer)
	default:
		return nil, nil, u.ErrInvalidTypeAvatar
	}

	if err != nil {
		return nil, nil, u.ErrInvalidTypeAvatar
	}

	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	if width != height {
		return nil, nil, u.ErrInvalidResolutionAvatar
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
			err512 = u.ErrInternal
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
			err52 = u.ErrInternal
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

	return buf52, buf512, nil
}

func isNonAnimatedGIF(reader io.Reader) (bool, error) {
	img, err := gif.DecodeAll(reader)
	if err != nil {
		return false, err
	}
	return len(img.Image) == 1, nil
}
