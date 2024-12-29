package shared

import (
	"math280h/wisp/db"
)

var DBClient *db.PrismaClient //nolint:gochecknoglobals // This is the database client

func InitDB() {
	DBClient = db.NewClient()
	if err := DBClient.Prisma.Connect(); err != nil {
		panic(err)
	}
}
