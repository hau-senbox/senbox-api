package gateway

import (
	"encoding/json"
	"fmt"
	"sen-global-api/pkg/consulapi/gateway/dto"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/hashicorp/consul/api"
)

type DepartmentGateway interface {
	GetDepartmentsByUser(context *gin.Context) ([]*dto.DepartmentGateway, error)
}

type departmentGateway struct {
	serviceName string
	consul      *api.Client
}

func NewDepartmentGateway(serviceName string, consulClient *api.Client) DepartmentGateway {
	return &departmentGateway{
		serviceName: serviceName,
		consul:      consulClient,
	}
}

func (dg *departmentGateway) GetDepartmentsByUser(context *gin.Context) ([]*dto.DepartmentGateway, error) {
	// Lấy token từ context (được set ở SecuredMiddleware)
	token, exists := context.Get("token")
	if !exists {
		return nil, fmt.Errorf("token not found in context")
	}

	tokenStr, ok := token.(string)
	if !ok || tokenStr == "" {
		return nil, fmt.Errorf("invalid token in context")
	}

	client, err := NewGatewayClient(dg.serviceName, tokenStr, dg.consul, nil)
	if err != nil {
		return nil, err
	}

	appLanguage, _ := context.Get("app_language")

	headers := make(map[string]string)
	if lang, ok := appLanguage.(uint); ok {
		headers["X-App-Language"] = strconv.Itoa(int(lang))
	}

	resp, err := client.Call("GET", "/api/v1/gateway/departments", nil, headers)
	if err != nil {
		return nil, err
	}

	var gwResp dto.APIGateWayResponse[[]*dto.DepartmentGateway]
	if err := json.Unmarshal(resp, &gwResp); err != nil {
		return nil, fmt.Errorf("unmarshal response fail: %w", err)
	}

	if gwResp.StatusCode != 200 {
		return nil, fmt.Errorf("call gateway get departments by user fail: %s", gwResp.Message)
	}

	return gwResp.Data, nil
}

func (dg *departmentGateway) GetDepartmentsByOrganization(context *gin.Context, orgID string) ([]*dto.DepartmentGateway, error) {
	token, exists := context.Get("token")
	if !exists {
		return nil, fmt.Errorf("token not found in context")
	}

	if !exists {
		return nil, fmt.Errorf("token not found in context")
	}

	tokenStr, ok := token.(string)
	if !ok || tokenStr == "" {
		return nil, fmt.Errorf("invalid token in context")
	}

	client, err := NewGatewayClient(dg.serviceName, tokenStr, dg.consul, nil)
	if err != nil {
		return nil, err
	}

	appLanguage, _ := context.Get("app_language")

	headers := make(map[string]string)
	if lang, ok := appLanguage.(uint); ok {
		headers["X-App-Language"] = strconv.Itoa(int(lang))
	}

	resp, err := client.Call("GET", fmt.Sprintf("/api/v1/gateway/departments/organization/%s", orgID), nil, headers)
	if err != nil {
		return nil, err
	}

	var gwResp dto.APIGateWayResponse[[]*dto.DepartmentGateway]
	if err := json.Unmarshal(resp, &gwResp); err != nil {
		return nil, fmt.Errorf("unmarshal response fail: %w", err)
	}

	if gwResp.StatusCode != 200 {
		return nil, fmt.Errorf("call gateway gett departments by organization fail: %s", gwResp.Message)
	}

	return gwResp.Data, nil

}
