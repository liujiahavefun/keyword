package main

import (
    "container/list"
)

type node struct {
    root bool
    b []byte
    output bool
    index int
    counter int
    child [256]*node
    fails [256]*node
    suffix *node
    fail *node
}

type Matcher struct {
    trie []node
    extent int
    root *node
    counter int
    patterns []string
}


func NewMatcher(patterns [][]byte) *Matcher {
    m := new(Matcher)
    m.buildTrie(patterns)
    return m
}

func NewStringMatcher(patterns []string) *Matcher {
    m := new(Matcher)
    var d [][]byte
    for _, s := range patterns {
        d = append(d, []byte(s))
    }
    m.buildTrie(d)
    m.patterns = patterns
    return m
}

func (m *Matcher) findPath(b []byte) *node {
    n := &m.trie[0]
    for n != nil && len(b) > 0 {
        n = n.child[int(b[0])]
        b = b[1:]
    }
    return n
}

func (m *Matcher) getFreeNode() *node {
    m.extent += 1
    if m.extent == 1 {
        m.root = &m.trie[0]
        m.root.root = true
    }
    return &m.trie[m.extent-1]
}

func (m *Matcher) buildTrie(dictionary [][]byte) {
    max := 1
    for _, blice := range dictionary {
        max += len(blice)
    }
    m.trie = make([]node, max)
    m.getFreeNode()

    for i, blice := range dictionary {
        n := m.root
        var path []byte
        for _, b := range blice {
            path = append(path, b)
            c := n.child[int(b)]

            if c == nil {
                c = m.getFreeNode()
                n.child[int(b)] = c
                c.b = make([]byte, len(path))
                copy(c.b, path)
                if len(path) == 1 {
                    c.fail = m.root
                }
                c.suffix = m.root
            }
            n = c
        }

        n.output = true
        n.index = i
    }

    l := new(list.List)
    l.PushBack(m.root)

    for l.Len() > 0 {
        n := l.Remove(l.Front()).(*node)
        for i := 0; i < 256; i++ {
            c := n.child[i]
            if c != nil {
                l.PushBack(c)
                for j := 1; j < len(c.b); j++ {
                    c.fail = m.findPath(c.b[j:])
                    if c.fail != nil {
                        break
                    }
                }
                if c.fail == nil {
                    c.fail = m.root
                }
                for j := 1; j < len(c.b); j++ {
                    s := m.findPath(c.b[j:])
                    if s != nil && s.output {
                        c.suffix = s
                        break
                    }
                }
            }
        }
    }

    for i := 0; i < m.extent; i++ {
        for c := 0; c < 256; c++ {
            n := &m.trie[i]
            for n.child[c] == nil && !n.root {
                n = n.fail
            }

            m.trie[i].fails[c] = n
        }
    }

    m.trie = m.trie[:m.extent]
}

func (m *Matcher) Match(in []byte) map[int]int {
    hits := make(map[int]int)
    n := m.root
    for _, b := range in {
        c := int(b)
        if !n.root && n.child[c] == nil {
            n = n.fails[c]
        }
        if n.child[c] != nil {
            f := n.child[c]
            n = f
            if f.output {
                hits[f.index] += 1
            }
            for !f.suffix.root {
                f = f.suffix
                if _, exist := hits[f.index]; !exist {
                    hits[f.index] += 1
                }else {
                    break
                }
            }
        }
    }

    return hits
}

func (m *Matcher) Index2Pattern(hits map[int]int) map[string]int {
    ret := make(map[string]int)
    for index, count := range hits {
        if index < len(m.patterns) {
            ret[m.patterns[index]] = count
        }
    }
    return ret
}