package types

import "encoding/json"

type SubspacesMsgRequest struct {
	CreateSubspace          json.RawMessage `json:"create_subspace"`
	EditSubspace            json.RawMessage `json:"edit_subspace"`
	DeleteSubspace          json.RawMessage `json:"delete_subspace"`
	CreateUserGroup         json.RawMessage `json:"create_user_group"`
	EditUserGroup           json.RawMessage `json:"edit_user_group"`
	SetUserGroupPermissions json.RawMessage `json:"set_user_group_permissions"`
	DeleteUserGroup         json.RawMessage `json:"delete_user_group"`
	AddUserToUserGroup      json.RawMessage `json:"add_user_to_user_group"`
	RemoveUserFromUserGroup json.RawMessage `json:"remove_user_from_user_group"`
	SetUserPermissions      json.RawMessage `json:"set_user_permissions"`
}

type SubspacesQueryRequest struct {
	Subspaces        json.RawMessage `json:"subspaces"`
	Subspace         json.RawMessage `json:"subspace"`
	UserGroups       json.RawMessage `json:"user_groups"`
	UserGroup        json.RawMessage `json:"user_group"`
	UserGroupMembers json.RawMessage `json:"user_group_members"`
	UserPermissions  json.RawMessage `json:"user_permissions"`
}
