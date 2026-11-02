import { useState, type ReactNode } from 'react'

type Screen = 'feed' | 'detail' | 'profile' | 'list'

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
  }
  return <svg aria-hidden="true" width={size} height={size} viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.8" strokeLinecap="round" strokeLinejoin="round">{paths[name]}</svg>
}

function Avatar({ user = users[0], small = false, onClick }: { user?: typeof users[0]; small?: boolean; onClick?: () => void }) {
  return <button className={`avatar ${user.tone} ${small ? 'avatar-sm' : ''}`} onClick={onClick} aria-label={`Ver perfil de ${user.name}`}>{user.initials}</button>
}

function Brand({ onClick }: { onClick: () => void }) {
  return <button className="brand" onClick={onClick} aria-label="Ir para o feed"><span className="brand-mark"><span /></span><span>Lume</span></button>
}

function Header({ screen, navigate, openLog }: { screen: Screen; navigate: (s: Screen) => void; openLog: () => void }) {
  return <header className="topbar">
    <div className="topbar-inner">
      <Brand onClick={() => navigate('feed')} />
      <div className="search"><Icon name="search" /><input aria-label="Busca universal" placeholder="Busque livros, autores ou leitores" /></div>
      <nav className="main-nav" aria-label="Navegação principal">
        <button className={screen === 'feed' ? 'active' : ''} onClick={() => navigate('feed')}>Feed</button>
        <button className={screen === 'list' ? 'active' : ''} onClick={() => navigate('list')}>Listas</button>
        <button onClick={() => navigate('feed')}>Explorar</button>
      </nav>
      <button className="button primary compact desktop-log" onClick={openLog}><Icon name="plus" /> Registrar</button>
      <Avatar small onClick={() => navigate('profile')} />
    </div>
  </header>
}

function MobileNav({ screen, navigate, openLog }: { screen: Screen; navigate: (s: Screen) => void; openLog: () => void }) {
  return <nav className="mobile-nav" aria-label="Navegação mobile">
    <button className={screen === 'feed' ? 'active' : ''} onClick={() => navigate('feed')}><Icon name="home" /><span>Feed</span></button>
    <button onClick={() => navigate('feed')}><Icon name="compass" /><span>Explorar</span></button>
    <button className="mobile-add" onClick={openLog}><Icon name="plus" size={22} /></button>
    <button className={screen === 'list' ? 'active' : ''} onClick={() => navigate('list')}><Icon name="list" /><span>Listas</span></button>
    <button className={screen === 'profile' ? 'active' : ''} onClick={() => navigate('profile')}><Icon name="user" /><span>Perfil</span></button>
  </nav>
}

function Stars({ value = 4.5, interactive = false }: { value?: number; interactive?: boolean }) {
  const [rating, setRating] = useState(value)
  return <div className={`stars ${interactive ? 'stars-interactive' : ''}`} aria-label={`${rating} de 5 estrelas`}>
    {[1, 2, 3, 4, 5].map(n => <button key={n} type="button" onClick={() => interactive && setRating(n)} tabIndex={interactive ? 0 : -1} aria-label={`${n} estrelas`}><span className={n <= rating ? 'filled' : ''}>★</span></button>)}
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

function FeedCard({ index, navigate }: { index: number; navigate: (s: Screen) => void }) {
  const user = users[index % users.length]
  const book = books[index]
  return <article className="feed-card">
    <div className="activity-head">
      <div className="user-line"><Avatar user={user} small onClick={() => navigate('profile')} /><div><button onClick={() => navigate('profile')}><strong>{user.name}</strong></button><p>terminou uma leitura · {index === 0 ? 'há 2 h' : index === 1 ? 'ontem' : 'terça'}</p></div></div>
      <button className="icon-button" aria-label="Mais opções"><Icon name="more" /></button>
    </div>
    <div className="activity-body">
      <BookCover book={book} onClick={() => navigate('detail')} />
      <div className="review">
        <button className="book-title" onClick={() => navigate('detail')}>{book.title}</button>
        <p className="author">{book.author} · 2024</p>
        <Stars value={index === 1 ? 4 : 5} />
        <p className="review-text">{index === 0 ? 'Um livro que não pede licença. A narrativa fala sobre identidade, luto e as marcas que herdamos — com uma delicadeza que corta.' : index === 1 ? 'A terra aqui não é cenário, é personagem. Uma história sobre voz, ancestralidade e tudo o que permanece.' : 'Li devagar para não chegar ao fim. Estranho, inquietante e impossível de esquecer.'}</p>
        {index === 2 && <span className="spoiler">Resenha contém spoilers</span>}
      </div>
    </div>
    <div className="activity-actions"><button><Icon name="heart" /> {28 + index * 13}</button><button><Icon name="comment" /> {4 + index * 2}</button><span>Leitura em {index === 0 ? '6 dias' : '11 dias'}</span></div>
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

function Feed({ navigate }: { navigate: (s: Screen) => void }) {
  return <main className="page feed-layout">
    <section className="feed-column">
      <div className="feed-intro"><div><p className="eyebrow">Quinta-feira, 12 de junho</p><h1>Boa noite, Clara.</h1><p>Veja o que sua rede anda lendo.</p></div><div className="reading-goal"><span><strong>24</strong>/36</span><div><p>Meta de leitura</p><div className="progress"><span /></div></div></div></div>
      <div className="filter-row"><button className="active">Sua rede</button><button>Atividade global</button><button>Resenhas</button></div>
      {[0, 1, 2].map(i => <FeedCard key={i} index={i} navigate={navigate} />)}
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

function Profile({ navigate }: { navigate: (s: Screen) => void }) {
  const [tab, setTab] = useState('Estatísticas')
  return <main>
    <section className="profile-hero">
      <div className="page profile-head">
        <Avatar user={users[0]} />
        <div className="profile-copy"><div><h1>Clara Nunes</h1><button className="button secondary compact">Editar perfil</button></div><p className="handle">@claralê · São Paulo, SP</p><p>Leio para encontrar nomes para coisas que já senti. Ficção brasileira, mulheres estranhas e finais que doem.</p><div className="profile-numbers"><span><strong>38</strong> lidos em 2025</span><span><strong>286</strong> livros</span><span><strong>12</strong> listas</span><span><strong>1.8k</strong> seguidores</span></div></div>
      </div>
    </section>
    <div className="page profile-page">
      <section className="top-four"><div className="section-title"><div><p className="eyebrow">Estante pessoal</p><h2>Quatro para uma ilha deserta</h2></div><button onClick={() => navigate('list')}>Ver estante <Icon name="arrow" /></button></div><div className="top-four-grid">{books.slice(0, 4).map(b => <BookCover key={b.title} book={b} onClick={() => navigate('detail')} />)}</div></section>
      <nav className="profile-tabs">{['Diário', 'Resenhas', 'Listas', 'Estatísticas'].map(t => <button className={tab === t ? 'active' : ''} key={t} onClick={() => setTab(t)}>{t}</button>)}</nav>
      {tab === 'Estatísticas' ? <Stats /> : <Diary navigate={navigate} tab={tab} />}
    </div>
  </main>
}

function Stats() {
  return <section className="stats-dashboard">
    <div className="stats-header"><div><p className="eyebrow">Seu ano em livros</p><h2>2025, até agora</h2></div><button className="button secondary compact">2025 <span>⌄</span></button></div>
    <div className="stat-kpis"><div><Icon name="book" /><span><strong>38</strong><small>livros lidos</small></span><em>+12% vs. 2024</em></div><div><Icon name="pages" /><span><strong>9.482</strong><small>páginas lidas</small></span><em>26 por dia</em></div><div><Icon name="calendar" /><span><strong>21</strong><small>dias em sequência</small></span><em>recorde: 34</em></div></div>
    <div className="chart-grid">
      <div className="chart-card"><div className="section-title"><h3>Gêneros favoritos</h3><span>Por livros lidos</span></div><div className="donut-row"><div className="donut"><span><strong>38</strong><small>livros</small></span></div><ul><li><i className="dot green" />Ficção literária <strong>42%</strong></li><li><i className="dot orange" />Romance <strong>27%</strong></li><li><i className="dot blue" />Contos <strong>18%</strong></li><li><i className="dot gray" />Outros <strong>13%</strong></li></ul></div></div>
      <div className="chart-card"><div className="section-title"><h3>Ritmo de leitura</h3><span>Livros por mês</span></div><div className="bar-chart">{[4, 7, 5, 9, 6, 7].map((n, i) => <div key={i}><i><b style={{ height: `${n * 10}%` }} /></i><span>{['Jan', 'Fev', 'Mar', 'Abr', 'Mai', 'Jun'][i]}</span></div>)}</div></div>
    </div>
    <div className="author-highlight"><span className="author-monogram">CM</span><div><p className="eyebrow">Autora mais lida</p><h3>Carla Madeira</h3><p>3 livros · 782 páginas</p></div><div className="author-books">{books.slice(3, 6).map(b => <img key={b.title} src={b.cover} alt="" />)}</div></div>
  </section>
}

function Diary({ navigate, tab }: { navigate: (s: Screen) => void; tab: string }) {
  return <section className="diary-empty"><Icon name={tab === 'Listas' ? 'list' : 'book'} size={32} /><h2>{tab} de Clara</h2><p>Esta visualização faz parte do protótipo. Explore as estatísticas ou abra um livro.</p><button className="button primary" onClick={() => navigate('detail')}>Explorar livro em destaque</button></section>
}

function ListPage({ navigate }: { navigate: (s: Screen) => void }) {
  return <main className="page list-page">
    <section className="list-header">
      <div className="list-art">{books.slice(0, 4).map(b => <img key={b.title} src={b.cover} alt="" />)}</div>
      <div><p className="eyebrow">Lista editorial · 24 livros</p><h1>Livros para sair de uma ressaca literária</h1><p>Histórias que agarram desde a primeira página — para quando você quer voltar a sentir tudo que uma boa leitura pode provocar.</p><div className="list-owner"><Avatar user={users[2]} small onClick={() => navigate('profile')} /><span>Criada por <button onClick={() => navigate('profile')}>Bia Almeida</button> · atualizada há 3 dias</span></div></div>
    </section>
    <div className="list-toolbar"><div><strong>24 livros</strong><span>Ordenados pela curadoria</span></div><button className="button secondary compact"><Icon name="list" /> Ordenar</button></div>
    <section className="books-grid">{books.map((b, i) => <div className="grid-book" key={b.title}><BookCover book={b} index={i} onClick={() => navigate('detail')} /><button onClick={() => navigate('detail')}><strong>{b.title}</strong><span>{b.author}</span></button><small><span>★</span> {b.rating} · {12 + i * 4}k leituras</small></div>)}</section>
  </main>
}

function LogModal({ close }: { close: () => void }) {
  const [spoiler, setSpoiler] = useState(false)
  const [liked, setLiked] = useState(false)
  const [saved, setSaved] = useState(false)
  if (saved) return <div className="modal-backdrop" onMouseDown={close}><div className="success-modal" onMouseDown={e => e.stopPropagation()}><span><Icon name="check" size={28} /></span><h2>Leitura registrada</h2><p>O Avesso da Pele agora faz parte do seu diário.</p><button className="button primary" onClick={close}>Voltar ao livro</button></div></div>
  return <div className="modal-backdrop" onMouseDown={close}>
    <section className="log-modal" role="dialog" aria-modal="true" aria-labelledby="log-title" onMouseDown={e => e.stopPropagation()}>
      <div className="modal-head"><div><p className="eyebrow">Adicionar ao diário</p><h2 id="log-title">Registrar leitura</h2></div><button className="icon-button" onClick={close} aria-label="Fechar"><Icon name="close" /></button></div>
      <div className="modal-book"><img src={books[0].cover} alt="" /><div><strong>{books[0].title}</strong><span>Jeferson Tenório</span></div></div>
      <div className="form-grid">
        <label><span>Data de término</span><span className="input-shell"><Icon name="calendar" /><input type="date" defaultValue="2025-06-12" /></span></label>
        <label><span>Sua avaliação</span><Stars value={4} interactive /></label>
      </div>
      <button className={`like-toggle ${liked ? 'liked' : ''}`} onClick={() => setLiked(!liked)}><Icon name="heart" /><span><strong>Você gostou?</strong><small>Marque como um dos favoritos</small></span><i /></button>
      <label className="review-field"><span>Sua resenha <small>Markdown suportado</small></span><textarea rows={5} placeholder="O que ficou com você depois desta leitura?" /><small>0 / 4.000</small></label>
      <div className="tag-field"><span>Tags pessoais</span><div><span>ficção brasileira <button aria-label="Remover tag"><Icon name="close" size={12} /></button></span><input placeholder="+ Adicionar tag" /></div></div>
      <label className="spoiler-toggle"><button type="button" className={spoiler ? 'on' : ''} onClick={() => setSpoiler(!spoiler)}><i /></button><span><strong>Contém spoilers</strong><small>A resenha ficará protegida no feed</small></span></label>
      <div className="modal-actions"><button className="button ghost" onClick={close}>Cancelar</button><button className="button primary" onClick={() => setSaved(true)}>Salvar no diário</button></div>
    </section>
  </div>
}

export default function App() {
  const [screen, setScreen] = useState<Screen>('feed')
  const [logOpen, setLogOpen] = useState(false)
  const navigate = (next: Screen) => { setScreen(next); window.scrollTo({ top: 0, behavior: 'smooth' }) }
  return <div className="app">
    <Header screen={screen} navigate={navigate} openLog={() => setLogOpen(true)} />
    {screen === 'feed' && <Feed navigate={navigate} />}
    {screen === 'detail' && <Detail navigate={navigate} openLog={() => setLogOpen(true)} />}
    {screen === 'profile' && <Profile navigate={navigate} />}
    {screen === 'list' && <ListPage navigate={navigate} />}
    <MobileNav screen={screen} navigate={navigate} openLog={() => setLogOpen(true)} />
    {logOpen && <LogModal close={() => setLogOpen(false)} />}
  </div>
}
