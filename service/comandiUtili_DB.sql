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

SELECT COUNT(*) FROM follows WHERE idFollower = 3;

SELECT file from images WHERE idImage IN (1,2,3,4)

DELETE FROM images WHERE idImage = 2

INSERT INTO likes (idLiker, idImageLiked) VALUES(2,1)

SELECT * FROM likes
SELECT * FROM comments
SELECT * FROM users

SELECT COUNT (*) FROM likes AS l, images AS i WHERE i.idOwner = 2 AND l.idImageLiked = i.idImage

SELECT * FROM comments WHERE idImageCommented = 6

INSERT INTO comments (idUserWriter, idImageCommented, text) VALUES(1, 1, "Commento da user 1")

SELECT idUser FROM users WHERE username = "abviuaerv"