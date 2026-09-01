# 🛡️ Go Sentinel — REST API with JWT Authentication (Golang)

**Go Sentinel** adalah implementasi REST API backend yang aman menggunakan bahasa pemrograman **Golang**, web framework **Fiber v2**, ORM **GORM**, dan database **PostgreSQL 16**.

Project ini dibangun sebagai referensi arsitektur bersih (*Clean Architecture-lite*) dan perbandingan langsung dengan implementasi ekosistem **Java / Spring Boot (Spring Security + JPA)**.

---

## 🔐 Logika Autentikasi & Cara Kerja JWT

Aplikasi ini menggunakan mekanisme **Stateless Authentication** berbasis **JSON Web Token (JWT)** standar RFC 7519. Server tidak menyimpan sesi user di memory/database (*zero session lookup*), melainkan mempercayai payload yang sudah ditandatangani secara kriptografis (*cryptographically signed*).

### 1. Alur Login & Pembuatan Token (*Token Issuance*)

```text
[ Client ]                        [ Fiber Server ]                   [ PostgreSQL ]
    │                                    │                                 │
    ├── 1. POST /api/v1/auth/login ─────>│                                 │
    │      { email, password }           ├── 2. Cari user by email ───────>│
    │                                    │<── 3. Return user data (hash) ──┤
    │                                    │                                 │
    │                                    ├── 4. bcrypt.CompareHashAndPassword()
    │                                    │      (Verifikasi password cocok)
    │                                    │                                 │
    │                                    ├── 5. jwtService.GenerateToken() │
    │                                    │      - Header: Alg: HS256, Typ: JWT
    │                                    │      - Payload Claims: email, role, exp (24h)
    │                                    │      - Signature: HMAC-SHA256(secretKey)
    │                                    │                                 │
    │<── 6. Return JSON { token, ... } ──┤                                 │
```

---

### 2. Alur Request Terproteksi (*Protected Request & Claims Extraction*)

```text
[ Client ]                     [ JWTMiddleware ]                [ Handler / Controller ]
    │                                  │                                  │
    ├── 1. GET /api/v1/demo/hello ────>│                                  │
    │      Header: Authorization:      │                                  │
    │      Bearer <jwt_token>          ├── 2. Ekstrak string token        │
    │                                  ├── 3. Validasi signature (Secret) │
    │                                  ├── 4. Cek Expiration Time (exp)   │
    │                                  │                                  │
    │                                  │── ❌ Token Invalid/Expired ──────> Return HTTP 401
    │                                  │                                  │
    │                                  ├── 5. Simpan Claims ke Context:   │
    │                                  │      c.Locals("userEmail", email)│
    │                                  │      c.Locals("userRole", role)  │
    │                                  │                                  │
    │                                  ├── 6. Lanjutkan (c.Next()) ──────>│
    │                                  │                                  ├── 7. Ambil user dari c.Locals
    │<── 8. Return HTTP 200 OK ────────┴──────────────────────────────────┤
```

---

### 3. Logika Role-Based Access Control (RBAC)

* Middleware `RequireRole("ADMIN")` bertindak sebagai *gatekeeper* sekunder.
* Setelah `JWTMiddleware` memvalidasi keabsahan token dan menyuntikkan `userRole` ke `c.Locals`, `RequireRole` akan memeriksa:
  * Jika `userRole == "ADMIN"` $\rightarrow$ Akses diberikan (`c.Next()`).
  * Jika `userRole != "ADMIN"` $\rightarrow$ Request langsung diputus dengan **HTTP 403 Forbidden**.

---

## 🚀 Tech Stack

| Komponen | Teknologi | Keterangan |
|---|---|---|
| **Language** | Go 1.22+ | Bahasa utama |
| **Web Framework** | Fiber v2 | HTTP engine berbasis Fasthttp (super cepat & hemat memori) |
| **Database** | PostgreSQL 16 | Relational database |
| **ORM** | GORM | Object Relational Mapping & Auto Migration |
| **JWT** | golang-jwt/jwt v5 | Autentikasi token standar RFC 7519 (HMAC-SHA256) |
| **Password Hashing** | Bcrypt (golang.org/x/crypto) | Hashing password dengan salt otomatis (anti brute-force) |
| **Environment** | godotenv | Memuat konfigurasi dari file `.env` |
| **Live Reload** | Air (air-verse/air) | Re-compile otomatis saat pengembangan |

---

## 📁 Struktur Folder Project

```text
go-sentinel/
├── cmd/
│   └── main.go                 # Entry point aplikasi & Dependency Injection wiring
├── internal/
│   ├── config/                 # Config loader & PostgreSQL connection
│   │   ├── config.go
│   │   └── database.go
│   ├── dto/                    # Data Transfer Objects (Request/Response payload)
│   │   └── auth_dto.go
│   ├── handler/                # HTTP Controllers (Fiber Request/Response handler)
│   │   ├── auth_handler.go
│   │   └── demo_handler.go
│   ├── middleware/             # JWT & Role Authorization middleware
│   │   └── jwt_middleware.go
│   ├── model/                  # Database Entities (GORM Structs)
│   │   └── user.go
│   ├── repository/             # Data Access Layer (GORM queries)
│   │   └── user_repository.go
│   └── service/                # Business Logic & JWT Services
│       ├── auth_service.go
│       ├── jwt_service.go
│       └── jwt_service_test.go # Unit tests
├── .air.toml                   # Konfigurasi live reload Air
├── .env                        # Konfigurasi lokal (rahasia)
├── .env.example                # Template konfigurasi
├── .gitignore
├── go.mod
├── go.sum
└── README.md
```

---

## ⚙️ Persiapan & Instalasi

### 1. Prasyarat
* **Go 1.22+**
* **PostgreSQL 16**
* **Air** (Opsional, untuk live reload): `go install github.com/air-verse/air@latest`

### 2. Setup Database
Buat database di PostgreSQL:
```bash
createdb gosentinel_db
```

### 3. Konfigurasi Environment (`.env`)
Salin template `.env.example` menjadi `.env`:
```bash
cp .env.example .env
```
Sesuaikan kredensial di dalam file `.env`:
```env
DB_HOST=localhost
DB_PORT=5432
DB_NAME=gosentinel_db
DB_USER=postgres
DB_PASSWORD=
DB_SSLMODE=disable

JWT_SECRET=your_super_secret_jwt_key_here_minimum_32_chars
JWT_EXPIRATION_HOURS=24

APP_PORT=8081
```

---

## 🏃 Cara Menjalankan Aplikasi

### Mode Development (Live Reload dengan Air):
```bash
air
```

### Mode Standar (Go Run):
```bash
go run cmd/main.go
```

Server akan aktif di: `http://localhost:8081`

---

## 🧪 Menjalankan Unit Test

```bash
go test -v ./...
```

---

## 📡 Daftar Endpoint API

### 1. Health Check
* **GET** `/health`
* **Response:**
  ```json
  {
    "status": "UP",
    "message": "Go Sentinel is running smoothly"
  }
  ```

### 2. Autentikasi (Public)
* **POST** `/api/v1/auth/register`
  * **Body:**
    ```json
    {
      "name": "Faiz",
      "email": "faiz@example.com",
      "password": "password123"
    }
    ```
* **POST** `/api/v1/auth/login`
  * **Body:**
    ```json
    {
      "email": "faiz@example.com",
      "password": "password123"
    }
    ```
  * **Response:**
    ```json
    {
      "message": "Login berhasil",
      "data": {
        "token": "eyJhbGciOiJIUzI1NiIsIn...",
        "email": "faiz@example.com",
        "role": "USER"
      }
    }
    ```

### 3. Protected Routes (Header: `Authorization: Bearer <token>`)
* **GET** `/api/v1/demo/hello` (Memerlukan Token Valid)
* **GET** `/api/v1/demo/admin` (Memerlukan Token Valid dengan Role `ADMIN`)

---

## ⚖️ Perbandingan: Go Sentinel vs Spring Boot

| Fitur / Konsep | Spring Boot (Java) | Go Sentinel (Go) |
|---|---|---|
| **Web Server** | Embedded Tomcat / Netty | Fasthttp (Fiber v2) |
| **Startup Time** | ~2 - 4 detik | **< 10 milidetik** ⚡ |
| **Memory Footprint** | ~200MB - 350MB | **~15MB - 30MB** |
| **Filter / Middleware** | `OncePerRequestFilter` | `fiber.Handler` closure |
| **Context User** | `SecurityContextHolder` | `c.Locals("userEmail", email)` |
| **Dependency Injection**| `@Autowired` / Spring IoC | Pure Manual Constructor DI |
| **Error Handling** | `try-catch` / `@ExceptionHandler` | Multiple Returns `(data, err)` |
| **Binary Output** | `.jar` (Memerlukan JVM) | Standalone Native Binary (Tanpa VM) |
