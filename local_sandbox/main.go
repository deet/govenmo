// Package local_sandbox emulates the Venmo sandbox API for posting payments and charges.
// It does does emulate retrieving payments or users.
package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strconv"
	"strings"
)

var decoder = schema.NewDecoder()

type SandboxRequest struct {
	AccessToken string `schema:"access_token"`
	UserId      string `schema:"user_id"`
	Email       string `schema:"email"`
	Phone       string `schema:"phone"`
	Amount      string `schema:"amount"`
	Note        string `schema:"note"`
	Audience    string `schema:"audience"`
}

type VenmoError struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

type ErrorResponse struct {
	Error VenmoError `json:"error"`
}

func main() {
	r := mux.NewRouter()
	//r.HandleFunc("/payments", NewPaymentHandler).Methods("POST")
	r.HandleFunc("/payments", PaymentsIndex).Methods("POST")
	r.HandleFunc("/payments/1111111111111111111", GetPaymentHandler).Methods("GET")
	r.HandleFunc("/me", MeHandler).Methods("GET")

	r.PathPrefix("/").HandlerFunc(ProxyHandler)

	http.Handle("/", r)

	addr := ":4000"
	log.Println("Listening on", addr)
	http.ListenAndServe(addr, nil)
}

func sandboxUser(request *SandboxRequest) bool {
	return request.UserId == "145434160922624933" || strings.ToLower(request.Email) == "venmo@venmo.com" || request.Phone == "15555555555"
}

func ProxyHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Proxying to real sandbox:", r.URL.Path)
	realSandboxUrl, _ := url.Parse("https://sandbox-api.venmo.com/v1/")
	httputil.NewSingleHostReverseProxy(realSandboxUrl).ServeHTTP(w, r)
}

func WriteError(w http.ResponseWriter, status, code int, message string) {
	b, err := json.Marshal(ErrorResponse{VenmoError{Message: message, Code: code}})
	if err != nil {
		http.Error(w, "Sandbox JSON error", 500)
		return
	}
	log.Println("Returning error:", string(b))
	http.Error(w, string(b), status)
}

func PaymentsIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	err := r.ParseForm()

	if err != nil {
		log.Println("Error parsing request:", err)
		http.Error(w, "Error parsing request.", 500)
		return
	}

	var request *SandboxRequest = &SandboxRequest{}
	decoder.IgnoreUnknownKeys(false)
	err = decoder.Decode(request, r.PostForm)
	if err != nil {
		http.Error(w, "Sandbox error. Could not decode request: "+err.Error(), 500)
		return
	}
	log.Printf("Incoming form: %+v \n", r.PostForm)

	if r.FormValue("access_token") != "" {
		request.AccessToken = r.FormValue("access_token")
	}

	if request.AccessToken == "" {
		log.Println("Missing access token.")
		WriteError(w, 401, 261, "You did not pass a valid OAuth access token.")
		return
	}

	sourceFile := ""
	amount, _ := strconv.ParseFloat(request.Amount, 64)
	switch amount {
	case 0.10:
		log.Println("0.10 settled payment")
		sourceFile = "responses/payment/settled.json"
		if !sandboxUser(request) {
			sourceFile = "responses/regular_user_error.json"
			w.WriteHeader(400)
		}
	case 0.20:
		log.Println("0.30 failed payment")
		sourceFile = "responses/payment/failed.json"
		if !sandboxUser(request) {
			sourceFile = "responses/regular_user_error.json"
			w.WriteHeader(400)
		}
	case 0.30:
		log.Println("0.30 pending payment")
		sourceFile = "responses/payment/pending.json"
		if sandboxUser(request) {
			sourceFile = "responses/payment/pending_error.json"
			w.WriteHeader(400)
		}
	case -0.10:
		log.Println("-0.10 settled charge")
		sourceFile = "responses/payment/settled_charge.json"
	case -0.20:
		log.Println("-0.20 pending charge")
		sourceFile = "responses/payment/pending_charge.json"
		if sandboxUser(request) {
			sourceFile = "responses/payment/pending_error.json"
			w.WriteHeader(400)
		}
	default:
		log.Println("Invalid amount:", request.Amount)
		sourceFile = "responses/invalid_amount.json"
		w.WriteHeader(400)
	}

	file, err := os.Open(sourceFile)
	if err != nil {
		http.Error(w, "SANDBOX ERROR", 500)
		log.Println("Sandbox error:", err)
		return
	}
	io.Copy(w, file)
}

func MeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	log.Println("/me request")

	err := r.ParseForm()

	if err != nil {
		log.Println("Error parsing request:", err)
		http.Error(w, "Error parsing request.", 500)
		return
	}

	var request *SandboxRequest = &SandboxRequest{}
	decoder.IgnoreUnknownKeys(false)
	err = decoder.Decode(request, r.PostForm)
	if err != nil {
		http.Error(w, "Sandbox error. Could not decode request: "+err.Error(), 500)
		return
	}
	log.Printf("Incoming form: %+v \n", r.PostForm)

	if r.FormValue("access_token") != "" {
		request.AccessToken = r.FormValue("access_token")
	}

	if request.AccessToken == "" {
		log.Println("Missing access token.")
		WriteError(w, 401, 261, "You did not pass a valid OAuth access token.")
		return
	}

	sourceFile := "responses/users/me.json"

	file, err := os.Open(sourceFile)
	if err != nil {
		http.Error(w, "SANDBOX ERROR", 500)
		log.Println("Sandbox error:", err)
		return
	}
	io.Copy(w, file)
}

func GetPaymentHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	log.Println("GET /payments/1111111111111111111 request")

	err := r.ParseForm()

	if err != nil {
		log.Println("Error parsing request:", err)
		http.Error(w, "Error parsing request.", 500)
		return
	}

	var request *SandboxRequest = &SandboxRequest{}
	decoder.IgnoreUnknownKeys(false)
	err = decoder.Decode(request, r.PostForm)
	if err != nil {
		http.Error(w, "Sandbox error. Could not decode request: "+err.Error(), 500)
		return
	}
	log.Printf("Incoming form: %+v \n", r.PostForm)

	if r.FormValue("access_token") != "" {
		request.AccessToken = r.FormValue("access_token")
	}

	if request.AccessToken == "" {
		log.Println("Missing access token.")
		WriteError(w, 401, 261, "You did not pass a valid OAuth access token.")
		return
	}

	sourceFile := "responses/payment/get.json"

	file, err := os.Open(sourceFile)
	if err != nil {
		http.Error(w, "SANDBOX ERROR", 500)
		log.Println("Sandbox error:", err)
		return
	}
	io.Copy(w, file)
}
