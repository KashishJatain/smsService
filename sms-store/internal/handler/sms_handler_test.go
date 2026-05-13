package handler_test
import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"


	"github.com/sms/sms-store/internal/handler"
	"github.com/sms/sms-store/internal/model"
	"github.com/sms/sms-store/internal/service"
)

type mockRepo struct{
	records []model.SmsRecord
	shouldFail bool
}
func (m *mockRepo) Save(_ context.Context, _ *model.SmsRecord) error {return nil}
func (m *mockRepo) FindByUserID(_ context.Context, _ string) ([]model.SmsRecord,error){
	if m.shouldFail{
		return nil,context.DeadlineExceeded
	}
	return m.records,nil
}
func newHandler(records []model.SmsRecord,fail bool) http.Handler{
	repo := &mockRepo{records:records,shouldFail: fail}
	svc := service.NewSmsService(repo)
	h := handler.NewSmsHandler(svc)
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)
	return mux
}
func TestGetUserMessages_OK(t*testing.T){
	records :=[]model.SmsRecord{
		{MessageID:"msg-1",UserID:"user-42"},
	}
	h := newHandler(records,false)
	req := httptest.NewRequest(http.MethodGet,"/v1/user/user-42/messages",nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w,req)
	if w.Code != http.StatusOK{
		t.Fatalf("expected 200, got %d",w.Code)
	}
	var resp model.SmsHistoryResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil{
		t.Fatalf("failed to decode response: %v", err)
	}if resp.UserID != "user-42" {
		t.Errorf("expected userId=user-42,got %s", resp.UserID)
	}
	if resp.Count != 1 {
		t.Errorf("expected count=1, got %d", resp.Count)
	}}

func TestGetUserMessages_EmptyHistory(t*testing.T){
	h := newHandler(nil, false)
	req := httptest.NewRequest(http.MethodGet,"/v1/user/user-99/messages",nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var resp model.SmsHistoryResponse
	_ =json.NewDecoder(w.Body).Decode(&resp)
	if resp.Count != 0 {
		t.Errorf("expected count=0, got %d", resp.Count)
	}
	if resp.Messages == nil{
		t.Error("expected empty slice, got nil")
	}
}

func TestGetUserMessages_RepoError_Returns500(t *testing.T){
	h := newHandler(nil, true)
	req := httptest.NewRequest(http.MethodGet, "/v1/user/user-42/messages", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)
	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", w.Code)
	}
}
func TestGetUserMessages_WrongMethod_Returns405(t*testing.T){
	h := newHandler(nil, false)
	req := httptest.NewRequest(http.MethodPost, "/v1/user/user-42/messages", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)
	if w.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405,got %d", w.Code)
	}
}
func TestHealth_Returns200(t *testing.T){
	h := newHandler(nil, false)
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200,got %d", w.Code)
	}
}