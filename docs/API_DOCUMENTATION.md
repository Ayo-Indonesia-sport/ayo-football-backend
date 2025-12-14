# AYO Football API Documentation

## Deskripsi Sistem

Sistem backend API untuk mengelola tim-tim sepakbola dibawah naungan Perusahaan XYZ. API ini dibangun menggunakan **Golang** dengan **GIN Framework** dan mengikuti prinsip **Clean Architecture** dan **SOLID**.

## Fitur Utama

1. **Pengelolaan Tim Sepakbola** - CRUD operasi untuk data tim
2. **Pengelolaan Pemain** - CRUD operasi untuk data pemain dengan validasi nomor punggung unik per tim
3. **Pengelolaan Jadwal Pertandingan** - Penjadwalan pertandingan antar tim
4. **Pencatatan Hasil Pertandingan** - Pencatatan skor dan pencetak gol
5. **Laporan/Report** - Laporan hasil pertandingan, top scorer, dan akumulasi kemenangan

## Tech Stack

- **Language**: Go 1.23+
- **Framework**: Gin v1.9.1
- **Database**: PostgreSQL (dengan dukungan MySQL)
- **ORM**: GORM v1.25.5
- **Authentication**: JWT (JSON Web Token)
- **Architecture**: Clean Architecture

---

## Base URL

### Production (Live)
```
https://ayo-football-api-production.up.railway.app
```

API Version 1:
```
https://ayo-football-api-production.up.railway.app/api/v1
```

### Local Development
```
http://localhost:8080
```

API Version 1:
```
http://localhost:8080/api/v1
```

---

## Authentication

### JWT Token

API menggunakan JWT untuk autentikasi. Sertakan token di header:

```
Authorization: Bearer <your_jwt_token>
```

### Default Admin Credentials

```
Email: admin@ayofootball.com
Password: Admin@123
```

### Role-Based Access

| Role | Akses |
|------|-------|
| `admin` | Full access (CRUD semua data) |
| `user` | Read-only access |

---

## Response Format

### Success Response
```json
{
  "success": true,
  "message": "Operation successful",
  "data": { ... },
  "meta": {
    "current_page": 1,
    "per_page": 10,
    "total_items": 100,
    "total_pages": 10
  }
}
```

### Error Response
```json
{
  "success": false,
  "message": "Error description",
  "error": "Detailed error message"
}
```

---

## API Endpoints

### 1. Health Check

#### GET /health
Cek status API.

**Response:**
```json
{
  "status": "healthy",
  "service": "ayo-football-api"
}
```

---

### 2. Authentication

#### POST /api/v1/auth/login
Login dan dapatkan JWT token.

**Request Body:**
```json
{
  "email": "admin@ayofootball.com",
  "password": "Admin@123"
}
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Login successful",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": "8c9acfdd-eb81-4370-9577-c56cc403e2d7",
      "email": "admin@ayofootball.com",
      "name": "Admin",
      "role": "admin"
    }
  }
}
```

#### POST /api/v1/auth/register
Registrasi user baru.

**Request Body:**
```json
{
  "name": "John Doe",
  "email": "john@example.com",
  "password": "password123"
}
```

**Response (201 Created):**
```json
{
  "success": true,
  "message": "User registered successfully",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "email": "john@example.com",
    "name": "John Doe",
    "role": "user"
  }
}
```

#### GET /api/v1/auth/profile
Dapatkan profil user yang sedang login.

**Headers:**
```
Authorization: Bearer <token>
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Profile retrieved successfully",
  "data": {
    "id": "8c9acfdd-eb81-4370-9577-c56cc403e2d7",
    "email": "admin@ayofootball.com",
    "name": "Admin",
    "role": "admin"
  }
}
```

---

### 3. Teams (Pengelolaan Tim)

Informasi yang dicatat: **nama tim, logo tim, tahun berdiri, alamat markas tim, kota markas tim**

#### GET /api/v1/teams
Dapatkan semua tim dengan pagination.

**Query Parameters:**
| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| page | int | 1 | Nomor halaman |
| limit | int | 10 | Jumlah item per halaman (max: 100) |
| search | string | - | Cari berdasarkan nama atau kota |

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Teams retrieved successfully",
  "data": [
    {
      "id": "f21a2c88-7eec-4024-97ed-6b3351dab67b",
      "name": "Manchester United",
      "logo": "https://example.com/mu-logo.png",
      "founded_year": 1878,
      "address": "Sir Matt Busby Way",
      "city": "Manchester",
      "created_at": "2025-12-14T09:00:57Z",
      "updated_at": "2025-12-14T09:00:57Z"
    }
  ],
  "meta": {
    "current_page": 1,
    "per_page": 10,
    "total_items": 1,
    "total_pages": 1
  }
}
```

#### GET /api/v1/teams/:id
Dapatkan tim berdasarkan ID.

**Query Parameters:**
| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| with_players | boolean | false | Sertakan daftar pemain |

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Team retrieved successfully",
  "data": {
    "id": "f21a2c88-7eec-4024-97ed-6b3351dab67b",
    "name": "Manchester United",
    "logo": "https://example.com/mu-logo.png",
    "founded_year": 1878,
    "address": "Sir Matt Busby Way",
    "city": "Manchester",
    "players": [
      {
        "id": "765c50ad-0fd3-448d-b737-6211eec03050",
        "name": "Marcus Rashford",
        "position": "forward",
        "jersey_number": 10
      }
    ],
    "created_at": "2025-12-14T09:00:57Z",
    "updated_at": "2025-12-14T09:00:57Z"
  }
}
```

#### POST /api/v1/teams
Tambah tim baru (Admin only).

**Headers:**
```
Authorization: Bearer <admin_token>
Content-Type: application/json
```

**Request Body:**
```json
{
  "name": "Persija Jakarta",
  "logo": "https://example.com/persija-logo.png",
  "founded_year": 1928,
  "address": "Jl. Casablanca No.1",
  "city": "Jakarta"
}
```

**Validation Rules:**
| Field | Rule |
|-------|------|
| name | Required, 2-255 karakter |
| logo | Optional, valid URL, max 500 karakter |
| founded_year | Required, 1800-2100 |
| address | Optional, max 500 karakter |
| city | Required, 2-100 karakter |

**Response (201 Created):**
```json
{
  "success": true,
  "message": "Team created successfully",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "name": "Persija Jakarta",
    "logo": "https://example.com/persija-logo.png",
    "founded_year": 1928,
    "address": "Jl. Casablanca No.1",
    "city": "Jakarta",
    "created_at": "2025-12-14T10:00:00Z",
    "updated_at": "2025-12-14T10:00:00Z"
  }
}
```

#### PUT /api/v1/teams/:id
Update data tim (Admin only).

**Headers:**
```
Authorization: Bearer <admin_token>
Content-Type: application/json
```

**Request Body:**
```json
{
  "name": "Persija Jakarta Updated",
  "city": "DKI Jakarta"
}
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Team updated successfully",
  "data": { ... }
}
```

#### DELETE /api/v1/teams/:id
Hapus tim - **Soft Delete** (Admin only).

**Headers:**
```
Authorization: Bearer <admin_token>
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Team deleted successfully",
  "data": null
}
```

---

### 4. Players (Pengelolaan Pemain)

Informasi yang dicatat: **nama pemain, tinggi badan, berat badan, posisi pemain, nomor punggung**

**Aturan:**
- 1 pemain hanya dapat bernaung pada 1 tim
- 1 tim dapat memiliki banyak pemain
- **Nomor punggung harus unik dalam 1 tim**

#### Posisi Pemain

| Value | Indonesian Name |
|-------|-----------------|
| `forward` | Penyerang |
| `midfielder` | Gelandang |
| `defender` | Bertahan |
| `goalkeeper` | Penjaga Gawang |

#### GET /api/v1/players
Dapatkan semua pemain dengan pagination.

**Query Parameters:**
| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| page | int | 1 | Nomor halaman |
| limit | int | 10 | Jumlah item per halaman (max: 100) |
| search | string | - | Cari berdasarkan nama pemain |
| team_id | uuid | - | Filter berdasarkan tim |

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Players retrieved successfully",
  "data": [
    {
      "id": "765c50ad-0fd3-448d-b737-6211eec03050",
      "team_id": "f21a2c88-7eec-4024-97ed-6b3351dab67b",
      "name": "Marcus Rashford",
      "height": 180,
      "weight": 70,
      "position": "forward",
      "position_name": "Penyerang",
      "jersey_number": 10,
      "team": {
        "id": "f21a2c88-7eec-4024-97ed-6b3351dab67b",
        "name": "Manchester United",
        "logo": "https://example.com/mu-logo.png",
        "city": "Manchester"
      },
      "created_at": "2025-12-14T09:01:31Z",
      "updated_at": "2025-12-14T09:01:31Z"
    }
  ],
  "meta": {
    "current_page": 1,
    "per_page": 10,
    "total_items": 1,
    "total_pages": 1
  }
}
```

#### GET /api/v1/players/:id
Dapatkan pemain berdasarkan ID.

#### POST /api/v1/players
Tambah pemain baru (Admin only).

**Headers:**
```
Authorization: Bearer <admin_token>
Content-Type: application/json
```

**Request Body:**
```json
{
  "team_id": "f21a2c88-7eec-4024-97ed-6b3351dab67b",
  "name": "Marcus Rashford",
  "height": 180,
  "weight": 70,
  "position": "forward",
  "jersey_number": 10
}
```

**Validation Rules:**
| Field | Rule |
|-------|------|
| team_id | Required, valid UUID, tim harus ada |
| name | Required, 2-255 karakter |
| height | Required, 100-250 cm |
| weight | Required, 30-200 kg |
| position | Required, salah satu dari: forward, midfielder, defender, goalkeeper |
| jersey_number | Required, 1-99, **unik dalam 1 tim** |

**Response (201 Created):**
```json
{
  "success": true,
  "message": "Player created successfully",
  "data": {
    "id": "765c50ad-0fd3-448d-b737-6211eec03050",
    "team_id": "f21a2c88-7eec-4024-97ed-6b3351dab67b",
    "name": "Marcus Rashford",
    "height": 180,
    "weight": 70,
    "position": "forward",
    "position_name": "Penyerang",
    "jersey_number": 10,
    "created_at": "2025-12-14T09:01:31Z",
    "updated_at": "2025-12-14T09:01:31Z"
  }
}
```

**Error - Nomor Punggung Sudah Digunakan (409 Conflict):**
```json
{
  "success": false,
  "message": "Jersey number is already taken by another player in this team",
  "error": null
}
```

#### PUT /api/v1/players/:id
Update data pemain (Admin only).

#### DELETE /api/v1/players/:id
Hapus pemain - **Soft Delete** (Admin only).

---

### 5. Matches (Pengelolaan Jadwal Pertandingan)

Informasi yang dicatat: **tanggal pertandingan, waktu pertandingan, tim tuan rumah, tim tamu**

#### Status Pertandingan

| Value | Description |
|-------|-------------|
| `scheduled` | Pertandingan terjadwal |
| `ongoing` | Pertandingan sedang berlangsung |
| `completed` | Pertandingan selesai |
| `cancelled` | Pertandingan dibatalkan |

#### GET /api/v1/matches
Dapatkan semua pertandingan dengan pagination.

**Query Parameters:**
| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| page | int | 1 | Nomor halaman |
| limit | int | 10 | Jumlah item per halaman (max: 100) |
| team_id | uuid | - | Filter berdasarkan tim |
| status | string | - | Filter berdasarkan status |
| start_date | date | - | Filter tanggal mulai (YYYY-MM-DD) |
| end_date | date | - | Filter tanggal akhir (YYYY-MM-DD) |

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Matches retrieved successfully",
  "data": [
    {
      "id": "80470462-42b4-4779-b20d-02b4f30fa5c1",
      "match_date": "2025-12-20",
      "match_time": "15:00",
      "home_team_id": "f21a2c88-7eec-4024-97ed-6b3351dab67b",
      "away_team_id": "5316c5a8-0f42-4b21-8649-a8b0e9bd2f30",
      "home_score": 2,
      "away_score": 1,
      "status": "completed",
      "status_name": "Completed",
      "home_team": {
        "id": "f21a2c88-7eec-4024-97ed-6b3351dab67b",
        "name": "Manchester United",
        "logo": "https://example.com/mu-logo.png",
        "city": "Manchester"
      },
      "away_team": {
        "id": "5316c5a8-0f42-4b21-8649-a8b0e9bd2f30",
        "name": "Liverpool FC",
        "logo": "https://example.com/lfc-logo.png",
        "city": "Liverpool"
      },
      "match_result": "home_win",
      "result_display": "Home Team Win",
      "created_at": "2025-12-14T09:01:31Z",
      "updated_at": "2025-12-14T09:01:43Z"
    }
  ],
  "meta": {
    "current_page": 1,
    "per_page": 10,
    "total_items": 1,
    "total_pages": 1
  }
}
```

#### GET /api/v1/matches/:id
Dapatkan detail pertandingan berdasarkan ID (termasuk goals).

#### POST /api/v1/matches
Tambah jadwal pertandingan baru (Admin only).

**Headers:**
```
Authorization: Bearer <admin_token>
Content-Type: application/json
```

**Request Body:**
```json
{
  "match_date": "2025-12-20",
  "match_time": "15:00",
  "home_team_id": "f21a2c88-7eec-4024-97ed-6b3351dab67b",
  "away_team_id": "5316c5a8-0f42-4b21-8649-a8b0e9bd2f30"
}
```

**Validation Rules:**
| Field | Rule |
|-------|------|
| match_date | Required, format YYYY-MM-DD |
| match_time | Required, format HH:MM |
| home_team_id | Required, valid UUID, harus berbeda dari away_team_id |
| away_team_id | Required, valid UUID, harus berbeda dari home_team_id |

**Response (201 Created):**
```json
{
  "success": true,
  "message": "Match created successfully",
  "data": {
    "id": "80470462-42b4-4779-b20d-02b4f30fa5c1",
    "match_date": "2025-12-20",
    "match_time": "15:00",
    "home_team_id": "f21a2c88-7eec-4024-97ed-6b3351dab67b",
    "away_team_id": "5316c5a8-0f42-4b21-8649-a8b0e9bd2f30",
    "home_score": null,
    "away_score": null,
    "status": "scheduled",
    "status_name": "Scheduled",
    "result_display": "Not Played",
    "created_at": "2025-12-14T09:01:31Z",
    "updated_at": "2025-12-14T09:01:31Z"
  }
}
```

#### PUT /api/v1/matches/:id
Update data pertandingan (Admin only).

#### DELETE /api/v1/matches/:id
Hapus pertandingan - **Soft Delete** (Admin only).

---

### 6. Record Match Result (Pencatatan Hasil Pertandingan)

Informasi yang dicatat: **total skor akhir, pemain yang mencetak gol, waktu terjadinya gol**

#### POST /api/v1/matches/:id/result
Catat hasil pertandingan (Admin only).

**Headers:**
```
Authorization: Bearer <admin_token>
Content-Type: application/json
```

**Request Body:**
```json
{
  "home_score": 2,
  "away_score": 1,
  "goals": [
    {
      "player_id": "765c50ad-0fd3-448d-b737-6211eec03050",
      "team_id": "f21a2c88-7eec-4024-97ed-6b3351dab67b",
      "minute": 23,
      "is_own_goal": false
    },
    {
      "player_id": "765c50ad-0fd3-448d-b737-6211eec03050",
      "team_id": "f21a2c88-7eec-4024-97ed-6b3351dab67b",
      "minute": 78,
      "is_own_goal": false
    }
  ]
}
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Match result recorded successfully",
  "data": {
    "id": "80470462-42b4-4779-b20d-02b4f30fa5c1",
    "match_date": "2025-12-20",
    "match_time": "15:00",
    "home_team_id": "f21a2c88-7eec-4024-97ed-6b3351dab67b",
    "away_team_id": "5316c5a8-0f42-4b21-8649-a8b0e9bd2f30",
    "home_score": 2,
    "away_score": 1,
    "status": "completed",
    "status_name": "Completed",
    "home_team": {
      "id": "f21a2c88-7eec-4024-97ed-6b3351dab67b",
      "name": "Manchester United",
      "logo": "https://example.com/mu-logo.png",
      "city": "Manchester"
    },
    "away_team": {
      "id": "5316c5a8-0f42-4b21-8649-a8b0e9bd2f30",
      "name": "Liverpool FC",
      "logo": "https://example.com/lfc-logo.png",
      "city": "Liverpool"
    },
    "goals": [
      {
        "id": "bc85a968-a800-4a1d-9cba-324a1b8b5b28",
        "match_id": "80470462-42b4-4779-b20d-02b4f30fa5c1",
        "player_id": "765c50ad-0fd3-448d-b737-6211eec03050",
        "player_name": "Marcus Rashford",
        "team_id": "f21a2c88-7eec-4024-97ed-6b3351dab67b",
        "team_name": "Manchester United",
        "minute": 23,
        "is_own_goal": false
      },
      {
        "id": "0e4d5f1d-a372-4595-9f5d-b6d1e5f52bb7",
        "match_id": "80470462-42b4-4779-b20d-02b4f30fa5c1",
        "player_id": "765c50ad-0fd3-448d-b737-6211eec03050",
        "player_name": "Marcus Rashford",
        "team_id": "f21a2c88-7eec-4024-97ed-6b3351dab67b",
        "team_name": "Manchester United",
        "minute": 78,
        "is_own_goal": false
      }
    ],
    "match_result": "home_win",
    "result_display": "Home Team Win",
    "created_at": "2025-12-14T09:01:31Z",
    "updated_at": "2025-12-14T09:01:43Z"
  }
}
```

---

### 7. Reports (Data Report)

Informasi yang ditampilkan:
- Jadwal pertandingan
- Tim home & away
- Skor akhir
- Status akhir pertandingan (Tim Home Menang/Tim Away Menang/Draw)
- Pemain pencetak gol terbanyak
- Akumulasi total kemenangan tim home
- Akumulasi total kemenangan tim away

#### GET /api/v1/reports/matches
Dapatkan semua laporan pertandingan yang sudah selesai.

**Query Parameters:**
| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| page | int | 1 | Nomor halaman |
| limit | int | 10 | Jumlah item per halaman (max: 100) |

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Match reports retrieved successfully",
  "data": [
    {
      "match": {
        "id": "80470462-42b4-4779-b20d-02b4f30fa5c1",
        "match_date": "2025-12-20",
        "match_time": "15:00",
        "status": "completed"
      },
      "home_team": {
        "id": "f21a2c88-7eec-4024-97ed-6b3351dab67b",
        "name": "Manchester United",
        "logo": "https://example.com/mu-logo.png",
        "city": "Manchester"
      },
      "away_team": {
        "id": "5316c5a8-0f42-4b21-8649-a8b0e9bd2f30",
        "name": "Liverpool FC",
        "logo": "https://example.com/lfc-logo.png",
        "city": "Liverpool"
      },
      "home_score": 2,
      "away_score": 1,
      "match_result": "home_win",
      "match_result_display": "Home Team Win",
      "goals": [
        {
          "id": "bc85a968-a800-4a1d-9cba-324a1b8b5b28",
          "player_id": "765c50ad-0fd3-448d-b737-6211eec03050",
          "player_name": "Marcus Rashford",
          "team_id": "f21a2c88-7eec-4024-97ed-6b3351dab67b",
          "minute": 23,
          "is_own_goal": false
        }
      ],
      "top_scorer": {
        "player_id": "765c50ad-0fd3-448d-b737-6211eec03050",
        "player_name": "Marcus Rashford",
        "team_id": "f21a2c88-7eec-4024-97ed-6b3351dab67b",
        "team_name": "Manchester United",
        "goal_count": 2
      },
      "home_team_total_wins": 1,
      "away_team_total_wins": 0
    }
  ],
  "meta": {
    "current_page": 1,
    "per_page": 10,
    "total_items": 1,
    "total_pages": 1
  }
}
```

#### GET /api/v1/reports/matches/:id
Dapatkan laporan detail untuk pertandingan tertentu.

#### GET /api/v1/reports/top-scorers
Dapatkan daftar top scorer (pencetak gol terbanyak).

**Query Parameters:**
| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| limit | int | 10 | Jumlah top scorer (max: 100) |

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Top scorers retrieved successfully",
  "data": [
    {
      "player_id": "765c50ad-0fd3-448d-b737-6211eec03050",
      "player_name": "Marcus Rashford",
      "team_id": "f21a2c88-7eec-4024-97ed-6b3351dab67b",
      "team_name": "Manchester United",
      "goal_count": 2
    }
  ]
}
```

---

## Error Codes

| HTTP Code | Description |
|-----------|-------------|
| 200 | OK - Request berhasil |
| 201 | Created - Data berhasil dibuat |
| 400 | Bad Request - Request tidak valid |
| 401 | Unauthorized - Token tidak valid atau tidak ada |
| 403 | Forbidden - Tidak memiliki akses (bukan admin) |
| 404 | Not Found - Data tidak ditemukan |
| 409 | Conflict - Data konflik (misal: nomor punggung sudah digunakan) |
| 500 | Internal Server Error - Error server |

### Contoh Error Responses

**400 Bad Request:**
```json
{
  "success": false,
  "message": "Invalid request body",
  "error": "Key: 'CreateTeamRequest.Name' Error:Field validation for 'Name' failed on the 'required' tag"
}
```

**401 Unauthorized:**
```json
{
  "success": false,
  "message": "Authorization header is required",
  "error": null
}
```

**403 Forbidden:**
```json
{
  "success": false,
  "message": "Admin access required",
  "error": null
}
```

**404 Not Found:**
```json
{
  "success": false,
  "message": "Team not found",
  "error": null
}
```

**409 Conflict:**
```json
{
  "success": false,
  "message": "Jersey number is already taken by another player in this team",
  "error": null
}
```

---

## Soft Delete

Semua operasi DELETE menggunakan mekanisme **Soft Delete**. Data tidak benar-benar dihapus dari database, melainkan diberi tanda `deleted_at` dengan timestamp. Data yang sudah di-soft delete tidak akan muncul di query biasa.

---

## Menjalankan API

### Prerequisites

- Go 1.23+
- PostgreSQL 15+

### Setup

1. Clone repository
2. Copy `.env.example` ke `.env` dan sesuaikan konfigurasi
3. Buat database PostgreSQL
4. Jalankan aplikasi:

```bash
go run cmd/api/main.go
```

### Environment Variables

```env
# Server
SERVER_PORT=8080
GIN_MODE=debug

# Database
DB_DRIVER=postgres
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=ayo_football
DB_SSLMODE=disable

# JWT
JWT_SECRET=your-super-secret-jwt-key
JWT_EXPIRATION_HOURS=24

# Admin
ADMIN_EMAIL=admin@ayofootball.com
ADMIN_PASSWORD=Admin@123
```

---

## Postman Collection

Import file `docs/postman_collection.json` ke Postman untuk testing semua endpoints.

### Cara Penggunaan:

1. Import collection ke Postman
2. Jalankan request "Login" terlebih dahulu untuk mendapatkan token
3. Token akan otomatis disimpan ke collection variable
4. Request lain akan menggunakan token tersebut secara otomatis

---

## Quick Test (Live API)

### Health Check
```bash
curl https://ayo-football-api-production.up.railway.app/health
```

### Login
```bash
curl -X POST https://ayo-football-api-production.up.railway.app/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@ayofootball.com","password":"Admin@123"}'
```

### Get All Teams
```bash
curl https://ayo-football-api-production.up.railway.app/api/v1/teams
```

### Create Team (dengan token)
```bash
curl -X POST https://ayo-football-api-production.up.railway.app/api/v1/teams \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name":"Persija Jakarta","logo":"https://example.com/logo.png","founded_year":1928,"address":"GBK Stadium","city":"Jakarta"}'
```

---

## Repository

- **GitHub**: https://github.com/Ayo-Indonesia-sport/ayo-football-backend
- **Live API**: https://ayo-football-api-production.up.railway.app

---

## Author

**AYO Indonesia Sport** - Technical Test for GO Software Developer
