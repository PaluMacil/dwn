# gamescrape

This is a commandline tool for scraping gameboardgeek.com. 
There are a few wiki pages with information on scraping the 
site, so the owners seem to consider this to be a permitted 
activity.

## Commands

### default

The subcommand `default` may be omitted or included. This 
command is used to build an index of game IDs and game names. 
If there are errors, they are available in the errors.json 
file.

```
gamescrape
```

Available options include:

- `-wait=#`: Set the wait to a whole number of seconds to wait 
between requests. If set to `0`, this tool will not pause and 
will also run multiple consecutive requests. The default value 
for this option is `5`.
- `-limit=500`: Set the max number of pages to check (the 
default, 0 is no limit)

After scraping, the error.json file is truncated and rewritten 
with any errors that occurred during this iteration of 
scraping.

### retry

The subcommand `retry` loads errors from `error.json` and 
retries scraping for each page.

```
gamescrape retry
```

While the `-wait` flag still works the same as before, the 
`-limit` flag is not recommended since items that aren't scanned will be lost when the errors from this scrape are written to 
disk.

## Roadmap

A postback url flag will be created to allow for the dwn 
webserver to call this utility with a way to notify the server 
of task completion.