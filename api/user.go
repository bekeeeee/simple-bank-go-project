package api

import (
	"fmt"
	db "github/bekeeeee/simplebank/db/sqlc"
	"github/bekeeeee/simplebank/util"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

type userResponse struct{
	Username          string    `json:"username"`
	FullName          string    `json:"full_name"`
	Email             string    `json:"email"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
	CreatedAt         time.Time `json:"created_at"`
}

type createUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
	FullName string `json:"full_name" binding:"required"`
	Email string `json:"email" binding:"required,email"`
}


func newUserResponse(user db.User) userResponse{
	return userResponse{
		Username: user.Username,
		FullName: user.FullName,
		Email: user.Email,
		PasswordChangedAt: user.PasswordChangedAt,
		CreatedAt: user.CreatedAt,
	}
}
func (server *Server) createUser(ctx *gin.Context){
	 var req createUserRequest
	 if err := ctx.ShouldBindJSON(&req); err!=nil {
		ctx.JSON(http.StatusBadRequest,errorResponse(err))
		return
	 }
	 hashedPassword,err:= util.HashPassword(req.Password)
	 if err != nil {
		ctx.JSON(http.StatusInternalServerError,errorResponse(err))
		return
	 }
	 arg := db.CreateUserParams{
		Username: req.Username,
		HashedPassword: hashedPassword,
		FullName: req.FullName,
		Email: req.Email,
		}
	user,err := server.store.CreateUser(ctx,arg) 
	if err != nil {
		if pqErr,ok := err.(*pq.Error); ok{
			log.Println(pqErr.Code.Name())
			switch pqErr.Code.Name(){
			case "unique_violation":
				ctx.JSON(http.StatusForbidden,errorResponse(err))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError,errorResponse(err))
		return;
	}
	rsp := newUserResponse(user)
	ctx.JSON(http.StatusOK,rsp)

}

// type getUserRequest struct {
// 	Username int64 `uri:"id" binding:"required,min=1"`
// }

// func (server *Server) getAccount(ctx *gin.Context){
// 	var req getAccountRequest
// 	if err := ctx.ShouldBindUri(&req); err!=nil {
// 	   ctx.JSON(http.StatusBadRequest,errorResponse(err))
// 	   return
// 	}
// 	account,err := server.store.GetAccount(ctx,req.ID);
	
// 	if err != nil {
// 		if err == sql.ErrNoRows{
// 		ctx.JSON(http.StatusNotFound,errorResponse(err))
// 			return
// 		}
// 		ctx.JSON(http.StatusInternalServerError,errorResponse(err))
// 		return;
// 	}
// 	ctx.JSON(http.StatusOK,account)
// }


type loginUserResponse struct {
	AccessToken string `json:"access_token"`
	User userResponse `json:"user"`

}

type loginUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`


}

func (server *Server) loginUser(ctx *gin.Context){
	var req loginUserRequest
	if err := ctx.ShouldBindJSON(&req); err!=nil {
		ctx.JSON(http.StatusBadRequest,errorResponse(err))
		return
	}
	user,err := server.store.GetUser(ctx,req.Username)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized,errorResponse(err))
		return
	}
	if err := util.CheckPassword(req.Password,user.HashedPassword); err != nil {
		ctx.JSON(http.StatusUnauthorized,errorResponse(err))
		return
	}
	accessToken,err := server.tokenMaker.CreateToken(user.Username,server.config.AccessTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError,errorResponse(err))
		return
	}
	rsp := loginUserResponse{
		AccessToken: accessToken,
		User: newUserResponse(user),
	}
	fmt.Println("rsp",rsp)
	 ctx.JSON(http.StatusOK,rsp)
}