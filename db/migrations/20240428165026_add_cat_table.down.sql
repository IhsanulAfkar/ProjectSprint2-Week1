-- Drop foreign key constraint
ALTER TABLE "cat" DROP CONSTRAINT "cat_userId_fkey";

-- Drop the cat table
DROP TABLE cat;

-- Drop the Gender ENUM type
DROP TYPE Gender;

-- Drop the CatRace ENUM type
DROP TYPE CatRace;