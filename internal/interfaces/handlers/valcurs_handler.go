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
	SetState(id string, date string, name string)
	GetState(id string) (date string, name string)
}

func NewValCursHandler(uc usecase.ValCursUseCase) ValCursHandler {
	return &valCursHandler{uc: uc, HandlerState: &HandlerState{}}
}

type valCursHandler struct {
	uc usecase.ValCursUseCase
	*HandlerState
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
	// берём date, name из состояния handler от ClientIP
	date, name := h.GetState(c.ClientIP())
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
	object, err := h.uc.GetByDateAndName(c.Request.Context(), c.ClientIP(), date, name)
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
	m    sync.Map
	date string //	default(time.Now() date)
	name string //	default("Foreign Currency Market")
}

// Функция установки состояния handler
func (h *HandlerState) SetState(id string, date string, name string) {
	h.Mutex.Lock()
	defer h.Mutex.Unlock()
	h.m.Store(id, []string{date, name})
	_, err := time.Parse("02.01.2006", date)
	if err != nil {
		return
	}
	h.date = date
	h.name = name
}

// Функция просмотра состояния handler
func (h *HandlerState) GetState(id string) (date string, name string) {
	h.Mutex.Lock()
	defer h.Mutex.Unlock()
	if str, ok := h.m.Load(id); ok {
		return str.([]string)[0], str.([]string)[1]
	}
	return time.Now().Format("02.01.2006"), "Foreign Currency Market"
}
