package handler
import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"
	"github.com/sms/sms-store/internal/model"
	"github.com/sms/sms-store/internal/service"
)

type SmsHandler struct{
	svc *service.SmsService
}
func NewSmsHandler(svc *service.SmsService) *SmsHandler{
	return &SmsHandler{svc: svc}
}
func (h *SmsHandler) RegisterRoutes(mux *http.ServeMux){
	mux.HandleFunc("/v1/user/",h.routeUserMessages)
	mux.HandleFunc("/health", h.Health)
}
func (h *SmsHandler) routeUserMessages(w http.ResponseWriter, r *http.Request){
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) != 4 || parts[0] != "v1" || parts[1] != "user" || parts[3] != "messages"{
		writeError(w, http.StatusNotFound, "route not found", "")
		return
	}
	userID := parts[2]
	if userID == "" {
		writeError(w, http.StatusBadRequest, "userId is required", "")
		return
	}
	switch r.Method {
	case http.MethodGet:
		h.GetUserMessages(w, r, userID)
	default:
		writeError(w, http.StatusMethodNotAllowed, "method not allowed", "")
	}
}
func (h*SmsHandler) GetUserMessages(w http.ResponseWriter,r*http.Request,userID string){
	slog.Info("GET /v1/user/:userId/messages", "userId", userID)
	records, err := h.svc.GetHistory(r.Context(), userID)
	if err != nil {
		slog.Error("Failed to fetch SMS history", "userId", userID, "err", err)
		writeError(w, http.StatusInternalServerError, "Failed to fetch messages", err.Error())
		return
	}
	if records== nil{
		records= []model.SmsRecord{}
	}
	resp := model.SmsHistoryResponse{
		UserID:userID,
		Count: len(records),
		Messages:records,
	}
	writeJSON(w,http.StatusOK,resp)
}
func (h*SmsHandler) Health(w http.ResponseWriter,r *http.Request){
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok", "service": "sms-store"})
}
func writeJSON(w http.ResponseWriter,status int, payload any){
	w.Header().Set("Content-Type","application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		slog.Error("Failed to write JSON response", "err", err)
	}
}
func writeError(w http.ResponseWriter, status int, errMsg, detail string){
	writeJSON(w,status,model.ErrorResponse{Error:errMsg,Message:detail})
}
