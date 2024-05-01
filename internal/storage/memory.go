package storage

type MemRepo struct {
	storage map[string]Metric
}

func NewMemRepo() *MemRepo {
	return &MemRepo{storage: make(map[string]Metric)}
}

func (repo *MemRepo) Get(pk string) (Metric, bool) {
	m, ok := repo.storage[pk]
	return m, ok
}

func (repo *MemRepo) Save(pk string, metric Metric) error {
	repo.storage[pk] = metric
	return nil
}
