package rest

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/eduardomiani/bestroute/route"
)

// ErrorResp represents the structure returned when any error occurs
type ErrorResp struct {
	Error  string `json:"error"`
	Detail string `json:"detail,omitempty"`
}

// handlerResponse generic response function to any request
func handlerResponse(objResp interface{}, status int, w http.ResponseWriter) {
	log.Printf("Response status: %v\n%#v\n", status, objResp)
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(objResp)
}

// BestRouteHandler main handler of this application
func BestRouteHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	switch r.Method {
	case "GET":
		findBestRouteHandler(w, r)
	case "POST":
		createNewRouteHandler(w, r)
	default:
		handlerResponse(
			ErrorResp{Error: "Method not allowed"},
			http.StatusMethodNotAllowed,
			w)
	}
}

// findBestRouteHandler handler executed for method GET
// Tries find the best route, given a src and dst
func findBestRouteHandler(w http.ResponseWriter, r *http.Request) {
	from := r.FormValue("from")
	to := r.FormValue("to")
	limit := 1
	limit, err := strconv.Atoi(r.FormValue("limit"))
	if err != nil {
		limit = 1
	}

	if from == "" {
		handlerResponse(
			ErrorResp{Error: "'From' param is required"},
			http.StatusBadRequest,
			w,
		)
		return
	}
	if to == "" {
		handlerResponse(
			ErrorResp{Error: "'To' param is required"},
			http.StatusBadRequest,
			w,
		)
		return
	}

	routes := route.FindBestRoute(from, to, limit)
	handlerResponse(
		&routes,
		http.StatusOK,
		w,
	)
}

// createNewRouteHandler creates a new route into the routes file
// If the route already exists, it is updated
func createNewRouteHandler(w http.ResponseWriter, r *http.Request) {
	var rt route.Route
	if err := json.NewDecoder(r.Body).Decode(&rt); err != nil {
		handlerResponse(
			ErrorResp{Error: "Invalid payload", Detail: err.Error()},
			http.StatusBadRequest,
			w,
		)
		return
	}

	if err := validateRoute(&rt); err != nil {
		handlerResponse(
			ErrorResp{Error: "Invalid payload", Detail: err.Error()},
			http.StatusBadRequest,
			w,
		)
		return
	}

	created, err := route.Add(&rt)
	if err != nil {
		handlerResponse(
			ErrorResp{Error: err.Error()},
			http.StatusInternalServerError,
			w,
		)
		return
	}

	respStatus := http.StatusCreated
	if !created {
		respStatus = http.StatusOK
	}
	log.Printf("Response status: %v\n", respStatus)
	w.WriteHeader(respStatus)
}

// validateRoute validates a route
func validateRoute(r *route.Route) error {
	if r.From == "" {
		return fmt.Errorf("'From' is required")
	}
	if r.To == "" {
		return fmt.Errorf("'To' is required")
	}
	if r.Price <= 0 {
		return fmt.Errorf("Invalid price")
	}
	return nil
}
