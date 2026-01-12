package api

import (
	"fmt"
	"net/http"
	"strconv"
)

func HandleHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func HandleSum(w http.ResponseWriter, r *http.Request) {
	// read x and y from query parameters
	xs := r.URL.Query().Get("x")
	ys := r.URL.Query().Get("y")

	// convert to integers
	x, err1 := strconv.Atoi(xs)
	y, err2 := strconv.Atoi(ys)

	if err1 != nil || err2 != nil {
		http.Error(w, "Invalid parameters", http.StatusBadRequest)
		return
	}

	// compute sum
	sum := x + y

	// return result
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"sum": %d}`, sum)
}
