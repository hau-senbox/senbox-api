package repository

import (
	"errors"
	"sen-global-api/internal/domain/entity"
	"sen-global-api/internal/domain/request"

	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type RolePolicyRepository struct {
	DBConn *gorm.DB
}

func NewRolePolicyRepository(dbConn *gorm.DB) *RolePolicyRepository {
	return &RolePolicyRepository{DBConn: dbConn}
}

func (receiver *RolePolicyRepository) GetAll() ([]entity.SRolePolicy, error) {
	var policies []entity.SRolePolicy
	err := receiver.DBConn.Table("s_role_policy").Find(&policies).Error
	if err != nil {
		log.Error("RolePolicyRepository.GetAll: " + err.Error())
		return nil, errors.New("failed to get all policies")
	}

	return policies, err
}

func (receiver *RolePolicyRepository) GetByID(req request.GetRolePolicyByIdRequest) (*entity.SRolePolicy, error) {
	var policy entity.SRolePolicy
	err := receiver.DBConn.Where("id = ?", req.ID).First(&policy).Error
	if err != nil {
		log.Error("RolePolicyRepository.GetByID: " + err.Error())
		return nil, errors.New("failed to get policy")
	}
	return &policy, nil
}

func (receiver *RolePolicyRepository) GetByName(req request.GetRolePolicyByNameRequest) (*entity.SRolePolicy, error) {
	var policy entity.SRolePolicy
	err := receiver.DBConn.Where("policy_name = ?", req.PolicyName).First(&policy).Error
	if err != nil {
		log.Error("RolePolicyRepository.GetByName: " + err.Error())
		return nil, errors.New("failed to get policy")
	}
	return &policy, nil
}

func (receiver *RolePolicyRepository) CreateRolePolicy(req request.CreateRolePolicyRequest) error {
	policy, _ := receiver.GetByName(request.GetRolePolicyByNameRequest{PolicyName: req.PolicyName})

	if policy != nil {
		log.Error("RolePolicyRepository.CreateRolePolicy: " + policy.PolicyName)
		return errors.New("policy already existed")
	}

	policyReq := entity.SRolePolicy{
		PolicyName:  req.PolicyName,
		Description: req.Description,
	}
	policyResult := receiver.DBConn.Create(&policyReq)

	if policyResult.Error != nil {
		log.Error("RolePolicyRepository.CreateRolePolicy: " + policyResult.Error.Error())
		return errors.New("failed to create policy")
	}

	if req.Roles != nil {
		for _, roleId := range *req.Roles {
			var role entity.SRole
			err := receiver.DBConn.Where("id = ?", roleId).First(&role).Error
			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return errors.New("role not found")
				}
				return errors.New("failed to get role")
			}

			result := receiver.DBConn.Create(&entity.SRolePolicyRoles{
				RoleId:   role.ID,
				PolicyId: policyReq.ID,
			})

			if result.Error != nil {
				log.Error("RolePolicyRepository.CreateRolePolicy: " + result.Error.Error())
				return errors.New("failed to create role policy")
			}
		}
	}

	if req.RoleClaims != nil {
		for _, roleClaimId := range *req.RoleClaims {
			var roleClaim entity.SRoleClaim
			err := receiver.DBConn.Where("id = ?", roleClaimId).First(&roleClaim).Error
			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return errors.New("role claim not found")
				}
				return errors.New("failed to get role claim")
			}

			result := receiver.DBConn.Create(&entity.SRolePolicyClaims{
				ClaimId:  roleClaim.ID,
				PolicyId: policyReq.ID,
			})

			if result.Error != nil {
				log.Error("RolePolicyRepository.CreateRolePolicy: " + result.Error.Error())
				return errors.New("failed to create role policy")
			}
		}
	}

	return nil
}

func (receiver *RolePolicyRepository) UpdateRolePolicy(req request.UpdateRolePolicyRequest) error {
	updateResult := receiver.DBConn.Model(&entity.SRole{}).Where("id = ?", req.ID).
		Updates(map[string]interface{}{
			"policy_name": req.PolicyName,
			"description": req.Description,
		})

	if updateResult.Error != nil {
		log.Error("RolePolicyRepository.UpdateRolePolicy: " + updateResult.Error.Error())
		return errors.New("failed to update policy")
	}

	if req.Roles != nil {
		removeResult := receiver.DBConn.Exec("DELETE FROM s_role_policy_roles WHERE policy_id = ?", req.ID)

		if removeResult.Error != nil {
			log.Error("RolePolicyRepository.UpdateRolePolicy: " + removeResult.Error.Error())
			return errors.New("failed to remove role policy")
		}

		for _, roleId := range *req.Roles {
			var role entity.SRole
			err := receiver.DBConn.Where("id = ?", roleId).First(&role).Error
			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return errors.New("role not found")
				}
				return errors.New("failed to get role")
			}

			result := receiver.DBConn.Create(&entity.SRolePolicyRoles{
				RoleId:   role.ID,
				PolicyId: int64(req.ID),
			})

			if result.Error != nil {
				log.Error("RolePolicyRepository.UpdateRolePolicy: " + result.Error.Error())
				return errors.New("failed to create role policy")
			}
		}
	}

	if req.RoleClaims != nil {
		removeResult := receiver.DBConn.Exec("DELETE FROM s_role_policy_claims WHERE policy_id = ?", req.ID)

		if removeResult.Error != nil {
			log.Error("RolePolicyRepository.UpdateRolePolicy: " + removeResult.Error.Error())
			return errors.New("failed to remove role policy")
		}

		for _, roleClaimId := range *req.RoleClaims {
			var roleClaim entity.SRoleClaim
			err := receiver.DBConn.Where("id = ?", roleClaimId).First(&roleClaim).Error
			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return errors.New("role claim not found")
				}
				return errors.New("failed to get role claim")
			}

			result := receiver.DBConn.Create(&entity.SRolePolicyClaims{
				ClaimId:  roleClaim.ID,
				PolicyId: int64(req.ID),
			})

			if result.Error != nil {
				log.Error("RolePolicyRepository.UpdateRolePolicy: " + result.Error.Error())
				return errors.New("failed to create role policy")
			}
		}
	}

	return nil
}

func (receiver *RolePolicyRepository) DeleteRolePolicy(req request.DeleteRolePolicyRequest) error {
	result := receiver.DBConn.Delete(&entity.SRolePolicy{}, req.ID)
	if result.Error != nil {
		log.Error("RolePolicyRepository.DeleteRolePolicy: " + result.Error.Error())
		return errors.New("failed to delete policy")
	}
	return nil
}
