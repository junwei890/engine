package utils

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"mime"
	"net/http"
	"net/url"
	"slices"
	"strconv"
	"strings"

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
	client := &http.Client{}

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
		return []byte{}, errors.New("content type not text/html")
	}

	page, err := io.ReadAll(res.Body)
	if err != nil {
		return []byte{}, err
	}

	return page, nil
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
							fullURL := domain.ResolveReference(structure).String()

							if comp := slices.Contains(response.Links, fullURL); !comp {
								response.Links = append(response.Links, fullURL)
							}
						} else {
							if comp := slices.Contains(response.Links, attr.Val); !comp {
								response.Links = append(response.Links, attr.Val)
							}
						}
					}
				}
			}
		}
	}

	return response, nil
}

func GetRobots(rawURL string) ([]byte, error) {
	client := &http.Client{}

	res, err := client.Get(fmt.Sprintf("%srobots.txt", rawURL))
	if err != nil {
		return []byte{}, err
	}
	defer res.Body.Close()

	if res.StatusCode == 403 {
		return []byte{}, errors.New("can't scrape")
	}
	if res.StatusCode == 404 {
		return []byte{}, nil
	}

	mediaType, _, err := mime.ParseMediaType(res.Header.Get("Content-Type"))
	if err != nil {
		return []byte{}, err
	}
	if mediaType != "text/plain" {
		return []byte{}, errors.New("content type not text/plain")
	}

	textFile, err := io.ReadAll(res.Body)
	if err != nil {
		return []byte{}, err
	}

	return textFile, nil
}

type Rules struct {
	Allowed    []string
	Disallowed []string
	Delay      int
}

func ParseRobots(normURL string, textFile []byte) (Rules, error) {
	rules := Rules{}

	scanner := bufio.NewScanner(bytes.NewReader(textFile))

	applicable := false

	for scanner.Scan() {
		if strings.TrimSpace(scanner.Text()) == "" || strings.HasPrefix(strings.TrimSpace(scanner.Text()), "#") {
			continue
		}

		line := strings.Split(scanner.Text(), ":")
		key := strings.TrimSpace(line[0])
		value := strings.TrimSpace(line[1])

		if key == "User-agent" {
			if value == "*" {
				applicable = true
			} else {
				applicable = false
			}
		}

		if applicable == true {
			switch key {
			case "Allow":
				if strings.HasPrefix(value, "/") {
					rules.Allowed = append(rules.Allowed, fmt.Sprintf("%s%s", normURL, value))
				}
			case "Disallow":
				if strings.HasPrefix(value, "/") {
					rules.Disallowed = append(rules.Disallowed, fmt.Sprintf("%s%s", normURL, value))
				}
			case "Crawl-delay":
				delay, err := strconv.Atoi(value)
				if err != nil {
					rules.Delay = 0
				} else {
					rules.Delay = delay
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Println(err)
	}

	return rules, nil
}

type Queue []string

type QueueOps interface {
	Enqueue(string)
	Dequeue() string
	Peek() string
	Empty() bool
	Size() int
}

func (q *Queue) Enqueue(url string) {
	*q = append(*q, url)
}

func (q *Queue) Dequeue() string {
	popped := (*q)[0]

	*q = slices.Delete(*q, 0, 1)
	return popped
}

func (q *Queue) Peek() string {
	return (*q)[0]
}

func (q *Queue) Empty() bool {
	return len(*q) == 0
}

func (q *Queue) Size() int {
	return len(*q)
}
