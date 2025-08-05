# Web crawler
This is a **web crawler** I wrote. When executed, the crawler reads a `links.txt` file and spawns a goroutine for each website (up to 1000). For each website, a request is sent for its `robots.txt` file, the program then parses it and returns the crawling rules (allowed/disallowed routes, crawl delay). The crawler then crawls the website, parses the HTML, extracts all content and stores this content in a SQLite database.

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

In `links.txt`, paste in the websites you want to crawl. Note that some websites prohibit crawling, there may be websites that you have pasted in that might not be crawled, this is handled for you.

Once you have all your links, you can run either command below to start crawling:
```
go build && ./crawler
```
or
```
go run .
```
