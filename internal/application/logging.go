package service

import (
	gateway "api-gateway/internal/domain"
	"github.com/ofm-microservices/ofm-common/pkg/logging"
)

// LogOperationFailure keeps expected client-visible domain outcomes at warn
// while preserving error severity for infrastructure and unknown failures.
func LogOperationFailure(log logging.Logger, message string, err error, fields ...logging.Field) {
	fields = append(fields, logging.Err(err))
	if gateway.IsBusinessError(err) {
		log.Warn(message, fields...)
		return
	}
	log.Error(message, fields...)
}
