package main

import (
	"fmt"
	"math/rand"
	"slices"
	"strings"
)

const (
	wg1 int = 4
	wg2 int = 2
	wg3 int = 1
)

const mapRepitation int = 100

func gen(d map[string]int, p []Pred, s string) (map[string]int, []Pred) {
	var res []int
	var output string
	var p1 int //p3 -> p2 -> p1 -> nextword
	var p2 int
	var p3 int
	var rep int = 0

	word1, word2, word3 := getInit(s)
	if v, ok := d[word1]; ok {
		p1 = v
	} else {
		p1 = rand.Intn(len(p))
	}
	if v, ok := d[word2]; ok {
		p2 = v
	} else {
		p2 = rand.Intn(len(p))
	}
	if v, ok := d[word3]; ok {
		p3 = v
	} else {
		p3 = rand.Intn(len(p))
	}

	for rep < mapRepitation {
		//predict the next word
		possibleID := []int{}
		possibleFreq := []int{}
		for i := range p[p1].W1 {
			if idx := slices.Index(possibleID, p[p1].W1[i]); idx != -1 {
				possibleFreq[idx] += (p[p1].F1[i] * wg1)
			} else {
				possibleID = append(possibleID, p[p1].W1[i])
				possibleFreq = append(possibleFreq, p[p1].F1[i]*wg1)
			}
		}
		for i := range p[p2].W2 {
			if idx := slices.Index(possibleID, p[p2].W2[i]); idx != -1 {
				possibleFreq[idx] += (p[p2].F2[i] * wg2)
			} else {
				possibleID = append(possibleID, p[p2].W2[i])
				possibleFreq = append(possibleFreq, p[p2].F2[i]*wg2)
			}
		}
		for i := range p[p3].W3 {
			if idx := slices.Index(possibleID, p[p3].W3[i]); idx != -1 {
				possibleFreq[idx] += (p[p3].F3[i] * wg3)
			} else {
				possibleID = append(possibleID, p[p3].W3[i])
				possibleFreq = append(possibleFreq, p[p3].F3[i]*wg3)
			}
		}
		apply := possibleID[slices.Index(possibleFreq, slices.Max(possibleFreq))]
		res = append(res, apply)
	}

	idDic := decodeMap(d)

	output = decodeWords(res, idDic)
	fmt.Println(output)

	return d, p
}

func decodeMap(d map[string]int) map[int]string {
	var m map[int]string
	for word, id := range d {
		if id >= 0 && id < len(d) {
			m[id] = word
		}
	}
	return m
}

func decodeWords(r []int, d map[int]string) string {
	output := ""
	for i := range r {
		output += d[r[i]]
		output += " "
	}
	return output
}

func getInit(sentence string) (string, string, string) {
	s := strings.Fields(sentence)
	n := len(s)

	// それぞれの単語を保持する変数（デフォルトは空文字列）
	var w1, w2, w3 string

	// 末尾から順に安全に取り出す
	if n >= 1 {
		w1 = s[n-1] // 1つ前（最後の単語） -> "jumps"
	}
	if n >= 2 {
		w2 = s[n-2] // 2つ前            -> "fox"
	}
	if n >= 3 {
		w3 = s[n-3] // 3つ前            -> "brown"
	}
	return w1, w2, w3
}
