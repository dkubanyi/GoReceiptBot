DROP TABLE IF EXISTS "public"."users";
CREATE TABLE "public"."users"
(
    "id"         uuid         NOT NULL,
    "user_id"    varchar(255) NOT NULL,
    "chat_id"    varchar(255) NOT NULL,
    "username"   varchar(255) NOT NULL,
    "first_name" varchar(255),
    "last_name"  varchar(255),
    PRIMARY KEY ("id")
);

CREATE UNIQUE INDEX "public"."idx_user_id_chat_id" ON users (user_id, chat_id);