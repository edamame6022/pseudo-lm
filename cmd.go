package main

import (
	"flag"
	"fmt"
	"os"
)

func cmd() (string, func(d map[string]int, p []Pred, path string) (map[string]int, []Pred)) {
	// 引数なし、または "help", "-h", "--help" の場合は使用方法を表示
	if len(os.Args) < 2 || os.Args[1] == "help" || os.Args[1] == "-h" || os.Args[1] == "--help" {
		printUsage()
		return "", nil
	}

	// サブコマンドごとに FlagSet を定義
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
		return args[0], learn

	case "gen":
		genCmd.Parse(os.Args[2:])
		args := genCmd.Args()
		if len(args) < 1 {
			fmt.Println("error: The argument is required")
			fmt.Println("usage: mycli gen <arg>")
			os.Exit(1)
		}
		return args[0], gen

	default:
		fmt.Printf("unknown command: %s\n\n", os.Args[1])
		printUsage()
		os.Exit(1)
		return "", nil
	}
}

// 全体の使用方法を表示
func printUsage() {
	fmt.Println("usage:")
	fmt.Println("  mycli <command> [arguments]")
	fmt.Println("\ncommands:")
	fmt.Println("  learn <filepath>   Learn from the file, update the dictionary, and recollect the word predictions.")
	fmt.Println("  gen   <sentence>   Generate the following words using the dictionary and the predictions.")
	fmt.Println("  help               Display this help")
}
