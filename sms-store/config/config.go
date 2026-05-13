package config

import (
	"os"
	"strings"
)

type AppConfig struct {
	ServerPort string
	MongoURI string
	MongoDBName  string
	KafkaBrokers []string
	KafkaTopic string
	KafkaGroupID string
}

func Load() AppConfig {
	return AppConfig{
		ServerPort:   getEnv("SERVER_PORT","8081"),
		MongoURI:  getEnv("MONGO_URI","mongodb://localhost:27017"),
		MongoDBName: getEnv("MONGO_DB_NAME","sms_store"),
		KafkaBrokers: strings.Split(getEnv("KAFKA_BOOTSTRAP_SERVERS","localhost:9092"), ","),
		KafkaTopic:getEnv("KAFKA_TOPIC", "sms-events"),
		KafkaGroupID: getEnv("KAFKA_GROUP_ID","sms-store-consumer-group"),
	}
}
func getEnv(key, defaultValue string) string{
	if v := os.Getenv(key); v != ""{
		return v
	}
	return defaultValue
}