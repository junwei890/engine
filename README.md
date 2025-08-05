# Web crawler
This is a **web crawler** I wrote in Go.

When executed, the crawler reads a `links.txt` file containing all the domains to be crawled and spawns a goroutine (up to a 1000) to crawl each domain. For each domain, a request is sent for its `robots.txt` file, the program then parses it and returns crawling rules. The crawler then crawls the domain, making sure to avoid disallowed routes and to crawl with the specified crawl delay. For each route, the crawler parses the HTML, extracts all content and stores this content in a SQLite database.

## Requirements
To try out the crawler locally, you need:
- [Go](https://go.dev/doc/install) installed
- [Goose](https://github.com/pressly/goose) installed
- Access to [Turso](https://docs.turso.tech/quickstart)

## Usage
Clone the repo using `git clone https://github.com/junwei890/crawler.git`

Create a `.env` file and paste your database URL and auth token in. It should look something like this:
```
DB_URL=<database URL>?authToken=<auth token>
```

Run database migrations on Turso, I've provided some scripts that you could run:
```
./scripts/upmigration.sh
```

In `links.txt`, paste in the websites you want to crawl. Note that some domains prohibit crawling, there may be domains that you have pasted in that might not be crawled, this is handled for you.

Once you have all your links, you can run either command below to start crawling:
```
go build && ./crawler
```
or
```
go run .
```
