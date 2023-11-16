# aggregator

A simple RSS aggregator API written in Go.

## Stack

* Go
* Postgres

## About

Allow a user to add source of feeds. Then fetch the posts from the feed concurrently in the background at a 
set interval.

The feeds are marked updated when fetched.

You can also fetch posts by user, that checks the `feeds_lookup` table to reference the posts by user.


## TODO:

* Support pagination of the endpoints that can return many items
* Support different options for sorting and filtering posts using query parameters
* Classify different types of feeds and posts (e.g. blog, podcast, video, etc.)
* Scrape lists of feeds themselves from a third-party site that aggregates feed URLs
* Add support for other types of feeds (e.g. Atom, JSON, etc.)
* Add integration tests that use the API to create, read, update, and delete feeds and posts (Could be done via Postman)
* Add bookmarking or "liking" to posts
* Create a simple web UI that uses the API
