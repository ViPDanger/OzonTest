package gRPC

import (
	context "context"
	"fmt"
	"time"

	"github.com/ViPDanger/OzonTest/internal/domain/entity"
	"github.com/ViPDanger/OzonTest/internal/interfaces/handlers"
	"github.com/ViPDanger/OzonTest/internal/interfaces/mapper"
	"github.com/ViPDanger/OzonTest/internal/usecase"
	"github.com/ViPDanger/OzonTest/proto"
)

func NewMockXMLDailyServer(uc usecase.ValCursUseCase, h handlers.ValCursHandler) proto.MockXMLDailyServer {
	return &mockXMLDailyServer{uc: uc, h: h}
}

type mockXMLDailyServer struct {
	proto.UnimplementedMockXMLDailyServer
	uc usecase.ValCursUseCase
	h  handlers.ValCursHandler
}

// Добавление mock данных
func (m *mockXMLDailyServer) AddMockData(ctx context.Context, req *proto.AddValCursRequest) (*proto.AddValCursResponse, error) {

	e := mapper.ValCursProtoToEntity(req.ValCurs)
	id, err := m.uc.Insert(ctx, &entity.XMLDailyMockResponse{Code: int(req.HttpStatus), ValuteCurs: e})
	if err != nil {
		return nil, fmt.Errorf("AddValCurs failed: %w", err)
	}
	return &proto.AddValCursResponse{Message: "added with id " + id}, nil
}

// Удаление mock данных
func (m *mockXMLDailyServer) DeleteMockData(ctx context.Context, req *proto.DeleteValCursRequest) (*proto.DeleteValCursResponse, error) {
	_, err := time.Parse("02.01.2006", req.GetDate())
	if err != nil {
		return nil, fmt.Errorf("DeleteValCurs failed: %w", err)
	}
	err = m.uc.DeleteByDate(ctx, req.GetDate())
	if err != nil {
		return nil, fmt.Errorf("DeleteValCurs failed: %w", err)
	}
	return &proto.DeleteValCursResponse{Message: "val deleted"}, nil
}
