package gateway

import (
	"encoding/json"
	"fmt"
	"sen-global-api/pkg/consulapi/gateway/dto/response"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/hashicorp/consul/api"
)

type DepartmentGateway interface {
	GetDepartmentsByUser(context *gin.Context) ([]*response.DepartmentGateway, error)
	GetDepartmentsByOrganization(context *gin.Context, orgID string) ([]*response.DepartmentGateway, error)
	AssignParentDepartmentGroup(context *gin.Context, parentID string, organizationID string) error
	AssignStudentDepartmentGroup(context *gin.Context, studentID string, organizationID string) error
	AssignTeacherDepartmentGroup(context *gin.Context, teacherID string, organizationID string) error
	AssignStaffDepartmentGroup(context *gin.Context, staffID string, organizationID string) error
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

func (dg *departmentGateway) GetDepartmentsByUser(context *gin.Context) ([]*response.DepartmentGateway, error) {
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

	var gwResp response.APIGateWayResponse[[]*response.DepartmentGateway]
	if err := json.Unmarshal(resp, &gwResp); err != nil {
		return nil, fmt.Errorf("unmarshal response fail: %w", err)
	}

	if gwResp.StatusCode != 200 {
		return nil, fmt.Errorf("call gateway get departments by user fail: %s", gwResp.Message)
	}

	return gwResp.Data, nil
}

func (dg *departmentGateway) GetDepartmentsByOrganization(context *gin.Context, orgID string) ([]*response.DepartmentGateway, error) {
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

	var gwResp response.APIGateWayResponse[[]*response.DepartmentGateway]
	if err := json.Unmarshal(resp, &gwResp); err != nil {
		return nil, fmt.Errorf("unmarshal response fail: %w", err)
	}

	if gwResp.StatusCode != 200 {
		return nil, fmt.Errorf("call gateway gett departments by organization fail: %s", gwResp.Message)
	}

	return gwResp.Data, nil

}

func (dg *departmentGateway) AssignParentDepartmentGroup(context *gin.Context, parentID string, organizationID string) error {
	token, exists := context.Get("token")
	if !exists {
		return fmt.Errorf("token not found in context")
	}

	tokenStr, ok := token.(string)
	if !ok || tokenStr == "" {
		return fmt.Errorf("invalid token in context")
	}

	client, err := NewGatewayClient(dg.serviceName, tokenStr, dg.consul, nil)
	if err != nil {
		return fmt.Errorf("create gateway client fail: %w", err)
	}

	appLanguage, _ := context.Get("app_language")

	headers := make(map[string]string)
	if lang, ok := appLanguage.(uint); ok {
		headers["X-App-Language"] = strconv.Itoa(int(lang))
	}

	type AssignParentDepartmentGroupRequest struct {
		ParentID       string `json:"parent_id"`
		OrganizationID string `json:"organization_id"`
	}

	req := AssignParentDepartmentGroupRequest{
		ParentID:       parentID,
		OrganizationID: organizationID,
	}

	resp, err := client.Call("POST", "/api/v1/gateway/departments/assign/parent-group", req, headers)
	if err != nil {
		return fmt.Errorf("call gateway assign parent department group fail: %w", err)
	}

	var gwResp response.APIGateWayResponse[string]
	if err := json.Unmarshal(resp, &gwResp); err != nil {
		return fmt.Errorf("unmarshal response fail: %w", err)
	}

	if gwResp.StatusCode != 200 {
		return fmt.Errorf("call gateway assign parent department group fail: %s", gwResp.Message)
	}

	return nil
}

func (dg *departmentGateway) AssignStudentDepartmentGroup(context *gin.Context, studentID string, organizationID string) error {
	token, exists := context.Get("token")
	if !exists {
		return fmt.Errorf("token not found in context")
	}

	tokenStr, ok := token.(string)
	if !ok || tokenStr == "" {
		return fmt.Errorf("invalid token in context")
	}

	client, err := NewGatewayClient(dg.serviceName, tokenStr, dg.consul, nil)
	if err != nil {
		return fmt.Errorf("create gateway client fail: %w", err)
	}

	appLanguage, _ := context.Get("app_language")

	headers := make(map[string]string)
	if lang, ok := appLanguage.(uint); ok {
		headers["X-App-Language"] = strconv.Itoa(int(lang))
	}

	type AssignStudentDepartmentGroupRequest struct {
		StudentID      string `json:"student_id"`
		OrganizationID string `json:"organization_id"`
	}

	req := AssignStudentDepartmentGroupRequest{
		StudentID:      studentID,
		OrganizationID: organizationID,
	}

	resp, err := client.Call("POST", "/api/v1/gateway/departments/assign/student-group", req, headers)
	if err != nil {
		return fmt.Errorf("call gateway assign student department group fail: %w", err)
	}

	var gwResp response.APIGateWayResponse[string]
	if err := json.Unmarshal(resp, &gwResp); err != nil {
		return fmt.Errorf("unmarshal response fail: %w", err)
	}

	if gwResp.StatusCode != 200 {
		return fmt.Errorf("call gateway assign student department group fail: %s", gwResp.Message)
	}

	return nil
}

func (dg *departmentGateway) AssignTeacherDepartmentGroup(context *gin.Context, teacherID string, organizationID string) error {
	token, exists := context.Get("token")
	if !exists {
		return fmt.Errorf("token not found in context")
	}

	tokenStr, ok := token.(string)
	if !ok || tokenStr == "" {
		return fmt.Errorf("invalid token in context")
	}

	client, err := NewGatewayClient(dg.serviceName, tokenStr, dg.consul, nil)
	if err != nil {
		return fmt.Errorf("create gateway client fail: %w", err)
	}

	appLanguage, _ := context.Get("app_language")

	headers := make(map[string]string)
	if lang, ok := appLanguage.(uint); ok {
		headers["X-App-Language"] = strconv.Itoa(int(lang))
	}

	type AssignTeacherDepartmentGroupRequest struct {
		TeacherID      string `json:"teacher_id"`
		OrganizationID string `json:"organization_id"`
	}

	req := AssignTeacherDepartmentGroupRequest{
		TeacherID:      teacherID,
		OrganizationID: organizationID,
	}

	resp, err := client.Call("POST", "/api/v1/gateway/departments/assign/teacher-group", req, headers)
	if err != nil {
		return fmt.Errorf("call gateway assign teacher department group fail: %w", err)
	}

	var gwResp response.APIGateWayResponse[string]
	if err := json.Unmarshal(resp, &gwResp); err != nil {
		return fmt.Errorf("unmarshal response fail: %w", err)
	}

	if gwResp.StatusCode != 200 {
		return fmt.Errorf("call gateway assign teacher department group fail: %s", gwResp.Message)
	}

	return nil
}

func (dg *departmentGateway) AssignStaffDepartmentGroup(context *gin.Context, staffID string, organizationID string) error {
	token, exists := context.Get("token")
	if !exists {
		return fmt.Errorf("token not found in context")
	}

	tokenStr, ok := token.(string)
	if !ok || tokenStr == "" {
		return fmt.Errorf("invalid token in context")
	}

	client, err := NewGatewayClient(dg.serviceName, tokenStr, dg.consul, nil)
	if err != nil {
		return fmt.Errorf("create gateway client fail: %w", err)
	}

	appLanguage, _ := context.Get("app_language")

	headers := make(map[string]string)
	if lang, ok := appLanguage.(uint); ok {
		headers["X-App-Language"] = strconv.Itoa(int(lang))
	}

	type AssignStaffDepartmentGroupRequest struct {
		StaffID        string `json:"staff_id"`
		OrganizationID string `json:"organization_id"`
	}

	req := AssignStaffDepartmentGroupRequest{
		StaffID:        staffID,
		OrganizationID: organizationID,
	}

	resp, err := client.Call("POST", "/api/v1/gateway/departments/assign/staff-group", req, headers)
	if err != nil {
		return fmt.Errorf("call gateway assign staff department group fail: %w", err)
	}

	var gwResp response.APIGateWayResponse[string]
	if err := json.Unmarshal(resp, &gwResp); err != nil {
		return fmt.Errorf("unmarshal response fail: %w", err)
	}

	if gwResp.StatusCode != 200 {
		return fmt.Errorf("call gateway assign teacher department group fail: %s", gwResp.Message)
	}

	return nil
}
