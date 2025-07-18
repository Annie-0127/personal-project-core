package auth

import (
	"context"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"personal-project-core/constant"
	svCtx "personal-project-core/context"
	"personal-project-core/helpers"
	"time"
)

type AuthServer struct {
	sc *svCtx.ServiceContext
}

func NewAuthServer(sc *svCtx.ServiceContext) *AuthServer {
	return &AuthServer{
		sc: sc,
	}
}

type JWTClaims struct {
	JWTPayload
	jwt.StandardClaims
}

type JWTPayload struct {
	ID       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	Username string    `json:"username"`
	Email    string    `json:"email"`
	Roles    []string  `json:"roles"`
}

func (a *AuthServer) PostAuthLogin(ctx context.Context, request PostAuthLoginRequestObject) (PostAuthLoginResponseObject, error) {
	cond := map[string]any{
		"username":  request.Body.Username,
		"is_active": constant.ACTIVE,
	}

	user, err := a.sc.UserRepo.FindOne(ctx, cond, []string{})
	if err != nil {
		a.sc.Log.Error("Failed to get user", zap.Error(err), zap.Any("request", request.Body))
		return PostAuthLogin400JSONResponse{
			Message: "Failed to get user",
			Code:    constant.PREFIX_USER + constant.ERROR_DATA_NOT_FOUND,
		}, nil
	}

	valid := helpers.CheckPasswordHash(request.Body.Password, user.Password)
	if !valid {
		a.sc.Log.Info("Invalid password", zap.Any("request", request.Body))
		return PostAuthLogin400JSONResponse{
			Message: "Invalid password",
			Code:    constant.PREFIX_USER + constant.ERROR_LOGIN_FAILED,
		}, nil
	}

	expiresAt := time.Now().Add(a.sc.Config.JWT.ExpiredAt)

	claim := JWTClaims{
		JWTPayload{
			ID:       user.Id,
			Name:     user.Name,
			Username: user.Username,
			Email:    user.Email,
			Roles:    user.RoleIds,
		},
		jwt.StandardClaims{
			ExpiresAt: expiresAt.Unix(),
		},
	}
	jwtSecret := a.sc.Config.JWT.Secret
	token, err := helpers.CreateToken(jwtSecret, claim)
	if err != nil {
		a.sc.Log.Error("Failed to create token", zap.Error(err))
		return PostAuthLogin400JSONResponse{
			Message: "Failed to create token",
			Code:    constant.PREFIX_USER + constant.ERROR_LOGIN_FAILED,
		}, nil
	}
	response := PostAuthLogin200JSONResponse{
		Message: "Logged in successfully",
	}
	response.Data.Name = user.Name
	response.Data.Username = user.Username
	response.Data.Email = user.Email
	response.Data.Token = token
	response.Data.UserId = user.Id.String()
	return response, nil
}
