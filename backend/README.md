# Lume API

API RESTful de alta performance para a **Lume**, uma plataforma social de avaliação e diário de livros.

- **Linguagem:** Go 1.21+ (desenvolvida e validada com Go 1.22)
- **Framework:** [Fiber](https://gofiber.io/) v2
- **Banco:** PostgreSQL (pgx nativo — `pgxpool`, sem ORM) — compatível com Supabase
- **Autenticação:** JWT (`golang-jwt/jwt/v5`) com middleware do Fiber
- **Integração externa:** Open Library (padrão, sem chave) com fallback para Google Books API (chave opcional)
- **Contratos de API:** modelados pelo frontend React anexado (`App.tsx` — app "Lume")

---

## Arquitetura

```
lume_api/
├── cmd/
│   ├── api/                  # Bootstrap do servidor (config, pool, migrações, worker pool, graceful shutdown)
│   └── migrate/              # CLI de migrações (make migrate-up / migrate-down)
├── migrations/               # Migrações SQL versionadas (0001..0008 .up/.down)
├── internal/
│   ├── config/               # Configuração via variáveis de ambiente (.env)
│   ├── database/             # pgxpool + runner de migrações (transacional, versionado)
│   ├── models/               # Entidades e DTOs (User, Book, Log, Shelf, List, Social, Comment…)
│   ├── dto/                  # Contratos de request/response
│   ├── middleware/           # Auth JWT, request-id, error handler padronizado, logging (slog)
│   ├── repositories/         # Camada de persistência (users, books, logs, shelves, lists, feed, stats)
│   ├── services/             # Regras de negócio, busca local-first + fallback externo
│   ├── handlers/             # Handlers HTTP (Feed, Perfil, Livro, Diário, Listas…)
│   ├── router/               # Rotas /api/v1 com CORS e rate limiting
│   ├── async/                # Worker pool (goroutines + channel) para recálculo de estatísticas
│   └── utils/                # Cursor pagination, JWT, bcrypt, respostas HTTP
├── .env.example
├── docker-compose.yml        # PostgreSQL local de desenvolvimento
└── Makefile
```

### Decisões técnicas

| Decisão | Escolha | Motivo |
|---|---|---|
| Acesso a dados | pgxpool nativo | Zero overhead de ORM; tipos nativos `[]text`, `jsonb` |
| IDs | UUID (`gen_random_uuid()`) | Únicos, gerados no banco, seguros para expor na API |
| Migrations | SQL versionado, runner próprio | Transacional (up+registro atômico), roda no boot e via CLI |
| Busca de livros | Local-first → Open Library → Google Books | Menor latência; o banco local vira cache do catálogo |
| Feed | Cursor opaco (`created_at+id` encodeado) | Paginação estável e `O(1)` por página, sem OFFSET |
| Recálculo de stats | Assíncrono (worker pool + canal) | Escrita rápida; stats anuais cacheadas em `user_year_stats` (jsonb) |
| Erros | Formato único `{status, message, error}` | Consumo previsível no frontend |
| Senhas | bcrypt | Hash com custo configurável |

## Entidades

`users`, `books`, `logs` (diário/avaliação, rating 0.5–5.0 em passos de 0.5), `shelves` (tbr/reading/read/abandoned), `lists`, `list_items` (ordenação por `position`), `social` (seguidores), `log_likes`, `comments`, `user_year_stats` (cache de estatísticas).

## Endpoints principais (`/api/v1`)

| Área | Rota |
|---|---|
| Auth | `POST /auth/register` · `POST /auth/login` · `GET /auth/me` |
| Perfil | `GET /users/:username` · `PUT /users/me` · `GET /users/:username/logs` · `GET /users/:username/shelf` · `GET /users/:username/stats?year=` |
| Social | `POST/DELETE /users/:id/follow` · `GET /users/:id/followers` · `GET /users/:id/following` |
| Livros | `GET /books/search?q=` · `GET /books/:id` (média, histograma, resenhas) · `GET /books/popular` |
| Diário | `POST /logs` · `GET/PUT/DELETE /logs/:id` · `POST/DELETE /logs/:id/like` · `POST /logs/:id/comments` · `GET /logs/:id/comments` |
| Estante | `PUT /shelf` |
| Listas | `POST/GET /lists` · `GET/PUT/DELETE /lists/:id` · `POST/DELETE /lists/:id/items` · `GET /lists/trending` |
| Feed | `GET /feed?cursor=&limit=` (seguidos) · `GET /feed/global` · `GET /feed/reviews` |
| Saúde | `GET /healthz` |

## Como executar

```bash
# 1. Banco local
cp .env.example .env
docker compose up -d postgres        # ou aponte DATABASE_URL para o Supabase

# 2. Migrações (aplicadas automaticamente no boot, ou manualmente)
make migrate-up

# 3. Rodar
make run                             # http://localhost:3000
```

**Smoke test rápido** (após subir e migrar):

```bash
curl -X POST localhost:3000/api/v1/auth/register -H 'Content-Type: application/json' \
  -d '{"username":"clara","email":"clara@mail.com","password":"secret123"}'

curl -X POST localhost:3000/api/v1/auth/login -H 'Content-Type: application/json' \
  -d '{"email":"clara@mail.com","password":"secret123"}'          # → {token}

curl 'localhost:3000/api/v1/books/search?q=torto%20arado'          # busca local-first + persistência

TOKEN=<seu-token>
curl localhost:3000/api/v1/feed -H "Authorization: Bearer $TOKEN"  # feed com cursor na resposta
```

## Produção (Supabase)

- Use a *connection string* do painel do Supabase com `sslmode=require`.
- `JWT_SECRET`: gere com `openssl rand -hex 32` e mantenha no seu gerenciador de segredos.
- As migrações são compatíveis com Supabase (extensão `pg_trgm` disponível na plataforma).

> **Nota:** o smoke test de ponta a ponta (banco + migrações + fluxo curl) não foi executado no ambiente de build por indisponibilidade do daemon Docker/servidor PostgreSQL local; o projeto compila e o `go vet` passam sem erros. Execute os passos acima na sua máquina para validar o fluxo completo.
