package web

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/uiratan/fullcycle-archdev-event-driven-architecture/balances/internal/usecase/find_account"
)

type WebBalanceHandler struct {
	FindAccountUseCase find_account.FindAccountUseCase
}

func NewWebBalanceHandler(findAccountUseCase find_account.FindAccountUseCase) *WebBalanceHandler {
	return &WebBalanceHandler{
		FindAccountUseCase: findAccountUseCase,
	}
}

func (h *WebBalanceHandler) FindAccount(w http.ResponseWriter, r *http.Request) {
	accountID := chi.URLParam(r, "account_id")
	fmt.Println("accountID", accountID)
	account, err := h.FindAccountUseCase.Execute(find_account.FindAccountInputDTO{AccountID: accountID})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	accountJson, err := json.Marshal(account)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	fmt.Fprint(w, string(accountJson))
}
