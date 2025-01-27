CREATE TABLE IF NOT EXISTS lotto.otp_user (
	id int NOT NULL AUTO_INCREMENT,
	uuid varchar(255) NOT NULL UNIQUE,
	otp varchar(255) NOT NULL,
	created_at datetime NULL,
	updated_at datetime NULL,
	PRIMARY KEY (id)
)
ENGINE=InnoDB
DEFAULT CHARSET=utf8mb3
COLLATE=utf8mb3_general_ci;

DROP TABLE lotto.otp_user;

SELECT * FROM lotto.otp_user;

INSERT INTO lotto.otp_user
(id, uuid, otp, created_at, updated_at)
VALUES(1, '4980f3a6fae54e5aa14780617bb2f045', '200139', '2025-01-22 16:56:45', '2025-01-22 16:56:45');