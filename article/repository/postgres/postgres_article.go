package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/makidai/go-clean-arch-backend/article/repository"
	"github.com/makidai/go-clean-arch-backend/domain"
	"github.com/sirupsen/logrus"
)

type postgresArticleRepository struct {
	Conn *sql.DB
}

func NewPostgresArticleRepository(Conn *sql.DB) domain.ArticleRepository {
	return &postgresArticleRepository{Conn}
}

func (m *postgresArticleRepository) fetch(ctx context.Context, query string, args ...interface{}) (result []domain.Article, err error) {
	rows, err := m.Conn.QueryContext(ctx, query, args...)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	defer func() {
		errRow := rows.Close()
		if errRow != nil {
			logrus.Error(errRow)
		}
	}()

	result = make([]domain.Article, 0)
	for rows.Next() {
		t := domain.Article{}
		authorID := int64(0)
		err = rows.Scan(
			&t.ID,
			&t.Title,
			&t.Content,
			&authorID,
			&t.UpdatedAt,
			&t.CreatedAt,
		)

		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		t.Author = domain.Author{
			ID: authorID,
		}
		result = append(result, t)
	}

	return result, nil
}

func (m *postgresArticleRepository) Fetch(ctx context.Context, cursor string, num int64) (res []domain.Article, nextCursor string, err error) {
	query := `SELECT id, title, content, author_id, updated_at, created_at
  				FROM article WHERE created_at > ? ORDER BY created_at LIMIT ? `

	decodedCursor, err := repository.DecodeCursor(cursor)
	if err != nil && cursor != "" {
		return nil, "", domain.ErrBadParamInput
	}

	res, err = m.fetch(ctx, query, decodedCursor, num)
	if err != nil {
		return nil, "", err
	}

	if len(res) == int(num) {
		nextCursor = repository.EncodeCursor(res[len(res)-1].CreatedAt)
	}

	return
}

func (m *postgresArticleRepository) GetByID(ctx context.Context, id int64) (res domain.Article, err error) {
	query := `SELECT id,title,content, author_id, updated_at, created_at
  				FROM article WHERE ID = ?`

	list, err := m.fetch(ctx, query, id)
	if err != nil {
		return domain.Article{}, err
	}

	if len(list) > 0 {
		res = list[0]
	} else {
		return res, domain.ErrNotFound
	}

	return
}

func (m *postgresArticleRepository) GetByTitle(ctx context.Context, title string) (res domain.Article, err error) {
	query := `SELECT id,title,content, author_id, updated_at, created_at
  				FROM article WHERE title = ?`

	list, err := m.fetch(ctx, query, title)
	if err != nil {
		return
	}

	if len(list) > 0 {
		res = list[0]
	} else {
		return res, domain.ErrNotFound
	}
	return
}

func (m *postgresArticleRepository) Store(ctx context.Context, a *domain.Article) (err error) {
	query := `INSERT  article SET title=? , content=? , author_id=?, updated_at=? , created_at=?`
	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return
	}

	res, err := stmt.ExecContext(ctx, a.Title, a.Content, a.Author.ID, a.UpdatedAt, a.CreatedAt)
	if err != nil {
		return
	}
	lastID, err := res.LastInsertId()
	if err != nil {
		return
	}
	a.ID = lastID
	return
}

func (m *postgresArticleRepository) Delete(ctx context.Context, id int64) (err error) {
	query := "DELETE FROM article WHERE id = ?"

	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return
	}

	res, err := stmt.ExecContext(ctx, id)
	if err != nil {
		return
	}

	rowsAfected, err := res.RowsAffected()
	if err != nil {
		return
	}

	if rowsAfected != 1 {
		err = fmt.Errorf("Weird  Behavior. Total Affected: %d", rowsAfected)
		return
	}

	return
}

func (m *postgresArticleRepository) Update(ctx context.Context, ar *domain.Article) (err error) {
	query := `UPDATE article set title=?, content=?, author_id=?, updated_at=? WHERE ID = ?`

	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return
	}

	res, err := stmt.ExecContext(ctx, ar.Title, ar.Content, ar.Author.ID, ar.UpdatedAt, ar.ID)
	if err != nil {
		return
	}
	affect, err := res.RowsAffected()
	if err != nil {
		return
	}
	if affect != 1 {
		err = fmt.Errorf("Weird  Behavior. Total Affected: %d", affect)
		return
	}

	return
}
