package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"

	"verified-job-platform/models"
	"verified-job-platform/utils"
)

var db *sql.DB

// SetDB sets the database connection for handlers
func SetDB(database *sql.DB) {
	db = database
}

// registerHandler handles user registration
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var req models.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate input
	if req.Email == "" || req.Password == "" || req.Name == "" {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Email, password, and name are required")
		return
	}

	if req.Role != "jobseeker" && req.Role != "recruiter" {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Role must be 'jobseeker' or 'recruiter'")
		return
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to hash password")
		return
	}

	// Insert user into database
	result, err := db.Exec(`
		INSERT INTO users (email, password, name, role)
		VALUES (?, ?, ?, ?)
	`, req.Email, hashedPassword, req.Name, req.Role)

	if err != nil { // handle unique constraint error
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			utils.WriteErrorResponse(w, http.StatusConflict, "Email already exists")
			return
		}
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to create user")
		return
	}

	userID, _ := result.LastInsertId()

	// Get the created user
	var user models.User
	err = db.QueryRow(`
		SELECT id, email, name, role, created_at
		FROM users WHERE id = ?
	`, userID).Scan(&user.ID, &user.Email, &user.Name, &user.Role, &user.CreatedAt)

	if err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to retrieve user")
		return
	}

	utils.WriteSuccessResponse(w, user)
}

// loginHandler handles user login
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate input
	if req.Email == "" || req.Password == "" {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Email and password are required")
		return
	}

	// Get user from database
	var user models.User
	var hashedPassword string
	err := db.QueryRow(`
		SELECT id, email, password, name, role, created_at
		FROM users WHERE email = ?
	`, req.Email).Scan(&user.ID, &user.Email, &hashedPassword, &user.Name, &user.Role, &user.CreatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			utils.WriteErrorResponse(w, http.StatusUnauthorized, "Invalid email or password")
			return
		}
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Database error")
		return
	}

	// Check password
	if !utils.CheckPasswordHash(req.Password, hashedPassword) {
		utils.WriteErrorResponse(w, http.StatusUnauthorized, "Invalid email or password")
		return
	}

	// Generate token
	token := utils.CreateToken(user.ID)

	response := models.LoginResponse{
		User:  user,
		Token: token,
	}

	utils.WriteSuccessResponse(w, response)
}
