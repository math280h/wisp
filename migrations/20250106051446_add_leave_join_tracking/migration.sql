-- RedefineTables
PRAGMA defer_foreign_keys=ON;
PRAGMA foreign_keys=OFF;
CREATE TABLE "new_infractions" (
    "id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    "user_id" INTEGER NOT NULL,
    "reason" TEXT NOT NULL,
    "type" TEXT NOT NULL,
    "points" INTEGER NOT NULL,
    "moderator_id" TEXT NOT NULL,
    "moderator_username" TEXT NOT NULL,
    "created_at" TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT "infractions_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE RESTRICT ON UPDATE CASCADE
);
INSERT INTO "new_infractions" ("created_at", "id", "moderator_id", "moderator_username", "points", "reason", "type", "user_id") SELECT "created_at", "id", "moderator_id", "moderator_username", "points", "reason", "type", "user_id" FROM "infractions";
DROP TABLE "infractions";
ALTER TABLE "new_infractions" RENAME TO "infractions";
CREATE TABLE "new_messages" (
    "id" TEXT NOT NULL PRIMARY KEY,
    "content" TEXT NOT NULL,
    "author_id" INTEGER NOT NULL,
    "channel_id" TEXT NOT NULL,
    "created_at" TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT "messages_author_id_fkey" FOREIGN KEY ("author_id") REFERENCES "users" ("id") ON DELETE RESTRICT ON UPDATE CASCADE
);
INSERT INTO "new_messages" ("author_id", "channel_id", "content", "created_at", "id", "updated_at") SELECT "author_id", "channel_id", "content", "created_at", "id", "updated_at" FROM "messages";
DROP TABLE "messages";
ALTER TABLE "new_messages" RENAME TO "messages";
CREATE UNIQUE INDEX "messages_id_key" ON "messages"("id");
CREATE TABLE "new_reports" (
    "id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    "channel_id" TEXT NOT NULL,
    "channel_name" TEXT NOT NULL,
    "user_id" INTEGER NOT NULL,
    "status" TEXT NOT NULL DEFAULT 'pending',
    "created_at" TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT "reports_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE RESTRICT ON UPDATE CASCADE
);
INSERT INTO "new_reports" ("channel_id", "channel_name", "created_at", "id", "status", "updated_at", "user_id") SELECT "channel_id", "channel_name", "created_at", "id", "status", "updated_at", "user_id" FROM "reports";
DROP TABLE "reports";
ALTER TABLE "new_reports" RENAME TO "reports";
CREATE TABLE "new_suggestion_votes" (
    "id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    "suggestion_id" INTEGER NOT NULL,
    "user_id" INTEGER NOT NULL,
    "sentiment" TEXT NOT NULL DEFAULT 'vote_up',
    "created_at" TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT "suggestion_votes_suggestion_id_fkey" FOREIGN KEY ("suggestion_id") REFERENCES "suggestions" ("id") ON DELETE RESTRICT ON UPDATE CASCADE,
    CONSTRAINT "suggestion_votes_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE RESTRICT ON UPDATE CASCADE
);
INSERT INTO "new_suggestion_votes" ("created_at", "id", "sentiment", "suggestion_id", "user_id") SELECT "created_at", "id", "sentiment", "suggestion_id", "user_id" FROM "suggestion_votes";
DROP TABLE "suggestion_votes";
ALTER TABLE "new_suggestion_votes" RENAME TO "suggestion_votes";
CREATE TABLE "new_suggestions" (
    "id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    "suggestion" TEXT NOT NULL,
    "user_id" INTEGER NOT NULL,
    "embed_id" TEXT,
    "status" TEXT NOT NULL DEFAULT 'pending',
    "created_at" TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT "suggestions_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE RESTRICT ON UPDATE CASCADE
);
INSERT INTO "new_suggestions" ("created_at", "embed_id", "id", "status", "suggestion", "user_id") SELECT "created_at", "embed_id", "id", "status", "suggestion", "user_id" FROM "suggestions";
DROP TABLE "suggestions";
ALTER TABLE "new_suggestions" RENAME TO "suggestions";
CREATE TABLE "new_users" (
    "id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    "user_id" TEXT NOT NULL,
    "nickname" TEXT NOT NULL,
    "points" INTEGER NOT NULL DEFAULT 0,
    "last_join" TEXT,
    "last_leave" TEXT,
    "leave_count" INTEGER NOT NULL DEFAULT 0,
    "created_at" TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
);
INSERT INTO "new_users" ("created_at", "id", "nickname", "points", "user_id") SELECT "created_at", "id", "nickname", "points", "user_id" FROM "users";
DROP TABLE "users";
ALTER TABLE "new_users" RENAME TO "users";
CREATE UNIQUE INDEX "users_user_id_key" ON "users"("user_id");
PRAGMA foreign_keys=ON;
PRAGMA defer_foreign_keys=OFF;
