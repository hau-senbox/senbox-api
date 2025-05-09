package usecase

import (
	"errors"
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/entity"
	"sen-global-api/internal/domain/request"
	"strings"

	"github.com/gin-gonic/gin"
)

type GetUserFromTokenUseCase struct {
	repository.UserEntityRepository
	repository.SessionRepository
}

func (c *GetUserFromTokenUseCase) GetUserFromToken(context *gin.Context) (*entity.SUserEntity, error) {
	authorization := context.GetHeader("Authorization")
	if authorization == "" {
		return nil, errors.New("no authorization header")
	}

	if len(authorization) == 0 {
		return nil, errors.New("no authorization header")
	}

	tokenString := strings.Split(authorization, " ")[1]

	userId, err := c.ExtractUserIdFromToken(tokenString)
	if err != nil {
		return nil, err
	}

	return c.GetByID(request.GetUserEntityByIdRequest{ID: *userId})
}
