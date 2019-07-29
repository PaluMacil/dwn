package core

// Elevated permissions
const (
	PermissionResetUserPassword   = "RESET_USER_PASSWORD"
	PermissionSetAdmin            = "SET_ADMIN"
	PermissionProxyAsUser         = "PROXY_AS_USER"
	PermissionViewAppSettings     = "VIEW_APP_SETTINGS"
	PermissionManageAppDeployment = "MANAGE_APP_DEPLOYMENT"
	PermissionStopServer          = "STOP_SERVER"
	PermissionManageIndexes       = "MANAGE_INDEXES"
	PermissionLogging             = "LOGGING"
)

const (
	PermissionPostComments      = "POST_COMMENTS"
	PermissionEditUserInfo      = "EDIT_USER_INFO"
	PermissionUnlockUser        = "UNLOCK_USER"
	PermissionViewUsers         = "VIEW_USERS"
	PermissionEditGroups        = "EDIT_GROUPS"
	PermissionViewGroups        = "VIEW_GROUPS"
	PermissionEditAllVisitPages = "EDIT_ALL_VISIT_PAGES"
	PermissionManageIOTDevices  = "MANAGE_IOT_DEVICES"
)

// Visit module
const (
	PermissionSetDefaultParty = "SET_DEFAULT_PARTY"
)

// Rental Module
const (
	PermissionViewAllTenantInfo = "VIEW_ALL_TENANT_INFO"
	PermissionManageAllTenants  = "MANAGE_ALL_TENANTS"
)

// Projects Module
const (
	PermissionListProjects = "LIST_PROJECTS"
)

var Permissions = []string{
	PermissionResetUserPassword,
	PermissionSetAdmin,
	PermissionProxyAsUser,
	PermissionViewAppSettings,
	PermissionManageAppDeployment,
	PermissionStopServer,
	PermissionManageIndexes,
	PermissionLogging,

	PermissionPostComments,
	PermissionEditUserInfo,
	PermissionUnlockUser,
	PermissionViewUsers,
	PermissionEditGroups,
	PermissionViewGroups,
	PermissionEditAllVisitPages,
	PermissionManageIOTDevices,

	PermissionSetDefaultParty,

	PermissionViewAllTenantInfo,
	PermissionManageAllTenants,

	PermissionListProjects,
}