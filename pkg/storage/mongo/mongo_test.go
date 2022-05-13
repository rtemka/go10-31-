package mongo

import (
	"GoNews/pkg/storage"
	"log"
	"os"
	"testing"
)

var testMongoDB *Mongo

const (
	testDbName       = "GoNewsTest"
	testDbCollection = "posts"
)

func TestMain(m *testing.M) {

	var err error

	dbUrl := os.Getenv("MONGO_DB_TEST_URL")
	if dbUrl == "" {
		log.Fatal("environment variable MONGO_DB_TEST_URL must be set")
	}

	testMongoDB, err = New(dbUrl, testDbName, testDbCollection)
	if err != nil {
		log.Fatal(err)
	}

	// восстанавливаем состояние тестовой БД
	err = testMongoDB.testCleanUp(testDbName, testDbCollection)
	if err != nil {
		testMongoDB.Close()
		log.Fatal(err)
	}

	// прогоняем тесты
	exitCode := m.Run()

	// восстанавливаем состояние тестовой БД
	err = testMongoDB.testCleanUp(testDbName, testDbCollection)
	if err != nil {
		testMongoDB.Close()
		log.Fatal(err)
	}

	testMongoDB.Close()

	os.Exit(exitCode)
}

func TestPosts(t *testing.T) {
	posts, err := testMongoDB.Posts()
	if err != nil {
		t.Fatalf("expected to get db data, got error [%v]\n", err)
	}

	if len(posts) != len(testData) {
		t.Fatalf("expected to get %d db data in total, got [%d]\n", len(testData), len(posts))
	}
}

func TestAddPost(t *testing.T) {
	newpost := storage.Post{
		Id:        3,
		Title:     "Mongo test post 1",
		Content:   "Lorem ipsum",
		Author:    storage.Author{Id: 3, Name: "Test Author"},
		CreatedAt: 0,
	}
	err := testMongoDB.AddPost(newpost)
	if err != nil {
		t.Fatalf("expected to get db data, got error [%v]\n", err)
	}

	post, err := testMongoDB.getPostById(newpost.Id)
	if err != nil {
		t.Fatalf("expected to get db data, got error [%v]\n", err)
	}

	if post != newpost {
		t.Fatalf("expected to get created db data %v, got %v\n", newpost, post)
	}
}

func TestUpdatePost(t *testing.T) {
	newpost := storage.Post{
		Id:        1,
		Title:     "Mongo updated test post 1",
		Content:   "Lorem ipsum dolor sit",
		Author:    storage.Author{Id: 2, Name: "Author 2"},
		CreatedAt: 0,
	}
	err := testMongoDB.UpdatePost(newpost)
	if err != nil {
		t.Fatalf("expected to get db data, got error [%v]\n", err)
	}

	post, err := testMongoDB.getPostById(newpost.Id)
	if err != nil {
		t.Fatalf("expected to get db data, got error [%v]\n", err)
	}

	if post != newpost {
		t.Fatalf("expected to get updated db data %v, got %v\n", newpost, post)
	}
}

func TestDeletePost(t *testing.T) {

	err := testMongoDB.DeletePost(storage.Post{Id: 1})
	if err != nil {
		t.Fatalf("expected to get db data, got error [%v]\n", err)
	}

	post, err := testMongoDB.getPostById(1)
	if err != nil && err != ErrNoDocuments {
		t.Fatalf("got error while deleting from DB [%v]\n", err)
	}
	if post != (storage.Post{}) {
		t.Fatalf("expected to get no documents, got [%v]\n", post)
	}
}
