package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
)

func learn(d map[string]int, p []Pred, path string) (map[string]int, []Pred) {
	first := true
	second := true
	third := true
	var p1, p2, p3, p4 int

	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanWords)

	for scanner.Scan() {
		word := scanner.Text()

		/*
			if word == "split word" {
				continue
			}
		*/

		// 単語登録・番号保持
		if v, ok := d[word]; ok {
			p1 = v
		} else {
			next := len(d)
			d[word] = next
			p = append(p, Pred{N: next})
			p1 = next
		}

		if first {
			first = false
			p2 = p1
			continue
		}

		// p2 (1つ前) の 1つ後 -> p1
		if i := slices.Index(p[p2].W1, p1); i != -1 {
			p[p2].F1[i] += 1
		} else {
			p[p2].W1 = append(p[p2].W1, p1)
			p[p2].F1 = append(p[p2].F1, 1)
		}

		if second {
			second = false
			p3 = p2
			p2 = p1
			continue
		}

		// p3 (2つ前) の 1つ後 -> p2
		if i := slices.Index(p[p3].W1, p2); i != -1 {
			p[p3].F1[i] += 1
		} else {
			p[p3].W1 = append(p[p3].W1, p2)
			p[p3].F1 = append(p[p3].F1, 1)
		}

		// p3 (2つ前) の 2つ後 -> p1
		if i := slices.Index(p[p3].W2, p1); i != -1 {
			p[p3].F2[i] += 1
		} else {
			p[p3].W2 = append(p[p3].W2, p1)
			p[p3].F2 = append(p[p3].F2, 1)
		}

		if third {
			third = false
			p4 = p3
			p3 = p2
			p2 = p1
			continue
		}

		// p4 (3つ前) の 1つ後 -> p3
		if i := slices.Index(p[p4].W1, p3); i != -1 {
			p[p4].F1[i] += 1
		} else {
			p[p4].W1 = append(p[p4].W1, p3)
			p[p4].F1 = append(p[p4].F1, 1)
		}

		// p4 (3つ前) の 2つ後 -> p2
		if i := slices.Index(p[p4].W2, p2); i != -1 {
			p[p4].F2[i] += 1
		} else {
			p[p4].W2 = append(p[p4].W2, p2)
			p[p4].F2 = append(p[p4].F2, 1)
		}

		// p4 (3つ前) の 3つ後 -> p1
		if i := slices.Index(p[p4].W3, p1); i != -1 {
			p[p4].F3[i] += 1
		} else {
			p[p4].W3 = append(p[p4].W3, p1)
			p[p4].F3 = append(p[p4].F3, 1)
		}

		p4 = p3
		p3 = p2
		p2 = p1
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "file read error: %v\n", err)
	}

	return d, p
}
