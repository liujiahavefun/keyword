package main

type KeywordMatch struct {
    Keyword      string  `json:"kw"`
}

type Response struct {
    Errno       int  `json:"errno"`
    Errmsg      string  `json:"errmsg"`
    Data        []string
}

type MatchItem struct {
    Pattern string `json:"kw"`
    Count int `json:"hits"`
    Cat string `json:"cat"`
}

type FuzzyMatchItem struct {
    Pattern string `json:"kw"`
    Cat string `json:"cat"`
    Sub string `json:sub`
    Score float32 `json:score`
    Count int `json:"hits"`
}