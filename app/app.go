package app

import (
    "database/sql"
    "log"
    "net/http"

    "github.com/gorilla/mux"
    "github.com/Abiti0233/memo-app/infrastructures"
    "github.com/Abiti0233/memo-app/interfaces/httphandlers"
    "github.com/Abiti0233/memo-app/interfaces/middlewares/authmiddleware"
    "github.com/Abiti0233/memo-app/usecases"
    "github.com/Abiti0233/memo-app/configs"
    "github.com/Abiti0233/memo-app/utils"
    "golang.org/x/oauth2"
    "golang.org/x/oauth2/google"
)

type App struct {
    Router          *mux.Router
    DB              *sql.DB
    UserUseCase     usecases.UserUseCase
    MemoUseCase     usecases.MemoUseCase
    BookmarkUseCase usecases.BookmarkUseCase
}

func (a *App) Initialize(db *sql.DB, config *configs.Config) {
    a.DB = db
    a.Router = mux.NewRouter()

    // リポジトリの初期化
    userRepo := infrastructures.NewUserRepository(a.DB)
    memoRepo := infrastructures.NewMemoRepository(a.DB)
    bookmarkRepo := infrastructures.NewBookmarkRepository(a.DB)
    categoryRepo := infrastructures.NewCategoryRepository(a.DB)

    // ユースケースの初期化
    userUseCase := usecases.NewUserUseCase(userRepo)
    memoUseCase := usecases.NewMemoUseCase(memoRepo, categoryRepo)
    bookmarkUseCase := usecases.NewBookmarkUseCase(bookmarkRepo)

    a.UserUseCase = userUseCase
    a.MemoUseCase = memoUseCase
    a.BookmarkUseCase = bookmarkUseCase

    // OAuth設定の初期化
    oauthConfig := &oauth2.Config{
        ClientID:     config.OAuthClientID,
        ClientSecret: config.OAuthClientSecret,
        RedirectURL:  config.OAuthRedirectURL,
        Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
        Endpoint:     google.Endpoint,
    }

    // ハンドラーの初期化
    authHandler := httphandlers.NewAuthHandler(a.UserUseCase, oauthConfig, config.JWTSecret)
    memoHandler := httphandlers.NewMemoHandler(a.MemoUseCase)
    bookmarkHandler := httphandlers.NewBookmarkHandler(a.BookmarkUseCase)

    // ミドルウェアの初期化
    jwtMiddleware := authmiddleware.JWTAuthentication(config)

    // ルートの設定
    a.initializeRoutes(authHandler, memoHandler, bookmarkHandler, jwtMiddleware)
}

func (a *App) initializeRoutes(authHandler *httphandlers.AuthHandler, memoHandler *httphandlers.MemoHandler, bookmarkHandler *httphandlers.BookmarkHandler, jwtMiddleware func(http.Handler) http.Handler) {
    api := a.Router.PathPrefix("/api/v1").Subrouter()

    // 認証が不要なルート
    api.HandleFunc("/auth/login", authHandler.Login).Methods("GET")
    api.HandleFunc("/auth/callback", authHandler.Callback).Methods("GET")

    // 認証が必要なルート
    authenticated := api.PathPrefix("/").Subrouter()
    authenticated.Use(jwtMiddleware)

    // メモ関連のルート
    authenticated.HandleFunc("/memos", memoHandler.CreateMemo).Methods("POST")
    authenticated.HandleFunc("/memos/{memoId}", memoHandler.UpdateMemo).Methods("PUT")
    authenticated.HandleFunc("/memos/{memoId}", memoHandler.DeleteMemo).Methods("DELETE")
    authenticated.HandleFunc("/memos/{memoId}", memoHandler.GetMemo).Methods("GET")
    authenticated.HandleFunc("/memos", memoHandler.ListMemos).Methods("GET")
    authenticated.HandleFunc("/memos/{memoId}/archive", memoHandler.ArchiveMemo).Methods("PATCH")

    // ブックマーク関連のルート
    authenticated.HandleFunc("/bookmarks", bookmarkHandler.CreateBookmark).Methods("POST")
    authenticated.HandleFunc("/bookmarks/{bookmarkId}", bookmarkHandler.DeleteBookmark).Methods("DELETE")
    authenticated.HandleFunc("/bookmarks/{bookmarkId}", bookmarkHandler.GetBookmark).Methods("GET")
    authenticated.HandleFunc("/bookmarks", bookmarkHandler.ListBookmarks).Methods("GET")
}

func (a *App) Run(addr string) {
    log.Printf("Server running on %s", addr)
    log.Fatal(http.ListenAndServe(addr, a.Router))
}
