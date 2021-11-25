package server

import (
	"app/api/llog"
	"app/api/presentation/handler"
	"app/api/presentation/middleware"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

type Server interface {
	Serve()
	Route(appHandler *handler.AppHandler)
}

type server struct {
	Handler *mux.Router
	Addr    string
}

func New(addr string) Server {
	r := mux.NewRouter()

	srv := &server{
		Addr:    addr,
		Handler: r,
	}
	return srv
}

func (s *server) Serve() {
	llog.Info(fmt.Sprintf("server running %s...", s.Addr))
	llog.Fatal(http.ListenAndServe(s.Addr, s.Handler))
}

func (s *server) Route(appHandler *handler.AppHandler) {

	s.Handler.Use(middleware.CommonMiddleware)

	authRouter := s.Handler.PathPrefix("/").Subrouter()
	authRouter.Use(middleware.AuthMiddleware)

	adminRouter := s.Handler.PathPrefix("/").Subrouter()
	// TODO: middleware
	//adminRouter.Use(middleware.AdminMiddleware)

	s.Handler.HandleFunc("/ping", pingHandler).Methods(http.MethodGet, http.MethodOptions)

	s.Handler.HandleFunc("/login", appHandler.AuthHandler.Login).Methods(http.MethodPost, http.MethodOptions)
	s.Handler.HandleFunc("/account", appHandler.UserHandler.Create).Methods(http.MethodPost, http.MethodOptions)
	s.Handler.HandleFunc("/users", appHandler.UserHandler.GetAll).Methods(http.MethodGet, http.MethodOptions)
	s.Handler.HandleFunc("/users/{id}", appHandler.UserHandler.GetByID).Methods(http.MethodGet, http.MethodOptions)
	s.Handler.HandleFunc("/users/{id}/img", appHandler.FileHandler.GetUserIcon).Methods(http.MethodGet, http.MethodOptions)
	s.Handler.HandleFunc("/users/{id}/follows", appHandler.UserHandler.GetFollows).Methods(http.MethodGet, http.MethodOptions)
	s.Handler.HandleFunc("/users/{id}/followers", appHandler.UserHandler.GetFollowers).Methods(http.MethodGet, http.MethodOptions)
	s.Handler.HandleFunc("/categories", appHandler.CategoryHandler.GetAll).Methods(http.MethodGet, http.MethodOptions)
	s.Handler.HandleFunc("/tags", appHandler.TagHandler.GetAll).Methods(http.MethodGet, http.MethodOptions)
	s.Handler.HandleFunc("/tags/{id}", appHandler.TagHandler.GetByID).Methods(http.MethodGet, http.MethodOptions)
	s.Handler.HandleFunc("/threads", appHandler.ThreadHandler.GetAll).Methods(http.MethodGet, http.MethodOptions)
	s.Handler.HandleFunc("/threads/{id}", appHandler.ThreadHandler.GetByID).Methods(http.MethodGet, http.MethodOptions)
	s.Handler.HandleFunc("/threads/{id}/members", appHandler.ThreadHandler.GetMembersByThreadID).Methods(http.MethodGet, http.MethodOptions)

	{
		authRouter.HandleFunc("/logout", appHandler.AuthHandler.Logout).Methods(http.MethodDelete, http.MethodOptions)

		authRouter.HandleFunc("/account", appHandler.UserHandler.GetMe).Methods(http.MethodGet, http.MethodOptions)
		authRouter.HandleFunc("/account", appHandler.UserHandler.DeleteMe).Methods(http.MethodDelete, http.MethodOptions)

		authRouter.HandleFunc("/account/img", appHandler.FileHandler.SetUserIcon).Methods(http.MethodPost, http.MethodOptions)
		authRouter.HandleFunc("/account/profile", appHandler.UserHandler.UpdateProfile).Methods(http.MethodPut, http.MethodOptions)
		authRouter.HandleFunc("/account/user-id", appHandler.UserHandler.UpdateUserID).Methods(http.MethodPut, http.MethodOptions)
		authRouter.HandleFunc("/account/password", appHandler.UserHandler.UpdatePassword).Methods(http.MethodPut, http.MethodOptions)

		authRouter.HandleFunc("/account/tags", appHandler.TagHandler.AddTagToUser).Methods(http.MethodPost, http.MethodOptions)
		authRouter.HandleFunc("/account/tags/{tagID}", appHandler.TagHandler.RemoveTagFromUser).Methods(http.MethodDelete, http.MethodOptions)

		authRouter.HandleFunc("/users/{followedUUID}/follows", appHandler.UserHandler.Follow).Methods(http.MethodPost, http.MethodOptions)
		authRouter.HandleFunc("/users/{followedUUID}/follows", appHandler.UserHandler.Unfollow).Methods(http.MethodDelete, http.MethodOptions)

		authRouter.HandleFunc("/tags", appHandler.TagHandler.Create).Methods(http.MethodPost, http.MethodOptions)

		// TODO: thread search
		authRouter.HandleFunc("/threads", appHandler.ThreadHandler.Create).Methods(http.MethodPost, http.MethodOptions)
		authRouter.HandleFunc("/threads/{id}", appHandler.ThreadHandler.Update).Methods(http.MethodPut, http.MethodOptions)
		authRouter.HandleFunc("/threads/{id}", appHandler.ThreadHandler.Delete).Methods(http.MethodDelete, http.MethodOptions)

		authRouter.HandleFunc("/threads/{id}/members", appHandler.ThreadHandler.Join).Methods(http.MethodPost, http.MethodOptions)
		authRouter.HandleFunc("/threads/{id}/members", appHandler.ThreadHandler.Leave).Methods(http.MethodDelete, http.MethodOptions)

		authRouter.HandleFunc("/threads/{id}/members/{userID}", appHandler.ThreadHandler.ForceToLeave).Methods(http.MethodDelete, http.MethodOptions)

		authRouter.HandleFunc("/threads/{threadID}/messages", appHandler.MessageHandler.GetByThreadID).Methods(http.MethodGet, http.MethodOptions)
		authRouter.HandleFunc("/threads/{threadID}/messages", appHandler.MessageHandler.Create).Methods(http.MethodPost, http.MethodOptions)
		authRouter.HandleFunc("/threads/{threadID}/messages/{messageID}", appHandler.MessageHandler.AddFavorite).Methods(http.MethodPost, http.MethodOptions)

		authRouter.HandleFunc("/threads/{threadID}/files", appHandler.FileHandler.Upload).Methods(http.MethodPost, http.MethodOptions)
		authRouter.HandleFunc("/threads/{threadID}/files/{fileID}", appHandler.FileHandler.Download).Methods(http.MethodGet, http.MethodOptions)

		authRouter.HandleFunc("/threads/{threadID}/tags", appHandler.TagHandler.AddTagToThread).Methods(http.MethodPost, http.MethodOptions)

		authRouter.HandleFunc("/threads/{threadID}/tags/{tagID}", appHandler.TagHandler.RemoveTagFromThread).Methods(http.MethodDelete, http.MethodOptions)

		authRouter.HandleFunc("/ws", appHandler.SocketHandler.WebsocketConnect).Methods(http.MethodGet, http.MethodOptions)

	}

	{
		adminRouter.HandleFunc("/categories", appHandler.CategoryHandler.Create).Methods(http.MethodPost, http.MethodOptions)
		adminRouter.HandleFunc("/categories/{id}", appHandler.CategoryHandler.Update).Methods(http.MethodPut, http.MethodOptions)
		adminRouter.HandleFunc("/categories/{id}", appHandler.CategoryHandler.Delete).Methods(http.MethodDelete, http.MethodOptions)
	}

}

func pingHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("pong"))
}
