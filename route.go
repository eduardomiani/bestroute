package main

import (
	"sort"
	"strconv"
	"strings"
)

// Route represents a found Route
type Route struct {
	Plan  string
	Price int
}

func findBestRoute(from, to string, limit int) []Route {
	rotas := make([]string, 0)
	rotas = append(rotas, from)
	rotas = find(rotas, to)
	return parseResults(rotas, limit)
}

// find sweeps the file searching all possible routes, given a source route and a destination
func find(rotas []string, f string) []string {
	var dones int
	cop := make([]string, 0, len(rotas))
	for i, v := range rotas {
		if strings.HasPrefix(v, "done-") {
			cop = append(cop, v)
			dones++
			continue
		}
		parts := strings.Split(v, "-")
		var o string
		if len(parts) <= 1 {
			o = parts[len(parts)-1]
		} else {
			o = parts[len(parts)-2]
		}
		for _, rota := range parsedRoutes[o] {
			nova := rotas[i] + "-" + rota[1] + "-" + rota[2]
			if rota[1] == f || strings.Contains(rotas[i], rota[1]) {
				nova = "done-" + nova
			}
			cop = append(cop, nova)
		}
	}
	if dones < len(cop) {
		return find(cop, f)
	} else {
		return cop
	}
}

// parseResults handles the found routes and returns a Route array
func parseResults(routes []string, limit int) []Route {
	results := make([]Route, 0, len(routes))
	for _, r := range routes {
		parts := strings.Split(r, "-")[1:]
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
		route := Route{
			Plan:  routeString,
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
		return len(results[i].Plan) < len(results[j].Plan)
	})
	if limit > len(results) {
		limit = len(results)
	}
	return results[:limit]
}
