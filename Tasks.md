# Tasks — Roadmap Belajar Golang REST API + JWT

Checklist bertahap dari instalasi Go hingga production. Tandai [x] saat selesai.

---

## [x] Fase 1 — Setup Environment Go

### Install Go
- [x] Install Go 1.22+ via Homebrew: `brew install go`
- [x] Verifikasi instalasi: `go version`
- [x] Pahami GOPATH vs Go Modules (perbedaan dengan Java classpath)
- [x] Tambahkan Go binary ke PATH di ~/.zshrc:
      `export PATH=$PATH:$(go env GOPATH)/bin`

### Tools
- [x] Install air (live reload, seperti Spring DevTools): `go install github.com/air-verse/air@latest`
- [ ] Install sqlc atau migrate (opsional, untuk SQL migrations)
- [x] Setup VS Code dengan extension: Go (by Google)

---

## [x] Fase 2 — Dasar-dasar Go (Wajib sebelum mulai project)

### Sintaks Dasar
- [x] Pahami deklarasi variabel: `var` vs `:=`
- [x] Pahami tipe data dasar: `string`, `int`, `bool`, `float64`
- [x] Pahami `struct` — ekuivalen `class` di Java
- [x] Pahami `interface` di Go — berbeda dengan Java (implicit, bukan explicit)
- [x] Pahami pointer: `*` dan `&` — tidak ada di Java
- [x] Pahami error handling di Go: `if err != nil` (tidak ada try-catch)

### Ekuivalen Java → Go
| Java | Go |
|---|---|
| `class User {}` | `type User struct {}` |
| `interface UserService` | `type UserService interface{}` |
| `try { } catch (Exception e) {}` | `result, err := doSomething(); if err != nil {}` |
| `List<User>` | `[]User` (slice) |
| `Map<String, String>` | `map[string]string` |
| `null` | `nil` |
| `Optional<User>` | `*User` (pointer, bisa nil) |

### Package & Module
- [x] Pahami struktur `go.mod` — ekuivalen `pom.xml`
- [x] Pahami `package main` dan `package internal`
- [x] Pahami cara import package: `import "github.com/gofiber/fiber/v2"`
- [x] Install dependencies: `go get <package>` (ekuivalen `mvn dependency:resolve`)

---

## [x] Fase 3 — Init Project Go Sentinel

### Setup Project
- [x] Inisialisasi Go module: `go mod init github.com/alfaizuna/go-sentinel`
- [x] Buat struktur folder sesuai PRD.md
- [x] Install semua dependency:
  ```bash
  go get github.com/gofiber/fiber/v2
  go get github.com/golang-jwt/jwt/v5
  go get gorm.io/gorm
  go get gorm.io/driver/postgres
  go get github.com/joho/godotenv
  go get golang.org/x/crypto/bcrypt
  ```

### Config & Database
- [x] Buat file `.env` dengan variabel:
  ```
  DB_HOST=localhost
  DB_PORT=5432
  DB_NAME=gosentinel_db
  DB_USER=postgres
  DB_PASSWORD=
  JWT_SECRET=404E635266556A586E3272357538782F
  JWT_EXPIRATION_HOURS=24
  APP_PORT=8081
  ```
- [x] Buat `.env.example` (copy dari `.env` tanpa nilai sensitif)
- [x] Buat `internal/config/config.go` — load env ke struct Config
- [x] Buat database: `createdb gosentinel_db`
- [x] Koneksikan GORM ke PostgreSQL

---

## [x] Fase 4 — Model & Repository

### Model
- [x] Buat `internal/model/user.go`:
  ```go
  type Role string
  const (
      RoleUser  Role = "USER"
      RoleAdmin Role = "ADMIN"
  )
  type User struct {
      gorm.Model          // ID, CreatedAt, UpdatedAt, DeletedAt
      Name     string     `gorm:"not null"`
      Email    string     `gorm:"uniqueIndex;not null"`
      Password string     `gorm:"not null"`
      Role     Role       `gorm:"default:USER"`
  }
  ```
- [x] Jalankan GORM AutoMigrate untuk membuat tabel `users`
- [x] Pahami perbedaan: `gorm.Model` di Go vs `@Entity` + Hibernate di Java

### Repository
- [x] Buat `internal/repository/user_repository.go`:
  - [x] `FindByEmail(email string) (*model.User, error)`
  - [x] `ExistsByEmail(email string) bool`
  - [x] `Save(user *model.User) error`
- [x] Pahami: GORM di Go lebih eksplisit dibanding Spring Data JPA

---

## [x] Fase 5 — JWT Service

- [x] Buat `internal/service/jwt_service.go`:
  - [x] `GenerateToken(email string) (string, error)` — buat JWT token baru
  - [x] `ValidateToken(tokenString string) (*Claims, error)` — validasi & ekstrak claims
- [x] Pahami struct `Claims` yang embed `jwt.RegisteredClaims`
- [x] Pahami perbedaan: di Go tidak ada `UserDetails` interface, cukup string email
- [x] Test manual generate & parse token

---

## [x] Fase 6 — Auth Service & Handler

### Auth Service
- [x] Buat `internal/service/auth_service.go`:
  - [x] `Register(req dto.RegisterRequest) (*dto.AuthResponse, error)`
    - Hash password dengan `bcrypt.GenerateFromPassword()`
    - Simpan user ke database
    - Generate JWT token
  - [x] `Login(req dto.AuthRequest) (*dto.AuthResponse, error)`
    - Cari user by email
    - Bandingkan password dengan `bcrypt.CompareHashAndPassword()`
    - Generate JWT token

### DTO
- [x] Buat `internal/dto/register_request.go`: `{ Name, Email, Password }`
- [x] Buat `internal/dto/auth_request.go`: `{ Email, Password }`
- [x] Buat `internal/dto/auth_response.go`: `{ Token }`

### Auth Handler
- [x] Buat `internal/handler/auth_handler.go`:
  - [x] `Register(c *fiber.Ctx) error` — `POST /api/v1/auth/register`
  - [x] `Login(c *fiber.Ctx) error` — `POST /api/v1/auth/login`
- [x] Pahami cara parse JSON body di Fiber: `c.BodyParser(&req)`
- [x] Pahami cara return JSON di Fiber: `c.JSON(response)`

---

## [x] Fase 7 — JWT Middleware & Protected Routes

- [x] Buat `internal/middleware/jwt_middleware.go`:
  - [x] Ambil header `Authorization: Bearer <token>`
  - [x] Panggil `jwtService.ValidateToken(token)`
  - [x] Simpan email/user ke `c.Locals("userEmail", email)` (ekuivalen SecurityContextHolder)
  - [x] Return 401 jika token tidak ada atau invalid
- [x] Buat `internal/handler/demo_handler.go`:
  - [x] `Hello(c *fiber.Ctx) error` — ambil user dari `c.Locals`
  - [x] `AdminOnly(c *fiber.Ctx) error` — cek role user
- [x] Setup routing di `cmd/main.go`:
  ```go
  app := fiber.New()
  // Public
  app.Post("/api/v1/auth/register", authHandler.Register)
  app.Post("/api/v1/auth/login", authHandler.Login)
  // Protected
  app.Use(middleware.JwtMiddleware)
  app.Get("/api/v1/demo/hello", demoHandler.Hello)
  app.Get("/api/v1/demo/admin", demoHandler.AdminOnly)
  ```

---

## [x] Fase 8 — Testing

- [x] Test register via curl:
  ```bash
  curl -X POST http://localhost:8081/api/v1/auth/register \
    -H "Content-Type: application/json" \
    -d '{"name":"John","email":"john@example.com","password":"secret123"}'
  ```
- [x] Test login & dapat token
- [x] Test protected endpoint dengan token valid
- [x] Test protected endpoint tanpa token — harus 401
- [x] Bandingkan response time Go vs Spring Boot

---

## [x] Fase 9 — Polish & README

- [x] Buat `README.md` yang lengkap (mirip baseapp)
- [x] Tambahkan `.gitignore` (exclude `.env`, binary, dll.)
- [x] Tulis unit test untuk `jwt_service.go`
- [ ] Push ke GitHub: `git push -u origin main`

---

## Referensi Belajar Go

| Topik | Sumber |
|---|---|
| Go Tour (resmi) | https://tour.golang.org |
| Go by Example | https://gobyexample.com |
| Fiber Docs | https://docs.gofiber.io |
| GORM Docs | https://gorm.io/docs |
| golang-jwt | https://github.com/golang-jwt/jwt |
| Effective Go | https://go.dev/doc/effective_go |
