package CustomErrors

import "errors"

var (
	ErrEmailExists             = errors.New("users_email_key")
	ErrEmailNotExists          = errors.New("user_email_not_exists")
	ErrUnsupportedImageFormat  = errors.New("unsupported image format")
	ErrAnimatedGIFNotSupported = errors.New("animated GIF is not supported")
	ErrImageMustBeSquare       = errors.New("image must be square (1:1 aspect ratio)")
	ErrFileReadError           = errors.New("error reading file")
	ErrImageDecodingError      = errors.New("error decoding image")
	ErrFileSaveError           = errors.New("could not save file")
	ErrWebPEncodingError       = errors.New("error saving WebP")
)
