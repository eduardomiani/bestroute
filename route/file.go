package route

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

var File string

// parseFile parses a routes file and returns a map of fromRoute to your respective routes matrix.
// If there is already a parsedRoutes variable in memory, it is returned.
// If the force flag is true, then the file is parsed anyway.
func ParseFile(file io.Reader, force bool) (map[string][][]string, error) {
	if !force && len(parsedRoutes) > 0 {
		log.Println("Using parsedRoutes cache")
		return parsedRoutes, nil
	}
	log.Printf("Loading file...\n")
	parsedRoutes = make(map[string][][]string, 0)
	reader := csv.NewReader(file)
	lines, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("Error parsing routes file: %v", err)
	}
	for i, line := range lines {
		if len(line) != 3 {
			return nil, fmt.Errorf("line %v: Invalid line format", i+1)
		}
		from := strings.TrimSpace(line[0])
		if from == "" {
			return nil, fmt.Errorf("line %v: Invalid src route", i+1)
		}
		to := strings.TrimSpace(line[1])
		if to == "" {
			return nil, fmt.Errorf("line %v: Invalid dst route", i+1)
		}
		priceString := strings.TrimSpace(line[2])
		_, err := strconv.Atoi(priceString)
		if err != nil {
			return nil, fmt.Errorf("line %v: Invalid price", i+1)
		}
		if _, ok := parsedRoutes[from]; !ok {
			parsedRoutes[from] = make([][]string, 0)
		}
		parsedRoutes[from] = append(parsedRoutes[from], []string{strings.ToUpper(from), strings.ToUpper(to), priceString})
	}
	log.Println("File successfully parsed!")
	return parsedRoutes, nil
}

// Add adds a new route into the program
// If this route already exists, it is updated
func Add(r *Route) (bool, error) {
	log.Printf("Adding route %s-%s:%d\n", r.From, r.To, r.Price)
	created, err := addRoute(r)
	if err != nil {
		return false, err
	}
	if err = os.Rename(File+".aux", File); err != nil {
		return false, err
	}

	file, err := os.Open(File)
	if err != nil {
		return false, err
	}
	defer file.Close()
	ParseFile(file, true)
	return created, nil
}

func addRoute(r *Route) (bool, error) {
	file, err := os.OpenFile(File, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return false, err
	}
	defer file.Close()
	reader := csv.NewReader(file)
	lines, err := reader.ReadAll()
	if err != nil {
		return false, err
	}

	var (
		updated bool
		from    = strings.ToUpper(r.From)
		to      = strings.ToUpper(r.To)
		price   = strconv.Itoa(r.Price)
	)
	for i := range lines {
		if strings.ToUpper(lines[i][0]) == from &&
			strings.ToUpper(lines[i][1]) == to {
			lines[i][2] = price
			updated = true
			break
		}
		if updated {
			break
		}
	}
	if !updated {
		lines = append(lines, []string{from, to, price})
	}

	newFile, err := os.Create(File + ".aux")
	if err != nil {
		return false, err
	}
	defer newFile.Close()
	writer := csv.NewWriter(newFile)
	if err = writer.WriteAll(lines); err != nil {
		return false, err
	}
	writer.Flush()
	return !updated, nil
}
