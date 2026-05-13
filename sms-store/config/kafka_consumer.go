package config
import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"
	"github.com/segmentio/kafka-go"
	"github.com/sms/sms-store/internal/model"
	"github.com/sms/sms-store/internal/service"
)

type KafkaConsumer struct{
	reader *kafka.Reader
	svc  *service.SmsService
}

func NewKafkaConsumer(brokers []string, topic, groupID string,svc *service.SmsService) *KafkaConsumer{
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  brokers,
		Topic:  topic,
		GroupID: groupID,
		MinBytes: 1,
		MaxBytes: 10e6,
		MaxWait: 500* time.Millisecond,
		CommitInterval: time.Second,
		StartOffset: kafka.FirstOffset,
	})
	return &KafkaConsumer{reader:reader,svc:svc}
}
func (c *KafkaConsumer) Start(ctx context.Context){
	slog.Info("Kafka consumer started")
	defer func(){
		if err := c.reader.Close();err != nil {
			slog.Error("Error closing Kafka reader","err",err)
		}
		slog.Info("Kafka consumer is stopped")
	}()
	for {
		msg, err := c.reader.FetchMessage(ctx)
		if err != nil{
			if ctx.Err() != nil{
				return
			}
			slog.Error("Error fetching Kafka message","err", err)
			time.Sleep(time.Second)
			continue
		}
		if processErr := c.processMessage(ctx, msg); processErr !=nil{
			slog.Error("Failed to process message; skipping",
				"offset", msg.Offset,"err",processErr)
		}
		if commitErr := c.reader.CommitMessages(ctx, msg); commitErr != nil {
			slog.Error("Failed to commit Kafka offset", "err", commitErr)
		}
	}}
func (c *KafkaConsumer) processMessage(ctx context.Context, msg kafka.Message) error{
	slog.Debug("Received Kafka message",
		"topic", msg.Topic,
		"partition",msg.Partition,
		"offset",msg.Offset,
	)
	var event model.SmsEvent
	if err := json.Unmarshal(msg.Value, &event); err != nil {
		return fmt.Errorf("failed to unmarshal SmsEvent: %w",err)
	}
	return c.svc.StoreEvent(ctx, &event)
}