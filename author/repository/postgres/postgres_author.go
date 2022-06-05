package postgres

import (
	"context"
	"database/sql"

	"github.com/makidai/go-clean-arch-backend/domain"
)

type postgresAuthorRepo struct {
	Conn *sql.DB
}

func NewPostgresAuthorRepository(Conn *sql.DB) domain.AuthorRepository {
	return &postgresAuthorRepo{Conn}
}

func (m *postgresAuthorRepo) getOne(ctx context.Context, query string, args ...interface{}) (res domain.Author, err error) {
	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return domain.Author{}, err
	}
	row := stmt.QueryRowContext(ctx, args...)
	res = domain.Author{}

	err = row.Scan(
		&res.ID,
		&res.Name,
		&res.CreatedAt,
		&res.UpdatedAt,
	)
	return
}

func (m *postgresAuthorRepo) GetByID(ctx context.Context, id int64) (domain.Author, error) {
	query := `SELECT id, name, created_at, updated_at FROM author WHERE id=?`
	return m.getOne(ctx, query, id)
}
