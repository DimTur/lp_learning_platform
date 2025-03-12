package ssomodels

type CreateLearningGroup struct {
	Name        string   `json:"name" validate:"required,min=3,max=100"`
	CreatedBy   string   `json:"created_by" validate:"required"`
	ModifiedBy  string   `json:"modified_by" validate:"required"`
	GroupAdmins []string `json:"group_admins" validate:"required"`
	Learners    []string `json:"learners" validate:"required"`
}

type CreateLearningGroupResp struct {
	Success bool
}

type GetLgByID struct {
	UserID string `json:"user_id" validate:"required"`
	LgId   string `json:"learning_group_id" validate:"required"`
}

type GetLgByIDResp struct {
	Id          string
	Name        string
	CreatedBy   string
	ModifiedBy  string
	Learners    []*Learner
	GroupAdmins []*GroupAdmins
}

type Learner struct {
	Id    string
	Email string
	Name  string
}

type GroupAdmins struct {
	Id    string
	Email string
	Name  string
}

type UpdateLearningGroup struct {
	UserID      string   `json:"user_id" validate:"required"`
	LgId        string   `json:"learning_group_id" validate:"required"`
	Name        string   `json:"name,omitempty"`
	ModifiedBy  string   `json:"modified_by,omitempty"`
	GroupAdmins []string `json:"group_admins,omitempty"`
	Learners    []string `json:"learners,omitempty"`
}

type UpdateLearningGroupResp struct {
	Success bool
}

type DelLgByID struct {
	UserID string `json:"user_id" validate:"required"`
	LgID   string `json:"learning_group_id" validate:"required"`
}

type DelLgByIDResp struct {
	Success bool
}

type GetLGroups struct {
	UserID string `json:"user_id" validate:"required"`
}

type GetLGroupsResp struct {
	LearningGroups []*LearningGroup
}

type LearningGroup struct {
	Id         string
	Name       string
	CreatedBy  string
	ModifiedBy string
	Created    string
	Updated    string
}

type IsGroupAdmin struct {
	UserID string `json:"user_id" validate:"required"`
	LgID   string `json:"learning_group_id" validate:"required"`
}

type IsGroupAdminResp struct {
	IsGroupAdmin bool
}

type UserIsGroupAdminIn struct {
	UserID string `json:"user_id" validate:"required"`
}

type UserIsLearnerIn struct {
	UserID string `json:"user_id" validate:"required"`
}

type GetLearners struct {
	Learners []string `json:"learners"`
}
