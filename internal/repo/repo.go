package repo

import (
	"TemplatestPGSQL/internal/config"
	"TemplatestPGSQL/internal/dto"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
)

type Repository interface {
	CreateTask(ctx context.Context, task Task) error
	GetAllTasks(ctx context.Context) ([]Task, error)
	GetTaskByID(ctx context.Context, id string) (*Task, error)
	GetLastTaskByUserID(ctx context.Context, id string) (*Task, error)
	GetTasksByUserName(ctx context.Context, name string) ([]Task, error)
	GetAllTasksByUserID(ctx context.Context, id string) ([]Task, error)
	UpdateStatusByID(ctx context.Context, id string, status string) error
	DeleteTaskByID(ctx context.Context, id string) error

	CreateUser(ctx context.Context, user User) error
}

type repository struct {
	pool *pgxpool.Pool
}

func NewRepository(ctx context.Context, cfg config.Memory) (*repository, error) {
	connString := fmt.Sprintf(
		`user=%s password=%s host=%s port=%d dbname=%s sslmode=%s 
        pool_max_conns=%d pool_max_conn_lifetime=%s pool_max_conn_idle_time=%s`,
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Name,
		cfg.SSLMode,
		cfg.PoolMaxConns,
		cfg.PoolMaxConnLifetime.String(),
		cfg.PoolMaxConnIdleTime.String(),
	)

	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse DB config")
	}

	config.ConnConfig.DefaultQueryExecMode = pgx.QueryExecModeCacheDescribe

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create DB connection pool")
	}

	return &repository{pool}, nil
}

func (r *repository) InitTables(ctx context.Context) error {
	_, err := r.pool.Exec(ctx, InitQuery)
	if err != nil {
		return errors.Wrap(err, "failed to initialise tables")
	}
	return nil
}

func (r *repository) GetAllTasks(ctx context.Context) ([]Task, error) {
	pgRows, err := r.pool.Query(ctx, GetAllTasksQuery)
	if err != nil {
		return nil, errors.Wrap(err, "failed to query all tasks")
	}

	defer pgRows.Close()
	tasks, err := pgx.CollectRows(pgRows, pgx.RowToStructByName[Task])
	if err != nil {
		return nil, errors.Wrap(err, "failed to convert all tasks")
	}

	return tasks, nil
}

func (r *repository) GetTaskByID(ctx context.Context, id string) (*Task, error) {
	pgRow, err := r.pool.Query(ctx, GetTaskByIdQuery, id)
	if err != nil {
		return nil, errors.Wrap(err, "failed to query task")
	}

	defer pgRow.Close()
	task, err := pgx.CollectOneRow(pgRow, pgx.RowToStructByName[Task])
	if err != nil {
		return nil, errors.Wrap(err, "failed to convert task")
	}

	return &task, nil
}

func (r *repository) GetLastTaskByUserID(ctx context.Context, id string) (*Task, error) {
	user, err := r.getUserByName(ctx, id)
	if err != nil {
		return nil, errors.Wrap(err, "failed to check user exist")
	}
	if user == nil {
		return nil, dto.ErrNotFound
	}

	pgRow, err := r.pool.Query(ctx, GetLastTaskByUserIdQuery, id)
	if err != nil {
		return nil, errors.Wrap(err, "failed to query task")
	}

	defer pgRow.Close()
	task, err := pgx.CollectOneRow(pgRow, pgx.RowToStructByName[Task])
	if err != nil {
		return nil, errors.Wrap(err, "failed to convert task")
	}

	return &task, nil
}

func (r *repository) GetAllTasksByUserID(ctx context.Context, id string) ([]Task, error) {
	user, err := r.getUserByID(ctx, id)
	if err != nil {
		return nil, errors.Wrap(err, "failed to check user exist")
	}
	if user == nil {
		return nil, dto.ErrNotFound
	}

	pgRows, err := r.pool.Query(ctx, GetAllTasksByUserIdQuery, id)
	if err != nil {
		return nil, errors.Wrap(err, "failed to query task")
	}

	defer pgRows.Close()
	tasks, err := pgx.CollectRows(pgRows, pgx.RowToStructByName[Task])
	if err != nil {
		return nil, errors.Wrap(err, "failed to convert task")
	}

	return tasks, nil
}

func (r *repository) GetTasksByUserName(ctx context.Context, name string) ([]Task, error) {
	user, err := r.getUserByName(ctx, name)
	if err != nil {
		return nil, errors.Wrap(err, "failed to check user exist")
	}
	if user == nil {
		return nil, dto.ErrNotFound
	}

	pgRows, err := r.pool.Query(ctx, GetAllTasksByUserNameQuery, name)
	if err != nil {
		return nil, errors.Wrap(err, "failed to query task")
	}
	defer pgRows.Close()
	tasks, err := pgx.CollectRows(pgRows, pgx.RowToStructByName[Task])
	if err != nil {
		return nil, errors.Wrap(err, "failed to convert task")
	}

	return tasks, nil
}

// status MUST BE IN ('new', 'in_progress', 'done')
func (r *repository) UpdateStatusByID(ctx context.Context, id string, status string) error {
	cmdTag, err := r.pool.Exec(ctx, UpdateTaskStatusByIDQuery, id, status)
	if err != nil {
		return errors.Wrap(err, "failed to query task")
	}
	if cmdTag.RowsAffected() == 0 {
		return dto.ErrNotFound
	}
	return nil
}

func (r *repository) DeleteTaskByID(ctx context.Context, id string) error {
	cmdTag, err := r.pool.Exec(ctx, DeleteTaskByIdQuery, id)
	if err != nil {
		return errors.Wrap(err, "failed to delete task")
	}

	if cmdTag.RowsAffected() == 0 {
		return dto.ErrNotFound
	}

	return nil
}

func (r *repository) CreateUser(ctx context.Context, user User) error {
	_, err := r.pool.Exec(ctx, CreateUserQuery, user.Name, user.Password)
	if err != nil {
		return errors.Wrap(err, "failed to create user")
	}
	return nil
}

func (r *repository) CreateTask(ctx context.Context, task Task) error {
	_, err := r.pool.Exec(ctx, CreateTaskQuery, task.UserID, task.Title, task.Data)
	if err != nil {
		return errors.Wrap(err, "failed to create task")
	}
	return nil
}

func (r *repository) getUserByID(ctx context.Context, id string) (*User, error) {
	pgRow, err := r.pool.Query(ctx, GetUserByIdQuery, id)
	if err != nil {
		return nil, errors.Wrap(err, "failed to query user")
	}

	defer pgRow.Close()
	User, err := pgx.CollectOneRow(pgRow, pgx.RowToStructByName[User])
	if err != nil {
		return nil, errors.Wrap(err, "failed to convert user")
	}

	return &User, nil
}

func (r *repository) getUserByName(ctx context.Context, name string) (*User, error) {
	pgRow, err := r.pool.Query(ctx, GetUserByNameQuery, name)
	if err != nil {
		return nil, errors.Wrap(err, "failed to query user")
	}

	defer pgRow.Close()
	User, err := pgx.CollectOneRow(pgRow, pgx.RowToStructByName[User])
	if err != nil {
		return nil, errors.Wrap(err, "failed to convert user")
	}

	return &User, nil
}
