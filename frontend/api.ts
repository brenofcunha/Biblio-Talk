const API_URL = (import.meta.env.VITE_API_URL ?? 'http://localhost:3000').replace(/\/$/, '')
const TOKEN_KEY = 'biblio-talk-token'
const LEGACY_TOKEN_KEY = 'biblio-talk-token'

function decodeTokenPayload(token: string) {
  const encoded = token.split('.')[1]
  if (!encoded) throw new Error('Token JWT inválido')
  const normalized = encoded.replace(/-/g, '+').replace(/_/g, '/')
  const padded = normalized.padEnd(Math.ceil(normalized.length / 4) * 4, '=')
  return JSON.parse(atob(padded)) as { exp?: number }
}

export type User = {
  id: string
  username: string
  email?: string
  bio?: string | null
  avatarUrl?: string | null
  top4BookIds?: string[]
}

export type Book = {
  id: string
  title: string
  authors: string[]
  description?: string | null
  coverUrl?: string | null
  publishYear?: number | null
  publisher?: string | null
  pageCount?: number | null
  isbn?: string | null
  genres?: string[]
}

export type BookDetails = Book & {
  averageRating?: number
  ratingsCount?: number
  histogram?: Record<string, number>
  reviews?: ReadingLog[]
}

export type FeedItem = {
  id: string
  createdAt: string
  user: { username: string }
  bookId?: string
  book: { title: string; authors: string[]; coverUrl?: string | null }
  rating: number
  reviewText?: string | null
  containsSpoilers: boolean
  likesCount: number
  commentsCount?: number
}

export type FeedPage = { items: FeedItem[]; nextCursor?: string; limit: number }

export type ReadingLog = FeedItem & {
  userId: string
  bookId: string
  readDate?: string | null
  isReread: boolean
}

export type ReadingList = {
  id: string
  userId: string
  title: string
  description?: string | null
  isPrivate: boolean
  createdAt: string
}

export type ShelfEntry = {
  id: string
  userId: string
  bookId: string
  status: string
  updatedAt: string
}

type APIResponse<T> = {
  success: boolean
  message?: string
  data?: T
  error?: { code?: string; details?: string }
}

export type AuthResult = { user: User; token: string }

function getToken() {
  const token = sessionStorage.getItem(TOKEN_KEY)
  if (!token) {
    localStorage.removeItem(LEGACY_TOKEN_KEY)
    return null
  }
  try {
    const payload = decodeTokenPayload(token)
    if (payload.exp && payload.exp * 1000 <= Date.now()) {
      clearToken()
      return null
    }
  } catch {
    clearToken()
    return null
  }
  return token
}

async function request<T>(path: string, init: RequestInit = {}): Promise<T> {
  const headers = new Headers(init.headers)
  headers.set('Content-Type', 'application/json')
  const token = getToken()
  if (token) headers.set('Authorization', `Bearer ${token}`)

  const response = await fetch(`${API_URL}${path}`, { ...init, headers })
  const body = await response.json() as APIResponse<T>
  if (response.status === 401) clearToken()
  if (!response.ok || !body.success) {
    throw new Error(body.message ?? body.error?.details ?? 'Não foi possível concluir a operação.')
  }
  return body.data as T
}

export function saveToken(token: string) {
  sessionStorage.setItem(TOKEN_KEY, token)
}

export function clearToken() {
  sessionStorage.removeItem(TOKEN_KEY)
  localStorage.removeItem(LEGACY_TOKEN_KEY)
}

export function hasToken() {
  return Boolean(getToken())
}

export function logout() {
  return request<null>('/api/v1/auth/logout', { method: 'POST' })
}

export function login(email: string, password: string) {
  return request<AuthResult>('/api/v1/auth/login', {
    method: 'POST',
    body: JSON.stringify({ email, password }),
  })
}

export function register(username: string, email: string, password: string) {
  return request<AuthResult>('/api/v1/auth/register', {
    method: 'POST',
    body: JSON.stringify({ username, email, password }),
  })
}

export function getCurrentUser() {
  return request<User>('/api/v1/auth/me')
}

export function searchBooks(query: string) {
  return request<Book[]>(`/api/v1/books/search?q=${encodeURIComponent(query)}`)
}

export function getFeed() {
  return request<FeedPage>('/api/v1/feed?limit=20')
}

export function getMoreFeed(cursor: string) {
  return request<FeedPage>(`/api/v1/feed?limit=20&cursor=${encodeURIComponent(cursor)}`)
}

export function getGlobalFeed() {
  return request<FeedItem[]>('/api/v1/feed/global?limit=20')
}

export function getReviewFeed() {
  return request<FeedItem[]>('/api/v1/feed/reviews?limit=20')
}

export function getPopularBooks() {
  return request<Array<Book & { averageRating?: number }>>('/api/v1/books/popular?limit=8')
}

export function getBook(id: string) {
  return request<BookDetails>(`/api/v1/books/${id}`)
}

export function getProfile(username: string) {
  return request<User & { stats?: Record<string, unknown> }>(`/api/v1/users/${encodeURIComponent(username)}`)
}

export function getUserLogs(username: string) {
  return request<ReadingLog[]>(`/api/v1/users/${encodeURIComponent(username)}/logs`)
}

export function getUserShelf(username: string, status = 'read') {
  return request<ShelfEntry[]>(`/api/v1/users/${encodeURIComponent(username)}/shelf?status=${status}`)
}

export function getAllUserShelves(username: string) {
  return Promise.all(['tbr', 'reading', 'read', 'abandoned'].map(status => getUserShelf(username, status)))
    .then(groups => groups.flat())
}

export function getUserStats(username: string, year = new Date().getFullYear()) {
  return request<Record<string, unknown>>(`/api/v1/users/${encodeURIComponent(username)}/stats?year=${year}`)
}

export function createReadingLog(payload: {
  bookId: string
  rating: number
  reviewText?: string
  readDate?: string
  containsSpoilers: boolean
  liked: boolean
  tags?: string[]
}) {
  return request<unknown>('/api/v1/logs/', {
    method: 'POST',
    body: JSON.stringify(payload),
  })
}

export function updateShelf(bookId: string, status: string) {
  return request<unknown>('/api/v1/shelf/', {
    method: 'PUT',
    body: JSON.stringify({ bookId, status }),
  })
}

export function getMyLists() {
  return request<ReadingList[]>('/api/v1/lists/')
}

export function getTrendingLists() {
  return request<ReadingList[]>('/api/v1/lists/trending?limit=10')
}

export function createReadingList(title: string, description = '', isPrivate = false) {
  return request<ReadingList>('/api/v1/lists/', {
    method: 'POST',
    body: JSON.stringify({ title, description: description || undefined, isPrivate }),
  })
}

export function addBookToList(listId: string, bookId: string, position = 1) {
  return request<unknown>(`/api/v1/lists/${listId}/items`, {
    method: 'POST',
    body: JSON.stringify({ bookId, position }),
  })
}

export function likeLog(logId: string) {
  return request<unknown>(`/api/v1/logs/${logId}/like`, { method: 'POST' })
}

export function unlikeLog(logId: string) {
  return request<unknown>(`/api/v1/logs/${logId}/like`, { method: 'DELETE' })
}

export function createLog(payload: {
  bookId: string
  rating: number
  reviewText?: string
  readDate?: string
  containsSpoilers: boolean
  liked: boolean
}) {
  return createReadingLog(payload)
}
