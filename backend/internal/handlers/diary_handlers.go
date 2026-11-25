package handlers

import (
	"github.com/gofiber/fiber/v2"
	"lume_api/internal/dto"
	"lume_api/internal/utils"
)

func (h *Handler) CreateReadingSession(c *fiber.Ctx) error {
	var req dto.ReadingSessionRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.Fail(c, fiber.StatusBadRequest, "invalidBody", "JSON inválido", err.Error())
	}
	ctx, cancel := h.withTimeout(c)
	defer cancel()
	session, err := h.Diary.CreateReadingSession(ctx, c.Locals("userId").(string), req)
	if err != nil {
		return mapError(c, err, "Falha ao criar sessão de leitura")
	}
	return utils.Success(c, fiber.StatusCreated, session, "Sessão de leitura criada")
}

func (h *Handler) GetReadingSession(c *fiber.Ctx) error {
	ctx, cancel := h.withTimeout(c)
	defer cancel()
	session, err := h.Diary.GetReadingSession(ctx, c.Locals("userId").(string), c.Params("id"))
	if err != nil {
		return mapError(c, err, "Sessão de leitura não encontrada")
	}
	return utils.Success(c, fiber.StatusOK, session, "Sessão carregada")
}

func (h *Handler) ListReadingSessions(c *fiber.Ctx) error {
	ctx, cancel := h.withTimeout(c)
	defer cancel()
	sessions, err := h.Diary.ListReadingSessions(ctx, c.Locals("userId").(string), parseLimit(c, 50))
	if err != nil {
		return mapError(c, err, "Falha ao listar sessões")
	}
	return utils.Success(c, fiber.StatusOK, sessions, "Sessões carregadas")
}

func (h *Handler) UpdateReadingSession(c *fiber.Ctx) error {
	var req dto.ReadingSessionRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.Fail(c, fiber.StatusBadRequest, "invalidBody", "JSON inválido", err.Error())
	}
	ctx, cancel := h.withTimeout(c)
	defer cancel()
	session, err := h.Diary.UpdateReadingSession(ctx, c.Locals("userId").(string), c.Params("id"), req)
	if err != nil {
		return mapError(c, err, "Falha ao atualizar sessão")
	}
	return utils.Success(c, fiber.StatusOK, session, "Sessão atualizada")
}

func (h *Handler) DeleteReadingSession(c *fiber.Ctx) error {
	ctx, cancel := h.withTimeout(c)
	defer cancel()
	if err := h.Diary.DeleteReadingSession(ctx, c.Locals("userId").(string), c.Params("id")); err != nil {
		return mapError(c, err, "Falha ao excluir sessão")
	}
	return utils.Success(c, fiber.StatusOK, nil, "Sessão excluída")
}

func (h *Handler) GetReadingHistory(c *fiber.Ctx) error {
	ctx, cancel := h.withTimeout(c)
	defer cancel()
	history, err := h.Diary.ReadingHistory(ctx, c.Locals("userId").(string), c.Params("id"))
	if err != nil {
		return mapError(c, err, "Falha ao carregar histórico")
	}
	return utils.Success(c, fiber.StatusOK, history, "Histórico carregado")
}

func (h *Handler) RecordReadingProgress(c *fiber.Ctx) error {
	var req dto.ReadingProgressRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.Fail(c, fiber.StatusBadRequest, "invalidBody", "JSON inválido", err.Error())
	}
	ctx, cancel := h.withTimeout(c)
	defer cancel()
	item, err := h.Diary.RecordProgress(ctx, c.Locals("userId").(string), c.Params("id"), req)
	if err != nil {
		return mapError(c, err, "Falha ao registrar progresso")
	}
	return utils.Success(c, fiber.StatusOK, item, "Progresso registrado ou atualizado")
}

func (h *Handler) ListReadingProgress(c *fiber.Ctx) error {
	ctx, cancel := h.withTimeout(c)
	defer cancel()
	progress, err := h.Diary.ReadingProgress(ctx, c.Locals("userId").(string), c.Params("id"))
	if err != nil {
		return mapError(c, err, "Falha ao carregar progresso")
	}
	return utils.Success(c, fiber.StatusOK, progress, "Progresso carregado")
}

func (h *Handler) SetReadingStatus(c *fiber.Ctx) error {
	var req dto.ReadingSessionStatusRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.Fail(c, fiber.StatusBadRequest, "invalidBody", "JSON inválido", err.Error())
	}
	ctx, cancel := h.withTimeout(c)
	defer cancel()
	session, err := h.Diary.SetReadingStatus(ctx, c.Locals("userId").(string), c.Params("id"), req.Status)
	if err != nil {
		return mapError(c, err, "Falha ao alterar status")
	}
	return utils.Success(c, fiber.StatusOK, session, "Status atualizado")
}

func (h *Handler) CreateNote(c *fiber.Ctx) error {
	var req dto.PrivateNoteRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.Fail(c, fiber.StatusBadRequest, "invalidBody", "JSON inválido", err.Error())
	}
	ctx, cancel := h.withTimeout(c)
	defer cancel()
	note, err := h.Diary.CreateNote(ctx, c.Locals("userId").(string), req)
	if err != nil {
		return mapError(c, err, "Falha ao criar anotação")
	}
	return utils.Success(c, fiber.StatusCreated, note, "Anotação privada criada")
}

func (h *Handler) GetNote(c *fiber.Ctx) error {
	ctx, cancel := h.withTimeout(c)
	defer cancel()
	userID, _ := c.Locals("userId").(string)
	note, err := h.Diary.GetNote(ctx, userID, c.Params("id"))
	if err != nil {
		return mapError(c, err, "Anotação não encontrada")
	}
	return utils.Success(c, fiber.StatusOK, note, "Anotação carregada")
}

func (h *Handler) ListNotes(c *fiber.Ctx) error {
	ctx, cancel := h.withTimeout(c)
	defer cancel()
	notes, err := h.Diary.ListNotes(ctx, c.Locals("userId").(string), c.Query("q"), parseLimit(c, 50))
	if err != nil {
		return mapError(c, err, "Falha ao listar anotações")
	}
	return utils.Success(c, fiber.StatusOK, notes, "Anotações carregadas")
}

func (h *Handler) UpdateNote(c *fiber.Ctx) error {
	var req dto.UpdatePrivateNoteRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.Fail(c, fiber.StatusBadRequest, "invalidBody", "JSON inválido", err.Error())
	}
	ctx, cancel := h.withTimeout(c)
	defer cancel()
	note, err := h.Diary.UpdateNote(ctx, c.Locals("userId").(string), c.Params("id"), req.Content)
	if err != nil {
		return mapError(c, err, "Falha ao atualizar anotação")
	}
	return utils.Success(c, fiber.StatusOK, note, "Anotação atualizada")
}

func (h *Handler) DeleteNote(c *fiber.Ctx) error {
	ctx, cancel := h.withTimeout(c)
	defer cancel()
	if err := h.Diary.DeleteNote(ctx, c.Locals("userId").(string), c.Params("id")); err != nil {
		return mapError(c, err, "Falha ao excluir anotação")
	}
	return utils.Success(c, fiber.StatusOK, nil, "Anotação excluída")
}

func (h *Handler) PublishNote(c *fiber.Ctx) error {
	ctx, cancel := h.withTimeout(c)
	defer cancel()
	note, err := h.Diary.PublishNote(ctx, c.Locals("userId").(string), c.Params("id"))
	if err != nil {
		return mapError(c, err, "Falha ao publicar anotação")
	}
	return utils.Success(c, fiber.StatusOK, note, "Anotação publicada")
}

func (h *Handler) CreateNoteTag(c *fiber.Ctx) error {
	var req dto.NoteTagRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.Fail(c, fiber.StatusBadRequest, "invalidBody", "JSON inválido", err.Error())
	}
	ctx, cancel := h.withTimeout(c)
	defer cancel()
	tag, err := h.Diary.CreateNoteTag(ctx, c.Locals("userId").(string), c.Params("id"), req.Tag)
	if err != nil {
		return mapError(c, err, "Falha ao adicionar etiqueta")
	}
	return utils.Success(c, fiber.StatusCreated, tag, "Etiqueta adicionada")
}

func (h *Handler) DeleteNoteTag(c *fiber.Ctx) error {
	ctx, cancel := h.withTimeout(c)
	defer cancel()
	if err := h.Diary.DeleteNoteTag(ctx, c.Locals("userId").(string), c.Params("id"), c.Params("tagId")); err != nil {
		return mapError(c, err, "Falha ao remover etiqueta")
	}
	return utils.Success(c, fiber.StatusOK, nil, "Etiqueta removida")
}

func (h *Handler) CreateNoteQuote(c *fiber.Ctx) error {
	var req dto.NoteQuoteRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.Fail(c, fiber.StatusBadRequest, "invalidBody", "JSON inválido", err.Error())
	}
	ctx, cancel := h.withTimeout(c)
	defer cancel()
	quote, err := h.Diary.CreateNoteQuote(ctx, c.Locals("userId").(string), c.Params("id"), req)
	if err != nil {
		return mapError(c, err, "Falha ao criar citação")
	}
	return utils.Success(c, fiber.StatusCreated, quote, "Citação criada")
}

func (h *Handler) UpdateNoteQuote(c *fiber.Ctx) error {
	var req dto.NoteQuoteRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.Fail(c, fiber.StatusBadRequest, "invalidBody", "JSON inválido", err.Error())
	}
	ctx, cancel := h.withTimeout(c)
	defer cancel()
	quote, err := h.Diary.UpdateNoteQuote(ctx, c.Locals("userId").(string), c.Params("id"), c.Params("quoteId"), req)
	if err != nil {
		return mapError(c, err, "Falha ao atualizar citação")
	}
	return utils.Success(c, fiber.StatusOK, quote, "Citação atualizada")
}

func (h *Handler) DeleteNoteQuote(c *fiber.Ctx) error {
	ctx, cancel := h.withTimeout(c)
	defer cancel()
	if err := h.Diary.DeleteNoteQuote(ctx, c.Locals("userId").(string), c.Params("id"), c.Params("quoteId")); err != nil {
		return mapError(c, err, "Falha ao excluir citação")
	}
	return utils.Success(c, fiber.StatusOK, nil, "Citação excluída")
}

func (h *Handler) CreateReviewDraft(c *fiber.Ctx) error {
	var req dto.ReviewDraftRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.Fail(c, fiber.StatusBadRequest, "invalidBody", "JSON inválido", err.Error())
	}
	ctx, cancel := h.withTimeout(c)
	defer cancel()
	review, err := h.Diary.CreateReviewDraft(ctx, c.Locals("userId").(string), req)
	if err != nil {
		return mapError(c, err, "Falha ao criar rascunho")
	}
	return utils.Success(c, fiber.StatusCreated, review, "Rascunho criado")
}

func (h *Handler) GetReviewDraft(c *fiber.Ctx) error {
	ctx, cancel := h.withTimeout(c)
	defer cancel()
	review, err := h.Diary.GetReviewDraft(ctx, c.Locals("userId").(string), c.Params("id"))
	if err != nil {
		return mapError(c, err, "Rascunho não encontrado")
	}
	return utils.Success(c, fiber.StatusOK, review, "Rascunho carregado")
}

func (h *Handler) GetPublishedReview(c *fiber.Ctx) error {
	ctx, cancel := h.withTimeout(c)
	defer cancel()
	review, err := h.Diary.GetPublishedReview(ctx, c.Params("id"))
	if err != nil {
		return mapError(c, err, "Resenha publicada não encontrada")
	}
	return utils.Success(c, fiber.StatusOK, review, "Resenha carregada")
}

func (h *Handler) UpdateReviewDraft(c *fiber.Ctx) error {
	var req dto.ReviewDraftRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.Fail(c, fiber.StatusBadRequest, "invalidBody", "JSON inválido", err.Error())
	}
	ctx, cancel := h.withTimeout(c)
	defer cancel()
	review, err := h.Diary.UpdateReviewDraft(ctx, c.Locals("userId").(string), c.Params("id"), req)
	if err != nil {
		return mapError(c, err, "Falha ao atualizar rascunho")
	}
	return utils.Success(c, fiber.StatusOK, review, "Rascunho atualizado")
}

func (h *Handler) DeleteReview(c *fiber.Ctx) error {
	ctx, cancel := h.withTimeout(c)
	defer cancel()
	if err := h.Diary.DeleteReview(ctx, c.Locals("userId").(string), c.Params("id")); err != nil {
		return mapError(c, err, "Falha ao excluir resenha")
	}
	return utils.Success(c, fiber.StatusOK, nil, "Resenha excluída")
}

func (h *Handler) PublishReview(c *fiber.Ctx) error {
	ctx, cancel := h.withTimeout(c)
	defer cancel()
	review, err := h.Diary.SetReviewPublished(ctx, c.Locals("userId").(string), c.Params("id"), true)
	if err != nil {
		return mapError(c, err, "Falha ao publicar resenha")
	}
	return utils.Success(c, fiber.StatusOK, review, "Resenha publicada")
}

func (h *Handler) UnpublishReview(c *fiber.Ctx) error {
	ctx, cancel := h.withTimeout(c)
	defer cancel()
	review, err := h.Diary.SetReviewPublished(ctx, c.Locals("userId").(string), c.Params("id"), false)
	if err != nil {
		return mapError(c, err, "Falha ao retirar publicação")
	}
	return utils.Success(c, fiber.StatusOK, review, "Resenha despublicada")
}

func (h *Handler) UpdateReview(c *fiber.Ctx) error {
	var req dto.UpdateReviewRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.Fail(c, fiber.StatusBadRequest, "invalidBody", "JSON inválido", err.Error())
	}
	ctx, cancel := h.withTimeout(c)
	defer cancel()
	review, err := h.Diary.UpdateReview(ctx, c.Locals("userId").(string), c.Params("id"), req)
	if err != nil {
		return mapError(c, err, "Falha ao atualizar resenha")
	}
	return utils.Success(c, fiber.StatusOK, review, "Resenha atualizada")
}

func (h *Handler) UpdateReviewPrivacy(c *fiber.Ctx) error {
	var req dto.ReviewPrivacyRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.Fail(c, fiber.StatusBadRequest, "invalidBody", "JSON inválido", err.Error())
	}
	ctx, cancel := h.withTimeout(c)
	defer cancel()
	review, err := h.Diary.UpdateReviewPrivacy(ctx, c.Locals("userId").(string), c.Params("id"), req)
	if err != nil {
		return mapError(c, err, "Falha ao atualizar configurações")
	}
	return utils.Success(c, fiber.StatusOK, review, "Configurações atualizadas")
}

func (h *Handler) GetReviewHistory(c *fiber.Ctx) error {
	ctx, cancel := h.withTimeout(c)
	defer cancel()
	history, err := h.Diary.ReviewHistory(ctx, c.Locals("userId").(string), c.Params("id"))
	if err != nil {
		return mapError(c, err, "Falha ao carregar histórico da resenha")
	}
	return utils.Success(c, fiber.StatusOK, history, "Histórico carregado")
}

func (h *Handler) CreateReviewFromNote(c *fiber.Ctx) error {
	ctx, cancel := h.withTimeout(c)
	defer cancel()
	review, err := h.Diary.CreateReviewFromNote(ctx, c.Locals("userId").(string), c.Params("id"))
	if err != nil {
		return mapError(c, err, "Falha ao criar resenha a partir da anotação")
	}
	return utils.Success(c, fiber.StatusCreated, review, "Rascunho de resenha criado")
}

func (h *Handler) CreateReviewFromLog(c *fiber.Ctx) error {
	ctx, cancel := h.withTimeout(c)
	defer cancel()
	review, err := h.Diary.CreateReviewFromLog(ctx, c.Locals("userId").(string), c.Params("id"))
	if err != nil {
		return mapError(c, err, "Falha ao criar resenha a partir da leitura")
	}
	return utils.Success(c, fiber.StatusCreated, review, "Rascunho de resenha criado")
}

func (h *Handler) GetReadingGoals(c *fiber.Ctx) error {
	ctx, cancel := h.withTimeout(c)
	defer cancel()
	goals, err := h.Diary.ReadingGoals(ctx, c.Locals("userId").(string), parseYear(c))
	if err != nil {
		return mapError(c, err, "Falha ao carregar metas")
	}
	return utils.Success(c, fiber.StatusOK, goals, "Metas carregadas")
}

func (h *Handler) SaveReadingGoal(c *fiber.Ctx) error {
	var req dto.ReadingGoalRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.Fail(c, fiber.StatusBadRequest, "invalidBody", "JSON inválido", err.Error())
	}
	ctx, cancel := h.withTimeout(c)
	defer cancel()
	goal, err := h.Diary.SaveReadingGoal(ctx, c.Locals("userId").(string), req)
	if err != nil {
		return mapError(c, err, "Falha ao salvar meta")
	}
	return utils.Success(c, fiber.StatusOK, goal, "Meta salva")
}

func (h *Handler) ImportReadingHistory(c *fiber.Ctx) error {
	var req dto.ReadingHistoryImportRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.Fail(c, fiber.StatusBadRequest, "invalidBody", "JSON inválido", err.Error())
	}
	ctx, cancel := h.withTimeout(c)
	defer cancel()
	job, err := h.Diary.ImportReadingHistory(ctx, c.Locals("userId").(string), req)
	if err != nil {
		return mapError(c, err, "Falha ao importar histórico")
	}
	return utils.Success(c, fiber.StatusCreated, job, "Histórico importado")
}

func (h *Handler) GetReadingHistoryImport(c *fiber.Ctx) error {
	ctx, cancel := h.withTimeout(c)
	defer cancel()
	job, err := h.Diary.ReadingHistoryImport(ctx, c.Locals("userId").(string), c.Params("jobId"))
	if err != nil {
		return mapError(c, err, "Importação não encontrada")
	}
	return utils.Success(c, fiber.StatusOK, job, "Importação carregada")
}
