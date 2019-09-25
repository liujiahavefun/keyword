package main

import (
    "bufio"
    "errors"
    "io"
    "os"
    "strings"
    "fmt"
    "github.com/mozillazg/go-pinyin"
    "unicode"
)

func WordToPinyi(word string, ignoreAlphaAndNumber bool, ignoreOther bool) []string {
    src := []rune(word)
    ret := make([]string, 0)
    conf := pinyin.NewArgs()
    for _, ch := range src {
        py := pinyin.LazyPinyin(string(ch), conf)
        if len(py) > 0 && len(py[0]) > 0 {
            ret = append(ret, py[0])
        }else if !ignoreAlphaAndNumber {
            if unicode.IsDigit(ch) || unicode.IsNumber(ch) || unicode.IsLetter(ch) {
                ret = append(ret, string(ch))
            }else if !ignoreOther && unicode.IsGraphic(ch) {
                ret = append(ret, string(ch))
            }
        }
    }

    return ret
}

func Similarity(src, kw string) float32 {
    kwrune := []rune(kw)
    srcrune := []rune(src)

    match := 0
    for _, kwch := range kwrune {
        for _, srcch := range srcrune {
            if kwch == srcch {
                match++
                if match == len(srcrune) {
                    return 1.0
                }
                break
            }
        }
    }
    score := float32(match)/float32(len(srcrune))
    // 补偿一下拼音完全一样的。。。。
    if score > 0.1 && score < 0.5 {
        if strings.Join(WordToPinyi(src, false, true), "") == strings.Join( WordToPinyi(kw, false, true), "") {
            score += 0.15
        }
    }

    return score
}


type InputString struct {
    Src string
    RuneSrc []rune
    PinyinSrc []string
    MaskSrc []int
}

func NewInputString(src string) *InputString {
    input := &InputString{
        Src : src,
    }
    input.RuneSrc = []rune(src)
    input.PinyinSrc = make([]string, len(input.RuneSrc))
    input.MaskSrc = make([]int, len(input.RuneSrc))

    conf := pinyin.NewArgs()
    for i, ch := range input.RuneSrc {
        py := pinyin.LazyPinyin(string(ch), conf)
        if len(py) > 0 && len(py[0]) > 0 {
            input.PinyinSrc[i] = py[0]
            input.MaskSrc[i] = 1
        }else {
            if unicode.IsDigit(ch) || unicode.IsNumber(ch) || unicode.IsLetter(ch) {
                input.PinyinSrc[i] = string(ch)
                input.MaskSrc[i] = 1
            }else {
                input.PinyinSrc[i] = string(ch)
                input.MaskSrc[i] = 0
            }
        }
    }
    return input
}

func (self InputString) GetPinyinPattern() (py string) {
    py = ""
    for i, src := range self.PinyinSrc {
        if self.MaskSrc[i] == 1 {
            py += src
        }
    }
    return py
}

func (self InputString) SearchSubStrings(py []string) (subs []string) {
    subs = make([]string, 0)
    if len(py) == 0 {
        return subs
    }

    subPy := strings.Join(py, "")
    fullPy := self.GetPinyinPattern()
    start := 0
    end := 0
    lastPos := 0

    mask := make([]int, 0)
    for index, pysrc := range self.PinyinSrc {
        if self.MaskSrc[index] == 0 {
            continue
        }
        for i := 0; i < len(pysrc); i++ {
            mask = append(mask, index)
        }
    }

    for {
        if lastPos > len(fullPy) -1 {
            return subs
        }
        pos := strings.Index(fullPy[lastPos:], subPy)
        if pos == -1 {
            return subs
        }
        pos += lastPos
        if pos >= 0 && pos < len(mask) && pos+len(subPy)-1 < len(mask) {
            start = mask[pos]
            end = mask[pos+len(subPy)-1]
            if start >= 0 && end < len(self.RuneSrc) {
                s := self.RuneSrc[start:end+1]
                subs = append(subs, string(s))
            }
        }
        lastPos = pos + len(subPy)
    }

    return subs
}

type KeyWord struct {
    Kw string
    Cat string
    Splits []string
    Pinyins [][]string
}

func NewKeyWord(word, cat string) *KeyWord {
    kw := &KeyWord{
        Kw : word,
        Cat : cat,
    }

    return kw
}

func (self *KeyWord) Prepare(dict map[string]string) {
    self.Splits = self.SplitChar(self.Kw, dict)
    self.Pinyins = make([][]string, len(self.Splits))
    for i, word := range self.Splits {
        self.Pinyins[i] = WordToPinyi(word, false, true)
    }
}

func (self KeyWord) GetPatterns() []string {
    patterns := make([]string, len(self.Pinyins))
    for i, pinyin := range self.Pinyins {
        patterns[i] = strings.Join(pinyin, "")
    }
    return patterns
}

func (self *KeyWord) SplitChar(s string, dict map[string]string) (splits []string) {
    if len(s) == 0 {
        return nil
    }
    splits = make([]string, 0)
    for _, ch := range s {
        var new_split []string
        if val, exist := dict[string(ch)]; exist {
            new_split = make([]string, len(splits))
            copy(new_split, splits)
            if len(new_split) == 0 {
                new_split = append(new_split, val)
            }else {
                for i, _ := range new_split {
                    new_split[i] = new_split[i] + val
                }
            }
        }

        if len(splits) == 0 {
            splits = append(splits, string(ch))
        }else {
            for i, _ := range splits {
                splits[i] = splits[i] + string(ch)
            }
        }
        if len(new_split) > 0 {
            splits = append(splits, new_split...)
        }
    }
    return splits
}

type FuzzyResult struct {
    Kw string
    Cat string
    Sub string
    Score float32
    Count int
}

func (self FuzzyResult) String() string {
    return fmt.Sprintf("keyword [%v]:[%v], found match string [%v] %d times with score %f", self.Kw, self.Cat, self.Sub, self.Count, self.Score)
}

type FuzzyMatcher struct {
    kws []*KeyWord
    rmap map[string][]*KeyWord
    splitDict map[string]string
    matcher *Matcher
}

func NewFuzzyMatcher() *FuzzyMatcher {
    matcher := &FuzzyMatcher{}
    return matcher
}

func (self *FuzzyMatcher) Prepare(kwords map[string]string, chaiziFile string) (err error) {
    err = self.LoadSplitDict(chaiziFile)
    if err != nil {
        return err
    }

    err = self.LoadKeywords(kwords)
    if err != nil {
        return err
    }

    patterns := make([]string, 0)
    for k, _ := range self.rmap {
        patterns = append(patterns, k)
    }

    self.matcher = NewStringMatcher(patterns)

    return nil
}

func (self *FuzzyMatcher) LoadKeywords(kwords map[string]string) error {
    if len(kwords) == 0 {
        return errors.New("no keywords")
    }
    self.kws = make([]*KeyWord, 0)
    self.rmap = make(map[string][]*KeyWord)
    total := 0
    for k, cat := range kwords {
        kw := NewKeyWord(k, cat)
        kw.Prepare(self.splitDict)
        self.kws = append(self.kws, kw)

        for _, pattern := range kw.GetPatterns() {
            if _, exist := self.rmap[pattern]; exist {
                self.rmap[pattern] = append(self.rmap[pattern], kw)
            }else {
                self.rmap[pattern] = make([]*KeyWord, 0)
                self.rmap[pattern] = append(self.rmap[pattern], kw)
            }
        }
        total += 1
        if total % 3000 == 0 {
            fmt.Printf("prepared %d keywords...\n", total)
        }
    }
    return nil
}

func (self *FuzzyMatcher) LoadSplitDict(chaiziFile string) error {
    self.splitDict = make(map[string]string)
    f, err := os.Open(chaiziFile)
    if err != nil {
        return err
    }
    defer f.Close()
    rd := bufio.NewReader(f)
    for {
        line, err := rd.ReadString('\n')
        if err != nil || io.EOF == err {
            break
        }
        line = strings.Replace(line, "\n", "", -1)
        parts := strings.Split(line, "\t")
        if len(parts) < 2 {
            break
        }
        self.splitDict[parts[0]] = strings.Replace(parts[1], " ", "", -1)
    }
    return nil
}

func (self FuzzyMatcher) Match(content string, threshold float32) (data map[string]map[string]*FuzzyResult, err error) {
    input := NewInputString(content)
    hits := self.matcher.Match([]byte(input.GetPinyinPattern()))
    results := self.matcher.Index2Pattern(hits)
    fuzzyResults := make(map[string]map[string]*FuzzyResult)
    for result, _ := range results {
        if kws, exist := self.rmap[result]; exist {
            for _, kw := range kws {
                for _, _py := range kw.Pinyins {
                    py := strings.Join(_py, "")
                    if py == result {
                        subs := input.SearchSubStrings(_py)
                        for _, sub := range subs {
                            if len(sub) > 0 {
                                if _, exist := fuzzyResults[kw.Kw]; exist {
                                    if _, exist = fuzzyResults[kw.Kw][sub]; exist {
                                        fuzzyResults[kw.Kw][sub].Count += 1
                                        continue
                                    }
                                }
                                sim := Similarity(sub, kw.Kw)
                                fuzzy := &FuzzyResult{
                                    Kw : kw.Kw,
                                    Cat : kw.Cat,
                                    Sub : sub,
                                    Score : sim,
                                    Count : 1,
                                }
                                if sim > threshold {
                                    if _, exist := fuzzyResults[kw.Kw]; exist {
                                        fuzzyResults[kw.Kw][sub] = fuzzy
                                    }else {
                                        fuzzyResults[kw.Kw] = make(map[string]*FuzzyResult, 0)
                                        fuzzyResults[kw.Kw][sub] = fuzzy
                                    }
                                }
                            }
                        }
                    }
                }
            }
        }
    }
    /*
    for k, v := range fuzzyResults {
        fmt.Printf("关键字 %v\n", k)
        for sub, rr := range v {
            fmt.Printf("匹配子串\"%v\", [%v]\n", sub, rr)
        }
        fmt.Println()
    }
    */
    return fuzzyResults, nil
}
