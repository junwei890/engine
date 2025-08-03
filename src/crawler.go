package src

import (
	"net/url"
	"time"

	"github.com/junwei890/engine/utils"
)

func Crawler(startURL string) error {
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

		time.Sleep(time.Duration(rules.Delay) * time.Second)
	}

	return nil
}
