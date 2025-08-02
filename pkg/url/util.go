package url

import (
	"net/url"
	"strings"
)

func Normalize(rawUrl string) (string, error) {
	urlStruct, err := url.Parse(rawUrl)
	if err != nil {
		return "", err
	}

	return urlStruct.Host + strings.TrimRight(urlStruct.Path, "/"), nil
}
