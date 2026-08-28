# PRD — Go Sentinel: REST API dengan JWT Authentication (Golang)

## 1. Latar Belakang

**Go Sentinel** adalah project belajar membangun REST API yang aman menggunakan **Golang** dengan JWT Authentication. Project ini adalah versi Golang dari `baseapp` (Spring Boot), sehingga konsep yang sudah dipelajari (JWT, Auth Flow, Role-based Access) dapat dibandingkan langsung antara dua bahasa/ekosistem yang berbeda.

---

## 2. Tujuan Project

| Tujuan | Keterangan |
|---|---|
| **Belajar Go** | Memahami idiom Go: struct, interface, goroutine, context |
| **Perbandingan** | Membandingkan implementasi JWT di Go vs Java (Spring Boot) |
| **Template** | Base project untuk REST API Go yang production-ready |

---

## 3. Tech Stack

| Teknologi | Pilihan | Keterangan |
|---|---|---|
| **Language** | Go 1.22+ | Bahasa utama |
| **Web Framework** | Fiber v2 | Cepat, mirip Express.js |
| **JWT Library** | golang-jwt/jwt v5 | Paling populer di ekosistem Go |
| **ORM** | GORM | ORM paling populer di Go, mirip Hibernate |
| **Database** | PostgreSQL 16 | Sama dengan Spring Boot project |
| **Config** | godotenv | Membaca .env file |
| **Password Hash** | bcrypt (stdlib) | Bawaan Go, ekuivalen BCryptPasswordEncoder |
| **Build Tool** | Go Modules (go.mod) | Standar Go, ekuivalen pom.xml |

---

## 4. Struktur Folder Project

```
go-sentinel/
├── cmd/
│   └── main.go
├── internal/
│   ├── config/config.go
│   ├── model/user.go
│   ├── repository/user_repository.go
│   ├── service/
│   │   ├── jwt_service.go
│   │   └── auth_service.go
│   ├── handler/
│   │   ├── auth_handler.go
│   │   └── demo_handler.go
│   ├── middleware/jwt_middleware.go
│   └── dto/
├── .env
├── .env.example
├── .gitignore
├── go.mod
├── go.sum
├── PRD.md
├── Tasks.md
└── README.md
```

---

## 5. Fitur yang Dibangun

### Phase 1 — Auth & Security
| Fitur | Endpoint | Status |
|---|---|---|
| Register | POST /api/v1/auth/register | Todo |
| Login | POST /api/v1/auth/login | Todo |
| Protected endpoint | GET /api/v1/demo/hello | Todo |
| Admin-only endpoint | GET /api/v1/demo/admin | Todo |

### Phase 2 — User Management
| Fitur | Endpoint | Status |
|---|---|---|
| Profil sendiri | GET /api/v1/users/me | Todo |
| Update profil | PUT /api/v1/users/me | Todo |

---

## 6. Perbandingan dengan Spring Boot

| Aspek | Spring Boot (baseapp) | Go (go-sentinel) |
|---|---|---|
| Filter/Middleware | OncePerRequestFilter | fiber.Handler |
| Context | SecurityContextHolder | fiber.Ctx |
| DI | Otomatis (Spring IoC) | Manual/Constructor |
| ORM | Hibernate JPA | GORM |
| Config | application.yaml | .env + godotenv |
| Port | 8080 | 8081 |
