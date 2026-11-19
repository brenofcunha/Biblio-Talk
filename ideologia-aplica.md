# Ideologia e formato da aplicação

A aplicação é uma rede social voltada a leitores que permite organizar o que desejam ler, o que estão lendo e o que já leram.

As anotações e os registros pessoais devem ajudar o usuário a acompanhar sua leitura sem serem compartilhados por padrão. Compartilhar experiências literárias, avaliações e resenhas deve ser uma escolha opcional.

## Princípios principais

- O diário pessoal é privado por padrão.
- Anotações podem existir sem publicação.
- Resenhas são públicas somente quando o usuário optar por publicá-las.
- Livros podem ser organizados em estantes pessoais.
- Listas podem ser pessoais ou públicas.
- Avaliações e experiências devem poder ser compartilhadas na rede.
- Privacidade deve ser visível e controlável em cada conteúdo.
- Dados reais devem substituir dados fictícios em todas as áreas.
- A rede social complementa a biblioteca pessoal; ela não substitui o diário privado.

## O que já existe

- Catálogo integrado com Open Library e Google Books.
- Busca por título e autor.
- Cadastro e login de usuários.
- Registro de leituras no diário.
- Avaliação por estrelas.
- Resenha associada à leitura.
- Marcação de spoilers.
- Estantes:
  - Quero ler;
  - Lendo;
  - Lido;
  - Abandonado.
- Criação de listas pessoais.
- Adição de livros às listas.
- Perfis públicos básicos.
- Diário pessoal.
- Feed social.
- Seguidores e usuários seguidos.
- Curtidas e comentários.
- Estatísticas anuais de leitura.
- Template de busca no cabeçalho, no Feed e no modal de registro.

## Experiência ideal

### Minha biblioteca

O usuário deve organizar os livros em:

- Quero ler;
- Estou lendo;
- Já li;
- Abandonados.

Cada livro pode ter:

- progresso de leitura;
- página atual;
- data de início;
- data de término;
- avaliação;
- indicação de releitura;
- classificação pessoal;
- notas privadas.

### Anotações privadas

As anotações devem funcionar como um diário privado por padrão.

Cada leitura pode conter:

- notas por página ou capítulo;
- citações;
- ideias destacadas;
- impressões não finalizadas;
- palavras-chave;
- reflexões pessoais;
- exportação ou edição posterior.

O usuário nunca deve ser obrigado a publicar uma nota.

### Resenhas opcionais

O sistema deve separar claramente:

- leitura;
- anotação privada;
- avaliação;
- resenha opcional;
- publicação pública.

Somente o conteúdo explicitamente publicado deve aparecer no feed.

Cada resenha deve permitir:

- salvar como rascunho;
- informar spoilers;
- publicar ou retirar a publicação;
- editar depois da publicação;
- controlar comentários e curtidas;
- manter a leitura privada.

### Rede social

A rede deve servir para:

- seguir leitores;
- descobrir leituras públicas;
- visualizar resenhas;
- comparar avaliações;
- entrar em grupos ou comunidades;
- encontrar usuários com interesses semelhantes;
- comentar resenhas;
- criar listas públicas;
- receber recomendações literárias.

Cada usuário deve controlar o que é público:

- perfil;
- diário;
- progresso;
- estante;
- notas;
- comentários;
- listas;
- estatísticas.

## Listas

Listas podem ser privadas ou públicas.

Exemplos:

- próximos livros;
- leituras do mês;
- autores para conhecer;
- leituras concluídas;
- temas semelhantes;
- desafios anuais;
- recomendações para outros leitores.

Uma lista deve permitir:

- título e descrição;
- ordem dos livros;
- data de criação;
- capas;
- notas por item;
- privacidade;
- compartilhamento;
- seguidores;
- curadoria opcional.

## O que falta implementar

### Privacidade e compartilhamento

- Marcar diário como privado ou público.
- Separar notas privadas de resenhas públicas.
- Definir visibilidade por estante.
- Definir visibilidade do perfil e das estatísticas.
- Impedir que o feed exponha notas privadas.
- Permitir excluir permanentemente uma anotação.
- Registrar quem pode visualizar cada conteúdo.
- Implementar consentimento e exclusão de dados.

### Leituras e progresso

- Página atual e percentual concluído.
- Número de páginas da leitura.
- Datas de início e término.
- Status atual separado da estante.
- Histórico de alterações de status.
- Progresso diário ou semanal.
- Indicação de releitura.
- Tempo estimado de leitura.
- Objetivos mensais e anuais.
- Importação de histórico de leitura.

### Anotações

- Criar notas vinculadas ao livro.
- Separar nota por página ou capítulo.
- Permitir citações com referência.
- Editar e organizar notas.
- Marcar notas como rascunho.
- Publicar uma anotação isoladamente.
- Transformar uma anotação privada em resenha.
- Pesquisar dentro das próprias notas.
- Criar tags e categorias pessoais.
- Exportar em Markdown ou PDF.

### Registro de leitura

- Mostrar a busca antes da seleção do livro.
- Confirmar visualmente o livro selecionado.
- Permitir registrar uma leitura sem avaliação.
- Permitir registrar nota enquanto o livro está sendo lido.
- Permitir editar ou excluir leituras.
- Evitar registros duplicados.
- Tornar o compartilhamento uma ação final e opcional.
- Permitir escolher estante e lista na finalização.

### Estantes e listas

- Padronizar os status entre frontend e API.
- Exibir contagens reais nas estantes.
- Permitir mover um livro entre estantes.
- Permitir ordenar livros dentro das listas.
- Permitir notas pessoais em itens da lista.
- Permitir duplicar ou seguir listas públicas.
- Mostrar progresso da leitura dentro da lista.
- Permitir arquivar ou excluir listas.
- Confirmar antes de remover um livro.
- Permitir criar uma lista vazia e preenchê-la depois.

### Resenhas

- Separar resenha de nota privada no banco de dados.
- Permitir salvar rascunho.
- Permitir editar resenhas após publicação.
- Permitir remover uma publicação sem apagar o diário.
- Filtrar conteúdo com spoilers.
- Permitir comentários e respostas.
- Mostrar data de publicação e edição.
- Permitir ocultar curtidas em perfis privados.

### Feed

- Separar visualmente feed social, diário pessoal e resenhas.
- Impedir publicação automática de notas privadas.
- Permitir configurar o que é compartilhado.
- Filtrar por gênero, autor e período.
- Permitir silenciar usuários.
- Permitir denunciar abuso.
- Remover dados fictícios ou exemplos não reais.
- Preparar recomendações baseadas em leituras públicas.
- Mostrar se um usuário quer ler, está lendo ou já terminou o livro.

### Perfis

- Banner e descrição opcional.
- Cidade ou região, quando o usuário desejar.
- Livros favoritos.
- Desafio anual.
- Estatísticas públicas.
- Estantes públicas ou privadas.
- Lista de listas públicas.
- Contagem de seguidores.
- Configuração de perfil privado.
- Exportação dos dados do usuário.

### Descoberta

- Recomendações por gêneros selecionados.
- Livros semelhantes aos já lidos.
- Ranking real baseado em leituras da comunidade.
- Usuários com gosto literário parecido.
- Clubes ou listas coletivas.
- Pesquisa por autor, ISBN e gênero.
- Sugestão da próxima leitura.

### Backend

- Entidades explícitas para sessões de leitura, notas e itens de lista.
- Campos de privacidade para logs, listas e estantes.
- Validação centralizada de status.
- Permissão de edição somente pelo proprietário.
- Endpoints de rascunho e publicação.
- Paginação em feeds e listas.
- Auditoria de alterações de privacidade.
- Limites de tamanho para notas e resenhas.
- Testes de propriedade e vazamento de dados.
- Índices para busca textual e filtros sociais.
- Notificações para seguidores e comentários.

### Interface

- Tela única do livro com abas:
  - detalhes;
  - progresso pessoal;
  - minhas notas;
  - resenha pública;
  - comunidade.
- Indicadores claros de conteúdo privado e público.
- Estado vazio em todas as estantes.
- Modal dividido em etapas.
- Confirmação antes de publicar uma resenha.
- Área independente para notas.
- Atalhos para começar, terminar e abandonar uma leitura.
- Indicador visual de salvamento e sincronização.
- Responsividade completa para celular.

## Diretriz final

Leitura, progresso e notas devem ser privados por padrão. Resenhas, recomendações e experiências literárias devem ser compartilhadas apenas quando o usuário escolher fazê-lo explicitamente.
