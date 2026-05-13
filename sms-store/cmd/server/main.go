package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sms/sms-store/config"
	"github.com/sms/sms-store/internal/handler"
	"github.com/sms/sms-store/internal/repository"
	"github.com/sms/sms-store/internal/service"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func main(){
	logger := slog.New(slog.NewJSONHandler(os.Stdout,&slog.HandlerOptions{Level:slog.LevelDebug}))
	slog.SetDefault(logger)
	cfg := config.Load()
	slog.Info("Starting sms-store service","port",cfg.ServerPort)
	mongoClient, err := connectMongo(cfg.MongoURI)
	if err !=nil{
		slog.Error("Failed to connect to MongoDB","err",err)
		os.Exit(1)
	}
	defer func(){
		ctx,cancel := context.WithTimeout(context.Background(),5*time.Second)
		defer cancel()
		_ = mongoClient.Disconnect(ctx)
	}()
	db := mongoClient.Database(cfg.MongoDBName)
	repo,err := repository.NewMongoSmsRepository(db)
	if err != nil {
		slog.Error("Failed to initialise repository", "err", err)
		os.Exit(1)
	}
	svc := service.NewSmsService(repo)
	h:= handler.NewSmsHandler(svc)
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)
	consumer := config.NewKafkaConsumer(cfg.KafkaBrokers,cfg.KafkaTopic,cfg.KafkaGroupID,svc)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go consumer.Start(ctx)
	srv := &http.Server{
		Addr: fmt.Sprintf(":%s", cfg.ServerPort),
		Handler:loggingMiddleware(mux),
		ReadTimeout:15*time.Second,
		WriteTimeout:15* time.Second,
		IdleTimeout: 60* time.Second,
	}
	go func(){
		slog.Info("HTTP server listening","addr",srv.Addr)
		if err := srv.ListenAndServe();err != nil && err != http.ErrServerClosed {
			slog.Error("HTTP server error","err",err)
			os.Exit(1)
		}
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit,syscall.SIGINT,syscall.SIGTERM)
	<-quit
	slog.Info("Shutdown signal received, draining...")
	cancel()
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(),10*time.Second)
	defer shutdownCancel()
	if err := srv.Shutdown(shutdownCtx); err != nil{
		slog.Error("HTTP server forced to shutdown","err",err)
	}
	slog.Info("sms-store stopped cleanly")
}
func connectMongo(uri string) (*mongo.Client, error){
	ctx, cancel := context.WithTimeout(context.Background(),15*time.Second)
	defer cancel()
	client,err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, fmt.Errorf("mongo.Connect: %w", err)
	}
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, fmt.Errorf("mongo ping failed: %w", err)
	}
	slog.Info("Connected to MongoDB", "uri", uri)
	return client, nil
}
func loggingMiddleware(next http.Handler) http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter,r*http.Request){
		start := time.Now()
		slog.Info("Incoming request",
			"method",r.Method,
			"path",r.URL.Path,
			"remoteAddr", r.RemoteAddr,
		)
		next.ServeHTTP(w, r)
		slog.Info("Request completed",
			"method", r.Method,
			"path", r.URL.Path,
			"duration", time.Since(start).String(),
		)
	})}