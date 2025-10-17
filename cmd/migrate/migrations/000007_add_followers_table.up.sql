-- Any user_id can follow the follower_id and vice versa
CREATE TABLE IF NOT EXISTS followers (
    user_id bigint NOT NULL,
    follower_id bigint NOT NULL,
    created_at timestamp(0) with time zone NOT NULL DEFAULT now(),

    PRIMARY KEY (user_id, follower_id), -- composite primary key (two or more columns, common in many-to-many relationships) this also means that a user cannot follow the same user twice since a table cannot have duplicate primary keys
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE, -- on delete cascade is not ideal for soft deletes
    FOREIGN KEY (follower_id) REFERENCES users (id) ON DELETE CASCADE
); 