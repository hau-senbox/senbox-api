package usecase

import (
	"errors"
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/entity"
	"sen-global-api/internal/domain/request"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ChildUseCase struct {
	childRepo repository.ChildRepository
}

func NewChildUseCase(
	childRepo repository.ChildRepository,
) *ChildUseCase {
	return &ChildUseCase{
		childRepo: childRepo,
	}
}

func (uc *ChildUseCase) CreateChild(req request.CreateChildRequest, ctx *gin.Context) error {
	userIDRaw, exists := ctx.Get("user_id")
	if !exists {
		return errors.New("unauthorized: user_id not found in context")
	}

	var userID uuid.UUID
	switch v := userIDRaw.(type) {
	case uuid.UUID:
		userID = v
	case string:
		parsed, err := uuid.Parse(v)
		if err != nil {
			return errors.New("invalid user_id format")
		}
		userID = parsed
	default:
		return errors.New("invalid user_id type in context")
	}

	child := &entity.SChild{
		ChildName: req.ChildName,
		Age:       req.Age,
		ParentID:  userID,
	}

	return uc.childRepo.Create(child)
}

func (uc *ChildUseCase) UpdateChild(child *entity.SChild) error {
	return uc.childRepo.Update(child)
}
