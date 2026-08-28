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

## [ ] Fase 2 — Dasar-dasar Go (Wajib sebelum mulai project)

### Sintaks Dasar
- [ ] Pahami deklarasi variabel: `var` vs `:=`
- [ ] Pahami tipe data dasar: `string`, `int`, `bool`, `float64`
- [ ] Pahami `struct` — ekuivalen `class` di Java
- [ ] Pahami `interface` di Go — berbeda dengan Java (implicit, bukan explicit)
- [ ] Pahami pointer: `*` dan `&` — tidak ada di Java
- [ ] Pahami error handling di Go: `if err != nil` (tidak ada try-catch)

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
- [ ] Pahami struktur `go.mod` — ekuivalen `pom.xml`
- [ ] Pahami `package main` dan `package internal`
- [ ] Pahami cara import package: `import "github.com/gofiber/fiber/v2"`
- [ ] Install dependencies: `go get <package>` (ekuivalen `mvn dependency:resolve`)

---

## [ ] Fase 3 — Init Project Go Sentinel

### Setup Project
- [ ] Inisialisasi Go module: `go mod init github.com/alfaizuna/go-sentinel`
- [ ] Buat struktur folder sesuai PRD.md
- [ ] Install semua dependency:
  ```bash
  go get github.com/gofiber/fiber/v2
  go get github.com/golang-jwt/jwt/v5
  go get gorm.io/gorm
  go get gorm.io/driver/postgres
  go get github.com/joho/godotenv
  go get golang.org/x/crypto/bcrypt
  ```

### Config & Database
- [ ] Buat file `.env` dengan variabel:
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
- [ ] Buat `.env.example` (copy dari `.env` tanpa nilai sensitif)
- [ ] Buat `internal/config/config.go` — load env ke struct Config
- [ ] Buat database: `createdb gosentinel_db`
- [ ] Koneksikan GORM ke PostgreSQL

---

## [ ] Fase 4 — Model & Repository

### Model
- [ ] Buat `internal/model/user.go`:
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
- [ ] Jalankan GORM AutoMigrate untuk membuat tabel `users`
- [ ] Pahami perbedaan: `gorm.Model` di Go vs `@Entity` + Hibernate di Java

### Repository
- [ ] Buat `internal/repository/user_repository.go`:
  - [ ] `FindByEmail(email string) (*model.User, error)`
  - [ ] `ExistsByEmail(email string) bool`
  - [ ] `Save(user *model.User) error`
- [ ] Pahami: GORM di Go lebih eksplisit dibanding Spring Data JPA

---

## [ ] Fase 5 — JWT Service

- [ ] Buat `internal/service/jwt_service.go`:
  - [ ] `GenerateToken(email string) (string, error)` — buat JWT token baru
  - [ ] `ValidateToken(tokenString string) (*Claims, error)` — validasi & ekstrak claims
- [ ] Pahami struct `Claims` yang embed `jwt.RegisteredClaims`
- [ ] Pahami perbedaan: di Go tidak ada `UserDetails` interface, cukup string email
- [ ] Test manual generate & parse token

---

## [ ] Fase 6 — Auth Service & Handler

### Auth Service
- [ ] Buat `internal/service/auth_service.go`:
  - [ ] `Register(req dto.RegisterRequest) (*dto.AuthResponse, error)`
    - Hash password dengan `bcrypt.GenerateFromPassword()`
    - Simpan user ke database
    - Generate JWT token
  - [ ] `Login(req dto.AuthRequest) (*dto.AuthResponse, error)`
    - Cari user by email
    - Bandingkan password dengan `bcrypt.CompareHashAndPassword()`
    - Generate JWT token

### DTO
- [ ] Buat `internal/dto/register_request.go`: `{ Name, Email, Password }`
- [ ] Buat `internal/dto/auth_request.go`: `{ Email, Password }`
- [ ] Buat `internal/dto/auth_response.go`: `{ Token }`

### Auth Handler
- [ ] Buat `internal/handler/auth_handler.go`:
  - [ ] `Register(c *fiber.Ctx) error` — `POST /api/v1/auth/register`
  - [ ] `Login(c *fiber.Ctx) error` — `POST /api/v1/auth/login`
- [ ] Pahami cara parse JSON body di Fiber: `c.BodyParser(&req)`
- [ ] Pahami cara return JSON di Fiber: `c.JSON(response)`

---

## [ ] Fase 7 — JWT Middleware & Protected Routes

- [ ] Buat `internal/middleware/jwt_middleware.go`:
  - [ ] Ambil header `Authorization: Bearer <token>`
  - [ ] Panggil `jwtService.ValidateToken(token)`
  - [ ] Simpan email/user ke `c.Locals("userEmail", email)` (ekuivalen SecurityContextHolder)
  - [ ] Return 401 jika token tidak ada atau invalid
- [ ] Buat `internal/handler/demo_handler.go`:
  - [ ] `Hello(c *fiber.Ctx) error` — ambil user dari `c.Locals`
  - [ ] `AdminOnly(c *fiber.Ctx) error` — cek role user
- [ ] Setup routing di `cmd/main.go`:
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

## [ ] Fase 8 — Testing

- [ ] Test register via curl:
  ```bash
  curl -X POST http://localhost:8081/api/v1/auth/register \
    -H "Content-Type: application/json" \
    -d '{"name":"John","email":"john@example.com","password":"secret123"}'
  ```
- [ ] Test login & dapat token
- [ ] Test protected endpoint dengan token valid
- [ ] Test protected endpoint tanpa token — harus 401
- [ ] Bandingkan response time Go vs Spring Boot:
  ```bash
  # Benchmark sederhana
  time curl -X POST http://localhost:8080/api/v1/auth/login ...  # Spring Boot
  time curl -X POST http://localhost:8081/api/v1/auth/login ...  # Go
  ```

---

## [ ] Fase 9 — Polish & README

- [ ] Buat `README.md` yang lengkap (mirip baseapp)
- [ ] Tambahkan `.gitignore` (exclude `.env`, binary, dll.)
- [ ] Tulis unit test untuk `jwt_service.go`
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
