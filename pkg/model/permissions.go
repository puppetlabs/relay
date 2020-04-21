package model

type PermissionGrantAccess string

const (
	PermissionGrantAccessAll  PermissionGrantAccess = "all"
	PermissionGrantAccessSome PermissionGrantAccess = "some"
	PermissionGrantAccessNone PermissionGrantAccess = "none"
)

type PermissionName string

const (
	PermissionNameWorkflowCreate  PermissionName = "workflow.create"
	PermissionNameWorkflowDelete  PermissionName = "workflow.delete"
	PermissionNameWorkflowEdit    PermissionName = "workflow.edit"
	PermissionNameWorkflowRun     PermissionName = "workflow.run"
	PermissionNameWorkflowView    PermissionName = "workflow.view"
	PermissionNAmeWorkflowApprove PermissionName = "workflow.approve"
)

type PermissionSummary struct {
	Name PermissionName `json:"name"`
}

type PermissionGrant struct {
	Grant      PermissionGrantAccess `json:"granted"`
	Permission PermissionSummary     `json:"permission"`
}

type EntityAccess struct {
	PermissionGrants *[]PermissionGrant `json:"permission_grants"`
}
