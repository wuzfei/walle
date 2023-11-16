package repository

import (
	"context"
	"fmt"
	"github.com/wuzfei/go-helper/rand"
	"golang.org/x/crypto/bcrypt"
	"testing"
	"yema.dev/internal/model"
	"yema.dev/pkg/db"
	log2 "yema.dev/pkg/log"
)

var logConfig = log2.Config{
	//Encoder: "console",
	Level:       "debug",
	Output:      "console",
	Development: true,
}

var dbConfig = db.Config{
	Driver:   "sqlite3",
	Dsn:      "/Users/wuxin/worker/yema.dev/yema_dev.db",
	LogLevel: "info",
}

func _userRepo() UserRepository {
	log := log2.NewLog(&logConfig)
	db, err := db.NewDB(&dbConfig, log)
	if err != nil {
		panic(err)
	}
	return NewUserRepository(NewRepository(db, log))
}

func TestUserRepository(t *testing.T) {
	repo := _userRepo()
	ctx := context.Background()
	_pwd, _ := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)
	m := model.User{
		Username: "yema",
		Email:    fmt.Sprintf("e_%d@yema.dev", rand.Int(1000, 9999)),
		Password: _pwd,
		Status:   1,
	}
	err := repo.Create(ctx, &m)
	if err != nil {
		t.Fatal("create user error", err)
	}

	m.Username = "yema_update"
	m.Status = 0

	err = repo.Update(ctx, &m, "username", "status")
	if err != nil {
		t.Fatal("update user error", err)
	}

	_m, err := repo.GetByID(ctx, m.ID)
	if err != nil {
		t.Fatal("GetByID user error", err)
	}

	if _m.Username != "yema_update" || _m.Status != 0 {
		t.Fatal("update user check error", err)
	}

	_, _, err = repo.List(ctx, "yema")

	if err != nil {
		t.Fatal("List user error", err)
	}

}
