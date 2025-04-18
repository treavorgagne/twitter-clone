use twitter;

DROP VIEW IF EXISTS users_stats;
DROP VIEW IF EXISTS tweets_stats;
DROP VIEW IF EXISTS comments_stats;
DROP TABLE IF EXISTS likes_comment;
DROP TABLE IF EXISTS likes_tweet;
DROP TABLE IF EXISTS comments;
DROP TABLE IF EXISTS tweets;
DROP TABLE IF EXISTS follows;
DROP TABLE IF EXISTS users;

-- users
CREATE TABLE users (
   user_id BIGINT NOT NULL AUTO_INCREMENT COMMENT "unique id for user",
   username varchar(36) NOT NULL CHECK (username <> '') COMMENT "the username for a user",
   created_at datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT "date when user was added",
   PRIMARY KEY(user_id)
);

-- follows
CREATE TABLE follows (
   user_id BIGINT NOT NULL COMMENT "unique id the person whos following another user",
   follow_id BIGINT NOT NULL COMMENT "unique id of the user being followed",
   created_at datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT "date when user was followed the other user",
   PRIMARY KEY (user_id, follow_id),
   FOREIGN KEY(user_id) references users(user_id) ON DELETE CASCADE,
   FOREIGN KEY(follow_id) references users(user_id) ON DELETE CASCADE,
  	CHECK (user_id <> follow_id)
);

-- user counts
CREATE VIEW users_stats AS
SELECT
   u.user_id,
   u.username,
   u.created_at,
   IFNULL(followers.total_followers, 0) AS total_followers,
   IFNULL(following.following_total, 0) AS following_total
FROM users u
LEFT JOIN (
   SELECT
      user_id,
      COUNT(*) AS following_total
   FROM follows
   GROUP BY user_id
) AS following ON u.user_id = following.user_id
LEFT JOIN (
   SELECT
      follow_id AS user_id,
      COUNT(*) AS total_followers
   FROM follows
   GROUP BY follow_id
) AS followers ON u.user_id = followers.user_id;

-- tweets
CREATE TABLE tweets (
	tweet_id BIGINT NOT NULL AUTO_INCREMENT COMMENT "unique id for the tweet",
	body varchar(256) NOT NULL CHECK (body <> '') COMMENT "tweet body",
	user_id BIGINT NOT NULL COMMENT "unique id for user who made the tweet",
	create_at datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT "date when tweet was published",
	PRIMARY KEY(tweet_id),
	FOREIGN KEY(user_id) references users(user_id) ON DELETE CASCADE
);

-- likes_tweet
CREATE TABLE likes_tweet (
	tweet_id BIGINT NOT NULL COMMENT "unique id for the tweet being liked",
	user_id BIGINT NOT NULL COMMENT "unique id of the user liking the tweet",
   PRIMARY KEY (tweet_id, user_id),
   FOREIGN KEY(tweet_id) references tweets(tweet_id) ON DELETE CASCADE,
   FOREIGN KEY(user_id) references users(user_id) ON DELETE CASCADE
);

-- tweets total likes
CREATE VIEW tweets_stats AS
SELECT
   t.tweet_id,
   t.body,
   t.user_id,
   t.create_at,
   IFNULL(likes_count.tweet_total_likes, 0) AS tweet_total_likes
FROM tweets t
LEFT JOIN (
   SELECT
      tweet_id,
      COUNT(*) AS tweet_total_likes
   FROM likes_tweet
   GROUP BY tweet_id
) AS likes_count ON t.tweet_id = likes_count.tweet_id;

-- comments
CREATE TABLE comments (
   comment_id BIGINT NOT NULL AUTO_INCREMENT COMMENT 'unique id for the comment',
   tweet_id BIGINT NOT NULL COMMENT 'unique id for the tweet',
   body VARCHAR(256) NOT NULL CHECK (body <> '') COMMENT 'tweet body',
   user_id BIGINT NOT NULL COMMENT 'unique id for user who made the comment',
   create_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'date when comment was published',
   PRIMARY KEY (comment_id),
   FOREIGN KEY (tweet_id) REFERENCES tweets(tweet_id) ON DELETE CASCADE,
   FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE
);

-- likes_comment
CREATE TABLE likes_comment (
	comment_id BIGINT NOT NULL COMMENT "unique id for the comment being liked",
	user_id BIGINT NOT NULL COMMENT "unique id of the user liking the tweet comment",
   PRIMARY KEY (comment_id, user_id),
   FOREIGN KEY(comment_id) references comments(comment_id) ON DELETE CASCADE,
   FOREIGN KEY(user_id) references users(user_id) ON DELETE CASCADE
);

-- comments total likes
DROP VIEW IF EXISTS comments_stats;
CREATE VIEW comments_stats AS
SELECT
   c.comment_id,
   c.tweet_id,
   c.body,
   c.user_id,
   c.create_at,
   IFNULL(likes_count.comment_total_likes, 0) AS comment_total_likes
FROM comments c
LEFT JOIN (
   SELECT
      comment_id,
      COUNT(*) AS comment_total_likes
   FROM likes_comment
   GROUP BY comment_id
) AS likes_count ON c.comment_id = likes_count.comment_id;

-- run optional to populate data
insert into users (username) values ("tmoney"), ("tdawg"), ("tbone");
insert into follows (user_id, follow_id) values (1,2), (1,3), (2,1), (2,3), (3,1), (3,2);
insert into tweets (body, user_id) values ("tmoney's first tweet", 1), ("tdawg's first tweet", 2), ("tbone's first tweet", 1);
insert into likes_tweet (tweet_id, user_id) values (1, 2), (1, 3), (2,1), (2,3), (3,1), (3,2);
insert into comments (tweet_id, body, user_id) values (1, "nice post", 3), (2, "nice post", 1), (3, "nice post", 2);
insert into likes_comment (comment_id, user_id) values (1, 2), (2, 3), (3, 1);