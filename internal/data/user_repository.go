package data

import (
	"context"
	"database/sql"

	"go.mocker.com/internal/models"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (ur *UserRepository) Create(ctx context.Context, user *models.User) error {
	stmt, err := ur.db.PrepareContext(ctx, "insert into users (email, password) values (?, ?)")
	if err != nil {
		return err
	}

	_, err = stmt.ExecContext(ctx, user.Email, user.Password)
	if err != nil {
		return err
	}
	return nil
}

func (ur *UserRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	stmt, err := ur.db.PrepareContext(ctx, "select password from users where email=?")
	if err != nil {
		return nil, err
	}

	var user models.User
	row := stmt.QueryRowContext(ctx, email)
	if err = row.Scan(&user.ID, &user.Email, &user.Password); err != nil {
		return nil, err
	}
	return &user, nil
}

func (ur *UserRepository) DeleteByID(ctx context.Context, id uint64) error {
	stmt, err := ur.db.PrepareContext(ctx, "delete from users where id=?")
	if err != nil {
		return err
	}

	_, err = stmt.ExecContext(ctx, id)
	if err != nil {
		return err
	}
	return nil
}
