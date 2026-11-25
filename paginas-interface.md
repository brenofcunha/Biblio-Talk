# Guia de criação da interface do Biblio-Talk

## Objetivo do documento

Este documento é o briefing para criar uma interface nova do zero. Não é uma
lista de correções da interface existente nem exige reaproveitar telas,
componentes, nomes internos ou organização visual atuais.

O produto é uma rede social para leitores com biblioteca pessoal, acompanhamento
de leitura, notas privadas, resenhas publicáveis e interação com outros
leitores. A interface deve tornar inequívoco o que é privado e o que será
compartilhado.

As integrações descritas abaixo usam como referência os contratos em
[`endpoints-api.md`](./endpoints-api.md). Quando uma necessidade visual ainda
não tem endpoint, ela está marcada como **bloqueio da API**: não desenhar uma
ação que pareça funcional até o contrato existir.

## Quantidade de páginas

Planejar **22 páginas responsivas no total**:

- **21 páginas para a primeira experiência completa**: 5 fluxos de acesso e
  recuperação e 16 páginas do produto (algumas públicas, outras autenticadas);
- **1 página futura** de importação de histórico, condicionada à definição do
  contrato de importação.

Modais, menus, abas e painéis laterais são estados/componentes dessas páginas,
não páginas adicionais. Páginas futuras sem API implementada estão listadas
separadamente ao final.

## Mapa das páginas

| Nº | Página | Acesso | Prioridade |
|---:|---|---|---|
| 1 | Entrar | Público | Obrigatória |
| 2 | Criar conta | Público | Obrigatória |
| 3 | Confirmar e-mail | Público | Obrigatória |
| 4 | Solicitar redefinição de senha | Público | Obrigatória |
| 5 | Definir nova senha | Público | Obrigatória |
| 6 | Início / minha leitura | Autenticado | Obrigatória |
| 7 | Feed da comunidade | Autenticado/público | Obrigatória |
| 8 | Busca de livros | Autenticado/público | Obrigatória |
| 9 | Detalhes do livro | Autenticado/público | Obrigatória |
| 10 | Minha biblioteca e diário | Autenticado | Obrigatória |
| 11 | Sessão de leitura | Autenticado | Obrigatória |
| 12 | Notas privadas | Autenticado | Obrigatória |
| 13 | Editor e rascunhos de resenha | Autenticado | Obrigatória |
| 14 | Metas e estatísticas | Autenticado | Obrigatória |
| 15 | Minhas listas / detalhe da lista | Autenticado/público | Obrigatória |
| 16 | Perfil público | Público | Obrigatória |
| 17 | Seguidores e seguindo | Público | Obrigatória |
| 18 | Editar meu perfil | Autenticado | Obrigatória |
| 19 | Privacidade | Autenticado | Obrigatória |
| 20 | Sessões e segurança | Autenticado | Obrigatória |
| 21 | Dados e encerramento da conta | Autenticado | Obrigatória |
| 22 | Importar histórico de leitura | Autenticado | Futura; depende de contrato de importação utilizável |

> **Atenção à página 22:** o endpoint atual aceita IDs de livros já existentes
> (`bookId`) em JSON, de 1 a 500 entradas. Não há upload/importação direta de
> arquivos de outros serviços nem correspondência automática de livros. A tela
> precisa ser desenhada como fluxo de revisão/mapeamento, ou ficar fora da
> primeira entrega até esse contrato ser definido.

## Especificação página por página

### 1. Entrar

**Objetivo:** autenticar a pessoa e levar para a área privada.

**Elementos necessários**

- Campo de e-mail e senha, alternância para exibir senha e validação inline.
- Ação principal “Entrar”, link para cadastro e link de recuperação.
- Mensagens distintas para credenciais inválidas, conta não verificada e
  indisponibilidade temporária.

**Integração**

- `POST /api/v1/auth/login` com `email` e `password`.
- Após sucesso, usar o token de acesso e o cookie de refresh HttpOnly; carregar
  `GET /api/v1/auth/me` ou `GET /api/v1/users/me`.
- Se a API responder que o e-mail não foi confirmado, oferecer ir à página 3
  e solicitar novo e-mail via `POST /api/v1/auth/email/verify-request`.

### 2. Criar conta

**Objetivo:** cadastrar uma conta e explicar que a confirmação de e-mail é
obrigatória para entrar.

**Elementos necessários**

- Nome de usuário, e-mail, senha e confirmação visual dos requisitos de senha.
- Estados de envio, sucesso, conflito de usuário/e-mail e falha no envio do
  e-mail de confirmação.
- Não exibir a conta como pronta para uso antes da confirmação.

**Integração**

- `POST /api/v1/auth/register` com `username`, `email` e `password`.
- Usar `verificationRequired` e `emailSent` para montar a confirmação e as
  alternativas de reenvio.

### 3. Confirmar e-mail

**Objetivo:** concluir ou repetir a confirmação recebida por e-mail.

**Elementos necessários**

- Estado inicial “verificando”, sucesso, token expirado/inválido e falha de
  conexão.
- Ação para reenviar confirmação solicitando o e-mail.

**Integração**

- Ler `verifyToken` da URL e enviar `POST /api/v1/auth/email/verify` com
  `token`.
- Reenvio por `POST /api/v1/auth/email/verify-request` com `email`.

### 4. Solicitar redefinição de senha

**Objetivo:** solicitar o envio de link para redefinir a senha.

**Elementos necessários**

- Campo de e-mail, confirmação neutra que não revele se a conta existe e opção
  de voltar ao login.
- Estados de envio e indisponibilidade de e-mail.

**Integração**

- `POST /api/v1/auth/password/reset-request` com `email`.

### 5. Definir nova senha

**Objetivo:** definir uma senha nova a partir do link de recuperação.

**Elementos necessários**

- Nova senha, confirmação, requisitos visíveis, sucesso e token inválido ou
  expirado.
- Após sucesso, direcionar para entrar novamente.

**Integração**

- Ler `resetToken` da URL.
- `POST /api/v1/auth/password/reset-confirm` com `token` e `password`.

### 6. Início / minha leitura

**Objetivo:** responder “o que estou lendo e qual é minha próxima ação?”.

**Elementos necessários**

- Leituras em andamento com capa, página atual, percentual quando conhecido e
  atalho para continuar.
- Próximas leituras e livros concluídos recentes.
- Resumo das metas somente quando calculado a partir de dados reais; distinguir
  metas configuradas de progresso alcançado.
- Atalhos para buscar livro, registrar leitura, notas, diário e listas.
- Estado vazio com ação para começar; nunca usar números demonstrativos como
  se fossem dados reais.

**Integração**

- `GET /api/v1/users/me`.
- `GET /api/v1/users/me/reading-sessions`.
- `GET /api/v1/users/me/goals?year={ano}`.
- `GET /api/v1/books/popular` para conteúdo de descoberta.
- Abrir a página 11 para uma leitura selecionada.

**Limite atual:** a API ainda não retorna progresso calculado das metas. Não
representar uma meta como concluída ou calcular progresso de livros sem uma
regra acordada.

### 7. Feed da comunidade

**Objetivo:** acompanhar pessoas e descobrir conteúdo publicado.

**Elementos necessários**

- Abas “Seguindo”, “Global” e “Resenhas”.
- Cards distinguindo atividade de leitura, avaliação e resenha; indicar autor,
  livro, data e spoiler.
- Curtida e comentários apenas quando permitidos; estados de carregamento,
  vazio, erro e paginação/cursor.
- Menu contextual para denunciar ou silenciar somente quando essas ações
  estiverem implementadas.

**Integração**

- `GET /api/v1/feed`, com `limit` e `cursor`.
- `GET /api/v1/feed/global?limit={n}`.
- `GET /api/v1/feed/reviews?limit={n}`.
- `POST` / `DELETE /api/v1/logs/:id/like`.
- `POST` / `GET /api/v1/logs/:id/comments`.

**Limites atuais:** feed global e feed de resenhas não oferecem filtros por
gênero, autor e período; o controle de comentários/curtidas das novas resenhas
não está conectado às rotas de interação. Cursor está disponível no feed de
quem a pessoa segue; global e resenhas aceitam `limit`, mas não oferecem
paginação por cursor. Não desenhar filtros, ações ou rolagem infinita como
funcionais onde o contrato não fornece suporte.

### 8. Busca de livros

**Objetivo:** localizar livros e abrir seus detalhes.

**Elementos necessários**

- Campo de busca com resultados em lista ou grade, capa, título, autores e
  metadados disponíveis.
- Estados de consulta, sem resultados, erro e carregamento incremental.
- Filtros avançados somente se disponíveis no endpoint.

**Integração**

- `GET /api/v1/books/search?q={termo}&limit={n}`.
- Selecionar um resultado abre a página 9.

**Limite atual:** a busca integrada aceita texto livre; filtros avançados por
ISBN, gênero, ano e autor ainda não existem.

### 9. Detalhes do livro

**Objetivo:** concentrar informações do livro, ações de leitura e conteúdo
público.

**Elementos necessários**

- Capa, título, autores, descrição e metadados retornados pela API.
- Avaliação da comunidade e distribuição quando disponíveis.
- Ações para colocar/remover da estante, iniciar/atualizar leitura e abrir
  notas próprias.
- Seções separadas para “Minha leitura”, “Minhas notas” e “Resenhas públicas”.
- Resenhas com indicador de spoiler e publicação; vazio quando não houver.
- Confirmação explícita antes de publicar conteúdo.

**Integração**

- `GET /api/v1/books/:id`.
- `PUT /api/v1/shelf` com `bookId` e status; `DELETE /api/v1/shelf/:bookId`.
- `POST /api/v1/reading-sessions` para iniciar leitura.
- `GET /api/v1/users/me/notes?q={termo}&limit={n}` para notas próprias e
  `POST /api/v1/notes` para criar.
- `POST /api/v1/reviews/drafts` para iniciar rascunho.

**Limite atual:** `GET /api/v1/books/:id` inclui registros antigos de leitura
com avaliações/textos, mas não consulta a nova entidade de resenhas publicada.
Para mostrar todas as resenhas, a API precisa incluir as duas origens ou
oferecer um endpoint de resenhas por livro. Não tratar a avaliação agregada
como se fosse uma resenha.

### 10. Minha biblioteca e diário

**Objetivo:** consultar e organizar as leituras pessoais sem expor conteúdo
privado.

**Elementos necessários**

- Abas de status: quero ler, lendo, lido e abandonado.
- Busca e filtros pelos dados disponíveis; cada item mostra capa, status,
  progresso e datas quando presentes.
- Ações de continuar, mudar status, abrir detalhes, remover da estante e
  confirmar exclusões.
- Separação visual entre sessão de leitura, registro/avaliação, nota privada e
  resenha publicada.

**Integração**

- `GET /api/v1/users/me/reading-sessions?limit={n}`.
- `PUT /api/v1/reading-sessions/:id/status`.
- `DELETE /api/v1/shelf/:bookId`.
- A listagem de sessões aceita limite, mas não busca/filtros no servidor;
  filtrar apenas os registros carregados e não prometer catálogo completo
  paginado.
- Estante pública de outro usuário não deve ser usada como substituto da
  biblioteca privada.

**Bloqueio da API:** não existe `GET /api/v1/users/me/logs`. A rota
`GET /api/v1/users/:username/logs` obedece à privacidade pública e não serve
para listar logs do próprio usuário quando o diário está privado. Para mostrar
todos os registros antigos do diário pessoal, criar endpoint autenticado
próprio ou definir a migração desses dados para sessões.

### 11. Sessão de leitura

**Objetivo:** acompanhar uma leitura individual e seu histórico.

**Elementos necessários**

- Livro, status, página atual, total de páginas, percentual, início, conclusão
  e releitura.
- Formulário rápido para registrar progresso diário e nota pessoal.
- Histórico de alterações/status ordenado por data.
- Ações para atualizar, concluir, abandonar ou excluir, com confirmação na
  exclusão.

**Integração**

- `GET /api/v1/reading-sessions/:id`.
- `PUT /api/v1/reading-sessions/:id` para a sessão completa.
- `PUT /api/v1/reading-sessions/:id/status`.
- `POST` / `GET /api/v1/reading-sessions/:id/progress`.
- `GET /api/v1/reading-sessions/:id/history`.
- `DELETE /api/v1/reading-sessions/:id`.

### 12. Notas privadas

**Objetivo:** criar, pesquisar e organizar anotações separadas das publicações.

**Elementos necessários**

- Lista pesquisável de notas com livro relacionado, trecho, tags e data de
  atualização.
- Editor de conteúdo; criar/editar/excluir com confirmação.
- Gestão de tags e citações com página/capítulo.
- Estado explícito “privada” e uma ação separada para publicar.
- Não incluir nota privada em feeds, cards de perfil ou busca pública.

**Integração**

- `GET /api/v1/users/me/notes?q={termo}&limit={n}`.
- `POST /api/v1/notes`, `GET /api/v1/notes/:id`, `PUT /api/v1/notes/:id` e
  `DELETE /api/v1/notes/:id`.
- `POST /api/v1/notes/:id/tags` e `DELETE /api/v1/notes/:id/tags/:tagId`.
- `POST /api/v1/notes/:id/quotes`, `PUT /api/v1/notes/:id/quotes/:quoteId` e
  `DELETE /api/v1/notes/:id/quotes/:quoteId`.
- `POST /api/v1/notes/:id/publish`.

**Limite atual:** não existe ação/endpoint para retirar uma nota publicada.
Projetar a publicação, mas não oferecer “despublicar” até a API suportar a
operação.

### 13. Editor e rascunhos de resenha

**Objetivo:** compor uma resenha em separado da sessão de leitura e decidir
explicitamente quando publicá-la.

**Elementos necessários**

- Editor com avaliação opcional, texto, spoiler e livro relacionado.
- Estados de rascunho, publicada e despublicada.
- Pré-visualização do conteúdo público, opção de permitir curtidas/comentários
  e histórico de edição/publicação.
- Confirmação antes de publicar ou excluir.

**Integração**

- Criar: `POST /api/v1/reviews/drafts`.
- Consultar/editar/excluir: `GET /api/v1/reviews/drafts/:id`,
  `PUT /api/v1/reviews/drafts/:id` e `DELETE /api/v1/reviews/drafts/:id`.
- Publicar/despublicar: `POST /api/v1/reviews/:id/publish` e
  `POST /api/v1/reviews/:id/unpublish`.
- Editar resenha: `PUT /api/v1/reviews/:id`.
- Preferências: `PUT /api/v1/reviews/:id/privacy`.
- Histórico: `GET /api/v1/reviews/:id/history`.

**Bloqueio da API:** falta endpoint para listar todos os rascunhos e resenhas
do próprio usuário. Sem ele, a página não consegue montar um arquivo de
rascunhos após navegação ou novo login. Criar, por exemplo,
`GET /api/v1/users/me/reviews` com filtro de status.

### 14. Metas e estatísticas

**Objetivo:** definir metas e acompanhar estatísticas reais, sem inventar
progresso.

**Elementos necessários**

- Formulário para meta anual/mensal de livros e/ou páginas.
- Seleção de ano e período; gráficos somente a partir de dados da API.
- Separação visual entre alvo configurado e resultado medido.
- Aviso/estado vazio quando não houver histórico suficiente.

**Integração**

- `GET /api/v1/users/me/goals?year={ano}` e
  `PUT /api/v1/users/me/goals`.
- `GET /api/v1/users/:username/stats?year={ano}` para estatísticas anuais
  autorizadas.
- `PUT /api/v1/users/me/annual-challenge` e
  `GET /api/v1/users/:username/annual-challenge` para desafio anual.

**Limite atual:** a API salva metas, mas ainda não calcula progresso para
livros/páginas nem oferece períodos semanal/mensal. Exibir metas como
configuração até que essas informações sejam fornecidas.

**Bloqueio da API:** `GET /api/v1/users/:username/stats` é uma rota pública
protegida pela configuração de privacidade; ela não permite que o próprio
usuário consulte estatísticas privadas. Criar `GET /api/v1/users/me/stats`
antes de exigir estatísticas privadas no painel pessoal.

### 15. Minhas listas / detalhe da lista

**Objetivo:** criar coleções de livros, controlar privacidade e organizar itens.

**Elementos necessários**

- Visão das listas próprias com título, descrição, privacidade e quantidade de
  livros.
- Modelo de detalhe da lista com itens, capa, posição e nota por item.
- Criar lista vazia, editar dados, adicionar/atualizar/remover item e excluir
  lista.
- Confirmar exclusão e mostrar estados para lista vazia/privada/não encontrada.
- Para listas de outras pessoas, respeitar a visibilidade e ocultar itens
  privados.

**Integração**

- `GET /api/v1/lists`, `POST /api/v1/lists`, `GET /api/v1/lists/:id`,
  `PUT /api/v1/lists/:id` e `DELETE /api/v1/lists/:id`.
- `POST /api/v1/lists/:id/items` com `bookId`, `position` e `notes`.
- `DELETE /api/v1/lists/:id/items` com `bookId` no corpo.
- `GET /api/v1/users/:username/lists` para listas públicas do perfil.
- `GET /api/v1/lists/trending` para descoberta.

**Limites atuais:** arquivamento e reordenação atômica ainda não existem.
Seguir, compartilhar, duplicar e curadoria também dependem de endpoints novos.

### 16. Perfil público

**Objetivo:** apresentar a identidade e atividade pública de um leitor,
respeitando suas preferências de privacidade.

**Elementos necessários**

- Avatar, banner, nome, bio, localização e botão seguir/deixar de seguir.
- Abas para diário, resenhas, estante, listas, favoritos, desafio e estatísticas;
  esconder ou substituir por estado privado quando a API não autorizar.
- Nunca exibir e-mail ou dados privados.
- Estado de perfil não encontrado/privado e estado sem conteúdo.

**Integração**

- `GET /api/v1/users/:username`.
- `GET /api/v1/users/:username/logs`, `/shelf`, `/lists`, `/favorites`,
  `/annual-challenge` e `/stats`.
- `POST /api/v1/users/:id/follow` e
  `DELETE /api/v1/users/:id/follow`.

**Bloqueio da API:** não existe endpoint para listar resenhas públicas por
usuário. Criar `GET /api/v1/users/:username/reviews` antes de habilitar a aba
Resenhas no perfil.

### 17. Seguidores e seguindo

**Objetivo:** navegar pelas relações sociais de um perfil.

**Elementos necessários**

- Duas abas: seguidores e seguindo.
- Avatar, nome de usuário e ação seguir/deixar de seguir quando autenticado.
- Estados vazios, carregamento, perfil privado e usuário não encontrado.

**Integração**

- `GET /api/v1/users/:id/followers`.
- `GET /api/v1/users/:id/following`.
- Ações de follow/unfollow da página 16.

### 18. Editar meu perfil

**Objetivo:** editar apenas os dados do perfil e os destaques pessoais.

**Elementos necessários**

- Nome público, bio, avatar, banner e localização.
- Escolha de até quatro livros favoritos.
- Configuração do desafio anual e sua visibilidade.
- Pré-visualização, validação de limites e confirmação do salvamento.

**Integração**

- `PUT /api/v1/users/me/profile`.
- `PUT /api/v1/users/me/favorites` com `bookIds`.
- `PUT /api/v1/users/me/annual-challenge`.
- Carregar dados com `GET /api/v1/users/me`.

### 19. Privacidade

**Objetivo:** oferecer controle compreensível e explícito sobre o que outras
pessoas podem ver.

**Elementos necessários**

- Controles separados para perfil, diário, progresso, estante, notas,
  comentários, listas, estatísticas e favoritos.
- Explicação do impacto de cada configuração e estado de salvamento.
- Indicadores “público” e “privado” próximos às ações de publicação.

**Integração**

- Carregar preferências por `GET /api/v1/users/me`.
- Salvar o objeto completo com `PUT /api/v1/users/me/privacy`; os nove campos
  são obrigatórios.

### 20. Sessões e segurança

**Objetivo:** permitir que a pessoa revise e encerre acessos à conta.

**Elementos necessários**

- Lista de sessões ativas, identificação visual da sessão atual e validade.
- Ações de revogar sessão individual ou todas as outras sessões.
- Confirmação para revogação em massa e encaminhamento ao login se a sessão
  atual for revogada.

**Integração**

- `GET /api/v1/auth/sessions`.
- `DELETE /api/v1/auth/sessions/:id`.
- `DELETE /api/v1/auth/sessions` revoga todas as outras e preserva a atual.
- `PUT /api/v1/auth/password` para alterar senha.

**Limite atual:** a API não registra nome do dispositivo, navegador ou
localização da sessão. Não desenhar esses dados como se fossem fornecidos.

### 21. Dados e encerramento da conta

**Objetivo:** permitir acesso aos dados exportados e exclusão consciente da
conta.

**Elementos necessários**

- Explicação do conteúdo da exportação e ação para solicitá-la/baixá-la.
- Seção separada para exclusão definitiva, com consequências, confirmação
  explícita e proteção contra clique acidental.
- Estado de processamento, sucesso e erro.

**Integração**

- `GET /api/v1/users/me/data-export`.
- Antes de excluir, exibir confirmação; executar `DELETE /api/v1/users/me`.
- Registrar/atualizar consentimentos aplicáveis por
  `POST /api/v1/users/me/consent`.

### 22. Importar histórico de leitura

**Objetivo:** revisar os registros que serão importados antes de gravá-los.

**Elementos necessários**

- Fonte/arquivo, pré-visualização, correspondência com livros do catálogo,
  validação de datas/status e resumo do lote.
- Mostrar que a importação é atômica: se houver erro, nenhum item do lote será
  gravado.
- Resultado com quantidade e leituras importadas.

**Integração atual**

- `POST /api/v1/users/me/reading-history/import` com 1–500 entradas que já
  tenham `bookId` válido.
- `GET /api/v1/users/me/reading-history/import/:jobId` para consultar o
  resultado.

**Bloqueio da API:** antes de desenhar upload de serviços externos, definir
 formato aceito, resolução de ISBN/título para `bookId` e tratamento de
 correspondências ambíguas.

## Componentes e estados que devem ser desenhados

Esses componentes são compartilhados pelas 22 páginas e não aumentam a
contagem:

- Navegação principal para desktop e mobile, com acesso a início, feed, busca,
  biblioteca, listas e perfil.
- Cartão de livro, cartão de leitura, cartão de resenha e avatar.
- Avaliação opcional por estrelas; nunca obrigar avaliação ao registrar uma
  sessão.
- Etiquetas visíveis de privado, público, rascunho e spoiler.
- Formulários com validação, salvamento, erro recuperável e confirmação.
- Estados padrão de carregamento, falha, vazio, conteúdo privado e recurso
  indisponível.
- Diálogo de confirmação para excluir leitura, nota, resenha, lista ou conta.
- Notificações/toasts acessíveis e feedback de operações idempotentes.

## Ações em modal ou painel, não páginas adicionais

- Criar/editar sessão de leitura.
- Registrar progresso diário.
- Criar/editar nota ou citação.
- Publicar/despublicar resenha.
- Adicionar livro à estante/lista ou mover entre estantes.
- Confirmar exclusões.
- Selecionar livros favoritos.
- Confirmar revogação de sessões.

## Páginas futuras — desenhar somente após ampliar a API

Estas áreas podem entrar numa fase posterior; não devem ser simuladas com
botões sem ação:

- Recomendações, busca avançada, ranking real e afinidade entre leitores.
- Comunidades/clubes.
- Central de notificações.
- Denúncia, moderação e silenciamento de usuários.
- Seguir, compartilhar, duplicar e curar listas; arquivar e reordenar listas.
- Edição/exclusão de comentários, respostas e curtidas em comentários.
- Filtros avançados e períodos adicionais de estatística.
- Página de histórico/auditoria de alterações de privacidade.

## Critérios de entrega para o designer

Para cada uma das 21 páginas centrais — e para a página 22 quando ela entrar
no escopo — entregar:

1. Wireframe desktop e mobile.
2. Estados carregando, sucesso, vazio, erro e acesso negado quando aplicável.
3. Hierarquia, campos, botões, validação e comportamento de navegação.
4. Anotação dos endpoints/ações que alimentam cada bloco.
5. Distinção visual consistente entre conteúdo privado, público, rascunho e
   publicado.
6. Confirmações para operações destrutivas e feedback após chamadas da API.

Não presumir dados fictícios ou capacidades não descritas neste guia. Se uma
tela depender de endpoint marcado como bloqueio, documentar a dependência e
deixar a ação indisponível até a API estar pronta.
