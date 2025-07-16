# HealthBuddy
> Web application where people share sports & wellness goals and achievements in a public feed, support each other with comments, and manage personal profiles.

---

## 1. Project Overview & Setup

| Layer | Tech | Responsibility |
|-------|------|----------------|
| Client | Flutter 3.19 (Web / Mobile) | UI, Riverpod state |
| Gateway | Nginx + Go | TLS, CORS, rate-limit |
| Auth Service | Go 1.22 + Gin | JWT, register / login |
| Feed Service | Go 1.22 + Gin | Posts & comments |
| Profile Service | Go 1.22 + Gin | User bio & avatars |
| Data | PostgreSQL 15 | per-service DBs |
| Cache | Redis 7 | sessions & hot posts |

### Quick Start
```bash
git clone https://github.com/your-org/health-buddy.git
cd health-buddy
cp .env.example .env          # fill secrets
docker compose up --build     # backend at http://localhost:8080
# Flutter
cd flutter_app && flutter run -d chrome   # or iOS / Android
```

## 2. Features presentation

### Registration
![Screencastfrom2025-07-1612-55-42-ezgif com-video-to-gif-converter](https://github.com/user-attachments/assets/1bc24f2b-b93d-42c7-ad79-45939b20c79b)

### Open profile and edit bio

![Screencastfrom2025-07-1612-55-42-ezgif com-video-to-gif-converter (1)](https://github.com/user-attachments/assets/ec858441-0e44-42f7-b27e-1adb09a9c32f)

### Write new post

![Screencastfrom2025-07-1612-55-42-ezgif com-video-to-gif-converter (2)](https://github.com/user-attachments/assets/eb717271-cf2a-4cee-b741-7838bb7465cd)

## 3. API Documentation

# HealthBuddy API Reference

## Base URLs
- **Gateway:**  https://api.healthbuddy.app
- **Local:**    http://localhost:8080

## Global Headers
- **Content-Type:**  application/json
- **Authorization:** Bearer `<jwt>`

---

## AUTH SERVICE (/auth)

### POST /auth/register
- **Body:** `{ username, email, password }`
- **Responses:**
  - `201`: `{ token }`
  - `400`: Invalid input
  - `500`: Server error

### POST /auth/login
- **Body:** `{ username, password }`
- **Responses:**
  - `200`: `{ token }`
  - `401`: Wrong credentials

### GET /auth/users/{id}
- **Headers:** Authorization
- **Responses:**
  - `200`: `{ id, username, email, created_at }`
  - `404`: Not found

### DELETE /auth/users/{id}
- **Responses:**
  - `204`: No content
  - `404`: Not found

---

## FEED SERVICE (/feed) • JWT required via X-User-ID

### GET /feed/health
- **Responses:**
  - `200`: `{ status: "ok" }`
  - `503`: `{ status: "down" }`

### POST /feed/publications
- **Body:** `{ title (≤300), content (≤10 000) }`
- **Responses:**
  - `201`: `{ post_id, user_id, title, content, created_at }`

### GET /feed/publications
- **Responses:**
  - `200`: `[ PublicationResponse… ]` (newest first)

### GET /feed/publications/{id}
- **Responses:**
  - `200`: `PublicationResponse`
  - `404`: Not found

### PUT /feed/publications/{id}
- **Body:** `{ title, content }`
- **Responses:**
  - `200`: Updated object
  - `403`: Forbidden
  - `404`: Not found

### DELETE /feed/publications/{id}
- **Responses:**
  - `204`: No content
  - `403`: Forbidden
  - `404`: Not found

### GET /feed/users/{userID}/publications
- **Responses:**
  - `200`: `[ PublicationResponse… ]`

### Comments

#### POST /feed/comments
- **Body:** `{ post_id, content (≤10 000) }`
- **Responses:**
  - `201`: `{ comment_id, user_id, content, created_at }`

#### GET /feed/comments?post_id={postID}
- **Responses:**
  - `200`: `[ CommentResponse… ]`
  - `400`: Missing param

#### GET /feed/comments/{id}
- **Responses:**
  - `200`: `CommentResponse`
  - `404`: Not found

#### PUT /feed/comments/{id}
- **Body:** `{ content }`
- **Responses:**
  - `200`: Updated object
  - `403`: Forbidden
  - `404`: Not found

#### DELETE /feed/comments/{id}
- **Responses:**
  - `204`: No content
  - `403`: Forbidden
  - `404`: Not found

---

## PROFILE SERVICE (/profile) • JWT required

### GET /profile/health
- **Responses:**
  - `200`: `{ status: "ok" }`

### POST /profile
- **Headers:** X-User-ID, Content-Type
- **Body:** `{ name, bio?, avatar_url? }`
- **Responses:**
  - `201`: ProfileResponse (without posts)

### GET /profile
- **Responses:**
  - `200`: `{ user_id, name, bio, avatar, created_at, posts: [PublicationResponse…] }`

### PUT /profile
- **Body:** any subset of `{ name, bio, avatar_url }`
- **Responses:**
  - `200`: Updated ProfileResponse

### DELETE /profile
- **Responses:**
  - `204`: No content (cascades to Auth deletion)

---

## Status Codes & Messages

- `200 OK`: Request succeeded
- `201 Created`: Resource created
- `204 No Content`: Deletion succeeded
- `400 Bad Request`: Validation / JSON error
- `401 Unauthorized`: Invalid or missing JWT
- `403 Forbidden`: Not owner
- `404 Not Found`: Resource does not exist
- `500 Internal`: Server or DB failure

---

*Generated from Postman collections*


## 4. Architecture Diagrams & Explanations

<img width="1429" height="699" alt="image" src="https://github.com/user-attachments/assets/c47e17c3-7379-4da6-8b5a-a62f086295e3" />

## Folder Layout
```
health-buddy/                 # repo root
├── backend/
│   ├── services/
│   │   ├── auth_service/     # Auth micro-service
│   │   ├── feed_service/     # Feed micro-service
│   │   └── profile_service/  # Profile micro-service
│   └── gateway_service/      # Gateway
├── front-end/                # Flutter application
│   ├── pubspec.yaml
│   └── .gitignore
├── docker-compose.yml        # Orchestrates all services
└── README.md                 # Initial version
```

## Build & Run

