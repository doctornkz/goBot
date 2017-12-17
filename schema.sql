CREATE TABLE "num_messages" (
    "userid" INTEGER NOT NULL PRIMARY KEY,
    "count" INTEGER
);
CREATE TABLE "user" (
    "userid" INTEGER,
    "username" TEXT NOT NULL DEFAULT ('""'),
    "firstname" TEXT,
    "lastname" TEXT
);
CREATE TABLE "messages" (
    "userid" INTEGER NOT NULL DEFAULT (0),
    "date" INTEGER NOT NULL DEFAULT (0),
    "text" TEXT NOT NULL
);
CREATE TABLE "categories" (
    "categoryid" INTEGER NOT NULL DEFAULT (0),
    "name" TEXT NOT NULL DEFAULT (' ')
);
CREATE TABLE words (
    "word" TEXT NOT NULL PRIMARY KEY DEFAULT (' '),
    "categoryid" INTEGER NOT NULL DEFAULT (0)
, "userid" INTEGER  NOT NULL  DEFAULT (0));
;
CREATE INDEX category_idx1 ON words(categoryid);
