# Web crawler
This is a web crawler I built that is polite and stack-safe.

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

In `links.txt`, paste in the websites you want to crawl. Note that some websites prohibit crawling, there may be websites that you have pasted in that might not be crawled. You can crawl up to 1000 websites concurrently.

Once you have all your links, you can run either command below to start crawling:
```
go build && ./crawler
```
or
```
go run .
```

## Design decisions
### Robots.txt parser
This is a simple parser I wrote that takes care of the most important parts of a `robots.txt` file, namely, the routes that you're allowed to scrape, routes that you're not allowed to scrape and crawl delay

Notes:
- If the GET request for the site's robots.txt file returns a 403, the website will not be crawled.
- If somehow there are identical `Allowed` and `Disallowed` routes, the route under `Allowed` takes precedence.
- If the syntax for routes under either `Allowed` and `Disallowed` are invalid, the line in the `robots.txt` file is ignored, meaning the route will be crawled.
- If `/route/` is under `Disallowed` while `/route/maps` is under `Allowed`, all subpaths to `route` can't be crawled except for maps.
- Pattern matching is supported, the route `/*world` is matched with `/helloworld`, for example.

### Crawling algorithm
A breadth first traversal was chosen over a recursive depth first traversal to prevent stack overflow. Go is **not** Tail Call Optimized, meaning on every recursive call, a new stack frame is allocated instead of reusing the current one. There **may** be sites so large that our program exceeds the stack limit.

Though stack safety is nice, choosing breadth first over depth first traversal means that I can't concurrently crawl routes within a domain, because I could break out of the loop in the subsequent iteration before the current iteration queues more links.

## Possible extensions
These are the possible extensions I could pursue:
- [ ] Extracting keywords from content and putting them in an inverted index
- [ ] Converting these keywords to vectors and using cosine similarity to query relevant documents

## Contributions
Are welcome.
