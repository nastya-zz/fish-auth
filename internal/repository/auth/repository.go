package auth

import (
	"auth/internal/client/db"
	"auth/internal/model"
	"auth/internal/repository"
	"auth/internal/repository/auth/converter"
	modelRepo "auth/internal/repository/auth/model"
	"auth/pkg/logger"
	"context"
	"errors"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	pgx "github.com/jackc/pgx/v4"
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

func (r repo) Create(ctx context.Context, user *model.CreateUser) (*model.User, error) {
	const op = "auth.Create"

	builder := sq.Insert(tableName).
		PlaceholderFormat(sq.Dollar).
		Columns(emailColumn, nameColumn, passwordColumn, roleColumn).
		Values(user.Email, user.Name, user.Password, user.Role).
		Suffix("RETURNING *")

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	q := db.Query{
		Name:     op,
		QueryRaw: query,
	}

	var createdUser modelRepo.User
	err = r.db.DB().
		QueryRowContext(ctx, q, args...).
		Scan(&createdUser.ID, &createdUser.Name, &createdUser.Email, &createdUser.Role, &createdUser.Password, &createdUser.IsVerified, &createdUser.CreatedAt, &createdUser.UpdatedAt, &createdUser.LastLogin, &createdUser.IsBlocked)
	if err != nil {
		logger.Error("error in create user", "error", err)
		return nil, fmt.Errorf("error in create user %w", err)
	}

	logger.Info("user created successfully", "user_id", createdUser.ID, "email", createdUser.Email)
	return converter.ToUserFromRepo(&createdUser), nil
}

func (r repo) Update(ctx context.Context, updateUser *model.UpdateUser) (*model.UpdateUser, error) {
	const op = "auth.Update"
	logger.Info("updating user", "user_id", updateUser.ID, "email", updateUser.Email)
	uuId, _ := model.GetUuid(updateUser.ID)
	builder := sq.Update(tableName).
		PlaceholderFormat(sq.Dollar).
		Set(updatedAtColumn, time.Now()).
		Set(isVerifiedColumn, updateUser.IsVerified).
		Where(sq.Eq{idColumn: uuId}).
		Suffix("RETURNING id, name, email, is_verified")

	if updateUser.Email != "" {
		builder = builder.Set(emailColumn, updateUser.Email)
	}
	if updateUser.Password != "" {
		builder = builder.Set(passwordColumn, updateUser.Password)
	}
	if updateUser.Name != "" {
		builder = builder.Set(nameColumn, updateUser.Name)
	}

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("cannot build query user with id: %d", updateUser.ID)
	}

	q := db.Query{
		Name:     op,
		QueryRaw: query,
	}

	var updatedUser modelRepo.UpdateUser
	err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(&updatedUser.ID, &updatedUser.Name, &updatedUser.Email, &updatedUser.IsVerified)
	if errors.Is(err, pgx.ErrNoRows) {
		logger.Error("user not found for update", "error", err, "user_id", updateUser.ID)
		return nil, fmt.Errorf("cannot update user with id: %s", updateUser.ID)
	}
	if err != nil {
		logger.Error("error in update user", "error", err, "user_id", updateUser.ID)
		return nil, fmt.Errorf("cannot update user %w", err)
	}

	return &model.UpdateUser{
		ID:         updatedUser.ID.String(),
		Name:       updatedUser.Name,
		Email:      updatedUser.Email,
		IsVerified: updatedUser.IsVerified,
	}, nil
}

func (r repo) Delete(ctx context.Context, id string) error {
	const op = "auth.Delete"
	uuId, _ := model.GetUuid(id)
	builder := sq.Delete(tableName).PlaceholderFormat(sq.Dollar).Where(sq.Eq{idColumn: uuId})

	query, args, err := builder.ToSql()
	if err != nil {
		return err
	}

	q := db.Query{
		Name:     op,
		QueryRaw: query,
	}

	row, err := r.db.DB().ExecContext(ctx, q, args...)
	if row == nil {
		return fmt.Errorf("cannot delete user with id: %s, err %w", id, err)
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
	const op = "auth.BlockUser"

	builder := sq.Update(tableName).
		PlaceholderFormat(sq.Dollar).
		Set(isBlockedColumn, true).
		Where(sq.Eq{idColumn: id}).Suffix("RETURNING id")

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
		logger.Error("user not found for blocking", "error", err, "user_id", id)
		return fmt.Errorf("cannot block user with id: %s", id)
	}
	if err != nil {
		logger.Error("error in block user", "error", err, "user_id", id)
		return fmt.Errorf("cannot block user %w", err)
	}

	return nil
}
func (r repo) UpdateRole(ctx context.Context, id, role string) error {
	const op = "auth.repository.UpdateRole"

	builder := sq.Update(tableName).
		PlaceholderFormat(sq.Dollar).
		Set(roleColumn, role).
		Where(sq.Eq{idColumn: id}).Suffix("RETURNING id")

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
		logger.Error("user not found for update role", "error", err, "user_id", id)
		return fmt.Errorf("cannot update role user with id: %s", id)
	}
	if err != nil {
		logger.Error("cannot update role user with id", "error", err, "user_id", id)
		return fmt.Errorf("cannot update role user with id %w", err)
	}

	return nil
}
