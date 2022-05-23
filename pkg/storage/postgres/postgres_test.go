package postgres

import (
	"GoNews/pkg/storage"
	"log"
	"os"
	"testing"
)

var db *Postgres

const (
	postsNum = 2
)

func TestMain(m *testing.M) {
	var err error

	dbUrl := os.Getenv("POSTGRES_DB_TEST_URL")
	if dbUrl == "" {
		log.Fatal("environment variable POSTGRES_DB_TEST_URL must be set")
	}

	db, err = New(dbUrl)
	if err != nil {
		log.Fatal(err)
	}

	exitCode := m.Run()

	err = db.testCleanUp()
	if err != nil {
		db.Close()
		log.Fatal(err)
	}

	db.Close()

	os.Exit(exitCode)
}

func TestPostgres_Posts(t *testing.T) {
	posts, err := db.Posts()
	if err != nil {
		t.Fatalf("postgres.Posts() = error %v\n", err)
	}

	if len(posts) != postsNum {
		t.Fatalf("postgres.Posts() = %d posts in total, want %d\n", len(posts), postsNum)
	}
}

func TestPostgres_AddPost(t *testing.T) {
	newpost := storage.Post{
		Id:        3,
		Author:    storage.Author{Id: 1, Name: "Иван Иванов"},
		Title:     "Test title1",
		Content:   "Test content1",
		CreatedAt: 0,
	}
	err := db.AddPost(newpost)
	if err != nil {
		t.Fatalf("postgres.AddPost() = error %v\n", err)
	}

	post, err := db.getPost(newpost.Id)
	if err != nil {
		t.Fatalf("postgres.getPost() = error %v\n", err)
	}

	if post != newpost {
		t.Fatalf("postgres.AddPost() = %v, want %v\n", post, newpost)
	}
}

func TestPostgres_UpdatePost(t *testing.T) {
	newpost := storage.Post{
		Id:        1,
		Author:    storage.Author{Id: 2, Name: "Петр Петров"},
		Title:     "Updated title",
		Content:   "Updated content",
		CreatedAt: 0,
	}
	err := db.UpdatePost(newpost)
	if err != nil {
		t.Fatalf("postgres.UpdatePost() = error %v\n", err)
	}

	post, err := db.getPost(newpost.Id)
	if err != nil {
		t.Fatalf("postgres.getPost() = error %v\n", err)
	}

	if post != newpost {
		t.Fatalf("postgres.UpdatePost() = %v, want %v\n", post, newpost)
	}
}

func TestPostgres_DeletePost(t *testing.T) {
	p := storage.Post{Id: 1}
	err := db.DeletePost(p)
	if err != nil {
		t.Fatalf("postgres.DeletePost() = error %v\n", err)
	}

	post, err := db.getPost(p.Id)
	if err != nil && err != ErrNoRows {
		t.Fatalf("postgres.getPost() = error %v\n", err)
	}
	if post != (storage.Post{}) {
		t.Fatalf("postgres.DeletePost() = %v, want nothing\n", post)
	}
}
