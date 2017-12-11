CREATE TABLE user (
    "userid" INTEGER,
    "username" TEXT NOT NULL DEFAULT ('""'),
    "firstname" TEXT,
    "lastname" TEXT
);
CREATE TABLE "num_messages" (
    "userid" INTEGER NOT NULL PRIMARY KEY,
    "count" INTEGER
);
CREATE TABLE "messages" (
    "userid" INTEGER NOT NULL DEFAULT (0),
    "date" INTEGER NOT NULL DEFAULT (0),
    "text" TEXT NOT NULL
);
