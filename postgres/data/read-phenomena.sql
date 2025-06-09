-- Active: 1742649723095@@127.0.0.1@5432@bank

DELETE FROM accounts;
SELECT pg_get_serial_sequence('accounts', 'id');
ALTER SEQUENCE accounts_id_seq RESTART WITH 1;

INSERT INTO accounts  (owner, balance, currency)
VALUES
  ('test', 100, 'usd'),
  ('test2', 100, 'usd'),
  ('test3', 100, 'usd');


-- transaction isolation:
--   - read UNCOMMITTED
--   - read COMMITTED
--   - REPEATABLE READ
--   - SERIALIZABLE
SHOW TRANSACTION ISOLATION LEVEL;

-------------
-- read UNCOMMITTED
--

-- session 1
begin;
SET TRANSACTION ISOLATION LEVEL READ UNCOMMITTED;
SHOW TRANSACTION ISOLATION LEVEL;


-- session 2
begin;
SET TRANSACTION ISOLATION LEVEL READ UNCOMMITTED;
SHOW TRANSACTION ISOLATION LEVEL;

-- session 1
SELECT * FROM accounts;

-- session 2
SELECT * FROM accounts WHERE id = 1;

-- session 1
update accounts set balance = balance - 10 where id =1 RETURNING *;

-- session 2
SELECT * FROM accounts WHERE id = 1;
-- here is note dirty read OK

-- session 1
commit;

-- session 2
SELECT * FROM accounts WHERE id = 1;
commit;
-- OK



-------------
-- read COMMITTED
--

-- session 1
begin;
SET TRANSACTION ISOLATION LEVEL READ COMMITTED;
SHOW TRANSACTION ISOLATION LEVEL;


-- session 2
begin;
SET TRANSACTION ISOLATION LEVEL READ COMMITTED;
SHOW TRANSACTION ISOLATION LEVEL;

-- session 1
SELECT * FROM accounts;

-- session 2
SELECT * FROM accounts WHERE id = 1;
SELECT * FROM accounts WHERE  balance >= 90;

-- session 1
update accounts set balance = balance - 10 where id =1 RETURNING *;

-- session 2
SELECT * FROM accounts WHERE id = 1;
-- here is note dirty read OK because session 1 not commited

-- session 1
commit;

-- session 2
SELECT * FROM accounts WHERE id = 1;
-- ok
SELECT * FROM accounts WHERE  balance >= 90;
-- phantom read
commit;





-------------
-- REPEATABLE READ;
--

-- session 1
begin;
SET TRANSACTION ISOLATION LEVEL REPEATABLE READ;
SHOW TRANSACTION ISOLATION LEVEL;


-- session 2
begin;
SET TRANSACTION ISOLATION LEVEL REPEATABLE READ;
SHOW TRANSACTION ISOLATION LEVEL;

-- session 1
SELECT * FROM accounts;

-- session 2
SELECT * FROM accounts WHERE id = 1;
SELECT * FROM accounts WHERE  balance >= 80;

-- session 1
update accounts set balance = balance - 10 where id =1 RETURNING *;
commit;

-- session 2
SELECT * FROM accounts WHERE id = 1;
--  non-repeatble read
SELECT * FROM accounts WHERE  balance >= 80;
-- phantom read prevented
update accounts set balance = balance - 10 where id =1 RETURNING *;
-- ERROR:  could not serialize access due to concurrent update
commit;



-------------
-- REPEATABLE READ;
--

-- session 1
begin;
SET TRANSACTION ISOLATION LEVEL REPEATABLE READ;
SHOW TRANSACTION ISOLATION LEVEL;


-- session 2
begin;
SET TRANSACTION ISOLATION LEVEL REPEATABLE READ;
SHOW TRANSACTION ISOLATION LEVEL;

-- session 1
SELECT * FROM accounts;
SELECT sum(ballance) FROM accounts;
INSERT INTO accounts  (owner, balance, currency) VALUES  ('sum', 270, 'usd');

-- session 2
SELECT * FROM accounts;
SELECT sum(ballance) FROM accounts;
INSERT INTO accounts  (owner, balance, currency) VALUES  ('sum', 270, 'usd');

-- session 1
commit;

-- session 2
COMMIT;
SELECT * FROM accounts;
-- serialization anomaly - two times 270 sum rows
