CREATE TABLE IF NOT EXISTS userDB (
    user_id     SERIAL PRIMARY KEY, 
    user_name   VARCHAR(30),
    surname     VARCHAR(30),
    age         INT,
    gender      VARCHAR(30),
    country     VARCHAR(30)
);