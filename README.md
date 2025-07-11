# 🏠 Housing Survey API

A RESTful backend API built with [Golang Fiber](https://docs.gofiber.io/), [GORM](https://gorm.io/index.html), and [PostgreSQL](https://www.postgresql.org/) designed for scalable and role-based data collection with JWT authentication, audit logging, and image metadata support.

---

## 🚀 Features

- ✅ JWT-based authentication (`login`, `logout`)
- ✅ Role-based access control (`SuperAdmin`, `Admin`, `Surveyor`, etc.)
- ✅ Audit logging for all actions
- ✅ Upload + comment handling with GORM + PostgreSQL
- ✅ Docker and non-Docker support
- ✅ Hot reload with [Air](https://github.com/air-verse/air)

---

## 📦 Prerequisites

Install the following:

| Tool           | Link                                                                 |
|----------------|----------------------------------------------------------------------|
| Go 1.21.x      | https://go.dev/dl/                                                   |
| Git            | https://git-scm.com/                                                 |
| PostgreSQL     | https://www.postgresql.org/download/                                 |
| Docker & Compose | https://docs.docker.com/desktop/install/                            |
| Air (optional) | https://github.com/air-verse/air#installation                        |

---

## 📁 Clone the Repository

```bash
git clone https://github.com/your-username/housing-survey-api.git
cd housing-survey-api
```

---

## ⚙️ Environment Setup

Copy the example env file and modify as needed:

```bash
cp .env.example .env
```

Sample `.env` values:

```env
DB_HOST=db
DB_PORT=5432
DB_USER=survey_user
DB_PASS=survey_pass
DB_NAME=survey_db
JWT_SECRET=supersecretkey
APP_PORT=8080
DB_SEED=true
```

---

## 🐳 Run with Docker (Recommended)

### 1. Build & Run Containers

```bash
docker-compose up --build
```

- API is accessible at: [http://localhost:8080](http://localhost:8080)
- On first run, seeds roles, admin users, and sample data if `DB_SEED=true`.

### 2. Access PostgreSQL (inside container)

```bash
docker-compose exec db psql -U survey_user -d survey_db
```

---

## 🛠️ Run Without Docker

### 1. Install Dependencies

Ensure Go 1.21 is being used:

```bash
go version  # should show go1.21.x
```

Install dependencies:

```bash
go mod tidy
```

### 2. Setup PostgreSQL

Manually create the database and user if not using Docker:

```bash
psql -U postgres

CREATE DATABASE survey_db;
CREATE USER survey_user WITH PASSWORD 'survey_pass';
GRANT ALL PRIVILEGES ON DATABASE survey_db TO survey_user;
```

### 3. Run the App

```bash
go run cmd/main.go
```

Or, if using [Air](https://github.com/air-verse/air):

```bash
air
```

---

## 🧪 API Usage

Use tools like [Postman](https://www.postman.com/) or `curl`:

### 🔐 Login

```bash
curl --location 'http://localhost:8080/api/login' --header 'Content-Type: application/json' --data '{
  "email": "superuser@gmail.com",
  "password": "3jutaRUMAH$"
}'
```

Copy the returned `token` and use it in your requests:

```bash
Authorization: Bearer <your_token_here>
```

### 📤 Create Comment (Public, No Auth)

```bash
curl -X POST http://localhost:8080/api/v1/comments -H "Content-Type: application/json" -d '{"content":"My public comment"}'
```

### 🔐 Create Survey (Auth Required)

```bash
curl -X POST http://localhost:8080/api/v1/surveys -H "Authorization: Bearer <your_token>" -H "Content-Type: application/json" -d '{"address":"Jl. Testing No.1","coordinate":"-6.2,106.8","type":"House"}'
```

---

## 🛠 Accessing Database GUI

- Use tools like [pgAdmin](https://www.pgadmin.org/) or [TablePlus](https://tableplus.com/)
- Host: `localhost`
- Port: `5432`
- User: `survey_user`
- Password: `survey_pass`
- Database: `survey_db`

---

## 🔄 Common Commands

### Rebuild Docker:

```bash
docker-compose down -v
docker-compose up --build
```

### Stop Docker:

```bash
docker-compose down
```

### Seed database manually:

```bash
DB_SEED=true go run cmd/main.go
```

---

## 📚 Docs & References

- [Fiber Documentation](https://docs.gofiber.io/)
- [GORM ORM](https://gorm.io/)
- [JWT Auth Guide](https://jwt.io/introduction/)
- [Air Hot Reload](https://github.com/air-verse/air)

---

## 🙋 Troubleshooting

| Issue                                         | Solution                                                                 |
|----------------------------------------------|--------------------------------------------------------------------------|
| `uuid_generate_v4()` does not exist           | Run `CREATE EXTENSION IF NOT EXISTS "uuid-ossp";` inside PostgreSQL.     |
| Hot reload doesn't work                       | Ensure `air` is installed and `.air.toml` is properly configured         |
| Cannot access containerized Postgres          | Use `docker-compose exec db psql ...`                                    |
| Module requires Go 1.23+                      | Downgrade that module or use Go 1.23+ in Docker                          |

---

## 📄 License

MIT

---

## 👨‍💻 Author

Made with ❤️ by [PKP Backend](https://github.com/pkpbackend) — contributions welcome!