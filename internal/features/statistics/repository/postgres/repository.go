package statistics_postgres_repository

import core_postgres_pool "github.com/AlexeyStrekozov/effective_mobile/internal/core/repositorty/postgres/pool"

type StatisticsRepository struct {
	pool core_postgres_pool.Pool
}

func NewStatisticsRepository(
	pool core_postgres_pool.Pool,
) *StatisticsRepository {
	return &StatisticsRepository{
		pool: pool,
	}
}
