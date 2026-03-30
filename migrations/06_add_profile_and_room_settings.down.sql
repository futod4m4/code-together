ALTER TABLE users
    DROP COLUMN IF EXISTS avatar_url,
    DROP COLUMN IF EXISTS github_url,
    DROP COLUMN IF EXISTS bio;

ALTER TABLE rooms
    DROP COLUMN IF EXISTS description,
    DROP COLUMN IF EXISTS is_private;
