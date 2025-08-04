package src

import (
	"context"
	"log"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/junwei890/crawler/internal/database"
	"github.com/junwei890/crawler/utils"
)

func Init(queries *database.Queries) error {
	file, err := os.ReadFile("links.txt")
	if err != nil {
		return err
	}
	links := strings.Fields(string(file))

	wg := &sync.WaitGroup{}
	channel := make(chan struct{}, 1000)

	for _, link := range links {
		wg.Add(1)
		channel <- struct{}{}
		go func() {
			defer func() {
				<-channel
				wg.Done()
			}()
			if err := crawler(link, queries); err != nil {
				log.Println(err)
				return
			}
		}()
	}
	wg.Wait()

	return nil
}

func crawler(startURL string, queries *database.Queries) error {
	file, err := utils.GetRobots(startURL)
	if err != nil {
		return err
	}

	normURL, err := utils.Normalize(startURL)
	if err != nil {
		return err
	}

	rules, err := utils.ParseRobots(normURL, file)
	if err != nil {
		return err
	}

	dom, err := url.Parse(startURL)
	if err != nil {
		return err
	}

	visited := map[string]struct{}{}
	queue := &utils.Queue{}
	queue.Enqueue(startURL)

	for {
		if comp := queue.CheckEmpty(); comp {
			log.Println("here")
			break
		}

		popped, err := queue.Dequeue()
		if err != nil {
			return err
		}

		ok, err := utils.CheckDomain(dom, popped)
		if err != nil {
			continue
		}
		if !ok {
			continue
		}

		currURL, err := utils.Normalize(popped)
		if err != nil {
			continue
		}

		ok = utils.CheckAbility(visited, rules, currURL)
		if !ok {
			continue
		}

		page, err := utils.GetHTML(popped)
		if err != nil {
			continue
		}

		res, err := utils.ParseHTML(dom, page)
		if err != nil {
			continue
		}

		for _, link := range res.Links {
			queue.Enqueue(link)
		}

		returned, err := queries.InsertData(context.TODO(), database.InsertDataParams{
			Url:       popped,
			Content:   strings.Join(res.Content, " "),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		})

		log.Println(returned)

		time.Sleep(time.Duration(rules.Delay) * time.Second)
	}

	return nil
}
