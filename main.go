package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"strconv"
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

		fmt.Printf("%v %v %v\n", a,b,operation)

		req := CalculatorRequest{
			A:a,
			B:b,
			Operation:operation,
		}

        reqJson, err := json.Marshal(req)

		fmt.Printf("%v \n", reqJson)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		resp, err := http.Post("http://127.0.0.1:8000", "application/json", bytes.NewBuffer(reqJson))

		fmt.Printf("%v \n", resp.Body)
		fmt.Printf("%v \n", resp.StatusCode)

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

		fmt.Println(string(body))

		var respData CalculatorResponse
		err = json.Unmarshal(body, &respData)

		fmt.Println(respData)

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
    http.HandleFunc("/", handler)

    fmt.Println("Server started on Port 9000")
    http.ListenAndServe(":9000", nil)
}