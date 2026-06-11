// Package http contains the HTTP transport adapters exposed by api-gateway.
//
// @title OFM API Gateway
// @version 1.0
// @description Public HTTP entrypoint for the OFM microservices MVP.
// @BasePath /v1
// @schemes http
package http

import gateway "api-gateway/internal/domain"

var (
	_ = gateway.SignUpRequest{}
	_ = gateway.SignUpResult{}
)

// SignUpBadRequestResponse describes validation failures returned by the signup
// endpoint.
type SignUpBadRequestResponse struct {
	Error string `json:"error" example:"invalid email"`
}

// SignUpConflictResponse describes duplicate or already-active registration
// attempts returned by the signup endpoint.
type SignUpConflictResponse struct {
	Error         string `json:"error" example:"registration conflict"`
	State         string `json:"state" example:"found_completed"`
	UsernameTaken bool   `json:"username_taken" example:"true"`
	EmailTaken    bool   `json:"email_taken" example:"false"`
}

// SignUpInternalErrorResponse describes unexpected failures returned by the
// signup endpoint.
type SignUpInternalErrorResponse struct {
	Error string `json:"error" example:"internal server error"`
}
