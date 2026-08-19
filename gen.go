package main

import (
	"fmt"
	"math/rand"
	"slices"
	"strings"
)

const (
	wg1 int = 1
	wg2 int = 1
	wg3 int = 1
)

const mapRepitation int = 100

func gen(d map[string]int, p []Pred, s string) (map[string]int, []Pred) {
	if len(p) == 0 || len(d) == 0 {
		fmt.Println("Error: Dictionary or Prediction data is empty.")
		return d, p
	}

	var res []int
	var p1, p2, p3 int

	// 末尾から単語を取得 (w1: 1つ前, w2: 2つ前, w3: 3つ前)
	word1, word2, word3 := getInit(s)

	p1 = getValidID(word1, d, len(p))
	p2 = getValidID(word2, d, len(p))
	p3 = getValidID(word3, d, len(p))

	for rep := 0; rep < mapRepitation; rep++ {
		possibleID := []int{}
		possibleFreq := []int{}

		// 1. 直前単語 (p1) の 1つ後 (W1)
		if p1 >= 0 && p1 < len(p) {
			for i, targetID := range p[p1].W1 {
				addOrUpdateScore(&possibleID, &possibleFreq, targetID, p[p1].F1[i]*wg1)
			}
		}

		// 2. 2つ前単語 (p2) の 2つ後 (W2)
		if p2 >= 0 && p2 < len(p) {
			for i, targetID := range p[p2].W2 {
				addOrUpdateScore(&possibleID, &possibleFreq, targetID, p[p2].F2[i]*wg2)
			}
		}

		// 3. 3つ前単語 (p3) の 3つ後 (W3)
		if p3 >= 0 && p3 < len(p) {
			for i, targetID := range p[p3].W3 {
				addOrUpdateScore(&possibleID, &possibleFreq, targetID, p[p3].F3[i]*wg3)
			}
		}

		// 候補が存在しない場合（行き止まり）はランダムな単語を選びフォールバック
		var apply int
		if len(possibleFreq) == 0 {
			apply = rand.Intn(len(p))
		} else {
			// 最高スコア固定ではなく、確率選択を適用
			apply = sampleWeighted(possibleID, possibleFreq)
		}

		res = append(res, apply)

		// 次のループに向けて状態を1つずつシフト
		p3 = p2
		p2 = p1
		p1 = apply
	}

	idDic := decodeMap(d)
	output := decodeWords(res, idDic)
	fmt.Println(output)

	return d, p
}

// 辞書からIDを取得（存在しなければランダム）
func getValidID(word string, d map[string]int, pLen int) int {
	if word != "" {
		if v, ok := d[word]; ok {
			return v
		}
	}
	return rand.Intn(pLen)
}

// スコアの追加・更新ロジックを関数化
func addOrUpdateScore(ids *[]int, freqs *[]int, targetID, score int) {
	if idx := slices.Index(*ids, targetID); idx != -1 {
		(*freqs)[idx] += score
	} else {
		*ids = append(*ids, targetID)
		*freqs = append(*freqs, score)
	}
}

// nilマップにならぬよう make で作成
func decodeMap(d map[string]int) map[int]string {
	m := make(map[int]string, len(d))
	for word, id := range d {
		if id >= 0 && id < len(d) {
			m[id] = word
		}
	}
	return m
}

// strings.Builder を使いメモリ効率よく結合
func decodeWords(r []int, d map[int]string) string {
	var sb strings.Builder
	for i, id := range r {
		if word, ok := d[id]; ok {
			sb.WriteString(word)
		} else {
			sb.WriteString("?")
		}
		if i < len(r)-1 {
			sb.WriteString(" ")
		}
	}
	return sb.String()
}

func getInit(sentence string) (string, string, string) {
	s := strings.Fields(sentence)
	n := len(s)

	var w1, w2, w3 string
	if n >= 1 {
		w1 = s[n-1]
	}
	if n >= 2 {
		w2 = s[n-2]
	}
	if n >= 3 {
		w3 = s[n-3]
	}
	return w1, w2, w3
}

// スコア（頻度）の高さに応じた確率で単語IDを選択する
func sampleWeighted(ids []int, freqs []int) int {
	// 1. スコアの総和を計算
	total := 0
	for _, f := range freqs {
		total += f
	}

	// すべてのスコアが0以下の場合は均等ランダム
	if total <= 0 {
		return ids[rand.Intn(len(ids))]
	}

	// 2. 0 〜 total-1 の範囲で乱数を生成
	r := rand.Intn(total)

	// 3. 乱数がどの単語の領域に入るか判定
	for i, f := range freqs {
		if r < f {
			return ids[i]
		}
		r -= f
	}

	return ids[len(ids)-1]
}
