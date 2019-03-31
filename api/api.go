package api

type Repository interface{}

type api struct {
	repo Repository
}

func New(r Repository) *api {
	return &api{
		repo: r,
	}
}
