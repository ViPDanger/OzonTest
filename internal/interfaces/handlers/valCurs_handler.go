package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/ViPDanger/OzonTest/internal/interfaces/mapper"
	"github.com/ViPDanger/OzonTest/internal/usecase"
	"github.com/gin-gonic/gin"
)

type ValCursHandler interface {
	GetByDateAndName(c *gin.Context)
	SetState(date string, name string)
	GetState() (date string, name string)
}

func NewValCursHandler(uc usecase.ValCursUseCase) ValCursHandler {
	return &valCursHandler{uc: uc, HandlerState: &HandlerState{date: time.Now().Format("02.01.2006")}}
}

type valCursHandler struct {
	uc usecase.ValCursUseCase
	*HandlerState
}

// GetByDateAndName godoc
// @Summary      Получить курсы валют на дату
// @Description  Возвращает курсы валют в XML-формате на указанную дату. Если заголовок date_req не указан, используется текущая дата.
// @Accept       xml
// @Produce      xml
// @Param        date_req header string false "Дата в формате DD.MM.YYYY"
// @Success      200 {object} dto.ValCursDTO
// @Failure      500 {object} map[string]string
// @Router       /currencies/by-date [get]
func (h *valCursHandler) GetByDateAndName(c *gin.Context) {
	// проверка Usecase
	if h.uc == nil {
		err := errors.New("valCursHandler.GetByDateAndName(): nil pointer Usecase")
		_ = c.Error(err)
		c.XML(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// берём date, name из
	date, name := h.GetState()
	if date == "" {
		date = time.Now().Format("02.01.2006")
	}
	_, err := time.Parse("02.01.2006", date)
	if err != nil {
		err = fmt.Errorf("valCursHandler.GetByDateAndName()/%w", err)
		_ = c.Error(err)
		c.XML(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// обращение к Usecase
	object, err := h.uc.GetByDateAndName(c.Request.Context(), date, name)
	if err != nil {
		err = fmt.Errorf("valCursHandler.GetByDateAndName()/%w", err)
		_ = c.Error(err)
		c.XML(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if object == nil {
		// Можно было и http.StatusNoContent с message в body, но условия задачи есть условия задачи.
		c.XML(http.StatusInternalServerError, gin.H{"error": "No Content"})
		return
	}
	// ВЫВОД
	c.XML(http.StatusOK, mapper.ValCursToDTO(*object))
}

type HandlerState struct {
	sync.Mutex
	date string
	name string
}

func (h *HandlerState) SetState(date string, name string) {
	h.Mutex.Lock()
	defer h.Mutex.Unlock()
	_, err := time.Parse("02.01.2006", date)
	if err != nil {
		return
	}
	h.date = date
	h.name = name
}

func (h *HandlerState) GetState() (date string, name string) {
	h.Mutex.Lock()
	defer h.Mutex.Unlock()
	return h.date, h.name
}
