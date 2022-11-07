/* Project 3 */


/* Question 1 */
create table password
(
    passwordId int NOT NULL AUTO_INCREMENT,
    changedDate date default now(),
    password varchar(255),
    active boolean,
    userId int NOT NULL,
    FOREIGN KEY (userId)
        REFERENCES user (userId)
        ON DELETE CASCADE,
    PRIMARY KEY (passwordId)
);

/* Question 2 */
create table user
(
    userId int NOT NULL AUTO_INCREMENT,
    firstName varchar(100),
    lastName varchar(100),
    city varchar(100),
    age int,
    PRIMARY KEY (userId)
);

/* Question 3 */
select * from password where active != 0;

/* Question 4 + 5 */
START TRANSACTION;

update password set active = 0 where active = 1 and userId=1234;

insert into password(password, active, userId) values ('myEncryptedPassword', 1, 1234);

COMMIT;
