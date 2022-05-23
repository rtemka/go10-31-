DROP TABLE IF EXISTS authors, posts;

-- таблица авторы
CREATE TABLE IF NOT EXISTS authors (
	id SERIAL PRIMARY KEY,
	name TEXT NOT NULL
);

-- таблица публикации
CREATE TABLE IF NOT EXISTS posts (
	id SERIAL PRIMARY KEY,
	title TEXT NOT NULL,
	content TEXT NOT NULL,
	author_id INTEGER DEFAULT 0,
	created_at BIGINT NOT NULL DEFAULT extract(epoch from now()),
	FOREIGN KEY(author_id) REFERENCES authors(id)
);

INSERT INTO authors(id, name) VALUES(1, 'Иван Иванов'), (2, 'Петр Петров');
INSERT INTO posts(id, title, content, author_id, created_at) 
VALUES (1, 'Постгрес Публикация номер 1', 'Lorem ipsum 1', 1, 1652355804),
(2, 'Постгрес Публикация номер 2', 'Lorem ipsum 2', 2, 1652355830);