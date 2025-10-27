package gateway

import (
	"encoding/json"
	"fmt"
	"sen-global-api/pkg/consulapi/gateway/dto/request"
	"sen-global-api/pkg/consulapi/gateway/dto/response"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/hashicorp/consul/api"
)

type ProfileGateway interface {
	GenerateStudentCode(ctx *gin.Context, ownerID string, createdIndex int) (*string, error)
	GenerateTeacherCode(ctx *gin.Context, ownerID string, createdIndex int) (*string, error)
	GenerateStaffCode(ctx *gin.Context, ownerID string, createdIndex int) (*string, error)
	GenerateParentCode(ctx *gin.Context, ownerID string, createdIndex int) (*string, error)
	GenerateUserCode(ctx *gin.Context, ownerID string, createdIndex int) (*string, error)
	GenerateChildCode(ctx *gin.Context, ownerID string, createdIndex int) (*string, error)
}

type profileGateway struct {
	serviceName string
	consul      *api.Client
}

func NewProfileGateway(serviceName string, consulClient *api.Client) ProfileGateway {
	return &profileGateway{
		serviceName: serviceName,
		consul:      consulClient,
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

	resp, err := client.Call("POST", "/api/v1/admin/profiles/owner-code/student/generate", req, headers)
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

	resp, err := client.Call("POST", "/api/v1/admin/profiles/owner-code/teacher/generate", req, headers)
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

	resp, err := client.Call("POST", "/api/v1/admin/profiles/owner-code/staff/generate", req, headers)
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

	resp, err := client.Call("POST", "/api/v1/admin/profiles/owner-code/parent/generate", req, headers)
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

	resp, err := client.Call("POST", "/api/v1/admin/profiles/owner-code/user/generate", req, headers)
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

	resp, err := client.Call("POST", "/api/v1/admin/profiles/owner-code/child/generate", req, headers)
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
