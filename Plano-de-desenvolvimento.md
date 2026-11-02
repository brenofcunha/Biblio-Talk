# Plano de desenvolvimento

## Objetivo

Transformar o protótipo do Biblio-Talk/Lume em uma aplicação funcional para
leitores, com frontend React, API em Go, PostgreSQL e publicação em nuvem.

## Diagnóstico atual

| Área | Situação |
| --- | --- |
| API | Implementada em Go/Fiber, com PostgreSQL, JWT, migrações SQL e endpoints de usuários, livros, diário, estantes, listas e feed. |
| Feed | API implementada com paginação por cursor. O carregamento assíncrono ainda precisa ser integrado ao frontend. |
| Banco | PostgreSQL local disponível via Docker Compose. |
| Supabase | Pode ser usado como PostgreSQL remoto via `DATABASE_URL`; Supabase Auth ainda não está integrado. |
| Frontend | Protótipo React/TypeScript com dados estáticos e sem integração com a API. |
| Deploy | Ainda não há configuração para Vercel, Railway, CI/CD ou provisionamento de produção. |
| Testes | `go test ./...` executa, mas ainda faltam testes automatizados de comportamento e integração. |
| Estrutura | A API de referência está em `C:\Users\Breno\Downloads\lume_api`, fora do repositório principal. |

## Arquitetura alvo do MVP

- **Frontend:** React + Vite publicado na Vercel.
- **Backend:** Go/Fiber publicado no Railway.
- **Banco:** PostgreSQL do Supabase.
- **Autenticação:** JWT gerado e validado pela API Go.
- **Comunicação:** API REST com paginação por cursor.
- **Processamento assíncrono:** worker interno para estatísticas.
- **Feed:** carregamento assíncrono no frontend, sem fila externa inicialmente.

O Supabase será usado inicialmente como banco PostgreSQL. A migração para
Supabase Auth deve ser considerada somente após o MVP, para evitar manter dois
modelos de autenticação simultaneamente.

## Etapas de desenvolvimento

### 1. Consolidar a estrutura do projeto

- [ ] Mover ou copiar a API para dentro do repositório principal.
- [ ] Criar a estrutura `backend/` e `frontend/`.
- [ ] Criar o projeto React/Vite com `package.json`.
- [ ] Adicionar `.gitignore` para `.env`, artefatos de build e dependências.
- [ ] Atualizar o README com os comandos da nova estrutura.

Estrutura esperada:

```text
Biblio-Talk/
├── frontend/
├── backend/
│   ├── cmd/
│   ├── internal/
│   ├── migrations/
│   ├── go.mod
│   └── Dockerfile
├── README.md
└── Plano-de-desenvolvimento.md
```

### 2. Integrar o frontend com a API

- [ ] Criar cliente HTTP centralizado.
- [ ] Configurar `VITE_API_URL`.
- [ ] Implementar armazenamento e envio do JWT.
- [ ] Criar telas de login e cadastro.
- [ ] Adicionar estados de carregamento, erro e conteúdo vazio.
- [ ] Integrar livros, perfis, listas, estantes e diário.
- [ ] Substituir os dados fixos do protótipo por respostas da API.

### 3. Implementar o carregamento assíncrono do feed

A API já disponibiliza:

- `GET /api/v1/feed`;
- `GET /api/v1/feed/global`;
- `GET /api/v1/feed/reviews`;
- paginação por cursor no feed de usuários seguidos.

Implementação prevista:

- [ ] Solicitar o feed após a autenticação.
- [ ] Exibir skeleton durante o carregamento.
- [ ] Tratar erros de rede e permitir nova tentativa.
- [ ] Implementar “carregar mais” usando o cursor retornado.
- [ ] Implementar atualização manual do feed.
- [ ] Avaliar polling, WebSocket ou SSE somente se houver necessidade de tempo
      real.

O worker existente deve continuar responsável por estatísticas. Ele não deve
ser usado para atrasar a resposta de leitura do feed.

### 4. Definir e consolidar a autenticação

Estratégia recomendada para o MVP:

- [ ] Manter cadastro e login na API Go.
- [ ] Manter JWT próprio.
- [ ] Validar o token pelo middleware existente.
- [ ] Revisar expiração, revogação e armazenamento seguro no frontend.
- [ ] Padronizar respostas de erro de autenticação.

Após o MVP, avaliar Supabase Auth como uma migração separada. Não misturar JWT
próprio e Supabase Auth na mesma etapa.

### 5. Preparar o Supabase

- [ ] Criar o projeto PostgreSQL no Supabase.
- [ ] Configurar `DATABASE_URL` com `sslmode=require`.
- [ ] Executar e validar todas as migrações.
- [ ] Confirmar extensões PostgreSQL necessárias.
- [ ] Configurar backups e retenção.
- [ ] Remover credenciais fixas de desenvolvimento da configuração de produção.
- [ ] Configurar usuários e permissões de banco com menor privilégio.

### 6. Preparar o deploy da API no Railway

- [ ] Criar `Dockerfile` para o backend.
- [ ] Usar a porta definida pela variável `PORT`.
- [ ] Configurar o health check em `/healthz`.
- [ ] Configurar graceful shutdown.
- [ ] Definir as variáveis de ambiente no Railway:

```text
PORT
DATABASE_URL
JWT_SECRET
JWT_EXPIRATION
CORS_ALLOWED_ORIGIN
OPEN_LIBRARY_BASE_URL
GOOGLE_BOOKS_BASE_URL
GOOGLE_BOOKS_API_KEY
```

- [ ] Definir uma estratégia segura para executar migrações no deploy.
- [ ] Configurar logs e monitoramento básico.
- [ ] Validar o fluxo de cadastro, login, busca e feed em ambiente publicado.

O `docker-compose.yml` deve permanecer voltado ao desenvolvimento local.

### 7. Preparar o deploy do frontend na Vercel

- [ ] Configurar o projeto frontend na Vercel.
- [ ] Definir `VITE_API_URL` apontando para a API no Railway.
- [ ] Configurar domínio permitido no CORS da API.
- [ ] Garantir fallback das rotas do frontend após refresh.
- [ ] Configurar variáveis de ambiente de produção.
- [ ] Validar carregamento das imagens externas e recursos estáticos.
- [ ] Configurar build e deploy automático a partir da branch `develop` ou de
      uma branch de produção.

### 8. Criar testes e qualidade

- [ ] Testar autenticação e autorização.
- [ ] Testar validação dos DTOs.
- [ ] Testar handlers principais.
- [ ] Testar paginação por cursor do feed.
- [ ] Testar execução das migrações.
- [ ] Criar testes de integração com PostgreSQL.
- [ ] Criar smoke test de cadastro, login, busca e feed.
- [ ] Adicionar build, testes e lint ao CI.

## Ordem recomendada para o MVP

1. Consolidar frontend e backend no mesmo repositório.
2. Criar o projeto React/Vite.
3. Integrar login e cadastro.
4. Integrar busca de livros.
5. Integrar criação e visualização do diário.
6. Implementar carregamento assíncrono do feed.
7. Adicionar testes de integração.
8. Migrar o banco local para Supabase.
9. Publicar a API no Railway.
10. Publicar o frontend na Vercel.
11. Configurar domínio, CORS, logs e monitoramento.

## Critérios de conclusão do MVP

- [ ] Usuário consegue criar conta e fazer login.
- [ ] Usuário consegue pesquisar livros.
- [ ] Usuário consegue registrar uma leitura.
- [ ] Usuário consegue visualizar e paginar o feed.
- [ ] Usuário consegue interagir com resenhas.
- [ ] Frontend e API estão publicados e comunicando-se por HTTPS.
- [ ] Banco de produção está no Supabase.
- [ ] Health check e logs estão disponíveis.
- [ ] Fluxos principais possuem testes automatizados.
- [ ] Nenhum segredo está versionado.

## Riscos e cuidados

- Não versionar arquivos `.env` ou chaves de APIs.
- Não usar o segredo JWT de desenvolvimento em produção.
- Não executar migrações destrutivas automaticamente no deploy.
- Não introduzir WebSocket ou filas antes de medir a necessidade.
- Garantir índices e limites de paginação para evitar consultas caras no feed.
- Monitorar limites e disponibilidade das APIs externas de livros.