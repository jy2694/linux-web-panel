package database

import (
	"context"
	"log"
	"time"

	"github.com/linux-web-panel/backend/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

func ConnectMongoDB(uri string) (*mongo.Database, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	// 연결 테스트
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	db := client.Database("linux_web_panel")
	
	// 초기 관리자 계정 생성
	if err := createInitialAdmin(db); err != nil {
		log.Printf("Warning: Failed to create initial admin account: %v", err)
	}

	log.Println("Connected to MongoDB!")
	return db, nil
}

func createInitialAdmin(db *mongo.Database) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 기존 관리자 계정 확인
	var existingAdmin models.User
	err := db.Collection("users").FindOne(ctx, bson.M{"username": "admin"}).Decode(&existingAdmin)
	if err == nil {
		// 이미 관리자 계정이 존재함
		return nil
	}

	// 비밀번호 해싱
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("admin"), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// 관리자 계정 생성
	admin := models.User{
		Username:  "admin",
		Password:  string(hashedPassword),
		Role:      "admin",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	_, err = db.Collection("users").InsertOne(ctx, admin)
	return err
}
