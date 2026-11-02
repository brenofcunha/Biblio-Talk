package handlers

import (
	"context"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"lume_api/internal/config"
	"lume_api/internal/dto"
	"lume_api/internal/models"
	"lume_api/internal/repositories"
	"lume_api/internal/services"
	"lume_api/internal/utils"
)

type Handler struct {
	Cfg      config.Config
	DB       *pgxpool.Pool
	UserRepo *repositories.UserRepository
	Auth     *services.AuthService
	Users    *services.UserService
	Books    *services.BookService
	Logs     *services.LogService
	Shelves  *services.ShelfService
	Lists    *services.ListService
	Feed     *services.FeedService
	Stats    *services.StatsService
}

func New(cfg config.Config, db *pgxpool.Pool, userRepo *repositories.UserRepository, auth *services.AuthService, users *services.UserService, books *services.BookService, logs *services.LogService, shelves *services.ShelfService, lists *services.ListService, feed *services.FeedService, stats *services.StatsService) *Handler {
	return &Handler{Cfg: cfg, DB: db, UserRepo: userRepo, Auth: auth, Users: users, Books: books, Logs: logs, Shelves: shelves, Lists: lists, Feed: feed, Stats: stats}
}

func (h *Handler) withTimeout(c *fiber.Ctx) (context.Context, context.CancelFunc) {
	return context.WithTimeout(c.UserContext(), time.Duration(h.Cfg.RequestTimeoutSec)*time.Second)
}

func parseLimit(c *fiber.Ctx, d int) int {
	if c.Query("limit") == "" {
		return d
	}
	v, err := strconv.Atoi(c.Query("limit"))
	if err != nil || v <= 0 || v > 100 {
		return d
	}
	return v
}

func parseYear(c *fiber.Ctx) int {
	if c.Query("year") == "" {
		return time.Now().Year()
	}
	y, err := strconv.Atoi(c.Query("year"))
	if err != nil || y < 2000 || y > 2100 {
		return time.Now().Year()
	}
	return y
}

func mapError(c *fiber.Ctx, err error, fallback string) error {
	if err == nil {
		return nil
	}
	if err == pgx.ErrNoRows {
		return utils.Fail(c, fiber.StatusNotFound, "notFound", fallback, "recurso não encontrado")
	}
	return utils.Fail(c, fiber.StatusBadRequest, "badRequest", fallback, err.Error())
}

func (h *Handler) Healthz(c *fiber.Ctx) error {
	ctx, cancel := h.withTimeout(c)
	defer cancel()
	if err := h.DB.Ping(ctx); err != nil {
		return utils.Fail(c, fiber.StatusServiceUnavailable, "dbUnavailable", "Banco indisponível", err.Error())
	}
	return utils.Success(c, fiber.StatusOK, map[string]string{"status": "ok"}, "Serviço saudável")
}

func (h *Handler) Register(c *fiber.Ctx) error {
	var req dto.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.Fail(c, fiber.StatusBadRequest, "invalidBody", "JSON inválido", err.Error())
	}
	ctx, cancel := h.withTimeout(c)
	defer cancel()
	u, token, err := h.Auth.Register(ctx, req.Username, req.Email, req.Password)
	if err != nil {
		return mapError(c, err, "Falha ao registrar usuário")
	}
	return utils.Success(c, fiber.StatusCreated, map[string]interface{}{"user": u, "token": token}, "Usuário registrado com sucesso")
}

func (h *Handler) Login(c *fiber.Ctx) error {
	var req dto.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.Fail(c, fiber.StatusBadRequest, "invalidBody", "JSON inválido", err.Error())
	}
	ctx, cancel := h.withTimeout(c)
	defer cancel()
	u, token, err := h.Auth.Login(ctx, req.Email, req.Password)
	if err != nil {
		return mapError(c, err, "Falha de autenticação")
	}
	return utils.Success(c, fiber.StatusOK, map[string]interface{}{"user": u, "token": token}, "Login realizado com sucesso")
}

func (h *Handler) Me(c *fiber.Ctx) error {
	ctx, cancel := h.withTimeout(c)
	defer cancel()
	u, err := h.Users.GetPublicProfile(ctx, c.Locals("username").(string), time.Now().Year())
	if err != nil {
		return mapError(c, err, "Falha ao carregar usuário")
	}
	return utils.Success(c, fiber.StatusOK, u, "Perfil carregado")
}

func (h *Handler) PublicProfile(c *fiber.Ctx) error {
	ctx, cancel := h.withTimeout(c)
	defer cancel()
	u, err := h.Users.GetPublicProfile(ctx, c.Params("username"), parseYear(c))
	if err != nil {
		return mapError(c, err, "Falha ao carregar perfil")
	}
	return utils.Success(c, fiber.StatusOK, u, "Perfil carregado")
}

func (h *Handler) UpdateMe(c *fiber.Ctx) error {
	var req dto.UpdateMeRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.Fail(c, fiber.StatusBadRequest, "invalidBody", "JSON inválido", err.Error())
	}
	ctx, cancel := h.withTimeout(c)
	defer cancel()
	if err := h.Users.UpdateMe(ctx, c.Locals("userId").(string), req.Bio, req.AvatarURL, req.Top4BookIDs); err != nil {
		return mapError(c, err, "Falha ao atualizar perfil")
	}
	return utils.Success(c, fiber.StatusOK, nil, "Perfil atualizado com sucesso")
}

func (h *Handler) Follow(c *fiber.Ctx) error {
	ctx, cancel := h.withTimeout(c)
	defer cancel()
	if err := h.UserRepo.Follow(ctx, c.Locals("userId").(string), c.Params("id")); err != nil {
		return mapError(c, err, "Falha ao seguir usuário")
	}
	return utils.Success(c, fiber.StatusOK, nil, "Agora você segue este usuário")
}

func (h *Handler) Unfollow(c *fiber.Ctx) error {
	ctx, cancel := h.withTimeout(c)
	defer cancel()
	if err := h.UserRepo.Unfollow(ctx, c.Locals("userId").(string), c.Params("id")); err != nil {
		return mapError(c, err, "Falha ao deixar de seguir usuário")
	}
	return utils.Success(c, fiber.StatusOK, nil, "Você deixou de seguir este usuário")
}

func (h *Handler) Followers(c *fiber.Ctx) error {
	ctx, cancel := h.withTimeout(c)
	defer cancel()
	data, err := h.UserRepo.ListFollowers(ctx, c.Params("id"))
	if err != nil {
		return mapError(c, err, "Falha ao listar seguidores")
	}
	return utils.Success(c, fiber.StatusOK, data, "Seguidores carregados")
}

func (h *Handler) Following(c *fiber.Ctx) error {
	ctx, cancel := h.withTimeout(c)
	defer cancel()
	data, err := h.UserRepo.ListFollowing(ctx, c.Params("id"))
	if err != nil {
		return mapError(c, err, "Falha ao listar seguindo")
	}
	return utils.Success(c, fiber.StatusOK, data, "Seguindo carregado")
}

func (h *Handler) SearchBooks(c *fiber.Ctx) error {
	q := c.Query("q")
	if q == "" {
		return utils.Fail(c, fiber.StatusBadRequest, "invalidQuery", "Parâmetro q é obrigatório", "use /books/search?q=... ")
	}
	ctx, cancel := h.withTimeout(c)
	defer cancel()
	books, err := h.Books.Search(ctx, q, parseLimit(c, 10))
	if err != nil {
		return mapError(c, err, "Falha na busca de livros")
	}
	return utils.Success(c, fiber.StatusOK, books, "Busca concluída")
}

func (h *Handler) BookDetail(c *fiber.Ctx) error {
	ctx, cancel := h.withTimeout(c)
	defer cancel()
	data, err := h.Books.GetDetail(ctx, c.Params("id"))
	if err != nil {
		return mapError(c, err, "Falha ao carregar livro")
	}
	return utils.Success(c, fiber.StatusOK, data, "Livro carregado")
}

func (h *Handler) PopularBooks(c *fiber.Ctx) error {
	ctx, cancel := h.withTimeout(c)
	defer cancel()
	data, err := h.Books.Popular(ctx, parseLimit(c, 3))
	if err != nil {
		return mapError(c, err, "Falha ao listar populares")
	}
	return utils.Success(c, fiber.StatusOK, data, "Populares carregados")
}

func (h *Handler) CreateLog(c *fiber.Ctx) error {
	var req dto.CreateLogRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.Fail(c, fiber.StatusBadRequest, "invalidBody", "JSON inválido", err.Error())
	}
	ctx, cancel := h.withTimeout(c)
	defer cancel()
	created, err := h.Logs.Create(ctx, c.Locals("userId").(string), c.Locals("username").(string), models.Log{
		BookID: req.BookID, Rating: req.Rating, ReviewText: req.ReviewText, IsReread: req.IsReread, ContainsSpoilers: req.ContainsSpoilers,
	}, req.ReadDate)
	if err != nil {
		return mapError(c, err, "Falha ao criar log")
	}
	return utils.Success(c, fiber.StatusCreated, created, "Leitura registrada com sucesso")
}

func (h *Handler) GetLog(c *fiber.Ctx) error {
	ctx, cancel := h.withTimeout(c)
	defer cancel()
	log, err := h.Logs.Get(ctx, c.Params("id"))
	if err != nil {
		return mapError(c, err, "Falha ao carregar log")
	}
	return utils.Success(c, fiber.StatusOK, log, "Log carregado")
}

func (h *Handler) UpdateLog(c *fiber.Ctx) error {
	var req dto.UpdateLogRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.Fail(c, fiber.StatusBadRequest, "invalidBody", "JSON inválido", err.Error())
	}
	ctx, cancel := h.withTimeout(c)
	defer cancel()
	err := h.Logs.Update(ctx, c.Locals("userId").(string), c.Locals("username").(string), c.Params("id"), models.Log{Rating: req.Rating, ReviewText: req.ReviewText, IsReread: req.IsReread, ContainsSpoilers: req.ContainsSpoilers}, req.ReadDate)
	if err != nil {
		return mapError(c, err, "Falha ao atualizar log")
	}
	return utils.Success(c, fiber.StatusOK, nil, "Log atualizado")
}

func (h *Handler) DeleteLog(c *fiber.Ctx) error {
	ctx, cancel := h.withTimeout(c)
	defer cancel()
	if err := h.Logs.Delete(ctx, c.Locals("userId").(string), c.Locals("username").(string), c.Params("id")); err != nil {
		return mapError(c, err, "Falha ao excluir log")
	}
	return utils.Success(c, fiber.StatusOK, nil, "Log excluído")
}

func (h *Handler) LikeLog(c *fiber.Ctx) error {
	ctx, cancel := h.withTimeout(c)
	defer cancel()
	if err := h.Logs.Like(ctx, c.Locals("userId").(string), c.Params("id")); err != nil {
		return mapError(c, err, "Falha ao curtir log")
	}
	return utils.Success(c, fiber.StatusOK, nil, "Log curtido")
}

func (h *Handler) UnlikeLog(c *fiber.Ctx) error {
	ctx, cancel := h.withTimeout(c)
	defer cancel()
	if err := h.Logs.Unlike(ctx, c.Locals("userId").(string), c.Params("id")); err != nil {
		return mapError(c, err, "Falha ao remover curtida")
	}
	return utils.Success(c, fiber.StatusOK, nil, "Curtida removida")
}

func (h *Handler) AddComment(c *fiber.Ctx) error {
	var req dto.CreateCommentRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.Fail(c, fiber.StatusBadRequest, "invalidBody", "JSON inválido", err.Error())
	}
	ctx, cancel := h.withTimeout(c)
	defer cancel()
	comment, err := h.Logs.AddComment(ctx, c.Locals("userId").(string), c.Params("id"), req.CommentText)
	if err != nil {
		return mapError(c, err, "Falha ao comentar")
	}
	return utils.Success(c, fiber.StatusCreated, comment, "Comentário adicionado")
}

func (h *Handler) ListComments(c *fiber.Ctx) error {
	ctx, cancel := h.withTimeout(c)
	defer cancel()
	comments, err := h.Logs.Comments(ctx, c.Params("id"))
	if err != nil {
		return mapError(c, err, "Falha ao listar comentários")
	}
	return utils.Success(c, fiber.StatusOK, comments, "Comentários carregados")
}

func (h *Handler) UserLogs(c *fiber.Ctx) error {
	ctx, cancel := h.withTimeout(c)
	defer cancel()
	logs, err := h.Logs.UserLogs(ctx, c.Params("username"), parseLimit(c, 20))
	if err != nil {
		return mapError(c, err, "Falha ao listar diário")
	}
	return utils.Success(c, fiber.StatusOK, logs, "Diário carregado")
}

func (h *Handler) UpsertShelf(c *fiber.Ctx) error {
	var req dto.ShelfUpsertRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.Fail(c, fiber.StatusBadRequest, "invalidBody", "JSON inválido", err.Error())
	}
	ctx, cancel := h.withTimeout(c)
	defer cancel()
	if err := h.Shelves.Upsert(ctx, c.Locals("userId").(string), req.BookID, req.Status); err != nil {
		return mapError(c, err, "Falha ao atualizar estante")
	}
	return utils.Success(c, fiber.StatusOK, nil, "Estante atualizada")
}

func (h *Handler) UserShelf(c *fiber.Ctx) error {
	status := c.Query("status", "reading")
	ctx, cancel := h.withTimeout(c)
	defer cancel()
	data, err := h.Shelves.ByStatus(ctx, c.Params("username"), status)
	if err != nil {
		return mapError(c, err, "Falha ao listar estante")
	}
	return utils.Success(c, fiber.StatusOK, data, "Estante carregada")
}

func (h *Handler) CreateList(c *fiber.Ctx) error {
	var req dto.CreateListRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.Fail(c, fiber.StatusBadRequest, "invalidBody", "JSON inválido", err.Error())
	}
	ctx, cancel := h.withTimeout(c)
	defer cancel()
	list, err := h.Lists.Create(ctx, c.Locals("userId").(string), req.Title, req.Description, req.IsPrivate)
	if err != nil {
		return mapError(c, err, "Falha ao criar lista")
	}
	return utils.Success(c, fiber.StatusCreated, list, "Lista criada")
}

func (h *Handler) MyLists(c *fiber.Ctx) error {
	ctx, cancel := h.withTimeout(c)
	defer cancel()
	data, err := h.Lists.Mine(ctx, c.Locals("userId").(string))
	if err != nil {
		return mapError(c, err, "Falha ao listar listas")
	}
	return utils.Success(c, fiber.StatusOK, data, "Listas carregadas")
}

func (h *Handler) ListDetail(c *fiber.Ctx) error {
	ctx, cancel := h.withTimeout(c)
	defer cancel()
	data, err := h.Lists.Detail(ctx, c.Params("id"))
	if err != nil {
		return mapError(c, err, "Falha ao carregar lista")
	}
	return utils.Success(c, fiber.StatusOK, data, "Lista carregada")
}

func (h *Handler) UpdateList(c *fiber.Ctx) error {
	var req dto.UpdateListRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.Fail(c, fiber.StatusBadRequest, "invalidBody", "JSON inválido", err.Error())
	}
	ctx, cancel := h.withTimeout(c)
	defer cancel()
	if err := h.Lists.Update(ctx, c.Params("id"), c.Locals("userId").(string), req.Title, req.Description, req.IsPrivate); err != nil {
		return mapError(c, err, "Falha ao atualizar lista")
	}
	return utils.Success(c, fiber.StatusOK, nil, "Lista atualizada")
}

func (h *Handler) DeleteList(c *fiber.Ctx) error {
	ctx, cancel := h.withTimeout(c)
	defer cancel()
	if err := h.Lists.Delete(ctx, c.Params("id"), c.Locals("userId").(string)); err != nil {
		return mapError(c, err, "Falha ao excluir lista")
	}
	return utils.Success(c, fiber.StatusOK, nil, "Lista excluída")
}

func (h *Handler) UpsertListItem(c *fiber.Ctx) error {
	var req dto.UpsertListItemRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.Fail(c, fiber.StatusBadRequest, "invalidBody", "JSON inválido", err.Error())
	}
	ctx, cancel := h.withTimeout(c)
	defer cancel()
	if err := h.Lists.UpsertItem(ctx, c.Params("id"), req.BookID, req.Position, req.Notes); err != nil {
		return mapError(c, err, "Falha ao adicionar item")
	}
	return utils.Success(c, fiber.StatusOK, nil, "Item salvo")
}

func (h *Handler) DeleteListItem(c *fiber.Ctx) error {
	var req dto.UpsertListItemRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.Fail(c, fiber.StatusBadRequest, "invalidBody", "JSON inválido", err.Error())
	}
	ctx, cancel := h.withTimeout(c)
	defer cancel()
	if err := h.Lists.DeleteItem(ctx, c.Params("id"), req.BookID); err != nil {
		return mapError(c, err, "Falha ao remover item")
	}
	return utils.Success(c, fiber.StatusOK, nil, "Item removido")
}

func (h *Handler) TrendingLists(c *fiber.Ctx) error {
	ctx, cancel := h.withTimeout(c)
	defer cancel()
	data, err := h.Lists.Trending(ctx, parseLimit(c, 10))
	if err != nil {
		return mapError(c, err, "Falha ao carregar listas em alta")
	}
	return utils.Success(c, fiber.StatusOK, data, "Listas em alta")
}

func (h *Handler) FeedFollowing(c *fiber.Ctx) error {
	ctx, cancel := h.withTimeout(c)
	defer cancel()
	limit := parseLimit(c, 20)
	data, nextCursor, err := h.Feed.Following(ctx, c.Locals("userId").(string), limit, c.Query("cursor"))
	if err != nil {
		return mapError(c, err, "Falha ao carregar feed")
	}
	return utils.Success(c, fiber.StatusOK, map[string]interface{}{"items": data, "nextCursor": nextCursor, "limit": limit}, "Feed carregado")
}

func (h *Handler) FeedGlobal(c *fiber.Ctx) error {
	ctx, cancel := h.withTimeout(c)
	defer cancel()
	data, err := h.Feed.Global(ctx, parseLimit(c, 20))
	if err != nil {
		return mapError(c, err, "Falha ao carregar feed global")
	}
	return utils.Success(c, fiber.StatusOK, data, "Feed global carregado")
}

func (h *Handler) FeedReviews(c *fiber.Ctx) error {
	ctx, cancel := h.withTimeout(c)
	defer cancel()
	data, err := h.Feed.Reviews(ctx, parseLimit(c, 20))
	if err != nil {
		return mapError(c, err, "Falha ao carregar feed de resenhas")
	}
	return utils.Success(c, fiber.StatusOK, data, "Feed de resenhas carregado")
}

func (h *Handler) UserStats(c *fiber.Ctx) error {
	ctx, cancel := h.withTimeout(c)
	defer cancel()
	stats, err := h.Stats.Year(ctx, c.Params("username"), parseYear(c))
	if err != nil {
		return mapError(c, err, "Falha ao carregar estatísticas")
	}
	return utils.Success(c, fiber.StatusOK, stats, "Estatísticas carregadas")
}
