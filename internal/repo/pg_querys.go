package repo

const (
	InitQuery = `CREATE TABLE IF NOT EXISTS users (
						id SERIAL PRIMARY KEY,
						username TEXT UNIQUE NOT NULL,
						password TEXT NOT NULL, -- для упрощения, в продакшн разработке пароли в открытом виде не хранятся, хранится хеш или зашифрованные паро
						created_at TIMESTAMP DEFAULT now()
				);

				CREATE TABLE IF NOT EXISTS tasks (
						id SERIAL PRIMARY KEY,
						user_id INT REFERENCES users(id) ON DELETE CASCADE,
						title TEXT NOT NULL,
						description TEXT DEFAULT 'new',
						status TEXT CHECK (status IN ('new', 'in_progress', 'done')) DEFAULT 'new',
						created_at TIMESTAMP DEFAULT now()
				);

				CREATE INDEX IF NOT EXISTS idx_tasks_user_id ON tasks(user_id);
`

	GetAllTasksQuery           = `SELECT * FROM tasks;`
	GetTaskByIDQuery           = `SELECT * FROM tasks WHERE id = $1;`
	GetAllTasksByUserIdQuery   = `SELECT * FROM tasks WHERE user_id = $1;`
	GetLastTaskByUserIdQuery   = `SELECT * FROM tasks WHERE user_id = $1 limit 1;`
	GetAllTasksByUserNameQuery = `SELECT t.* FROM user AS u LEFT JOIN tasks AS t ON t.user_id = u.id WHERE u.username = $1;`

	CreateTaskQuery = `INSERT INTO tasks (user_id, title, description) SELECT $1, $2, $3 
					   WHERE EXISTS (SELECT 1 FROM users WHERE id = $1) RETURNING id;`

	UpdateTaskStatusByIDQuery = `UPDATE tasks SET status = $1 WHERE user_id = $2;`

	DeleteTaskByIDQuery = `DELETE FROM tasks WHERE id = $1;`

	CreateUserQuery = `INSERT INTO users (username, password) VALUES ($1, $2) RETURNING id;`
)
