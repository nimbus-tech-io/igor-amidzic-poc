package provider

import (
	"goapp/internal/repo"
)

func Repository() repo.Repository {
	return repo.NewNeonRepo()
}

