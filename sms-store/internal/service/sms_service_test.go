package service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/sms/sms-store/internal/model"
	"github.com/sms/sms-store/internal/service"
)

type mockRepo struct {
	savedRecords  []*model.SmsRecord
	shouldFail  bool
	returnRecords []model.SmsRecord
}

func (m *mockRepo) Save(_ context.Context, record *model.SmsRecord) error{
	if m.shouldFail{
		return errors.New("mock DB error")
	}
	m.savedRecords = append(m.savedRecords, record)
	return nil
}
func (m *mockRepo) FindByUserID(_ context.Context, userID string) ([]model.SmsRecord,error){
	if m.shouldFail {
		return nil,errors.New("mock DB error")
	}
	return m.returnRecords, nil
}
func sampleEvent() *model.SmsEvent {
	return &model.SmsEvent{
		MessageID:  "msg-001",
		UserID:    "user-123",
		PhoneNumber:   "+919876543210",
		Message:   "Test message",
		Status:   "SUCCESS",
		VendorResponse: "MOCK-REF-001",
		Timestamp: time.Now().UTC().Format(time.RFC3339Nano),
	}
}
func TestStoreEvent_Success(t *testing.T) {
	repo := &mockRepo{}
	svc := service.NewSmsService(repo)
	err := svc.StoreEvent(context.Background(), sampleEvent())
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if len(repo.savedRecords) != 1 {
		t.Fatalf("expected 1 saved record,got %d", len(repo.savedRecords))
	}
	saved := repo.savedRecords[0]
	if saved.MessageID != "msg-001" {
		t.Errorf("expected messageId=msg-001,got %s",saved.MessageID)
	}
	if saved.UserID != "user-123" {
		t.Errorf("expected userId=user-123,got %s",saved.UserID)
	}if saved.Status != model.StatusSuccess {
		t.Errorf("expected status=SUCCESS,got %s",saved.Status)
	}
}

func TestStoreEvent_RepoFailure(t *testing.T){
	repo := &mockRepo{shouldFail: true}
	svc := service.NewSmsService(repo)
	err := svc.StoreEvent(context.Background(), sampleEvent())
	if err == nil {
		t.Fatal("expected error but got nil")
	}
}
func TestStoreEvent_BadTimestamp_UsesNow(t *testing.T){
	repo := &mockRepo{}
	svc := service.NewSmsService(repo)
	event := sampleEvent()
	event.Timestamp = "not-a-real-timestamp"
	before := time.Now().Add(-time.Second)
	err := svc.StoreEvent(context.Background(), event)
	after := time.Now().Add(time.Second)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	saved := repo.savedRecords[0]
	if saved.Timestamp.Before(before) || saved.Timestamp.After(after) {
		t.Errorf("timestamp %v is outside expected range",saved.Timestamp)
	}
}
func TestGetHistory_ReturnsList(t *testing.T){
	records := []model.SmsRecord{
		{MessageID: "msg-001", UserID:"user-123"},
		{MessageID: "msg-002", UserID: "user-123"},
	}
	repo := &mockRepo{returnRecords: records}
	svc := service.NewSmsService(repo)

	results, err := svc.GetHistory(context.Background(),"user-123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 2 {
		t.Errorf("expected 2 records, got %d", len(results))
	}
}
func TestGetHistory_EmptyUserID_ReturnsError(t *testing.T) {
	repo := &mockRepo{}
	svc := service.NewSmsService(repo)
	_, err := svc.GetHistory(context.Background(), "")
	if err == nil {
		t.Fatal("expected error for empty userID, got nil")
	}
}
func TestGetHistory_RepoFailure(t *testing.T) {
	repo := &mockRepo{shouldFail: true}
	svc := service.NewSmsService(repo)
	_, err := svc.GetHistory(context.Background(), "user-123")
	if err == nil {
		t.Fatal("expected error but got nil")
	}
}