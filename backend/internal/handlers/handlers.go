package handlers

import (
	"context"
	"errors"
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
	Diary    *services.DiaryService
}

func New(cfg config.Config, db *pgxpool.Pool, userRepo *repositories.UserRepository, auth *services.AuthService, users *services.UserService, books *services.BookService, logs *services.LogService, shelves *services.ShelfService, lists *services.ListService, feed *services.FeedService, stats *services.StatsService, diary *services.DiaryService) *Handler {
	return &Handler{Cfg: cfg, DB: db, UserRepo: userRepo, Auth: auth, Users: users, Books: books, Logs: logs, Shelves: shelves, Lists: lists, Feed: feed, Stats: stats, Diary: diary}
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
	u, err := h.Auth.Register(ctx, req.Username, req.Email, req.Password)
	if err != nil {
		if u != nil && errors.Is(err, services.ErrEmailUnavailable) {
			return utils.Success(c, fiber.StatusCreated, map[string]interface{}{"user": u, "verificationRequired": true, "emailSent": false}, "Conta criada, mas não foi possível enviar o e-mail. Solicite uma nova confirmação.")
		}
		return authError(c, err, "Falha ao registrar usuário")
	}
	return utils.Success(c, fiber.StatusCreated, map[string]interface{}{"user": u, "verificationRequired": true, "emailSent": true}, "Conta criada. Confirme seu e-mail para entrar.")
}

func (h *Handler) Login(c *fiber.Ctx) error {
	var req dto.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.Fail(c, fiber.StatusBadRequest, "invalidBody", "JSON inválido", err.Error())
	}
	ctx, cancel := h.withTimeout(c)
	defer cancel()
	session, err := h.Auth.Login(ctx, req.Email, req.Password)
	if err != nil {
		return authError(c, err, "Falha de autenticação")
	}
	h.setRefreshCookie(c, session.RefreshToken)
	return utils.Success(c, fiber.StatusOK, map[string]interface{}{"user": session.User, "token": session.AccessToken}, "Login realizado com sucesso")
}

func (h *Handler) Logout(c *fiber.Ctx) error {
	h.clearRefreshCookie(c)
	if err := h.Auth.RevokeSession(c.UserContext(), c.Locals("sessionId").(string)); err != nil {
		return authError(c, err, "Não foi possível encerrar a sessão")
	}
	return utils.Success(c, fiber.StatusOK, nil, "Logout realizado com sucesso")
}

func (h *Handler) ListAuthSessions(c *fiber.Ctx) error {
	ctx, cancel := h.withTimeout(c)
	defer cancel()
	sessions, err := h.Auth.ListActiveSessions(ctx, c.Locals("userId").(string), c.Locals("sessionId").(string))
	if err != nil {
		return authError(c, err, "Não foi possível listar as sessões")
	}
	return utils.Success(c, fiber.StatusOK, sessions, "Sessões carregadas")
}

func (h *Handler) RevokeAuthSession(c *fiber.Ctx) error {
	ctx, cancel := h.withTimeout(c)
	defer cancel()
	userID := c.Locals("userId").(string)
	sessionID := c.Params("id")
	if err := h.Auth.RevokeUserSession(ctx, userID, sessionID); err != nil {
		return authError(c, err, "Não foi possível revogar a sessão")
	}
	if sessionID == c.Locals("sessionId").(string) {
		h.clearRefreshCookie(c)
	}
	return utils.Success(c, fiber.StatusOK, nil, "Sessão revogada")
}

func (h *Handler) RevokeOtherAuthSessions(c *fiber.Ctx) error {
	ctx, cancel := h.withTimeout(c)
	defer cancel()
	revoked, err := h.Auth.RevokeOtherSessions(ctx, c.Locals("userId").(string), c.Locals("sessionId").(string))
	if err != nil {
		return authError(c, err, "Não foi possível revogar as outras sessões")
	}
	return utils.Success(c, fiber.StatusOK, map[string]int64{"revokedCount": revoked}, "Outras sessões revogadas")
}

func (h *Handler) Refresh(c *fiber.Ctx) error {
	ctx, cancel := h.withTimeout(c)
	defer cancel()
	session, err := h.Auth.Refresh(ctx, c.Cookies("biblio-talk-refresh"))
	if err != nil {
		h.clearRefreshCookie(c)
		return authError(c, err, "Não foi possível renovar a sessão")
	}
	h.setRefreshCookie(c, session.RefreshToken)
	return utils.Success(c, fiber.StatusOK, map[string]interface{}{"user": session.User, "token": session.AccessToken}, "Sessão renovada")
}

func (h *Handler) RequestPasswordReset(c *fiber.Ctx) error {
	var req dto.EmailRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.Fail(c, fiber.StatusBadRequest, "invalidBody", "JSON inválido", err.Error())
	}
	if err := services.ValidateEmail(req.Email); err != nil {
		return authError(c, err, "E-mail inválido")
	}
	ctx, cancel := h.withTimeout(c)
	defer cancel()
	if err := h.Auth.RequestPasswordReset(ctx, req.Email); err != nil {
		return authError(c, err, "Não foi possível solicitar a redefinição")
	}
	return utils.Success(c, fiber.StatusOK, nil, "Se o e-mail estiver cadastrado, enviaremos um link para redefinir a senha.")
}

func (h *Handler) ResetPassword(c *fiber.Ctx) error {
	var req dto.ResetPasswordRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.Fail(c, fiber.StatusBadRequest, "invalidBody", "JSON inválido", err.Error())
	}
	ctx, cancel := h.withTimeout(c)
	defer cancel()
	if err := h.Auth.ResetPassword(ctx, req.Token, req.Password); err != nil {
		return authError(c, err, "Não foi possível redefinir a senha")
	}
	return utils.Success(c, fiber.StatusOK, nil, "Senha redefinida. Entre novamente.")
}

func (h *Handler) ChangePassword(c *fiber.Ctx) error {
	var req dto.ChangePasswordRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.Fail(c, fiber.StatusBadRequest, "invalidBody", "JSON inválido", err.Error())
	}
	ctx, cancel := h.withTimeout(c)
	defer cancel()
	err := h.Auth.ChangePassword(ctx, c.Locals("userId").(string), c.Locals("sessionId").(string), req.CurrentPassword, req.NewPassword)
	if err != nil {
		return authError(c, err, "Não foi possível trocar a senha")
	}
	return utils.Success(c, fiber.StatusOK, nil, "Senha alterada. Outras sessões foram encerradas.")
}

func (h *Handler) RequestEmailVerification(c *fiber.Ctx) error {
	var req dto.EmailRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.Fail(c, fiber.StatusBadRequest, "invalidBody", "JSON inválido", err.Error())
	}
	if err := services.ValidateEmail(req.Email); err != nil {
		return authError(c, err, "E-mail inválido")
	}
	ctx, cancel := h.withTimeout(c)
	defer cancel()
	if err := h.Auth.RequestEmailVerification(ctx, req.Email); err != nil {
		return authError(c, err, "Não foi possível enviar a confirmação")
	}
	return utils.Success(c, fiber.StatusOK, nil, "Se a conta puder ser confirmada, enviaremos um novo link.")
}

func (h *Handler) VerifyEmail(c *fiber.Ctx) error {
	var req dto.TokenRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.Fail(c, fiber.StatusBadRequest, "invalidBody", "JSON inválido", err.Error())
	}
	ctx, cancel := h.withTimeout(c)
	defer cancel()
	if err := h.Auth.VerifyEmail(ctx, req.Token); err != nil {
		return authError(c, err, "Não foi possível confirmar o e-mail")
	}
	return utils.Success(c, fiber.StatusOK, nil, "E-mail confirmado. Você já pode entrar.")
}

func authError(c *fiber.Ctx, err error, fallback string) error {
	switch {
	case errors.Is(err, services.ErrValidation):
		return utils.Fail(c, fiber.StatusBadRequest, "invalidAuthInput", err.Error(), "")
	case errors.Is(err, services.ErrInvalidCredentials):
		return utils.Fail(c, fiber.StatusUnauthorized, "invalidCredentials", "E-mail ou senha inválidos", "")
	case errors.Is(err, services.ErrConflict):
		return utils.Fail(c, fiber.StatusConflict, "accountConflict", "Nome de usuário ou e-mail já cadastrado", "")
	case errors.Is(err, services.ErrEmailNotVerified):
		return utils.Fail(c, fiber.StatusForbidden, "emailNotVerified", "Confirme seu e-mail antes de entrar", "")
	case errors.Is(err, services.ErrInvalidToken):
		return utils.Fail(c, fiber.StatusBadRequest, "invalidToken", "Link inválido ou expirado", "")
	case errors.Is(err, services.ErrInvalidSession):
		return utils.Fail(c, fiber.StatusUnauthorized, "invalidSession", "Sessão inválida ou expirada", "")
	case errors.Is(err, services.ErrEmailUnavailable):
		return utils.Fail(c, fiber.StatusServiceUnavailable, "emailUnavailable", "Não foi possível enviar o e-mail. Tente novamente mais tarde.", "")
	default:
		return utils.Fail(c, fiber.StatusInternalServerError, "authError", fallback, "")
	}
}

func (h *Handler) setRefreshCookie(c *fiber.Ctx, token string) {
	lifetime := int(h.Cfg.RefreshTokenExpiration.Seconds())
	c.Cookie(&fiber.Cookie{
		Name:     "biblio-talk-refresh",
		Value:    token,
		Path:     "/api/v1/auth",
		HTTPOnly: true,
		Secure:   h.Cfg.RefreshCookieSecure,
		SameSite: "Lax",
		MaxAge:   lifetime,
		Expires:  time.Now().Add(h.Cfg.RefreshTokenExpiration),
	})
}

func (h *Handler) clearRefreshCookie(c *fiber.Ctx) {
	c.Cookie(&fiber.Cookie{
		Name:     "biblio-talk-refresh",
		Value:    "",
		Path:     "/api/v1/auth",
		HTTPOnly: true,
		Secure:   h.Cfg.RefreshCookieSecure,
		SameSite: "Lax",
		MaxAge:   -1,
		Expires:  time.Unix(1, 0),
	})
}

func (h *Handler) Me(c *fiber.Ctx) error {
	ctx, cancel := h.withTimeout(c)
	defer cancel()
	u, err := h.Users.GetMyProfile(ctx, c.Locals("userId").(string), time.Now().Year())
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

func (h *Handler) UpdateProfile(c *fiber.Ctx) error {
	var req dto.UpdateProfileRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.Fail(c, fiber.StatusBadRequest, "invalidBody", "JSON inválido", err.Error())
	}
	ctx, cancel := h.withTimeout(c)
	defer cancel()
	err := h.Users.UpdateProfile(ctx, c.Locals("userId").(string), models.User{
		DisplayName: req.DisplayName, Bio: req.Bio, AvatarURL: req.AvatarURL, BannerURL: req.BannerURL, Location: req.Location,
	})
	if err != nil {
		return userError(c, err, "Falha ao atualizar perfil")
	}
	return utils.Success(c, fiber.StatusOK, nil, "Perfil atualizado")
}

func (h *Handler) UpdatePrivacy(c *fiber.Ctx) error {
	var req dto.UpdatePrivacyRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.Fail(c, fiber.StatusBadRequest, "invalidBody", "JSON inválido", err.Error())
	}
	ctx, cancel := h.withTimeout(c)
	defer cancel()
	settings := models.PrivacySettings{
		Profile: req.Profile, Diary: req.Diary, Progress: req.Progress, Shelf: req.Shelf,
		Notes: req.Notes, Comments: req.Comments, Lists: req.Lists, Stats: req.Stats, Favorites: req.Favorites,
	}
	if err := h.Users.UpdatePrivacy(ctx, c.Locals("userId").(string), settings); err != nil {
		return userError(c, err, "Configuração de privacidade inválida")
	}
	return utils.Success(c, fiber.StatusOK, settings, "Privacidade atualizada")
}

func (h *Handler) ExportMyData(c *fiber.Ctx) error {
	ctx, cancel := h.withTimeout(c)
	defer cancel()
	data, err := h.Users.Export(ctx, c.Locals("userId").(string))
	if err != nil {
		return userError(c, err, "Falha ao exportar dados")
	}
	c.Set(fiber.HeaderContentDisposition, `attachment; filename="biblio-talk-data.json"`)
	return utils.Success(c, fiber.StatusOK, data, "Exportação de dados concluída")
}

func (h *Handler) DeleteMe(c *fiber.Ctx) error {
	ctx, cancel := h.withTimeout(c)
	defer cancel()
	if err := h.Users.DeleteAccount(ctx, c.Locals("userId").(string)); err != nil {
		return userError(c, err, "Falha ao excluir conta")
	}
	h.clearRefreshCookie(c)
	return utils.Success(c, fiber.StatusOK, nil, "Conta e dados associados foram excluídos")
}

func (h *Handler) RecordConsent(c *fiber.Ctx) error {
	var req dto.ConsentRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.Fail(c, fiber.StatusBadRequest, "invalidBody", "JSON inválido", err.Error())
	}
	ctx, cancel := h.withTimeout(c)
	defer cancel()
	if req.Granted == nil {
		return utils.Fail(c, fiber.StatusBadRequest, "invalidUserData", "Consentimento inválido", "informe granted como true ou false")
	}
	consent, err := h.Users.RecordConsent(ctx, c.Locals("userId").(string), models.Consent{
		ConsentType: req.ConsentType, Version: req.Version, Granted: *req.Granted,
	})
	if err != nil {
		return userError(c, err, "Consentimento inválido")
	}
	return utils.Success(c, fiber.StatusCreated, consent, "Consentimento registrado")
}

func (h *Handler) UpdateFavorites(c *fiber.Ctx) error {
	var req dto.UpdateFavoritesRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.Fail(c, fiber.StatusBadRequest, "invalidBody", "JSON inválido", err.Error())
	}
	ctx, cancel := h.withTimeout(c)
	defer cancel()
	if err := h.Users.UpdateFavorites(ctx, c.Locals("userId").(string), req.BookIDs); err != nil {
		return userError(c, err, "Não foi possível atualizar livros favoritos")
	}
	return utils.Success(c, fiber.StatusOK, req.BookIDs, "Livros favoritos atualizados")
}

func (h *Handler) PublicFavorites(c *fiber.Ctx) error {
	ctx, cancel := h.withTimeout(c)
	defer cancel()
	data, err := h.Users.PublicFavorites(ctx, c.Params("username"))
	if err != nil {
		return mapError(c, err, "Favoritos não encontrados")
	}
	return utils.Success(c, fiber.StatusOK, data, "Livros favoritos carregados")
}

func (h *Handler) PublicUserLists(c *fiber.Ctx) error {
	ctx, cancel := h.withTimeout(c)
	defer cancel()
	data, err := h.Users.PublicLists(ctx, c.Params("username"))
	if err != nil {
		return mapError(c, err, "Listas públicas não encontradas")
	}
	return utils.Success(c, fiber.StatusOK, data, "Listas públicas carregadas")
}

func (h *Handler) UpdateAnnualChallenge(c *fiber.Ctx) error {
	var req dto.AnnualChallengeRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.Fail(c, fiber.StatusBadRequest, "invalidBody", "JSON inválido", err.Error())
	}
	ctx, cancel := h.withTimeout(c)
	defer cancel()
	challenge := models.AnnualChallenge{Year: req.Year, TargetBooks: req.TargetBooks, IsPublic: req.IsPublic}
	if err := h.Users.SaveAnnualChallenge(ctx, c.Locals("userId").(string), challenge); err != nil {
		return userError(c, err, "Desafio anual inválido")
	}
	return utils.Success(c, fiber.StatusOK, challenge, "Desafio anual atualizado")
}

func (h *Handler) PublicAnnualChallenge(c *fiber.Ctx) error {
	year := parseYear(c)
	if c.Query("year") != "" {
		parsed, err := strconv.Atoi(c.Query("year"))
		if err != nil {
			return utils.Fail(c, fiber.StatusBadRequest, "invalidYear", "Ano inválido", "use um ano entre 2000 e 2100")
		}
		year = parsed
	}
	ctx, cancel := h.withTimeout(c)
	defer cancel()
	challenge, err := h.Users.PublicAnnualChallenge(ctx, c.Params("username"), year)
	if err != nil {
		return userError(c, err, "Desafio anual não encontrado")
	}
	return utils.Success(c, fiber.StatusOK, challenge, "Desafio anual carregado")
}

func userError(c *fiber.Ctx, err error, fallback string) error {
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return utils.Fail(c, fiber.StatusNotFound, "notFound", fallback, "")
	case errors.Is(err, services.ErrInvalidProfile), errors.Is(err, services.ErrInvalidPrivacy),
		errors.Is(err, services.ErrInvalidConsent), errors.Is(err, services.ErrInvalidChallenge):
		return utils.Fail(c, fiber.StatusBadRequest, "invalidUserData", err.Error(), "")
	default:
		return utils.Fail(c, fiber.StatusInternalServerError, "userOperationFailed", fallback, "")
	}
}

func (h *Handler) Follow(c *fiber.Ctx) error {
	ctx, cancel := h.withTimeout(c)
	defer cancel()
	created, err := h.Users.Follow(ctx, c.Locals("userId").(string), c.Params("id"))
	if err != nil {
		return socialError(c, err, "Falha ao seguir usuário")
	}
	status := fiber.StatusCreated
	message := "Agora você segue este usuário"
	if !created {
		status = fiber.StatusOK
		message = "Você já segue este usuário"
	}
	return utils.Success(c, status, map[string]bool{"following": true}, message)
}

func (h *Handler) Unfollow(c *fiber.Ctx) error {
	ctx, cancel := h.withTimeout(c)
	defer cancel()
	removed, err := h.Users.Unfollow(ctx, c.Locals("userId").(string), c.Params("id"))
	if err != nil {
		return socialError(c, err, "Falha ao deixar de seguir usuário")
	}
	message := "Você deixou de seguir este usuário"
	if !removed {
		message = "Você não seguia este usuário"
	}
	return utils.Success(c, fiber.StatusOK, map[string]bool{"following": false}, message)
}

func socialError(c *fiber.Ctx, err error, fallback string) error {
	switch {
	case errors.Is(err, services.ErrSelfFollow):
		return utils.Fail(c, fiber.StatusBadRequest, "selfFollow", err.Error(), "")
	case errors.Is(err, services.ErrSocialUserNotFound), errors.Is(err, pgx.ErrNoRows):
		return utils.Fail(c, fiber.StatusNotFound, "userNotFound", "Usuário não encontrado", "")
	default:
		return utils.Fail(c, fiber.StatusInternalServerError, "socialOperationFailed", fallback, "")
	}
}

func (h *Handler) Followers(c *fiber.Ctx) error {
	ctx, cancel := h.withTimeout(c)
	defer cancel()
	canView, err := h.Users.CanViewByID(ctx, c.Params("id"), "profile")
	if err != nil {
		return mapError(c, err, "Usuário não encontrado")
	}
	if !canView {
		return utils.Fail(c, fiber.StatusNotFound, "notFound", "Usuário não encontrado", "")
	}
	data, err := h.UserRepo.ListFollowers(ctx, c.Params("id"))
	if err != nil {
		return mapError(c, err, "Falha ao listar seguidores")
	}
	return utils.Success(c, fiber.StatusOK, data, "Seguidores carregados")
}

func (h *Handler) Following(c *fiber.Ctx) error {
	ctx, cancel := h.withTimeout(c)
	defer cancel()
	canView, err := h.Users.CanViewByID(ctx, c.Params("id"), "profile")
	if err != nil {
		return mapError(c, err, "Usuário não encontrado")
	}
	if !canView {
		return utils.Fail(c, fiber.StatusNotFound, "notFound", "Usuário não encontrado", "")
	}
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
	log, err := h.Logs.Get(ctx, c.Locals("userId").(string), c.Params("id"))
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
	canView, err := h.Users.CanView(ctx, c.Params("username"), "diary")
	if err != nil {
		return mapError(c, err, "Usuário não encontrado")
	}
	if !canView {
		return utils.Fail(c, fiber.StatusNotFound, "notFound", "Diário não encontrado", "")
	}
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

func (h *Handler) RemoveShelfBook(c *fiber.Ctx) error {
	ctx, cancel := h.withTimeout(c)
	defer cancel()
	if err := h.Shelves.Remove(ctx, c.Locals("userId").(string), c.Params("bookId")); err != nil {
		return mapError(c, err, "Falha ao remover livro da estante")
	}
	return utils.Success(c, fiber.StatusOK, nil, "Livro removido da estante")
}

func (h *Handler) UserShelf(c *fiber.Ctx) error {
	status := c.Query("status", "reading")
	ctx, cancel := h.withTimeout(c)
	defer cancel()
	canView, err := h.Users.CanView(ctx, c.Params("username"), "shelf")
	if err != nil {
		return mapError(c, err, "Usuário não encontrado")
	}
	if !canView {
		return utils.Fail(c, fiber.StatusNotFound, "notFound", "Estante não encontrada", "")
	}
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
	viewerID, _ := c.Locals("userId").(string)
	data, err := h.Lists.Detail(ctx, c.Params("id"), viewerID)
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
	if err := h.Lists.UpsertItem(ctx, c.Params("id"), c.Locals("userId").(string), req.BookID, req.Position, req.Notes); err != nil {
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
	if err := h.Lists.DeleteItem(ctx, c.Params("id"), c.Locals("userId").(string), req.BookID); err != nil {
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
	canView, err := h.Users.CanView(ctx, c.Params("username"), "stats")
	if err != nil {
		return mapError(c, err, "Usuário não encontrado")
	}
	if !canView {
		return utils.Fail(c, fiber.StatusNotFound, "notFound", "Estatísticas não encontradas", "")
	}
	stats, err := h.Stats.Year(ctx, c.Params("username"), parseYear(c))
	if err != nil {
		return mapError(c, err, "Falha ao carregar estatísticas")
	}
	return utils.Success(c, fiber.StatusOK, stats, "Estatísticas carregadas")
}
