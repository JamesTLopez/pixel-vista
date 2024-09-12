CREATE TABLE
    IF NOT EXISTS credit_prices (
        id serial primary key,
        product_id text null,
        name text,
        price text,
        created_at timestamp not null default now ()
    )