CREATE TABLE "users" (
    "id" BIGSERIAL PRIMARY KEY,
    "username" TEXT NOT NULL UNIQUE,
    "password_hash" TEXT NOT NULL,
    "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE "calls" (
    "id" BIGSERIAL PRIMARY KEY,
    "client_name" TEXT NOT NULL,
    "phone_number" TEXT CHECK (phone_number ~ '^(\+?\d{1,3}|\d)?[\d\-]{7,15}$'),
    "description" TEXT NOT NULL,
    "status" TEXT DEFAULT 'открыта' CHECK (status IN ('открыта', 'закрыта')),
    "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "user_id" BIGINT NOT NULL,
    CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);