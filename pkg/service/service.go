package service

import (
	"github.com/rs/zerolog"
	"scbunn.org/tmp/gps-tracking-service/pkg/models"
)

type Service struct {
	address   string
	telemetry models.TelemetryReaderWriterChecker
	logger    *zerolog.Logger
}

func New(addr string, telemetry models.TelemetryReaderWriterChecker, log *zerolog.Logger) *Service {
	return &Service{
		address:   addr,
		telemetry: telemetry,
		logger:    log,
	}
}
