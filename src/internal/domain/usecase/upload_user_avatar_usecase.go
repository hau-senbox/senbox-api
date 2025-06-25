package usecase

import (
	"fmt"
	"sen-global-api/helper"
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/entity"
	"sen-global-api/internal/domain/request"
	"sen-global-api/pkg/uploader"
)

type UploadUserAvatarUseCase struct {
	UploadImageUseCase
	DeleteImageUseCase
	repository.UserEntityRepository
}

func (receiver *UploadUserAvatarUseCase) UploadAvatar(userID string, data []byte, fileName string) (*string, *entity.SImage, error) {
	user, err := receiver.UserEntityRepository.GetByID(request.GetUserEntityByIDRequest{ID: userID})
	if err != nil {
		return nil, nil, err
	}

	if user.Avatar != "" {
		err = receiver.DeleteImageUseCase.DeleteImage(user.Avatar)
		if err != nil {
			return nil, nil, err
		}

		err = receiver.UserEntityRepository.UpdateUserAvatar(userID, "", "")
		if err != nil {
			return nil, nil, err
		}
	}

	url, image, err := receiver.UploadImageUseCase.UploadImage(
		data,
		"avatar",
		fileName,
		fmt.Sprintf("%s_avatar", helper.Slugify(user.Username)),
		uploader.UploadPrivate,
	)
	if err != nil {
		return nil, nil, err
	}

	err = receiver.UserEntityRepository.UpdateUserAvatar(userID, image.Key, *url)
	if err != nil {
		return nil, nil, err
	}

	return url, image, nil
}
