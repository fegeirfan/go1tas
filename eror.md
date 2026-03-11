# Dokumentasi Error dan Perkembangan

## Error yang Pernah Terjadi

### 1. Error: "github.com/segmentio/kafka-go requires Go >= 1.23"

**Penyebab:**
- Package `github.com/segmentio/kafka-go` versi terbaru memerlukan Go 1.23
- Project menggunakan Go 1.21

**Solusi:**
- Mengupgrade Go version dari 1.21 ke 1.23 di `go.mod` dan `Dockerfile`

**Perubahan:**
```go
// go.mod
go 1.21 → go 1.23

// Dockerfile
FROM golang:1.21-alpine → FROM golang:1.23-alpine
```

### 2. Error: "malformed go.sum: wrong number of fields"

**Penyebab:**
- File `go.sum` yang ada corrupt atau tidak valid
- Terdapat masalah dengan format go.sum

**Solusi:**
- Menghapus baris yang cek go.sum dan langsung menjalankan `go mod download`
- Mengubah di Dockerfile:
```dockerfile
# Sebelum
RUN if [ -f go.sum ]; then go mod download; else echo "// go.sum is intentionally empty" > go.sum && go mod download; fi

# Sesudah
RUN rm -f go.sum && go mod download
```

### 3. Error: Docker pull bitnami/kafka and bitnami/zookeeper

**Penyebab:**
- Image bitnami/kafka dan bitnami/zookeeper tidak ditemukan atau sudah deprecated

**Solusi:**
- Menghapus service kafka dan zookeeper dari docker-compose.yml
- Hanya menggunakan PostgreSQL dan Aplikasi saja

### 4. Error: Docker compose "version attribute is obsolete"

**Penyebab:**
- Atribut `version` di docker-compose.yml sudah tidak diperlukan di Docker Compose v2

**Solusi:**
- Menghapus baris `version: "3.8"` dari docker-compose.yml (sudah dihapus)

## Perkembangan Terbaru

### Struktur Docker Baru:
- **postgres** - PostgreSQL 15 untuk database
- **app** - Go 1.23 application

### Fitur yang Tersedia:
- User registration dan login dengan JWT
- Role-based access (admin/user)
- CRUD ticket
- Admin monitoring

### Fitur yang Dapat Ditambahkan Nanti:
- Redis untuk caching
- Kafka untuk messaging/queue
