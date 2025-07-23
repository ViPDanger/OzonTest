package app_test

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"os/signal"
	"syscall"
	"testing"

	"github.com/ViPDanger/OzonTest/internal/app"
	"github.com/ViPDanger/OzonTest/proto"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const host = "localhost:8085"
const mongoURI = "mongodb://127.0.0.1:27017"
const mongodbName = "xml_daily"
const mongoUser = "admin"
const mongoPassword = "admin"
const grpcHost = ":2025"

var testValCurs []proto.ValCurs

func setup(ctx context.Context) (*app.Application, proto.MockXMLDailyClient) {
	app, err := app.Run(ctx, host, mongoURI, mongodbName, mongoUser, mongoPassword, grpcHost)
	if err != nil {
		fmt.Printf("main(): error to run app.Run() %v", err.Error())
		return nil, nil
	}

	conn, err := grpc.NewClient(grpcHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	client := proto.NewMockXMLDailyClient(conn)
	go func() {
		defer conn.Close()
		<-ctx.Done()
	}()
	return app, client
}

// Тест должен вывести OK
func TestOK(t *testing.T) {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill, syscall.SIGALRM)
	defer cancel()
	app, grpcClient := setup(ctx)
	if app == nil || grpcClient == nil {
		t.Error("Nil Pointer gin server or gRPC client")
	}
	_, _ = grpcClient.Reset(ctx, &proto.ResetRequest{})                                                           // Очистить бд
	_, _ = grpcClient.AddValCurs(ctx, &proto.AddValCursRequest{ValCurs: &testValCurs[0]})                         // Добавить значение в бд
	_, _ = grpcClient.SetState(ctx, &proto.SetStateRequest{Date: testValCurs[0].Date, Name: testValCurs[0].Name}) // Установить HandlerState
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "http://"+host+"/", bytes.NewBuffer([]byte{}))
	app.GetGinServer().Handler.ServeHTTP(w, req)
	body, _ := io.ReadAll(w.Body)
	assert.Equal(t, w.Code, http.StatusOK, string(body))
}

// Тест должен вывести StatusInternalServerError
func TestInternalServerError(t *testing.T) {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill, syscall.SIGALRM)
	defer cancel()
	app, grpcClient := setup(ctx)
	if app == nil || grpcClient == nil {
		t.Error("Nil Pointer gin server or gRPC client")
	}
	_, _ = grpcClient.Reset(ctx, &proto.ResetRequest{}) // Очистить БД
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "http://"+host+"/", bytes.NewBuffer([]byte{}))
	app.GetGinServer().Handler.ServeHTTP(w, req)
	body, _ := io.ReadAll(w.Body)
	assert.Equal(t, w.Code, http.StatusInternalServerError, string(body))
}

func TestSwitchState(t *testing.T) {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill, syscall.SIGALRM)
	defer cancel()
	app, grpcClient := setup(ctx)
	if app == nil || grpcClient == nil {
		t.Error("Nil Pointer gin server or gRPC client")
	}
	_, _ = grpcClient.Reset(ctx, &proto.ResetRequest{})                                                           // Очистить бд
	_, _ = grpcClient.AddValCurs(ctx, &proto.AddValCursRequest{ValCurs: &testValCurs[0]})                         // Добавить значение в бд
	_, _ = grpcClient.AddValCurs(ctx, &proto.AddValCursRequest{ValCurs: &testValCurs[1]})                         // Добавить значение в бд
	_, _ = grpcClient.SetState(ctx, &proto.SetStateRequest{Date: testValCurs[0].Date, Name: testValCurs[0].Name}) // Установить HandlerState
	//работаем с состояние testValCurs[0]
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "http://"+host+"/", bytes.NewBuffer([]byte{}))
	app.GetGinServer().Handler.ServeHTTP(w, req)
	body, _ := io.ReadAll(w.Body)
	assert.Equal(t, w.Code, http.StatusOK, string(body))
	// изменить состояние на testValCurs[1]
	w = httptest.NewRecorder()
	_, _ = grpcClient.SetState(ctx, &proto.SetStateRequest{Date: testValCurs[1].Date, Name: testValCurs[1].Name}) // Установить HandlerState
	app.GetGinServer().Handler.ServeHTTP(w, req)
	body, _ = io.ReadAll(w.Body)
	assert.Equal(t, w.Code, http.StatusOK, string(body))

}

// переменные для тестирования
func init() {
	testValCurs = append(make([]proto.ValCurs, 0),
		proto.ValCurs{
			Date: "02.03.2002",
			Name: "Foreign Currency Market",
			Valutes: []*proto.Valute{
				{Id: "R01010", NumCode: 36, CharCode: "AUD", Nominal: 1, Name: "Австралийский доллар", Value: 16.0102, VunitRate: 16.0102},
				{Id: "R01035", NumCode: 826, CharCode: "GBP", Nominal: 1, Name: "Фунт стерлингов", Value: 43.8254, VunitRate: 43.8254},
				{Id: "R01090", NumCode: 974, CharCode: "BYR", Nominal: 1000, Name: "Белорусских рублей", Value: 18.4290, VunitRate: 0.018429},
				{Id: "R01215", NumCode: 208, CharCode: "DKK", Nominal: 10, Name: "Датских крон", Value: 36.1010, VunitRate: 3.6101},
				{Id: "R01235", NumCode: 840, CharCode: "USD", Nominal: 1, Name: "Доллар США", Value: 30.9436, VunitRate: 30.9436},
				{Id: "R01239", NumCode: 978, CharCode: "EUR", Nominal: 1, Name: "Евро", Value: 26.8343, VunitRate: 26.8343},
				{Id: "R01310", NumCode: 352, CharCode: "ISK", Nominal: 100, Name: "Исландских крон", Value: 30.7958, VunitRate: 0.307958},
				{Id: "R01335", NumCode: 398, CharCode: "KZT", Nominal: 100, Name: "Тенге", Value: 20.3393, VunitRate: 0.203393},
				{Id: "R01350", NumCode: 124, CharCode: "CAD", Nominal: 1, Name: "Канадский доллар", Value: 19.3240, VunitRate: 19.324},
				{Id: "R01535", NumCode: 578, CharCode: "NOK", Nominal: 10, Name: "Норвежских крон", Value: 34.7853, VunitRate: 3.47853},
				{Id: "R01589", NumCode: 960, CharCode: "XDR", Nominal: 1, Name: "СДР (специальные права заимствования)", Value: 38.4205, VunitRate: 38.4205},
				{Id: "R01625", NumCode: 702, CharCode: "SGD", Nominal: 1, Name: "Сингапурский доллар", Value: 16.8878, VunitRate: 16.8878},
				{Id: "R01700", NumCode: 792, CharCode: "TRL", Nominal: 1000000, Name: "Турецких лир", Value: 22.2616, VunitRate: 0.0000222616},
				{Id: "R01720", NumCode: 980, CharCode: "UAH", Nominal: 10, Name: "Гривен", Value: 58.1090, VunitRate: 5.8109},
				{Id: "R01770", NumCode: 752, CharCode: "SEK", Nominal: 10, Name: "Шведских крон", Value: 29.5924, VunitRate: 2.95924},
				{Id: "R01775", NumCode: 756, CharCode: "CHF", Nominal: 1, Name: "Швейцарский франк", Value: 18.1861, VunitRate: 18.1861},
				{Id: "R01820", NumCode: 392, CharCode: "JPY", Nominal: 100, Name: "Иен", Value: 23.1527, VunitRate: 0.231527},
			},
		},
		proto.ValCurs{
			Date: "30.04.2005",
			Name: "Foreign Currency Market",
			Valutes: []*proto.Valute{
				{Id: "R01010", NumCode: 36, CharCode: "AUD", Nominal: 1, Name: "Австралийский доллар", Value: 21.7098, VunitRate: 21.7098},
				{Id: "R01035", NumCode: 826, CharCode: "GBP", Nominal: 1, Name: "Фунт стерлингов", Value: 53.1873, VunitRate: 53.1873},
				{Id: "R01090", NumCode: 974, CharCode: "BYR", Nominal: 1000, Name: "Белорусских рублей", Value: 12.8987, VunitRate: 0.0128987},
				{Id: "R01215", NumCode: 208, CharCode: "DKK", Nominal: 10, Name: "Датских крон", Value: 48.3675, VunitRate: 4.83675},
				{Id: "R01235", NumCode: 840, CharCode: "USD", Nominal: 1, Name: "Доллар США", Value: 27.7726, VunitRate: 27.7726},
				{Id: "R01239", NumCode: 978, CharCode: "EUR", Nominal: 1, Name: "Евро", Value: 36.0072, VunitRate: 36.0072},
				{Id: "R01310", NumCode: 352, CharCode: "ISK", Nominal: 100, Name: "Исландских крон", Value: 43.8676, VunitRate: 0.438676},
				{Id: "R01335", NumCode: 398, CharCode: "KZT", Nominal: 100, Name: "Тенге", Value: 21.1391, VunitRate: 0.211391},
				{Id: "R01350", NumCode: 124, CharCode: "CAD", Nominal: 1, Name: "Канадский доллар", Value: 22.2199, VunitRate: 22.2199},
				{Id: "R01535", NumCode: 578, CharCode: "NOK", Nominal: 10, Name: "Норвежских крон", Value: 44.1993, VunitRate: 4.41993},
				{Id: "R01589", NumCode: 960, CharCode: "XDR", Nominal: 1, Name: "СДР (специальные права заимствования)", Value: 42.0349, VunitRate: 42.0349},
				{Id: "R01625", NumCode: 702, CharCode: "SGD", Nominal: 1, Name: "Сингапурский доллар", Value: 16.9407, VunitRate: 16.9407},
				{Id: "R01700J", NumCode: 949, CharCode: "TRY", Nominal: 1, Name: "Турецкая лира", Value: 19.8376, VunitRate: 19.8376},
				{Id: "R01720", NumCode: 980, CharCode: "UAH", Nominal: 10, Name: "Гривен", Value: 54.6974, VunitRate: 5.46974},
				{Id: "R01770", NumCode: 752, CharCode: "SEK", Nominal: 10, Name: "Шведских крон", Value: 39.2784, VunitRate: 3.92784},
				{Id: "R01775", NumCode: 756, CharCode: "CHF", Nominal: 1, Name: "Швейцарский франк", Value: 23.4269, VunitRate: 23.4269},
				{Id: "R01820", NumCode: 392, CharCode: "JPY", Nominal: 100, Name: "Иен", Value: 26.3873, VunitRate: 0.263873},
			},
		},
	)

}
