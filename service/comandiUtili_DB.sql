-- SQLite

INSERT INTO users (username, biography) VALUES("paper23", "");
INSERT INTO users (username, biography) VALUES("sapienza", "");
INSERT INTO users (username, biography) VALUES("provaa", "");

INSERT INTO bans (idUser, idBanned) VALUES(3,2);

SELECT * FROM users;
SELECT * FROM images;
SELECT * FROM follows;

SELECT images.idImage, users.idUser, users.username
FROM images
JOIN users ON images.idOwner = users.idUser
JOIN follows ON users.idUser = follows.idFollowed
WHERE follows.idFollower = 1
ORDER BY images.dateTime DESC
LIMIT 100;