
CREATE TABLE "matcher" (
  "pkId" SERIAL PRIMARY KEY,
 "id" uuid UNIQUE NOT NULL DEFAULT (uuid_generate_v4()),
  "userId" Int NOT NULL,
  "matchUserId" Int NOT NULL,
  "userCatId" Int NOT NULL,
  "matchCatId" Int NOT NULL,
  "message" varchar(120) NOT NULL,
  "isApproved" bool DEFAULT false,
  "isValid" bool DEFAULT true,
  "created_at" timestamp NOT NULL DEFAULT (now()),
  "updated_at" timestamp NOT NULL DEFAULT (now())
);

ALTER TABLE "matcher" ADD FOREIGN KEY ("userId") REFERENCES "user" ("pkId");
-- workaround
-- ALTER TABLE "matcher" ADD FOREIGN KEY ("userCatId") REFERENCES "cat" ("pkId");

-- ALTER TABLE "matcher" ADD FOREIGN KEY ("matchCatId") REFERENCES "cat" ("pkId");
