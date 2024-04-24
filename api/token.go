package api

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type renewAccessTokenResponse struct {
	AccessToken string `json:"access_token"`
	AccessTokenExpiresAt time.Time `json:"access_token_expires_at"`
}

type renewAccessTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

func (server *Server) refreshAccessToken(ctx *gin.Context){
	/*
	1- get refresh token from request
	2- get token payload
	3- get the session and make some validations on it
	4- generate a nes access token

	*/
	var req renewAccessTokenRequest
	if err := ctx.ShouldBindJSON(&req); err!=nil {
		ctx.JSON(http.StatusBadRequest,errorResponse(err))
		return
	}

	refreshPayload,err := server.tokenMaker.VerifyToken(req.RefreshToken)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized,errorResponse(err))
		return
	}

	session,err := server.store.GetSession(ctx,refreshPayload.ID)
	if err != nil {
		if err == sql.ErrNoRows{
			ctx.JSON(http.StatusNotFound,errorResponse(err))
			return
		}
		ctx.JSON(http.StatusUnauthorized,errorResponse(err))
		return
	}

	if session.IsBlocked{
		err:= fmt.Errorf("blocked session")
		ctx.JSON(http.StatusUnauthorized,errorResponse(err))
		return
	}

	if session.Username != refreshPayload.Username{
		err:= fmt.Errorf("incorrect session user")
		ctx.JSON(http.StatusUnauthorized,errorResponse(err))
		return
	}

	if session.RefreshToken != req.RefreshToken{
		err:= fmt.Errorf("mismatch session token")
		ctx.JSON(http.StatusUnauthorized,errorResponse(err))
		return
	}

	if time.Now().After(session.ExpiresAt){
		err:= fmt.Errorf("expired session")
		ctx.JSON(http.StatusUnauthorized,errorResponse(err))
		return
	
	}
	accessToken,accessPayload,err := server.tokenMaker.CreateToken(refreshPayload.Username,server.config.AccessTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError,errorResponse(err))
		return
	}

	rsp := renewAccessTokenResponse{
		AccessToken: accessToken,
		AccessTokenExpiresAt: accessPayload.ExpiredAt,
		
	}
	 ctx.JSON(http.StatusOK,rsp)
}