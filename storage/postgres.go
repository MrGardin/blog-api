package storage

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Pool struct {
	pool *pgxpool.Pool
}

func (p *Pool) GetAllUsers(ctx context.Context) ([]User, error) {
	var users []User
	rows, err := p.pool.Query(ctx, "SELECT id, first_name, second_name, email, created_at FROM users ORDER BY id")
	if err != nil {
		return nil, fmt.Errorf("error querying users: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var user User
		//Scan data
		err = rows.Scan(&user.Id, &user.FirstName, &user.SecondName, &user.Email, &user.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("error scan rows")
		}
		users = append(users, user)
	}
	//Check errors rows
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating users: %w", err)
	}
	return users, nil
}
func (p *Pool) GetUserById(ctx context.Context, id int) (*User, error) {
	var user User
	err := p.pool.QueryRow(ctx,
		"SELECT id, first_name, second_name, email, created_at FROM users WHERE id = $1", id).Scan(&user.Id, &user.FirstName, &user.SecondName, &user.Email, &user.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("error getUserById: %w", err)
	}
	return &user, nil

}
func NewPool(ctx context.Context, connString string) (*Pool, error) {
	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, fmt.Errorf("error with config of pool: %w", err)
	}
	config.MaxConns = 20
	config.MinConns = 5

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("error created pool: %w", err)
	}

	err = pool.Ping(ctx)
	if err != nil {
		return nil, fmt.Errorf("error ping of database: %w", err)
	}
	return &Pool{pool}, nil
}

func (p *Pool) Close() {
	if p.pool != nil {
		p.pool.Close()
	}
}

func (p *Pool) GetVersion(ctx context.Context) (string, error) {
	var version string
	err := p.pool.QueryRow(ctx, "SELECT version()").Scan(&version)
	if err != nil {
		return "", fmt.Errorf("error with select version database %w", err)
	}
	return version, nil
}
func (p *Pool) CreateNewUser(ctx context.Context, firstName, secondName, email string) (int, error) {
	var id int
	err := p.pool.QueryRow(ctx, "INSERT INTO users (first_name, second_name, email) VALUES ($1, $2, $3) RETURNING id", firstName, secondName, email).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("error with createNewUser %w", err)
	}
	return id, nil
}
func (p *Pool) UpdateUser(ctx context.Context, id int, firstName, secondName, email string) (int, error) {
	_, err := p.pool.Exec(ctx, "UPDATE users SET first_name = $1, second_name = $2, email = $3 WHERE id = $4", firstName, secondName, email, id)
	if err != nil {
		return 0, fmt.Errorf("error with update user %w", err)
	}
	return id, nil
}
func (p *Pool) DeleteUser(ctx context.Context, id int) error {
	_, err := p.pool.Exec(ctx, "DELETE FROM users WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("error with delete user %w", err)
	}
	return nil
}
func (p *Pool) CreateNewPost(ctx context.Context, title, content string, userId int) (int, error) {
	var id int
	err := p.pool.QueryRow(ctx, "INSERT INTO posts (title, content, user_id) VALUES($1, $2, $3) RETURNING id", title, content, userId).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("error with create post %w", err)
	}
	return id, nil
}

// Получить все посты
func (p *Pool) GetAllPosts(ctx context.Context) ([]Post, error) {
	var posts []Post
	rows, err := p.pool.Query(ctx, "SELECT id, title, content, user_id, created_at, updated_at FROM posts ORDER BY id")
	if err != nil {
		return nil, fmt.Errorf("error with get rows in GetAllPosts: %w", err)
	}
	var post Post
	for rows.Next() {
		err = rows.Scan(&post.ID, &post.Title, &post.Content, &post.UserID, &post.CreatedAt, &post.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("error with scan post: %w", err)
		}
		posts = append(posts, post)
	}
	return posts, nil
}

// Получить пост по id
func (p *Pool) GetPostById(ctx context.Context, id int) (*Post, error) {
	var post Post
	err := p.pool.QueryRow(ctx, "SELECT id, title, content, user_id, created_at, updated_at FROM posts WHERE id = $1", id).Scan(&post.ID, &post.Title, &post.Content, &post.UserID, &post.CreatedAt, &post.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("error with get post: %w", err)
	}
	return &post, nil
}
func (p *Pool) GetUsersPost(ctx context.Context, UserID int) ([]Post, error) {
	var posts []Post
	rows, err := p.pool.Query(ctx, "SELECT id, title, content, user_id, created_at, updated_at FROM posts WHERE user_id = $1 ORDER BY created_at DESC", UserID)
	if err != nil {
		return nil, fmt.Errorf("error with create rows in get UsersPost")
	}
	defer rows.Close()

	for rows.Next() {
		var post Post
		err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.UserID,
			&post.CreatedAt, &post.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("error scan")
		}
		posts = append(posts, post)
	}
	return posts, nil
}
