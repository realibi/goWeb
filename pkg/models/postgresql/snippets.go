package postgresql

import (
	"database/sql"
	"errors"
	"realibi.com/pkg/models"
)

type SnippetModel struct{
	DB *sql.DB
}

func (m *SnippetModel) Insert(title, content, expires string) error {
	stmt := `INSERT INTO snippets (title, content, created, expires)
VALUES($1, $2, CURRENT_DATE, (CURRENT_DATE + cast($3 as interval)))`

	_, err := m.DB.Exec(stmt, title, content, expires)

	if err != nil {
		println(err.Error())
		return err
	}

	return nil
}

func (m *SnippetModel) Get(id int) (*models.Snippet, error) {
	s := &models.Snippet{}
	err := m.DB.QueryRow("SELECT id, title, content, created, expires FROM snippets\n\tWHERE id = $1", id).Scan(&s.Id, &s.Title, &s.Content, &s.Created, &s.Expires)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}
	return s, nil
}

func (m *SnippetModel) Latest() ([]*models.Snippet, error) {
	stmt := `SELECT id, title, content, created, expires FROM snippets
			 ORDER BY created DESC LIMIT 10`

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	snippets := []*models.Snippet{}

	for rows.Next(){
		s := &models.Snippet{}
		err = rows.Scan(&s.Id, &s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}

		snippets = append(snippets, s)
	}

	if err = rows.Err(); err != nil{
		return nil, err
	}

	return snippets, nil
}