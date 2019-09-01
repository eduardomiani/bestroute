package route

import (
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
	"strings"
)

// parseFile parses a routes file and returns a map of fromRoute to your respective routes matrix.
// If there is already a parsedRoutes variable in memory, it is returned.
// If the force flag is true, then the file is parsed anyway.
func ParseFile(file io.Reader, force bool) (map[string][][]string, error) {
	if !force && len(parsedRoutes) > 0 {
		fmt.Println("Using parsedRoutes cache")
		return parsedRoutes, nil
	}
	fmt.Println("Loading file...")
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
		parsedRoutes[from] = append(parsedRoutes[from], []string{from, to, priceString})
	}
	fmt.Println("File successfully parsed!")
	return parsedRoutes, nil
}
