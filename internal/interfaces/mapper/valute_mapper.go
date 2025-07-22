package mapper

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/ViPDanger/OzonTest/internal/domain/entity"
	"github.com/ViPDanger/OzonTest/internal/interfaces/dto"
	"github.com/ViPDanger/OzonTest/internal/proto"
)

// ==========================Value mapper=============================
func ValuteToDTO(e entity.Valute) dto.ValuteDTO {
	return dto.ValuteDTO{
		ID:        e.ID,
		NumCode:   fmt.Sprintf("%03d", e.NumCode),
		CharCode:  e.CharCode,
		Nominal:   fmt.Sprintf("%d", e.Nominal),
		Name:      e.Name,
		Value:     strings.Replace(fmt.Sprintf("%5.4f", e.Value), ".", ",", 1),
		VunitRate: strings.Replace(fmt.Sprintf("%.4f", e.VunitRate), ".", ",", 1), // При необходимости можно добавить функцию по обработе показываемого значения.
	}
}

func ValuteToEntity(d dto.ValuteDTO) (r entity.Valute) {
	numCode, _ := strconv.Atoi(d.NumCode)
	nominal, _ := strconv.Atoi(d.Nominal)
	value, _ := strconv.ParseFloat(strings.Replace(d.Value, ",", ".", 1), 64)
	vunitRate, _ := strconv.ParseFloat(strings.Replace(d.VunitRate, ",", ".", 1), 64)

	return entity.Valute{
		ID:        d.ID,
		NumCode:   numCode,
		CharCode:  d.CharCode,
		Nominal:   nominal,
		Name:      d.Name,
		Value:     value,
		VunitRate: vunitRate,
	}
}

func ValuteProtoToEntity(p *proto.Valute) entity.Valute {
	return entity.Valute{
		ID:        p.GetId(),
		NumCode:   int(p.GetNumCode()),
		CharCode:  p.GetCharCode(),
		Nominal:   int(p.GetNominal()),
		Name:      p.GetName(),
		Value:     p.GetValue(),
		VunitRate: p.GetVunitRate(),
	}
}

// ========================ValCurs mapper============================
func ValCursToEntity(d dto.ValCursDTO) entity.ValuteCurs {
	valutes := make([]entity.Valute, len(d.Valutes))
	for i := range d.Valutes {
		valutes[i] = ValuteToEntity(d.Valutes[i])
	}
	return entity.ValuteCurs{
		Date:    d.Date,
		Name:    d.Name,
		Valutes: valutes,
	}
}

func ValCursToDTO(e entity.ValuteCurs) dto.ValCursDTO {
	valutes := make([]dto.ValuteDTO, 0, len(e.Valutes))
	for _, v := range e.Valutes {
		valutes = append(valutes, ValuteToDTO(v))
	}
	return dto.ValCursDTO{
		Date:    e.Date,
		Name:    e.Name,
		Valutes: valutes,
	}
}

func ValCursProtoToEntity(p *proto.ValCurs) entity.ValuteCurs {
	valutes := make([]entity.Valute, 0, len(p.GetValutes()))
	for _, v := range p.GetValutes() {
		valutes = append(valutes, ValuteProtoToEntity(v))
	}
	return entity.ValuteCurs{
		Date:    p.GetDate(),
		Name:    p.GetName(),
		Valutes: valutes,
	}
}
