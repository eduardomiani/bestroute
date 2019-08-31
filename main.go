package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

var (
	interfac, port string
)

func init() {
	flag.StringVar(&interfac, "it", "console", "Interface: type 'rest' to run a web app")
	flag.StringVar(&interfac, "p", "8080", "port")
}

func main() {
	flag.Parse()
	args := os.Args
	if len(args) <= 1 {
		log.Fatal("File is required")
	}

	file, err := os.Open(args[len(args)-1])
	if err != nil {
		log.Fatal(err)
	}
	_, err = parseFile(file, false)
	if err != nil {
		log.Fatal(err)
	}

	if interfac == "rest" {
		fmt.Println("Not implemented yet")
		os.Exit(1)
	} else {
		scanner := bufio.NewScanner(os.Stdin)
		fmt.Printf("\n>>> Please enter the route: ")
		for scanner.Scan() {
			input := scanner.Text()
			if input == "" {
				fmt.Printf("\n>>> Please enter the route: ")
				continue
			}
			input = strings.ToUpper(input)
			parts := strings.Split(input, "-")
			results := findBestRoute(parts[0], parts[1], 1)
			if len(results) > 0 {
				fmt.Printf("\nBest route: %s > $%d\n", results[0].Route, results[0].Price)
			} else {
				fmt.Printf("\nNo route was found. Please try another route\n")
			}
			fmt.Printf("\n>>> Please enter the route: ")
		}
	}
}
