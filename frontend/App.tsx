import { useEffect, useRef, useState, type FormEvent, type ReactNode } from 'react'
import { addBookToList, clearToken, createReadingList, createReadingLog, getAllUserShelves, getBook, getCurrentUser, getFeed, getGlobalFeed, getMoreFeed, getMyLists, getPopularBooks, getReviewFeed, getProfile, getUserLogs, getUserStats, hasToken, login, logout, register, saveToken, searchBooks, updateShelf, type Book, type BookDetails, type FeedItem, type ReadingList, type ReadingLog, type ShelfEntry, type User } from './api'

type Screen = 'dashboard' | 'feed' | 'detail' | 'profile' | 'list' | 'diary' | 'shelves'

type AvatarUser = { name: string; initials: string; tone: 'coral' | 'lime' | 'blue' }

function avatarFor(name: string): AvatarUser {
  const initials = name.trim().split(/\s+/).slice(0, 2).map(part => part[0]).join('').toUpperCase() || '?'
  const tones: AvatarUser['tone'][] = ['coral', 'lime', 'blue']
  const tone = tones[Array.from(name).reduce((sum, char) => sum + char.charCodeAt(0), 0) % tones.length]
  return { name, initials, tone }
}

function Icon({ name, size = 18 }: { name: string; size?: number }) {
  const paths: Record<string, ReactNode> = {
    search: <><circle cx="11" cy="11" r="7" /><path d="m20 20-4-4" /></>,
    home: <><path d="m3 11 9-8 9 8" /><path d="M5 10v10h14V10M9 20v-6h6v6" /></>,
    compass: <><circle cx="12" cy="12" r="9" /><path d="m15 9-2 4-4 2 2-4z" /></>,
    list: <><path d="M8 6h13M8 12h13M8 18h13" /><circle cx="3" cy="6" r=".5" /><circle cx="3" cy="12" r=".5" /><circle cx="3" cy="18" r=".5" /></>,
    heart: <path d="M20.8 4.7a5.5 5.5 0 0 0-7.8 0L12 5.8l-1.1-1.1a5.5 5.5 0 0 0-7.8 7.8l1.1 1.1L12 21l7.8-7.4a5.5 5.5 0 0 0 1-8.9Z" />,
    comment: <path d="M21 15a4 4 0 0 1-4 4H8l-5 3V7a4 4 0 0 1 4-4h10a4 4 0 0 1 4 4Z" />,
    more: <><circle cx="5" cy="12" r="1" /><circle cx="12" cy="12" r="1" /><circle cx="19" cy="12" r="1" /></>,
    plus: <path d="M12 5v14M5 12h14" />,
    book: <><path d="M4 4h11a3 3 0 0 1 3 3v13H7a3 3 0 0 1-3-3Z" /><path d="M7 4v16" /></>,
    bookmark: <path d="M6 3h12v18l-6-4-6 4z" />,
    close: <path d="m6 6 12 12M18 6 6 18" />,
    calendar: <><rect x="3" y="5" width="18" height="16" rx="2" /><path d="M16 3v4M8 3v4M3 10h18" /></>,
    user: <><circle cx="12" cy="8" r="4" /><path d="M4 21a8 8 0 0 1 16 0" /></>,
    arrow: <path d="m9 18 6-6-6-6" />,
    check: <path d="m5 12 4 4L19 6" />,
    pages: <><path d="M6 3h10l3 3v15H6z" /><path d="M16 3v4h4M9 12h6M9 16h6" /></>,
    refresh: <><path d="M20 11a8 8 0 0 0-14.9-3M4 5v4h4" /><path d="M4 13a8 8 0 0 0 14.9 3M20 19v-4h-4" /></>,
  }
  return <svg aria-hidden="true" width={size} height={size} viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.8" strokeLinecap="round" strokeLinejoin="round">{paths[name]}</svg>
}

function AuthScreen({ onAuthenticated }: { onAuthenticated: (result: { user: User; token: string }) => void }) {
  const [mode, setMode] = useState<'login' | 'register'>('login')
  const [username, setUsername] = useState('')
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(false)

  async function submit(event: FormEvent) {
    event.preventDefault()
    setError('')
    setLoading(true)
    try {
      const result = mode === 'login'
        ? await login(email, password)
        : await register(username, email, password)
      saveToken(result.token)
      onAuthenticated(result)
    } catch (requestError) {
      setError(requestError instanceof Error ? requestError.message : 'Não foi possível autenticar.')
    } finally {
      setLoading(false)
    }
  }

  return <main className="auth-page">
    <section className="auth-visual">
      <button className="brand large" type="button"><span className="brand-mark"><span /></span><span>Lume</span></button>
      <div><p className="eyebrow">Sua história como leitor</p><h1>Todo livro deixa uma marca.</h1><p>Guarde suas leituras, descubra novas histórias e encontre pessoas que sentem os livros como você.</p></div>
      <blockquote>“A leitura é uma forma de estar em muitos lugares sem deixar de estar em si.”<span>— Ana Martins Marques</span></blockquote>
    </section>
    <section className="auth-panel"><div className="auth-box"><p className="eyebrow">{mode === 'login' ? 'Bem-vindo de volta' : 'Comece seu diário'}</p><h2>{mode === 'login' ? 'Entre na sua conta' : 'Crie sua conta'}</h2><p>{mode === 'login' ? 'Seu diário e suas próximas leituras estão esperando.' : 'É gratuito. Você só precisa de um minuto.'}</p>
      <form className="auth-form" onSubmit={submit}>
        {mode === 'register' && <label><span>Nome de usuário</span><input required value={username} onChange={event => setUsername(event.target.value)} /></label>}
        <label><span>E-mail</span><input required type="email" value={email} onChange={event => setEmail(event.target.value)} /></label>
        <label><span>Senha</span><input required minLength={8} type="password" value={password} onChange={event => setPassword(event.target.value)} /></label>
        {error && <p className="form-error" role="alert">{error}</p>}
        <button className="primary-button wide" disabled={loading}>{loading ? 'Aguarde...' : mode === 'login' ? 'Entrar no Lume' : 'Criar minha conta'}</button>
      </form>
      <p className="auth-switch">{mode === 'login' ? 'Ainda não tem uma conta?' : 'Já faz parte do Lume?'} <button type="button" onClick={() => { setMode(mode === 'login' ? 'register' : 'login'); setError('') }}>{mode === 'login' ? 'Cadastre-se' : 'Entrar'}</button></p>
    </div></section>
  </main>
}

function Avatar({ user, small = false, onClick }: { user: AvatarUser; small?: boolean; onClick?: () => void }) {
  return <button className={`avatar ${user.tone} ${small ? 'avatar-sm' : ''}`} onClick={onClick} aria-label={`Ver perfil de ${user.name}`}>{user.initials}</button>
}

function Brand({ onClick }: { onClick: () => void }) {
  return <button className="brand" onClick={onClick} aria-label="Ir para o início"><span className="brand-mark"><span /></span><span>Lume</span></button>
}

function Header({ screen, navigate, openLog, onSearch, onLogout }: { screen: Screen; navigate: (s: Screen) => void; openLog: () => void; onSearch: (query: string) => void; onLogout: () => void }) {
  const [query, setQuery] = useState('')
  const searchInput = useRef<HTMLInputElement>(null)
  return <header className="topbar">
    <div className="topbar-inner">
      <Brand onClick={() => navigate('feed')} />
      <form className="search" onSubmit={event => { event.preventDefault(); if (query.trim()) onSearch(query.trim()) }}><Icon name="search" /><input ref={searchInput} aria-label="Busca universal" placeholder="Busque livros ou autores" value={query} onChange={event => setQuery(event.target.value)} /></form>
      <nav className="main-nav" aria-label="Navegação principal">
        <button className={screen === 'dashboard' ? 'active' : ''} onClick={() => navigate('dashboard')}>Feed</button>
        <button className={screen === 'feed' ? 'active' : ''} onClick={() => navigate('feed')}>Comunidade</button>
        <button className={screen === 'shelves' ? 'active' : ''} onClick={() => navigate('shelves')}>Estantes</button>
        <button className={screen === 'list' ? 'active' : ''} onClick={() => navigate('list')}>Listas</button>
        <button className={screen === 'diary' ? 'active' : ''} onClick={() => navigate('diary')}>Diário</button>
        <button onClick={() => { navigate('dashboard'); searchInput.current?.focus() }}>Explorar</button>
      </nav>
      <button className="button primary compact desktop-log" onClick={openLog}><Icon name="plus" /> Registrar</button>
      <Avatar user={avatarFor('Você')} small onClick={() => navigate('profile')} />
      <button className="icon-button" onClick={onLogout} aria-label="Sair">Sair</button>
    </div>
  </header>
}

function MobileNav({ screen, navigate, openLog }: { screen: Screen; navigate: (s: Screen) => void; openLog: () => void }) {
  return <nav className="mobile-nav" aria-label="Navegação mobile">
    <button className={screen === 'dashboard' ? 'active' : ''} onClick={() => navigate('dashboard')}><Icon name="home" /><span>Feed</span></button>
    <button className={screen === 'feed' ? 'active' : ''} onClick={() => navigate('feed')}><Icon name="compass" /><span>Comunidade</span></button>
    <button className="mobile-add" onClick={openLog}><Icon name="plus" size={22} /></button>
    <button className={screen === 'list' ? 'active' : ''} onClick={() => navigate('list')}><Icon name="list" /><span>Listas</span></button>
    <button className={screen === 'profile' ? 'active' : ''} onClick={() => navigate('profile')}><Icon name="user" /><span>Perfil</span></button>
  </nav>
}

function Stars({ value = 4.5, interactive = false, onChange }: { value?: number; interactive?: boolean; onChange?: (value: number) => void }) {
  const [rating, setRating] = useState(value)
  return <div className={`stars ${interactive ? 'stars-interactive' : ''}`} aria-label={`${rating} de 5 estrelas`}>
    {[1, 2, 3, 4, 5].map(n => <button key={n} type="button" onClick={() => { if (interactive) { setRating(n); onChange?.(n) } }} tabIndex={interactive ? 0 : -1} aria-label={`${n} estrelas`}><span className={n <= rating ? 'filled' : ''}>★</span></button>)}
    {interactive && <strong>{rating.toFixed(1)}</strong>}
  </div>
}

function BookCover({ book, index, onClick, large = false }: { book: Book; index?: number; onClick: () => void; large?: boolean }) {
  const author = book.authors?.join(', ') || 'Autor desconhecido'
  return <div className={`book-cover-wrap ${large ? 'cover-large' : ''}`}>
    {index !== undefined && <span className="rank">{index + 1}</span>}
    <button className="book-cover" onClick={onClick} aria-label={`Ver ${book.title}`}>
      {book.coverUrl ? <img src={book.coverUrl} alt={`Capa de ${book.title}`} /> : <span className="cover-placeholder"><Icon name="book" /></span>}
      <span className="cover-shade" />
      <span className="cover-copy"><small>{author}</small><strong>{book.title}</strong></span>
      <span className="quick-actions"><span><Icon name="check" /> Lido</span><Icon name="heart" /></span>
    </button>
  </div>
}

function FeedCard({ item, navigate }: { item: FeedItem; navigate: (s: Screen) => void }) {
  const user = avatarFor(item.user.username)
  const book: Book = { id: item.bookId || item.id, title: item.book.title, authors: item.book.authors, coverUrl: item.book.coverUrl }
  return <article className="feed-card">
    <div className="activity-head">
      <div className="user-line"><Avatar user={user} small onClick={() => navigate('profile')} /><div><button onClick={() => navigate('profile')}><strong>@{item.user.username}</strong></button><p>registrou uma leitura · {new Date(item.createdAt).toLocaleDateString('pt-BR')}</p></div></div>
      <button className="icon-button" aria-label="Mais opções"><Icon name="more" /></button>
    </div>
    <div className="activity-body">
      <BookCover book={book} onClick={() => navigate('detail')} />
      <div className="review">
        <button className="book-title" onClick={() => navigate('detail')}>{book.title}</button>
        <p className="author">{book.authors?.join(', ') || 'Autor desconhecido'}</p>
        <Stars value={item.rating} />
        <p className="review-text">{item.reviewText || 'Registrou esta leitura no diário.'}</p>
        {item.containsSpoilers && <span className="spoiler">Resenha contém spoilers</span>}
      </div>
    </div>
    <div className="activity-actions"><button><Icon name="heart" /> {item.likesCount}</button><button><Icon name="comment" /> {item.commentsCount ?? 0}</button></div>
  </article>
}

function Sidebar({ navigate, popularBooks, onBookSelect }: { navigate: (s: Screen) => void; popularBooks: Book[]; onBookSelect: (book: Book) => void }) {
  return <aside className="sidebar">
    <section className="side-section">
      <div className="section-heading"><h3>Populares esta semana</h3><button onClick={() => navigate('list')}>Ver todos</button></div>
      {popularBooks.slice(0, 3).map((book, i) => <button className="mini-book" key={book.id} onClick={() => onBookSelect(book)}><span className="mini-rank">0{i + 1}</span>{book.coverUrl ? <img src={book.coverUrl} alt="" /> : <span className="mini-cover-placeholder"><Icon name="book" /></span>}<span><strong>{book.title}</strong><small>{book.authors?.join(', ') || 'Autor desconhecido'}</small><em><span>★</span> Dados da API</em></span></button>)}
      {!popularBooks.length && <p className="empty-copy">Ainda não há leituras suficientes para gerar tendências.</p>}
    </section>
  </aside>
}

function Dashboard({ user, navigate, popularBooks, featuredBook, openLog, onBookSelect }: { user: User; navigate: (s: Screen) => void; popularBooks: Book[]; featuredBook: Book | null; openLog: () => void; onBookSelect: (book: Book) => void }) {
  const catalog = popularBooks
  return <main className="page dashboard">
    <section className="dashboard-main">
      <div className="welcome-row"><div><p className="eyebrow">Sua vida em livros</p><h1>Boa noite, {user.username}.</h1><p>Continue de onde parou ou encontre sua próxima leitura.</p></div><div className="goal-card"><span>0<small>/24</small></span><div><b>Meta de leitura</b><i><em style={{ width: '0%' }} /></i></div></div></div>
      <section className="content-block"><div className="section-head"><div><p className="eyebrow">Descobertas</p><h2>Em alta na comunidade</h2></div><button className="text-link" onClick={() => navigate('feed')}>Ver comunidade</button></div>
        {catalog.length ? <div className="book-grid featured">{catalog.slice(0, 4).map((book, index) => <article className="book-card" key={book.id}><button className="book-cover" onClick={() => onBookSelect(book)}><span className="rank">{index + 1}</span>{book.coverUrl ? <img src={book.coverUrl} alt={`Capa de ${book.title}`} /> : <span className="cover-placeholder"><Icon name="book" /></span>}<span className="cover-overlay"><b>Ver livro</b><i>♡</i></span></button><button className="book-card-link" onClick={() => onBookSelect(book)}><strong>{book.title}</strong><span>{book.authors?.join(', ')}</span></button><small><b>★</b> {book.publishYear ?? '—'}</small></article>)}</div> : <section className="diary-empty"><h2>Sem livros em alta ainda</h2><p>As tendências aparecerão quando houver leituras reais na comunidade.</p></section>}
      </section>
      <section className="content-block" id="catalogo"><div className="section-head"><div><p className="eyebrow">Catálogo</p><h2>Explore todos os livros</h2></div><span>{catalog.length} títulos</span></div>
        {catalog.length ? <div className="book-grid">{catalog.map(book => <article className="book-card" key={book.id}><BookCover book={book} onClick={() => onBookSelect(book)} /><button className="book-card-link" onClick={() => onBookSelect(book)}><strong>{book.title}</strong><span>{book.authors?.join(', ')}</span></button></article>)}</div> : <p className="empty-copy">Busque um livro para adicioná-lo ao catálogo.</p>}
      </section>
      {featuredBook && <section className="content-block featured-book-template"><div className="section-head"><div><p className="eyebrow">Em alta</p><h2>{featuredBook.title}</h2></div><button className="text-link" onClick={() => onBookSelect(featuredBook)}>Ver detalhes</button></div><div className="featured-book-copy"><BookCover book={featuredBook} large onClick={() => onBookSelect(featuredBook)} /><div><p className="author">{featuredBook.authors?.join(', ') || 'Autor desconhecido'}</p><p className="review-text">Livro encontrado no catálogo externo e salvo pela API para uso no aplicativo.</p><button className="button primary compact" onClick={openLog}>Registrar leitura</button></div></div></section>}
    </section>
    <aside className="dashboard-aside">
      <div className="aside-card profile-mini"><span className="profile-avatar">{user.username.slice(0, 2).toUpperCase()}</span><div><strong>{user.username}</strong><small>@{user.username}</small></div><button onClick={() => navigate('profile')}>Ver perfil</button></div>
      <div className="aside-card"><div className="aside-title"><h3>Seu diário</h3><button onClick={() => navigate('detail')}>Abrir</button></div><div className="mini-empty"><Icon name="book" /><p>Seu diário começa com uma leitura.</p><button onClick={openLog}>Registrar leitura</button></div></div>
      <div className="aside-card"><div className="aside-title"><h3>Suas estantes</h3><button onClick={() => navigate('shelves')}>Gerenciar</button></div><p className="empty-copy">Os livros das suas estantes aparecerão aqui após o primeiro registro.</p></div>
    </aside>
  </main>
}

function Feed({ navigate, items, loading, error, onRetry, onRefresh, onLoadMore, hasMore, mode, onModeChange, popularBooks, onBookSelect }: { navigate: (s: Screen) => void; items: FeedItem[]; loading: boolean; error: string; onRetry: () => void; onRefresh: () => void; onLoadMore: () => void; hasMore: boolean; mode: 'following' | 'global' | 'reviews'; onModeChange: (mode: 'following' | 'global' | 'reviews') => void; popularBooks: Book[]; onBookSelect: (book: Book) => void }) {
  const [query, setQuery] = useState('')
  const [results, setResults] = useState<Book[]>([])
  const [searchLoading, setSearchLoading] = useState(false)
  const [searchError, setSearchError] = useState('')

  async function search(event: FormEvent) {
    event.preventDefault()
    if (!query.trim()) return
    setSearchLoading(true)
    setSearchError('')
    try {
      setResults(await searchBooks(query.trim()))
    } catch (requestError) {
      setResults([])
      setSearchError(requestError instanceof Error ? requestError.message : 'Falha ao buscar livros.')
    } finally {
      setSearchLoading(false)
    }
  }

  return <main className="page feed-layout">
    <section className="feed-column">
      <div className="feed-intro"><div><p className="eyebrow">Seu diário social</p><h1>Veja o que sua rede anda lendo.</h1><p>Atualizado diretamente pela API.</p></div><button className="button secondary compact" onClick={onRefresh} disabled={loading}><Icon name="refresh" /> {loading ? 'Atualizando...' : 'Atualizar'}</button></div>
      <form className="feed-book-search" onSubmit={search}><Icon name="search" /><input aria-label="Buscar livros no feed" value={query} onChange={event => setQuery(event.target.value)} placeholder="Buscar livros ou autores" /><button className="button primary compact" disabled={searchLoading || !query.trim()}>{searchLoading ? 'Buscando...' : 'Buscar'}</button></form>
      {searchError && <p className="form-error" role="alert">{searchError}</p>}
      {results.length > 0 && <section className="feed-search-results"><div className="section-title"><h2>Possíveis livros</h2><span>{results.length} resultados da API</span></div><div className="feed-book-grid">{results.map(book => <article className="feed-book-result" key={book.id}><button className="feed-book-cover" onClick={() => onBookSelect(book)}>{book.coverUrl ? <img src={book.coverUrl} alt={`Capa de ${book.title}`} /> : <span className="cover-placeholder"><Icon name="book" /></span>}</button><div><button className="book-title" onClick={() => onBookSelect(book)}>{book.title}</button><p className="author">{book.authors?.join(', ') || 'Autor desconhecido'}</p><p className="feed-book-description">{book.description || 'Descrição não fornecida pela fonte externa.'}</p><button className="text-link" onClick={() => onBookSelect(book)}>Ver detalhes</button></div></article>)}</div></section>}
      <div className="filter-row"><button className={mode === 'following' ? 'active' : ''} onClick={() => onModeChange('following')}>Sua rede</button><button className={mode === 'global' ? 'active' : ''} onClick={() => onModeChange('global')}>Atividade global</button><button className={mode === 'reviews' ? 'active' : ''} onClick={() => onModeChange('reviews')}>Resenhas</button></div>
      {loading && !items.length && <section className="feed-skeletons" aria-label="Carregando feed">{[1, 2, 3].map(index => <div className="feed-skeleton" key={index}><span className="skeleton-avatar" /><div><i className="skeleton-line short" /><i className="skeleton-line" /><i className="skeleton-line medium" /><div className="skeleton-book" /></div></div>)}</section>}
      {loading && items.length > 0 && <p className="feed-refreshing" role="status">Atualizando o feed...</p>}
      {error && <section className="diary-empty"><p role="alert">{error}</p><button className="button primary" onClick={onRetry}>Tentar novamente</button></section>}
      {!loading && !error && items.length === 0 && <section className="diary-empty"><h2>O feed está vazio</h2><p>Siga leitores ou registre uma leitura para começar.</p></section>}
      {!loading && !error && items.map(item => <FeedCard key={item.id} item={item} navigate={navigate} />)}
      {!loading && !error && hasMore && <button className="button secondary" onClick={onLoadMore}>Carregar mais</button>}
    </section>
    <Sidebar navigate={navigate} popularBooks={popularBooks} onBookSelect={onBookSelect} />
  </main>
}

function Detail({ book, navigate, openLog }: { book: BookDetails | null; navigate: (s: Screen) => void; openLog: () => void }) {
  if (!book) return <main className="page diary-empty"><h2>Selecione um livro</h2><p>Use a busca para carregar um livro real da API.</p></main>
  const author = book.authors?.join(', ') || 'Autor desconhecido'
  const rating = book.averageRating ?? 0
  return <main>
    <section className="book-hero">
      <div className="hero-backdrop" style={book.coverUrl ? { backgroundImage: `url(${book.coverUrl})` } : undefined} />
      <div className="page hero-content">
        <BookCover book={book} large onClick={() => {}} />
        <div className="book-meta">
          <p className="eyebrow">Dados do catálogo</p>
          <h1>{book.title}</h1>
          <p className="author-link">por {author}</p>
          <p className="book-description">{book.description || 'Descrição não fornecida pela fonte externa.'}</p>
          <div className="metadata"><span>{book.publisher || 'Editora não informada'}</span><span>{book.pageCount ? `${book.pageCount} páginas` : 'Páginas não informadas'}</span><span>{book.publishYear || 'Ano não informado'}</span></div>
          <div className="book-actions">
            <button className="button primary" onClick={openLog}><Icon name="plus" /> Registrar leitura</button>
            <button className="button icon-only" aria-label="Adicionar à lista de desejos"><Icon name="bookmark" /></button>
          </div>
          <div className="quick-rate"><span>Sua nota</span><Stars value={0} interactive /></div>
        </div>
      </div>
    </section>
    <div className="page detail-grid">
      <section>
        <div className="stats-card">
          <div className="score"><strong>{rating.toFixed(1)}</strong><Stars value={rating} /><span>{book.ratingsCount ?? 0} avaliações</span></div>
          <div className="histogram">
            {[1, 2, 3, 4, 5].map(star => <div key={star}><span>{star}</span><i><b style={{ height: `${Math.min(100, (book.histogram?.[String(star)] ?? 0) * 10)}%` }} /></i></div>)}
          </div>
        </div>
        <div className="content-section"><div className="section-title"><h2>Sobre o livro</h2></div><p>Detalhes adicionais serão exibidos quando fornecidos pela API do catálogo.</p></div>
        <div className="content-section">
          <div className="section-title"><h2>Resenhas da comunidade</h2></div>
          {!book.reviews?.length && <p className="empty-copy">Ainda não há resenhas reais para este livro.</p>}
        </div>
      </section>
      <aside className="detail-aside"><h3>Informações</h3><dl><div><dt>Publicado em</dt><dd>{book.publishYear || 'Não informado'}</dd></div><div><dt>ISBN</dt><dd>{book.isbn || 'Não informado'}</dd></div><div><dt>Gêneros</dt><dd>{book.genres?.join(', ') || 'Não informado'}</dd></div></dl></aside>
    </div>
  </main>
}

function Profile({ user, logs, stats, navigate }: { user: User; logs: ReadingLog[]; stats?: Record<string, unknown>; navigate: (s: Screen) => void }) {
  const [tab, setTab] = useState('Estatísticas')
  return <main>
    <section className="profile-hero">
      <div className="page profile-head">
        <Avatar user={avatarFor(user.username)} />
        <div className="profile-copy"><div><h1>{user.username}</h1><button className="button secondary compact">Editar perfil</button></div><p className="handle">@{user.username}</p><p>{user.bio || 'Registre suas leituras e compartilhe suas descobertas com a comunidade.'}</p><div className="profile-numbers"><span><strong>{Number(stats?.booksRead ?? logs.length)}</strong> lidos no ano</span><span><strong>{logs.length}</strong> entradas no diário</span><span><strong>{user.top4BookIds?.length ?? 0}</strong> favoritos</span></div></div>
      </div>
    </section>
    <div className="page profile-page">
      <section className="top-four"><div className="section-title"><div><p className="eyebrow">Estante pessoal</p><h2>Quatro para uma ilha deserta</h2></div><button onClick={() => navigate('shelves')}>Ver estante <Icon name="arrow" /></button></div>{logs.length ? <div className="top-four-grid">{logs.slice(0, 4).map(log => <BookCover key={log.id} book={{ id: log.bookId, ...log.book }} onClick={() => navigate('detail')} />)}</div> : <p className="empty-copy">Sua seleção aparecerá após registrar leituras reais.</p>}</section>
      <nav className="profile-tabs">{['Diário', 'Resenhas', 'Listas', 'Estatísticas'].map(t => <button className={tab === t ? 'active' : ''} key={t} onClick={() => setTab(t)}>{t}</button>)}</nav>
      {tab === 'Estatísticas' ? <Stats stats={stats} /> : <Diary navigate={navigate} tab={tab} logs={logs} />}
    </div>
  </main>
}

function Stats({ stats }: { stats?: Record<string, unknown> }) {
  const booksRead = Number(stats?.booksRead ?? 0)
  const pagesRead = Number(stats?.pagesRead ?? 0)
  const genres = Object.entries((stats?.genresPercent as Record<string, number> | undefined) ?? {})
  const monthlyBooks = Object.entries((stats?.booksPerMonth as Record<string, number> | undefined) ?? {})
  return <section className="stats-dashboard">
    <div className="stats-header"><div><p className="eyebrow">Seu ano em livros</p><h2>{new Date().getFullYear()}, até agora</h2></div></div>
    <div className="stat-kpis"><div><Icon name="book" /><span><strong>{booksRead}</strong><small>livros lidos</small></span><em>dados da API</em></div><div><Icon name="pages" /><span><strong>{pagesRead.toLocaleString('pt-BR')}</strong><small>páginas lidas</small></span><em>dados da API</em></div><div><Icon name="calendar" /><span><strong>{Number(stats?.currentStreakDays ?? 0)}</strong><small>dias em sequência</small></span><em>dados da API</em></div></div>
    <div className="chart-grid">
      <div className="chart-card"><div className="section-title"><h3>Gêneros favoritos</h3><span>Por livros lidos</span></div>{genres.length ? <ul>{genres.map(([genre, percent]) => <li key={genre}><i className="dot green" />{genre} <strong>{percent}%</strong></li>)}</ul> : <p className="empty-copy">Sem dados reais de gêneros ainda.</p>}</div>
      <div className="chart-card"><div className="section-title"><h3>Ritmo de leitura</h3><span>Livros por mês</span></div>{monthlyBooks.length ? <div className="bar-chart">{monthlyBooks.map(([month, count]) => <div key={month}><i><b style={{ height: `${Math.min(100, Number(count) * 10)}%` }} /></i><span>{month}</span></div>)}</div> : <p className="empty-copy">Sem leituras reais para exibir.</p>}</div>
    </div>
  </section>
}

function Diary({ navigate, tab, logs }: { navigate: (s: Screen) => void; tab: string; logs: ReadingLog[] }) {
  if (!logs.length) return <section className="diary-empty"><Icon name="book" size={32} /><h2>{tab} vazio</h2><p>Registre uma leitura para começar seu diário.</p><button className="button primary" onClick={() => navigate('detail')}>Explorar livro em destaque</button></section>
  return <section className="diary-list">{logs.slice(0, 10).map(log => <article className="diary-entry" key={log.id}><time>{log.readDate ? new Date(log.readDate).toLocaleDateString('pt-BR') : 'Sem data'}</time><div className="diary-cover">{log.book.coverUrl ? <img src={log.book.coverUrl} alt="" /> : <Icon name="book" />}</div><div><h2>{log.book.title}</h2><p>{log.book.authors?.join(', ')}</p><span className="rating">{'★'.repeat(Math.round(log.rating))} <b>{log.rating.toFixed(1)}</b></span>{log.reviewText && <blockquote>{log.reviewText}</blockquote>}</div></article>)}</section>
}

function ListPage({ navigate, lists, onListsChange }: { navigate: (s: Screen) => void; lists: ReadingList[]; onListsChange: () => void }) {
  const visibleLists = lists.length ? lists : []
  const [title, setTitle] = useState('')
  const [query, setQuery] = useState('')
  const [selectedListId, setSelectedListId] = useState(lists[0]?.id ?? '')
  const [results, setResults] = useState<Book[]>([])
  const [status, setStatus] = useState('')
  const [loading, setLoading] = useState(false)

  useEffect(() => {
    if (!selectedListId && lists.length) setSelectedListId(lists[0].id)
  }, [lists, selectedListId])

  return <main className="page list-page">
    <section className="list-header">
      <div><p className="eyebrow">Suas listas · {visibleLists.length}</p><h1>{visibleLists[0]?.title ?? 'Suas listas de leitura'}</h1><p>{visibleLists[0]?.description ?? 'Crie listas para organizar suas próximas leituras e recomendações.'}</p><div className="list-owner"><Avatar user={avatarFor('Você')} small onClick={() => navigate('profile')} /><span>Listas carregadas pela API Go</span></div></div>
    </section>
    <div className="list-toolbar"><div><strong>{visibleLists.length} listas</strong><span>Conteúdo da sua conta</span></div><form className="inline-list-create" onSubmit={async event => { event.preventDefault(); if (!title.trim()) return; setLoading(true); setStatus(''); try { await createReadingList(title.trim()); setTitle(''); onListsChange(); setStatus('Lista criada.') } catch (error) { setStatus(error instanceof Error ? error.message : 'Não foi possível criar a lista.') } finally { setLoading(false) } }}><input aria-label="Nome da lista" value={title} onChange={event => setTitle(event.target.value)} placeholder="Nome da nova lista" /><button className="button primary compact" disabled={loading || !title.trim()}>{loading ? 'Criando...' : 'Criar lista'}</button></form></div>
    {!visibleLists.length && <section className="diary-empty"><Icon name="list" size={32} /><h2>Nenhuma lista criada</h2><p>As listas criadas pela API aparecerão aqui.</p></section>}
    <section className="books-grid">{visibleLists.map(list => <article className="grid-book" key={list.id}><div className="book-cover-wrap"><div className="book-cover"><span className="cover-placeholder"><Icon name="list" /></span></div></div><strong>{list.title}</strong><span>{list.description || 'Lista de leitura'}</span><small>{list.isPrivate ? 'Privada' : 'Pública'}</small></article>)}</section>
    {visibleLists.length > 0 && <section className="list-book-search">
      <div className="section-title"><div><p className="eyebrow">Adicionar livros</p><h2>Buscar na API</h2></div></div>
      <form className="list-search-form" onSubmit={async event => { event.preventDefault(); if (!query.trim()) return; setLoading(true); setStatus(''); try { setResults(await searchBooks(query.trim())) } catch (error) { setStatus(error instanceof Error ? error.message : 'Falha ao buscar livros.') } finally { setLoading(false) } }}><input aria-label="Buscar livro" value={query} onChange={event => setQuery(event.target.value)} placeholder="Título ou autor" /><label><span>Adicionar à lista</span><select value={selectedListId} onChange={event => setSelectedListId(event.target.value)}>{visibleLists.map(list => <option value={list.id} key={list.id}>{list.title}</option>)}</select></label><button className="button secondary compact" disabled={loading || !query.trim()}>{loading ? 'Buscando...' : 'Buscar'}</button></form>
      {status && <p className="form-error" role="status">{status}</p>}
      <div className="books-grid">{results.map(book => <article className="grid-book" key={book.id}><div className="book-cover-wrap"><div className="book-cover">{book.coverUrl ? <img src={book.coverUrl} alt={`Capa de ${book.title}`} /> : <span className="cover-placeholder"><Icon name="book" /></span>}</div></div><strong>{book.title}</strong><span>{book.authors?.join(', ') || 'Autor desconhecido'}</span><button className="button primary compact" disabled={!selectedListId} onClick={async () => { setLoading(true); setStatus(''); try { await addBookToList(selectedListId, book.id, 1); setStatus(`“${book.title}” adicionado à lista.`); onListsChange() } catch (error) { setStatus(error instanceof Error ? error.message : 'Não foi possível adicionar o livro.') } finally { setLoading(false) } }}>Adicionar</button></article>)}</div>
    </section>}
  </main>
}

function DiaryPage({ logs }: { logs: ReadingLog[] }) {
  return <main className="page inner-page"><header className="page-title"><div><p className="eyebrow">Histórico de leitura</p><h1>Seu diário</h1><p>Leituras carregadas da API Go.</p></div><div className="diary-kpis"><span><strong>{logs.length}</strong> entradas</span></div></header><Diary navigate={() => undefined} tab="Diário" logs={logs} /></main>
}

function ShelvesPage({ shelves }: { shelves: ShelfEntry[] }) {
  return <main className="page inner-page"><header className="page-title"><div><p className="eyebrow">Organização</p><h1>Suas estantes</h1><p>Livros organizados por status na API Go.</p></div></header>{[
    ['tbr', 'Quero ler'],
    ['reading', 'Lendo'],
    ['read', 'Lidos'],
    ['abandoned', 'Abandonados'],
  ].map(([status, label]) => { const entries = shelves.filter(shelf => shelf.status === status); return <section className="shelf-section" key={status}><p className="eyebrow">{label}</p>{entries.length ? <p>{entries.length} livro(s) nesta estante.</p> : <div className="empty-state"><h3>Estante vazia</h3></div>}</section> })}</main>
}

function LogModal({ initialBook, lists, close, onSaved }: { initialBook: Book | null; lists: ReadingList[]; close: () => void; onSaved: () => void }) {
  const [book, setBook] = useState<Book | null>(initialBook)
  const [searchQuery, setSearchQuery] = useState('')
  const [searchResults, setSearchResults] = useState<Book[]>([])
  const [searchLoading, setSearchLoading] = useState(false)
  const [spoiler, setSpoiler] = useState(false)
  const [liked, setLiked] = useState(false)
  const [rating, setRating] = useState(4)
  const [date, setDate] = useState(new Date().toISOString().slice(0, 10))
  const [review, setReview] = useState('')
  const [shelfStatus, setShelfStatus] = useState('read')
  const [listId, setListId] = useState('')
  const [saving, setSaving] = useState(false)
  const [error, setError] = useState('')

  async function searchBook(event: FormEvent) {
    event.preventDefault()
    if (!searchQuery.trim()) return
    setSearchLoading(true)
    setError('')
    try {
      setSearchResults(await searchBooks(searchQuery.trim()))
    } catch (requestError) {
      setSearchResults([])
      setError(requestError instanceof Error ? requestError.message : 'Falha ao buscar livros.')
    } finally {
      setSearchLoading(false)
    }
  }

  async function save() {
    if (!book) {
      setError('Selecione um livro antes de salvar.')
      return
    }
    setSaving(true)
    setError('')
    try {
      await createReadingLog({ bookId: book.id, rating, readDate: date, reviewText: review || undefined, containsSpoilers: spoiler, liked })
      await updateShelf(book.id, shelfStatus)
      if (listId) await addBookToList(listId, book.id, 1)
      onSaved()
      close()
    } catch (requestError) {
      setError(requestError instanceof Error ? requestError.message : 'Não foi possível registrar a leitura.')
    } finally {
      setSaving(false)
    }
  }
  return <div className="modal-backdrop" onMouseDown={close}>
    <section className="log-modal" role="dialog" aria-modal="true" aria-labelledby="log-title" onMouseDown={e => e.stopPropagation()}>
      <div className="modal-head"><div><p className="eyebrow">Adicionar ao diário</p><h2 id="log-title">Registrar leitura</h2></div><button className="icon-button" onClick={close} aria-label="Fechar"><Icon name="close" /></button></div>
      {!book ? <div className="log-book-search">
        <p className="eyebrow">1. Selecione o livro</p>
        <form className="log-search-form" onSubmit={searchBook}><Icon name="search" /><input aria-label="Buscar livro para registrar" value={searchQuery} onChange={event => setSearchQuery(event.target.value)} placeholder="Título ou autor" /><button className="button primary compact" disabled={searchLoading || !searchQuery.trim()}>{searchLoading ? 'Buscando...' : 'Buscar'}</button></form>
        <div className="log-search-results">{searchResults.map(result => <button className="log-search-result" key={result.id} onClick={() => setBook(result)}>{result.coverUrl ? <img src={result.coverUrl} alt="" /> : <span className="cover-placeholder"><Icon name="book" /></span>}<span><strong>{result.title}</strong><small>{result.authors?.join(', ') || 'Autor desconhecido'}</small>{result.publishYear && <em>{result.publishYear}</em>}</span></button>)}</div>
      </div> : <div className="modal-book">{book.coverUrl ? <img src={book.coverUrl} alt="" /> : <Icon name="book" />}<div><strong>{book.title}</strong><span>{book.authors?.join(', ') || 'Autor desconhecido'}</span></div><button className="text-link" type="button" onClick={() => setBook(null)}>Trocar livro</button></div>}
      {book && <>
        <div className="log-section-title"><p className="eyebrow">2. Registro da leitura</p></div>
        <div className="form-grid">
          <label><span>Data de término</span><span className="input-shell"><Icon name="calendar" /><input type="date" value={date} onChange={event => setDate(event.target.value)} /></span></label>
          <label><span>Sua avaliação</span><Stars value={rating} interactive onChange={setRating} /></label>
        </div>
        <button className={`like-toggle ${liked ? 'liked' : ''}`} onClick={() => setLiked(!liked)}><Icon name="heart" /><span><strong>Você gostou?</strong><small>Marque como favorito</small></span><i /></button>
        <div className="log-section-title"><p className="eyebrow">3. Resenha opcional</p></div>
        <label className="review-field"><span>Sua resenha <small>Markdown suportado</small></span><textarea rows={5} maxLength={4000} value={review} onChange={event => setReview(event.target.value)} placeholder="O que ficou com você depois desta leitura?" /><small>{review.length} / 4.000</small></label>
        <label className="spoiler-toggle"><button type="button" className={spoiler ? 'on' : ''} onClick={() => setSpoiler(!spoiler)}><i /></button><span><strong>Contém spoilers</strong><small>A resenha ficará protegida no feed</small></span></label>
        <div className="log-section-title"><p className="eyebrow">4. Organização</p></div>
        <div className="form-grid organization-grid">
          <label><span>Estante</span><select value={shelfStatus} onChange={event => setShelfStatus(event.target.value)}><option value="tbr">Quero ler</option><option value="reading">Lendo</option><option value="read">Lido</option><option value="abandoned">Abandonado</option></select></label>
          <label><span>Adicionar a uma lista</span><select value={listId} onChange={event => setListId(event.target.value)}><option value="">Não adicionar</option>{lists.map(list => <option value={list.id} key={list.id}>{list.title}</option>)}</select></label>
        </div>
      </>}
      {error && <p role="alert" style={{ color: '#ff7a45' }}>{error}</p>}
      <div className="modal-actions"><button className="button ghost" onClick={close}>Cancelar</button><button className="button primary" disabled={saving} onClick={save}>{saving ? 'Salvando...' : 'Salvar no diário'}</button></div>
    </section>
  </div>
}

function SearchResults({ results, loading, error, close, onBookSelect }: { results: Book[]; loading: boolean; error: string; close: () => void; onBookSelect: (book: Book) => void }) {
  return <main className="page list-page">
    <div className="list-toolbar"><div><strong>Resultados da busca</strong><span>{loading ? 'Buscando...' : `${results.length} livro(s) encontrado(s)`}</span></div><button className="button secondary compact" onClick={close}>Voltar</button></div>
    {error && <p role="alert" style={{ color: '#ff7a45' }}>{error}</p>}
    {!loading && !error && results.length === 0 && <section className="diary-empty"><Icon name="book" size={32} /><h2>Nenhum livro encontrado</h2><p>Tente buscar por outro título ou autor.</p></section>}
    <section className="books-grid">{results.map(book => <button className="grid-book" key={book.id} onClick={() => onBookSelect(book)}><div className="book-cover-wrap"><div className="book-cover">{book.coverUrl ? <img src={book.coverUrl} alt={`Capa de ${book.title}`} /> : <span className="cover-placeholder"><Icon name="book" /></span>}</div></div><strong>{book.title}</strong><span>{book.authors?.join(', ') || 'Autor desconhecido'}</span></button>)}</section>
  </main>
}

export default function App() {
  const [screen, setScreen] = useState<Screen>('dashboard')
  const [logOpen, setLogOpen] = useState(false)
  const [user, setUser] = useState<User | null>(null)
  const [authLoading, setAuthLoading] = useState(hasToken())
  const [searchState, setSearchState] = useState<{ results: Book[]; loading: boolean; error: string } | null>(null)
  const [feedItems, setFeedItems] = useState<FeedItem[]>([])
  const [feedCursor, setFeedCursor] = useState<string>()
  const [feedMode, setFeedMode] = useState<'following' | 'global' | 'reviews'>('following')
  const [feedLoading, setFeedLoading] = useState(false)
  const [feedError, setFeedError] = useState('')
  const [popularBooks, setPopularBooks] = useState<Book[]>([])
  const [featuredBook, setFeaturedBook] = useState<Book | null>(null)
  const [selectedBook, setSelectedBook] = useState<BookDetails | null>(null)
  const [profile, setProfile] = useState<User | null>(null)
  const [profileStats, setProfileStats] = useState<Record<string, unknown>>({})
  const [logs, setLogs] = useState<ReadingLog[]>([])
  const [shelves, setShelves] = useState<ShelfEntry[]>([])
  const [lists, setLists] = useState<ReadingList[]>([])
  const [libraryError, setLibraryError] = useState('')

  useEffect(() => {
    if (!hasToken()) return
    getCurrentUser()
      .then(setUser)
      .catch(() => clearToken())
      .finally(() => setAuthLoading(false))
  }, [])

  async function loadFeed(mode = feedMode, append = false) {
    setFeedLoading(true)
    setFeedError('')
    try {
      if (mode === 'following') {
        const page = append && feedCursor ? await getMoreFeed(feedCursor) : await getFeed()
        setFeedItems(current => append ? [...current, ...page.items] : page.items)
        setFeedCursor(page.nextCursor || undefined)
      } else {
        const items = mode === 'global' ? await getGlobalFeed() : await getReviewFeed()
        setFeedItems(items)
        setFeedCursor(undefined)
      }
    } catch (error) {
      setFeedError(error instanceof Error ? error.message : 'Falha ao carregar o feed.')
    } finally {
      setFeedLoading(false)
    }
  }

  useEffect(() => {
    if (user) {
      loadFeed()
      getPopularBooks().then(setPopularBooks).catch(() => setPopularBooks([]))
      searchBooks('Qual e a tua obra').then(results => setFeaturedBook(results.find(book => book.title.toLowerCase().includes('qual e a tua obra')) ?? results[0] ?? null)).catch(() => setFeaturedBook(null))
      setLibraryError('')
      Promise.all([
        getProfile(user.username),
        getUserLogs(user.username),
        getAllUserShelves(user.username),
        getMyLists(),
        getUserStats(user.username),
      ]).then(([nextProfile, nextLogs, nextShelves, nextLists, nextStats]) => {
        setProfile(nextProfile)
        setProfileStats(nextStats)
        setLogs(nextLogs)
        setShelves(nextShelves)
        setLists(nextLists)
      }).catch(error => {
        setLibraryError(error instanceof Error ? error.message : 'Falha ao carregar sua biblioteca.')
      })
    }
  }, [user])

  function changeFeedMode(mode: 'following' | 'global' | 'reviews') {
    setFeedMode(mode)
    loadFeed(mode)
  }

  const navigate = (next: Screen) => { setScreen(next); window.scrollTo({ top: 0, behavior: 'smooth' }) }
  async function handleSearch(query: string) {
    setSearchState({ results: [], loading: true, error: '' })
    try {
      setSearchState({ results: await searchBooks(query), loading: false, error: '' })
    } catch (error) {
      setSearchState({ results: [], loading: false, error: error instanceof Error ? error.message : 'Falha ao buscar livros.' })
    }

  }

  async function handleBookSelect(book: Book) {
    try {
      setSelectedBook(await getBook(book.id))
    } catch (error) {
      setSelectedBook(book)
      setLibraryError(error instanceof Error ? error.message : 'Falha ao carregar os detalhes do livro.')
    }
    setSearchState(null)
    navigate('detail')
  }

  async function handleLogout() {
    try {
      await logout()
    } catch (error) {
      console.error('Não foi possível revogar a sessão no servidor.', error)
    } finally {
      clearToken()
      setUser(null)
      setSearchState(null)
    }
  }

  if (authLoading) return <main className="page" style={{ paddingTop: 96 }}><p>Carregando sua sessão...</p></main>
  if (!user) return <AuthScreen onAuthenticated={result => setUser(result.user)} />

  return <div className="app">
    <Header screen={screen} navigate={navigate} openLog={() => setLogOpen(true)} onSearch={handleSearch} onLogout={() => { void handleLogout() }} />
    {searchState ? <SearchResults {...searchState} close={() => setSearchState(null)} onBookSelect={handleBookSelect} /> : screen === 'dashboard' && <Dashboard user={user} navigate={navigate} popularBooks={popularBooks} featuredBook={featuredBook} onBookSelect={handleBookSelect} openLog={() => setLogOpen(true)} />}
    {libraryError && !searchState && <p className="page form-error" role="alert">{libraryError}</p>}
    {!searchState && screen === 'feed' && <Feed navigate={navigate} items={feedItems} loading={feedLoading} error={feedError} onRetry={() => loadFeed()} onRefresh={() => loadFeed()} onLoadMore={() => loadFeed('following', true)} hasMore={Boolean(feedCursor)} mode={feedMode} onModeChange={changeFeedMode} popularBooks={popularBooks} onBookSelect={handleBookSelect} />}
    {!searchState && screen === 'detail' && <Detail book={selectedBook} navigate={navigate} openLog={() => setLogOpen(true)} />}
    {!searchState && screen === 'profile' && <Profile user={profile ?? user} logs={logs} stats={profileStats} navigate={navigate} />}
    {!searchState && screen === 'list' && <ListPage navigate={navigate} lists={lists} onListsChange={async () => { try { setLists(await getMyLists()) } catch (error) { setLibraryError(error instanceof Error ? error.message : 'Falha ao recarregar listas.') } }} />}
    {!searchState && screen === 'diary' && <DiaryPage logs={logs} />}
    {!searchState && screen === 'shelves' && <ShelvesPage shelves={shelves} />}
    <MobileNav screen={screen} navigate={navigate} openLog={() => setLogOpen(true)} />
    {logOpen && <LogModal initialBook={(selectedBook as Book | null) || featuredBook || popularBooks[0] || null} lists={lists} close={() => setLogOpen(false)} onSaved={() => { loadFeed(); void getMyLists().then(setLists).catch(() => undefined); void getAllUserShelves(user.username).then(setShelves).catch(() => undefined) }} />}
  </div>
}
