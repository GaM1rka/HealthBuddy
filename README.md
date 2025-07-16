# HealthBuddy
> Web application where people share sports & wellness goals and achievements in a public feed, support each other with comments, and manage personal profiles.

---

## 1. Project Overview & Setup (1 pt)

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

## 2. Screenshots & GIFs
| Feature     | Preview                    |
|-------------|----------------------------|
| Onboarding  | docs/img/onboarding.png    |
| Create Goal | docs/img/create_goal.gif   |
| Comments    | docs/img/comments.png      |
| Profile     | docs/img/profile.gif       |

## 3. API Documentation (1 pt)
Interactive Swagger UI → http://localhost:8080/docs
Static reference → docs/api_documentation.md

## 4. Architecture Diagrams & Explanations

<img width="1429" height="699" alt="image" src="https://github.com/user-attachments/assets/c47e17c3-7379-4da6-8b5a-a62f086295e3" />

## Folder Layout

health-buddy/                 # repo root
├── backend/
│   ├── services/
│   │   ├── auth_service/     # Auth micro-service (Dockerfile present)
│   │   ├── feed_service/     # Feed micro-service  (Dockerfile present)
│   │   └── profile_service/  # Profile micro-service (Dockerfile present)
│   └── gateway_service/      # Gateway (CORS & structure recently updated)
├── front-end/                # Flutter application
│   ├── pubspec.yaml
│   └── .gitignore
├── docker-compose.yml        # Orchestrates all services
└── README.md                 # Initial version
