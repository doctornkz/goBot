CREATE TABLE "num_messages" (
    "userid" INTEGER NOT NULL PRIMARY KEY,
    "count" INTEGER
);
CREATE TABLE user (
    "userid" INTEGER,
    "username" TEXT NOT NULL DEFAULT ('""'),
    "firstname" TEXT,
    "lastname" TEXT
);
