CREATE TABLE
    IF NOT EXISTS images (
        id serial primary key,
        user_id uuid references auth.users,
        status int not null default 1,
        prompt text not null,
        image_url text,
        batch_id uuid not null,
        deleted boolean not null default 'false',
        created_at timestamp not null default now (),
        deleted_at timestamp
    )