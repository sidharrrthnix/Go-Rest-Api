# Go HTTP API - School Management System

A RESTful HTTP API built with Go for managing school teachers and students. Features comprehensive CRUD operations, filtering, sorting, and batch operations.

## 🚀 Features

- **RESTful API** with full CRUD operations
- **Teacher Management**: Create, read, update, delete teachers
- **Student Management**: Basic student endpoints (framework ready)
- **Advanced Filtering & Sorting**: Query teachers by various fields
- **Batch Operations**: Update/delete multiple records at once
- **MariaDB Integration**: Persistent data storage
- **Docker Setup**: Easy development environment
- **Adminer UI**: Web-based database management
- **TLS Support**: HTTPS with self-signed certificates
- **Middleware Stack**: CORS, rate limiting, compression, security headers

## 📋 API Endpoints

### Teachers

#### Single Teacher Operations
- `GET /teachers/{id}` - Get teacher by ID
- `PUT /teachers/{id}` - Update teacher (full update)
- `PATCH /teachers/{id}` - Partial update teacher
- `DELETE /teachers/{id}` - Delete teacher

#### Batch Operations
- `GET /teachers` - Get all teachers (with optional filters & sorting)
- `POST /teachers` - Create multiple teachers
- `PATCH /teachers` - Batch partial update teachers
- `DELETE /teachers` - Batch delete teachers

### Students (Framework Ready)
- `GET /students` - Get all students
- `POST /students` - Create students
- `GET /students/{id}` - Get student by ID
- `PUT /students/{id}` - Update student
- `DELETE /students/{id}` - Delete student

## 🛠️ Quick Start

### Prerequisites
- Docker & Docker Compose
- Go 1.25+ (for local development)
- Git

### 1. Clone & Setup

```bash
git clone <your-repo-url>
cd http-api
```

### 2. Start Database

```bash
docker compose up -d
```

This starts:
- MariaDB on `localhost:3306`
- Adminer on `http://localhost:8081`

### 3. Run the API

```bash
go run .
```

The API will be available at `https://localhost:8080`

### 4. Access Adminer

Open `http://localhost:8081` in your browser:
- **Server**: `mariadb`
- **Username**: `api_user`
- **Password**: `api_password`
- **Database**: `SCHOOL_DB`

## 📚 API Usage Examples

### Get All Teachers
```bash
curl -k https://localhost:8080/teachers
```

### Filter Teachers
```bash
# By subject
curl -k "https://localhost:8080/teachers?subject=Mathematics"

# By class
curl -k "https://localhost:8080/teachers?class=10A"

# Multiple filters
curl -k "https://localhost:8080/teachers?subject=Science&class=10B"
```

### Sort Teachers
```bash
# Sort by first name ascending
curl -k "https://localhost:8080/teachers?sortby=first_name:asc"

# Sort by last name descending
curl -k "https://localhost:8080/teachers?sortby=last_name:desc"

# Multiple sorting
curl -k "https://localhost:8080/teachers?sortby=class:asc&sortby=last_name:asc"
```

### Create Teachers
```bash
curl -k -X POST https://localhost:8080/teachers \
  -H "Content-Type: application/json" \
  -d '[
    {
      "firstName": "John",
      "lastName": "Doe",
      "email": "john.doe@school.com",
      "class": "10A",
      "subject": "Mathematics"
    },
    {
      "firstName": "Jane",
      "lastName": "Smith",
      "email": "jane.smith@school.com",
      "class": "10B",
      "subject": "Science"
    }
  ]'
```

### Update Single Teacher (PUT - Full Update)
```bash
curl -k -X PUT https://localhost:8080/teachers/1 \
  -H "Content-Type: application/json" \
  -d '{
    "firstName": "Updated",
    "lastName": "Name",
    "email": "updated.email@school.com",
    "class": "11A",
    "subject": "Advanced Math"
  }'
```

### Partial Update Single Teacher (PATCH)
```bash
curl -k -X PATCH https://localhost:8080/teachers/1 \
  -H "Content-Type: application/json" \
  -d '{
    "email": "newemail@school.com",
    "subject": "Calculus"
  }'
```

### Batch Update Teachers (PATCH)
```bash
curl -k -X PATCH https://localhost:8080/teachers \
  -H "Content-Type: application/json" \
  -d '[
    {
      "id": 1,
      "email": "teacher1.updated@school.com",
      "class": "12A"
    },
    {
      "id": 2,
      "subject": "Advanced Physics",
      "class": "12B"
    }
  ]'
```

### Delete Single Teacher
```bash
curl -k -X DELETE https://localhost:8080/teachers/1
```

### Batch Delete Teachers
```bash
curl -k -X DELETE https://localhost:8080/teachers \
  -H "Content-Type: application/json" \
  -d '[2, 3, 4]'
```

## 🏗️ Project Structure

```
http-api/
├── docker-compose.yaml      # Database and Adminer setup
├── go.mod                   # Go module definition
├── go.sum                   # Go dependencies
├── server.go                # Main server entry point
├── sample-env              # Environment variables template
├── handlers/               # HTTP request handlers
│   ├── root.go            # Root endpoint handler
│   ├── teacher.go         # Teacher CRUD operations
│   └── student.go         # Student CRUD operations (framework)
├── models/                 # Data models
│   ├── teacher.go         # Teacher struct
│   ├── student.go         # Student struct
│   └── user.go            # User struct (framework)
├── middlewares/            # HTTP middlewares
│   ├── cors.go            # CORS handling
│   ├── rate_limiter.go    # Rate limiting
│   ├── compression.go     # Response compression
│   ├── security_headers.go # Security headers
│   └── response_time.go   # Response time logging
├── respository/            # Data access layer
│   └── sqlconnect/        # Database connection and CRUD
│       ├── sqlconfig.go   # Database configuration
│       └── teachers_crud.go # Teacher database operations
├── router/                 # Route definitions
│   └── router.go          # HTTP router setup
├── utils/                  # Utility functions
│   ├── config.go          # Configuration loading
│   ├── env.go             # Environment variable helpers
│   ├── error_handler.go   # Error handling utilities
│   └── applyMiddleware.go # Middleware application
├── init/                   # Database initialization
│   └── 01-init.sql        # Database schema and seed data
└── scripts/                # Utility scripts
    └── check-db.sh        # Database health check
```

## ⚙️ Configuration

### Environment Variables

Copy `sample-env` to `.env` and adjust values as needed:

```bash
cp sample-env .env
```

| Variable | Default | Description |
|----------|---------|-------------|
| `DB_HOST` | `localhost` | Database host |
| `DB_PORT` | `3306` | Database port |
| `DB_USER` | `api_user` | Database username |
| `DB_PASSWORD` | `api_password` | Database password |
| `DB_NAME` | `SCHOOL_DB` | Database name |
| `DB_MAX_OPEN_CONNS` | `10` | Max open connections |
| `DB_MAX_IDLE_CONNS` | `5` | Max idle connections |
| `DB_CONN_MAX_LIFETIME` | `3m` | Connection max lifetime |
| `SERVER_PORT` | `8080` | Server port |
| `TLS_CERT_FILE` | `cert.pem` | TLS certificate file |
| `TLS_KEY_FILE` | `key.pem` | TLS private key file |

### TLS Certificates

For development, you can generate self-signed certificates:

```bash
# Generate private key
openssl genpkey -algorithm RSA -out key.pem -pkcs8

# Generate certificate
openssl req -new -x509 -key key.pem -out cert.pem -days 365 -subj "/C=US/ST=State/L=City/O=Org/CN=localhost"
```

## 🧪 Testing

### Database Health Check

```bash
./scripts/check-db.sh
```

### API Testing

Use the provided curl commands above or import the Postman collection (create one with the examples above).

## 📦 Dependencies

- **Go 1.25+**
- **MariaDB 11.4+**
- **Docker & Docker Compose**

### Go Modules

- `github.com/go-sql-driver/mysql` - MySQL/MariaDB driver
- `github.com/joho/godotenv` - Environment variable loading

## 🔒 Security Features

- **TLS/HTTPS** encryption
- **CORS** protection
- **Rate limiting** (10 requests/minute)
- **Security headers** (CSP, HSTS, etc.)
- **Input validation**
- **SQL injection protection** (prepared statements)

## 🐛 Troubleshooting

### Database Connection Issues

1. Ensure Docker containers are running: `docker compose ps`
2. Check database logs: `docker compose logs mariadb`
3. Verify environment variables in `.env` file

### TLS Certificate Issues

For development, browsers will show security warnings. Either:
- Click "Advanced" → "Proceed to localhost (unsafe)" in Chrome
- Add the certificate to your system's trusted certificates
- Use `curl -k` to skip verification

### Common Errors

- **"Method not allowed"**: Check HTTP method (GET, POST, PUT, PATCH, DELETE)
- **"Invalid ID"**: Ensure ID is a valid integer in URL path
- **"Teacher not found"**: Verify the teacher ID exists in database
- **"Invalid request body"**: Check JSON syntax and required fields

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch: `git checkout -b feature-name`
3. Make your changes
4. Test thoroughly
5. Submit a pull request

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.

## 🔗 Links

- [Go Documentation](https://golang.org/doc/)
- [MariaDB Documentation](https://mariadb.com/kb/en/documentation/)
- [Docker Compose](https://docs.docker.com/compose/)
- [Adminer](https://www.adminer.org/)
