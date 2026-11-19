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
| `POST` | `/auth/register` | Não | Implementado |
| `POST` | `/auth/login` | Não | Implementado |
| `GET` | `/auth/me` | Bearer | Implementado |
| `POST` | `/auth/logout` | Bearer | Implementado |

Ainda não existem recuperação de senha, confirmação de email, troca de senha,
revogação de todas as sessões ou gerenciamento de dispositivos.

### Usuários e perfil

| Método | Endpoint | Autenticação | Estado |
|---|---|---|---|
| `GET` | `/users/:username` | Não | Implementado |
| `PUT` | `/users/me` | Bearer | Parcial |
| `GET` | `/users/:username/logs` | Não | Parcial |
| `GET` | `/users/:username/shelf` | Não | Parcial |
| `GET` | `/users/:username/stats` | Não | Parcial |

`PUT /users/me` altera biografia, avatar e quatro livros favoritos. Não há,
por enquanto, campos para banner, cidade/região, desafio anual, privacidade,
conteúdos públicos ou exportação de dados.

### Redes sociais

| Método | Endpoint | Autenticação | Estado |
|---|---|---|---|
| `POST` | `/users/:id/follow` | Bearer | Implementado |
| `DELETE` | `/users/:id/follow` | Bearer | Implementado |
| `GET` | `/users/:id/followers` | Não | Implementado |
| `GET` | `/users/:id/following` | Não | Implementado |

### Livros

| Método | Endpoint | Autenticação | Estado |
|---|---|---|---|
| `GET` | `/books/search?q=:q` | Não | Implementado |
| `GET` | `/books/:id` | Não | Implementado |
| `GET` | `/books/popular` | Não | Implementado |

A busca já consulta o catálogo local e integra Open Library/Google Books.
Também há ISBN, gêneros e página no modelo do livro, mas não há um conjunto
completo de filtros ou endpoints de descoberta.

### Leituras e diário

| Método | Endpoint | Autenticação | Estado |
|---|---|---|---|
| `POST` | `/logs` | Bearer | Parcial |
| `GET` | `/logs/:id` | Bearer | Parcial |
| `PUT` | `/logs/:id` | Bearer | Parcial |
| `DELETE` | `/logs/:id` | Bearer | Parcial |

Essas rotas permitem criar, consultar, editar e excluir registros com livro,
avaliação, resenha, data, releitura e marcação de spoilers. Elas ainda não
separam leitura privada, nota privada, rascunho e resenha publicada.

Campos atuais insuficientes para a experiência ideal:

- página atual e percentual concluído;
- número de páginas lidas;
- datas de início e término;
- status separado da estante;
- histórico de status;
- progresso diário ou semanal;
- tempo estimado;
- objetivos mensais e anuais;
- visibilidade individual da leitura;
- registro sem avaliação.

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
| `GET` | `/users/:username/shelf` | Não | Parcial |

Os status existentes são `tbr`, `reading`, `read` e `abandoned`. A API não
possui endpoint específico para remover um livro da estante, listar contagens
ou mover o livro entre estantes sem uma operação de atualização de status.
Também não há visibilidade por estante.

### Listas

| Método | Endpoint | Autenticação | Estado |
|---|---|---|---|
| `GET` | `/lists/trending` | Não | Parcial |
| `POST` | `/lists` | Bearer | Parcial |
| `GET` | `/lists` | Bearer | Parcial |
| `GET` | `/lists/:id` | Não | Parcial |
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
resenhas. O objetivo exige separar claramente feed social, diário pessoal e
resenhas, além de privacidade por conteúdo e filtros por gênero, autor e período.

### Estatísticas

| Método | Endpoint | Autenticação | Estado |
|---|---|---|---|
| `GET` | `/users/:username/stats` | Não | Parcial |

Já existe cache anual calculado pelo worker assíncrono. Faltam períodos
mensais/semanal, objetivos, progresso da comunidade e visibilidade pública.

## Endpoints que precisam ser criados

As rotas abaixo são propostas. A convenção pode mudar durante a implementação,
mas os recursos precisam existir em algum endpoint REST explícito.

### Privacidade e conta

| Prioridade | Método | Endpoint proposto | Objetivo |
|---|---|---|---|
| P0 | `PUT` | `/users/me/privacy` | Configurar privacidade de perfil, diário, progresso, estante, notas, comentários, listas e estatísticas |
| P0 | `GET` | `/users/me/data-export` | Iniciar ou obter exportação dos dados do usuário |
| P0 | `DELETE` | `/users/me` | Solicitar exclusão da conta e dos dados pessoais |
| P0 | `POST` | `/users/me/consent` | Registrar consentimento para processamento ou recursos opcionais |
| P1 | `POST` | `/auth/password/reset-request` | Solicitar recuperação de senha |
| P1 | `POST` | `/auth/password/reset-confirm` | Confirmar recuperação de senha |
| P1 | `PUT` | `/auth/password` | Trocar senha autenticada |
| P1 | `POST` | `/auth/email/verify-request` | Solicitar confirmação de email |
| P1 | `POST` | `/auth/email/verify` | Confirmar email |
| P1 | `GET` | `/auth/sessions` | Listar sessões ativas |
| P1 | `DELETE` | `/auth/sessions/:id` | Revogar uma sessão específica |
| P1 | `DELETE` | `/auth/sessions` | Revogar todas as outras sessões |

### Perfil completo

| Prioridade | Método | Endpoint proposto | Objetivo |
|---|---|---|---|
| P1 | `PUT` | `/users/me/profile` | Atualizar banner, descrição, cidade/região e dados públicos opcionais |
| P1 | `GET` | `/users/:username/lists` | Listar as listas públicas do usuário |
| P1 | `GET` | `/users/:username/favorites` | Consultar livros favoritos |
| P1 | `PUT` | `/users/me/favorites` | Atualizar livros favoritos |
| P1 | `PUT` | `/users/me/annual-challenge` | Definir meta anual de leitura |
| P1 | `GET` | `/users/:username/annual-challenge` | Consultar desafio público |

### Leituras e progresso

| Prioridade | Método | Endpoint proposto | Objetivo |
|---|---|---|---|
| P0 | `POST` | `/reading-sessions` | Criar uma sessão privada de leitura |
| P0 | `GET` | `/reading-sessions/:id` | Consultar uma sessão |
| P0 | `PUT` | `/reading-sessions/:id` | Atualizar páginas, percentual, datas e releitura |
| P0 | `DELETE` | `/reading-sessions/:id` | Excluir uma sessão |
| P0 | `GET` | `/users/me/reading-sessions` | Listar leituras privadas do usuário |
| P0 | `GET` | `/reading-sessions/:id/history` | Consultar histórico de status e alterações |
| P0 | `POST` | `/reading-sessions/:id/progress` | Registrar progresso diário ou semanal |
| P0 | `GET` | `/reading-sessions/:id/progress` | Listar progresso registrado |
| P1 | `PUT` | `/reading-sessions/:id/status` | Alterar quero ler, lendo, lido ou abandonado |
| P1 | `GET` | `/users/me/goals` | Consultar metas mensais e anuais |
| P1 | `PUT` | `/users/me/goals` | Criar ou atualizar metas mensais/anuais |
| P1 | `POST` | `/users/me/reading-history/import` | Importar histórico de leitura |
| P1 | `GET` | `/users/me/reading-history/import/:jobId` | Consultar resultado ou progresso da importação |

### Anotações privadas

| Prioridade | Método | Endpoint proposto | Objetivo |
|---|---|---|---|
| P0 | `POST` | `/notes` | Criar nota privada vinculada ao livro ou sessão |
| P0 | `GET` | `/notes/:id` | Consultar uma nota própria |
| P0 | `PUT` | `/notes/:id` | Editar uma nota |
| P0 | `DELETE` | `/notes/:id` | Excluir permanentemente uma nota |
| P0 | `GET` | `/users/me/notes` | Listar e pesquisar notas próprias |
| P0 | `POST` | `/notes/:id/tags` | Adicionar tags ou categorias pessoais |
| P0 | `DELETE` | `/notes/:id/tags/:tagId` | Remover uma tag pessoal |
| P1 | `POST` | `/notes/:id/quotes` | Criar uma citação com página ou capítulo |
| P1 | `PUT` | `/notes/:id/quotes/:quoteId` | Editar uma citação |
| P1 | `DELETE` | `/notes/:id/quotes/:quoteId` | Excluir uma citação |
| P1 | `POST` | `/notes/:id/publish` | Publicar uma anotação isoladamente |
| P1 | `POST` | `/notes/:id/to-review` | Transformar uma anotação privada em resenha |

### Rascunhos e resenhas

| Prioridade | Método | Endpoint proposto | Objetivo |
|---|---|---|---|
| P0 | `POST` | `/reviews/drafts` | Criar rascunho de resenha separado da leitura |
| P0 | `GET` | `/reviews/drafts/:id` | Consultar rascunho próprio |
| P0 | `PUT` | `/reviews/drafts/:id` | Editar rascunho |
| P0 | `DELETE` | `/reviews/drafts/:id` | Excluir rascunho |
| P0 | `POST` | `/reviews/:id/publish` | Publicar ou reapublicar uma resenha |
| P0 | `POST` | `/reviews/:id/unpublish` | Retirar publicação sem apagar a leitura |
| P0 | `PUT` | `/reviews/:id` | Editar resenha já publicada |
| P0 | `DELETE` | `/reviews/:id` | Excluir permanentemente uma resenha |
| P0 | `PUT` | `/reviews/:id/privacy` | Configurar comentários e curtidas |
| P1 | `GET` | `/reviews/:id/history` | Consultar datas de publicação e edição |
| P1 | `POST` | `/logs/:id/to-review` | Criar revisão a partir de uma leitura |

### Estantes

| Prioridade | Método | Endpoint proposto | Objetivo |
|---|---|---|---|
| P0 | `DELETE` | `/shelf/:bookId` | Remover livro da estante |
| P0 | `POST` | `/shelf/:bookId/move` | Mover livro entre estantes |
| P1 | `GET` | `/users/me/shelf/counts` | Obter contagens reais por estante |
| P1 | `PUT` | `/users/me/shelf/privacy` | Definir visibilidade por estante |

### Listas

| Prioridade | Método | Endpoint proposto | Objetivo |
|---|---|---|---|
| P0 | `POST` | `/lists/:id/archive` | Arquivar uma lista |
| P0 | `POST` | `/lists/:id/unarchive` | Restaurar uma lista arquivada |
| P0 | `PUT` | `/lists/:id/items/reorder` | Reordenar livros atomicamente |
| P0 | `DELETE` | `/lists/:id/items/:bookId` | Remover um livro específico da lista |
| P0 | `PUT` | `/lists/:id/items/:bookId/notes` | Atualizar nota pessoal do item |
| P1 | `POST` | `/lists/:id/follow` | Seguir uma lista pública |
| P1 | `DELETE` | `/lists/:id/follow` | Deixar de seguir uma lista |
| P1 | `GET` | `/lists/:id/followers` | Consultar seguidores da lista |
| P1 | `POST` | `/lists/:id/duplicate` | Criar cópia própria da lista |
| P1 | `PUT` | `/lists/:id/privacy` | Alterar privacidade da lista |
| P1 | `POST` | `/lists/:id/share` | Criar vínculo ou convite de compartilhamento |
| P1 | `PUT` | `/lists/:id/curators/:userId` | Adicionar ou remover curadoria |

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
| `PUT /logs/:id` | Permitir edição pelo proprietário e impedir vazamento ou publicação automática |
| `DELETE /logs/:id` | Distinguir exclusão da leitura, da nota e da publicação |
| `PUT /shelf` | Centralizar status validados e padronizados entre frontend/API |
| `GET /users/:username/logs` | Respeitar privacidade e não expor notas ou progresso privados |
| `GET /users/:username/shelf` | Respeitar visibilidade da estante |
| `GET /users/:username/stats` | Respeitar visibilidade das estatísticas |
| `PUT /users/me` | Acrescer perfil público/opcional e visibilidade sem misturar dados privados |
| `POST /lists` | Permitir criar lista vazia e publicar somente explicitamente |
| `PUT /lists/:id/items` | Adicionar reordenação estável e remover resposta ambígua em relação ao item |
| `GET /lists/:id` | Ocultar listas privadas de usuários sem acesso |
| `GET /feed` | Não incluir notas privadas e separar publicação explícita |
| `GET /feed/global` | Filtrar, paginar e remover conteúdos privados |
| `GET /feed/reviews` | Filtrar spoilers e indicar conteúdo publicado/editado |
| `POST /logs/:id/comments` | Verificar controle de comentários da publicação |
| `GET /logs/:id/comments` | Respeitar privacidade da publicação e paginação |
| `GET /books/search` | Buscar por ISBN, gênero e autor além de título/texto livre |
| `GET /users/:username` | Exibir somente perfil e estatísticas autorizados |

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
