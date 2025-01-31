// interfaces/httphandlers/memo_handler.go
package httphandlers

import (
	"encoding/json"
	"net/http"

	"log"

	"github.com/Abiti0233/memo-app/domains/memo"
	"github.com/Abiti0233/memo-app/interfaces/middlewares/authmiddleware"
	"github.com/Abiti0233/memo-app/usecases"
	"github.com/gorilla/mux"
)

// writeJSONErrorはエラー時にJSON形式でレスポンスを返すためのヘルパー関数
func writeJSONError(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(map[string]string{
		"error": message,
	})
}

type MemoHandler struct {
	useCase usecases.MemoUseCase
}

func NewMemoHandler(useCase usecases.MemoUseCase) *MemoHandler {
	return &MemoHandler{useCase: useCase}
}

// CreateMemoRequest はメモ作成リクエストの構造体
type CreateMemoRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

// UpdateMemoRequest はメモ更新リクエストの構造体
type UpdateMemoRequest struct {
	Title      string `json:"title"`
	Content    string `json:"content"`
	IsArchived bool   `json:"isArchived"`
}

// CreateMemo はメモ作成のハンドラー
func (h *MemoHandler) CreateMemo(w http.ResponseWriter, r *http.Request) {
	var req CreateMemoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "ペイロードのデコードに失敗しました")
		return
	}

	userID, ok := r.Context().Value(authmiddleware.UserIDKey).(string)
	if !ok || userID == "" {
		writeJSONError(w, http.StatusUnauthorized, "認証情報がありません（Unauthorized）")
		return
	}

	m := &memo.Memo{
		Title:      req.Title,
		Content:    req.Content,
		IsArchived: false,
	}

	err := h.useCase.CreateMemo(userID, m)
	if err != nil {
		log.Printf("メモの作成に失敗: %v\n", err)
		writeJSONError(w, http.StatusInternalServerError, "メモの作成に失敗しました")
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(m); err != nil {
		log.Printf("レスポンスのエンコードに失敗: %v\n", err)
	}
}

// UpdateMemo はメモ更新のハンドラー
func (h *MemoHandler) UpdateMemo(w http.ResponseWriter, r *http.Request) {
	var req UpdateMemoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "リクエストのデコードに失敗しました")
		return
	}

	vars := mux.Vars(r)
	id := vars["memoId"]

	userID, ok := r.Context().Value(authmiddleware.UserIDKey).(string)
	if !ok || userID == "" {
		writeJSONError(w, http.StatusUnauthorized, "認証情報がありません（Unauthorized）")
		return
	}

	m := &memo.Memo{
		ID:         id,
		Title:      req.Title,
		Content:    req.Content,
		IsArchived: req.IsArchived,
	}

	err := h.useCase.UpdateMemo(userID, m)
	if err != nil {
		if err == usecases.ErrMemoNotFound {
			writeJSONError(w, http.StatusNotFound, "メモが見つかりません")
			return
		}
		log.Printf("メモ更新失敗: %v\n", err)
		writeJSONError(w, http.StatusInternalServerError, "メモ更新に失敗しました")
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(m); err != nil {
	}
}

// DeleteMemo はメモ削除のハンドラー
func (h *MemoHandler) DeleteMemo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["memoId"]

	userID, ok := r.Context().Value(authmiddleware.UserIDKey).(string)
	if !ok || userID == "" {
		writeJSONError(w, http.StatusUnauthorized, "認証情報がありません（Unauthorized）")
		return
	}

	err := h.useCase.DeleteMemo(userID, id)
	if err != nil {
		if err == usecases.ErrMemoNotFound {
			writeJSONError(w, http.StatusNotFound, "メモが見つかりません")
			return
		}
		log.Printf("メモ削除失敗: %v\n", err)
		writeJSONError(w, http.StatusInternalServerError, "メモ削除に失敗しました")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GetMemo はメモ取得のハンドラー
func (h *MemoHandler) GetMemo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["memoId"]

	userID, ok := r.Context().Value(authmiddleware.UserIDKey).(string)
	if !ok || userID == "" {
		writeJSONError(w, http.StatusUnauthorized, "認証情報がありません（Unauthorized）")
		return
	}

	m, err := h.useCase.GetMemoByID(userID, id)
	if err != nil {
		log.Printf("メモ取得失敗: %v\n", err)
		writeJSONError(w, http.StatusInternalServerError, "メモ取得に失敗しました")
		return
	}
	if m == nil {
		writeJSONError(w, http.StatusNotFound, "メモが見つかりません")
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(m); err != nil {
		log.Printf("レスポンスのエンコードに失敗: %v\n", err)
	}
}

// ListMemos はメモ一覧取得のハンドラー
func (h *MemoHandler) ListMemos(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(authmiddleware.UserIDKey).(string)
	log.Printf("userIDKey: %v\n", userID)
	log.Printf("ok: %v\n", ok)
	if !ok || userID == "" {
		log.Printf("認証情報がありません（Unauthorized）")
		writeJSONError(w, http.StatusUnauthorized, "認証情報がありません（Unauthorized）")
		return
	}

	memos, err := h.useCase.ListMemos(userID)
	if err != nil {
		log.Printf("メモ一覧取得失敗: %v\n", err)
		writeJSONError(w, http.StatusInternalServerError, "メモ一覧の取得に失敗しました")
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(memos); err != nil {
		log.Printf("レスポンスのエンコードに失敗: %v\n", err)
	}
}

// ArchiveMemo はメモのアーカイブ化/アーカイブ解除のハンドラー
func (h *MemoHandler) ArchiveMemo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["memoId"]

	userID, ok := r.Context().Value(authmiddleware.UserIDKey).(string)
	if !ok || userID == "" {
		writeJSONError(w, http.StatusUnauthorized, "認証情報がありません（Unauthorized）")
		return
	}

	// リクエストボディからアーカイブ状態を取得
	var req struct {
		IsArchived bool `json:"isArchived"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "リクエストのデコードに失敗しました")
		return
	}

	err := h.useCase.ArchiveMemo(userID, id, req.IsArchived)
	if err != nil {
		if err == usecases.ErrMemoNotFound {
			writeJSONError(w, http.StatusNotFound, "メモが見つかりません")
			return
		}
		log.Printf("アーカイブ失敗: %v\n", err)
		writeJSONError(w, http.StatusInternalServerError, "アーカイブに失敗しました")
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
}
