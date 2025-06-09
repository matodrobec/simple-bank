-- Active: 1742649766239@@127.0.0.1@3306@bank

-- read phenomena:
-- - dirty read
--   a transaction reads data written other concurrent uncommited transaction
-- - non-repeatble read
--   read the same row twice and see different value because modified by other commited transsaction
-- - phantom read
--   re-execute a query to find rows  and sees a different set of rows
-- - serialization anomaly
--   the result of a group of concurrent commited transactions is impossibel to achieve

-- transaction isolation:
--   - read UNCOMMITTED
--   - read COMMITTED
--   - REPEATABLE READ
          -- NO dirty read
          -- NO non-repeatble read
          -- NO phantom read
--   - SERIALIZABLE
          -- NO dirty read
          -- NO non-repeatble read
          -- NO phantom read
          -- NO serialization anomaly

DELETE from accounts;
INSERT INTO accounts VALUES (null, "test", 100, "usd", now()),
  (null, "test2", 100, "usd", now()),
  (null, "test3", 100, "usd", now());

-- docker exec -it bank-mysql-1 mysql -uroot -ptest bank

----------------------------------------------
--
-- read UNCOMMITTED
--

-- session 1
select @@transaction_isolation;
select @@global.transaction_isolation;

SET SESSION TRANSACTION ISOLATION LEVEL READ UNCOMMITTED;
select @@transaction_isolation;

-- session 2
SET SESSION TRANSACTION ISOLATION LEVEL READ UNCOMMITTED;
select @@transaction_isolation;
begin;

-- session 1
begin;
select * from accounts;


-- session 2
select * from accounts where id =1;

-- session 1
update accounts set balance = balance - 10 where id =1;

-- session 2
select * from accounts where id =1;


-- balance in both are 90 it is dearty read




----------------------------------------------
--
-- read COMMITTED
--

-- session 1
SET SESSION TRANSACTION ISOLATION LEVEL READ COMMITTED;
select @@transaction_isolation;

-- session 2
SET SESSION TRANSACTION ISOLATION LEVEL READ COMMITTED;
select @@transaction_isolation;
begin;

-- session 1
begin;
select * from accounts;

-- session 2
select * from accounts where id =1;

-- session 1
update accounts set balance = balance - 10 where id =1;

-- session 2
select * from accounts where id =1;
-- dirty read is OK
select * from accounts where balance >= 90;

-- session 1
commit;

-- session 2
select * from accounts where id =1;
-- non-repeatble phenomena read bacause deffirent balance value
select * from accounts where balance >= 90;
-- phantom read phenomena because show different rows



----------------------------------------------
--
-- repeatable read
--

-- session 1
SET SESSION TRANSACTION ISOLATION LEVEL REPEATABLE READ;
select @@transaction_isolation;

-- session 2
SET SESSION TRANSACTION ISOLATION LEVEL REPEATABLE READ;
select @@transaction_isolation;
begin;

-- session 1
begin;
select * from accounts;

-- session 2
select * from accounts where id =1;

-- session 1
update accounts set balance = balance - 10 where id =1;

-- session 2
select * from accounts where id =1;
-- dirty read is OK
select * from accounts where balance >= 80;

-- session 1
commit;

-- session 2
select * from accounts where id =1;
-- non-repeatble phenomena OK
select * from accounts where balance >= 80;
-- phantom read phenomena  OK
update accounts set balance = balance - 10 where id =1;
select * from accounts where id =1;
-- serialization anomaly  because balance should be 10 dolar less but it is 20 becasue session 1


----------------------------------------------
--
-- SERIALIZABLE LEVEL
--


-- session 1
SET SESSION TRANSACTION ISOLATION LEVEL SERIALIZABLE;
select @@transaction_isolation;

-- session 2
SET SESSION TRANSACTION ISOLATION LEVEL SERIALIZABLE;
select @@transaction_isolation;
begin;

-- session 1
begin;
select * from accounts;

-- session 2
select * from accounts where id =1;

-- session 1
update accounts set balance = balance - 10 where id =1;
-- wait to session 2
--ERROR 1205 (HY000): Lock wait timeout exceeded; try restarting transaction

