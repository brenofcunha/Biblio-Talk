# Contexto para configurar o Supabase — Biblio-Talk

## Objetivo deste arquivo

Este documento fornece ao Gemini o contexto técnico atual do Biblio-Talk para
auxiliar na execução da **Etapa 5 — Preparar o Supabase**. A aplicação também
aparece com o nome visual **Lume** em alguns arquivos e telas.

O objetivo da Etapa 5 é usar o PostgreSQL gerenciado do Supabase como banco da
API Go/Fiber, sem migrar a autenticação para o Supabase Auth neste momento.

## Estado atual do projeto

- **Etapas 1 a 4:** concluídas.
- **Etapa 5:** ainda pendente.
- **Etapa 6:** deploy da API no Railway ainda pendente.
- **Etapa 7:** deploy do frontend na Vercel ainda pendente.
- **Etapa 8:** testes automatizados e CI ainda pendentes.

O plano oficial está em [Plano-de-desenvolvimento.md](./Plano-de-desenvolvimento.md).

## Arquitetura atual

```text
frontend/  React + Vite
    |
    | REST + Authorization: Bearer <JWT>
    v
backend/   Go + Fiber
    |
    | pgxpool / SQL / migrações automáticas
    v
PostgreSQL local ou Supabase
```

### Componentes

- Frontend: React, Vite e TypeScript.
- Backend: Go, Fiber e `pgxpool`; não há ORM.
- Banco: PostgreSQL 16 no ambiente local.
- Autenticação: JWT próprio da API e bcrypt para senhas.
- Livros: Open Library como fonte principal e Google Books como fallback
  opcional.
- Estatísticas: worker interno que atualiza `user_year_stats`.
- Feed: consultas SQL com paginação por cursor.

O Supabase será usado inicialmente apenas como PostgreSQL. Não adicionar
Supabase Auth, triggers de autenticação ou dependência do SDK Supabase sem uma
decisão explícita para uma etapa futura.

## Configuração atual do banco

### Desenvolvimento local

O banco local é definido em [docker-compose.yml](./backend/docker-compose.yml):

```text
Host: localhost
Porta: 5432
Banco: lume
Usuário: postgres
Senha: postgres
SSL: desabilitado apenas localmente
```

### Produção no Supabase

A variável usada pela API é:

```text
DATABASE_URL=<connection-string-do-Supabase>
```

Para produção, a conexão deve usar `sslmode=require`. Não commitar a URL,
senha, token, chave de API ou qualquer outro segredo no repositório.

O arquivo de referência das variáveis é
[backend/.env.example](./backend/.env.example). Em produção, usar as
credenciais geradas no painel do Supabase e armazená-las no gerenciador de
segredos da plataforma de deploy.

## Migrações existentes

As migrações versionadas estão em [backend/migrations](./backend/migrations).
Elas são executadas automaticamente durante o boot da API e também podem ser
executadas pela CLI em `backend/cmd/migrate`.

Ordem atual:

1. `0001_create_users`
2. `0002_create_books`
3. `0003_create_logs`
4. `0004_create_shelves`
5. `0005_create_lists`
6. `0006_create_social`
7. `0007_create_likes_comments`
8. `0008_create_stats_and_indexes`

As migrações criam as seguintes estruturas principais:

- `users`: contas, senha bcrypt, perfil e quatro livros favoritos.
- `books`: catálogo local, IDs externos, autores, ISBN, editora, páginas,
  capas e gêneros.
- `logs`: diário de leitura, nota, resenha, data, spoiler e contagem de likes.
- `shelves`: estantes `tbr`, `reading`, `read` e `abandoned`.
- `lists` e `list_items`: listas pessoais e seus livros.
- `social`: relacionamentos de seguidores.
- `log_likes` e `comments`: interação social.
- `user_year_stats`: cache JSON de estatísticas anuais.

A migração `0008` habilita `pg_trgm` para busca por similaridade no título dos
livros e cria índices de feed, diário, livros, listas e comentários. Confirmar
que a extensão está disponível no projeto Supabase antes de considerar a etapa
concluída.

## Rotas relevantes da API

Todas as rotas usam o prefixo `/api/v1`, exceto `/healthz`.

```text
GET  /healthz

POST /auth/register
POST /auth/login
GET  /auth/me
POST /auth/logout

GET  /books/search
GET  /books/:id
GET  /books/popular

POST /logs
GET  /logs/:id
PUT  /logs/:id
DELETE /logs/:id

PUT  /shelf
CRUD /lists
GET  /feed
GET  /feed/global
GET  /feed/reviews
```

O banco não é acessado diretamente pelo frontend. Apenas a API Go deve
possuir a URL de conexão e as credenciais do PostgreSQL.

## Checklist da Etapa 5

Executar e registrar evidências para cada item:

- [ ] Criar o projeto PostgreSQL no Supabase.
- [ ] Configurar a `DATABASE_URL` de produção com `sslmode=require`.
- [ ] Executar as oito migrações em ordem.
- [ ] Validar tabelas, índices, chaves estrangeiras e constraints.
- [ ] Confirmar a extensão `pg_trgm`.
- [ ] Testar conexão e `GET /healthz` usando o banco Supabase.
- [ ] Confirmar que o worker consegue gravar estatísticas em
      `user_year_stats`.
- [ ] Configurar backups e retenção conforme o plano do Supabase.
- [ ] Criar uma credencial de banco com menor privilégio possível para a API.
- [ ] Remover qualquer credencial fixa de desenvolvimento da configuração de
      produção.
- [ ] Documentar o procedimento de migração e rollback sem expor segredos.

## Recomendações para a configuração

1. Criar o projeto Supabase e escolher a região mais próxima do futuro backend
   no Railway.
2. Usar a conexão PostgreSQL recomendada pelo Supabase para o backend. Validar
   se o modo escolhido é adequado para migrações e para `pgxpool`; não usar a
   chave pública `anon` como substituta de `DATABASE_URL`.
3. Executar primeiro as migrações em um projeto Supabase de validação, não no
   projeto de produção.
4. Confirmar o resultado com:

   ```powershell
   $env:DATABASE_URL = "<segredo-fornecido-localmente>"
   go run ./cmd/migrate
   go run ./cmd/api
   curl http://localhost:3000/healthz
   ```

5. Fazer um smoke test de cadastro, login, criação de livro/registro de
   leitura e consulta do feed.
6. Só depois configurar a mesma `DATABASE_URL` como segredo no Railway.

## Cuidados e decisões que o Gemini deve respeitar

- Não revelar, imprimir ou commitar senhas, tokens ou URLs completas com
  credenciais.
- Não alterar o modelo para Supabase Auth nesta etapa.
- Não remover as migrações SQL existentes sem comparar o schema resultante.
- Não executar `DROP DATABASE`, `DROP SCHEMA` ou comandos destrutivos no banco
  de produção.
- Não usar `sslmode=disable` em produção.
- Não marcar a Etapa 5 como concluída apenas por criar o projeto; é necessário
  validar migrações, extensão, conexão, health check e permissões.
- Se a conexão escolhida pelo Supabase não funcionar com `pgxpool` ou com o
  runner de migrações, explicar a causa e propor a alternativa de conexão,
  sem alterar o código às cegas.

## Pedido ao Gemini

Com base neste contexto, auxiliar a concluir a Etapa 5 de forma segura:

1. revisar as migrações e os requisitos do PostgreSQL;
2. orientar a criação e configuração do projeto Supabase;
3. indicar qual connection string deve ser usada em cada ambiente;
4. validar as migrações e a extensão `pg_trgm`;
5. orientar backups, retenção e menor privilégio;
6. preparar uma lista de evidências para atualizar
   [Plano-de-desenvolvimento.md](./Plano-de-desenvolvimento.md).

Antes de qualquer ação destrutiva ou alteração de schema, pedir confirmação
explícita e apresentar o impacto.
