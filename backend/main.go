package main

import (
	"database/sql"
	"log"
	"net/http"

	"verified-job-platform/handlers"

	_ "modernc.org/sqlite"
)

var db *sql.DB

func main() {
	// Initialize database
	initDB()
	defer db.Close() //

	// Set database connection for handlers
	handlers.SetDB(db)

	// Setup router
	r := http.NewServeMux() //

	// Apply CORS middleware to all routes
	r.HandleFunc("/", corsMiddleware(http.HandlerFunc(handleRoutes)))

	log.Println("Server starting on :8080")    //print the server starting message
	log.Fatal(http.ListenAndServe(":8080", r)) //print the error if any occurs while starting the server
}

// CORS middleware allows cross-origin requests from the frontend
func corsMiddleware(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" { //
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	}
}

func handleRoutes(w http.ResponseWriter, r *http.Request) {
	// Simple routing without external router
	switch {
	case r.URL.Path == "/register" && r.Method == "POST": // Register a new user
		handlers.RegisterHandler(w, r)
	case r.URL.Path == "/login" && r.Method == "POST": // Login an existing user
		handlers.LoginHandler(w, r)
	case r.URL.Path == "/jobs" && r.Method == "POST": // Create a new job
		handlers.CreateJobHandler(w, r)
	case r.URL.Path == "/jobs" && r.Method == "GET": // List all jobs
		handlers.ListJobsHandler(w, r)
	case r.URL.Path == "/jobs/apply" && r.Method == "POST": //apply for a job
		handlers.ApplyJobHandler(w, r)
	case r.URL.Path == "/jobs/applications" && r.Method == "GET": // Get job applications for a recruiter
		handlers.GetJobApplicationsHandler(w, r)
	case r.URL.Path == "/applications/update" && r.Method == "PATCH": // Update application status
		handlers.UpdateApplicationHandler(w, r)
	case r.URL.Path == "/jobs/rate" && r.Method == "POST":
		handlers.RateJobHandler(w, r)
	case r.URL.Path == "/jobs/ratings" && r.Method == "GET":
		handlers.GetJobRatingsHandler(w, r)
	case r.URL.Path == "/jobs/disable" && r.Method == "PATCH":
		handlers.DisableJobHandler(w, r)
	case r.URL.Path == "/jobs" && r.Method == "DELETE":
		handlers.DeleteJobHandler(w, r)
	default:
		// Serve static files
		http.FileServer(http.Dir("../frontend")).ServeHTTP(w, r)
	}
}

func initDB() {
	var err error
	db, err = sql.Open("sqlite", "verified_jobs.db")
	if err != nil {
		log.Fatal(err)
	}

	// Create tables
	createTables()
}

func createTables() {
	// Users table
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			email TEXT UNIQUE NOT NULL,
			password TEXT NOT NULL,
			name TEXT NOT NULL,
			role TEXT NOT NULL CHECK(role IN ('jobseeker', 'recruiter')),
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		log.Fatal(err)
	}

	// Jobs table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS jobs (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT NOT NULL,
			description TEXT NOT NULL,
			company TEXT NOT NULL,
			location TEXT NOT NULL,
			salary_min INTEGER,
			salary_max INTEGER,
			recruiter_id INTEGER NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			disabled BOOLEAN NOT NULL DEFAULT 0,
			FOREIGN KEY (recruiter_id) REFERENCES users (id)
		)
	`)
	if err != nil {
		log.Fatal(err)
	}

	// Applications table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS applications (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			job_id INTEGER NOT NULL,
			applicant_id INTEGER NOT NULL,
			status TEXT NOT NULL DEFAULT 'pending' CHECK(status IN ('pending', 'accepted', 'rejected')),
			cover_letter TEXT,
			resume_url TEXT,
			applied_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (job_id) REFERENCES jobs (id),
			FOREIGN KEY (applicant_id) REFERENCES users (id)
		)
	`)
	if err != nil {
		log.Fatal(err)
	}

	// Job Ratings table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS job_ratings (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			job_id INTEGER NOT NULL,
			user_id INTEGER NOT NULL,
			rating TEXT NOT NULL CHECK(rating IN ('genuine', 'scam')),
			UNIQUE(job_id, user_id),
			FOREIGN KEY (job_id) REFERENCES jobs (id),
			FOREIGN KEY (user_id) REFERENCES users (id)
		)
	`)
	if err != nil {
		log.Fatal(err)
	}
}
