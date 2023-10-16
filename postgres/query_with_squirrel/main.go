package query_with_squirrel

import (
	"context"
	"log"
	"time"

	sq "github.com/Masterminds/squirrel"
	deps "github.com/MikhailRibalkov/auth/pkg/auth_v1/pkg/auth_v1"
	"github.com/jackc/pgx/v4/pgxpool"
)

const (
	dbDSN = "host=localhost port=5432 dbname=auth user=auth-user password=auth-password sslmode=disable"
)

type PgClient struct {
	ctx context.Context
}

type PgInterface interface {
	CreateUser(user deps.CreateRequest)
	GetUserInfo(id int64)
	UpdateUser()
	DeleteUser(id int64)
}

func (pgc *PgClient) CreateUser(user *deps.CreateRequest) (id int64, err error) {
	// Создаем пул соединений с базой данных
	var usr = user.GetUser()
	pool, err := pgxpool.Connect(pgc.ctx, dbDSN)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
		return -1, err
	}
	defer pool.Close()

	// Делаем запрос на вставку записи в таблицу note
	builderInsert := sq.Insert("auth").
		PlaceholderFormat(sq.Dollar).
		Columns("name", "email", "role").
		Values(usr.Name, usr.Email, usr.Role).
		Suffix("RETURNING id")

	query, args, err := builderInsert.ToSql()
	if err != nil {
		log.Fatalf("failed to build query: %v", err)
		return -1, err
	}

	var authID int
	err = pool.QueryRow(pgc.ctx, query, args...).Scan(&authID)
	if err != nil {
		log.Fatalf("failed to insert note: %v", err)
		return -1, err
	}

	log.Printf("inserted note with id: %d", authID)
	return int64(authID), nil
}

func (pgc *PgClient) GetUserInfo(reqId *deps.GetRequest) (info deps.UserInfo, err error) {
	var id = reqId.GetId()
	pool, err := pgxpool.Connect(pgc.ctx, dbDSN)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
		return deps.UserInfo{}, err
	}
	defer pool.Close()

	// Делаем запрос на получение измененной записи из таблицы note
	builderSelectOne := sq.Select("id", "name", "email", "role", "created_at", "updated_at").
		From("auth").
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"id": id}).
		Limit(1)

	query, args, err := builderSelectOne.ToSql()
	if err != nil {
		log.Fatalf("failed to build query: %v", err)
		return deps.UserInfo{}, err
	}

	info = deps.UserInfo{}

	err = pool.QueryRow(pgc.ctx, query, args...).Scan(&info.Id, &info.Name, &info.Email, &info.Role, &info.CreatedAt, &info.UpdatedAt)
	if err != nil {
		log.Fatalf("failed to select notes: %v", err)
		return deps.UserInfo{}, err
	}

	log.Printf("id: %d, name: %s, email: %s, role %d, created_at: %v, updated_at: %v\n", info.Id, info.Name, info.Email, info.Role, info.CreatedAt, info.UpdatedAt)

	return info, nil
}

func (pgc *PgClient) UpdateUser(req *deps.UpdateRequest) (id int64, err error) {
	pool, err := pgxpool.Connect(pgc.ctx, dbDSN)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
		return -1, err
	}
	defer pool.Close()

	// Делаем запрос на обновление записи в таблице note
	builderUpdate := sq.Update("auth").
		PlaceholderFormat(sq.Dollar).
		Set("name", req.Name).
		Set("email", req.Email).
		Set("role", req.Role).
		Set("updated_at", time.Now()).
		Where(sq.Eq{"id": req.Id})

	query, args, err := builderUpdate.ToSql()
	if err != nil {
		log.Fatalf("failed to build query: %v", err)
		return -1, err
	}

	res, err := pool.Exec(pgc.ctx, query, args...)
	if err != nil {
		log.Fatalf("failed to update note: %v", err)
		return -1, err
	}

	log.Printf("updated %d rows", res.RowsAffected())
	return req.Id, nil
}

func (pgc *PgClient) DeleteUser(req *deps.DeleteRequest) (id int64, err error) {
	var delId = req.GetId()
	pool, err := pgxpool.Connect(pgc.ctx, dbDSN)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
		return -1, err
	}
	defer pool.Close()

	// Создаем запрос на удаление
	builderUpdate := sq.Delete("auth").
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"id": delId})

	query, args, err := builderUpdate.ToSql()
	if err != nil {
		log.Fatalf("failed to build query: %v", err)
		return -1, err
	}

	res, err := pool.Exec(pgc.ctx, query, args...)
	if err != nil {
		log.Fatalf("failed to update note: %v", err)
		return -1, err
	}

	log.Printf("updated %d rows", res.RowsAffected())
	return delId, nil
}
