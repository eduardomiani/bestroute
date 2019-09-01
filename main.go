package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/eduardomiani/bestroute/rest"
	"github.com/eduardomiani/bestroute/route"
)

var (
	interfac, port string
)

func init() {
	flag.StringVar(&interfac, "it", "console", "Interface: type 'rest' to run a web app")
	flag.StringVar(&port, "p", "8080", "port")
}

func main() {
	flag.Parse()
	args := os.Args
	if len(args) <= 1 {
		log.Fatal("File is required")
		os.Exit(1)
	}

	filename := args[len(args)-1]
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	defer file.Close()
	_, err = route.ParseFile(file, false)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	route.File = filename

	if interfac == "rest" {
		restInterface()
	} else {
		consoleInterface()
	}
}

// restInterface provides the rest interface
func restInterface() {
	http.HandleFunc("/api/v1/routes", rest.BestRouteHandler)
	log.Printf("Starting server on port %s...\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}

// consoleInterface provides the default interface of this program
func consoleInterface() {
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
		if len(parts) != 2 {
			fmt.Printf("\nInvalid route. Please enter a route in the pattern FROM-TO\n")
			fmt.Printf("\n>>> Please enter the route: ")
			continue
		}
		results := route.FindBestRoute(parts[0], parts[1], 1)
		if len(results) > 0 {
			fmt.Printf("\nBest route: %s > $%d\n", results[0].Route, results[0].Price)
		} else {
			fmt.Printf("\nNo route found. Please try another route\n")
		}
		fmt.Printf("\n>>> Please enter the route: ")
	}
}
