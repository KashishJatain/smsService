package model
import (
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SmsStatus string
const (
	StatusSuccess SmsStatus = "SUCCESS"
	StatusFail SmsStatus= "FAIL"
	StatusBlocked SmsStatus= "BLOCKED"
)
type SmsRecord struct {
	ID             primitive.ObjectID `bson:"_id,omitempty"      json:"id,omitempty"`
	MessageID      string             `bson:"message_id"    json:"messageId"`
	UserID         string             `bson:"user_id"    json:"userId"`
	PhoneNumber    string             `bson:"phone_number"     json:"phoneNumber"`
	Message        string             `bson:"message"     json:"message"`
	Status         SmsStatus          `bson:"status"    json:"status"`
	VendorResponse string             `bson:"vendor_response"    json:"vendorResponse"`
	Timestamp      time.Time          `bson:"timestamp"   json:"timestamp"`
	CreatedAt      time.Time          `bson:"created_at"   json:"createdAt"`
}
type SmsEvent struct {
	MessageID     string `json:"messageId"`
	UserID       string `json:"userId"`
	PhoneNumber    string `json:"phoneNumber"`
	Message    string `json:"message"`
	Status   string `json:"status"`
	VendorResponse string `json:"vendorResponse"`
	Timestamp     string `json:"timestamp"`
}
type SmsHistoryResponse struct {
	UserID   string      `json:"userId"`
	Count   int       `json:"count"`
	Messages []SmsRecord `json:"messages"`
}
type ErrorResponse struct {
	Error  string `json:"error"`
	Message string `json:"message"`
}
