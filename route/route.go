package route

import (
	"log"
	"sort"
	"strconv"
	"strings"
)

var parsedRoutes map[string][][]string

// Route represents a single route into the program
type Route struct {
	From  string `json:"from"`
	To    string `json:"to"`
	Price int    `json:"price"`
}

// RouteResp represents a found Route
type RouteResp struct {
	Route string `json:"route"`
	Price int    `json:"price"`
}

func FindBestRoute(from, to string, limit int) []RouteResp {
	from = strings.ToUpper(from)
	to = strings.ToUpper(to)
	log.Printf("Searching for route %s-%s...\n", from, to)
	rotas := make([]string, 0)
	rotas = append(rotas, from)
	rotas = find(rotas, to)
	return parseResults(rotas, limit)
}

// find sweeps the file searching all possible routes, given a source route and a destination
func find(routes []string, to string) []string {
	var done int
	routesAux := make([]string, 0, len(routes))
	for i, r := range routes {
		parts := strings.Split(r, "-")
		var from string
		if len(parts) <= 1 {
			from = parts[len(parts)-1]
		} else {
			from = parts[len(parts)-2]
			if from == to {
				routesAux = append(routesAux, r)
				done++
				continue
			}
		}
		for _, route := range parsedRoutes[from] {
			if strings.Contains(routes[i], route[1]) {
				continue
			}
			newRoute := routes[i] + "-" + route[1] + "-" + route[2]
			routesAux = append(routesAux, newRoute)
		}
	}
	if done < len(routesAux) {
		return find(routesAux, to)
	} else {
		return routesAux
	}
}

// parseResults handles the found routes and returns a Route array
func parseResults(routes []string, limit int) []RouteResp {
	results := make([]RouteResp, 0, len(routes))
	for _, r := range routes {
		parts := strings.Split(r, "-")
		var (
			routeString string
			total       int
		)
		for i, v := range parts {
			va, err := strconv.Atoi(v)
			if err == nil {
				total += va
			} else {
				routeString += v
				if i < (len(parts) - 2) {
					routeString += " - "
				}
			}
		}
		route := RouteResp{
			Route: routeString,
			Price: total,
		}
		results = append(results, route)
	}

	// Order results from best price and shortest route
	sort.Slice(results, func(i, j int) bool {
		if results[i].Price < results[j].Price {
			return true
		}
		if results[i].Price > results[j].Price {
			return false
		}
		return len(results[i].Route) < len(results[j].Route)
	})
	if limit > len(results) {
		limit = len(results)
	}
	return results[:limit]
}
