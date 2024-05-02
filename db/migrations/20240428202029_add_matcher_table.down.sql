-- Drop foreign keys
ALTER TABLE "matcher" DROP CONSTRAINT IF EXISTS "matcher_userId_fkey";
ALTER TABLE "matcher" DROP CONSTRAINT IF EXISTS "matcher_userCatId_fkey";
ALTER TABLE "matcher" DROP CONSTRAINT IF EXISTS "matcher_matchCatId_fkey";

-- Drop the table
DROP TABLE IF EXISTS "matcher";
