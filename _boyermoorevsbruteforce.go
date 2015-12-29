package repono

// http://play.golang.org/p/4_Zg_mjvc0

import (
	"bytes"
	"fmt"
)

var (
	p = []byte(`ABC ABCDAB ABCDABCDABDE`)
	b = []byte(`ABCDABD`)
)

func main() {
	n := bruteForceFind(p, b)
	fmt.Println(n)

	bf := makeByteFinder(b)
	n = bf.next(p)
	fmt.Println(n)
}

func bruteForceFind(p, b []byte) int {
	for i := 0; i < len(p); i++ {
		for j := 0; j < len(b); j++ {
			if b[j] != p[i+j] {
				break
			} else if j == len(b)-1 {
				return i
			}
		}
	}
	return -1
}

type byteFinder struct {
	pattern        []byte
	badCharSkip    [256]int
	goodSuffixSkip []int
}

func makeByteFinder(pattern []byte) *byteFinder {
	f := &byteFinder{
		pattern:        pattern,
		goodSuffixSkip: make([]int, len(pattern)),
	}
	last := len(pattern) - 1
	for i := range f.badCharSkip {
		f.badCharSkip[i] = len(pattern)
	}
	for i := 0; i < last; i++ {
		f.badCharSkip[pattern[i]] = last - i
	}
	lastPrefix := last
	for i := last; i >= 0; i-- {
		if bytes.HasPrefix(pattern, pattern[i+1:]) {
			lastPrefix = i + 1
		}
		f.goodSuffixSkip[i] = lastPrefix + last - i
	}
	for i := 0; i < last; i++ {
		lenSuffix := longestCommonSuffix(pattern, pattern[1:i+1])
		if pattern[i-lenSuffix] != pattern[last-lenSuffix] {
			f.goodSuffixSkip[last-lenSuffix] = lenSuffix + last - i
		}
	}
	return f
}

func longestCommonSuffix(a, b []byte) (i int) {
	for ; i < len(a) && i < len(b); i++ {
		if a[len(a)-1-i] != b[len(b)-1-i] {
			break
		}
	}
	return
}

func (f *byteFinder) next(text []byte) int {
	i := len(f.pattern) - 1
	for i < len(text) {
		j := len(f.pattern) - 1
		for j >= 0 && text[i] == f.pattern[j] {
			i--
			j--
		}
		if j < 0 {
			return i + 1
		}
		i += max(f.badCharSkip[text[i]], f.goodSuffixSkip[j])
	}
	return -1
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
