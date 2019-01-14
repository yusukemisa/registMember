package main

import (
	"context"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/yusukemisa/registMember/src/user"

	"github.com/yusukemisa/registMember/src/infra"
	"go.mercari.io/datastore"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	ctx := context.Background()
	datastoreClient, err := infra.NewDatastoreClient(ctx,
		datastore.WithProjectID(os.Getenv("ProjectID")),
		datastore.WithCredentialsFile(os.Getenv("GCP_CREDENTIALS")),
	)
	if err != nil {
		log.Fatal("failed to make new datastore client")
	}
	if err := user.RegistUser(ctx, datastoreClient, makeUser()); err != nil {
		log.Fatal("failed to create user:", err.Error())
	}
}

func makeUser() *user.User {
	u := &user.User{
		UserName:   familyNames[rand.Intn(16)+1] + firstNames[rand.Intn(16)+1],
		PassWord:   "12345678",
		Status:     user.GREEN,
		Registered: time.Now(),
	}
	if u.UserName == "月ノ美兎" {
		u.Status = user.BLACK
	}
	return u
}

var familyNames = map[int]string{
	1:  "静",
	2:  "勇気",
	3:  "月ノ",
	4:  "モイラと",
	5:  "家長",
	6:  "剣持",
	7:  "伏見",
	8:  "ギルザレン",
	9:  "鈴鹿",
	10: "鈴谷",
	11: "物述",
	12: "ミライ",
	13: "キズナ",
	14: "渋谷",
	15: "もこ田",
	16: "樋口",
}

var firstNames = map[int]string{
	1:  "凛",
	2:  "ちひろ",
	3:  "美兎",
	4:  "える",
	5:  "むぎ",
	6:  "力也",
	7:  "ガク",
	8:  "3世",
	9:  "詩子",
	10: "アキ",
	11: "アリス",
	12: "アカリ",
	13: "アイ",
	14: "ハジメ",
	15: "めめめ",
	16: "楓",
}
