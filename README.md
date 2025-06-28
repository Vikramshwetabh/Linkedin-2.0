# Verified Job Platform

A full-stack job platform where job seekers can register, log in, apply to jobs, and recruiters can post jobs and manage applicants.

## 🚀 Features

### For Job Seekers

- User registration and authentication
- Browse available job listings
- Apply to jobs with cover letter and resume URL
- Track application status

### For Recruiters

- User registration and authentication
- Post new job opportunities
- View and manage applications
- Accept or reject applicants

## 🛠️ Tech Stack

### Frontend

- **HTML5** - Structure
- **Tailwind CSS** - Styling and responsive design
- **Vanilla JavaScript** - Client-side functionality

### Backend

- **Go** - Server-side language
- **SQLite** - File-based database
- **net/http** - HTTP server
- **bcrypt** - Password hashing

## 📁 Project Structure

```
verified-job-platform/
├── backend/
│   ├── main.go              # Server entry point
│   ├── go.mod               # Go module file
│   ├── handlers/
│   │   ├── auth.go          # Authentication handlers
│   │   └── jobs.go          # Job-related handlers
│   ├── models/
│   │   └── models.go        # Data structures
│   └── utils/
│       └── auth.go          # Authentication utilities
├── frontend/
│   ├── index.html           # Landing page
│   ├── login.html           # Login page
│   ├── register.html        # Registration page
│   ├── jobseeker-dashboard.html  # Job seeker dashboard
│   └── recruiter-dashboard.html  # Recruiter dashboard
└── README.md
```

## 🚀 Getting Started

### Prerequisites

- Go 1.21 or higher
- Modern web browser

### Installation

1. **Clone the repository**

   ```bash
   git clone <repository-url>
   cd verified-job-platform
   ```

2. **Install Go dependencies**

   ```bash
   cd backend
   go mod tidy
   ```

3. **Run the server**

   ```bash
   go run main.go
   ```

4. **Access the application**
   - Open your browser and go to `http://localhost:8080`
   - The server will automatically serve the frontend files

## 📡 API Endpoints

### Authentication

- `POST /register` - User registration
- `POST /login` - User login

### Jobs

- `POST /jobs` - Create a new job (recruiters only)
- `GET /jobs` - List all jobs
- `POST /jobs/apply?job_id={id}` - Apply to a job (job seekers only)
- `GET /jobs/applications?job_id={id}` - Get applications for a job (recruiters only)
- `PATCH /applications/update?application_id={id}` - Update application status (recruiters only)

## 🔐 Authentication

The application uses a simple token-based authentication system:

- Tokens are stored in localStorage after login
- All API requests include the token in the Authorization header
- Format: `Bearer {token}`

## 🗄️ Database Schema

### Users Table

```sql
CREATE TABLE users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    email TEXT UNIQUE NOT NULL,
    password TEXT NOT NULL,
    name TEXT NOT NULL,
    role TEXT NOT NULL CHECK(role IN ('jobseeker', 'recruiter')),
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
```

### Jobs Table

```sql
CREATE TABLE jobs (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT NOT NULL,
    description TEXT NOT NULL,
    company TEXT NOT NULL,
    location TEXT NOT NULL,
    salary_min INTEGER,
    salary_max INTEGER,
    recruiter_id INTEGER NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (recruiter_id) REFERENCES users (id)
);
```

### Applications Table

```sql
CREATE TABLE applications (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    job_id INTEGER NOT NULL,
    applicant_id INTEGER NOT NULL,
    status TEXT NOT NULL DEFAULT 'pending' CHECK(status IN ('pending', 'accepted', 'rejected')),
    cover_letter TEXT,
    resume_url TEXT,
    applied_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (job_id) REFERENCES jobs (id),
    FOREIGN KEY (applicant_id) REFERENCES users (id)
);
```

## 🎨 Frontend Pages

### Landing Page (`index.html`)

- Welcome page with platform overview
- Navigation to login and register

### Authentication Pages

- **Login** (`login.html`) - User authentication
- **Register** (`register.html`) - User registration with role selection

### Dashboards

- **Job Seeker Dashboard** (`jobseeker-dashboard.html`)

  - Browse available jobs
  - Apply to jobs with cover letter
  - View application status

- **Recruiter Dashboard** (`recruiter-dashboard.html`)
  - Post new job opportunities
  - View applications for posted jobs
  - Accept or reject applicants

## 🔧 Development

### Adding New Features

1. Update the database schema if needed
2. Add new models in `backend/models/`
3. Create handlers in `backend/handlers/`
4. Update the main.go routing
5. Create frontend pages and JavaScript

### Testing the API

You can test the API endpoints using curl or Postman:

```bash
# Register a new user
curl -X POST http://localhost:8080/register \
  -H "Content-Type: application/json" \
  -d '{"name":"John Doe","email":"john@example.com","password":"password123","role":"jobseeker"}'

# Login
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{"email":"john@example.com","password":"password123"}'

# Create a job (requires recruiter token)
curl -X POST http://localhost:8080/jobs \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer {token}" \
  -d '{"title":"Software Engineer","company":"Tech Corp","location":"San Francisco","description":"We are looking for..."}'
```

## 🚀 Deployment

### Local Development

The application is ready to run locally with the provided setup.

### Production Deployment

For production deployment, consider:

- Using a production database (PostgreSQL, MySQL)
- Implementing proper JWT tokens
- Adding HTTPS
- Setting up proper CORS policies
- Using a reverse proxy (nginx)
- Containerizing with Docker

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Test thoroughly
5. Submit a pull request

## 📝 License

This project is open source and available under the [MIT License](LICENSE).

## 🆘 Support

If you encounter any issues or have questions:

1. Check the documentation
2. Review the code comments
3. Open an issue on GitHub

---

