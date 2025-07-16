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

flowchart TD
    subgraph Client
        A[Flutter App] -->|HTTPS / WebSocket| B(Gateway)
    end
    subgraph Back-end
        B --> C[Go Auth Service]
        B --> D[Go Feed Service]
        B --> E[Go Profile Service]
    end
    subgraph Infra
        C --> F[(PostgreSQL Auth)]
        D --> G[(PostgreSQL Feed)]
        E --> H[(PostgreSQL Profile)]
        B --> I[Redis Cache]
    end
