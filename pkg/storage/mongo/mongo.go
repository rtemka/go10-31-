package mongo

import (
	"GoNews/pkg/storage"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var ErrNoDocuments = mongo.ErrNoDocuments

// Mongo выполняет CRUD операции с БД
type Mongo struct {
	client *mongo.Client
	// для более реалистичных сценариев
	// тут наверное надо использовать что-то
	// типа мапы с именами баз данных и
	// соответствующими коллекциями
	databaseName   string
	collectionName string
}

// New выполняет подключение
// и возвращает объект для взаимодействия с БД
func New(connString string, dbName, collectionName string) (*Mongo, error) {

	mongoOpts := options.Client().ApplyURI(connString)

	client, err := mongo.Connect(context.Background(), mongoOpts)
	if err != nil {
		return nil, err
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		return nil, err
	}

	return &Mongo{
		client:         client,
		databaseName:   dbName,
		collectionName: collectionName,
	}, nil
}

// Close выполняет закрытие подключения к БД
func (m *Mongo) Close() {
	m.client.Disconnect(context.Background())
}

// AddPost создает пост в БД
func (m *Mongo) AddPost(post storage.Post) error {

	collection := m.client.Database(m.databaseName).Collection(m.collectionName)

	// если в коллекции отсутствует публикация по
	// заданному id, то выполняем INSERT, иначе UPDATE
	opts := options.Replace().SetUpsert(true)

	filter := bson.D{bson.E{Key: "_id", Value: post.Id}}

	_, err := collection.ReplaceOne(context.Background(), filter, post, opts)
	if err != nil {
		return err
	}

	return nil
}

// UpdatePost обновялет публикацию
func (m *Mongo) UpdatePost(post storage.Post) error {
	return m.AddPost(post)
}

// DeletePost удаляет публикацию
func (m *Mongo) DeletePost(post storage.Post) error {
	collection := m.client.Database(m.databaseName).Collection(m.collectionName)

	_, err := collection.DeleteOne(context.Background(), bson.D{bson.E{Key: "_id", Value: post.Id}})
	if err != nil {
		return err
	}

	return nil
}

// Posts возвращает список всех публикаций
func (m *Mongo) Posts() ([]storage.Post, error) {
	collection := m.client.Database(m.databaseName).Collection(m.collectionName)

	cur, err := collection.Find(context.Background(), bson.D{})
	if err != nil {
		return nil, err
	}

	defer cur.Close(context.Background())

	var posts []storage.Post

	err = cur.All(context.Background(), &posts)
	if err != nil {
		return nil, err
	}

	return posts, nil
}

func (m *Mongo) getPostById(id any) (storage.Post, error) {
	collection := m.client.Database(m.databaseName).Collection(m.collectionName)

	var p storage.Post

	err := collection.FindOne(context.Background(),
		bson.D{bson.E{Key: "_id", Value: id}}).Decode(&p)
	if err != nil {
		return p, err
	}

	return p, nil
}

func (m *Mongo) testCleanUp(dbName, collectionName string) error {
	err := m.dropDB(dbName)
	if err != nil {
		return err
	}

	return m.restoreTestDB(dbName, collectionName)
}

func (m *Mongo) dropDB(dbName string) error {
	return m.client.Database(dbName).Drop(context.Background())
}

func (m *Mongo) restoreTestDB(dbName, collectionName string) error {
	collection := m.client.Database(dbName).Collection(collectionName)

	_, err := collection.InsertMany(context.Background(), testData)
	if err != nil {
		return err
	}
	return nil
}

var testData = []any{
	storage.Post{Id: 1,
		Title:     "mongo post 1",
		Content:   "Lorem ipsum",
		CreatedAt: 1652431685,
		Author:    storage.Author{Id: 1, Name: "Author 1"}},
	storage.Post{Id: 2,
		Title:     "mongo post 2",
		Content:   "Lorem ipsum",
		CreatedAt: 1652431703,
		Author:    storage.Author{Id: 2, Name: "Author 2"}},
}
