CREATE TABLE IF NOT EXISTS lotto.profiles (
	id int NOT NULL AUTO_INCREMENT,
	uuid varchar(255) NOT NULL UNIQUE,
	firstname varchar(255) NOT NULL UNIQUE,
	lastname varchar(255) NOT NULL UNIQUE,
	email varchar(255) NOT NULL UNIQUE,
	phone varchar(255) NOT NULL UNIQUE,
	created_at datetime NULL,
	updated_at datetime NULL,
	PRIMARY KEY (id)
)
ENGINE=InnoDB
DEFAULT CHARSET=utf8mb3
COLLATE=utf8mb3_general_ci;

DROP TABLE lotto.profiles;

SELECT * FROM lotto.profiles p;

DELETE FROM lotto.profiles;

INSERT INTO lotto.profiles
(id, uuid, firstname, lastname, email, phone, created_at, updated_at)
VALUES(7, '4980f3a6fae54e5aa14780617bb2f045', 'firsttest', 'lasttest', 'l.firsttest@gmail.com', '0965462231', '2025-01-22 11:34:55', NULL);