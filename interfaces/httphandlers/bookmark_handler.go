// interfaces/httphandlers/bookmark_handler.go
package httphandlers

import (
	"encoding/json"
	"net/http"

	"github.com/Abiti0233/memo-app/usecases"
	"github.com/gorilla/mux"
)

type BookmarkHandler struct {
	useCase usecases.BookmarkUseCase
}

func NewBookmarkHandler(useCase usecases.BookmarkUseCase) *BookmarkHandler {
	return &BookmarkHandler{useCase: useCase}
}

// CreateBookmarkRequest はブックマーク作成リクエストの構造体
type CreateBookmarkRequest struct {
	MemoID string `json:"memoId"`
}

// CreateBookmark はブックマーク作成のハンドラー
func (h *BookmarkHandler) CreateBookmark(w http.ResponseWriter, r *http.Request) {
	var req CreateBookmarkRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "リクエスト形式が違う", http.StatusBadRequest)
		return
	}

	if req.MemoID == "" {
		http.Error(w, "メモIDは必須", http.StatusBadRequest)
		return
	}

	userID, ok := r.Context().Value("userID").(string)
	if !ok || userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	bookmark, err := h.useCase.CreateBookmark(userID, req.MemoID)
	if err != nil {
		if err == usecases.ErrBookmarkAlreadyExists {
			http.Error(w, "すでにブックマーク", http.StatusConflict)
			return
		}
		http.Error(w, "ブックマーク失敗", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(bookmark)
}

// DeleteBookmark はブックマーク削除のハンドラー
func (h *BookmarkHandler) DeleteBookmark(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["bookmarkId"]

	userID, ok := r.Context().Value("userID").(string)
	if !ok || userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	err := h.useCase.DeleteBookmark(userID, id)
	if err != nil {
		if err == usecases.ErrBookmarkNotFound {
			http.Error(w, "ブックマークが見つからない", http.StatusNotFound)
			return
		}
		http.Error(w, "ブックマーク削除失敗", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GetBookmark はブックマーク取得のハンドラー
func (h *BookmarkHandler) GetBookmark(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["bookmarkId"]

	userID, ok := r.Context().Value("userID").(string)
	if !ok || userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	bookmark, err := h.useCase.GetBookmarkByID(userID, id)
	if err != nil {
		http.Error(w, "取得失敗", http.StatusInternalServerError)
		return
	}
	if bookmark == nil {
		http.Error(w, "ブックマークが見つからない", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(bookmark)
}

// ListBookmarks はブックマーク一覧取得のハンドラー
func (h *BookmarkHandler) ListBookmarks(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(string)
	if !ok || userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	bookmarks, err := h.useCase.ListBookmarks(userID)
	if err != nil {
		http.Error(w, "ブックマークのリストに失敗", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(bookmarks)
}
