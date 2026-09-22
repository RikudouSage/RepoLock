package manager

type RepositoryPermission interface {
}

func NewRepositoryPermissionManager() RepositoryPermission {
	return &repositoryPermission{}
}

type repositoryPermission struct{}
