package api

import (
	"errors"
	"time"

	"github.com/Oabraham1/open-blogger/server/util"
	"github.com/gin-gonic/gin"
)

type RenewTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type RenewTokenResponse struct {
	AccessToken string    `json:"access_token"`
	ExpiresAt   time.Time `json:"expires_at"`
}

func (server *Server) RenewTokenRequest(ctx *gin.Context) {
	var req RenewTokenRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		server.BadRequestError(ctx)
		return
	}

	payload, err := server.Authenticator.VerifyToken(req.RefreshToken)
	if err != nil {
		server.UnauthorizedError(ctx)
		return
	}

	userSession, err := server.DataStore.GetSessionById(ctx, payload.ID)
	if err != nil {
		if errors.Is(err, util.ErrRecordNotFound) {
			server.NotFoundError(ctx)
			return
		}
		server.InternalServerError(ctx)
		return
	}

	if userSession.RefreshToken != req.RefreshToken {
		server.UnauthorizedError(ctx)
		return
	}

	if userSession.IsBlocked {
		server.ForbiddenError(ctx)
		return
	}

	if userSession.Username != payload.Username {
		server.UnauthorizedError(ctx)
		return
	}

	if userSession.ExpiresAt.Before(time.Now()) {
		server.UnauthorizedError(ctx)
		return
	}

	accessToken, accessPayload, err := server.Authenticator.CreateToken(userSession.Username, server.Configurations.AccessTokenDuration)
	if err != nil {
		server.InternalServerError(ctx)
		return
	}

	rsp := RenewTokenResponse{
		AccessToken: accessToken,
		ExpiresAt:   accessPayload.ExpiredAt,
	}
	server.ReturnOK(ctx, rsp)
}
