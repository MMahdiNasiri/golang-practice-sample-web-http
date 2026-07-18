package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"sample-web-http/internal/user"
)

type UserRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (u *UserRepo) Create(ctx context.Context, user *user.User) (*user.User, error) {
	err := u.db.QueryRowContext(ctx, "INSERT INTO users (user_name, password) VALUES ($1, $2) returning id, created_at", user.UserName, user.Password).Scan(&user.ID, &user.CreatedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *UserRepo) FindById(ctx context.Context, id int) (*user.User, error) {
	row := u.db.QueryRowContext(ctx, "SELECT id, user_name, created_at, updated_at FROM users WHERE id = $1", id)
	var userStruct user.User
	err := row.Scan(&userStruct.ID, &userStruct.UserName, &userStruct.CreatedAt, &userStruct.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &userStruct, nil
}

func (u *UserRepo) FindByUserName(ctx context.Context, userName string) (*user.User, error) {
	row := u.db.QueryRowContext(ctx, "SELECT id, user_name, created_at, updated_at FROM users WHERE user_name = $1", userName)
	var userStruct user.User
	err := row.Scan(&userStruct.ID, &userStruct.UserName, &userStruct.CreatedAt, &userStruct.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &userStruct, nil
}

func (u *UserRepo) FindByUserNameWithPassword(ctx context.Context, userName string) (*user.User, error) {
	row := u.db.QueryRowContext(ctx, "SELECT id, user_name, password, created_at, updated_at FROM users WHERE user_name = $1", userName)
	var userStruct user.User
	err := row.Scan(&userStruct.ID, &userStruct.UserName, &userStruct.Password, &userStruct.CreatedAt, &userStruct.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &userStruct, nil
}

func (u *UserRepo) Delete(ctx context.Context, id int) error {
	result, err := u.db.ExecContext(ctx, "DELETE FROM users WHERE id = $1", id)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("expected row affected")
	}
	return nil
}

func (u *UserRepo) List(ctx context.Context, page int) ([]*user.User, error) {
	rows, err := u.db.QueryContext(ctx, "SELECT id, user_name, created_at, updated_at From users limit 10 offset $1", page*10)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			return
		}
	}(rows)
	userList := make([]*user.User, 0)
	for rows.Next() {
		var userStruct user.User
		err = rows.Scan(&userStruct.ID, &userStruct.UserName, &userStruct.CreatedAt, &userStruct.UpdatedAt)
		if err != nil {
			return nil, err
		}
		userList = append(userList, &userStruct)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return userList, nil
}
