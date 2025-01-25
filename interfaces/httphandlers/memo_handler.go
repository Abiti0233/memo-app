// interfaces/httphandlers/memo_handler.go
package httphandlers

import (
	"encoding/json"
	"net/http"

	"github.com/Abiti0233/memo-app/domains/memo"
	"github.com/Abiti0233/memo-app/usecases"
	"github.com/gorilla/mux"
)

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
		http.Error(w, "ペイロードに失敗", http.StatusBadRequest)
		return
	}

	userID, ok := r.Context().Value("userID").(string)
	if !ok || userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	memo := &memo.Memo{
		Title:      req.Title,
		Content:    req.Content,
		IsArchived: false,
	}

	err := h.useCase.CreateMemo(userID, memo)
	if err != nil {
		http.Error(w, "メモの作成に失敗", http.StatusInternalServerError)
		return
	}

	// http.StatusCreatedは、HTTPステータスコード201を表す定数
	// 新しくリソースが作成されているから201を返す
	w.WriteHeader(http.StatusCreated)
	// json.NewEncoder(w).Encode(memo)は、memoをJSON形式にエンコードして、HTTPレスポンスに書き込む
	json.NewEncoder(w).Encode(memo)
}

// UpdateMemo はメモ更新のハンドラー
func (h *MemoHandler) UpdateMemo(w http.ResponseWriter, r *http.Request) {
	var req UpdateMemoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "リクエストが正しくない", http.StatusBadRequest)
		return
	}

	vars := mux.Vars(r)
	id := vars["memoId"]

	userID, ok := r.Context().Value("userID").(string)
	if !ok || userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	memo := &memo.Memo{
		ID:         id,
		Title:      req.Title,
		Content:    req.Content,
		IsArchived: req.IsArchived,
	}

	err := h.useCase.UpdateMemo(userID, memo)
	if err != nil {
		if err == usecases.ErrMemoNotFound {
			http.Error(w, "Memo not found", http.StatusNotFound)
			return
		}
		http.Error(w, "メモ作成失敗", http.StatusInternalServerError)
		return
	}
	// http.StatusOKは、HTTPステータスコード200を表す定数
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(memo)
}

// DeleteMemo はメモ削除のハンドラー
func (h *MemoHandler) DeleteMemo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["memoId"]

	userID, ok := r.Context().Value("userID").(string)
	if !ok || userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	err := h.useCase.DeleteMemo(userID, id)
	if err != nil {
		if err == usecases.ErrMemoNotFound {
			http.Error(w, "メモが見つからない", http.StatusNotFound)
			return
		}
		http.Error(w, "メモ削除に失敗", http.StatusInternalServerError)
		return
	}

	// http.StatusNoContentは、HTTPステータスコード204を表す定数
	w.WriteHeader(http.StatusNoContent)
}

// GetMemo はメモ取得のハンドラー
func (h *MemoHandler) GetMemo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["memoId"]

	userID, ok := r.Context().Value("userID").(string)
	if !ok || userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	memo, err := h.useCase.GetMemoByID(userID, id)
	if err != nil {
		http.Error(w, "メモ取得失敗", http.StatusInternalServerError)
		return
	}
	if memo == nil {
		http.Error(w, "メモが見つからない", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(memo)
}

// ListMemos はメモ一覧取得のハンドラー
func (h *MemoHandler) ListMemos(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(string)
	if !ok || userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	memos, err := h.useCase.ListMemos(userID)
	if err != nil {
		http.Error(w, "メモ一覧取得失敗", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(memos)
}

// ArchiveMemo はメモのアーカイブ化/アーカイブ解除のハンドラー
func (h *MemoHandler) ArchiveMemo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["memoId"]

	userID, ok := r.Context().Value("userID").(string)
	if !ok || userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// リクエストボディからアーカイブ状態を取得
	var req struct {
		IsArchived bool `json:"isArchived"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	err := h.useCase.ArchiveMemo(userID, id, req.IsArchived)
	if err != nil {
		if err == usecases.ErrMemoNotFound {
			http.Error(w, "メモが見つからない", http.StatusNotFound)
			return
		}
		http.Error(w, "アーカイブに失敗", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
