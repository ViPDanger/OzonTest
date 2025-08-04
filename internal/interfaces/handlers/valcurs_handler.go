package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/ViPDanger/OzonTest/internal/interfaces/mapper"
	"github.com/ViPDanger/OzonTest/internal/usecase"
	"github.com/gin-gonic/gin"
)

type ValCursHandler interface {
	GetByDateAndName(c *gin.Context)
}

func NewValCursHandler(uc usecase.ValCursUseCase) ValCursHandler {
	return &valCursHandler{uc: uc}
}

type valCursHandler struct {
	uc usecase.ValCursUseCase
}

// GetByDateAndName godoc
// @Summary      Получить валютные курсы по дате и имени источника
// @Description  Возвращает XML-список валют на заданную дату и по заданному имени источника. Дата и имя берутся из внутреннего состояния, которое можно изменить через gRPC SetState.
// @Tags         currency
// @Produce      xml
// @Success      200 {object} dto.ValCursDTO
// @Failure      500 {object} map[string]string "Ошибка запроса"
// @Router       /curs [get]
func (h *valCursHandler) GetByDateAndName(c *gin.Context) {
	// проверка Usecase
	if h.uc == nil {
		err := errors.New("valCursHandler.GetByDateAndName(): nil pointer Usecase")
		_ = c.Error(err)
		c.XML(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// валидируем date_req
	dateReq := c.Query("date_req")
	t, err := time.Parse("02/01/2006", dateReq)
	if err != nil {
		err = fmt.Errorf("valCursHandler.GetByDateAndName()/%w", err)
		_ = c.Error(err)
		c.XML(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// обращение к Usecase
	object, err := h.uc.GetByDate(c.Request.Context(), t.Format("02.01.2006"))
	if err != nil {
		err = fmt.Errorf("valCursHandler.GetByDateAndName()/%w", err)
		_ = c.Error(err)
		c.XML(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if object == nil {
		// Можно было и http.StatusNoContent с message в body, но условия задачи есть условия задачи.
		_ = c.Error(errors.New("No content on " + dateReq))
		c.Status(http.StatusInternalServerError)
		return
	}
	// ВЫВОД
	c.XML(object.Code, mapper.ValCursToDTO(object.ValuteCurs))
}
