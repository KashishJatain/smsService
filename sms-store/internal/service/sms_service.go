package service
import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/sms/sms-store/internal/model"
	"github.com/sms/sms-store/internal/repository"
)
type SmsService struct {
	repo repository.SmsRepository
}
func NewSmsService(repo repository.SmsRepository) *SmsService {
	return &SmsService{repo: repo}
}
func (s *SmsService) StoreEvent(ctx context.Context,event *model.SmsEvent) error {
	slog.Info("Storing SMS event",
		"messageId", event.MessageID,
		"userId", event.UserID,
		"status", event.Status,
	)
	ts, err := parseTimestamp(event.Timestamp)
	if err != nil {
		slog.Warn("Could not parse event timestamp,using now","raw",event.Timestamp,"err",err)
		ts = time.Now().UTC()
	}
	record := &model.SmsRecord{
		MessageID: event.MessageID,
		UserID:  event.UserID,
		PhoneNumber: event.PhoneNumber,
		Message:event.Message,
		Status:  model.SmsStatus(event.Status),
		VendorResponse: event.VendorResponse,
		Timestamp:  ts,
	}
	if err := s.repo.Save(ctx, record); err != nil {
		return fmt.Errorf("smsService.StoreEvent: %w",err)
	}
	slog.Info("SMS event stored successfully","messageId", event.MessageID)
	return nil
}
func (s *SmsService) GetHistory(ctx context.Context, userID string) ([]model.SmsRecord, error){
	if userID == "" {
		return nil, fmt.Errorf("userID must not be empty")
	}
	records, err := s.repo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("smsService.GetHistory: %w", err)
	}
	slog.Info("Fetched SMS history", "userId", userID, "count", len(records))
	return records, nil
}

func parseTimestamp(raw string) (time.Time, error){
	// Try RFC3339 / ISO-8601 (what Java's Instant.toString() produces)
	t, err := time.Parse(time.RFC3339Nano, raw)
	if err != nil {
		t, err = time.Parse(time.RFC3339, raw)
	}
	return t, err
}