CREATE TABLE "user" (
  "pkId" Int UNIQUE PRIMARY KEY NOT NULL,
  "id" uuid UNIQUE NOT NULL DEFAULT (uuid_generate_v4()),
  "email" varchar(255) UNIQUE NOT NULL,
  "password" varchar(255) NOT NULL,
  "accessToken" varchar(255) UNIQUE,
  "created_at" timestamp NOT NULL DEFAULT (now()),
  "updated_at" timestamp NOT NULL
);