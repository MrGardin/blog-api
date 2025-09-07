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
