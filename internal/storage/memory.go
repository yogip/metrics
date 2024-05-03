package storage

type MemRepo struct {
	storage map[string]interface{}
}

func NewMemRepo() *MemRepo {
	return &MemRepo{storage: make(map[string]interface{})}
}

func (repo *MemRepo) Get(pk string) (interface{}, bool) {
	m, ok := repo.storage[pk]
	return m, ok
}

func (repo *MemRepo) Save(pk string, metric interface{}) error {
	repo.storage[pk] = metric
	return nil
}
