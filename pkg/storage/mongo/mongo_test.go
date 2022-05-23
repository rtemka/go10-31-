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

func TestMongo_Posts(t *testing.T) {
	posts, err := testMongoDB.Posts()
	if err != nil {
		t.Fatalf("mongo.Posts() = error %v\n", err)
	}

	if len(posts) != len(testData) {
		t.Fatalf("mongo.Posts() = %d posts in total, want %d\n", len(posts), len(testData))
	}
}

func TestMongo_AddPost(t *testing.T) {
	newpost := storage.Post{
		Id:        3,
		Title:     "Mongo test post 1",
		Content:   "Lorem ipsum",
		Author:    storage.Author{Id: 3, Name: "Test Author"},
		CreatedAt: 0,
	}
	err := testMongoDB.AddPost(newpost)
	if err != nil {
		t.Fatalf("mongo.AddPost() = error %v\n", err)
	}

	post, err := testMongoDB.getPostById(newpost.Id)
	if err != nil {
		t.Fatalf("mongo.getPostById() = error %v\n", err)
	}

	if post != newpost {
		t.Fatalf("mongo.AddPost() = %v, want %v\n", post, newpost)
	}
}

func TestMongo_UpdatePost(t *testing.T) {
	newpost := storage.Post{
		Id:        1,
		Title:     "Mongo updated test post 1",
		Content:   "Lorem ipsum dolor sit",
		Author:    storage.Author{Id: 2, Name: "Author 2"},
		CreatedAt: 0,
	}
	err := testMongoDB.UpdatePost(newpost)
	if err != nil {
		t.Fatalf("mongo.UpdatePost() = error %v\n", err)
	}

	post, err := testMongoDB.getPostById(newpost.Id)
	if err != nil {
		t.Fatalf("mongo.getPostById() = error %v\n", err)
	}

	if post != newpost {
		t.Fatalf("mongo.UpdatePost() = %v, want %v\n", post, newpost)
	}
}

func TestMongo_DeletePost(t *testing.T) {

	err := testMongoDB.DeletePost(storage.Post{Id: 1})
	if err != nil {
		t.Fatalf("mongo.DeletePost() = error %v\n", err)
	}

	post, err := testMongoDB.getPostById(1)
	if err != nil && err != ErrNoDocuments {
		t.Fatalf("mongo.getPostById() = error %v\n", err)
	}
	if post != (storage.Post{}) {
		t.Fatalf("mongo.DeletePost() = %v, want nothing\n", post)
	}
}
