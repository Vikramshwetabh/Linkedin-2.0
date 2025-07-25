package handlers

import (
	"encoding/json" //use for encoding json
	"net/http"
	"strconv"

	"verified-job-platform/models"
	"verified-job-platform/utils"
)

// create JobHandler handles job creation
func CreateJobHandler(w http.ResponseWriter, r *http.Request) {
	// Get user ID from token (simplified authentication)
	userID, err := getUserIDFromRequest(r)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusUnauthorized, "Authentication required")
		return
	}

	// Check if user is a recruiter
	var userRole string
	err = db.QueryRow("SELECT role FROM users WHERE id = ?", userID).Scan(&userRole)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "User not found")
		return
	}

	if userRole != "recruiter" { // Check if user is a recruiter
		// If not a recruiter, return forbidden
		utils.WriteErrorResponse(w, http.StatusForbidden, "Only recruiters can create jobs")
		return
	}

	var req models.CreateJobRequest // Decode request body into CreateJobRequest struct
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate input
	if req.Title == "" || req.Description == "" || req.Company == "" || req.Location == "" {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Title, description, company, and location are required")
		return
	}

	// Insert job into database
	result, err := db.Exec(`
		INSERT INTO jobs (title, description, company, location, salary_min, salary_max, recruiter_id, disabled)
		VALUES (?, ?, ?, ?, ?, ?, ?, 0)
	`, req.Title, req.Description, req.Company, req.Location, req.SalaryMin, req.SalaryMax, userID)

	if err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to create job")
		return
	}

	jobID, _ := result.LastInsertId() // Get the last inserted job ID

	// Get the created job
	var job models.Job
	err = db.QueryRow(`
		SELECT id, title, description, company, location, salary_min, salary_max, recruiter_id, created_at, disabled
		FROM jobs WHERE id = ?
	`, jobID).Scan(&job.ID, &job.Title, &job.Description, &job.Company, &job.Location, &job.SalaryMin, &job.SalaryMax, &job.RecruiterID, &job.CreatedAt, &job.Disabled)

	if err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to retrieve job")
		return
	}

	utils.WriteSuccessResponse(w, job)
}

// listJobsHandler handles job listing
func ListJobsHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query(`
		SELECT id, title, description, company, location, salary_min, salary_max, recruiter_id, created_at, disabled
		FROM jobs
		ORDER BY created_at DESC
	`)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to fetch jobs")
		return
	}
	defer rows.Close()

	var jobs []models.Job
	for rows.Next() {
		var job models.Job
		err := rows.Scan(&job.ID, &job.Title, &job.Description, &job.Company, &job.Location, &job.SalaryMin, &job.SalaryMax, &job.RecruiterID, &job.CreatedAt, &job.Disabled)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to scan job")
			return
		}
		jobs = append(jobs, job)
	}

	utils.WriteSuccessResponse(w, jobs)
}

// applyJobHandler handles job applications
func ApplyJobHandler(w http.ResponseWriter, r *http.Request) {
	// Get user ID from token
	userID, err := getUserIDFromRequest(r)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusUnauthorized, "Authentication required")
		return
	}

	// Check if user is a jobseeker
	var userRole string
	err = db.QueryRow("SELECT role FROM users WHERE id = ?", userID).Scan(&userRole)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "User not found")
		return
	}

	if userRole != "jobseeker" {
		utils.WriteErrorResponse(w, http.StatusForbidden, "Only job seekers can apply to jobs")
		return
	}

	// Get job ID from query parameter
	jobIDStr := r.URL.Query().Get("job_id")
	if jobIDStr == "" {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Job ID is required")
		return
	}

	jobID, err := strconv.Atoi(jobIDStr)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid job ID")
		return
	}

	// Check if job exists
	var jobExists bool
	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM jobs WHERE id = ?)", jobID).Scan(&jobExists)
	if err != nil || !jobExists {
		utils.WriteErrorResponse(w, http.StatusNotFound, "Job not found")
		return
	}

	// Check if already applied
	var alreadyApplied bool
	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM applications WHERE job_id = ? AND applicant_id = ?)", jobID, userID).Scan(&alreadyApplied)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Database error")
		return
	}

	if alreadyApplied {
		utils.WriteErrorResponse(w, http.StatusConflict, "Already applied to this job")
		return
	}

	var req models.ApplyJobRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Insert application into database
	result, err := db.Exec(`
		INSERT INTO applications (job_id, applicant_id, cover_letter, resume_url)
		VALUES (?, ?, ?, ?)
	`, jobID, userID, req.CoverLetter, req.ResumeURL)

	if err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to submit application")
		return
	}

	applicationID, _ := result.LastInsertId()

	// Get the created application
	var application models.Application
	err = db.QueryRow(`
		SELECT id, job_id, applicant_id, status, cover_letter, resume_url, applied_at
		FROM applications WHERE id = ?
	`, applicationID).Scan(&application.ID, &application.JobID, &application.ApplicantID, &application.Status, &application.CoverLetter, &application.ResumeURL, &application.AppliedAt)

	if err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to retrieve application")
		return
	}

	utils.WriteSuccessResponse(w, application)
}

// getJobApplicationsHandler handles fetching applications for a job
func GetJobApplicationsHandler(w http.ResponseWriter, r *http.Request) {
	// Get user ID from token
	userID, err := getUserIDFromRequest(r)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusUnauthorized, "Authentication required")
		return
	}

	// Get job ID from query parameter
	jobIDStr := r.URL.Query().Get("job_id")
	if jobIDStr == "" {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Job ID is required")
		return
	}

	jobID, err := strconv.Atoi(jobIDStr)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid job ID")
		return
	}

	// Check if user is the recruiter who posted the job
	var recruiterID int
	err = db.QueryRow("SELECT recruiter_id FROM jobs WHERE id = ?", jobID).Scan(&recruiterID)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusNotFound, "Job not found")
		return
	}

	if recruiterID != userID {
		utils.WriteErrorResponse(w, http.StatusForbidden, "Only the job poster can view applications")
		return
	}

	// Get applications with applicant details
	rows, err := db.Query(`
		SELECT a.id, a.job_id, a.applicant_id, a.status, a.cover_letter, a.resume_url, a.applied_at,
		       u.name, u.email
		FROM applications a
		JOIN users u ON a.applicant_id = u.id
		WHERE a.job_id = ?
		ORDER BY a.applied_at DESC
	`, jobID)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to fetch applications")
		return
	}
	defer rows.Close()

	var applications []models.Application
	for rows.Next() {
		var app models.Application
		var applicantName, applicantEmail string
		err := rows.Scan(&app.ID, &app.JobID, &app.ApplicantID, &app.Status, &app.CoverLetter, &app.ResumeURL, &app.AppliedAt, &applicantName, &applicantEmail)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to scan application")
			return
		}
		app.Applicant = &models.User{
			ID:    app.ApplicantID,
			Name:  applicantName,
			Email: applicantEmail,
		}
		applications = append(applications, app)
	}

	utils.WriteSuccessResponse(w, applications)
}

// updateApplicationHandler handles updating application status
func UpdateApplicationHandler(w http.ResponseWriter, r *http.Request) {
	// Get user ID from token
	userID, err := getUserIDFromRequest(r)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusUnauthorized, "Authentication required")
		return
	}

	// Get application ID from query parameter
	applicationIDStr := r.URL.Query().Get("application_id")
	if applicationIDStr == "" {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Application ID is required")
		return
	}

	applicationID, err := strconv.Atoi(applicationIDStr)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid application ID")
		return
	}

	// Check if user is the recruiter who posted the job
	var jobID int
	err = db.QueryRow("SELECT job_id FROM applications WHERE id = ?", applicationID).Scan(&jobID)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusNotFound, "Application not found")
		return
	}

	var recruiterID int
	err = db.QueryRow("SELECT recruiter_id FROM jobs WHERE id = ?", jobID).Scan(&recruiterID)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusNotFound, "Job not found")
		return
	}

	if recruiterID != userID {
		utils.WriteErrorResponse(w, http.StatusForbidden, "Only the job poster can update applications")
		return
	}

	var req models.UpdateApplicationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate status
	if req.Status != "accepted" && req.Status != "rejected" {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Status must be 'accepted' or 'rejected'")
		return
	}

	// Update application status
	_, err = db.Exec("UPDATE applications SET status = ? WHERE id = ?", req.Status, applicationID)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to update application")
		return
	}

	utils.WriteSuccessResponse(w, map[string]string{"message": "Application status updated successfully"})
}

// RateJobHandler allows a user to rate a job as 'genuine' or 'scam'
func RateJobHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserIDFromRequest(r)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusUnauthorized, "Authentication required")
		return
	}

	var req struct {
		JobID  int    `json:"job_id"`
		Rating string `json:"rating"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	if req.Rating != "genuine" && req.Rating != "scam" {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Rating must be 'genuine' or 'scam'")
		return
	}
	// Upsert rating
	_, err = db.Exec(`INSERT INTO job_ratings (job_id, user_id, rating) VALUES (?, ?, ?)
		ON CONFLICT(job_id, user_id) DO UPDATE SET rating=excluded.rating`, req.JobID, userID, req.Rating)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to rate job")
		return
	}
	utils.WriteSuccessResponse(w, map[string]string{"message": "Rating submitted"})
}

// GetJobRatingsHandler returns the count of 'genuine' and 'scam' ratings for a job
func GetJobRatingsHandler(w http.ResponseWriter, r *http.Request) {
	jobIDStr := r.URL.Query().Get("job_id")
	if jobIDStr == "" {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Job ID is required")
		return
	}
	jobID, err := strconv.Atoi(jobIDStr)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid job ID")
		return
	}
	var genuineCount, scamCount int
	db.QueryRow(`SELECT COUNT(*) FROM job_ratings WHERE job_id = ? AND rating = 'genuine'`, jobID).Scan(&genuineCount)
	db.QueryRow(`SELECT COUNT(*) FROM job_ratings WHERE job_id = ? AND rating = 'scam'`, jobID).Scan(&scamCount)
	utils.WriteSuccessResponse(w, map[string]int{"genuine": genuineCount, "scam": scamCount})
}

// DisableJobHandler toggles the disabled state of a job
func DisableJobHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserIDFromRequest(r)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusUnauthorized, "Authentication required")
		return
	}
	jobIDStr := r.URL.Query().Get("id")
	if jobIDStr == "" {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Job ID is required")
		return
	}
	jobID, err := strconv.Atoi(jobIDStr)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid job ID")
		return
	}
	// Check if the user is the recruiter for this job
	var recruiterID int
	err = db.QueryRow("SELECT recruiter_id FROM jobs WHERE id = ?", jobID).Scan(&recruiterID)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusNotFound, "Job not found")
		return
	}
	if recruiterID != userID {
		utils.WriteErrorResponse(w, http.StatusForbidden, "You are not authorized to disable this job")
		return
	}
	// Toggle the disabled state
	_, err = db.Exec("UPDATE jobs SET disabled = NOT disabled WHERE id = ?", jobID)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to update job status")
		return
	}
	utils.WriteSuccessResponse(w, map[string]interface{}{"job_id": jobID, "disabled_toggled": true})
}

// DeleteJobHandler deletes a job
func DeleteJobHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserIDFromRequest(r)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusUnauthorized, "Authentication required")
		return
	}
	jobIDStr := r.URL.Query().Get("id")
	if jobIDStr == "" {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Job ID is required")
		return
	}
	jobID, err := strconv.Atoi(jobIDStr)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid job ID")
		return
	}
	// Check if the user is the recruiter for this job
	var recruiterID int
	err = db.QueryRow("SELECT recruiter_id FROM jobs WHERE id = ?", jobID).Scan(&recruiterID)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusNotFound, "Job not found")
		return
	}
	if recruiterID != userID {
		utils.WriteErrorResponse(w, http.StatusForbidden, "You are not authorized to delete this job")
		return
	}
	// Delete the job
	_, err = db.Exec("DELETE FROM jobs WHERE id = ?", jobID)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to delete job")
		return
	}
	utils.WriteSuccessResponse(w, map[string]interface{}{"job_id": jobID, "deleted": true})
}

// Helper function to get user ID from request (simplified)
func getUserIDFromRequest(r *http.Request) (int, error) {
	// In a real application, you would extract the token from the Authorization header
	// and validate it properly. This is a simplified version.
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		_, err := utils.GetUserIDFromToken("")
		return 0, err
	}

	// Remove "Bearer " prefix if present
	token := authHeader
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		token = authHeader[7:]
	}

	return utils.GetUserIDFromToken(token)
}
