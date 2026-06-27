# 🚗 SpotSync - Smart Parking & EV Charging Reservation System

![Go Version](https://img.shields.io/badge/Go-1.22%2B-00ADD8?logo=go)
![Echo](https://img.shields.io/badge/Echo-v4-blue)
![GORM](https://img.shields.io/badge/GORM-ORM-green)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-NeonDB-4169E1?logo=postgresql)
![JWT](https://img.shields.io/badge/Auth-JWT-black?logo=jsonwebtokens)
![License](https://img.shields.io/badge/License-MIT-yellow)

A centralized platform for busy airports and malls to manage parking zones, specifically handling the high-demand reservation of limited EV charging spots with concurrency-safe booking.

**Live API URL:** `https://your-deployed-url.com` *(Update after deployment)*

---

## ✨ Features

- 🔐 **JWT Authentication** - Secure login with bcrypt password hashing
- 👮 **Role-Based Access** - Driver and Admin roles with granular permissions
- 🅿️ **Parking Zone Management** - CRUD operations for parking zones (general, EV charging, covered)
- 📊 **Real-Time Availability** - Dynamic spot counting with capacity checks
- 🎯 **Concurrency-Safe Booking** - Database transactions with row-level locking prevents overbooking
- 🛡️ **Input Validation** - Request validation on all endpoints
- 📝 **Clean Architecture** - Strict separation of concerns (Handlers → Services → Repositories)
- 🌐 **CORS Ready** - Configured for frontend integration
- 💚 **Health Check** - Monitor API status
- 🔄 **Graceful Shutdown** - Clean resource cleanup on server stop

---

## 🛠️ Tech Stack

| Technology | Purpose |
|------------|---------|
| **Go 1.22+** | Programming Language |
| **Echo v4** | HTTP Framework |
| **GORM** | ORM |
| **PostgreSQL** | Database (NeonDB) |
| **JWT v5** | Authentication |
| **bcrypt** | Password Hashing |
| **Validator v10** | Request Validation |
| **godotenv** | Environment Variables |

---

## 🚀 Getting Started

### Prerequisites

- **Go** (version 1.22 or higher) - [Download](https://go.dev/dl/)
- **PostgreSQL Database** - We recommend [NeonDB](https://neon.tech) (free tier)
- **Git** - For version control

### 1. Clone the Repository

```bash
git clone https://github.com/Rahima-Akter/spotsync-golang.git
cd spotsync-golang
```
### 2. Set Up the Database

#### Option A: NeonDB (Recommended)

1. Create a free account at https://neon.tech
2. Create a new project named `spotsync`
3. Copy the PostgreSQL connection string.
4. Paste it into your `.env` file.

#### Option B: Local PostgreSQL

```sql
CREATE DATABASE spotsync;
```

### Copy the example env file from - .env.example
- Edit .env with your actual values

### 4. Install Dependencies
```bash
# Download all required packages
go mod download

# Or tidy up (if you've made changes)
go mod tidy
```
### 5. Run the Server
```bash
# Development mode (with live reload)
air

# Or directly with go run
go run ./cmd/server

# Or build and run
go build -o spotsync.exe ./cmd/server
./spotsync.exe
```
# 📡 API Endpoints
### 🌐 Base URL
```bash
http://localhost:8080/api/v1
```
### 🔑 Authentication
| Method | Endpoint         | Access | Description             |
| ------ | ---------------- | ------ | ----------------------- |
| `POST` | `/auth/register` | Public | Register a new user     |
| `POST` | `/auth/login`    | Public | Login and get JWT token |

### 🚕 Parking Zones
| Method   | Endpoint     | Access | Description                      |
| -------- | ------------ | ------ | -------------------------------- |
| `GET`    | `/zones`     | Public | List all zones with availability |
| `GET`    | `/zones/:id` | Public | Get single zone details          |
| `POST`   | `/zones`     | Admin  | Create a new parking zone        |
| `PUT`    | `/zones/:id` | Admin  | Update a parking zone            |
| `DELETE` | `/zones/:id` | Admin  | Delete a parking zone            |

### 🍴 Reservations
| Method   | Endpoint                        | Access        | Description             |
| -------- | ------------------------------- | ------------- | ----------------------- |
| `POST`   | `/reservations`                 | Authenticated | Reserve a parking spot  |
| `GET`    | `/reservations/my-reservations` | Authenticated | View your reservations  |
| `DELETE` | `/reservations/:id`             | Authenticated | Cancel your reservation |
| `GET`    | `/reservations`                 | Admin         | View all reservations   |

### 🩹 Health Check
| Method | Endpoint  | Access | Description      |
| ------ | --------- | ------ | ---------------- |
| `GET`  | `/health` | Public | API health check |

## 🏛️ Architecture

### Clean Architecture Layers

```
┌──────────────────────────────────────────────────────────────┐
│                        CLIENT LAYER                          │
│                (Postman, Frontend, Mobile App)               │
└──────────────────────────────────────────────────────────────┘
                              │
                              │ 
                              ▼
                          HTTP Request
┌──────────────────────────────────────────────────────────────┐
│                       HANDLER LAYER                          │
│ • Parses HTTP requests                                       │
│ • Validates DTOs                                             │
│ • Extracts JWT claims                                        │
│ • Returns JSON responses                                     │
│ • NO business logic or database access                       │
└──────────────────────────────────────────────────────────────┘
                              │
                              │ 
                              ▼
                        Calls Service
┌──────────────────────────────────────────────────────────────┐
│                       SERVICE LAYER                          │
│ • Business logic                                             │
│ • Password hashing (bcrypt)                                  │
│ • JWT generation                                             │
│ • Capacity / availability checks                             │
│ • Enforces business rules                                    │
│ • NO HTTP or database code                                   │
└──────────────────────────────────────────────────────────────┘
                              │
                              │ 
                              ▼
                        Calls Repository
┌──────────────────────────────────────────────────────────────┐
│                      REPOSITORY LAYER                        │
│ • Database CRUD operations                                   │
│ • Transactions & row locks                                   │
│ • GORM queries                                               │
│ • NO business logic                                          │
└──────────────────────────────────────────────────────────────┘
                              │
                              ▼

┌──────────────────────────────────────────────────────────────┐
│                         DATABASE                             │
│                    PostgreSQL (NeonDB)                       │
│          Tables: users, parking_zones, reservations          │
└──────────────────────────────────────────────────────────────┘
```

## 🔒 Concurrency Safety

### The "EV Spot Bottleneck" Solution

#### ❌ Problem: Race Condition

```text
Driver A reads: 19/20 spots occupied
Driver B reads: 19/20 spots occupied

Driver A creates a reservation ✅
Driver B creates a reservation ✅

Result:
21 cars are booked into a parking zone with only 20 spots. ❌
```

#### ✅ Solution: Row-Level Locking (`SELECT ... FOR UPDATE`)

```text
Driver A
────────
BEGIN TRANSACTION
↓
Lock Parking Zone #5 (FOR UPDATE)
↓
Active reservations = 19
↓
19 < 20 → Create reservation ✅
↓
COMMIT (Releases the lock)

                    │
                    ▼

Driver B
────────
Attempts to lock Parking Zone #5
↓
Waits until Driver A commits... ⏳
↓
Lock acquired
↓
Active reservations = 20
↓
20 == 20 → Reservation rejected ❌
↓
COMMIT
```

**Result:** The parking zone never exceeds its capacity, preventing overbooking even when multiple users make reservations simultaneously.

---

## 📁 Project Structure

```
spotsync/
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── config/
│   │   ├── config.go
│   │   └── database.go
│   ├── dto/
│   │   ├── auth_dto.go
│   │   ├── zone_dto.go
│   │   └── reservation_dto.go
│   ├── handler/
│   │   ├── auth_handler.go
│   │   ├── zone_handler.go
│   │   └── reservation_handler.go
│   ├── middleware/
│   │   ├── auth.go
│   │   ├── role.go
│   │   ├── cors.go
│   │   └── error_handler.go
│   ├── models/
│   │   ├── user.go
│   │   ├── parking_zone.go
│   │   └── reservation.go
│   ├── repository/
│   │   ├── user_repository.go
│   │   ├── zone_repository.go
│   │   └── reservation_repository.go
│   ├── router/
│   │   └── router.go
│   ├── service/
│   │   ├── auth_service.go
│   │   ├── zone_service.go
│   │   └── reservation_service.go
│   └── utils/
│       ├── response.go
│       └── errors.go
```

## 👤 Author
**Rahima Akter**<br/>
GitHub: [@Rahima-Akter](https://github.com/Rahima-Akter)

<br/>
Built with ❤️ in Go