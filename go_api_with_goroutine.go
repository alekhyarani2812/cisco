package api

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

// RequestBody represents the incoming request payload
type RequestBody struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

// ResponseBody represents the API response payload
type ResponseBody struct {
	Message string `json:"message"`
}

func Run() {
	http.HandleFunc("/process", processHandler)
	log.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// processHandler handles the incoming API request
func processHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	var reqBody RequestBody
	// Read the request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Unable to read request body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close() // Close the body after reading to avoid resource leaks

	// Unmarshal the JSON data into the struct
	err = json.Unmarshal(body, &reqBody)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		http.Error(w, "Error unmarshalling JSON:", http.StatusBadRequest)
		return
	}
	
	// Respond to the client immediately
	respBody := ResponseBody{Message: "Processing your request"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(respBody)
	
	// Run background task using a goroutine
	go func(email string) {
		if err := sendEmail(email); err != nil {
			log.Printf("Failed to send email to %s: %v\n", email, err)
		} else {
			log.Printf("Email successfully sent to %s\n", email)
		}
	}(reqBody.Email)
}

// sendEmail simulates sending an email (mock function)
func sendEmail(email string) error {
	// Simulate a delay
	log.Printf("Sending email to %s...\n", email)
	time.Sleep(2 * time.Second) // Simulated delay
	log.Printf("Email sent to %s\n", email)
	return nil
}
