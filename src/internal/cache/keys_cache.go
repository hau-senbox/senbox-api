package cache

import "sen-global-api/internal/domain/value"

// ==============================
// === Common Cache Key Utils ===
// ==============================

// entity data
func UserCacheKey(userID string) string {
	return value.ProfileCachePrefix + "user:" + userID
}

func StudentCacheKey(studentID string) string {
	return value.ProfileCachePrefix + "student:" + studentID
}

func TeacherCacheKey(teacherID string) string {
	return value.ProfileCachePrefix + "teacher:" + teacherID
}

func StaffCacheKey(staffID string) string {
	return value.ProfileCachePrefix + "staff:" + staffID
}

func ParentCacheKey(parentID string) string {
	return value.ProfileCachePrefix + "parent:" + parentID
}

func ChildCacheKey(childID string) string {
	return value.ProfileCachePrefix + "child:" + childID
}

// code mapping
func UserCodeCacheKey(userID string) string {
	return value.ProfileCachePrefix + "user_code:" + userID
}

func StudentCodeCacheKey(studentID string) string {
	return value.ProfileCachePrefix + "student_code:" + studentID
}

func TeacherCodeCacheKey(teacherID string) string {
	return value.ProfileCachePrefix + "teacher_code:" + teacherID
}

func StaffCodeCacheKey(staffID string) string {
	return value.ProfileCachePrefix + "staff_code:" + staffID
}

func ParentCodeCacheKey(parentID string) string {
	return value.ProfileCachePrefix + "parent_code:" + parentID
}

func ChildCodeCacheKey(childID string) string {
	return value.ProfileCachePrefix + "child_code:" + childID
}

// generic
func GenericCacheKey(prefix, id string) string {
	return value.ProfileCachePrefix + prefix + ":" + id
}
