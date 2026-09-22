package permissions

import (
	"go.chrastecky.dev/repolock/config/data"
	"go.chrastecky.dev/repolock/http/response"
	"go.chrastecky.dev/repolock/manager"
)

type Controller struct {
	orgPermissionManager  manager.OrganizationMembership
	repoPermissionManager manager.RepositoryPermission
	accessManager         manager.Access
	writer                response.Writer
	orgManager            manager.Organization
	repoManager           manager.Repository
	anyoneCanRequestJoin  bool
}

func NewController(
	orgPermissionManager manager.OrganizationMembership,
	repoPermissionManager manager.RepositoryPermission,
	orgManager manager.Organization,
	repoManager manager.Repository,
	accessManager manager.Access,
	writer response.Writer,
	config *data.GlobalConfig,
) *Controller {
	return &Controller{
		orgPermissionManager:  orgPermissionManager,
		repoPermissionManager: repoPermissionManager,
		accessManager:         accessManager,
		writer:                writer,
		orgManager:            orgManager,
		repoManager:           repoManager,
		anyoneCanRequestJoin:  config.AnyoneCanRequestJoin,
	}
}
