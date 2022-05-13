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

func TestPosts(t *testing.T) {
	posts, err := db.Posts()
	if err != nil {
		t.Fatalf("expected to get posts, got error [%v]\n", err)
	}

	if len(posts) != postsNum {
		t.Fatalf("expected to get %d posts in total, got [%d]\n", postsNum, len(posts))
	}
}

func TestAddPost(t *testing.T) {
	newpost := storage.Post{
		Id:        3,
		Author:    storage.Author{Id: 1, Name: "Иван Иванов"},
		Title:     "Test title1",
		Content:   "Test content1",
		CreatedAt: 0,
	}
	err := db.AddPost(newpost)
	if err != nil {
		t.Fatalf("expected to create post, got error [%v]\n", err)
	}

	post, err := db.getPost(newpost.Id)
	if err != nil {
		t.Fatalf("expected to get post, got error [%v]\n", err)
	}

	if post != newpost {
		t.Fatalf("expected to get created post %v, got %v\n", newpost, post)
	}
}

func TestUpdatePost(t *testing.T) {
	newpost := storage.Post{
		Id:        1,
		Author:    storage.Author{Id: 2, Name: "Петр Петров"},
		Title:     "Updated title",
		Content:   "Updated content",
		CreatedAt: 0,
	}
	err := db.UpdatePost(newpost)
	if err != nil {
		t.Fatalf("expected to update post, got error [%v]\n", err)
	}

	post, err := db.getPost(newpost.Id)
	if err != nil {
		t.Fatalf("expected to get posts, got error [%v]\n", err)
	}

	if post != newpost {
		t.Fatalf("expected to get updated post %v, got %v\n", newpost, post)
	}
}

func TestDeletePost(t *testing.T) {
	p := storage.Post{Id: 1}
	err := db.DeletePost(p)
	if err != nil {
		t.Fatalf("expected to delete post, got error [%v]\n", err)
	}

	post, err := db.getPost(p.Id)
	if err != nil && err != ErrNoRows {
		t.Fatalf("got error while deleting from DB [%v]\n", err)
	}
	if post != (storage.Post{}) {
		t.Fatalf("expected to get no rows, got [%v]\n", post)
	}
}
