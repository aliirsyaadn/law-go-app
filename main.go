package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Math struct {
	A      float64 `json:"a"`
	B      float64 `json:"b"`
	Result float64 `json:"result"`
}

type MathOperation struct {
	Math
	Operation string `json:"operation"`
}

type Error struct {
	Err string `json:"err"`
}

func JSONError(w http.ResponseWriter, info string, code int) {
    w.Header().Set("Content-Type", "application/json; charset=utf-8")
    w.Header().Set("X-Content-Type-Options", "nosniff")
    w.WriteHeader(code)
	err := Error{Err:info}
    json.NewEncoder(w).Encode(err)
}

func all(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var mathOp MathOperation
	err := json.NewDecoder(r.Body).Decode(&mathOp)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	if mathOp.Operation == "add" {
		mathOp.Result = mathOp.A + mathOp.B
		log.Printf("Add %v %v\n", mathOp.A, mathOp.B)
	} else if mathOp.Operation == "sub" {
		mathOp.Result = mathOp.A - mathOp.B
		log.Printf("Substract %v %v\n", mathOp.A, mathOp.B)
	} else if mathOp.Operation == "mul" {
		mathOp.Result = mathOp.A * mathOp.B
		log.Printf("Multiply %v %v\n", mathOp.A, mathOp.B)
	} else if mathOp.Operation == "div" {
		if mathOp.B == 0 {
			JSONError(w,"Divide by 0", 400)
			return
		}
		mathOp.Result = mathOp.A / mathOp.B
		log.Printf("Divide %v %v\n", mathOp.A, mathOp.B)
	} else {
		JSONError(w, "Operation not found", 400)
		return
	}
	
	err = json.NewEncoder(w).Encode(mathOp)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}
// Add
func add(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var math Math
	err := json.NewDecoder(r.Body).Decode(&math)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}	
	math.Result = math.A + math.B
	log.Printf("Add %v %v\n", math.A, math.B)

	err = json.NewEncoder(w).Encode(math)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

// Subtract
func sub(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var math Math
	err := json.NewDecoder(r.Body).Decode(&math)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	math.Result = math.A - math.B
	log.Printf("Substract %v %v\n", math.A, math.B)

	err = json.NewEncoder(w).Encode(math)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

// Multiply
func mul(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var math Math
	err := json.NewDecoder(r.Body).Decode(&math)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	math.Result = math.A * math.B
	log.Printf("Multiply %v %v\n", math.A, math.B)

	err = json.NewEncoder(w).Encode(math)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

// Divide
func div(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var math Math
	err := json.NewDecoder(r.Body).Decode(&math)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	if math.B == 0 {
		JSONError(w,"Divide by 0", 400)
		return
	}
	math.Result = math.A / math.B
	log.Printf("Divide %v %v\n", math.A, math.B)

	err = json.NewEncoder(w).Encode(math)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

// Main function
func main() {
	// Init router
	r := mux.NewRouter()
	r.HandleFunc("/", all).Methods("POST")
	r.HandleFunc("/add", add).Methods("POST")
	r.HandleFunc("/sub", sub).Methods("POST")
	r.HandleFunc("/mul", mul).Methods("POST")
	r.HandleFunc("/div", div).Methods("POST")

	// Start server
	log.Println("Server Started on Port 8000")
	log.Fatal(http.ListenAndServe(":8000", r))
}
