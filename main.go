package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Todo struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
}

type EchoHandler struct{}
type SleepHandler struct{}

type ErrorHandler struct {
	code    int
	message string
}

type Data struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (h *SleepHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	var data Data
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&data)
	if err != nil {
		code := http.StatusBadRequest
		h_resp(w, code, Todo{ID: code, Title: fmt.Sprintf("err=%v", err)})
		return
	}

	fmt.Printf("sleeping %d seconds ... zzz\n", data.Code)
	time.Sleep(time.Duration(data.Code) * time.Second)
	fmt.Printf("wake up  %d seconds\n", data.Code)
	fmt.Fprint(w, "I slept in %d seconds\n", data.Code)
}

func (h *EchoHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	var data Data
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&data)
	if err != nil {
		code := http.StatusBadRequest
		h_resp(w, code, Todo{ID: code, Title: fmt.Sprintf("err=%v", err)})
		return
	}


	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "code=%d, message=%s\n", data.Code, data.Message)
}

// curl --include  -X POST -H "Content-Type: application/json" -d '{"code": 123, "message": "hogefuga"}'   http://localhost:8888/500
func (h *ErrorHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	todos := []Todo{}

	var data Data
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&data)
	if err != nil {
		code := http.StatusBadRequest
		h_resp(w, code, Todo{ID: code, Title: fmt.Sprintf("err=%v", err)})
		return
	}
	// fmt.Printf("data=%v\n", data)

	todos = append(todos, Todo{
		ID:    h.code,
		Title: fmt.Sprintf("method: " + r.Method)})
	todos = append(todos, Todo{ID: data.Code, Title: data.Message})
	h_resp(w, data.Code, todos)
}

func h_resp(w http.ResponseWriter, code int, data any) {
	response, _ := json.Marshal(data)

	w.WriteHeader(code)
	w.Write(response)
}

func main() {

	server := http.Server{
		Addr:    ":8888",
		Handler: nil,
		ReadHeaderTimeout: 20 * time.Second,
	}

	http.Handle("/sleep", &SleepHandler{})
	http.Handle("/echo", &EchoHandler{})

	http.Handle("/500", &ErrorHandler{code: 500, message: "internal SERVER error"})

	err := server.ListenAndServe()
	if err != nil {
		fmt.Println(err)
	}
}
