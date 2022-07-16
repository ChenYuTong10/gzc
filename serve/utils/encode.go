package utils

import (
	"errors"
	"github.com/ChenYuTong10/chardet"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/encoding/unicode"
	"strings"
)

var ErrUnSupportedEncoding = errors.New("unsupported encoding")

// EncodeUTF8 dispatches to different handler according to the encoding.
// If the encoding has been UTF8, it returns directly. If the encoding is
// out of ANSI, UTF8, BOM UTF8, UTF16 BE/LE, it returns Unexpected Error.
func EncodeUTF8(stream []byte, encoding string) ([]byte, error) {

	switch strings.ToUpper(encoding) {
	case chardet.BOM_UTF8:
		return BomUTF8ToUTF8(stream), nil
	case chardet.ANSI:
		decodeStream, err := AnsiToUTF8(stream)
		if err != nil {
			// TODO: log error
			return nil, err
		}
		return decodeStream, nil
	case chardet.BOM_UTF16_BE:
		decodeStream, err := UTF16BEToUTF8(stream)
		if err != nil {
			// TODO: log error
			return nil, err
		}
		return decodeStream, nil
	case chardet.BOM_UTF16_LE:
		decodeStream, err := UTF16LEToUTF8(stream)
		if err != nil {
			// TODO: log error
			return nil, err
		}
		return decodeStream, nil
	case chardet.UTF8:
		return stream, nil
	}
	return nil, ErrUnSupportedEncoding
}

// AnsiToUTF8 transforms encoding from Ansi to UTF8.
// It is underlying with UTF8ToBomUTF8.
func AnsiToUTF8(stream []byte) ([]byte, error) {
	decoder := simplifiedchinese.GBK.NewDecoder()
	decodeStream, err := decoder.Bytes(stream)
	if err != nil {
		return nil, err
	}
	return decodeStream, nil
}

// UTF16BEToUTF8 transforms encoding from UTF16 BE to UTF8.
func UTF16BEToUTF8(stream []byte) ([]byte, error) {
	decoder := unicode.UTF16(unicode.BigEndian, unicode.ExpectBOM).NewDecoder()
	decodeStream, err := decoder.Bytes(stream)
	if err != nil {
		return nil, err
	}
	return decodeStream, nil
}

// UTF16LEToUTF8 transforms encoding from UTF16 LE to UTF8.
func UTF16LEToUTF8(stream []byte) ([]byte, error) {
	decoder := unicode.UTF16(unicode.LittleEndian, unicode.ExpectBOM).NewDecoder()
	decodeStream, err := decoder.Bytes(stream)
	if err != nil {
		return nil, err
	}
	return decodeStream, nil
}

// BomUTF8ToUTF8 cut first three bytes BOM prefix of the stream.
func BomUTF8ToUTF8(stream []byte) []byte {
	return stream[3:]
}
