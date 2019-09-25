package main

import (
    "net/url"
    "bufio"
    "strings"
    "os"
    "io"
    "fmt"
    "github.com/gin-gonic/gin"
)

type Key struct {
    word string
    cat  string
}

var (
    matcher *Matcher
    fuzzyMatcher *FuzzyMatcher
    key2cat = make(map[string]string)
)

func init() {
    patterns, err := ReadKeywords("words.txt")
    if err != nil {
        os.Exit(-1)
    }

    matcher = NewStringMatcher(patterns)

    fuzzyMatcher = NewFuzzyMatcher()
    err = fuzzyMatcher.Prepare(key2cat, "chaizi-jt-simple.txt")
    if err != nil {
        os.Exit(-1)
    }
}

func ReadKeywords(fileName string) (patterns []string, err error) {
    f, err := os.Open(fileName)
    if err != nil {
        return nil, err
    }
    patterns = make([]string, 0)
    buf := bufio.NewReader(f)
    for {
        line, err := buf.ReadString('\n')
        if err != nil {
            if err == io.EOF {
                return patterns, nil
            }
            return patterns, err
        }
        line = strings.TrimSpace(line)
        words := strings.Split(line, "\t")
        if len(words) != 2 {
            continue
        }

        patterns = append(patterns, words[0])
        key2cat[words[0]] = words[1]
    }
    return patterns, nil
}

func queryStringGetDefault(values url.Values, key, defaultValue string) string {
    v := values.Get(key)
    if v == "" {
        return defaultValue
    } else {
        return v
    }
}

func keywordMatchHandler(c *gin.Context)  {
    var jsonData gin.H = nil
    defer func() {
        c.Keys["jsonData"] = jsonData
    }()

    kw := c.Query("kw") // equals to c.Request.URL.Query().Get("kw")
    if kw == "" && c.Request.Method == "POST" {
        c.Request.ParseForm()
        kw = queryStringGetDefault(c.Request.Form, "kw", "")

        if kw == "" {
            var param KeywordMatch
            if err := c.ShouldBindJSON(&param); err == nil {
                kw = param.Keyword
            }
        }
    }

    if kw == "" {
        jsonData = gin.H{
            "errno": 1001,
            "errmsg": "invalid param",
            "data": "no keyword",
        }
        return
    }

    hits := matcher.Match([]byte(kw))
    patterns := matcher.Index2Pattern(hits)
    data := make([]*MatchItem, 0)
    for pattern, count := range patterns {
        item := &MatchItem{
            Pattern:pattern,
            Count:count,
        }
        if cat, exist := key2cat[pattern];exist {
            item.Cat = cat
        }else {
            item.Cat = "其它"
        }
        data = append(data, item)
    }

    jsonData = gin.H{
        "errno": 0,
        "errmsg": "ok",
        "data": data,
    }

    return
}

func fileKeywordMatchHandler(c *gin.Context) {
    var jsonData gin.H = nil
    defer func() {
        c.Keys["jsonData"] = jsonData
    }()

    file, err := c.FormFile("file")
    if err != nil {
        jsonData = gin.H{
            "errno": 1002,
            "errmsg": fmt.Sprintf("get file failed, %v", err),
            "data": nil,
        }
        return
    }

    fr := NewFileReader(file.Filename)
    if !fr.Support() {
        jsonData = gin.H{
            "errno": 1003,
            "errmsg": fmt.Sprintf("unsupported file type"),
            "data": nil,
        }
        return
    }

    content, err := fr.Read()
    if err != nil {
        jsonData = gin.H{
            "errno": 1004,
            "errmsg": fmt.Sprintf("parse file failed, %v", err),
            "data": nil,
        }
        return
    }

    data := make([]*MatchItem, 0)
    hits := matcher.Match([]byte(content))
    patterns := matcher.Index2Pattern(hits)
    for pattern, count := range patterns {
        item := &MatchItem{
            Pattern:pattern,
            Count:count,
        }
        if cat, exist := key2cat[pattern];exist {
            item.Cat = cat
        }else {
            item.Cat = "其它"
        }
        data = append(data, item)
    }

    jsonData = gin.H{
        "errno": 0,
        "errmsg": "ok",
        "data": data,
    }

    return
}

func fuzzyMatchHandler(c *gin.Context)  {
    var jsonData gin.H = nil
    defer func() {
        c.Keys["jsonData"] = jsonData
    }()

    kw := c.Query("kw")
    if kw == "" && c.Request.Method == "POST" {
        c.Request.ParseForm()
        kw = queryStringGetDefault(c.Request.Form, "kw", "")

        if kw == "" {
            var param KeywordMatch
            if err := c.ShouldBindJSON(&param); err == nil {
                kw = param.Keyword
            }
        }
    }

    if kw == "" {
        jsonData = gin.H{
            "errno": 1001,
            "errmsg": "invalid param",
            "data": "no keyword",
        }
        return
    }

    results, err := fuzzyMatcher.Match(kw, 0.2)
    if err != nil {
        jsonData = gin.H{
            "errno": 1002,
            "errmsg": "something bad happened...",
            "data": "no keyword",
        }
        return
    }

    data := make([]*FuzzyMatchItem, 0)
    for _, v := range results {
        for _, result := range v {
            item := &FuzzyMatchItem{
                Pattern:result.Kw,
                Cat:result.Cat,
                Sub:result.Sub,
                Score:result.Score,
                Count:result.Count,
            }
            data = append(data, item)
        }
    }

    jsonData = gin.H{
        "errno": 0,
        "errmsg": "ok",
        "data": data,
    }

    return
}