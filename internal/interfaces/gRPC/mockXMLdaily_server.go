package gRPC

import (
	context "context"
	"fmt"
	"strings"
	"time"

	"github.com/ViPDanger/OzonTest/internal/interfaces/handlers"
	"github.com/ViPDanger/OzonTest/internal/interfaces/mapper"
	"github.com/ViPDanger/OzonTest/internal/usecase"
	"github.com/ViPDanger/OzonTest/proto"
	"google.golang.org/grpc/peer"
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
func (m *mockXMLDailyServer) AddValCurs(ctx context.Context, req *proto.AddValCursRequest) (*proto.AddValCursResponse, error) {

	e := mapper.ValCursProtoToEntity(req.ValCurs)
	e.CreatorID = getIP(ctx)
	id, err := m.uc.Insert(ctx, &e)
	if err != nil {
		return nil, fmt.Errorf("AddValCurs failed: %w", err)
	}
	return &proto.AddValCursResponse{Message: "added with id " + id}, nil
}

// Установка состояния
func (m *mockXMLDailyServer) SetState(ctx context.Context, req *proto.SetStateRequest) (*proto.SetStateResponse, error) {

	_, err := time.Parse("02.01.2006", req.GetDate())
	if err != nil {
		return nil, fmt.Errorf("AddValCurs failed: %w", err)
	}
	m.h.SetState(getIP(ctx), req.GetDate(), req.GetName())
	return &proto.SetStateResponse{Message: "state set"}, nil
}

// Получение текущего состояния mock-данных

func (m *mockXMLDailyServer) GetState(ctx context.Context, req *proto.GetStateRequest) (*proto.GetStateResponse, error) {
	date, name := m.h.GetState(getIP(ctx))
	return &proto.GetStateResponse{Message: "state set on date:" + date + ", name: " + name}, nil
}

// Очистка всех данных (reset)
func (m *mockXMLDailyServer) Reset(ctx context.Context, req *proto.ResetRequest) (*proto.ResetResponse, error) {
	err := m.uc.Reset(ctx)
	if err != nil {
		return nil, fmt.Errorf("Reset failed: %w", err)
	}
	return &proto.ResetResponse{Message: "reset complete"}, nil
}

func (m *mockXMLDailyServer) DeleteValCurs(ctx context.Context, req *proto.DeleteValCursRequest) (*proto.DeleteValCursResponse, error) {
	_, err := time.Parse("02.01.2006", req.GetDate())
	if err != nil {
		return nil, fmt.Errorf("DeleteValCurs failed: %w", err)
	}
	err = m.uc.DeleteByDateAndName(ctx, getIP(ctx), req.GetDate(), req.GetName())
	if err != nil {
		return nil, fmt.Errorf("DeleteValCurs failed: %w", err)
	}
	return &proto.DeleteValCursResponse{Message: "val deleted"}, nil
}

func getIP(ctx context.Context) string {
	p, ok := peer.FromContext(ctx)
	ip, _, _ := strings.Cut(p.Addr.String(), ":")
	if ok {
		fmt.Println("gRPC ClientAddr: ", ip)
	}
	return ip
}
