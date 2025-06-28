package models

import "time"

// User represents a user in the system
type User struct {
	ID        int       `json:"id"`
	Email     string    `json:"email"`
	Password  string    `json:"-"` // Don't include password in JSON
	Name      string    `json:"name"`
	Role      string    `json:"role"` // "jobseeker" or "recruiter"
	CreatedAt time.Time `json:"created_at"`
}

// Job represents a job posting
type Job struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Company     string    `json:"company"`
	Location    string    `json:"location"`
	SalaryMin   *int      `json:"salary_min,omitempty"`
	SalaryMax   *int      `json:"salary_max,omitempty"`
	RecruiterID int       `json:"recruiter_id"`
	CreatedAt   time.Time `json:"created_at"`
}

// Application represents a job application
type Application struct {
	ID          int       `json:"id"`
	JobID       int       `json:"job_id"`
	ApplicantID int       `json:"applicant_id"`
	Status      string    `json:"status"` // "pending", "accepted", "rejected"
	CoverLetter string    `json:"cover_letter,omitempty"`
	ResumeURL   string    `json:"resume_url,omitempty"`
	AppliedAt   time.Time `json:"applied_at"`
	Applicant   *User     `json:"applicant,omitempty"`
	Job         *Job      `json:"job,omitempty"`
}

// RegisterRequest represents the registration request
type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
	Role     string `json:"role"`
}

// LoginRequest represents the login request
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginResponse represents the login response
type LoginResponse struct {
	User  User   `json:"user"`
	Token string `json:"token"`
}

// CreateJobRequest represents the job creation request
type CreateJobRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Company     string `json:"company"`
	Location    string `json:"location"`
	SalaryMin   *int   `json:"salary_min,omitempty"`
	SalaryMax   *int   `json:"salary_max,omitempty"`
}

// ApplyJobRequest represents the job application request
type ApplyJobRequest struct {
	CoverLetter string `json:"cover_letter,omitempty"`
	ResumeURL   string `json:"resume_url,omitempty"`
}

// UpdateApplicationRequest represents the application status update request
type UpdateApplicationRequest struct {
	Status string `json:"status"` // "accepted" or "rejected"
}

// APIResponse represents a generic API response
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}
