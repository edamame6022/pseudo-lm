package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
)

type Dic struct {
	Word []string `json:"word"`
	N    []int    `json:"n"`
}

type DicItem struct {
	Word string `json:"word"`
	N    int    `json:"n"`
}

type Pred struct {
	N  int   `json:"n"`
	F1 []int `json:"f1"`
	F2 []int `json:"f2"`
	F3 []int `json:"f3"`
	W1 []int `json:"w1"`
	W2 []int `json:"w2"`
	W3 []int `json:"w3"`
}

func main() {
	dic := loadDic("./lib/dictionary.json")
	dicmap := make(map[string]int)
	for i, s := range dic.Word {
		dicmap[s] = dic.N[i]
	}

	pred := loadPreds("./lib/prediction.json")

	s, proc := cmd()
	// ヘルプ表示時（proc == nil）は処理を行わずに安全終了
	if proc == nil {
		return
	}

	newDicmap, newPred := proc(dicmap, pred, s)
	newDic := mapToDic(newDicmap)

	if err := saveDic("./lib/dictionary.json", newDic); err != nil {
		panic(fmt.Sprintf("failed to save dictionary: %v", err))
	}
	if err := savePreds("./lib/prediction.json", newPred); err != nil {
		panic(fmt.Sprintf("failed to save prediction: %v", err))
	}
}

// 空ファイル (0バイト) やファイル不在を安全にハンドリングして読み込む
func loadDic(path string) Dic {
	file, err := os.Open(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return Dic{}
		}
		panic(err)
	}
	defer file.Close()

	// saveDic で []DicItem 形式として保存されているため、[]DicItem でデコードする
	var items []DicItem
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&items); err != nil {
		if errors.Is(err, io.EOF) {
			return Dic{} // 0バイトファイルの場合は空のDicを返す
		}
		panic(err)
	}

	words := make([]string, len(items))
	nList := make([]int, len(items))
	for i, item := range items {
		words[i] = item.Word
		nList[i] = item.N
	}

	return Dic{
		Word: words,
		N:    nList,
	}
}

// 空ファイル (0バイト) やファイル不在を安全にハンドリングして読み込む
func loadPreds(path string) []Pred {
	file, err := os.Open(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return []Pred{}
		}
		panic(err)
	}
	defer file.Close()

	var preds []Pred
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&preds); err != nil {
		if errors.Is(err, io.EOF) {
			return []Pred{} // 0バイトファイルの場合は空のスライスを返す
		}
		panic(err)
	}

	return preds
}

func mapToDic(m map[string]int) Dic {
	words := make([]string, len(m))
	nList := make([]int, len(m))

	for word, id := range m {
		if id >= 0 && id < len(m) {
			words[id] = word
			nList[id] = id
		}
	}

	return Dic{
		Word: words,
		N:    nList,
	}
}

func saveDic(path string, dic Dic) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	items := make([]DicItem, len(dic.Word))
	for i := range dic.Word {
		items[i] = DicItem{
			Word: dic.Word[i],
			N:    dic.N[i],
		}
	}

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(items)
}

func savePreds(path string, preds []Pred) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(preds)
}
