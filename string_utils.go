package main

import "strings"

type indentlevel struct {
	tabs int
	spaces int
}


func TrimIndent(s string) (string, indentlevel) {
	count := indentlevel{}
	for {
		if strings.HasPrefix(s, "\t") {
			s = strings.TrimPrefix(s, "\t")
			count.tabs += 1
		} else if strings.HasPrefix(s, " ") {
			s = strings.TrimPrefix(s, " ")
			count.spaces += 1
		} else {
			break
		}
	}
	return strings.TrimSpace(s), count
}

func Partition(words []string, split string) ([]string, []string, bool) {
	idx := IndexOf(words, split)
	if idx == -1 {
		return nil,nil,false
	}
	left := words[:idx]
	right := words[idx+1:]
	return left,right,true
}
	
func IndexOf(words []string, w string) int {
	for i := range words {
		if words[i] == w { return i }
	}
	return -1
}
func FirstRune(s string) rune {
	for _, r := range s {
		return r
	}
	return 0
}


