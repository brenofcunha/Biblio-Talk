# Biblio-Talk

Uma rede social para leitores compartilharem o que estão lendo, registrarem
leituras, avaliarem livros e descobrirem novas recomendações.

Nos arquivos de referência, a aplicação aparece com o nome **Lume**. O produto
combina uma API REST em Go com uma interface web em React, criando um espaço
para diário de leitura, estantes, listas editoriais e interação entre leitores.

## Visão geral

O Biblio-Talk foi pensado para acompanhar todo o ciclo de leitura:

- pesquisar livros e consultar informações do catálogo;
- registrar uma leitura com data, nota, resenha, tags e indicação de spoiler;
- organizar livros em estantes como “quero ler”, “lendo”, “lido” e
  “abandonado”;
- criar e acompanhar listas de livros;
- seguir leitores, curtir resenhas e comentar publicações;
- visualizar feed social, perfil e estatísticas anuais de leitura.

## Componentes

### API

A API está implementada em Go 1.26 e utiliza:

- [Fiber](https://gofiber.io/) para HTTP;
- PostgreSQL 16 como banco de dados;
- `pgxpool`, sem ORM;
- JWT e bcrypt para autenticação;
- migrações SQL versionadas;
- Open Library como fonte principal de livros e Google Books como fallback
  opcional;
- worker assíncrono para atualização de estatísticas.

O servidor aplica as migrações automaticamente durante o boot. Também existe
uma CLI para executá-las manualmente.

### Interface web

A interface anexada é um protótipo React/TypeScript com estilo visual escuro,
layout responsivo e navegação entre:

- feed;
- detalhe de livro;
- perfil;
- listas;
- modal de registro de leitura.

No estado atual, a interface usa dados de demonstração definidos no próprio
`App.tsx` e ainda não está conectada à API. A integração entre as duas camadas
deve substituir esses dados estáticos por chamadas aos endpoints descritos
abaixo.

## Arquitetura da API

```text
cmd/
├── api/                 # Inicialização do servidor
└── migrate/             # CLI de migrações

internal/
├── async/               # Worker de estatísticas
├── config/              # Variáveis de ambiente
├── database/            # Pool PostgreSQL e migrações
├── dto/                 # Contratos de entrada e saída
├── handlers/            # Handlers HTTP
├── middleware/          # JWT, CORS, logs e tratamento de erros
├── models/              # Entidades do domínio
├── repositories/        # Persistência
├── router/              # Rotas da API
├── services/            # Regras de negócio
└── utils/               # JWT, bcrypt, paginação e respostas HTTP

migrations/              # Migrações SQL up/down
```

## Principais endpoints

Todas as rotas, exceto `/healthz`, usam o prefixo `/api/v1`.

| Domínio | Endpoints |
| --- | --- |
| Saúde | `GET /healthz` |
| Autenticação | `POST /auth/register`, `POST /auth/login`, `GET /auth/me` |
| Usuários | `GET /users/:username`, `PUT /users/me` |
| Social | `POST/DELETE /users/:id/follow`, seguidores e seguindo |
| Livros | `GET /books/search`, `GET /books/:id`, `GET /books/popular` |
| Diário | `POST /logs`, `GET/PUT/DELETE /logs/:id` |
| Interação | likes e comentários em `/logs/:id` |
| Estantes | `PUT /shelf` |
| Listas | CRUD em `/lists` e gerenciamento de itens |
| Feed | `GET /feed`, `/feed/global` e `/feed/reviews` |

As respostas seguem um formato padronizado com status, mensagem e dados. As
rotas protegidas exigem o header:

```text
Authorization: Bearer <token>
```

## Pré-requisitos

Para executar a API localmente, instale:

- Go 1.26 ou superior;
- Docker e Docker Compose;
- Git;
- `make` (opcional; os comandos também podem ser executados diretamente com
  `go run`).

Para executar a interface, além dos arquivos React, será necessário um projeto
Vite com Node.js e npm. O anexo atual não inclui `package.json`; portanto, a
interface precisa ser colocada em um projeto React/Vite antes de ser iniciada.

## Inicialização local da API

### 1. Obter o código

```bash
git clone https://github.com/brenofcunha/Biblio-Talk.git
cd Biblio-Talk
```

O código da API de referência está no diretório `lume_api`. Se ele ainda estiver
fora deste repositório, entre no diretório onde a API foi salva:

```bash
cd caminho/para/lume_api
```

### 2. Configurar as variáveis de ambiente

Crie o arquivo `.env` a partir do exemplo:

```bash
cp .env.example .env
```

Para Windows PowerShell, use:

```powershell
Copy-Item .env.example .env
```

A configuração local padrão é:

```dotenv
PORT=3000
DATABASE_URL=postgres://postgres:postgres@localhost:5432/lume?sslmode=disable
JWT_SECRET=troque-este-segredo-em-producao
JWT_EXPIRATION=24h
CORS_ALLOWED_ORIGIN=http://localhost:5173
```

Open Library funciona sem chave. `GOOGLE_BOOKS_API_KEY` é opcional e só deve
ser preenchida se o fallback do Google Books for necessário.

### 3. Subir o PostgreSQL

```bash
docker compose up -d postgres
```

Verifique se o container está saudável:

```bash
docker compose ps
```

### 4. Aplicar as migrações

As migrações são aplicadas automaticamente quando a API inicia. Para executá-las
manualmente:

```bash
make migrate-up
```

Sem `make`:

```bash
go run ./cmd/migrate -direction=up
```

### 5. Iniciar a API

```bash
make run
```

Sem `make`:

```bash
go run ./cmd/api
```

A API ficará disponível em `http://localhost:3000`.

### 6. Validar a instalação

Teste o health check:

```bash
curl http://localhost:3000/healthz
```

Crie um usuário e faça login:

```bash
curl -X POST http://localhost:3000/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"clara","email":"clara@example.com","password":"secret123"}'

curl -X POST http://localhost:3000/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"clara@example.com","password":"secret123"}'
```

Para testar a busca de livros:

```bash
curl "http://localhost:3000/api/v1/books/search?q=torto%20arado"
```

## Inicialização da interface

Os arquivos anexados (`App.tsx`, `main.tsx`, `index.css` e `vite-env.d.ts`)
devem ser copiados para um projeto React + TypeScript criado com Vite.

Exemplo:

```bash
npm create vite@latest frontend -- --template react-ts
cd frontend
npm install
npm install tailwindcss @tailwindcss/vite
```

Depois, substitua os arquivos `src/App.tsx`, `src/main.tsx` e
`src/index.css` pelos arquivos da interface anexada. Configure o plugin do
Tailwind conforme a versão instalada e inicie o servidor:

```bash
npm run dev
```

A interface ficará disponível, normalmente, em
`http://localhost:5173`. Essa origem já está liberada pelo CORS padrão da API.

> A interface anexada é um protótipo visual. O botão de registro, os livros,
> usuários, avaliações e navegação usam estado local e dados fixos. Para
> transformá-la em uma aplicação completa, será necessário adicionar um cliente
> HTTP, autenticação JWT e os estados de carregamento/erro para consumir a API.

## Desenvolvimento

Comandos disponíveis na API:

```bash
make test          # executa os testes
make vet           # verifica problemas do código Go
make format        # formata os arquivos
make lint          # executa golangci-lint, se instalado
make build         # gera bin/lume-api
make dev-db-down   # encerra o PostgreSQL local
```

Para desligar a API, use `Ctrl+C`. Para remover o container e a rede do banco:

```bash
docker compose down
```

Para remover também os dados persistidos do PostgreSQL, use o comando abaixo
com atenção:

```bash
docker compose down -v
```

## Segurança e produção

- Nunca use o `JWT_SECRET` de desenvolvimento em produção.
- Armazene segredos em um gerenciador de secrets.
- Use uma connection string PostgreSQL com `sslmode=require` em ambientes
  externos, como Supabase.
- Restrinja `CORS_ALLOWED_ORIGIN` ao domínio real da interface.
- Não versione `.env` nem chaves de APIs.

## Licença

Este projeto está sob a licença MIT. Consulte [LICENSE](./LICENSE).
