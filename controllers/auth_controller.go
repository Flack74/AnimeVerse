package controller

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"animeverse/services"
)

type AuthRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name,omitempty"`
}

type AuthResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Token   string `json:"token,omitempty"`
	User    interface{} `json:"user,omitempty"`
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var req AuthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendAuthResponse(w, http.StatusBadRequest, false, "Invalid request", "", nil)
		return
	}

	if req.Email == "" || req.Password == "" {
		sendAuthResponse(w, http.StatusBadRequest, false, "Email and password required", "", nil)
		return
	}

	// Create user with Supabase-style ID (simulate)
	userID := "user_" + generateID()
	user, err := services.CreateOrUpdateUser(userID, req.Email, req.Name)
	if err != nil {
		sendAuthResponse(w, http.StatusInternalServerError, false, "Failed to create user", "", nil)
		return
	}

	// Generate JWT token
	token, err := generateJWT(userID, req.Email, req.Name)
	if err != nil {
		sendAuthResponse(w, http.StatusInternalServerError, false, "Failed to generate token", "", nil)
		return
	}

	sendAuthResponse(w, http.StatusCreated, true, "User registered successfully", token, user)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req AuthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendAuthResponse(w, http.StatusBadRequest, false, "Invalid request", "", nil)
		return
	}

	if req.Email == "" || req.Password == "" {
		sendAuthResponse(w, http.StatusBadRequest, false, "Email and password required", "", nil)
		return
	}

	// Simple auth check (in production, verify password hash)
	if req.Email == "demo@animeverse.com" && req.Password == "demo123" {
		userID := "user_demo"
		user, err := services.CreateOrUpdateUser(userID, req.Email, "Demo User")
		if err != nil {
			sendAuthResponse(w, http.StatusInternalServerError, false, "Failed to get user", "", nil)
			return
		}

		// Seed demo data if user is new
		if user.CreatedAt.After(user.UpdatedAt.Add(-time.Minute)) {
			services.SeedDemoUserData(userID)
		}

		token, err := generateJWT(userID, req.Email, "Demo User")
		if err != nil {
			sendAuthResponse(w, http.StatusInternalServerError, false, "Failed to generate token", "", nil)
			return
		}

		sendAuthResponse(w, http.StatusOK, true, "Login successful", token, user)
		return
	}

	sendAuthResponse(w, http.StatusUnauthorized, false, "Invalid credentials", "", nil)
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	// For JWT, logout is handled client-side by removing token
	sendAuthResponse(w, http.StatusOK, true, "Logged out successfully", "", nil)
}

func SupabaseOAuthHandler(w http.ResponseWriter, r *http.Request) {
	provider := r.URL.Query().Get("provider")
	if provider == "" {
		provider = "google"
	}

	// Redirect to Supabase OAuth
	supabaseURL := os.Getenv("SUPABASE_URL")
	redirectURL := supabaseURL + "/auth/v1/authorize?provider=" + provider + "&redirect_to=" + r.Header.Get("Referer")
	
	http.Redirect(w, r, redirectURL, http.StatusTemporaryRedirect)
}

func generateJWT(userID, email, name string) (string, error) {
	claims := jwt.MapClaims{
		"sub":   userID,
		"email": email,
		"name":  name,
		"exp":   time.Now().Add(time.Hour * 24 * 7).Unix(), // 7 days
		"iat":   time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("SUPABASE_JWT_SECRET")))
}

func generateID() string {
	return time.Now().Format("20060102150405")
}

func sendAuthResponse(w http.ResponseWriter, statusCode int, success bool, message, token string, user interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	
	response := AuthResponse{
		Success: success,
		Message: message,
		Token:   token,
		User:    user,
	}
	
	json.NewEncoder(w).Encode(response)
}