# Web crawler
This is a web crawler I wrote in Go.

## How it works

### Reading links
When the crawler is executed with `./crawler`, the program reads a `links.txt` file and spawns a Goroutine to crawl each domain. The number of Goroutines that can exist at any time is limited to a 1000, meaning the program can crawl up to a 1000 domains concurrently.

The limit on the number of Goroutines is done by using a **buffered channel**. When a Goroutine is spawned to crawl a domain, a wait group is added and an empty struct is sent onto the channel, when the program is done crawling the entirety of the domain, a signal is sent from the channel and the wait group is dropped, allowing for another Goroutine to crawl another domain if the channel was previously full.

This is repeated until all domains in the `links.txt` file have been crawled.

### Parsing a robots.txt file
When crawling a single domain, the first thing that is done is a GET request for the domain's `robots.txt file`. This file dictates which routes the program **can and cannot crawl** as well as the **crawl delay**. If the GET request returns a 403 status code, the program will not crawl the domain, if it returns a 404 then the program will crawl the domain. If a `robots.txt` file is returned, a couple of checks are done to confirm the validity of the returned file. The program will then parse it line by line to look for crawling rules that apply to it only, applicable crawling rules are then marshalled into a struct and returned.

### Parsing HTML
HTML content is retrieved using a GET request to the current route, early returns are set up if a 404 is returned or if the returned content is not HMTL. The HTML content retrieved is then tokenised. These tokens are then parsed, **only text between `<p>` tags** are extracted as relevant content. The program also looks for **links within starting `<a>` tags**. These pieces of information are then cleaned and marshalled into a struct before being returned.

### Data storage
Extracted content is stored in a **Turso database**. Since the crawler is only working with text, using a database like Turso (built on a **SQLite** fork) keeps the entire program lightweight.

### Crawling algorithm
A **Breadth First Traversal** was chosen for the crawling algorithm over a **Recursive Depth First Traversal**. Go is **not tail call optimised**, meaning on every recursive call, a new stack frame is allocated instead of reusing the current one. For a web crawler designed to crawl 1000 sites concurrently, using a recursive algorithm written in Go could lead to stack overflow. For a Breadth First Traversal, routes are put into queues, meaning the stack size grows only with the queue size, memory allocation for such an algorithm is thus much smaller compared to that of a recursive algorithm.

## Planned extensions

- [ ] An API that the crawler can send requests to to extract keywords from content and turn keywords into vector embeddings.
- [ ] A query engine that makes use of vector embeddings at search time.
