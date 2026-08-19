package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
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

const (
	dicPath  = "./lib/dictionary.json"
	predPath = "./lib/prediction.json"
)

func main() {
	mode, arg, proc := cmd()
	if mode == ModeHelp {
		return
	}

	// 1. モデル読み込み（初回1回のみで高速化）
	fmt.Print("Loading model into memory... ")
	dic := loadDic(dicPath)
	dicmap := make(map[string]int)
	for i, s := range dic.Word {
		dicmap[s] = dic.N[i]
	}
	pred := loadPreds(predPath)
	fmt.Printf("Done! (Dictionary: %d words, Predictions: %d entries)\n", len(dicmap), len(pred))

	// 2. 単発実行モード（従来の mycli gen "..." 等）
	if mode == ModeSingle {
		newDicmap, newPred := proc(dicmap, pred, arg)
		saveModel(newDicmap, newPred)
		return
	}

	// 3. インタラクティブ（待機）モード
	runInteractiveMode(dicmap, pred)
}

// 対話型ループ処理
func runInteractiveMode(dicmap map[string]int, pred []Pred) {
	fmt.Println("\n--- Interactive Standby Mode ---")
	fmt.Println("Commands:")
	fmt.Println("  gen <sentence>   - Generate text")
	fmt.Println("  learn <filepath> - Learn from file")
	fmt.Println("  exit / quit      - Save and Exit")
	fmt.Println("--------------------------------")

	scanner := bufio.NewScanner(os.Stdin)
	isModified := false

	for {
		fmt.Print("\nmycli> ")
		if !scanner.Scan() {
			break
		}

		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		command, inputArg := parseInteractiveInput(line)

		switch command {
		case "exit", "quit", "q":
			if isModified {
				fmt.Print("Saving updated model... ")
				saveModel(dicmap, pred)
				fmt.Println("Done!")
			}
			fmt.Println("Bye!")
			return

		case "gen":
			if inputArg == "" {
				fmt.Println("error: sentence is required. Example: gen my name is")
				continue
			}
			dicmap, pred = gen(dicmap, pred, inputArg)

		case "learn":
			if inputArg == "" {
				fmt.Println("error: filepath is required. Example: learn ./data.txt")
				continue
			}
			dicmap, pred = learn(dicmap, pred, inputArg)
			isModified = true
			fmt.Println("Successfully learned and updated in-memory model.")

		case "help":
			fmt.Println("Available commands: gen <sentence>, learn <filepath>, exit")

		default:
			fmt.Printf("Unknown command '%s'. Type 'gen', 'learn', or 'exit'.\n", command)
		}
	}

	// ループ終了後の読み込みエラーチェックを追加
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "read error: %v\n", err)
	}
}

// モデル保存の共通関数
func saveModel(dicmap map[string]int, pred []Pred) {
	newDic := mapToDic(dicmap)
	if err := saveDic(dicPath, newDic); err != nil {
		panic(fmt.Sprintf("failed to save dictionary: %v", err))
	}
	if err := savePreds(predPath, pred); err != nil {
		panic(fmt.Sprintf("failed to save prediction: %v", err))
	}
}

func loadDic(path string) Dic {
	file, err := os.Open(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return Dic{}
		}
		panic(err)
	}
	defer file.Close()

	var items []DicItem
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&items); err != nil {
		if errors.Is(err, io.EOF) {
			return Dic{}
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
			return []Pred{}
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
