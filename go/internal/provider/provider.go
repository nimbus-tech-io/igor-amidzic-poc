package provider

import (
	"goapp/internal/repo"
)

type AppProvider struct {
    repository repo.Repository
}

func NewAppProvider() *AppProvider {
    return &AppProvider{}
}

func (p *AppProvider) Repository() repo.Repository {
    if p.repository == nil {
        p.repository = repo.NewNeonRepo()
    }
    return p.repository
}
