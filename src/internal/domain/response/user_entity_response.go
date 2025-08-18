package response

import "time"

type UserEntityResponse struct {
	ID           string   `json:"id"`
	Username     string   `json:"username"`
	Nickname     string   `json:"nickname"`
	Fullname     string   `json:"fullname"`
	Phone        string   `json:"phone"`
	Email        string   `json:"email"`
	Dob          string   `json:"dob"`
	QRLogin      string   `json:"qr_login"`
	Avatar       string   `json:"avatar"`
	AvatarURL    string   `json:"avatar_url"`
	IsBlocked    bool     `json:"is_blocked"`
	BlockedAt    string   `json:"blocked_at"`
	Organization []string `json:"organizations"`
	CreatedAt    string   `json:"created_at"`
	CustomID     string   `json:"custom_id"`

	Roles   *[]RoleListResponseData `json:"roles"`
	Devices *[]string               `json:"devices"`

	UserOrganizationActive UserOrganizationActive `json:"user_organization_active"`
}

type UserEntityResponseV2 struct {
	ID                   string   `json:"id"`
	Username             string   `json:"username"`
	Nickname             string   `json:"nickname"`
	Fullname             string   `json:"fullname"`
	Phone                string   `json:"phone"`
	Email                string   `json:"email"`
	Dob                  string   `json:"dob"`
	QRLogin              string   `json:"qr_login"`
	Avatar               string   `json:"avatar"`
	AvatarURL            string   `json:"avatar_url"`
	IsBlocked            bool     `json:"is_blocked"`
	BlockedAt            string   `json:"blocked_at"`
	Organization         []string `json:"organizations"`
	CreatedAt            string   `json:"created_at"`
	IsDeactive           bool     `json:"is_deactive"`
	IsSuperAdmin         bool     `json:"is_super_admin"`
	OrganizationIdActive string   `json:"organization_id_active"`

	Roles   *[]RoleListResponseData `json:"roles"`
	Devices *[]string               `json:"devices"`

	OrganizationAdmin *OrganizationAdmin `json:"organization_admin"`
}

type UserEntityResponseData struct {
	ID        string   `json:"id"`
	Username  string   `json:"username"`
	Nickname  string   `json:"nickname"`
	Avatar    string   `json:"avatar"`
	AvatarURL string   `json:"avatar_url"`
	Roles     []string `json:"roles"`
}

type OrganizationAdmin struct {
	ID               string    `json:"id"`
	OrganizationName string    `json:"organization_name"`
	Avatar           string    `json:"avatar"`
	AvatarURL        string    `json:"avatar_url"`
	Address          string    `json:"address"`
	Description      string    `json:"description"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

type OrganizationActive struct {
	ID               string    `json:"id"`
	OrganizationName string    `json:"organization_name"`
	Avatar           string    `json:"avatar"`
	AvatarURL        string    `json:"avatar_url"`
	CreatedAt        time.Time `json:"created_at"`
}

type UserOrganizationActive struct {
	DefaultOrganization []OrganizationActive `json:"default_organization"`
	TeacherOrganization []OrganizationActive `json:"teacher_organization"`
	StaffOrganization   []OrganizationActive `json:"staff_organization"`
}

// type OrganizationIdActive struct {
// 	OrgTeacher
// }
