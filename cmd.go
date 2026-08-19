package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

// コマンド実行モード
type Mode int

const (
	ModeSingle      Mode = iota // 従来の単発実行 (mycli gen "hello")
	ModeInteractive             // 対話型常駐モード (mycli interactive)
	ModeHelp                    // ヘルプ表示
)

func cmd() (Mode, string, func(d map[string]int, p []Pred, path string) (map[string]int, []Pred)) {
	// 引数なし、または "interactive" の場合はインタラクティブモードへ
	if len(os.Args) < 2 {
		return ModeInteractive, "", nil
	}

	if os.Args[1] == "interactive" || os.Args[1] == "repl" {
		return ModeInteractive, "", nil
	}

	if os.Args[1] == "help" || os.Args[1] == "-h" || os.Args[1] == "--help" {
		printUsage()
		return ModeHelp, "", nil
	}

	learnCmd := flag.NewFlagSet("learn", flag.ExitOnError)
	genCmd := flag.NewFlagSet("gen", flag.ExitOnError)

	switch os.Args[1] {
	case "learn":
		learnCmd.Parse(os.Args[2:])
		args := learnCmd.Args()
		if len(args) < 1 {
			fmt.Println("error: The argument is required")
			fmt.Println("usage: mycli learn <arg>")
			os.Exit(1)
		}
		return ModeSingle, args[0], learn

	case "gen":
		genCmd.Parse(os.Args[2:])
		args := genCmd.Args()
		if len(args) < 1 {
			fmt.Println("error: The argument is required")
			fmt.Println("usage: mycli gen <arg>")
			os.Exit(1)
		}
		return ModeSingle, strings.Join(args, " "), gen

	default:
		fmt.Printf("unknown command: %s\n\n", os.Args[1])
		printUsage()
		os.Exit(1)
		return ModeHelp, "", nil
	}
}

// インタラクティブモード内で入力された1行を分解して解析する
func parseInteractiveInput(line string) (string, string) {
	parts := strings.Fields(line)
	if len(parts) == 0 {
		return "", ""
	}
	command := parts[0]
	arg := ""
	if len(parts) > 1 {
		arg = strings.Join(parts[1:], " ")
	}
	return command, arg
}

func printUsage() {
	fmt.Println("usage:")
	fmt.Println("  mycli <command> [arguments]")
	fmt.Println("  mycli                       (Starts in interactive mode)")
	fmt.Println("\ncommands:")
	fmt.Println("  interactive        Start standby interactive mode")
	fmt.Println("  learn <filepath>   Learn from the file and update the data")
	fmt.Println("  gen   <sentence>   Generate words starting with sentence")
	fmt.Println("  help               Display this help")
}
