package main

import (
    //"regexp"
    "testing"
    //"sort"
    //"fmt"
)

func assert(t *testing.T, b bool) {
    if !b {
        t.Fail()
    }
}

func TestNoPatterns(t *testing.T) {
    m := NewStringMatcher([]string{})
    hits := m.Match([]byte("老 东西 哈哈"))
    assert(t, len(hits) == 0)
}

func TestNoData(t *testing.T) {
    m := NewStringMatcher([]string{"老", "东西", "哈哈"})
    hits := m.Match([]byte(""))
    assert(t, len(hits) == 0)
}

func TestSuffixes(t *testing.T) {
    m := NewStringMatcher([]string{"牛逼", "真牛逼", "国米真牛逼"})
    hits := m.Match([]byte("国米真牛逼"))

    desired := []string{"牛逼", "真牛逼", "国米真牛逼"}
    assert(t, len(hits) == len(desired))
    patterns := m.Index2Pattern(hits)
    for _, pattern := range desired {
        if pattern, exist := patterns[pattern]; !exist {
            t.Error("match error, pattern: ", pattern, " not found")
        }
    }

    /*
    patterns := m.Index2Pattern(hits)
    sort.Sort(sort.StringSlice(patterns))
    sort.Sort(sort.StringSlice(desired))
    for i, pattern := range patterns {
        if pattern != desired[i] {
            t.Error("match error, pattern:", pattern, " desired of index ", i, " is ", desired[i])
        }
    }
    */

}

/*
func TestPrefixes(t *testing.T) {
    m := NewStringMatcher([]string{"国米", "国米真", "国米真牛"})
    hits := m.Match([]byte("国米真牛逼"))

    desired := []string{"国米", "国米真", "国米真牛"}
    assert(t, len(hits) == len(desired))
    patterns := m.Index2Pattern(hits)
    sort.Sort(sort.StringSlice(patterns))
    sort.Sort(sort.StringSlice(desired))
    for i, pattern := range patterns {
        if pattern != desired[i] {
            t.Error("match error, pattern:", pattern, " desired of index ", i, " is ", desired[i])
        }
    }
}

func TestInterior(t *testing.T) {
    m := NewStringMatcher([]string{"刘哈哈", "哈哈", "哈哈一", "哈哈一笑"})
    hits := m.Match([]byte("老刘哈哈一笑，这个老东西"))

    desired := []string{"刘哈哈", "哈哈", "哈哈一", "哈哈一笑"}
    assert(t, len(hits) == len(desired))
    patterns := m.Index2Pattern(hits)
    sort.Sort(sort.StringSlice(patterns))
    sort.Sort(sort.StringSlice(desired))
    for i, pattern := range patterns {
        if pattern != desired[i] {
            t.Error("match error, pattern:", pattern, " desired of index ", i, " is ", desired[i])
        }
    }
}

func TestMatchAtStart(t *testing.T) {
    m := NewStringMatcher([]string{"老刘哈哈", "老", "老刘"})
    hits := m.Match([]byte("老刘哈哈一笑，这个老东西"))

    desired := []string{"老刘哈哈", "老", "老刘"}
    assert(t, len(hits) == len(desired))
    patterns := m.Index2Pattern(hits)
    sort.Sort(sort.StringSlice(patterns))
    sort.Sort(sort.StringSlice(desired))
    for i, pattern := range patterns {
        if pattern != desired[i] {
            t.Error("match error, pattern:", pattern, " desired of index ", i, " is ", desired[i])
        }
    }
}

func TestMatchAtEnd(t *testing.T) {
    m := NewStringMatcher([]string{"东西", "老东西"})
    hits := m.Match([]byte("老刘哈哈一笑，这个老东西"))

    desired := []string{"东西", "老东西"}
    assert(t, len(hits) == len(desired))
    patterns := m.Index2Pattern(hits)
    sort.Sort(sort.StringSlice(patterns))
    sort.Sort(sort.StringSlice(desired))
    for i, pattern := range patterns {
        if pattern != desired[i] {
            t.Error("match error, pattern:", pattern, " desired of index ", i, " is ", desired[i])
        }
    }
}

func TestOverlappingPatterns(t *testing.T) {
    m := NewStringMatcher([]string{"哈哈 ", "哈 一", " 一笑"})
    hits := m.Match([]byte("老刘哈哈 一笑，这个老东西"))

    desired := []string{"哈哈 ", "哈 一", " 一笑"}
    assert(t, len(hits) == len(desired))
    patterns := m.Index2Pattern(hits)
    sort.Sort(sort.StringSlice(patterns))
    sort.Sort(sort.StringSlice(desired))
    for i, pattern := range patterns {
        if pattern != desired[i] {
            t.Error("match error, pattern:", pattern, " desired of index ", i, " is ", desired[i])
        }
    }
}

func TestMultipleMatches(t *testing.T) {
    m := NewStringMatcher([]string{"中国", "贸易体制", "人类命运共同体", "高级知识分子"})
    hits := m.Match([]byte("习近平强调，中国致力于促进更高水平对外开放，坚定支持多边贸易体制，将在更广领域扩大外资市场准入，积极打造一流营商环境。中国愿同各国深化服务贸易投资合作，促进贸易和投资自由化便利化，推动经济全球化朝着更加开放、包容、普惠、平衡、共赢的方向发展。本届交易会的主题是“开放、创新、智慧、融合”。希望各位代表和嘉宾深入交流，凝聚共识，加强合作，共同促进全球服务贸易繁荣发展，引领世界经济发展方向，造福各国人民，推动构建人类命运共同体。"))

    desired := []string{"中国", "贸易体制", "人类命运共同体"}
    assert(t, len(hits) == len(desired))
    patterns := m.Index2Pattern(hits)
    sort.Sort(sort.StringSlice(patterns))
    sort.Sort(sort.StringSlice(desired))
    for i, pattern := range patterns {
        if pattern != desired[i] {
            t.Error("match error, pattern:", pattern, " desired of index ", i, " is ", desired[i])
        }
    }
}

func TestSingleCharacterMatches(t *testing.T) {
    m := NewStringMatcher([]string{"老", "刘", "老"})
    hits := m.Match([]byte("老刘哈哈 一笑，这个老东西"))

    desired := []string{"老", "刘"}
    assert(t, len(hits) == len(desired))
    patterns := m.Index2Pattern(hits)
    sort.Sort(sort.StringSlice(patterns))
    sort.Sort(sort.StringSlice(desired))
    for i, pattern := range patterns {
        if pattern != desired[i] {
            t.Error("match error, pattern:", pattern, " desired of index ", i, " is ", desired[i])
        }
    }
}

func TestNothingMatches(t *testing.T) {
    m := NewStringMatcher([]string{"fuck", "哦", "世界和平"})
    hits := m.Match([]byte("老刘哈哈 一笑，这个老东西"))
    assert(t, len(hits) == 0)
}
*/