package usecase

import (
	"encoding/json"
	"fmt"
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/entity"
	"sen-global-api/internal/domain/request"
	"sen-global-api/internal/domain/response"
	"sort"
	"strconv"

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

func (uc *AnswerUseCase) GetTotalNrByKeyAndDb(input repository.GetSubmissionByConditionParam) (response.GetTotalNrByKeyAndDbResponse, error) {
	answers, err := uc.answerRepo.GetTotalByKeyAndDb(input)
	if err != nil {
		return response.GetTotalNrByKeyAndDbResponse{}, err
	}

	var total float32

	for _, ans := range answers {
		var strValue string
		if err := json.Unmarshal(ans.Response, &strValue); err != nil {
			continue
		}

		if value, err := strconv.ParseFloat(strValue, 32); err == nil {
			total += float32(value)
		}
	}

	return response.GetTotalNrByKeyAndDbResponse{
		Total: total,
	}, nil
}

func (uc *AnswerUseCase) GetChartTotalByDay(input repository.GetSubmissionByConditionParam) ([]response.ChartDataResponse, error) {
	answers, err := uc.answerRepo.GetChartTotalByKeyAndDb(input)
	if err != nil {
		return nil, err
	}

	result := make(map[string]float32)

	for _, ans := range answers {
		// Parse ngày yyyy-MM-dd
		day := ans.CreatedAt.Format("2006-01-02")

		// Parse từ RawMessage
		var strValue string
		if err := json.Unmarshal(ans.Response, &strValue); err != nil {
			continue
		}
		if value, err := strconv.ParseFloat(strValue, 32); err == nil {
			result[day] += float32(value)
		}
	}

	// Convert map to slice, đảm bảo có thứ tự
	var chartData []response.ChartDataResponse
	for day, total := range result {
		chartData = append(chartData, response.ChartDataResponse{
			X: day,
			Y: fmt.Sprintf("%.2f", total),
		})
	}

	// Sắp xếp theo ngày tăng dần (nếu cần)
	sort.Slice(chartData, func(i, j int) bool {
		return chartData[i].X < chartData[j].X
	})

	return chartData, nil
}
