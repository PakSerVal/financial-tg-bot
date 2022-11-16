package grpc

import (
	"context"

	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/clients/tg"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/logger"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/model"
	api "gitlab.ozon.dev/paksergey94/telegram-bot/pkg"
	"go.uber.org/zap"
)

const msgSendErr = "tg message send error"

type server struct {
	tgClient tg.Client

	api.UnimplementedReportServer
}

func New(tgClient tg.Client) api.ReportServer {
	return &server{
		tgClient: tgClient,
	}
}

func (s server) SendReport(ctx context.Context, in *api.SendReportIn) (*api.SendReportOut, error) {
	err := s.tgClient.SendMessage(model.MessageOut{
		Text: in.GetReport(),
	}, in.GetUserId())

	if err != nil {
		logger.Error(msgSendErr, zap.Error(err))

		return &api.SendReportOut{Ok: false}, err
	}

	return &api.SendReportOut{Ok: true}, nil
}
