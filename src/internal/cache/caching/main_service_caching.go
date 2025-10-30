package caching

import (
	"context"
	"sen-global-api/internal/cache"
	"sen-global-api/internal/domain/entity"
	"sen-global-api/internal/domain/response"
)

type CachingService interface {
	SetUserCache(ctx context.Context, user *entity.SUserEntity) error
	SetStudentCache(ctx context.Context, student *response.GetStudent4Gateway) error
	SetTeacherCache(ctx context.Context, teacher *response.GetTeacher4Gateway) error
	SetStaffCache(ctx context.Context, staff *response.GetStaff4Gateway) error
	SetParentCache(ctx context.Context, parent *response.GetParent4Gateway) error
	SetChildCache(ctx context.Context, child *response.ChildResponse) error
	SetTeacherByUserAndOrgCacheKey(ctx context.Context, userID, orgID string, teacher *response.GetTeacher4Gateway) error
	SetStaffByUserAndOrgCacheKey(ctx context.Context, userID, orgID string, staff *response.GetStaff4Gateway) error
	SetUserByTeacherCacheKey(ctx context.Context, teacherID string, user *entity.SUserEntity) error
	SetParentByUserCacheKey(ctx context.Context, userID string, parent *response.GetParent4Gateway) error

	InvalidUserCache(ctx context.Context, userID string) error
	InvalidTeacherCache(ctx context.Context, teacherID string) error
	InvalidParentCache(ctx context.Context, parentID string) error
	InvalidStaffCache(ctx context.Context, staffID string) error
	InvalidChildCache(ctx context.Context, childID string) error
	InvalidStudentCache(ctx context.Context, studentID string) error
	InvalidUserByTeacherCacheKey(ctx context.Context, teacherID string) error
	InvalidParentByUserCacheKey(ctx context.Context, userID string) error
	InvalidTeacherByUserAndOrgCacheKey(ctx context.Context, userID, orgID string) error
	InvalidStaffByUserAndOrgCacheKey(ctx context.Context, userID, orgID string) error
}

type cachingService struct {
	cache      *cache.RedisCache
	defaultTTL int
}

func NewCachingService(cache *cache.RedisCache, defaultTTL int) CachingService {
	return &cachingService{cache: cache, defaultTTL: defaultTTL}
}

func (s *cachingService) setByKey(ctx context.Context, key string, data any) error {
	return s.cache.Set(ctx, key, data, s.defaultTTL)
}

func (s *cachingService) deleteByKey(ctx context.Context, key string) error {
	return s.cache.Delete(ctx, key)
}

// ========================
// === SET CACHE ===
// ========================

func (s *cachingService) SetUserCache(ctx context.Context, user *entity.SUserEntity) error {
	if user == nil {

	}
	return s.setByKey(ctx, cache.UserCacheKey(user.ID.String()), user)
}

func (s *cachingService) SetStudentCache(ctx context.Context, student *response.GetStudent4Gateway) error {
	if student == nil {
		return nil
	}
	return s.setByKey(ctx, cache.StudentCacheKey(student.StudentID), student)
}

func (s *cachingService) SetTeacherCache(ctx context.Context, teacher *response.GetTeacher4Gateway) error {
	if teacher == nil {
		return nil
	}
	return s.setByKey(ctx, cache.TeacherCacheKey(teacher.TeacherID), teacher)
}

func (s *cachingService) SetParentCache(ctx context.Context, parent *response.GetParent4Gateway) error {
	if parent == nil {
		return nil
	}
	return s.setByKey(ctx, cache.ParentCacheKey(parent.ParentID), parent)
}

func (s *cachingService) SetStaffCache(ctx context.Context, staff *response.GetStaff4Gateway) error {
	if staff == nil {
		return nil
	}
	return s.setByKey(ctx, cache.StaffCacheKey(staff.StaffID), staff)
}

func (s *cachingService) SetChildCache(ctx context.Context, child *response.ChildResponse) error {
	if child == nil {
		return nil
	}
	return s.setByKey(ctx, cache.ChildCacheKey(child.ChildID), child)
}

func (s *cachingService) SetTeacherByUserAndOrgCacheKey(ctx context.Context, userID, orgID string, teacher *response.GetTeacher4Gateway) error {
	return s.setByKey(ctx, cache.TeacherByUserAndOrgCacheKey(userID, orgID), teacher)
}

func (s *cachingService) SetStaffByUserAndOrgCacheKey(ctx context.Context, userID, orgID string, staff *response.GetStaff4Gateway) error {
	return s.setByKey(ctx, cache.StaffByUserAndOrgCacheKey(userID, orgID), staff)
}

func (s *cachingService) SetUserByTeacherCacheKey(ctx context.Context, teacherID string, user *entity.SUserEntity) error {
	if user == nil {
		return nil
	}
	return s.setByKey(ctx, cache.UserByTeacherCacheKey(teacherID), user)
}

func (s *cachingService) SetParentByUserCacheKey(ctx context.Context, userID string, parent *response.GetParent4Gateway) error {
	if parent == nil {
		return nil
	}
	return s.setByKey(ctx, cache.ParentByUserCacheKey(userID), parent)
}

// ========================
// === INVALID CACHE ===
// ========================

func (s *cachingService) InvalidUserCache(ctx context.Context, userID string) error {
	return s.deleteByKey(ctx, cache.UserCacheKey(userID))
}

func (s *cachingService) InvalidStudentCache(ctx context.Context, studentID string) error {
	return s.deleteByKey(ctx, cache.StudentCacheKey(studentID))
}

func (s *cachingService) InvalidTeacherCache(ctx context.Context, teacherID string) error {
	return s.deleteByKey(ctx, cache.TeacherCacheKey(teacherID))
}

func (s *cachingService) InvalidParentCache(ctx context.Context, parentID string) error {
	return s.deleteByKey(ctx, cache.ParentCacheKey(parentID))
}

func (s *cachingService) InvalidStaffCache(ctx context.Context, staffID string) error {
	return s.deleteByKey(ctx, cache.StaffCacheKey(staffID))
}

func (s *cachingService) InvalidChildCache(ctx context.Context, childID string) error {
	return s.deleteByKey(ctx, cache.ChildCacheKey(childID))
}

func (s *cachingService) InvalidUserByTeacherCacheKey(ctx context.Context, teacherID string) error {
	return s.deleteByKey(ctx, cache.UserByTeacherCacheKey(teacherID))
}

func (s *cachingService) InvalidParentByUserCacheKey(ctx context.Context, userID string) error {
	return s.deleteByKey(ctx, cache.ParentByUserCacheKey(userID))
}

func (s *cachingService) InvalidTeacherByUserAndOrgCacheKey(ctx context.Context, userID, orgID string) error {
	return s.deleteByKey(ctx, cache.TeacherByUserAndOrgCacheKey(userID, orgID))
}

func (s *cachingService) InvalidStaffByUserAndOrgCacheKey(ctx context.Context, userID, orgID string) error {
	return s.deleteByKey(ctx, cache.StaffByUserAndOrgCacheKey(userID, orgID))
}
