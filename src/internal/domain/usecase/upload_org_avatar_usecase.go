package usecase

import (
	"fmt"
	"sen-global-api/helper"
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/entity"
	"sen-global-api/pkg/uploader"
)

type UploadOrgAvatarUseCase struct {
	DeleteImageUseCase
	UploadImageUseCase
	repository.OrganizationRepository
}

func (receiver *UploadOrgAvatarUseCase) UploadAvatar(orgID string, data []byte, fileName string) (*string, *entity.SImage, error) {
	org, err := receiver.OrganizationRepository.GetByID(orgID)
	if err != nil {
		return nil, nil, err
	}

	if org.Avatar != "" {
		err = receiver.DeleteImageUseCase.DeleteImage(org.Avatar)
		if err != nil {
			return nil, nil, err
		}

		err = receiver.OrganizationRepository.UpdateOrgAvatar(orgID, "", "")
		if err != nil {
			return nil, nil, err
		}
	}

	url, image, err := receiver.UploadImageUseCase.UploadImage(
		data,
		"avatar",
		fileName,
		fmt.Sprintf("%s_avatar", helper.Slugify(org.OrganizationName)),
		uploader.UploadPrivate,
		nil,
		nil,
		nil,
		nil,
	)
	if err != nil {
		return nil, nil, err
	}

	err = receiver.OrganizationRepository.UpdateOrgAvatar(orgID, image.Key, *url)
	if err != nil {
		return nil, nil, err
	}

	return url, image, nil
}
