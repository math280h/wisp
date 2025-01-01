package shared

import (
	"context"
	"math280h/wisp/db"

	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog/log"
)

var DBClient *db.PrismaClient //nolint:gochecknoglobals // This is the database client

func InitDB() {
	DBClient = db.NewClient()
	if err := DBClient.Prisma.Connect(); err != nil {
		panic(err)
	}
}

func GetUserIfExists(author *discordgo.User) (*db.UserModel, error) {
	userObj, err := DBClient.User.UpsertOne(
		db.User.UserID.Equals(author.ID),
	).Create(
		db.User.UserID.Set(author.ID),
		db.User.Nickname.Set(author.Username),
	).Update().Exec(context.Background())
	if err != nil {
		log.Error().Err(err).Msg("Failed to create user")
		return nil, err
	}

	return userObj, nil
}
