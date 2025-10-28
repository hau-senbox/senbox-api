package cached

import (
	"sen-global-api/internal/cache"
	"sen-global-api/pkg/consulapi/gateway"

	"github.com/gin-gonic/gin"
)

type CachedProfileGateway struct {
	inner gateway.ProfileGateway
	cache cache.Cache
	ttl   int
}

func NewCachedProfileGateway(inner gateway.ProfileGateway, cache cache.Cache, ttl int) gateway.ProfileGateway {
	return &CachedProfileGateway{
		inner: inner,
		cache: cache,
		ttl:   ttl,
	}
}

// -------------------------
// Generic helper
// -------------------------

func (c *CachedProfileGateway) getCode(
	ctx *gin.Context,
	cacheKey string,
	fn func(*gin.Context, string) (string, error),
	id string,
) (string, error) {
	var cached string
	if err := c.cache.Get(ctx.Request.Context(), cacheKey, &cached); err == nil && cached != "" {
		return cached, nil
	}

	code, err := fn(ctx, id)
	if err != nil {
		return "", err
	}

	// set lai
	if err := c.cache.Set(ctx.Request.Context(), cacheKey, code, c.ttl); err != nil {
		return "", err
	}
	return code, nil
}

// -------------------------
// Public methods
// -------------------------

func (c *CachedProfileGateway) GetStudentCode(ctx *gin.Context, studentID string) (string, error) {
	return c.getCode(ctx, cache.StudentCodeCacheKey(studentID), c.inner.GetStudentCode, studentID)
}

func (c *CachedProfileGateway) GetTeacherCode(ctx *gin.Context, teacherID string) (string, error) {
	return c.getCode(ctx, cache.TeacherCodeCacheKey(teacherID), c.inner.GetTeacherCode, teacherID)
}

func (c *CachedProfileGateway) GetParentCode(ctx *gin.Context, parentID string) (string, error) {
	return c.getCode(ctx, cache.ParentCodeCacheKey(parentID), c.inner.GetParentCode, parentID)
}

func (c *CachedProfileGateway) GetStaffCode(ctx *gin.Context, staffID string) (string, error) {
	return c.getCode(ctx, cache.StaffCodeCacheKey(staffID), c.inner.GetStaffCode, staffID)
}

func (c *CachedProfileGateway) GetChildCode(ctx *gin.Context, childID string) (string, error) {
	return c.getCode(ctx, cache.ChildCodeCacheKey(childID), c.inner.GetChildCode, childID)
}

func (c *CachedProfileGateway) GetUserCode(ctx *gin.Context, userID string) (string, error) {
	return c.getCode(ctx, cache.UserCodeCacheKey(userID), c.inner.GetUserCode, userID)
}

func (c *CachedProfileGateway) GenerateUserCode(ctx *gin.Context, userID string, createdIndex int) (*string, error) {
	return c.inner.GenerateUserCode(ctx, userID, createdIndex)
}

func (c *CachedProfileGateway) GenerateChildCode(ctx *gin.Context, childID string, createdIndex int) (*string, error) {
	return c.inner.GenerateChildCode(ctx, childID, createdIndex)
}

func (c *CachedProfileGateway) GenerateParentCode(ctx *gin.Context, parentID string, createdIndex int) (*string, error) {
	return c.inner.GenerateParentCode(ctx, parentID, createdIndex)
}

func (c *CachedProfileGateway) GenerateStaffCode(ctx *gin.Context, staffID string, createdIndex int) (*string, error) {
	return c.inner.GenerateStaffCode(ctx, staffID, createdIndex)
}

func (c *CachedProfileGateway) GenerateTeacherCode(ctx *gin.Context, teacherID string, createdIndex int) (*string, error) {
	return c.inner.GenerateTeacherCode(ctx, teacherID, createdIndex)
}

func (c *CachedProfileGateway) GenerateStudentCode(ctx *gin.Context, studentID string, createdIndex int) (*string, error) {
	return c.inner.GenerateStudentCode(ctx, studentID, createdIndex)
}
