package usecase

import (
	"fmt"
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/response"
)

type GetUserOrganizationActiveUsecase struct {
	*repository.OrganizationRepository
	*repository.TeacherApplicationRepository
	*repository.StaffApplicationRepository
}

func (receiver *GetUserOrganizationActiveUsecase) GetUserOrganizationActive(userID string) (*response.UserOrganizationActive, error) {
	// get teacher organizations
	teachers, err := receiver.TeacherApplicationRepository.GetByUserIDApproved(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get teacher applications: %w", err)
	}

	teacherOrgs := make([]response.OrganizationActive, 0, len(teachers))
	for _, t := range teachers {
		// get org info
		orgInfo, err := receiver.OrganizationRepository.GetByID(t.OrganizationID.String())
		if err != nil {
			return nil, err
		}

		teacherOrgs = append(teacherOrgs, response.OrganizationActive{
			ID:               orgInfo.ID.String(),
			OrganizationName: orgInfo.OrganizationName,
			Avatar:           orgInfo.Avatar,
			AvatarURL:        orgInfo.AvatarURL,
			CreatedAt:        orgInfo.CreatedAt,
		})
	}

	// get staff organizations
	staffs, err := receiver.StaffApplicationRepository.GetByUserIDApproved(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get staff applications: %w", err)
	}

	staffOrgs := make([]response.OrganizationActive, 0, len(staffs))
	for _, s := range staffs {
		orgInfo, err := receiver.OrganizationRepository.GetByID(s.OrganizationID.String())
		if err != nil {
			return nil, err
		}
		staffOrgs = append(staffOrgs, response.OrganizationActive{
			ID:               orgInfo.ID.String(),
			OrganizationName: orgInfo.OrganizationName,
			Avatar:           orgInfo.Avatar,
			AvatarURL:        orgInfo.AvatarURL,
			CreatedAt:        orgInfo.CreatedAt,
		})
	}

	// build final response
	return &response.UserOrganizationActive{
		TeacherOrganization: teacherOrgs,
		StaffOrganization:   staffOrgs,
	}, nil
}
