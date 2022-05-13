package storage

// Post содержит информацию о статье
type Post struct {
	Id        int    `bson:"_id"`
	Author    Author `bson:"author"`
	Title     string `bson:"title"`
	Content   string `bson:"content"`
	CreatedAt int64  `bson:"created_at"`
}

// Author содержит информацию об авторе
type Author struct {
	Id   int    `bson:"_id"`
	Name string `bson:"name"`
}

// Model задаёт контракт на работу с БД.
type Model interface {
	Posts() ([]Post, error) // получение всех публикаций
	AddPost(Post) error     // создание новой публикации
	UpdatePost(Post) error  // обновление публикации
	DeletePost(Post) error  // удаление публикации по ID
	Close()                 // закрытие подключения к БД
}
