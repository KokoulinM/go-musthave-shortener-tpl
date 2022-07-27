// Package models for working with the URL entity
package models

type ShortURL = string
type LongURL = string

type ShortURLs = map[ShortURL]LongURL
