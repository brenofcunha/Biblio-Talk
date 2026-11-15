import { useEffect, useState, type FormEvent, type ReactNode } from 'react'
import { clearToken, createReadingLog, getAllUserShelves, getCurrentUser, getFeed, getGlobalFeed, getMoreFeed, getMyLists, getPopularBooks, getReviewFeed, getProfile, getUserLogs, getUserStats, hasToken, login, logout, register, saveToken, searchBooks, type Book, type FeedItem, type ReadingList, type ReadingLog, type ShelfEntry, type User } from './api'

type Screen = 'dashboard' | 'feed' | 'detail' | 'profile' | 'list' | 'diary' | 'shelves'

const covers = [
  'https://images.unsplash.com/photo-1789938725934-0b97b0bfd365?auto=format&fit=crop&w=600&q=85',
  'https://images.unsplash.com/photo-1768229349400-64f2b0da1c3a?auto=format&fit=crop&w=600&q=85',
  'https://images.unsplash.com/photo-1773751392399-241afcb9fac2?auto=format&fit=crop&w=600&q=85',
  'https://images.unsplash.com/photo-1773761193222-4ad3b5ac23dc?auto=format&fit=crop&w=600&q=85',
  'https://images.unsplash.com/photo-1789938725904-553468f362ab?auto=format&fit=crop&w=600&q=85',
  'https://images.unsplash.com/photo-1768229349416-f551858d993a?auto=format&fit=crop&w=600&q=85',
  'https://images.unsplash.com/photo-1632516643720-e7f5d7d6ecc9?auto=format&fit=crop&w=600&q=85',
  'https://images.unsplash.com/photo-1782069243619-ed7c9dd3eb05?auto=format&fit=crop&w=600&q=85',
]

const books = [
  { title: 'O Avesso da Pele', author: 'Jeferson Tenório', cover: covers[0], rating: 4.7 },
  { title: 'Torto Arado', author: 'Itamar Vieira Junior', cover: covers[1], rating: 4.4 },
  { title: 'A Vegetariana', author: 'Han Kang', cover: covers[2], rating: 4.1 },
  { title: 'Tudo é Rio', author: 'Carla Madeira', cover: covers[3], rating: 4.3 },
  { title: 'A Cabeça do Santo', author: 'Socorro Acioli', cover: covers[4], rating: 4.2 },
  { title: 'Pessoas Normais', author: 'Sally Rooney', cover: covers[5], rating: 4.0 },
  { title: 'Véspera', author: 'Carla Madeira', cover: covers[6], rating: 4.1 },
  { title: 'A Palavra que Resta', author: 'Stênio Gardel', cover: covers[7], rating: 4.6 },
]

const users = [
  { name: 'Clara Nunes', handle: '@claralê', initials: 'CN', tone: 'coral' },
  { name: 'Miguel Costa', handle: '@miguelc', initials: 'MC', tone: 'lime' },
  { name: 'Bia Almeida', handle: '@bialmeida', initials: 'BA', tone: 'blue' },
]

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

function Avatar({ user = users[0], small = false, onClick }: { user?: typeof users[0]; small?: boolean; onClick?: () => void }) {
  return <button className={`avatar ${user.tone} ${small ? 'avatar-sm' : ''}`} onClick={onClick} aria-label={`Ver perfil de ${user.name}`}>{user.initials}</button>
}

function Brand({ onClick }: { onClick: () => void }) {
  return <button className="brand" onClick={onClick} aria-label="Ir para o início"><span className="brand-mark"><span /></span><span>Lume</span></button>
}

function Header({ screen, navigate, openLog, onSearch, onLogout }: { screen: Screen; navigate: (s: Screen) => void; openLog: () => void; onSearch: (query: string) => void; onLogout: () => void }) {
  const [query, setQuery] = useState('')
  return <header className="topbar">
    <div className="topbar-inner">
      <Brand onClick={() => navigate('feed')} />
      <form className="search" onSubmit={event => { event.preventDefault(); if (query.trim()) onSearch(query.trim()) }}><Icon name="search" /><input aria-label="Busca universal" placeholder="Busque livros, autores ou leitores" value={query} onChange={event => setQuery(event.target.value)} /></form>
      <nav className="main-nav" aria-label="Navegação principal">
        <button className={screen === 'dashboard' ? 'active' : ''} onClick={() => navigate('dashboard')}>Feed</button>
        <button className={screen === 'feed' ? 'active' : ''} onClick={() => navigate('feed')}>Comunidade</button>
        <button className={screen === 'shelves' ? 'active' : ''} onClick={() => navigate('shelves')}>Estantes</button>
        <button className={screen === 'list' ? 'active' : ''} onClick={() => navigate('list')}>Listas</button>
        <button className={screen === 'diary' ? 'active' : ''} onClick={() => navigate('diary')}>Diário</button>
        <button onClick={() => navigate('dashboard')}>Explorar</button>
      </nav>
      <button className="button primary compact desktop-log" onClick={openLog}><Icon name="plus" /> Registrar</button>
      <Avatar small onClick={() => navigate('profile')} />
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

function BookCover({ book = books[0], index, onClick, large = false }: { book?: typeof books[0]; index?: number; onClick: () => void; large?: boolean }) {
  return <div className={`book-cover-wrap ${large ? 'cover-large' : ''}`}>
    {index !== undefined && <span className="rank">{index + 1}</span>}
    <button className="book-cover" onClick={onClick} aria-label={`Ver ${book.title}`}>
      <img src={book.cover} alt={`Capa de ${book.title}`} />
      <span className="cover-shade" />
      <span className="cover-copy"><small>{book.author}</small><strong>{book.title}</strong></span>
      <span className="quick-actions"><span><Icon name="check" /> Lido</span><Icon name="heart" /></span>
    </button>
  </div>
}

function FeedCard({ item, navigate }: { item: FeedItem; navigate: (s: Screen) => void }) {
  const user = users[Array.from(item.user.username).reduce((sum, char) => sum + char.charCodeAt(0), 0) % users.length]
  const book = { title: item.book.title, author: item.book.authors?.join(', ') || 'Autor desconhecido', cover: item.book.coverUrl ?? covers[0], rating: item.rating }
  return <article className="feed-card">
    <div className="activity-head">
      <div className="user-line"><Avatar user={user} small onClick={() => navigate('profile')} /><div><button onClick={() => navigate('profile')}><strong>@{item.user.username}</strong></button><p>registrou uma leitura · {new Date(item.createdAt).toLocaleDateString('pt-BR')}</p></div></div>
      <button className="icon-button" aria-label="Mais opções"><Icon name="more" /></button>
    </div>
    <div className="activity-body">
      <BookCover book={book} onClick={() => navigate('detail')} />
      <div className="review">
        <button className="book-title" onClick={() => navigate('detail')}>{book.title}</button>
        <p className="author">{book.author} · 2024</p>
        <Stars value={item.rating} />
        <p className="review-text">{item.reviewText || 'Registrou esta leitura no diário.'}</p>
        {item.containsSpoilers && <span className="spoiler">Resenha contém spoilers</span>}
      </div>
    </div>
    <div className="activity-actions"><button><Icon name="heart" /> {item.likesCount}</button><button><Icon name="comment" /> {item.commentsCount ?? 0}</button></div>
  </article>
}

function Sidebar({ navigate }: { navigate: (s: Screen) => void }) {
  return <aside className="sidebar">
    <section className="side-section">
      <div className="section-heading"><h3>Populares esta semana</h3><button onClick={() => navigate('list')}>Ver todos</button></div>
      {books.slice(0, 3).map((book, i) => <button className="mini-book" key={book.title} onClick={() => navigate('detail')}><span className="mini-rank">0{i + 1}</span><img src={book.cover} alt="" /><span><strong>{book.title}</strong><small>{book.author}</small><em><span>★</span> {book.rating}</em></span></button>)}
    </section>
    <section className="side-section trending">
      <div className="section-heading"><h3>Listas em alta</h3></div>
      <button onClick={() => navigate('list')}><span className="list-collage">{books.slice(3, 6).map(b => <img key={b.title} src={b.cover} alt="" />)}</span><span><strong>Livros para sair de uma ressaca literária</strong><small>32 livros · por @bia.lê</small></span></button>
      <button onClick={() => navigate('list')}><span className="list-collage">{books.slice(5, 8).map(b => <img key={b.title} src={b.cover} alt="" />)}</span><span><strong>Brasil contemporâneo essencial</strong><small>18 livros · por @caiom</small></span></button>
    </section>
    <section className="network-card"><div><span className="pulse-dot" /><small>Sua rede esta semana</small></div><strong>74 livros</strong><p>lidos por 38 pessoas que você segue</p></section>
  </aside>
}

function Dashboard({ user, navigate, popularBooks, openLog }: { user: User; navigate: (s: Screen) => void; popularBooks: Book[]; openLog: () => void }) {
  const catalog = popularBooks.length ? popularBooks : books.map(book => ({ id: book.title, title: book.title, authors: [book.author], coverUrl: book.cover, publishYear: 2020 }))
  return <main className="page dashboard">
    <section className="dashboard-main">
      <div className="welcome-row"><div><p className="eyebrow">Sua vida em livros</p><h1>Boa noite, {user.username}.</h1><p>Continue de onde parou ou encontre sua próxima leitura.</p></div><div className="goal-card"><span>0<small>/24</small></span><div><b>Meta de leitura</b><i><em style={{ width: '0%' }} /></i></div></div></div>
      <section className="content-block"><div className="section-head"><div><p className="eyebrow">Descobertas</p><h2>Em alta na comunidade</h2></div><button className="text-link" onClick={() => navigate('feed')}>Ver comunidade</button></div>
        <div className="book-grid featured">{catalog.slice(0, 4).map((book, index) => <article className="book-card" key={book.id}><button className="book-cover" onClick={() => navigate('detail')}><span className="rank">{index + 1}</span><img src={book.coverUrl ?? covers[index]} alt={`Capa de ${book.title}`} /><span className="cover-overlay"><b>Ver livro</b><i>♡</i></span></button><button className="book-card-link" onClick={() => navigate('detail')}><strong>{book.title}</strong><span>{book.authors?.join(', ')}</span></button><small><b>★</b> {book.publishYear ?? '—'}</small></article>)}</div>
      </section>
      <section className="content-block" id="catalogo"><div className="section-head"><div><p className="eyebrow">Catálogo</p><h2>Explore todos os livros</h2></div><span>{catalog.length} títulos</span></div>
        <div className="book-grid">{catalog.map(book => <article className="book-card" key={book.id}><button className="book-cover" onClick={() => navigate('detail')}><img src={book.coverUrl ?? covers[0]} alt={`Capa de ${book.title}`} /></button><button className="book-card-link" onClick={() => navigate('detail')}><strong>{book.title}</strong><span>{book.authors?.join(', ')}</span></button></article>)}</div>
      </section>
    </section>
    <aside className="dashboard-aside">
      <div className="aside-card profile-mini"><span className="profile-avatar">{user.username.slice(0, 2).toUpperCase()}</span><div><strong>{user.username}</strong><small>@{user.username}</small></div><button onClick={() => navigate('profile')}>Ver perfil</button></div>
      <div className="aside-card"><div className="aside-title"><h3>Seu diário</h3><button onClick={() => navigate('detail')}>Abrir</button></div><div className="mini-empty"><Icon name="book" /><p>Seu diário começa com uma leitura.</p><button onClick={openLog}>Registrar leitura</button></div></div>
      <div className="aside-card"><div className="aside-title"><h3>Suas estantes</h3><button onClick={() => navigate('list')}>Gerenciar</button></div>{['Quero ler', 'Lendo', 'Lidos'].map((shelf, index) => <button className="shelf-line" key={shelf}><span><Icon name={index === 2 ? 'check' : 'bookmark'} /></span><b>{shelf}</b><small>0</small></button>)}</div>
    </aside>
  </main>
}

function Feed({ navigate, items, loading, error, onRetry, onRefresh, onLoadMore, hasMore, mode, onModeChange }: { navigate: (s: Screen) => void; items: FeedItem[]; loading: boolean; error: string; onRetry: () => void; onRefresh: () => void; onLoadMore: () => void; hasMore: boolean; mode: 'following' | 'global' | 'reviews'; onModeChange: (mode: 'following' | 'global' | 'reviews') => void }) {
  return <main className="page feed-layout">
    <section className="feed-column">
      <div className="feed-intro"><div><p className="eyebrow">Seu diário social</p><h1>Veja o que sua rede anda lendo.</h1><p>Atualizado diretamente pela API.</p></div><button className="button secondary compact" onClick={onRefresh} disabled={loading}><Icon name="refresh" /> {loading ? 'Atualizando...' : 'Atualizar'}</button></div>
      <div className="filter-row"><button className={mode === 'following' ? 'active' : ''} onClick={() => onModeChange('following')}>Sua rede</button><button className={mode === 'global' ? 'active' : ''} onClick={() => onModeChange('global')}>Atividade global</button><button className={mode === 'reviews' ? 'active' : ''} onClick={() => onModeChange('reviews')}>Resenhas</button></div>
      {loading && !items.length && <section className="feed-skeletons" aria-label="Carregando feed">{[1, 2, 3].map(index => <div className="feed-skeleton" key={index}><span className="skeleton-avatar" /><div><i className="skeleton-line short" /><i className="skeleton-line" /><i className="skeleton-line medium" /><div className="skeleton-book" /></div></div>)}</section>}
      {loading && items.length > 0 && <p className="feed-refreshing" role="status">Atualizando o feed...</p>}
      {error && <section className="diary-empty"><p role="alert">{error}</p><button className="button primary" onClick={onRetry}>Tentar novamente</button></section>}
      {!loading && !error && items.length === 0 && <section className="diary-empty"><h2>O feed está vazio</h2><p>Siga leitores ou registre uma leitura para começar.</p></section>}
      {!loading && !error && items.map(item => <FeedCard key={item.id} item={item} navigate={navigate} />)}
      {!loading && !error && hasMore && <button className="button secondary" onClick={onLoadMore}>Carregar mais</button>}
    </section>
    <Sidebar navigate={navigate} />
  </main>
}

function Detail({ navigate, openLog }: { navigate: (s: Screen) => void; openLog: () => void }) {
  const book = books[0]
  return <main>
    <section className="book-hero">
      <div className="hero-backdrop" style={{ backgroundImage: `url(${book.cover})` }} />
      <div className="page hero-content">
        <BookCover book={book} large onClick={() => {}} />
        <div className="book-meta">
          <p className="eyebrow">Ficção brasileira · 2020</p>
          <h1>{book.title}</h1>
          <button className="author-link" onClick={() => navigate('profile')}>por Jeferson Tenório</button>
          <p className="book-description">Um romance sobre identidade, violência e os laços entre pai e filho. Uma narrativa contundente que atravessa memória, afeto e resistência.</p>
          <div className="metadata"><span>Companhia das Letras</span><span>192 páginas</span><span>Português</span></div>
          <div className="book-actions">
            <button className="button primary" onClick={openLog}><Icon name="plus" /> Registrar leitura</button>
            <button className="button secondary"><Icon name="check" /> Marcar como lido</button>
            <button className="button icon-only" aria-label="Adicionar à lista de desejos"><Icon name="bookmark" /></button>
          </div>
          <div className="quick-rate"><span>Sua nota</span><Stars value={0} interactive /></div>
        </div>
      </div>
    </section>
    <div className="page detail-grid">
      <section>
        <div className="stats-card">
          <div className="score"><strong>4.7</strong><Stars value={5} /><span>8.492 avaliações</span></div>
          <div className="histogram">
            {[18, 28, 52, 84, 100].map((h, i) => <div key={h}><span>{i + 1}</span><i><b style={{ height: `${h}%` }} /></i></div>)}
          </div>
        </div>
        <div className="content-section"><div className="section-title"><h2>Sobre o livro</h2></div><p>Depois que o pai é assassinado numa abordagem policial, Pedro retorna à casa onde cresceu, em Porto Alegre. Em busca de respostas, ele refaz os caminhos do pai e recompõe uma história marcada por racismo, afetos e silêncios.</p><button className="text-link">Continuar lendo</button></div>
        <div className="content-section">
          <div className="section-title"><h2>Resenhas da comunidade</h2><button>Mais relevantes</button></div>
          {[users[1], users[2]].map((user, i) => <article className="community-review" key={user.name}><Avatar user={user} small onClick={() => navigate('profile')} /><div><div><button onClick={() => navigate('profile')}><strong>{user.name}</strong></button><Stars value={i ? 4 : 5} /></div><p>{i ? 'Uma prosa precisa, sem excessos. Terminei querendo conversar com alguém sobre tudo que ficou nas entrelinhas.' : 'Este livro me encontrou num momento muito específico. A voz de Pedro fica com a gente muito depois da última página.'}</p><span><Icon name="heart" /> {i ? 19 : 46} curtidas</span></div></article>)}
        </div>
      </section>
      <aside className="detail-aside"><h3>Leitores que você segue</h3><div className="avatar-stack">{users.map(u => <Avatar key={u.name} user={u} small onClick={() => navigate('profile')} />)}<span>+12</span></div><p>15 amigos já leram este livro. A média deles é <strong>4.8</strong>.</p><hr /><h3>Informações</h3><dl><div><dt>Publicado em</dt><dd>2020</dd></div><div><dt>ISBN</dt><dd>978-85-359-3392-9</dd></div><div><dt>Gêneros</dt><dd>Literatura brasileira, Romance</dd></div></dl></aside>
    </div>
  </main>
}

function Profile({ user, logs, stats, navigate }: { user: User; logs: ReadingLog[]; stats?: Record<string, unknown>; navigate: (s: Screen) => void }) {
  const [tab, setTab] = useState('Estatísticas')
  return <main>
    <section className="profile-hero">
      <div className="page profile-head">
        <Avatar user={users[0]} />
        <div className="profile-copy"><div><h1>{user.username}</h1><button className="button secondary compact">Editar perfil</button></div><p className="handle">@{user.username}</p><p>{user.bio || 'Registre suas leituras e compartilhe suas descobertas com a comunidade.'}</p><div className="profile-numbers"><span><strong>{Number(stats?.booksRead ?? logs.length)}</strong> lidos no ano</span><span><strong>{logs.length}</strong> entradas no diário</span><span><strong>{user.top4BookIds?.length ?? 0}</strong> favoritos</span></div></div>
      </div>
    </section>
    <div className="page profile-page">
      <section className="top-four"><div className="section-title"><div><p className="eyebrow">Estante pessoal</p><h2>Quatro para uma ilha deserta</h2></div><button onClick={() => navigate('list')}>Ver estante <Icon name="arrow" /></button></div><div className="top-four-grid">{books.slice(0, 4).map(b => <BookCover key={b.title} book={b} onClick={() => navigate('detail')} />)}</div></section>
      <nav className="profile-tabs">{['Diário', 'Resenhas', 'Listas', 'Estatísticas'].map(t => <button className={tab === t ? 'active' : ''} key={t} onClick={() => setTab(t)}>{t}</button>)}</nav>
      {tab === 'Estatísticas' ? <Stats stats={stats} /> : <Diary navigate={navigate} tab={tab} logs={logs} />}
    </div>
  </main>
}

function Stats({ stats }: { stats?: Record<string, unknown> }) {
  const booksRead = Number(stats?.booksRead ?? 0)
  const pagesRead = Number(stats?.pagesRead ?? 0)
  return <section className="stats-dashboard">
    <div className="stats-header"><div><p className="eyebrow">Seu ano em livros</p><h2>2025, até agora</h2></div><button className="button secondary compact">2025 <span>⌄</span></button></div>
    <div className="stat-kpis"><div><Icon name="book" /><span><strong>{booksRead}</strong><small>livros lidos</small></span><em>dados da API</em></div><div><Icon name="pages" /><span><strong>{pagesRead.toLocaleString('pt-BR')}</strong><small>páginas lidas</small></span><em>dados da API</em></div><div><Icon name="calendar" /><span><strong>{Number(stats?.currentStreakDays ?? 0)}</strong><small>dias em sequência</small></span><em>dados da API</em></div></div>
    <div className="chart-grid">
      <div className="chart-card"><div className="section-title"><h3>Gêneros favoritos</h3><span>Por livros lidos</span></div><div className="donut-row"><div className="donut"><span><strong>38</strong><small>livros</small></span></div><ul><li><i className="dot green" />Ficção literária <strong>42%</strong></li><li><i className="dot orange" />Romance <strong>27%</strong></li><li><i className="dot blue" />Contos <strong>18%</strong></li><li><i className="dot gray" />Outros <strong>13%</strong></li></ul></div></div>
      <div className="chart-card"><div className="section-title"><h3>Ritmo de leitura</h3><span>Livros por mês</span></div><div className="bar-chart">{[4, 7, 5, 9, 6, 7].map((n, i) => <div key={i}><i><b style={{ height: `${n * 10}%` }} /></i><span>{['Jan', 'Fev', 'Mar', 'Abr', 'Mai', 'Jun'][i]}</span></div>)}</div></div>
    </div>
    <div className="author-highlight"><span className="author-monogram">CM</span><div><p className="eyebrow">Autora mais lida</p><h3>Carla Madeira</h3><p>3 livros · 782 páginas</p></div><div className="author-books">{books.slice(3, 6).map(b => <img key={b.title} src={b.cover} alt="" />)}</div></div>
  </section>
}

function Diary({ navigate, tab, logs }: { navigate: (s: Screen) => void; tab: string; logs: ReadingLog[] }) {
  if (!logs.length) return <section className="diary-empty"><Icon name="book" size={32} /><h2>{tab} vazio</h2><p>Registre uma leitura para começar seu diário.</p><button className="button primary" onClick={() => navigate('detail')}>Explorar livro em destaque</button></section>
  return <section className="diary-list">{logs.slice(0, 10).map(log => <article className="diary-entry" key={log.id}><time>{log.readDate ? new Date(log.readDate).toLocaleDateString('pt-BR') : 'Sem data'}</time><div className="diary-cover"><img src={log.book.coverUrl ?? covers[0]} alt="" /></div><div><h2>{log.book.title}</h2><p>{log.book.authors?.join(', ')}</p><span className="rating">{'★'.repeat(Math.round(log.rating))} <b>{log.rating.toFixed(1)}</b></span>{log.reviewText && <blockquote>{log.reviewText}</blockquote>}</div></article>)}</section>
}

function ListPage({ navigate, lists }: { navigate: (s: Screen) => void; lists: ReadingList[] }) {
  const visibleLists = lists.length ? lists : []
  return <main className="page list-page">
    <section className="list-header">
      <div className="list-art">{books.slice(0, 4).map(b => <img key={b.title} src={b.cover} alt="" />)}</div>
      <div><p className="eyebrow">Suas listas · {visibleLists.length}</p><h1>{visibleLists[0]?.title ?? 'Suas listas de leitura'}</h1><p>{visibleLists[0]?.description ?? 'Crie listas para organizar suas próximas leituras e recomendações.'}</p><div className="list-owner"><Avatar user={users[2]} small onClick={() => navigate('profile')} /><span>Listas carregadas pela API Go</span></div></div>
    </section>
    <div className="list-toolbar"><div><strong>{visibleLists.length} listas</strong><span>Conteúdo da sua conta</span></div></div>
    {!visibleLists.length && <section className="diary-empty"><Icon name="list" size={32} /><h2>Nenhuma lista criada</h2><p>As listas criadas pela API aparecerão aqui.</p></section>}
    <section className="books-grid">{visibleLists.map(list => <article className="grid-book" key={list.id}><div className="book-cover-wrap"><div className="book-cover"><img src={covers[0]} alt="" /></div></div><strong>{list.title}</strong><span>{list.description || 'Lista de leitura'}</span><small>{list.isPrivate ? 'Privada' : 'Pública'}</small></article>)}</section>
  </main>
}

function DiaryPage({ logs }: { logs: ReadingLog[] }) {
  return <main className="page inner-page"><header className="page-title"><div><p className="eyebrow">Histórico de leitura</p><h1>Seu diário</h1><p>Leituras carregadas da API Go.</p></div><div className="diary-kpis"><span><strong>{logs.length}</strong> entradas</span></div></header><Diary navigate={() => undefined} tab="Diário" logs={logs} /></main>
}

function ShelvesPage({ shelves }: { shelves: ShelfEntry[] }) {
  return <main className="page inner-page"><header className="page-title"><div><p className="eyebrow">Organização</p><h1>Suas estantes</h1><p>Livros organizados por status na API Go.</p></div></header>{['want_to_read', 'reading', 'read', 'abandoned'].map(status => { const entries = shelves.filter(shelf => shelf.status === status); return <section className="shelf-section" key={status}><p className="eyebrow">{status}</p>{entries.length ? <p>{entries.length} livro(s) nesta estante.</p> : <div className="empty-state"><h3>Estante vazia</h3></div>}</section> })}</main>
}

function LogModal({ book, close, onSaved }: { book: Book; close: () => void; onSaved: () => void }) {
  const [spoiler, setSpoiler] = useState(false)
  const [liked, setLiked] = useState(false)
  const [rating, setRating] = useState(4)
  const [date, setDate] = useState(new Date().toISOString().slice(0, 10))
  const [review, setReview] = useState('')
  const [saving, setSaving] = useState(false)
  const [error, setError] = useState('')
  async function save() {
    setSaving(true)
    setError('')
    try {
      await createReadingLog({ bookId: book.id, rating, readDate: date, reviewText: review || undefined, containsSpoilers: spoiler, liked })
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
      <div className="modal-book"><img src={book.coverUrl ?? covers[0]} alt="" /><div><strong>{book.title}</strong><span>{book.authors?.join(', ') || 'Autor desconhecido'}</span></div></div>
      <div className="form-grid">
        <label><span>Data de término</span><span className="input-shell"><Icon name="calendar" /><input type="date" value={date} onChange={event => setDate(event.target.value)} /></span></label>
        <label><span>Sua avaliação</span><Stars value={rating} interactive onChange={setRating} /></label>
      </div>
      <button className={`like-toggle ${liked ? 'liked' : ''}`} onClick={() => setLiked(!liked)}><Icon name="heart" /><span><strong>Você gostou?</strong><small>Marque como um dos favoritos</small></span><i /></button>
      <label className="review-field"><span>Sua resenha <small>Markdown suportado</small></span><textarea rows={5} maxLength={4000} value={review} onChange={event => setReview(event.target.value)} placeholder="O que ficou com você depois desta leitura?" /><small>{review.length} / 4.000</small></label>
      <div className="tag-field"><span>Tags pessoais</span><div><span>ficção brasileira <button aria-label="Remover tag"><Icon name="close" size={12} /></button></span><input placeholder="+ Adicionar tag" /></div></div>
      <label className="spoiler-toggle"><button type="button" className={spoiler ? 'on' : ''} onClick={() => setSpoiler(!spoiler)}><i /></button><span><strong>Contém spoilers</strong><small>A resenha ficará protegida no feed</small></span></label>
      {error && <p role="alert" style={{ color: '#ff7a45' }}>{error}</p>}
      <div className="modal-actions"><button className="button ghost" onClick={close}>Cancelar</button><button className="button primary" disabled={saving} onClick={save}>{saving ? 'Salvando...' : 'Salvar no diário'}</button></div>
    </section>
  </div>
}

function SearchResults({ results, loading, error, close }: { results: Book[]; loading: boolean; error: string; close: () => void }) {
  return <main className="page list-page">
    <div className="list-toolbar"><div><strong>Resultados da busca</strong><span>{loading ? 'Buscando...' : `${results.length} livro(s) encontrado(s)`}</span></div><button className="button secondary compact" onClick={close}>Voltar</button></div>
    {error && <p role="alert" style={{ color: '#ff7a45' }}>{error}</p>}
    {!loading && !error && results.length === 0 && <section className="diary-empty"><Icon name="book" size={32} /><h2>Nenhum livro encontrado</h2><p>Tente buscar por outro título ou autor.</p></section>}
    <section className="books-grid">{results.map(book => <div className="grid-book" key={book.id}><div className="book-cover-wrap"><div className="book-cover"><img src={book.coverUrl ?? covers[0]} alt={`Capa de ${book.title}`} /></div></div><strong>{book.title}</strong><span>{book.authors?.join(', ') || 'Autor desconhecido'}</span></div>)}</section>
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
    {searchState ? <SearchResults {...searchState} close={() => setSearchState(null)} /> : screen === 'dashboard' && <Dashboard user={user} navigate={navigate} popularBooks={popularBooks} openLog={() => setLogOpen(true)} />}
    {libraryError && !searchState && <p className="page form-error" role="alert">{libraryError}</p>}
    {!searchState && screen === 'feed' && <Feed navigate={navigate} items={feedItems} loading={feedLoading} error={feedError} onRetry={() => loadFeed()} onRefresh={() => loadFeed()} onLoadMore={() => loadFeed('following', true)} hasMore={Boolean(feedCursor)} mode={feedMode} onModeChange={changeFeedMode} />}
    {!searchState && screen === 'detail' && <Detail navigate={navigate} openLog={() => setLogOpen(true)} />}
    {!searchState && screen === 'profile' && <Profile user={profile ?? user} logs={logs} stats={profileStats} navigate={navigate} />}
    {!searchState && screen === 'list' && <ListPage navigate={navigate} lists={lists} />}
    {!searchState && screen === 'diary' && <DiaryPage logs={logs} />}
    {!searchState && screen === 'shelves' && <ShelvesPage shelves={shelves} />}
    <MobileNav screen={screen} navigate={navigate} openLog={() => setLogOpen(true)} />
    {logOpen && <LogModal book={popularBooks[0] ?? { id: '', title: books[0].title, authors: [books[0].author], coverUrl: books[0].cover }} close={() => setLogOpen(false)} onSaved={() => loadFeed()} />}
  </div>
}
