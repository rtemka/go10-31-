package api

import (
	"GoNews/pkg/storage"
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"
)

// methods - ассоциативный массив, где
// в роли ключа предполагается метод GET, POST... и т.д.
// а в роли значения соответствующий обработчик
type methods map[string]http.Handler

// Api представляет собой интерфейс нашего сервиса,
// занимается запросами к БД, логгированием, а также
// содержит карту ресурсов, их доступных http-методов
// и обработчиков этих методов
type Api struct {
	db        storage.Model
	logger    *log.Logger
	resources map[string]methods
}

// New возвращает объект API нашего сервиса
func New(s storage.Model, log *log.Logger) *Api {
	api := Api{db: s, logger: log}

	// назаначаем обработчики соответствующим ресурсам
	api.resources = map[string]methods{
		"/posts": {
			http.MethodGet:    http.HandlerFunc(api.getPostsHandler),
			http.MethodPost:   http.HandlerFunc(api.postPostHandler),
			http.MethodPut:    http.HandlerFunc(api.putPostHandler),
			http.MethodDelete: http.HandlerFunc(api.deletePostHandler),
		},
	}

	return &api
}

// Mux возвращает мультиплексер для работы с API
func (api *Api) Mux() http.Handler {
	// создаем мультиплексер и назначаем обработчики
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.WriteHeader(http.StatusNoContent)
	})
	mux.Handle("/posts", api)
	return drainAndClose(mux)
}

// drainAndClose вспомогательная функция, опустошает
// и закрывает request body, это позволяет
// клиенту переиспользовать tcp-сессию
func drainAndClose(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
		_, _ = io.Copy(io.Discard, r.Body)
		_ = r.Body.Close()
	})
}

// writeResponse вспомогательная функция, которая
// устанавливает заголовки ответа, пишет тело сообщения в виде json(если есть)
func (api *Api) writeResponse(w http.ResponseWriter, reply any, code int) {
	if reply == nil {
		http.Error(w, "", code)
		return
	}

	buf := new(bytes.Buffer)

	err := json.NewEncoder(buf).Encode(reply)
	if err != nil {
		api.logger.Printf("error encoding response [%v] -> [%v]\n", reply, err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("Content-Length", strconv.Itoa(buf.Len()))
	w.WriteHeader(code)

	_, err = w.Write(buf.Bytes())
	if err != nil {
		api.logger.Printf("error writing response [%v]\n", err)
	}
}

func (api *Api) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	// если у нас имеется требуемый ресурс
	if resourceMethods, ok := api.resources[r.URL.Path]; ok {

		// и имеется требуемый обработчик
		if handler, ok := resourceMethods[r.Method]; ok {

			if handler != nil {
				handler.ServeHTTP(w, r)
			} else {
				// чтобы не вызывать панику на сервере, если вдруг
				// на место обработчика назначен nil
				api.logger.Printf("error setting handler: the handler [%s] for [%s] is nil\n",
					r.Method, r.URL.Path)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
			}

			return
		}

		// если не нашли обработчик,
		// проверяем может это OPTIONS, тогда
		// возвращаем список допустимых методов
		if r.Method == http.MethodOptions {
			w.Header().Add("Allow", resourceMethods.allowedMethods())
			api.writeResponse(w, nil, http.StatusOK)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
		return
	}

	// не нашли требуемый ресурс
	api.writeResponse(w, nil, http.StatusNotFound)
}

// allowedMethods возвращает список
// допустимых методов для ресурса в виде строки
func (m methods) allowedMethods() string {
	a := make([]string, 0, len(m))

	for k := range m {
		a = append(a, k)
	}

	sort.Strings(a)

	return strings.Join(a, ", ")
}

// getPostsHandler обработчик для метода GET
func (api *Api) getPostsHandler(w http.ResponseWriter, r *http.Request) {
	posts, err := api.db.Posts()
	if err != nil {
		api.logger.Printf("error fetching from database: [%v]\n", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	api.writeResponse(w, map[string]any{"data": posts}, http.StatusOK)
}

// postPostHandler обработчик для метода POST
func (api *Api) postPostHandler(w http.ResponseWriter, r *http.Request) {

	var post storage.Post

	err := json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		api.logger.Printf("error decoding request body [%v]\n", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	err = api.db.AddPost(post)
	if err != nil {
		api.logger.Printf("error posting to database: [%v]\n", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	api.writeResponse(w, nil, http.StatusCreated)

}

// putPostHandler обработчик для метода PUT
func (api *Api) putPostHandler(w http.ResponseWriter, r *http.Request) {

	var post storage.Post

	err := json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		api.logger.Printf("error decoding request body [%v]\n", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	err = api.db.UpdatePost(post)
	if err != nil {
		api.logger.Printf("error updating in database: [%v]\n", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	api.writeResponse(w, nil, http.StatusOK)

}

// putPostHandler обработчик для метода DELETE
func (api *Api) deletePostHandler(w http.ResponseWriter, r *http.Request) {

	var post storage.Post

	err := json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		api.logger.Printf("error decoding request body [%v]\n", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	err = api.db.DeletePost(post)
	if err != nil {
		api.logger.Printf("error updating in database: [%v]\n", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	api.writeResponse(w, nil, http.StatusOK)

}
