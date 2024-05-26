package userRepos

import (
	"context"
	userEntities "github.com/MercyClassic/go_saga/src/app/domain/entities/user"
	"github.com/MercyClassic/go_saga/src/app/infrastructure/db/client"
	userErrors "github.com/MercyClassic/go_saga/src/app/infrastructure/db/errors"
)

type UserRepositoryInterface interface {
	GetUserById(ctx context.Context, userId int) (*userEntities.User, error)
	GetUsers(ctx context.Context) ([]*userEntities.User, error)
	SaveUser(ctx context.Context, user *userEntities.User) error
	GetUserBalance(ctx context.Context, userId int) (float32, error)
	SetUserBalance(ctx context.Context, userId int, amount float32) error
}

type UserRepository struct {
	client client.Client
}

func NewUserRepository(client client.Client) *UserRepository {
	return &UserRepository{client: client}
}

func (r *UserRepository) GetUserById(ctx context.Context, userId int) (*userEntities.User, error) {
	q := `
		select 
		    id, name, username, balance 
		from account
		where id=$1;
	`
	user := userEntities.User{}
	err := r.client.QueryRow(ctx, q, userId).Scan(
		&user.Id,
		&user.Name,
		&user.Username,
		&user.Balance,
	)
	if err != nil {
		return nil, userErrors.ErrUserNotFound
	}
	return &user, nil
}

func (r *UserRepository) GetUsers(ctx context.Context) ([]*userEntities.User, error) {
	q := `
		select 
		    id, name, username, balance 
		from account;
	`
	userList := make([]*userEntities.User, 0, 100)

	rows, err := r.client.Query(ctx, q)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var user userEntities.User
		if err := rows.Scan(
			&user.Id,
			&user.Name,
			&user.Username,
			&user.Balance,
		); err != nil {
			return nil, err
		}
		userList = append(userList, &user)
	}
	return userList, nil
}

func (r *UserRepository) SaveUser(ctx context.Context, user *userEntities.User) error {
	q := `
		insert into account
		    (name, username)
		values ($1, $2)
		returning id;
	`
	err := r.client.QueryRow(ctx, q, user.Name, user.Username).Scan(&user.Id)
	if err != nil {
		return userErrors.ErrUserAlreadyExists
	}
	return nil
}

func (r *UserRepository) GetUserBalance(ctx context.Context, userId int) (float32, error) {
	q := `
		select balance
		from account
		where user_id = $1
		for update skip locked;
	`
	var balance float32
	err := r.client.QueryRow(ctx, q, userId).Scan(&balance)
	if err != nil {
		return 0, err
	}
	return balance, nil
}

func (r *UserRepository) SetUserBalance(
	ctx context.Context,
	userId int,
	amount float32,
) error {
	tx, err := r.client.Begin(ctx)
	if err != nil {
		return err
	}
	q := `
		update account 
		set balance = balance - $1
		where user_id = $2;
	`
	_, err = tx.Exec(ctx, q, userId, amount)
	if err != nil {
		return err
	}
	err = tx.Commit(ctx)
	if err != nil {
		return err
	}
	return nil
}
