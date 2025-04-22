package auth

import (
	"auth/internal/client/db"
	"auth/internal/model"
	"auth/internal/repository"
	"auth/internal/repository/auth/converter"
	modelRepo "auth/internal/repository/auth/model"
	"context"
	"errors"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4"
	"log"
	"time"
)

const (
	tableName = "users"

	idColumn         = "id"
	nameColumn       = "name"
	emailColumn      = "email"
	passwordColumn   = "password_hash"
	roleColumn       = "role"
	isVerifiedColumn = "is_verified"
	isBlockedColumn  = "is_blocked"
	lastLoginColumn  = "last_login"
	createdAtColumn  = "created_at"
	updatedAtColumn  = "updated_at"
)

type repo struct {
	db db.Client
}

func NewRepository(db db.Client) repository.AuthRepository {
	return &repo{db: db}
}

func (r repo) Get(ctx context.Context, id string) (*model.User, error) {
	const op = "auth.Get"

	builder := sq.Select(idColumn, nameColumn, emailColumn, roleColumn, createdAtColumn, updatedAtColumn, lastLoginColumn, isVerifiedColumn, isBlockedColumn).
		PlaceholderFormat(sq.Dollar).
		From(tableName).
		Where(sq.Eq{idColumn: id}).
		Limit(1)

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	q := db.Query{
		Name:     op,
		QueryRaw: query,
	}

	var user modelRepo.User
	err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(&user.ID, &user.Name, &user.Email, &user.Role, &user.CreatedAt, &user.UpdatedAt, &user.LastLogin, &user.IsVerified, &user.IsBlocked)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("cannot get user with id: %d", user.ID)
		}
	}
	return converter.ToUserFromRepo(&user), nil
}

func (r repo) Create(ctx context.Context, user *model.CreateUser) (string, error) {
	const op = "auth.Create"

	builder := sq.Insert(tableName).
		PlaceholderFormat(sq.Dollar).
		Columns(emailColumn, nameColumn, passwordColumn, roleColumn).
		Values(user.Email, user.Name, user.Password, user.Role).
		Suffix("RETURNING id")

	query, args, err := builder.ToSql()
	if err != nil {
		return "", err
	}

	q := db.Query{
		Name:     op,
		QueryRaw: query,
	}

	var id string
	err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(&id)
	if err != nil {
		log.Println(err)
		return "", fmt.Errorf("error in create user %w", err)
	}

	log.Println("id:", id)
	return id, nil
}

func (r repo) Update(ctx context.Context, updateUser *model.UpdateUser) error {
	const op = "auth.Update"
	log.Printf("updating user %+v", updateUser)

	builder := sq.Update(tableName).
		PlaceholderFormat(sq.Dollar).
		Set(updatedAtColumn, time.Now()).
		Where(sq.Eq{idColumn: updateUser.ID}).
		Suffix("RETURNING id")

	if updateUser.Email != "" {
		builder = builder.Set(emailColumn, updateUser.Email)
	}

	query, args, err := builder.ToSql()
	if err != nil {
		return fmt.Errorf("cannot build query user with id: %d", updateUser.ID)
	}

	q := db.Query{
		Name:     op,
		QueryRaw: query,
	}

	var id string
	err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(&id)
	if errors.Is(err, pgx.ErrNoRows) {
		log.Printf("error in update user with id: %s %d", err, id)

		return fmt.Errorf("cannot update user with id: %d", updateUser.ID)
	}
	if err != nil {
		log.Printf("error in update user with id: %s", err)
		return fmt.Errorf("cannot update user %w", err)
	}

	return nil
}

func (r repo) Delete(ctx context.Context, id string) error {
	const op = "auth.Delete"

	builder := sq.Delete(tableName).PlaceholderFormat(sq.Dollar).Where(sq.Eq{idColumn: id})

	query, args, err := builder.ToSql()
	if err != nil {
		return err
	}

	q := db.Query{
		Name:     op,
		QueryRaw: query,
	}

	row := r.db.DB().QueryRowContext(ctx, q, args...)
	if row == nil {
		return fmt.Errorf("cannot delete user with id: %d", id)
	}

	return nil
}

func (r repo) Login(ctx context.Context, email string) (*model.User, error) {
	const op = "auth.Login"

	builder := sq.Select(idColumn, nameColumn, roleColumn, passwordColumn, isBlockedColumn).
		PlaceholderFormat(sq.Dollar).
		From(tableName).
		Where(sq.Eq{emailColumn: email}).
		Limit(1)
	query, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	q := db.Query{
		Name:     op,
		QueryRaw: query,
	}

	var user modelRepo.User
	err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(&user.ID, &user.Name, &user.Role, &user.Password, &user.IsBlocked)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("cannot get user with email: %s", email)
		}
	}

	return converter.ToUserFromRepo(&user), nil
}

func (r repo) Block(ctx context.Context, id string) error {
	const op = "auth.Block"

	builder := sq.Update(idColumn).
		PlaceholderFormat(sq.Dollar).
		From(tableName).
		Where(sq.Eq{id: id})

	query, args, err := builder.ToSql()
	if err != nil {
		return err
	}

	q := db.Query{
		Name:     op,
		QueryRaw: query,
	}

	var deletedId string
	err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(&deletedId)
	if errors.Is(err, pgx.ErrNoRows) {
		log.Printf("error in delete user with id: %s %d", err, id)

		return fmt.Errorf("cannot delete user with id: %d", id)
	}
	if err != nil {
		log.Printf("error in delete user with id: %s", err)
		return fmt.Errorf("cannot delete user %w", err)
	}

	return nil
}
