package main

import (
	"bytes"
	"encoding/json"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type CalculatorRequest struct {
	A      float64 `json:"a"`
	B      float64 `json:"b"`
	Operation  string `json:"operation"`
}

type CalculatorResponse struct {
	A      float64 `json:"a"`
	B      float64 `json:"b"`
	Result float64 `json:"result"`
}

type ErrorResp struct {
	Error string
}

var apiEndPoint string

func ErrorRespWrap(w http.ResponseWriter, error string, code int) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)
	tmpl := template.Must(template.New("error").ParseFiles("view.html"))
	resp := ErrorResp{
		Error:error,
	}
	if err := tmpl.Execute(w, resp); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	return
}

func handler(w http.ResponseWriter, r *http.Request) {
    if r.Method == "GET" {
        var tmpl = template.Must(template.New("form").ParseFiles("view.html"))
        var err = tmpl.Execute(w, nil)

        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
			return
        }
        return
    } else if r.Method == "POST" {
        var tmpl = template.Must(template.New("final").ParseFiles("view.html"))

        if err := r.ParseForm(); err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        a, _:= strconv.ParseFloat(r.FormValue("a"), 64)
        b, _ := strconv.ParseFloat(r.FormValue("b"), 64)
		operation := r.Form.Get("operation")


		req := CalculatorRequest{
			A:a,
			B:b,
			Operation:operation,
		}

        reqJson, err := json.Marshal(req)


		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		resp, err := http.Post(apiEndPoint, "application/json", bytes.NewBuffer(reqJson))

		if resp.StatusCode == 400 {
			ErrorRespWrap(w, "Bad Request", 400)
			return
		}

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var respData CalculatorResponse
		err = json.Unmarshal(body, &respData)


		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

        if err := tmpl.Execute(w, respData); err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
			return
        }
        return
    }

    http.Error(w, "", http.StatusBadRequest)
}


func main() {
	err := godotenv.Load()

	if err != nil {
		log.Println("Error loading .env file")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "9000"
	}

	apiEndPoint = os.Getenv("API_ENDPOINT")
	if apiEndPoint == "" {
		port = "http://localhost:8000"
	}

    http.HandleFunc("/", handler)

    log.Println("Server Started on Port "+port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}