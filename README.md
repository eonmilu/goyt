# goyt [![Build status](https://travis-ci.org/eonmilu/goyt.svg?branch=master)](https://travis-ci.org/eonmilu/goyt/)

Library used by the [Your Time Extension](https://github.com/eonmilu/yourtime).

## Installation

Run:
`go get github.com/eonmilu/goyt`

## Requirements

You will need a database with the fields:

- TimemarkID int64
- Author     string
- AuthorURL  string
- Timemark   int64
- Content    string
- Votes      int64
- Date       int64

## Use

```go
yourtime = goyt.YourTime{
    AuthTokenURL:   YOUR_GOOGLE_AUTH_TOKEN_URL,
    GoogleClientID: YOUR_GOOGLE_AUTH_CLIENT_ID,
    DB:             YOUR_DATABASE,
}


myRouter.HandleFunc(YOUR_SEARCH_URL, yourtime.Search)
myRouter.HandleFunc(YOUR_INSERT_URL, yourtime.Insert)
myRouter.HandleFunc(YOUR_AUTH_URL, yourtime.Auth)
```
