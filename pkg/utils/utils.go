package utils

import (
	"bytes"
	"errors"
	"io"
	"log"
	"mime"
	"net/http"
	"net/url"
	"strings"
	"time"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

func Normalize(rawURL string) (string, error) {
	structure, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}

	return structure.Host + strings.TrimRight(structure.Path, "/"), nil
}

func GetHTML(rawURL string) ([]byte, error) {
	client := &http.Client{
		Timeout: 3 * time.Second,
	}

	res, err := client.Get(rawURL)
	if err != nil {
		return []byte{}, err
	}

	defer res.Body.Close()

	if res.StatusCode >= 400 && res.StatusCode < 500 {
		return []byte{}, errors.New("400+ status code")
	}

	mediaType, _, err := mime.ParseMediaType(res.Header.Get("Content-Type"))
	if err != nil {
		return []byte{}, err
	}

	if mediaType != "text/html" {
		return []byte{}, errors.New("content type not html")
	}

	content, err := io.ReadAll(res.Body)
	if err != nil {
		return []byte{}, err
	}

	return content, nil
}

type Response struct {
	Content []string
	Links   []string
}

func ParseHTML(domain *url.URL, page []byte) (Response, error) {
	response := Response{}

	tokens := html.NewTokenizer(bytes.NewReader(page))

	for {
		tn := tokens.Next()

		if tn == html.ErrorToken {
			if tokens.Err() == io.EOF {
				break
			}

			return response, errors.New("couldn't tokenize")
		}

		if tn == html.TextToken {
			t := tokens.Token()
			clean := strings.ToLower(strings.Join(strings.Fields(t.Data), " "))
			if clean != "" {
				response.Content = append(response.Content, clean)
			}
			continue
		}

		if tn == html.StartTagToken {
			t := tokens.Token()
			if t.Data == "a" && t.DataAtom == atom.A {
				for _, attr := range t.Attr {
					if attr.Key == "href" {
						structure, err := url.Parse(attr.Val)
						if err != nil {
							log.Println("invalid url")
							continue
						}

						if structure.Hostname() == "" {
							response.Links = append(response.Links, domain.ResolveReference(structure).String())
						} else {
							response.Links = append(response.Links, attr.Val)
						}
					}
				}
			}
		}
	}

	return response, nil
}
