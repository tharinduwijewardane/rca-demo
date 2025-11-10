package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"
)

var (
	deleteEnabled  bool
	deleteAttempts map[string]time.Time
)

// Request/Response structures
type IntegrationRequest struct {
	UserID   string            `json:"user_id"`
	Token    string            `json:"token"`
	Action   string            `json:"action"`
	Metadata map[string]string `json:"metadata,omitempty"`
}

type IntegrationResponse struct {
	Success     bool                   `json:"success"`
	Message     string                 `json:"message"`
	Data        map[string]interface{} `json:"data,omitempty"`
	ProcessedAt time.Time              `json:"processed_at"`
	RequestID   string                 `json:"request_id"`
}

type AuthResponse struct {
	Valid   bool   `json:"valid"`
	UserID  string `json:"user_id"`
	Message string `json:"message"`
}

type DatabaseResponse struct {
	Success bool                   `json:"success"`
	Data    map[string]interface{} `json:"data"`
}

type NotificationResponse struct {
	Sent    bool   `json:"sent"`
	Message string `json:"message"`
}

func main() {
	// Load environment configuration
	loadConfiguration()

	// Main integration endpoint
	http.HandleFunc("/api/process", handleIntegrationRequest)

	// Simulated external services
	http.HandleFunc("/auth/validate", handleAuthService)
	http.HandleFunc("/database/fetch", handleDatabaseService)
	http.HandleFunc("/notification/send", handleNotificationService)

	// Health check
	http.HandleFunc("/health", handleHealth)

	log.Printf("Integration Service starting on port 9090...")

	if err := http.ListenAndServe(":9090", nil); err != nil {
		log.Fatal(err)
	}
}

// loadConfiguration loads settings from environment variables
func loadConfiguration() {
	// Load user IDs that might have issues
	var err error
	deleteEnabled, err = strconv.ParseBool(os.Getenv("DELETE_ENABLED"))
	if err != nil {
		deleteEnabled = false // Default to false if not set or invalid
	}
}

// Main integration handler that orchestrates the flow
func handleIntegrationRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	requestID := generateRequestID()
	startTime := time.Now()

	log.Printf("[%s] Received integration request", requestID)

	// Parse request
	var req IntegrationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("[%s] Failed to parse request: %v", requestID, err)
		respondWithError(w, "Invalid request body", http.StatusBadRequest, requestID)
		return
	}

	log.Printf("[%s] Request parsed - UserID: %s, Action: %s", requestID, req.UserID, req.Action)

	// Step 1: Validate with Auth Service
	log.Printf("[%s] Step 1: Calling Auth Service...", requestID)
	authValid, err := callAuthService(req.Token, req.UserID)
	if err != nil {
		log.Printf("[%s] Auth service error: %v", requestID, err)
		respondWithError(w, "Auth service unavailable", http.StatusServiceUnavailable, requestID)
		return
	}

	if !authValid {
		log.Printf("[%s] Authentication failed", requestID)
		respondWithError(w, "Authentication failed", http.StatusUnauthorized, requestID)
		return
	}
	log.Printf("[%s] ✓ Auth validated successfully", requestID)

	// Step 2: Fetch data from Database
	log.Printf("[%s] Step 2: Calling Database Service...", requestID)
	dbData, err := callDatabaseService(req.UserID, req.Action)
	if err != nil {
		log.Printf("[%s] Database service error: %v", requestID, err)
		respondWithError(w, "Database service unavailable", http.StatusServiceUnavailable, requestID)
		return
	}
	log.Printf("[%s] ✓ Data fetched successfully", requestID)

	// Step 3: Send notification
	log.Printf("[%s] Step 3: Calling Notification Service...", requestID)
	_, err = callNotificationService(req.UserID, req.Action)
	if err != nil {
		log.Printf("[%s] Notification service error: %v", requestID, err)
		// We'll continue even if notification fails
		log.Printf("[%s] ⚠ Continuing despite notification failure", requestID)
	} else {
		log.Printf("[%s] ✓ Notification sent successfully", requestID)
	}

	// Step 4: Respond to client
	duration := time.Since(startTime)
	log.Printf("[%s] Processing complete in %v", requestID, duration)

	response := IntegrationResponse{
		Success:     true,
		Message:     fmt.Sprintf("Request processed successfully for action: %s", req.Action),
		Data:        dbData,
		ProcessedAt: time.Now(),
		RequestID:   requestID,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// Simulated Auth Service
func handleAuthService(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Simulate processing time
	time.Sleep(time.Duration(50+rand.Intn(100)) * time.Millisecond)

	var req struct {
		Token  string `json:"token"`
		UserID string `json:"user_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Simulate validation - tokens starting with "valid_" are considered valid
	isValid := len(req.Token) > 6 && req.Token[:6] == "valid_"

	response := AuthResponse{
		Valid:   isValid,
		UserID:  req.UserID,
		Message: "Token validated",
	}

	if !isValid {
		response.Message = "Invalid token"
	}

	log.Printf("[AUTH] Validated token for user %s: %v", req.UserID, isValid)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Simulated Database Service
func handleDatabaseService(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Simulate database query time
	time.Sleep(time.Duration(100+rand.Intn(200)) * time.Millisecond)

	userID := r.URL.Query().Get("user_id")
	action := r.URL.Query().Get("action")

	// Simulate fetched data
	data := map[string]interface{}{
		"user_id":         userID,
		"action":          action,
		"records":         rand.Intn(100) + 1,
		"last_access":     time.Now().Add(-24 * time.Hour).Format(time.RFC3339),
		"permissions":     []string{"read", "write", "execute"},
		"quota_remaining": rand.Intn(1000),
	}

	response := DatabaseResponse{
		Success: true,
		Data:    data,
	}

	log.Printf("[DATABASE] Fetched data for user %s, action %s", userID, action)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Simulated Notification Service
func handleNotificationService(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Simulate notification processing time
	time.Sleep(time.Duration(30+rand.Intn(70)) * time.Millisecond)

	var req struct {
		UserID  string `json:"user_id"`
		Message string `json:"message"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	response := NotificationResponse{
		Sent:    true,
		Message: fmt.Sprintf("Notification sent to user %s", req.UserID),
	}

	log.Printf("[NOTIFICATION] Sent notification to user %s", req.UserID)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Health check endpoint
func handleHealth(w http.ResponseWriter, r *http.Request) {
	health := map[string]interface{}{
		"status":    "healthy",
		"timestamp": time.Now().Format(time.RFC3339),
		"services": map[string]string{
			"auth":         "operational",
			"database":     "operational",
			"notification": "operational",
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(health)
}

// Helper function to call the auth service
func callAuthService(token, userID string) (bool, error) {
	client := &http.Client{Timeout: 5 * time.Second}

	reqBody := map[string]string{
		"token":   token,
		"user_id": userID,
	}

	jsonData, _ := json.Marshal(reqBody)
	resp, err := client.Post("http://localhost:9090/auth/validate", "application/json",
		bytes.NewReader(jsonData))

	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	var authResp AuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
		return false, err
	}

	return authResp.Valid, nil
}

// Helper function to call the database service
func callDatabaseService(userID, action string) (map[string]interface{}, error) {
	client := &http.Client{Timeout: 5 * time.Second}

	if !deleteEnabled && action == "delete_user" {
		return nil, fmt.Errorf("deletion is not enabled")
	}

	url := fmt.Sprintf("http://localhost:9090/database/fetch?user_id=%s&action=%s", userID, action)
	resp, err := client.Get(url)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var dbResp DatabaseResponse
	if err := json.NewDecoder(resp.Body).Decode(&dbResp); err != nil {
		return nil, err
	}

	return dbResp.Data, nil
}

// Helper function to call the notification service
func callNotificationService(userID, action string) (bool, error) {
	client := &http.Client{Timeout: 5 * time.Second}

	reqBody := map[string]string{
		"user_id": userID,
		"message": fmt.Sprintf("Action '%s' completed successfully", action),
	}

	jsonData, _ := json.Marshal(reqBody)
	resp, err := client.Post("http://localhost:9090/notification/send", "application/json",
		bytes.NewReader(jsonData))

	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	var notifResp NotificationResponse
	if err := json.NewDecoder(resp.Body).Decode(&notifResp); err != nil {
		return false, err
	}

	return notifResp.Sent, nil
}

// Helper function to respond with error
func respondWithError(w http.ResponseWriter, message string, statusCode int, requestID string) {
	response := IntegrationResponse{
		Success:     false,
		Message:     message,
		ProcessedAt: time.Now(),
		RequestID:   requestID,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}

// Helper function to generate request ID
func generateRequestID() string {
	return fmt.Sprintf("REQ-%d-%d", time.Now().Unix(), rand.Intn(10000))
}
