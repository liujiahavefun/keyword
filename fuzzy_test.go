package main

import (
    //"regexp"
    "testing"
    //"sort"
    "fmt"
)

func assert1(t *testing.T, b bool) {
    if !b {
        t.Fail()
    }
}

/*
func TestLoadDict(t *testing.T) {
    //LoadDict()
    for k, v := range SPLITDICT {
        fmt.Printf("%v : %v\n", k, v)
    }
}*/

/*
func TestSplitChar(t *testing.T) {
    r1 := SplitChar("小姐姐ab")
    fmt.Println(r1)

    r2 := SplitChar("大姐")
    fmt.Println(r2)

    r3 := SplitChar("中国人")
    fmt.Println(r3)
}
*/

/*
func TestTranslateToPinyi(t *testing.T) {
    words := make([]string, 0)
    words = append(words, "中国人")
    words = append(words, "小姐姐")
    words = append(words, "小姐姐ab")
    words = append(words, "*小&姐姐 *#")
    words = append(words, "hello")
    result := TranslateToPinyi(words, true)
    fmt.Println(result)

    result2 := TranslateToPinyi(words, false)
    fmt.Println(result2)
}
*/

/*
func TestWordToPinyi(t *testing.T) {
    words := make([]string, 0)
    //words = append(words, "中国人")
    words = append(words, "小姐姐")
    //words = append(words, "小姐姐ab")
    words = append(words, "*小&姐姐 *#123Ab")
    //words = append(words, "hello")

    fmt.Println("忽略所有数字字母和其它符号")
    for _, word := range words {
        result := WordToPinyi2(word, true, true)
        fmt.Println(word)
        fmt.Println(result)
    }

    fmt.Println("不忽略所有数字字母和其它符号")
    for _, word := range words {
        result := WordToPinyi2(word, false, false)
        fmt.Println(word)
        fmt.Println(result)
    }

    fmt.Println("不忽略所有数字字母，忽略其它符号")
    for _, word := range words {
        result := WordToPinyi2(word, false, true)
        fmt.Println(word)
        fmt.Println(result)
    }
}
*/

/*
func TestInputStringSearch(t *testing.T) {
    input := NewInputString("a1好看的小*姐姐小姐姐haha小 姐姐开发阶段卡咖啡店小姐姐")
    kwpy := WordToPinyi("小姐姐", false, true)
    r := input.SearchSubStrings(kwpy)

    fmt.Println(input)
    fmt.Println(kwpy)
    fmt.Println(r)

    kwpy = WordToPinyi("小姐", false, true)
    r = input.SearchSubStrings(kwpy)

    fmt.Println(input)
    fmt.Println(kwpy)
    fmt.Println(r)

    kwpy = WordToPinyi("a小姐", false, true)
    r = input.SearchSubStrings(kwpy)

    fmt.Println(input)
    fmt.Println(kwpy)
    fmt.Println(r)
}

 */

/*
func TestSimilarity(t *testing.T) {
    kw := "小姐姐"
    words := []string {
        "小姐",
        "小姐姐",
        "小*姐*姐",
        "fuck",
    }

    for _, word := range words {
        sim := Similarity(word, kw)
        fmt.Printf("src:%v, kw:%v, sim:%v\n", word, kw, sim)
    }
}
*/

/*
func TestLoadKeywords(t *testing.T) {
   kws := map[string]string {
       "小姐姐" : "漂亮",
       "fuck" : "傻叉",
   }

   keywords, rmap := LoadKeywords(kws)
   fmt.Println(keywords)
   fmt.Println(rmap)
}
*/

func TestFuzzyMatch(t *testing.T) {
    kws := map[string]string {
        "小姐姐" : "漂亮",
        "小妹妹" : "漂亮",
        "fuck" : "傻叉",
    }

    matcher := NewFuzzyMatcher()
    err := matcher.Prepare(kws, "chaizi-jt.txt")
    if err != nil {
        t.Fatal("fuzzy matcher initialized failed!")
    }

    /*
    content := "我喜欢小姐姐"
    fmt.Printf("content: %v\n", content)
    matcher.Match(content, 0.2)
    //fmt.Printf("finaly result: %v\n", r)

    content = "我喜欢小*姐 姐"
    fmt.Printf("content: %v\n", content)
    matcher.Match(content, 0.2)
    //fmt.Printf("finaly result: %v\n", r)

    content = "我喜欢xiao姐姐"
    fmt.Printf("content: %v\n", content)
    matcher.Match(content, 0.2)
    //fmt.Printf("finaly result: %v\n", r)

    content = "我喜欢小女且姐"
    fmt.Printf("content: %v\n", content)
    matcher.Match(content, 0.2)
    //fmt.Printf("finaly result: %v\n", r)

    content = "我tm谁也不喜欢，可以么额！！！！"
    fmt.Printf("content: %v\n", content)
    matcher.Match(content, 0.2)
    //fmt.Printf("finaly result: %v\n", r)

    content = "fuck小姐姐"
    fmt.Printf("content: %v\n", content)
    matcher.Match(content, 0.2)
    //fmt.Printf("finaly result: %v\n", r) */

    content := "a1好看的小*姐姐小姐姐haha小 姐姐开发阶段卡咖啡店小姐女且小妹@妹 六放假大姐夫 fuck!!!xiao姐姐小姐姐"
    fmt.Printf("content: %v\n", content)
    matcher.Match(content, 0.2)
}

/*
func TestFuck(t *testing.T) {
    a := []int{0,1,2,3,4,5,6}
    fmt.Println(a[0:])
    fmt.Println(a[6:])
}
 */