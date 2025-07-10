package usecase

import (
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/entity"

	"github.com/google/uuid"
)

type AnswerUseCase struct {
	answerRepo repository.AnswerRepository
}

func NewAnswerUseCase(answerRepo repository.AnswerRepository) *AnswerUseCase {
	return &AnswerUseCase{
		answerRepo: answerRepo,
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
func (uc *AnswerUseCase) GetAnswersByKeyAndDB(key, db string) ([]entity.SAnswer, error) {
	return uc.answerRepo.FindByKeyAndDB(key, db)
}
