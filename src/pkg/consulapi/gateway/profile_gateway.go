package gateway

import (
	"encoding/json"
	"fmt"
	"sen-global-api/pkg/consulapi/gateway/dto/request"
	"sen-global-api/pkg/consulapi/gateway/dto/response"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/hashicorp/consul/api"

	"github.com/hung-senbox/senbox-cache-service/pkg/cache/cached"
	"github.com/hung-senbox/senbox-cache-service/pkg/cache/caching"
)

type ProfileGateway interface {
	// generate owner codes
	GenerateStudentCode(ctx *gin.Context, studentID string, createdIndex int) (*string, error)
	GenerateTeacherCode(ctx *gin.Context, teacherID string, createdIndex int) (*string, error)
	GenerateStaffCode(ctx *gin.Context, staffID string, createdIndex int) (*string, error)
	GenerateParentCode(ctx *gin.Context, parentID string, createdIndex int) (*string, error)
	GenerateUserCode(ctx *gin.Context, userID string, createdIndex int) (*string, error)
	GenerateChildCode(ctx *gin.Context, childID string, createdIndex int) (*string, error)

	// get owner codes
	GetStudentCode(ctx *gin.Context, studentID string) (string, error)
	GetTeacherCode(ctx *gin.Context, teacherID string) (string, error)
	GetStaffCode(ctx *gin.Context, staffID string) (string, error)
	GetParentCode(ctx *gin.Context, parentID string) (string, error)
	GetUserCode(ctx *gin.Context, userID string) (string, error)
	GetChildCode(ctx *gin.Context, childID string) (string, error)
}

type profileGateway struct {
	serviceName           string
	consul                *api.Client
	cachedProfileGateway  cached.CachedProfileGateway
	cachingProfileGateway caching.CachingProfileService
}

func NewProfileGateway(serviceName string, consulClient *api.Client, cachedProfileGateway cached.CachedProfileGateway, cachingProfileGateway caching.CachingProfileService) ProfileGateway {
	return &profileGateway{
		serviceName:           serviceName,
		consul:                consulClient,
		cachedProfileGateway:  cachedProfileGateway,
		cachingProfileGateway: cachingProfileGateway,
	}
}

func (pg *profileGateway) GenerateStudentCode(ctx *gin.Context, ownerID string, createdIndex int) (*string, error) {
	// Lấy token từ context (được set ở SecuredMiddleware)
	token, exists := ctx.Get("token")
	if !exists {
		return nil, fmt.Errorf("token not found in context")
	}

	tokenStr, ok := token.(string)
	if !ok || tokenStr == "" {
		return nil, fmt.Errorf("invalid token in context")
	}

	client, err := NewGatewayClient(pg.serviceName, tokenStr, pg.consul, nil)
	if err != nil {
		return nil, err
	}

	appLanguage, _ := ctx.Get("app_language")

	headers := make(map[string]string)
	if lang, ok := appLanguage.(uint); ok {
		headers["X-App-Language"] = strconv.Itoa(int(lang))
	}

	req := request.GenerateOwnerCodeRequest{
		OwnerID:      ownerID,
		CreatedIndex: createdIndex,
	}

	resp, err := client.Call("POST", "/api/v1/gateway/profiles/owner-code/student/generate", req, headers)
	if err != nil {
		return nil, err
	}

	var gwResp response.APIGateWayResponse[*string]
	if err := json.Unmarshal(resp, &gwResp); err != nil {
		return nil, fmt.Errorf("unmarshal response fail: %w", err)
	}

	if gwResp.StatusCode != 200 {
		return nil, fmt.Errorf("call gateway generate student code fail: %s", gwResp.Message)
	}

	return gwResp.Data, nil
}

func (pg *profileGateway) GenerateTeacherCode(ctx *gin.Context, ownerID string, createdIndex int) (*string, error) {
	// Lấy token từ context (được set ở SecuredMiddleware)
	token, exists := ctx.Get("token")
	if !exists {
		return nil, fmt.Errorf("token not found in context")
	}

	tokenStr, ok := token.(string)
	if !ok || tokenStr == "" {
		return nil, fmt.Errorf("invalid token in context")
	}

	client, err := NewGatewayClient(pg.serviceName, tokenStr, pg.consul, nil)
	if err != nil {
		return nil, err
	}

	appLanguage, _ := ctx.Get("app_language")

	headers := make(map[string]string)
	if lang, ok := appLanguage.(uint); ok {
		headers["X-App-Language"] = strconv.Itoa(int(lang))
	}

	req := request.GenerateOwnerCodeRequest{
		OwnerID:      ownerID,
		CreatedIndex: createdIndex,
	}

	resp, err := client.Call("POST", "/api/v1/gateway/profiles/owner-code/teacher/generate", req, headers)
	if err != nil {
		return nil, err
	}

	var gwResp response.APIGateWayResponse[*string]
	if err := json.Unmarshal(resp, &gwResp); err != nil {
		return nil, fmt.Errorf("unmarshal response fail: %w", err)
	}

	if gwResp.StatusCode != 200 {
		return nil, fmt.Errorf("call gateway generate teacher code fail: %s", gwResp.Message)
	}

	return gwResp.Data, nil
}

func (pg *profileGateway) GenerateStaffCode(ctx *gin.Context, ownerID string, createdIndex int) (*string, error) {
	// Lấy token từ context (được set ở SecuredMiddleware)
	token, exists := ctx.Get("token")
	if !exists {
		return nil, fmt.Errorf("token not found in context")
	}

	tokenStr, ok := token.(string)
	if !ok || tokenStr == "" {
		return nil, fmt.Errorf("invalid token in context")
	}

	client, err := NewGatewayClient(pg.serviceName, tokenStr, pg.consul, nil)
	if err != nil {
		return nil, err
	}

	appLanguage, _ := ctx.Get("app_language")

	headers := make(map[string]string)
	if lang, ok := appLanguage.(uint); ok {
		headers["X-App-Language"] = strconv.Itoa(int(lang))
	}

	req := request.GenerateOwnerCodeRequest{
		OwnerID:      ownerID,
		CreatedIndex: createdIndex,
	}

	resp, err := client.Call("POST", "/api/v1/gateway/profiles/owner-code/staff/generate", req, headers)
	if err != nil {
		return nil, err
	}

	var gwResp response.APIGateWayResponse[*string]
	if err := json.Unmarshal(resp, &gwResp); err != nil {
		return nil, fmt.Errorf("unmarshal response fail: %w", err)
	}

	if gwResp.StatusCode != 200 {
		return nil, fmt.Errorf("call gateway generate staff code fail: %s", gwResp.Message)
	}

	return gwResp.Data, nil
}

func (pg *profileGateway) GenerateParentCode(ctx *gin.Context, ownerID string, createdIndex int) (*string, error) {
	// Lấy token từ context (được set ở SecuredMiddleware)
	token, exists := ctx.Get("token")
	if !exists {
		return nil, fmt.Errorf("token not found in context")
	}

	tokenStr, ok := token.(string)
	if !ok || tokenStr == "" {
		return nil, fmt.Errorf("invalid token in context")
	}

	client, err := NewGatewayClient(pg.serviceName, tokenStr, pg.consul, nil)
	if err != nil {
		return nil, err
	}

	appLanguage, _ := ctx.Get("app_language")

	headers := make(map[string]string)
	if lang, ok := appLanguage.(uint); ok {
		headers["X-App-Language"] = strconv.Itoa(int(lang))
	}

	req := request.GenerateOwnerCodeRequest{
		OwnerID:      ownerID,
		CreatedIndex: createdIndex,
	}

	resp, err := client.Call("POST", "/api/v1/gateway/profiles/owner-code/parent/generate", req, headers)
	if err != nil {
		return nil, err
	}

	var gwResp response.APIGateWayResponse[*string]
	if err := json.Unmarshal(resp, &gwResp); err != nil {
		return nil, fmt.Errorf("unmarshal response fail: %w", err)
	}

	if gwResp.StatusCode != 200 {
		return nil, fmt.Errorf("call gateway generate parent code fail: %s", gwResp.Message)
	}

	return gwResp.Data, nil
}

func (pg *profileGateway) GenerateUserCode(ctx *gin.Context, ownerID string, createdIndex int) (*string, error) {

	client, err := NewGatewayClient(pg.serviceName, "", pg.consul, nil)
	if err != nil {
		return nil, err
	}

	appLanguage, _ := ctx.Get("app_language")

	headers := make(map[string]string)
	if lang, ok := appLanguage.(uint); ok {
		headers["X-App-Language"] = strconv.Itoa(int(lang))
	}

	req := request.GenerateOwnerCodeRequest{
		OwnerID:      ownerID,
		CreatedIndex: createdIndex,
	}

	resp, err := client.Call("POST", "/api/v1/gateway/profiles/owner-code/user/generate", req, headers)
	if err != nil {
		return nil, err
	}

	var gwResp response.APIGateWayResponse[*string]
	if err := json.Unmarshal(resp, &gwResp); err != nil {
		return nil, fmt.Errorf("unmarshal response fail: %w", err)
	}

	if gwResp.StatusCode != 200 {
		return nil, fmt.Errorf("call gateway generate user code fail: %s", gwResp.Message)
	}

	return gwResp.Data, nil
}

func (pg *profileGateway) GenerateChildCode(ctx *gin.Context, ownerID string, createdIndex int) (*string, error) {
	// Lấy token từ context (được set ở SecuredMiddleware)
	token, exists := ctx.Get("token")
	if !exists {
		return nil, fmt.Errorf("token not found in context")
	}

	tokenStr, ok := token.(string)
	if !ok || tokenStr == "" {
		return nil, fmt.Errorf("invalid token in context")
	}

	client, err := NewGatewayClient(pg.serviceName, tokenStr, pg.consul, nil)
	if err != nil {
		return nil, err
	}

	appLanguage, _ := ctx.Get("app_language")

	headers := make(map[string]string)
	if lang, ok := appLanguage.(uint); ok {
		headers["X-App-Language"] = strconv.Itoa(int(lang))
	}

	req := request.GenerateOwnerCodeRequest{
		OwnerID:      ownerID,
		CreatedIndex: createdIndex,
	}

	resp, err := client.Call("POST", "/api/v1/gateway/profiles/owner-code/child/generate", req, headers)
	if err != nil {
		return nil, err
	}

	var gwResp response.APIGateWayResponse[*string]
	if err := json.Unmarshal(resp, &gwResp); err != nil {
		return nil, fmt.Errorf("unmarshal response fail: %w", err)
	}

	if gwResp.StatusCode != 200 {
		return nil, fmt.Errorf("call gateway generate child code fail: %s", gwResp.Message)
	}

	return gwResp.Data, nil
}

// =============== Get owner codes ===============
func (pg *profileGateway) GetStudentCode(ctx *gin.Context, ownerID string) (string, error) {

	cachedData, _ := pg.cachedProfileGateway.GetStudentCode(ctx, ownerID)
	if cachedData != "" {
		return cachedData, nil
	}

	// Lấy token từ context (được set ở SecuredMiddleware)
	token, exists := ctx.Get("token")
	if !exists {
		return "", fmt.Errorf("token not found in context")
	}

	tokenStr, ok := token.(string)
	if !ok || tokenStr == "" {
		return "", fmt.Errorf("invalid token in context")
	}

	client, err := NewGatewayClient(pg.serviceName, tokenStr, pg.consul, nil)
	if err != nil {
		return "", err
	}

	appLanguage, _ := ctx.Get("app_language")

	headers := make(map[string]string)
	if lang, ok := appLanguage.(uint); ok {
		headers["X-App-Language"] = strconv.Itoa(int(lang))
	}

	resp, err := client.Call("GET", fmt.Sprintf("/api/v1/gateway/profiles/owner-code/student/%s", ownerID), nil, headers)
	if err != nil {
		return "", err
	}

	var gwResp response.APIGateWayResponse[string]
	if err := json.Unmarshal(resp, &gwResp); err != nil {
		return "", fmt.Errorf("unmarshal response fail: %w", err)
	}

	if gwResp.StatusCode != 200 {
		return "", fmt.Errorf("call gateway get student code fail: %s", gwResp.Message)
	}

	pg.cachingProfileGateway.SetStudentCode(ctx, ownerID, gwResp.Data)

	return gwResp.Data, nil
}

func (pg *profileGateway) GetTeacherCode(ctx *gin.Context, ownerID string) (string, error) {

	cachedData, _ := pg.cachedProfileGateway.GetTeacherCode(ctx, ownerID)
	if cachedData != "" {
		return cachedData, nil
	}

	// Lấy token từ context (được set ở SecuredMiddleware)
	token, exists := ctx.Get("token")
	if !exists {
		return "", fmt.Errorf("token not found in context")
	}

	tokenStr, ok := token.(string)
	if !ok || tokenStr == "" {
		return "", fmt.Errorf("invalid token in context")
	}

	client, err := NewGatewayClient(pg.serviceName, tokenStr, pg.consul, nil)
	if err != nil {
		return "nil", err
	}

	appLanguage, _ := ctx.Get("app_language")

	headers := make(map[string]string)
	if lang, ok := appLanguage.(uint); ok {
		headers["X-App-Language"] = strconv.Itoa(int(lang))
	}

	resp, err := client.Call("GET", fmt.Sprintf("/api/v1/gateway/profiles/owner-code/teacher/%s", ownerID), nil, headers)
	if err != nil {
		return "", err
	}

	var gwResp response.APIGateWayResponse[string]
	if err := json.Unmarshal(resp, &gwResp); err != nil {
		return "", fmt.Errorf("unmarshal response fail: %w", err)
	}

	if gwResp.StatusCode != 200 {
		return "", fmt.Errorf("call gateway get teacher code fail: %s", gwResp.Message)
	}

	pg.cachingProfileGateway.SetTeacherCode(ctx, ownerID, gwResp.Data)

	return gwResp.Data, nil
}

func (pg *profileGateway) GetStaffCode(ctx *gin.Context, ownerID string) (string, error) {

	cachedData, _ := pg.cachedProfileGateway.GetStaffCode(ctx, ownerID)
	if cachedData != "" {
		return cachedData, nil
	}

	// Lấy token từ context (được set ở SecuredMiddleware)
	token, exists := ctx.Get("token")
	if !exists {
		return "", fmt.Errorf("token not found in context")
	}

	tokenStr, ok := token.(string)
	if !ok || tokenStr == "" {
		return "", fmt.Errorf("invalid token in context")
	}

	client, err := NewGatewayClient(pg.serviceName, tokenStr, pg.consul, nil)
	if err != nil {
		return "", err
	}

	appLanguage, _ := ctx.Get("app_language")

	headers := make(map[string]string)
	if lang, ok := appLanguage.(uint); ok {
		headers["X-App-Language"] = strconv.Itoa(int(lang))
	}

	resp, err := client.Call("GET", fmt.Sprintf("/api/v1/gateway/profiles/owner-code/staff/%s", ownerID), nil, headers)
	if err != nil {
		return "", err
	}

	var gwResp response.APIGateWayResponse[string]
	if err := json.Unmarshal(resp, &gwResp); err != nil {
		return "", fmt.Errorf("unmarshal response fail: %w", err)
	}

	if gwResp.StatusCode != 200 {
		return "", fmt.Errorf("call gateway get staff code fail: %s", gwResp.Message)
	}

	pg.cachingProfileGateway.SetStaffCode(ctx, ownerID, gwResp.Data)

	return gwResp.Data, nil
}

func (pg *profileGateway) GetParentCode(ctx *gin.Context, ownerID string) (string, error) {
	cachedData, _ := pg.cachedProfileGateway.GetParentCode(ctx, ownerID)
	if cachedData != "" {
		return cachedData, nil
	}

	// Lấy token từ context (được set ở SecuredMiddleware)
	token, exists := ctx.Get("token")
	if !exists {
		return "", fmt.Errorf("token not found in context")
	}
	tokenStr, ok := token.(string)
	if !ok || tokenStr == "" {
		return "", fmt.Errorf("invalid token in context")
	}

	client, err := NewGatewayClient(pg.serviceName, tokenStr, pg.consul, nil)
	if err != nil {
		return "", err
	}

	appLanguage, _ := ctx.Get("app_language")

	headers := make(map[string]string)
	if lang, ok := appLanguage.(uint); ok {
		headers["X-App-Language"] = strconv.Itoa(int(lang))
	}
	resp, err := client.Call("GET", fmt.Sprintf("/api/v1/gateway/profiles/owner-code/parent/%s", ownerID), nil, headers)
	if err != nil {
		return "", err
	}

	var gwResp response.APIGateWayResponse[string]
	if err := json.Unmarshal(resp, &gwResp); err != nil {
		return "", fmt.Errorf("unmarshal response fail: %w", err)
	}

	if gwResp.StatusCode != 200 {
		return "", fmt.Errorf("call gateway get parent code fail: %s", gwResp.Message)
	}

	pg.cachingProfileGateway.SetParentCode(ctx, ownerID, gwResp.Data)

	return gwResp.Data, nil
}

func (pg *profileGateway) GetUserCode(ctx *gin.Context, ownerID string) (string, error) {
	cachedData, _ := pg.cachedProfileGateway.GetUserCode(ctx, ownerID)
	if cachedData != "" {
		return cachedData, nil
	}

	// Lấy token từ context (được set ở SecuredMiddleware)
	token, exists := ctx.Get("token")
	if !exists {
		return "", fmt.Errorf("token not found in context")
	}

	tokenStr, ok := token.(string)
	if !ok || tokenStr == "" {
		return "", fmt.Errorf("invalid token in context")
	}

	client, err := NewGatewayClient(pg.serviceName, tokenStr, pg.consul, nil)
	if err != nil {
		return "", err
	}

	appLanguage, _ := ctx.Get("app_language")

	headers := make(map[string]string)
	if lang, ok := appLanguage.(uint); ok {
		headers["X-App-Language"] = strconv.Itoa(int(lang))
	}

	resp, err := client.Call("GET", fmt.Sprintf("/api/v1/gateway/profiles/owner-code/user/%s", ownerID), nil, headers)
	if err != nil {
		return "", err
	}

	var gwResp response.APIGateWayResponse[string]
	if err := json.Unmarshal(resp, &gwResp); err != nil {
		return "", fmt.Errorf("unmarshal response fail: %w", err)
	}

	if gwResp.StatusCode != 200 {
		return "", fmt.Errorf("call gateway get user code fail: %s", gwResp.Message)
	}

	pg.cachingProfileGateway.SetUserCode(ctx, ownerID, gwResp.Data)

	return gwResp.Data, nil
}

func (pg *profileGateway) GetChildCode(ctx *gin.Context, ownerID string) (string, error) {
	cachedData, _ := pg.cachedProfileGateway.GetChildCode(ctx, ownerID)
	if cachedData != "" {
		return cachedData, nil
	}

	// Lấy token từ context (được set ở SecuredMiddleware)
	token, exists := ctx.Get("token")
	if !exists {
		return "", fmt.Errorf("token not found in context")
	}

	tokenStr, ok := token.(string)
	if !ok || tokenStr == "" {
		return "", fmt.Errorf("invalid token in context")
	}

	client, err := NewGatewayClient(pg.serviceName, tokenStr, pg.consul, nil)
	if err != nil {
		return "", err
	}

	appLanguage, _ := ctx.Get("app_language")

	headers := make(map[string]string)
	if lang, ok := appLanguage.(uint); ok {
		headers["X-App-Language"] = strconv.Itoa(int(lang))
	}

	resp, err := client.Call("GET", fmt.Sprintf("/api/v1/gateway/profiles/owner-code/child/%s", ownerID), nil, headers)
	if err != nil {
		return "", err
	}

	var gwResp response.APIGateWayResponse[string]
	if err := json.Unmarshal(resp, &gwResp); err != nil {
		return "", fmt.Errorf("unmarshal response fail: %w", err)
	}

	if gwResp.StatusCode != 200 {
		return "", fmt.Errorf("call gateway get child code fail: %s", gwResp.Message)
	}

	pg.cachingProfileGateway.SetChildCode(ctx, ownerID, gwResp.Data)

	return gwResp.Data, nil
}
