// interfaces/httphandlers/auth_handler.go
package httphandlers

import (
    "context"
    "encoding/json"
		"time"
    "net/http"

    "golang.org/x/oauth2"
    "golang.org/x/oauth2/google"
    "github.com/google/uuid"
    "github.com/Abiti0233/memo-app/domains"
    "github.com/Abiti0233/memo-app/usecases"
    "github.com/Abiti0233/memo-app/utils"
)

type AuthHandler struct {
    userUseCase usecases.UserUseCase
    oauthConfig *oauth2.Config
    jwtSecret   string
}

func NewAuthHandler(userUseCase usecases.UserUseCase, oauthConfig *oauth2.Config, jwtSecret string) *AuthHandler {
    return &AuthHandler{
        userUseCase: userUseCase,
        oauthConfig: oauthConfig,
        jwtSecret:   jwtSecret,
    }
}

// Login はGoogle OAuth認証のためのリダイレクトハンドラー
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
    state := uuid.New().String()
    url := h.oauthConfig.AuthCodeURL(state, oauth2.AccessTypeOffline)
    http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

// Callback はGoogle OAuth認証後のコールバックハンドラー
// func (レシーバー) メソッド名（引数）→メソッドの構文
// レシーバーは、メソッドが属する構造体のインスタンスを指す。
// w http.ResponseWriterは、HTTPレスポンスを書き込むためのインターフェース
// r *http.Requestは、HTTPリクエスト情報を含む構造体へのポインタ
func (h *AuthHandler) Callback(w http.ResponseWriter, r *http.Request) {

	// 認可コードを取得
  code := r.URL.Query().Get("code")
  if code == "" {
    http.Error(w, "コードが見つかりませんでした。", http.StatusBadRequest)
    return
  }

	// Exchangeメソッドで取得した認可コードを使用してアクセストークンを取得
  token, err := h.oauthConfig.Exchange(context.Background(), code)
  if err != nil {
      http.Error(w, "トークンの交換に失敗しました。", http.StatusInternalServerError)
      return
  }
	// oauth2.ConfigのClientメソッドは指定されたコンテキストとトークンを使用して、認証されたHTTPクライアントを返すもの。
  client := h.oauthConfig.Client(context.Background(), token)
  resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
  if err != nil {
      http.Error(w, "ユーザー情報の取得に失敗しました。", http.StatusInternalServerError)
      return
  }
  defer resp.Body.Close()
  var userInfo struct {
      ID            string `json:"id"`
      Email         string `json:"email"`
      VerifiedEmail bool   `json:"verified_email"`
      Name          string `json:"name"`
  }
	// レスポンスのボディをデコードして、userInfoに格納
  if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
      http.Error(w, "ユーザー情報のデコードに失敗しました。", http.StatusInternalServerError)
      return
  }
  // ユーザーをデータベースに登録または取得
  user, err := h.userUseCase.GetUserByID(userInfo.ID)
  if err != nil {
      http.Error(w, "ユーザーの取得に失敗しました。", http.StatusInternalServerError)
      return
  }
  if user == nil {
      // 新規ユーザーとして登録
      newUser := &domains.User{
          ID:            userInfo.ID,
          Name:          userInfo.Name,
          Email:         userInfo.Email,
          EmailVerified: nil,
          CreatedAt:     timeNow(),
          UpdatedAt:     timeNow(),
      }
      if err := h.userUseCase.RegisterUser(newUser); err != nil {
          http.Error(w, "ユーザーの登録に失敗しました。", http.StatusInternalServerError)
          return
      }
      user = newUser
  }
  // JWTトークンの生成
  tokenStr, err := utils.GenerateJWT(user.ID, h.jwtSecret)
  if err != nil {
      http.Error(w, "トークンの生成に失敗しました。", http.StatusInternalServerError)
      return
  }
  // トークンをJSONで返却
  response := map[string]string{
      "token": tokenStr,
  }
  w.Header().Set("Content-Type", "application/json")
  json.NewEncoder(w).Encode(response)
}

func timeNow() time.Time {
    return time.Now().UTC()
}
