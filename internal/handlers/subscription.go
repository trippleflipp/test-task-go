package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"github.com/trippleflipp/test-task-go/internal/models"
)

type SubscriptionHandler struct {
	DB     *sqlx.DB
	Logger *logrus.Logger
}

func NewSubscriptionHandler(db *sqlx.DB, logger *logrus.Logger) *SubscriptionHandler {
	return &SubscriptionHandler{DB: db, Logger: logger}
}

// Create godoc
// @Summary      Создать подписку
// @Description  Добавляет новую запись о подписке в базу данных
// @Tags         subscriptions
// @Accept       json
// @Produce      json
// @Param        input body models.Subscription true "Данные подписки"
// @Success      201 {object} models.Subscription
// @Failure      400 {object} map[string]string
// @Router       /api/subscriptions [post]
func (h *SubscriptionHandler) Create(c *gin.Context) {
	var sub models.Subscription
	if err := c.ShouldBindJSON(&sub); err != nil {
		h.Logger.WithError(err).Warn("Failed to bind JSON create")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	query := `INSERT INTO subscriptions (service_name, price, user_id, start_date, end_date)
		VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at`

	err := h.DB.QueryRow(query, sub.ServiceName, sub.Price, sub.UserID, sub.StartDate, sub.EndDate).
		Scan(&sub.ID, &sub.CreatedAt)

	if err != nil {
		h.Logger.WithError(err).Error("Failed to insert subscription")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	h.Logger.WithFields(logrus.Fields{
		"id":      sub.ID,
		"user_id": sub.UserID,
	}).Info("Subscription created")
	c.JSON(http.StatusCreated, sub)
}

// Update godoc
// @Summary      Обновить подписку
// @Description  Обновляет данные существующей подписки по её ID
// @Tags         subscriptions
// @Accept       json
// @Produce      json
// @Param        id    path      int                  true  "ID подписки"
// @Param        input body      models.Subscription  true  "Новые данные подписки"
// @Success      200   {object}  map[string]string
// @Failure      400   {object}  map[string]string
// @Failure      404   {object}  map[string]string
// @Router       /api/subscriptions/{id} [put]
func (h *SubscriptionHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var sub models.Subscription

	if err := c.ShouldBindJSON(&sub); err != nil {
		h.Logger.WithError(err).Warn("Failed to bind JSON update")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	query := `UPDATE subscriptions
		SET service_name = $1, price = $2, user_id = $3, start_date = $4, end_date = $5
		WHERE id = $6`

	res, err := h.DB.Exec(query, sub.ServiceName, sub.Price, sub.UserID, sub.StartDate, sub.EndDate, id)
	if err != nil {
		h.Logger.WithError(err).Error("Failed to update subscription")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Subscription not found"})
		return
	}

	h.Logger.Infof("Subscription %s updated", id)
	c.JSON(http.StatusOK, gin.H{"status": "OK"})
}

// List godoc
// @Summary      Список всех подписок
// @Description  Возвращает список всех оформленных подписок
// @Tags         subscriptions
// @Produce      json
// @Success      200 {array} models.Subscription
// @Router       /api/subscriptions [get]
func (h *SubscriptionHandler) List(c *gin.Context) {
	var subs []models.Subscription
	err := h.DB.Select(&subs, "SELECT * FROM subscriptions ORDER BY created_at DESC")
	if err != nil {
		h.Logger.WithError(err).Error("Failed to fetch subscriptions")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	c.JSON(http.StatusOK, subs)
}

// GetByID godoc
// @Summary      Получить подписку по ID
// @Description  Возвращает полную информацию об одной подписке
// @Tags         subscriptions
// @Produce      json
// @Param        id   path      int  true  "ID подписки"
// @Success      200  {object}  models.Subscription
// @Failure      404  {object}  map[string]string
// @Router       /api/subscriptions/{id} [get]
func (h *SubscriptionHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	var sub models.Subscription

	err := h.DB.Get(&sub, "SELECT * FROM subscriptions WHERE id = $1", id)
	if err != nil {
		h.Logger.WithError(err).Warnf("Subscription %s not found", id)
		c.JSON(http.StatusNotFound, gin.H{"error": "Subscription not found"})
		return
	}

	c.JSON(http.StatusOK, sub)
}

// GetTotal godoc
// @Summary      Рассчитать общую стоимость
// @Description  Возвращает сумму цен подписок с фильтрацией по пользователю, сервису и датам
// @Tags         subscriptions
// @Produce      json
// @Param        user_id      query    string  false  "UUID пользователя"
// @Param        service_name query    string  false  "Название сервиса"
// @Param        from        query    string  false  "Дата от (MM-YYYY)"
// @Param        to          query    string  false  "Дата до (MM-YYYY)"
// @Success      200 {object} map[string]int
// @Router       /api/subscriptions/total [get]
func (h *SubscriptionHandler) GetTotal(c *gin.Context) {
	userID := c.Query("user_id")
	serviceName := c.Query("service_name")
	from := c.Query("from")
	to := c.Query("to")

	query := `SELECT COALESCE(SUM(price), 0) FROM subscriptions WHERE 1=1`
	var args []any
	argCount := 1

	if userID != "" {
		query += fmt.Sprintf(" AND user_id = $%d", argCount)
		args = append(args, userID)
		argCount++
	}
	if serviceName != "" {
		query += fmt.Sprintf(" AND service_name = $%d", argCount)
		args = append(args, serviceName)
	}
	if from != "" {
		t, err := time.Parse("01-2006", from)
		if err == nil {
			query += fmt.Sprintf(" AND start_date >= $%d", argCount)
			args = append(args, t.Format("2006-01-02"))
			argCount++
		}
	}
	if to != "" {
		t, err := time.Parse("01-2006", to)
		if err == nil {
			query += fmt.Sprintf(" AND start_date <= $%d", argCount)
			args = append(args, t.Format("2006-01-02"))
			argCount++
		}
	}

	var total int
	if err := h.DB.Get(&total, query, args...); err != nil {
		h.Logger.WithError(err).Error("Failed to calculate total cost")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	h.Logger.WithFields(logrus.Fields{
		"user_id": userID,
		"total":   total,
	}).Info("Total cost calculated")
	c.JSON(http.StatusOK, gin.H{"total_cost": total})
}

// Delete godoc
// @Summary      Удалить подписку
// @Description  Удаляет запись о подписке из базы данных по ID
// @Tags         subscriptions
// @Param        id   path      int  true  "ID подписки"
// @Success      204  "No Content"
// @Failure      404  {object}  map[string]string
// @Router       /api/subscriptions/{id} [delete]
func (h *SubscriptionHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	res, err := h.DB.Exec("DELETE FROM subscriptions WHERE id = $1", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
		return
	}

	h.Logger.Infof("Subscription %s deleted", id)
	c.Status(http.StatusNoContent)
}
