package handlers

import (
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/levinOo/go-crudl-task/internal/domain"

	"github.com/gin-gonic/gin"
)

// Структура создания подписки
type createSubInput struct {
	ServiceName string `json:"service_name" binding:"required"`
	Price       int64  `json:"price" binding:"required"`
	UserID      string `json:"user_id" binding:"required"`
	StartDate   string `json:"start_date" binding:"required"`
}

// Структура обновления подписки
type updateSubInput struct {
	Price   *int64  `json:"price"`
	EndDate *string `json:"end_date"`
}

// Парсинг даты
func parseDate(dateStr string) (time.Time, error) {
	t, err := time.Parse("01-2006", dateStr)
	if err != nil {
		return time.Time{}, err
	}
	return t, nil
}

// CreateSubscription - создание подписки
//
//	@Summary		Создание подписки
//	@Description	Создать новую подписку
//	@Tags			subscriptions
//	@Accept			json
//	@Produce		json
//	@Param			body	body		createSubInput		true	"Данные подписки"
//	@Success		201		{object}	map[string]string	"ID подписки"
//	@Failure		400		{object}	domain.ErrorResponse		"Неверное тело запроса"
//	@Failure		500		{object}	domain.ErrorResponse		"Внутренняя ошибка сервера"
//	@Router			/subscriptions [post]
func (h *Handler) createSubscription(c *gin.Context) {
	var input createSubInput

	// Читаем JSON
	if err := c.ShouldBindJSON(&input); err != nil {
		h.log.Warn("ошибка при чтении JSON", slog.String("error", err.Error()))
		newErrorResponse(c, http.StatusBadRequest, "Неверное тело запроса")
		return
	}

	// Парсим дату начала
	startDate, err := parseDate(input.StartDate)
	if err != nil {
		h.log.Warn("ошибка парсинга даты", slog.String("date", input.StartDate), slog.String("error", err.Error()))
		newErrorResponse(c, http.StatusBadRequest, "Неверный формат даты начала. Ожидается MM-YYYY")
		return
	}

	sub := domain.Subscription{
		ServiceName: input.ServiceName,
		Price:       int(input.Price),
		UserID:      input.UserID,
		StartDate:   startDate,
		EndDate:     nil,
	}

	// Вызываем слой сервис
	id, err := h.services.Create(c.Request.Context(), sub)
	if err != nil {
		h.log.Error("ошибка при создании подписки", slog.String("error", err.Error()))
		newErrorResponse(c, http.StatusInternalServerError, "Внутренняя ошибка сервера")
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": id})
}

// GetSubscription - получение подписки по ID
//
//	@Summary		Получение подписки
//	@Description	Получить информацию о подписке по её ID
//	@Tags			subscriptions
//	@Produce		json
//	@Param			id	path		string	true	"ID подписки"
//	@Success		200	{object}	domain.Subscription
//	@Failure		404	{object}	domain.ErrorResponse	"Подписка не найдена"
//	@Failure		500	{object}	domain.ErrorResponse	"Внутренняя ошибка сервера"
//	@Router			/subscriptions/{id} [get]
func (h *Handler) getSubscription(c *gin.Context) {
	// Достаем id из URL
	id := c.Param("id")
	if id == "" {
		h.log.Error("ID подписки не может быть пустым", slog.String("id", id))
		newErrorResponse(c, http.StatusBadRequest, "ID подписки не может быть пустым")
		return
	}

	// Вызываем слой сервис
	sub, err := h.services.Get(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, domain.ErrSubscriptionNotFound) {
			h.log.Error("подписка не найдена", slog.String("id", id), slog.String("error", err.Error()))
			newErrorResponse(c, http.StatusNotFound, "Подписка не найдена")
			return
		}

		h.log.Error("ошибка при получении подписки", slog.String("id", id), slog.String("error", err.Error()))
		newErrorResponse(c, http.StatusInternalServerError, "Внутренняя ошибка сервера")
		return
	}

	c.JSON(http.StatusOK, sub)
}

// UpdateSubscription - обновление (цена, дата окончания)
//
//	@Summary		Обновление данных подписки
//	@Description	Обновить цену или дату окончания подписки
//	@Tags			subscriptions
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string				true	"ID подписки"
//	@Param			body	body		updateSubInput		true	"Данные для обновления"
//	@Success		200		{object}	map[string]string	"Статус и сообщение"
//	@Failure		400		{object}	domain.ErrorResponse	"Неверные данные"
//	@Failure		404		{object}	domain.ErrorResponse	"Подписка не найдена"
//	@Failure		500		{object}	domain.ErrorResponse	"Внутренняя ошибка сервера"
//	@Router			/subscriptions/{id} [patch]
func (h *Handler) updateSubscription(c *gin.Context) {
	// Достаем id из URL
	id := c.Param("id")
	if id == "" {
		h.log.Error("ID подписки не может быть пустым", slog.String("id", id))
		newErrorResponse(c, http.StatusBadRequest, "ID подписки не может быть пустым")
		return
	}

	var input updateSubInput

	// Читаем JSON
	if err := c.ShouldBindJSON(&input); err != nil {
		h.log.Warn("ошибка при чтении JSON", slog.String("error", err.Error()))
		newErrorResponse(c, http.StatusBadRequest, "Неверное тело запроса")
		return
	}

	var endDate *time.Time
	// Если прислали дату окончания — парсим её
	if input.EndDate != nil {
		t, err := parseDate(*input.EndDate)
		if err != nil {
			h.log.Error("ошибка парсинга даты", slog.String("date", *input.EndDate), slog.String("error", err.Error()))
			newErrorResponse(c, http.StatusBadRequest, "Неверный формат даты окончания. Ожидается MM-YYYY")
			return
		}
		endDate = &t
	}

	updateData := domain.UpdateSubscriptionInput{
		Price:   input.Price,
		EndDate: endDate,
	}

	// Вызываем слой сервис
	err := h.services.Update(c.Request.Context(), id, updateData)
	if err != nil {
		if errors.Is(err, domain.ErrSubscriptionNotFound) {
			h.log.Error("подписка не найдена", slog.String("id", id), slog.String("error", err.Error()))
			newErrorResponse(c, http.StatusNotFound, "Подписка не найдена")
			return
		}
		if errors.Is(err, domain.ErrInvalidPeriod) {
			h.log.Error("ошибка при обновлении подписки", slog.String("id", id), slog.String("error", err.Error()))
			newErrorResponse(c, http.StatusBadRequest, "Дата окончания не может быть раньше даты начала")
			return
		}

		h.log.Error("ошибка при обновлении подписки", slog.String("id", id), slog.String("error", err.Error()))
		newErrorResponse(c, http.StatusInternalServerError, "Внутренняя ошибка сервера")
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok", "message": "Подписка обновлена"})
}

// DeleteSubscription - удаление
//
//	@Summary		Удаление подписки
//	@Description	Удалить подписку по ID
//	@Tags			subscriptions
//	@Produce		json
//	@Param			id	path	string	true	"ID подписки"
//	@Success		204	"Подписка успешно удалена"
//	@Failure		404	{object}	domain.ErrorResponse	"Подписка не найдена"
//	@Failure		500	{object}	domain.ErrorResponse	"Внутренняя ошибка сервера"
//	@Router			/subscriptions/{id} [delete]
func (h *Handler) deleteSubscription(c *gin.Context) {
	// Достаем id из URL
	id := c.Param("id")
	if id == "" {
		h.log.Error("ID подписки не может быть пустым", slog.String("id", id))
		newErrorResponse(c, http.StatusBadRequest, "ID подписки не может быть пустым")
		return
	}

	// Вызываем слой сервис
	err := h.services.Delete(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, domain.ErrSubscriptionNotFound) {
			h.log.Error("подписка не найдена", slog.String("id", id), slog.String("error", err.Error()))
			newErrorResponse(c, http.StatusNotFound, "Подписка не найдена")
			return
		}

		h.log.Error("ошибка при удалении подписки", slog.String("id", id), slog.String("error", err.Error()))
		newErrorResponse(c, http.StatusInternalServerError, "Внутренняя ошибка сервера")
		return
	}

	// Успешное удаление (204 No Content)
	c.Status(http.StatusNoContent)
}

// GetList - получение списка (с фильтрацией по user_id)
//
//	@Summary		Получение списка подписок
//	@Description	Получить список всех подписок пользователя
//	@Tags			subscriptions
//	@Produce		json
//	@Param			user_id	query		string	true	"UUID пользователя"
//	@Success		200		{array}		domain.Subscription
//	@Failure		500		{object}	domain.ErrorResponse	"Внутренняя ошибка сервера"
//	@Router			/subscriptions [get]
func (h *Handler) getList(c *gin.Context) {
	// Достаем id из URL
	userID := c.Query("user_id")
	if userID == "" {
		h.log.Error("ID пользователя не может быть пустым", slog.String("user_id", userID))
		newErrorResponse(c, http.StatusBadRequest, "ID пользователя не может быть пустым")
		return
	}

	// Вызываем слой сервис
	subs, err := h.services.List(c.Request.Context(), userID)
	if err != nil {
		h.log.Error("ошибка при получении списка", slog.String("error", err.Error()))
		newErrorResponse(c, http.StatusInternalServerError, "Внутренняя ошибка сервера")
		return
	}

	// Если список пустой, возвращаем пустой массив []
	if subs == nil {
		subs = []domain.Subscription{}
	}

	c.JSON(http.StatusOK, subs)
}

// GetTotalCost - подсчет суммарной стоимости подписок за выбранный период с фильтрацией
//
//	@Summary		Подсчитать суммарную стоимость подписок
//	@Description	Получить суммарную стоимость подписок за выбранный период с фильтрацией по user_id и названию подписки
//	@Tags			subscriptions
//	@Produce		json
//	@Param			user_id			query		string			true	"UUID пользователя"
//	@Param			service_name	query		string			false	"Название подписки"
//	@Param			start_date		query		string			true	"Начальная дата (формат MM-YYYY)"
//	@Param			end_date		query		string			true	"Конечная дата (формат MM-YYYY)"
//	@Success		200				{object}	map[string]int	"Суммарная стоимость"
//	@Failure		400				{object}	domain.ErrorResponse	"Неверные параметры"
//	@Failure		500				{object}	domain.ErrorResponse	"Внутренняя ошибка сервера"
//	@Router			/subscriptions/total-cost [get]
func (h *Handler) getTotalCost(c *gin.Context) {
	// Достаем значеня из query params
	userID := c.Query("user_id")
	serviceName := c.Query("service_name")
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	// Проверяем обязательные параметры
	if userID == "" {
		h.log.Warn("user_id не указан")
		newErrorResponse(c, http.StatusBadRequest, "user_id обязателен")
		return
	}

	if startDateStr == "" || endDateStr == "" {
		h.log.Warn("даты не указаны")
		newErrorResponse(c, http.StatusBadRequest, "start_date и end_date обязательны")
		return
	}

	// Парсим даты
	startDate, err := parseDate(startDateStr)
	if err != nil {
		h.log.Warn("ошибка парсинга start_date", slog.String("date", startDateStr), slog.String("error", err.Error()))
		newErrorResponse(c, http.StatusBadRequest, "Неверный формат start_date. Ожидается MM-YYYY")
		return
	}

	endDate, err := parseDate(endDateStr)
	if err != nil {
		h.log.Warn("ошибка парсинга end_date", slog.String("date", endDateStr), slog.String("error", err.Error()))
		newErrorResponse(c, http.StatusBadRequest, "Неверный формат end_date. Ожидается MM-YYYY")
		return
	}

	if startDate.After(endDate) {
		h.log.Warn("start_date после end_date")
		newErrorResponse(c, http.StatusBadRequest, "start_date не может быть после end_date")
		return
	}

	// Вызываем слой сервис
	total, err := h.services.GetTotalCost(c.Request.Context(), userID, serviceName, startDate, endDate)
	if err != nil {
		h.log.Error("ошибка при подсчете стоимости", slog.String("error", err.Error()))
		newErrorResponse(c, http.StatusInternalServerError, "Внутренняя ошибка сервера")
		return
	}

	c.JSON(http.StatusOK, gin.H{"total_cost": total})
}
