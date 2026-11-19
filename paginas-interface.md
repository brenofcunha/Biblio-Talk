# Páginas da interface - implementação atual e lacunas

Este documento compara a interface do frontend React com a experiência descrita em
[`ideologia-aplica.md`](./ideologia-aplica.md). O mapeamento foi feito a partir
de [`frontend/App.tsx`](./frontend/App.tsx) e considera como página ou modal
todo estado principal de navegação da aplicação.

## Legenda

- **Implementada**: existe uma tela ou estado funcional no frontend.
- **Parcial**: existe visualmente ou navega, mas não cobre o fluxo completo.
- **Faltante**: a página não existe ou é apenas um rótulo sem fluxo funcional.
- **Reformular**: existe, mas não representa adequadamente o objetivo da
  aplicação e precisa de alterações estruturais.

## Páginas já implementadas

| Página | Estado no frontend | Implementação atual | Avaliação |
|---|---|---|---|
| Login e cadastro | Tela única com alternância de modo | Email/senha e nome de usuário em `AuthScreen` | Parcial |
| Feed inicial | `dashboard` | Boas-vindas, meta fixa em 0/24, catálogo em alta e atalhos para diário/estantes | Reformular |
| Feed social/comunidade | `feed` | Seguidos, atividade global e resenhas; busca de livros; cursor no feed seguido | Parcial |
| Detalhes do livro | `detail` | Capa, descrição, metadados, ISBN, gêneros, avaliação, resenhas vazias e registro | Reformular |
| Perfil público | `profile` | Nome, bio, avatar, favoritos, abas, estatísticas e diário | Reformular |
| Diário | `diary` | Lista das entradas carregadas pela API e contagem | Parcial |
| Estantes | `shelves` | Colunas para quero ler, lendo, lido e abandonado | Parcial |
| Listas | `list` | Criação, seleção, privacidade exibida e busca para adicionar livros | Parcial |
| Resultados da busca | Estado `searchState` | Busca global por título/autor e seleção de livro | Parcial |
| Modal de registro de leitura | `LogModal` | Seleção do livro, data, avaliação, resenha, spoilers, estante e lista | Reformular |

## Componentes compartilhados existentes

| Componente/estado | Uso | Estado |
|---|---|---|
| Cabeçalho desktop | Busca, navegação, registrar e sair | Parcial |
| Navegação mobile | Feed, comunidade, registrar, listas e perfil | Parcial |
| Capa do livro | Feed, catálogo, perfil, listas e detalhes | Implementada |
| Avaliação por estrelas | Detalhes, feed e modal | Parcial |
| Estados vazios | Feed, diário, perfil, estantes e busca | Implementada visualmente |
| Skeleton/erro/tentativa novamente | Feed social | Implementada parcialmente |

## Páginas que faltam

### Prioridade P0

| Página/ferramenta | Objetivo | Motivo da criação |
|---|---|---|
| Configurações de privacidade | Controlar perfil, diário, progresso, estante, notas, comentários, listas e estatísticas | É um princípio central da ideologia |
| Área exclusiva de notas privadas | Criar, editar, pesquisar e excluir notas por livro/página/capítulo | O diário não pode misturar nota privada e publicação |
| Editor de resenha separada | Rascunho, edição, publicação, retirada da publicação e controle de interação | Resenha não pode ser inseparável da leitura |
| Painel de progresso da leitura | Página atual, percentual, datas, releitura e registro diário/semanal | O objetivo inclui acompanhar uma leitura, não apenas registrá-la |
| Histórico de status | Ver entradas quero ler, lendo, lido e abandonado | Permite entender mudanças no estado da leitura |
| Confirmação de exclusão/alteração | Excluir anotações, resenhas, leituras e itens com confirmação | Evita perda acidental e separa ações irreversíveis |

### Prioridade P1

| Página/ferramenta | Objetivo | Motivo da criação |
|---|---|---|
| Descoberta/recomendações | Gêneros, livros semelhantes, próximo livro e usuários compatíveis | A rede precisa ajudar a encontrar novas leituras |
| Ranking da comunidade | Ranking real baseado em leituras e avaliações | Substitui tendências fictícias por dados reais |
| Comunidades ou clubes | Criar, entrar e organizar listas coletivas | A ideologia prevê grupos de leitores |
| Perfil público de outro usuário | Exibir apenas dados permitidos e seguidores reais | Atualmente o perfil parece usar o usuário atual em vez de uma URL pública completa |
| Configurações de perfil | Banner, descrição, cidade/região, favoritos e desafio anual | Completar o perfil público |
| Exportação e exclusão de dados | Solicitar exportação e confirmar exclusão da conta | Privacidade e consentimento |
| Recuperação de senha | Solicitar e confirmar redefinição | Fluxo essencial de conta |
| Configuração de notificações | Notificações de seguidores, comentários e interações | O backend também precisará deste recurso |
| Central de denúncias e moderação | Denunciar abuso e acompanhar situação | Protege comunidade e usuários |
| Central de listas públicas | Buscar, seguir, duplicar e curar listas | Listas são experiência social e não apenas armazenagem pessoal |

### Prioridade P2

| Página/ferramenta | Objetivo | Motivo da criação |
|---|---|---|
| Importação de histórico | Selecionar arquivo/serviço e acompanhar importação | Permite migrar leituras existentes |
| Metas mensais e anuais | Criar e acompanhar objetivos | Completa desafios e estatísticas |
| Sessões ativas e segurança | Ver dispositivos/sessões e revogar acessos | Segurança da conta |
| Visualização de citações | Página ou painel por livro, nota e capítulo | Anotações devem apoiar citações referenciais |
| Arquivo de rascunhos | Resenhas e notas não finalizadas | Evita confundir conteúdo privado com publicado |
| Página de busca avançada | Autor, ISBN, gênero e período | Complementa a busca universal atual |

## Páginas que devem ser reformuladas

### 1. Tela inicial / Feed inicial

**Estado atual**

- Nomeia internamente `dashboard`, mas também é chamado de Feed pela navegação.
- A meta é exibida como `0/24` sem permitir edição.
- Usa textos e estados que devem ser completados por dados reais.
- Mistura catálogo, descoberta e resumos pessoais em uma única página.

**Reformulação sugerida**

Transformar a página em uma separação clara de duas áreas:

- **Minha leitura**: livro atual, progresso, próxima ação e objetivos;
- **Descoberta**: recomendações baseadas em gêneros e leituras públicas.

Mostrar no mesmo estado:

- meta anual configurável;
- progresso real da leitura;
- estantes com contagens;
- atalhos privados para diário e notas;
- indicadores inequívocos de privado/público.

Não usar mais `0/24` ou textos fixos como substituição de dados; caso não
exista histórico, manter o estado vazio orientativo.

### 2. Feed social/comunidade

**Estado atual**

- Possui modos seguindo/global/resenhas.
- A modalidade global e resenhas não têm paginação como o feed seguido.
- Não mostra filtros por gênero, autor ou período.
- Não permite silenciar, denunciar ou configurar compartilhamento.
- O card mostra "registrou uma leitura", mesmo em caso de resenha.

**Reformulação sugerida**

Dividir claramente em três visões:

1. **Minha rede**;
2. **Resenhas públicas**;
3. **Descobertas da comunidade**.

Cada publicação deve distinguir:

- leitura privada;
- avaliação;
- resenha;
- nota publicada isoladamente;
- data de publicação e edição.

Adicionar:

- filtros por gênero, autor e período;
- ocultação de spoilers;
- silenciar usuário;
- denunciar conteúdo;
- opções de curtida/comentário conforme privacidade;
- status público de quero ler/lendo/terminado.

### 3. Página única do livro

**Estado atual**

- Só apresenta detalhes, média, histograma, resenhas vazias e ação de registro.
- `quick-rate` permite selecionar nota sem estar ligado a uma leitura persistida.
- Não possui abas para componentes privados e públicos.

**Reformulação sugerida**

Criar cinco abas na mesma página:

| Aba | Conteúdo |
|---|---|
| Detalhes | Capa, descrição, metadados, ISBN, gêneros e comunidade |
| Minha leitura | Status, progresso, datas, percentual e releitura |
| Minhas notas | Notas privadas, páginas/capítulos, citações e busca interna |
| Minha resenha | Rascunho, publicação, spoilers e configuração de comentários |
| Comunidade | Resenhas públicas, avaliações e comparações |

Incluir na página:

- seleção ou movimento de estante;
- indicadores claros de privado/público;
- atalhos começar, atualizar progresso, terminar e abandonar;
- estado vazio para resenhas;
- confirmação antes de publicar.

### 4. Modal de registrar leitura

**Estado atual**

- Inicializa com um livro já escolhido em várias situações.
- Mistura seleção, avaliação, resenha, estante e lista em uma ação única.
- Assume `data de término` mesmo antes de o usuário indicar que terminou.
- Usa avaliação obrigatória no objeto enviado pela API.
- Não diferencia rascunho de publicação.

**Reformulação sugerida**

Substituir o modal único por fluxo em etapas:

1. **Buscar e selecionar o livro**;
2. **Escolher status**;
3. **Informar progresso e datas**;
4. **Adicionar nota privada opcional**;
5. **Decidir avaliação**;
6. **Criar ou publicar resenha opcional**;
7. **Organizar estante e listas**;
8. **Confirmar conteúdo privado e público**.

O último passo deve confirmar explicitamente o que será compartilhado.
O usuário deve conseguir salvar a leitura sem avaliação ou sem resenha.

### 5. Perfil

**Estado atual**

- A página usa dados do usuário atual.
- O botão editar não tem uma ação implementada.
- Abas Diário, Resenhas, Listas e Estatísticas usam a mesma fonte `logs`.
- A aba Estatísticas é exibida inicialmente.
- Não há configuração de privacidade, banner ou desafios.

**Reformulação sugerida**

Criar uma página pública e uma área privada separada:

- perfil público por `/users/:username`;
- painel privado em configurações;
- abas reais para Diário, Resenhas, Listas, Estantes e Estatísticas;
- contagem de seguidores;
- listas públicas;
- estatísticas públicas ou privadas;
- edição de banner, bio, região e favoritos;
- desafio anual configurável;
- exportação e exclusão de dados.

Quando houver privacidade restrita, esconder:

- diário;
- progresso;
- estante;
- notas;
- comentários;
- listas privadas;
- estatísticas.

### 6. Diário

**Estado atual**

- Exibe logs em uma lista simples.
- Não separa leitura, progresso, anotação e resenha.
- Não possui edição ou exclusão individual.
- Não permite pesquisa ou filtros.

**Reformulação sugerida**

Criar o diário como biblioteca pessoal:

- filtros por ano, status, gênero e estante;
- leitura atual em destaque;
- linha do tempo por entrada;
- progresso e página atual;
- área independente de notas;
- rascunhos;
- ação editar/excluir;
- indicador privado em cada registro;
- ação explícita de publicar apenas uma resenha.

### 7. Estantes

**Estado atual**

- As quatro estantes são seções separadas.
- Apenas são mostradas contagens como texto.
- Não são exibidos livros reais dentro das seções.
- Não há mover, remover ou editar status.
- Não há visibilidade individual.

**Reformulação sugerida**

Criar quatro abas ou cartões com contagem real e lista de capas.

Cada livro deve permitir:

- mover de estante;
- remover da estante;
- iniciar, concluir ou abandonar;
- atualizar progresso;
- abrir notas;
- indicar privacidade.

Também é necessário estado vazio real em todas as estantes.

### 8. Listas

**Estado atual**

- Cria listas por título.
- Exibe cards das listas e busca para adicionar livros.
- Não mostra os itens efetivamente carregados de uma lista.
- Não permite reordenar, remover individualmente ou editar notas.
- Não diferencia seguidores, curadores ou privacidade além do rótulo.

**Reformulação sugerida**

Criar uma página de coleção da lista:

- cabeçalho com título, descrição e privacidade;
- itens ordenados com capas;
- progresso e status de cada leitura;
- notas pessoais por item;
- reordenação por arrastar ou botões;
- remover com confirmação;
- arquivar/excluir;
- seguir/duplicar para listas públicas;
- curadores opcionais.

### 9. Resultados da busca

**Estado atual**

- Busca apenas título/autor pela interface atual.
- Não permite filtros por ISBN ou gênero.
- O resultado usa uma tela substituta, sem URL própria.
- Não possui estado persistente ou histórico.

**Reformulação sugerida**

Criar uma página ou diálogo persistente de busca avançada com:

- título;
- autor;
- ISBN;
- gênero;
- ano;
- usuário/autor literário;
- sugestão da próxima leitura.

Ao abrir um livro, continuar para a página única do livro com abas.

## Páginas novas complementares

### Recomendados

- Página de recomendações;
- Página de ranking;
- Página de comunidades;
- Página de um clube/comunidade;
- Página pública de lista;
- Página de seguidores e seguidos;
- Página de notificações;
- Página de edição de perfil;
- Página de privacidade;
- Página de segurança da conta;
- Página de metas;
- Página de histórico de privacidade;
- Página de exportação/exclusão de dados.

### Modais

- Publicar ou retirar resenha;
- Confirmar exclusão;
- Selecionar estante;
- Adicionar ou mover livro;
- Atualizar progresso;
- Criar citação;
- Denunciar conteúdo;
- Silenciar usuário;
- Compartilhar lista.

## Priorização da reformulação

### Fase 1 - separação obrigatória

1. Reformular o modal de registro em etapas.
2. Criar o editor e rascunho de resenha.
3. Criar a área privada de notas.
4. Reformular a página do livro com abas.
5. Criar configurações de privacidade.

### Fase 2 - experiência da biblioteca

1. Reformular o diário.
2. Reformular as estantes.
3. Reformular as listas.
4. Criar progresso e histórico de status.
5. Reformular o perfil.

### Fase 3 - rede social

1. Reformular o feed.
2. Criar descoberta e recomendações.
3. Criar ranking real.
4. Criar perfil público completo.
5. Criar comunidades e clubes.

### Fase 4 - conta e segurança

1. Criar configurações de perfil e privacidade.
2. Criar recuperação de senha e segurança.
3. Criar exportação e exclusão.
4. Criar notificações e denúncias.
5. Cobrir os fluxos com testes responsivos e de vazamento visual.

## Critérios mínimos de aceite

Uma reformulação será considerada consistente com a ideologia quando:

- diário, progresso e notas forem privados por padrão;
- nenhuma nota privada aparecer no feed;
- publicação de resenha depender de ação explícita;
- a interface mostrar claramente conteúdo privado/público;
- a página do livro centralizar detalhes, progresso, notas e resenha;
- estantes e listas tiverem ações reais de edição e organização;
- o perfil respeitar privacidade;
- dados fictícios forem substituídos por estados vazios ou dados da API;
- toda ação sensível tiver confirmação;
- todos os estados principais forem responsivos em celular.
