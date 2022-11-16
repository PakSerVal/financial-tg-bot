package queue_message

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/kafka"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/logger"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/model"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/service/report"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/utils"
	api "gitlab.ozon.dev/paksergey94/telegram-bot/pkg"
	"go.uber.org/zap"
)

const (
	commandToday = "today"
	commandMonth = "month"
	commandYear  = "year"
)

type handler struct {
	reportClient  api.ReportClient
	reportService report.Service
}

func NewHandler(reportClient api.ReportClient, reportService report.Service) kafka.Handler {
	return &handler{
		reportClient:  reportClient,
		reportService: reportService,
	}
}

func (h *handler) HandleMessage(ctx context.Context, msg model.ReportMsg) {
	logger.Info("queue message received", zap.Int64("userId", msg.UserId), zap.String("period", msg.Period))

	rep, err := h.makeReport(ctx, msg.Period, msg.UserId)
	if err != nil {
		logger.Error("can not handle message from queue",
			zap.String("period", msg.Period),
			zap.Int64("userId", msg.UserId),
			zap.Error(err))
		return
	}

	logger.Info("sending report to tg...", zap.Int64("userId", msg.UserId), zap.String("report", rep))
	_, err = h.reportClient.SendReport(ctx, &api.SendReportIn{
		UserId: msg.UserId,
		Report: rep,
	})

	if err != nil {
		logger.Error("report service client error",
			zap.Error(err))
	}
}

func (h *handler) makeReport(ctx context.Context, period string, userId int64) (string, error) {
	now := time.Now()
	switch period {
	case commandToday:
		return h.reportService.MakeReport(
			ctx,
			userId,
			utils.BeginOfDay(now),
			"сегодня",
		)
	case commandMonth:
		return h.reportService.MakeReport(
			ctx,
			userId,
			utils.BeginOfMonth(now),
			"в текущем месяце",
		)
	case commandYear:
		return h.reportService.MakeReport(
			ctx,
			userId,
			utils.BeginOfYear(now),
			"в этом году",
		)
	}

	return "", errors.New("invalid period")
}
