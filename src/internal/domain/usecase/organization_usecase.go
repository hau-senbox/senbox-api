package usecase

import (
	"sen-global-api/internal/data/repository"
	"sen-global-api/pkg/consulapi/gateway"

	"github.com/gin-gonic/gin"
)

type OrganizationUsecase struct {
	OrganizationRepository *repository.OrganizationRepository
	ProfileGateway         gateway.ProfileGateway
}

func NewOrganizationUsecase(organizationRepository *repository.OrganizationRepository, profileGateway gateway.ProfileGateway) *OrganizationUsecase {
	return &OrganizationUsecase{OrganizationRepository: organizationRepository, ProfileGateway: profileGateway}
}

func (receiver *OrganizationUsecase) GenerateOrganizationCode(ctx *gin.Context) {
	organizations, _ := receiver.OrganizationRepository.GetAllOrganizations()

	for _, organization := range organizations {
		_, _ = receiver.ProfileGateway.GenerateOrganizationCode(ctx, organization.ID.String(), organization.CreatedIndex)
	}
}
