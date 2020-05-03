package core

// Elevated permissions
const (
	PermissionResetUserPassword   = "RESET_USER_PASSWORD"
	PermissionSetAdmin            = "SET_ADMIN"
	PermissionProxyAsUser         = "PROXY_AS_USER"
	PermissionViewAppSettings     = "VIEW_APP_SETTINGS"
	PermissionChangeAppSettings   = "CHANGE_APP_SETTINGS"
	PermissionManageAppDeployment = "MANAGE_APP_DEPLOYMENT"
	PermissionStopServer          = "STOP_SERVER"
	PermissionManageIndexes       = "MANAGE_INDEXES"
	PermissionLogging             = "LOGGING"
)

const (
	PermissionPostComments     = "POST_COMMENTS"
	PermissionEditUserInfo     = "EDIT_USER_INFO"
	PermissionUnlockUser       = "UNLOCK_USER"
	PermissionViewUsers        = "VIEW_USERS"
	PermissionEditGroups       = "EDIT_GROUPS"
	PermissionViewGroups       = "VIEW_GROUPS"
	PermissionManageIOTDevices = "MANAGE_IOT_DEVICES"
)

// Visit module
const (
	PermissionSetDefaultParty = "SET_DEFAULT_PARTY"
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
	PermissionChangeAppSettings,
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
	PermissionManageIOTDevices,

	PermissionSetDefaultParty,

	PermissionListProjects,
}
