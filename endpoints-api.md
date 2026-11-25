# Endpoints da API - estado atual e lacunas

Este documento compara a ideologia e o objetivo real descritos em
[`ideologia-aplica.md`](./ideologia-aplica.md) com os endpoints registrados em
[`backend/internal/router/router.go`](./backend/internal/router/router.go).

Prefixo padrão: `/api/v1`, exceto `/healthz`.

## Legenda

- **Implementado**: rota existente no backend.
- **Parcial**: rota existente, mas sem todos os campos, regras ou entidades necessários ao objetivo.
- **Criar**: rota ainda inexistente.

## Endpoints já implementados

### Saúde

| Método | Endpoint | Autenticação | Estado |
|---|---|---|---|
| `GET` | `/healthz` | Não | Implementado |

### Autenticação

| Método | Endpoint | Autenticação | Estado |
|---|---|---|---|
| `POST` | `/auth/register` | Não | Implementado; exige confirmação de e-mail |
| `POST` | `/auth/login` | Não | Implementado; credenciais inválidas retornam `401` |
| `POST` | `/auth/refresh` | Cookie HttpOnly | Implementado; rotaciona o refresh token |
| `GET` | `/auth/me` | Bearer | Implementado |
| `POST` | `/auth/logout` | Bearer | Implementado; revoga a sessão no banco |
| `POST` | `/auth/password/reset-request` | Não | Implementado |
| `POST` | `/auth/password/reset-confirm` | Token de uso único | Implementado |
| `PUT` | `/auth/password` | Bearer | Implementado; encerra outras sessões |
| `POST` | `/auth/email/verify-request` | Não | Implementado |
| `POST` | `/auth/email/verify` | Token de uso único | Implementado |
| `GET` | `/auth/sessions` | Bearer | Implementado; lista sessões ativas e identifica a atual |
| `DELETE` | `/auth/sessions/:id` | Bearer | Implementado; revoga somente sessão do proprietário |
| `DELETE` | `/auth/sessions` | Bearer | Implementado; revoga todas as outras sessões |

Cadastro, recuperação e confirmação dependem de SMTP configurado no backend.
Novas contas só podem entrar após confirmar o e-mail. Senhas têm de 8 a 72
bytes; nomes de usuário têm de 3 a 30 caracteres alfanuméricos ou `_`.
Access tokens JWT duram 15 minutos por padrão; refresh tokens são rotacionados,
armazenados como hash no banco e enviados em cookie `HttpOnly` com duração
padrão de 30 dias. Contas já existentes no momento da migração 0010 são
consideradas verificadas para evitar bloqueio de usuários atuais.

### Usuários e perfil

| Método | Endpoint | Autenticação | Estado |
|---|---|---|---|
| `GET` | `/users/me` | Bearer | Implementado; retorna dados privados e preferências |
| `PUT` | `/users/me` | Bearer | Implementado |
| `PUT` | `/users/me/profile` | Bearer | Implementado |
| `PUT` | `/users/me/privacy` | Bearer | Implementado |
| `GET` | `/users/me/data-export` | Bearer | Implementado; exportação JSON |
| `DELETE` | `/users/me` | Bearer | Implementado; exclusão definitiva e cascata |
| `POST` | `/users/me/consent` | Bearer | Implementado; registra versão e concessão/retirada |
| `PUT` | `/users/me/favorites` | Bearer | Implementado; até quatro livros existentes |
| `PUT` | `/users/me/annual-challenge` | Bearer | Implementado |
| `GET` | `/users/:username` | Não | Implementado; respeita privacidade e não expõe e-mail |
| `GET` | `/users/:username/lists` | Não | Implementado; somente listas públicas |
| `GET` | `/users/:username/favorites` | Não | Implementado; respeita privacidade |
| `GET` | `/users/:username/annual-challenge` | Não | Implementado; somente desafio público |
| `GET` | `/users/:username/logs` | Não | Implementado; respeita privacidade do diário |
| `GET` | `/users/:username/shelf` | Não | Implementado; respeita privacidade da estante |
| `GET` | `/users/:username/stats` | Não | Implementado; respeita privacidade das estatísticas |

O perfil aceita `displayName`, `bio`, `avatarUrl`, `bannerUrl` e `location`.
Privacidade é configurada para `profile`, `diary`, `progress`, `shelf`, `notes`,
`comments`, `lists`, `stats` e `favorites`, cada qual como `public` ou `private`.
O padrão mantém públicos os recursos já públicos e deixa notas privadas.
Recursos de notas e progresso ainda sem endpoints não são expostos por esta API.
Consentimentos aceitos: `terms_of_service`, `privacy_policy` e `marketing`;
cada chamada registra uma nova evidência com versão, decisão e horário.
O desafio anual aceita anos de 2000 até o próximo ano e metas de 1 a 1000 livros.
`PUT /users/me/profile` recebe `displayName`, `bio`, `avatarUrl`, `bannerUrl` e
`location` (limites, respectivamente: 80, 500, 2048, 2048 e 120 caracteres).
`PUT /users/me/privacy` substitui o objeto completo de preferências; todos os
nove campos são obrigatórios. `PUT /users/me/favorites` recebe
`{"bookIds":["<uuid>"]}`. O desafio recebe `{"year":2026,"targetBooks":24,"isPublic":true}`.
`POST /users/me/consent` recebe `{"consentType":"marketing","version":"1.0","granted":false}`.
`GET /users/me/data-export` retorna um envelope JSON com `exportedAt` e os
dados da conta, sem hash de senha nem tokens de sessão. `DELETE /users/me`
exclui imediatamente a conta e os registros associados em cascata.

### Redes sociais

| Método | Endpoint | Autenticação | Estado |
|---|---|---|---|
| `POST` | `/users/:id/follow` | Bearer | Implementado; repetição idempotente |
| `DELETE` | `/users/:id/follow` | Bearer | Implementado; remoção idempotente |
| `GET` | `/users/:id/followers` | Não | Implementado; exige perfil público |
| `GET` | `/users/:id/following` | Não | Implementado; exige perfil público |

As listas retornam `id`, `username` e `followedAt`, ordenadas da relação mais
recente para a mais antiga. Seguir a própria conta retorna `400`; usuário
inexistente ou com perfil privado retorna `404`. Seguir pela segunda vez
continua em `200`, e deixar de seguir uma conta não seguida continua em `200`.

### Livros

| Método | Endpoint | Autenticação | Estado |
|---|---|---|---|
| `GET` | `/books/search?q=:q` | Não | Implementado |
| `GET` | `/books/:id` | Não | Implementado |
| `GET` | `/books/popular` | Não | Implementado |

A busca já consulta o catálogo local e integra Open Library/Google Books.
Também há ISBN, gêneros e página no modelo do livro, mas não há um conjunto
completo de filtros ou endpoints de descoberta.

O detalhe e o ranking de livros populares calculam média, contagem e histograma
apenas com registros de usuários que deixaram o diário público. O histograma
mantém cinco faixas; notas de meia estrela são agrupadas na faixa inteira
superior (por exemplo, 4,5 na faixa 5). A lista de registros recentes segue a
mesma regra de visibilidade.

### Leituras e diário

| Método | Endpoint | Autenticação | Estado |
|---|---|---|---|
| `POST` | `/logs` | Bearer | Implementado |
| `GET` | `/logs/:id` | Bearer | Implementado; somente proprietário |
| `PUT` | `/logs/:id` | Bearer | Implementado; somente proprietário |
| `DELETE` | `/logs/:id` | Bearer | Implementado; somente proprietário |

Essas rotas permitem criar, consultar, editar e excluir registros com livro,
avaliação, resenha, data, releitura e marcação de spoilers. Elas ainda não
aceitam leitura sem avaliação; para isso, use as sessões de leitura. Rascunhos
e resenhas publicadas são recursos separados e não são publicados
automaticamente ao registrar um log.

#### Sessões de leitura e progresso

| Método | Endpoint | Autenticação | Estado |
|---|---|---|---|
| `POST` | `/reading-sessions` | Bearer | Implementado |
| `GET` | `/reading-sessions/:id` | Bearer | Implementado; somente proprietário |
| `PUT` | `/reading-sessions/:id` | Bearer | Implementado; atualiza a sessão completa |
| `DELETE` | `/reading-sessions/:id` | Bearer | Implementado; remove progresso e histórico associados |
| `GET` | `/users/me/reading-sessions` | Bearer | Implementado |
| `GET` | `/reading-sessions/:id/history` | Bearer | Implementado |
| `POST` | `/reading-sessions/:id/progress` | Bearer | Implementado |
| `GET` | `/reading-sessions/:id/progress` | Bearer | Implementado |
| `PUT` | `/reading-sessions/:id/status` | Bearer | Implementado |

Uma sessão é privada e aceita os status `want_to_read`, `reading`, `read` e
`abandoned`. A criação e as alterações de status sincronizam a estante (`tbr`,
`reading`, `read`, `abandoned`) e registram histórico. O progresso aceita
`date` (`YYYY-MM-DD`), `currentPage`, `pagesRead` e `note`; há no máximo um
registro por sessão e dia, e uma nova chamada para o mesmo dia o atualiza.
`PUT /reading-sessions/:id` recebe a sessão completa (`bookId`, `status`,
`currentPage`, `totalPages`, `startedAt`, `finishedAt`, `isReread`). Quando
`totalPages` é conhecido, a resposta inclui `progressPercent` calculado.

#### Anotações privadas e citações

| Método | Endpoint | Autenticação | Estado |
|---|---|---|---|
| `POST` | `/notes` | Bearer | Implementado |
| `GET` | `/notes/:id` | Bearer opcional | Implementado; proprietário ou nota publicada |
| `PUT` | `/notes/:id` | Bearer | Implementado; somente proprietário |
| `DELETE` | `/notes/:id` | Bearer | Implementado; exclusão permanente |
| `GET` | `/users/me/notes` | Bearer | Implementado; aceita `q` e `limit` |
| `POST` | `/notes/:id/tags` | Bearer | Implementado |
| `DELETE` | `/notes/:id/tags/:tagId` | Bearer | Implementado |
| `POST` | `/notes/:id/quotes` | Bearer | Implementado |
| `PUT` | `/notes/:id/quotes/:quoteId` | Bearer | Implementado |
| `DELETE` | `/notes/:id/quotes/:quoteId` | Bearer | Implementado |
| `POST` | `/notes/:id/publish` | Bearer | Implementado; publicação explícita |
| `POST` | `/notes/:id/to-review` | Bearer | Implementado; cria rascunho de resenha |

Cada nota deve apontar para exatamente um `bookId` ou `readingSessionId`.
Notas são privadas por padrão; tags têm até 50 caracteres e citações aceitam
texto, página e capítulo.

#### Rascunhos e resenhas

| Método | Endpoint | Autenticação | Estado |
|---|---|---|---|
| `POST` | `/reviews/drafts` | Bearer | Implementado |
| `GET` | `/reviews/drafts/:id` | Bearer | Implementado; somente proprietário |
| `PUT` | `/reviews/drafts/:id` | Bearer | Implementado; somente proprietário |
| `DELETE` | `/reviews/drafts/:id` | Bearer | Implementado |
| `GET` | `/reviews/:id` | Não | Implementado; somente resenhas publicadas |
| `POST` | `/reviews/:id/publish` | Bearer | Implementado |
| `POST` | `/reviews/:id/unpublish` | Bearer | Implementado |
| `PUT` | `/reviews/:id` | Bearer | Implementado; somente proprietário |
| `DELETE` | `/reviews/:id` | Bearer | Implementado; exclusão permanente |
| `PUT` | `/reviews/:id/privacy` | Bearer | Implementado |
| `GET` | `/reviews/:id/history` | Bearer | Implementado; somente proprietário |
| `POST` | `/logs/:id/to-review` | Bearer | Implementado; cria rascunho |

Resenhas aceitam avaliação opcional de 0,5 a 5,0 em incrementos de 0,5.
Publicar exige texto não vazio. `PUT /reviews/:id/privacy` recebe
`allowComments` e `allowLikes`. Publicação, retirada e edições são registradas
no histórico. As preferências ficam salvas; as rotas de curtida/comentário
existentes continuam vinculadas a `/logs/:id` e ainda não aceitam IDs de
resenhas dessa nova entidade.

#### Metas e importação do histórico

| Método | Endpoint | Autenticação | Estado |
|---|---|---|---|
| `GET` | `/users/me/goals` | Bearer | Implementado; aceita `year` |
| `PUT` | `/users/me/goals` | Bearer | Implementado; cria ou atualiza uma meta mensal/anual |
| `POST` | `/users/me/reading-history/import` | Bearer | Implementado; importação atômica |
| `GET` | `/users/me/reading-history/import/:jobId` | Bearer | Implementado; somente proprietário |

Metas aceitam `year`, `type` (`monthly` ou `annual`), `month` para metas
mensais e um ou ambos `targetBooks`/`targetPages`. A importação recebe
`{"entries":[{"bookId":"<uuid>","status":"read","currentPage":250,"totalPages":250,"startedAt":"2026-01-01","finishedAt":"2026-01-12","isReread":false}]}`.
Ela aceita de 1 a 500 sessões vinculadas a livros já cadastrados, valida e
grava todo o lote em uma transação e retorna um `jobId` concluído; o endpoint
de consulta informa os IDs das sessões criadas.

### Curtidas e comentários

| Método | Endpoint | Autenticação | Estado |
|---|---|---|---|
| `POST` | `/logs/:id/like` | Bearer | Implementado |
| `DELETE` | `/logs/:id/like` | Bearer | Implementado |
| `POST` | `/logs/:id/comments` | Bearer | Parcial |
| `GET` | `/logs/:id/comments` | Não | Parcial |

O comentário possui texto e data, mas não há edição, exclusão, respostas,
moderação ou controle de comentários por publicação.

### Estantes

| Método | Endpoint | Autenticação | Estado |
|---|---|---|---|
| `PUT` | `/shelf` | Bearer | Parcial |
| `DELETE` | `/shelf/:bookId` | Bearer | Implementado; remove o livro da estante própria |
| `GET` | `/users/:username/shelf` | Não | Parcial |

Os status existentes são `tbr`, `reading`, `read` e `abandoned`. A API não
possui endpoint para listar contagens nem visibilidade por estante. Mover um
livro entre estantes é feito atualizando seu status com `PUT /shelf`; não é
necessário um endpoint `POST /shelf/:bookId/move`.

### Listas

| Método | Endpoint | Autenticação | Estado |
|---|---|---|---|
| `GET` | `/lists/trending` | Não | Parcial |
| `POST` | `/lists` | Bearer | Parcial |
| `GET` | `/lists` | Bearer | Parcial |
| `GET` | `/lists/:id` | Bearer opcional | Implementado; respeita privacidade da lista e do usuário |
| `PUT` | `/lists/:id` | Bearer | Parcial |
| `DELETE` | `/lists/:id` | Bearer | Parcial |
| `POST` | `/lists/:id/items` | Bearer | Parcial |
| `DELETE` | `/lists/:id/items` | Bearer | Parcial |

Há título, descrição, privacidade, itens, posição e notas por item. Faltam
arquivamento, data de publicação, compartilhamento, seguidores, curadoria,
duplicação, operação segura de reordenação e confirmação removível no cliente.

### Feed

| Método | Endpoint | Autenticação | Estado |
|---|---|---|---|
| `GET` | `/feed` | Bearer | Parcial |
| `GET` | `/feed/global` | Não | Parcial |
| `GET` | `/feed/reviews` | Não | Parcial |

O feed usa cursor em parte das consultas e possui variantes seguido/global/
resenhas. Os feeds públicos ignoram registros de diários privados; o feed de
resenhas inclui resenhas publicadas e resenhas antigas de diários públicos.
Ainda faltam um feed pessoal separado e filtros por gênero, autor e período.

### Estatísticas

| Método | Endpoint | Autenticação | Estado |
|---|---|---|---|
| `GET` | `/users/:username/stats` | Não | Parcial |

Já existe cache anual calculado pelo worker assíncrono e endpoints para
configurar metas mensais/anuais. Faltam períodos mensais/semanal, cálculo de
progresso das metas, progresso da comunidade e mais granularidade de visibilidade.

## Endpoints que precisam ser criados

As rotas abaixo são propostas ou implementadas. Recursos concluídos estão
marcados como `Implementado`; os demais continuam no roadmap.

### Sessões

| Prioridade | Estado | Método | Endpoint proposto | Objetivo |
|---|---|---|---|---|
| P1 | Implementado | `GET` | `/auth/sessions` | Listar sessões ativas |
| P1 | Implementado | `DELETE` | `/auth/sessions/:id` | Revogar uma sessão específica |
| P1 | Implementado | `DELETE` | `/auth/sessions` | Revogar todas as outras sessões |

### Perfil completo

| Estado | Método | Endpoint | Objetivo |
|---|---|---|---|
| Implementado | `PUT` | `/users/me/profile` | Atualizar banner, descrição, cidade/região e dados públicos opcionais |
| Implementado | `GET` | `/users/:username/lists` | Listar as listas públicas do usuário |
| Implementado | `GET` | `/users/:username/favorites` | Consultar livros favoritos |
| Implementado | `PUT` | `/users/me/favorites` | Atualizar livros favoritos |
| Implementado | `PUT` | `/users/me/annual-challenge` | Definir meta anual de leitura |
| Implementado | `GET` | `/users/:username/annual-challenge` | Consultar desafio público |

### Leituras e progresso

| Estado | Método | Endpoint | Objetivo |
|---|---|---|---|
| Implementado | `POST` | `/reading-sessions` | Criar uma sessão privada de leitura |
| Implementado | `GET` | `/reading-sessions/:id` | Consultar uma sessão |
| Implementado | `PUT` | `/reading-sessions/:id` | Atualizar páginas, datas e releitura |
| Implementado | `DELETE` | `/reading-sessions/:id` | Excluir uma sessão |
| Implementado | `GET` | `/users/me/reading-sessions` | Listar leituras privadas do usuário |
| Implementado | `GET` | `/reading-sessions/:id/history` | Consultar histórico de status e alterações |
| Implementado | `POST` | `/reading-sessions/:id/progress` | Registrar progresso diário |
| Implementado | `GET` | `/reading-sessions/:id/progress` | Listar progresso registrado |
| Implementado | `PUT` | `/reading-sessions/:id/status` | Alterar status da leitura |
| Implementado | `GET` | `/users/me/goals` | Consultar metas mensais e anuais |
| Implementado | `PUT` | `/users/me/goals` | Criar ou atualizar metas mensais/anuais |
| Implementado | `POST` | `/users/me/reading-history/import` | Importar histórico de leitura |
| Implementado | `GET` | `/users/me/reading-history/import/:jobId` | Consultar resultado da importação |

### Anotações privadas

| Estado | Método | Endpoint | Objetivo |
|---|---|---|---|
| Implementado | `POST` | `/notes` | Criar nota privada vinculada ao livro ou sessão |
| Implementado | `GET` | `/notes/:id` | Consultar nota própria |
| Implementado | `PUT` | `/notes/:id` | Editar uma nota |
| Implementado | `DELETE` | `/notes/:id` | Excluir permanentemente uma nota |
| Implementado | `GET` | `/users/me/notes` | Listar e pesquisar notas próprias |
| Implementado | `POST` | `/notes/:id/tags` | Adicionar tags pessoais |
| Implementado | `DELETE` | `/notes/:id/tags/:tagId` | Remover uma tag pessoal |
| Implementado | `POST` | `/notes/:id/quotes` | Criar uma citação com página ou capítulo |
| Implementado | `PUT` | `/notes/:id/quotes/:quoteId` | Editar uma citação |
| Implementado | `DELETE` | `/notes/:id/quotes/:quoteId` | Excluir uma citação |
| Implementado | `POST` | `/notes/:id/publish` | Publicar uma anotação explicitamente |
| Implementado | `POST` | `/notes/:id/to-review` | Criar rascunho a partir da anotação |

### Rascunhos e resenhas

| Estado | Método | Endpoint | Objetivo |
|---|---|---|---|
| Implementado | `POST` | `/reviews/drafts` | Criar rascunho separado da leitura |
| Implementado | `GET` | `/reviews/drafts/:id` | Consultar rascunho próprio |
| Implementado | `PUT` | `/reviews/drafts/:id` | Editar rascunho |
| Implementado | `DELETE` | `/reviews/drafts/:id` | Excluir rascunho |
| Implementado | `POST` | `/reviews/:id/publish` | Publicar ou reapublicar uma resenha |
| Implementado | `POST` | `/reviews/:id/unpublish` | Retirar publicação |
| Implementado | `PUT` | `/reviews/:id` | Editar resenha |
| Implementado | `DELETE` | `/reviews/:id` | Excluir permanentemente uma resenha |
| Implementado | `PUT` | `/reviews/:id/privacy` | Configurar comentários e curtidas |
| Implementado | `GET` | `/reviews/:id/history` | Consultar datas de publicação e edição |
| Implementado | `POST` | `/logs/:id/to-review` | Criar rascunho a partir de uma leitura |

### Estantes

| Prioridade | Estado | Método | Endpoint proposto | Objetivo |
|---|---|---|---|---|
| P0 | Implementado | `DELETE` | `/shelf/:bookId` | Remover livro da estante |
| — | Não necessário | — | `/shelf/:bookId/move` | `PUT /shelf` já atualiza o status e move o livro |
| P1 | Pendente | `GET` | `/users/me/shelf/counts` | Obter contagens reais por estante |
| P1 | Pendente | `PUT` | `/users/me/shelf/privacy` | Definir visibilidade por estante |

### Listas

| Prioridade | Estado | Método | Endpoint proposto | Objetivo |
|---|---|---|---|---|
| P0 | Pendente | `POST` | `/lists/:id/archive` | Arquivar uma lista |
| P0 | Pendente | `POST` | `/lists/:id/unarchive` | Restaurar uma lista arquivada |
| P0 | Pendente | `PUT` | `/lists/:id/items/reorder` | Reordenar livros atomicamente |
| P0 | Já coberto | `DELETE` | `/lists/:id/items` | Remove um livro específico via `bookId` no corpo; valida propriedade |
| P0 | Já coberto | `POST` | `/lists/:id/items` | Insere/atualiza item e nota pessoal; valida propriedade |
| P1 | Pendente | `POST` | `/lists/:id/follow` | Seguir uma lista pública |
| P1 | Pendente | `DELETE` | `/lists/:id/follow` | Deixar de seguir uma lista |
| P1 | Pendente | `GET` | `/lists/:id/followers` | Consultar seguidores da lista |
| P1 | Pendente | `POST` | `/lists/:id/duplicate` | Criar cópia própria da lista |
| P1 | Pendente | `PUT` | `/lists/:id/privacy` | Alterar privacidade da lista |
| P1 | Pendente | `POST` | `/lists/:id/share` | Criar vínculo ou convite de compartilhamento |
| P1 | Pendente | `PUT` | `/lists/:id/curators/:userId` | Adicionar ou remover curadoria |

### Feed e recomendações

| Prioridade | Método | Endpoint proposto | Objetivo |
|---|---|---|---|
| P0 | `GET` | `/feed/personal` | Consultar diário privado e progresso próprio |
| P0 | `GET` | `/feed/publications` | Consultar apenas resenhas/opiniões publicadas |
| P0 | `PUT` | `/users/me/feed/preferences` | Configurar o que será compartilhado no feed |
| P0 | `GET` | `/feed/reviews/search` | Filtrar feed por gênero, autor e período |
| P0 | `POST` | `/users/:id/mute` | Silenciar um usuário no feed |
| P0 | `DELETE` | `/users/:id/mute` | Reativar conteúdo de um usuário |
| P0 | `POST` | `/reports` | Denunciar abuso ou conteúdo inadequado |
| P1 | `GET` | `/recommendations/books` | Recomendar próximo livro |
| P1 | `GET` | `/recommendations/similar/:bookId` | Buscar livros semelhantes |
| P1 | `GET` | `/recommendations/users` | Encontrar usuários com gosto literário parecido |
| P1 | `GET` | `/rankings/books` | Ranking real baseado em leituras da comunidade |
| P1 | `PUT` | `/users/me/genres` | Atualizar gêneros preferidos |
| P1 | `GET` | `/users/:username/shared-current-status` | Mostrar quero ler, lendo ou concluído, respeitando privacidade |

### Comentários e interações

| Prioridade | Método | Endpoint proposto | Objetivo |
|---|---|---|---|
| P1 | `PUT` | `/logs/:id/comments/:commentId` | Editar comentário próprio |
| P1 | `DELETE` | `/logs/:id/comments/:commentId` | Excluir comentário próprio |
| P1 | `POST` | `/comments/:commentId/replies` | Criar resposta a comentário |
| P1 | `GET` | `/comments/:commentId/replies` | Listar respostas |
| P1 | `POST` | `/comments/:commentId/like` | Curtir comentário |
| P1 | `DELETE` | `/comments/:commentId/like` | Remover curtida de comentário |

### Descoberta e comunidade

| Prioridade | Método | Endpoint proposto | Objetivo |
|---|---|---|---|
| P1 | `GET` | `/books/search/advanced` | Buscar por título, autor, ISBN e gênero |
| P1 | `GET` | `/communities` | Listar comunidades ou clubes |
| P1 | `POST` | `/communities` | Criar comunidade |
| P1 | `GET` | `/communities/:id` | Consultar comunidade |
| P1 | `POST` | `/communities/:id/join` | Entrar em comunidade |
| P1 | `DELETE` | `/communities/:id/join` | Sair de comunidade |
| P1 | `GET` | `/communities/:id/lists` | Consultar listas coletivas da comunidade |

### Notificações e auditoria

| Prioridade | Método | Endpoint proposto | Objetivo |
|---|---|---|---|
| P1 | `GET` | `/notifications` | Listar notificações |
| P1 | `PUT` | `/notifications/:id/read` | Marcar uma notificação como lida |
| P1 | `POST` | `/notifications/read-all` | Marcar todas como lidas |
| P1 | `PUT` | `/users/me/notification-preferences` | Configurar notificações |
| P1 | `GET` | `/privacy/audit` | Consultar histórico de alterações de privacidade |
| P1 | `GET` | `/content/:id/visibility` | Consultar quem pode visualizar um conteúdo |
| P1 | `PUT` | `/content/:id/visibility` | Registrar ou alterar a visibilidade de um conteúdo |

## Endpoints existentes que precisam ser corrigidos ou ampliados

Não são apenas criação de novas rotas. Os contratos existentes precisam de
ajustes antes de suportar a ideologia:

| Endpoint atual | Lacuna necessária |
|---|---|
| `POST /logs` | Permitir leitura sem avaliação, separar rascunho, nota, avaliação e publicação, e evitar duplicações |
| `PUT /logs/:id` | Revisar se a edição deve separar nota/avaliação de publicação; edição já é restrita ao proprietário |
| `DELETE /logs/:id` | Definir se a exclusão do log também deve remover ou preservar conteúdo relacionado |
| `PUT /shelf` | Validar se o frontend usa os mesmos status padronizados da API |
| `GET /users/:username/logs` | Decidir se o endpoint deve retornar progresso além dos logs; logs privados já são filtrados |
| `GET /users/:username/shelf` | Acrescentar visibilidade por estante; privacidade geral da estante já é respeitada |
| `GET /users/:username/stats` | Acrescentar granularidade de visibilidade; privacidade geral das estatísticas já é respeitada |
| `POST /lists` | Publicação explícita requer decisão de produto; listas vazias já podem ser criadas |
| `PUT /lists/:id/items` | Reordenação atômica ainda não existe; a API atual usa `POST /lists/:id/items` |
| `GET /lists/:id` | Corrigido: listas privadas retornam não encontrado para terceiros; proprietário autenticado pode consultá-las |
| `POST /lists/:id/items`, `DELETE /lists/:id/items` | Corrigido: operações de itens agora verificam a propriedade da lista |
| `GET /feed` | Separar publicações explícitas do feed de leituras continua pendente |
| `GET /feed/global` | Cursor e filtros continuam pendentes; conteúdos de diários privados são excluídos |
| `GET /feed/reviews` | Filtrar spoilers e indicar conteúdo publicado/editado |
| `POST /logs/:id/comments` | Verificar controle de comentários da publicação |
| `GET /logs/:id/comments` | Respeitar privacidade da publicação e paginação |
| `GET /books/search` | Buscar por ISBN, gênero e autor além de título/texto livre |
| `GET /users/:username` | Perfil público já respeita privacidade; revisar campos adicionais ao ampliar o contrato |

## Prioridades sugeridas

### Fase 1 - privacidade e integridade

- Criar entidades reais de sessões de leitura e notas.
- Separar resenhas publicadas de registros privados.
- Implementar visibilidade de diário, estante, perfil, listas e estatísticas.
- Criar exclusão permanente de notas e dados.
- Impedir que o feed leia ou exponha conteúdo privado.

### Fase 2 - experiência completa de leitura

- Criar endpoints de progresso, datas, releitura e histórico.
- Implementar rascunhos e publicação explícita de resenhas.
- Criar remoção específica da estante e reordenação de listas.
- Permitir edição/exclusão de comentários e respostas.

### Fase 3 - descoberta e comunidade

- Implementar filtros do feed e recomendações.
- Criar ranking real e usuários com afinidade.
- Criar comunidades, clubes e listas coletivas.
- Implementar notificações e silenciamento.

### Fase 4 - auditoria e conformidade

- Implementar auditoria de mudanças de privacidade.
- Registrar consentimentos e controlar exclusão/exportação.
- Aplicar limites de tamanho para notas e resenhas.
- Validar propriedade de recursos e cobrir vazamentos com testes automatizados.
