package postgres

import (
	"GoNews/pkg/storage"
	"context"
	"os"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

var ErrNoRows = pgx.ErrNoRows

// Postgres выполняет CRUD операции с БД
type Postgres struct {
	db *pgxpool.Pool
}

// New выполняет подключение
// и возвращает объект для взаимодействия с БД
func New(connString string) (*Postgres, error) {

	pool, err := pgxpool.Connect(context.Background(), connString)
	if err != nil {
		return nil, err
	}

	return &Postgres{db: pool}, nil
}

// Close выполняет закрытие подключения к БД
func (p *Postgres) Close() {
	p.db.Close()
}

// AddPost создает пост в БД
func (p *Postgres) AddPost(post storage.Post) error {
	tx, err := p.db.Begin(context.Background())
	if err != nil {
		return err
	}

	defer tx.Rollback(context.Background())

	// добавляем в БД сначала автора, если
	// передан без id
	if post.Author.Id == 0 {
		post.Author.Id, err = p.addAuthor(tx, post.Author)
		if err != nil {
			return err
		}
	}

	stmt := `
		INSERT INTO posts(id, title, content, author_id, created_at)
		VALUES ($1, $2, $3, $4, $5);
	`

	_, err = tx.Exec(context.Background(), stmt,
		post.Id, post.Title, post.Content, post.Author.Id, post.CreatedAt)
	if err != nil {
		return err
	}

	return tx.Commit(context.Background())
}

// addAuthor добавляет автора публикации и возвращает его новый id
func (p *Postgres) addAuthor(tx pgx.Tx, a storage.Author) (int, error) {
	stmt := `
			INSERT INTO authors(name)
			VALUES ($1) RETURNING id;
	`
	var id int
	err := tx.QueryRow(context.Background(), stmt, a.Name).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

// Posts возвращает список всех публикаций
func (p *Postgres) Posts() ([]storage.Post, error) {

	stmt := `
		SELECT 
			p.id, 
			p.title, 
			p.content,
			p.created_at,  
			a.name,
			a.id  
		FROM
			posts AS p INNER JOIN authors AS a ON p.author_id = a.id;
	`

	rows, err := p.db.Query(context.Background(), stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []storage.Post

	for rows.Next() {
		var post storage.Post

		err = rows.Scan(
			&post.Id, &post.Title, &post.Content, &post.CreatedAt,
			&post.Author.Name, &post.Author.Id)
		if err != nil {
			return nil, err
		}

		posts = append(posts, post)
	}

	return posts, rows.Err()
}

// getPost возвращает публикацию по id
func (p *Postgres) getPost(id int) (storage.Post, error) {

	stmt := `
		SELECT 
			p.id, 
			p.title, 
			p.content,
			p.created_at,  
			a.name,
			a.id  
		FROM
			posts AS p INNER JOIN authors AS a ON p.author_id = a.id
		WHERE p.id = $1;
	`

	var post storage.Post
	err := p.db.QueryRow(context.Background(), stmt, id).Scan(
		&post.Id, &post.Title, &post.Content, &post.CreatedAt, &post.Author.Name, &post.Author.Id)
	if err != nil {
		return post, err
	}

	return post, nil
}

// UpdatePost обновялет публикацию
func (p *Postgres) UpdatePost(post storage.Post) error {

	stmt := `
		UPDATE posts
		SET title = $2,
			content = $3,
			author_id = $4,
			created_at = $5
		WHERE id = $1;
	`
	ctx := context.Background()

	tx, err := p.db.Begin(ctx)
	if err != nil {
		return err
	}

	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, stmt,
		post.Id, post.Title, post.Content, post.Author.Id, post.CreatedAt)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

// DeletePost удаляет публикацию
func (p *Postgres) DeletePost(post storage.Post) error {

	stmt := `
		DELETE FROM posts
		WHERE posts.id = $1;
	`
	ctx := context.Background()

	tx, err := p.db.Begin(ctx)
	if err != nil {
		return err
	}

	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, stmt, post.Id)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (p *Postgres) testCleanUp() error {

	b, err := os.ReadFile("testdata/testCleanUp.sql")
	if err != nil {
		return err
	}

	ctx := context.Background()

	tx, err := p.db.Begin(ctx)
	if err != nil {
		return err
	}

	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, string(b))
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}
