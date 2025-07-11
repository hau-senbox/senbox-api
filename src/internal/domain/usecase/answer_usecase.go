package usecase

import (
	"encoding/json"
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/entity"
	"sen-global-api/internal/domain/request"
	"sen-global-api/internal/domain/response"

	"github.com/google/uuid"
)

type AnswerUseCase struct {
	answerRepo repository.AnswerRepository
	userRepo   repository.UserEntityRepository
}

func NewAnswerUseCase(
	answerRepo repository.AnswerRepository,
	userRepo repository.UserEntityRepository,
) *AnswerUseCase {
	return &AnswerUseCase{
		answerRepo: answerRepo,
		userRepo:   userRepo,
	}
}

// Tạo câu trả lời mới
func (uc *AnswerUseCase) CreateAnswer(answer *entity.SAnswer) error {
	return uc.answerRepo.Create(answer)
}

// Lấy câu trả lời theo ID
func (uc *AnswerUseCase) GetAnswerByID(id uuid.UUID) (*entity.SAnswer, error) {
	return uc.answerRepo.FindByID(id)
}

// Lấy danh sách câu trả lời theo submissionID
func (uc *AnswerUseCase) GetAnswersBySubmissionID(submissionID string) ([]entity.SAnswer, error) {
	return uc.answerRepo.FindBySubmissionID(submissionID)
}

// Cập nhật câu trả lời
func (uc *AnswerUseCase) UpdateAnswer(answer *entity.SAnswer) (*entity.SAnswer, error) {
	return uc.answerRepo.Update(answer)
}

// Xóa câu trả lời theo ID
func (uc *AnswerUseCase) DeleteAnswer(id uuid.UUID) error {
	return uc.answerRepo.Delete(id)
}

// Lấy danh sách câu trả lời theo key và db
func (uc *AnswerUseCase) GetAnswersByKeyAndDB(input repository.GetSubmissionByConditionParam) ([]response.GetAnswerByKeyAndDbResponse, error) {
	answers, err := uc.answerRepo.FindByKeyAndDB(input)
	if err != nil {
		return nil, err
	}

	var result []response.GetAnswerByKeyAndDbResponse
	seenUserIDs := make(map[string]bool)
	for _, a := range answers {
		//bo qua neu trung user_id
		if seenUserIDs[a.UserID] {
			continue
		}
		seenUserIDs[a.UserID] = true

		var answerStr string
		_ = json.Unmarshal(a.Response, &answerStr)
		var UserNickName string
		user, err := uc.userRepo.GetByID(request.GetUserEntityByIDRequest{ID: a.UserID})
		if err != nil {
			continue
		}
		UserNickName = user.Nickname
		res := response.GetAnswerByKeyAndDbResponse{
			ID:           a.ID.String(),
			SubmissionID: a.SubmissionID,
			UserID:       a.UserID,
			UserNickName: UserNickName,
			Key:          a.Key,
			DB:           a.DB,
			Answer:       answerStr,
			CreatedAt:    a.CreatedAt,
		}
		result = append(result, res)
	}

	return result, nil
}
