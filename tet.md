# Dokumentasi API Docger

## Apa Itu Docger?

Docger adalah backend ticketing system yang dibangun dengan bahasa pemrograman Go (Golang), PostgreSQL, dan Docker. Sistem ini menyediakan API RESTful untuk mengelola ticket dukungan (support ticket) dengan fitur autentikasi dan role-based access control.

## Fitur Utama

- **Autentikasi Pengguna**: Sistem register dan login dengan JWT (JSON Web Token)
- **Manajemen Ticket**: CRUD lengkap untuk ticket (Create, Read, Update, Delete)
- **Role-Based Access**: Dua role pengguna - user dan admin
- **Admin Monitoring**: Admin dapat melihat semua ticket dan menugaskan ticket ke user tertentu

## Cara Menjalankan

```bash
docker compose up -d
```

Aplikasi akan berjalan di `http://localhost:8080`

## Tech Stack

- **Bahasa**: Go 1.23
- **Database**: PostgreSQL 15
- **Cache**: Redis 7
- **Admin Tool**: pgAdmin 4
- **Framework**: Gin
- **Auth**: JWT
- **Container**: Docker

## Docker Services

| Service   | Port   | Description        |
|-----------|--------|--------------------|
| app       | 8080   | Go application     |
| postgres  | 5432   | PostgreSQL database |
| pgadmin   | 5050   | pgAdmin (PostgreSQL GUI) |
| redis     | 6379   | Redis cache        |

### Akses pgAdmin
- URL: http://localhost:5050
- Email: admin@docger.com
- Password: admin123

## Testing API dengan Postman

### 1. Registrasi Pengguna Baru

**Endpoint:** `POST /api/auth/register`

```json
{
  "username": "nama_user",
  "email": "email@example.com",
  "password": "password123"
}
```

### 2. Login Pengguna

**Endpoint:** `POST /api/auth/login`

```json
{
  "username": "nama_user",
  "password": "password123"
}
```

Response akan berisi token JWT yang perlu digunakan untuk request selanjutnya.

### 3. Mendapatkan Profile Pengguna (Terproteksi)

**Endpoint:** `GET /api/profile`

Headers:
- `Authorization: Bearer <token_jwt>`

### 4. Membuat Ticket Baru

**Endpoint:** `POST /api/tickets`

Headers:
- `Authorization: Bearer <token_jwt>`

```json
{
  "title": "Judul Ticket",
  "description": "Deskripsi masalah yang dihadapi",
  "priority": "high"
}
```

Priority: `low`, `medium`, `high`

### 5. Melihat Ticket Milik Sendiri

**Endpoint:** `GET /api/tickets/my`

Headers:
- `Authorization: Bearer <token_jwt>`

### 6. Melihat Detail Ticket

**Endpoint:** `GET /api/tickets/:id`

Headers:
- `Authorization: Bearer <token_jwt>`

### 7. Mengupdate Ticket

**Endpoint:** `PUT /api/tickets/:id`

Headers:
- `Authorization: Bearer <token_jwt>`

```json
{
  "title": "Judul Baru",
  "description": "Deskripsi baru",
  "status": "in_progress",
  "priority": "high"
}
```

Status: `open`, `in_progress`, `resolved`, `closed`

### 8. Menghapus Ticket

**Endpoint:** `DELETE /api/tickets/:id`

Headers:
- `Authorization: Bearer <token_jwt>`

### 9. Melihat Semua User (Admin Only)

**Endpoint:** `GET /api/admin/users`

Headers:
- `Authorization: Bearer <token_admin>`

### 10. Melihat Semua Ticket (Admin Only)

**Endpoint:** `GET /api/admin/tickets`

Headers:
- `Authorization: Bearer <token_admin>`

### 11. Menugaskan Ticket ke User (Admin Only)

**Endpoint:** `POST /api/admin/tickets/:id/assign`

Headers:
- `Authorization: Bearer <token_admin>`

```json
{
  "assigned_to": 2
}
```

## Cara Menggunakan Token di Postman

1. Setelah login, copy token dari response
2. Di Postman, pilih tab **Headers**
3. Tambahkan key: `Authorization`
4. Value: `Bearer <paste_token_disini>`
5. Klik Send

## Struktur Project

```
docger/
├── cmd/
│   ├── main.go              # Entry point aplikasi utama
│   └── worker/main.go       # Worker service (opsional)
├── internal/
│   ├── handler/             # HTTP handlers
│   │   ├── user_handler.go
│   │   ├── ticket_handler.go
│   │   └── middleware.go
│   ├── model/               # Data models
│   │   ├── user.go
│   │   └── ticket.go
│   ├── repository/          # Database operations
│   │   ├── db.go
│   │   ├── user_repository.go
│   │   └── ticket_repository.go
│   └── service/             # Business logic
│       ├── user_service.go
│       └── ticket_service.go
├── database/
│   └── migration.sql        # Schema database
├── Dockerfile
├── docker-compose.yml
├── tet.md                    # Dokumentasi API (Indonesia)
└── eror.md                   # Dokumentasi error & perkembangan
```

## Default Credentials

Untuk membuat user admin, gunakan kode rahasia saat registrasi:

```json
{
  "username": "admin",
  "email": "admin@docger.com",
  "password": "admin123",
  "admin_secret": "admin123secret"
}
```

Untuk user biasa, cukup registrasi tanpa admin_secret:

```json
{
  "username": "user",
  "email": "user@example.com",
  "password": "user123"
}
```

## Environment Variables

| Variable     | Description                    | Default                          |
|--------------|--------------------------------|----------------------------------|
| DATABASE_URL | PostgreSQL connection string   | postgres://docger:docgerpass@postgres:5432/docger?sslmode=disable |
| JWT_SECRET   | JWT signing secret            | docger_secret_key_2024          |
| PORT         | Server port                    | 8080                             |
