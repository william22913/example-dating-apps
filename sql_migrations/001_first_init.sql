-- +migrate Up

-- +migrate StatementBegin

CREATE TYPE gender AS ENUM ('Male', 'Female', 'Any') ;

CREATE SEQUENCE IF NOT EXISTS "user_pkey_seq";

CREATE TABLE IF NOT EXISTS "user" (
    id BIGINT DEFAULT nextval('user_pkey_seq'::regclass),
    uuid_key UUID DEFAULT public.uuid_generate_v4(),
    phone_number VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    first_name VARCHAR(255),
    middle_name VARCHAR(255),
    last_name VARCHAR(255),
    birth_date TIMESTAMP WITHOUT TIME ZONE,
    gender gender,
    max_swipe INT DEFAULT 10,
    created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    created_by BIGINT DEFAULT 0,
    updated_by BIGINT DEFAULT 0,
    deleted BOOLEAN DEFAULT FALSE
);

ALTER TABLE "user" ADD CONSTRAINT pk_user_id PRIMARY KEY (id);
ALTER TABLE "user" ADD CONSTRAINT uq_user_phone_number UNIQUE (phone_number);

CREATE SEQUENCE IF NOT EXISTS "user_preferences_pkey_seq";

CREATE TABLE IF NOT EXISTS "user_preferences" (
    id BIGINT DEFAULT nextval('user_preferences_pkey_seq'::regclass),
    uuid_key UUID DEFAULT public.uuid_generate_v4(),
    user_id BIGINT NOT NULL,
    gender gender,
    min_age INT,
    max_age INT,
    created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    created_by BIGINT DEFAULT 0,
    updated_by BIGINT DEFAULT 0,
    deleted BOOLEAN DEFAULT FALSE
);

ALTER TABLE "user_preferences" ADD CONSTRAINT pk_user_preferences_id PRIMARY KEY (id);
ALTER TABLE "user_preferences" ADD CONSTRAINT fk_user_preferences_user_id_user_id FOREIGN KEY (user_id) REFERENCES "user"(id);

CREATE SEQUENCE IF NOT EXISTS "user_passions_pkey_seq";

CREATE TABLE IF NOT EXISTS "user_passions" (
    id BIGINT DEFAULT nextval('user_passions_pkey_seq'::regclass),
    uuid_key UUID DEFAULT public.uuid_generate_v4(),
    user_id BIGINT NOT NULL,
    tags TEXT[],
    created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    created_by BIGINT DEFAULT 0,
    updated_by BIGINT DEFAULT 0,
    deleted BOOLEAN DEFAULT FALSE
);

ALTER TABLE "user_passions" ADD CONSTRAINT pk_user_passions_id PRIMARY KEY (id);
Alter TABLE "user_passions" ADD CONSTRAINT fk_user_passions_user_id_user_id FOREIGN KEY (user_id) REFERENCES "user"(id);

CREATE SEQUENCE IF NOT EXISTS "user_premium_pkey_seq";

CREATE TABLE IF NOT EXISTS "user_premium" (
    id BIGINT DEFAULT nextval('user_premium_pkey_seq'::regclass),
    uuid_key UUID DEFAULT public.uuid_generate_v4(),
    user_id BIGINT NOT NULL,
    purchase_at TIMESTAMP WITHOUT TIME ZONE,
    price FLOAT8,
    ended_at TIMESTAMP WITHOUT TIME ZONE,
    created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    created_by BIGINT DEFAULT 0,
    updated_by BIGINT DEFAULT 0,
    deleted BOOLEAN DEFAULT FALSE
);

ALTER TABLE "user_premium" ADD CONSTRAINT pk_user_premium_id PRIMARY KEY (id);
ALTER TABLE "user_premium" ADD CONSTRAINT fk_user_premium_user_id_user_id FOREIGN KEY (user_id) REFERENCES "user"(id);

CREATE SEQUENCE IF NOT EXISTS "salt_pkey_seq";

CREATE TABLE IF NOT EXISTS "salt" (
    id BIGINT DEFAULT nextval('salt_pkey_seq'::regclass),
    uuid_key UUID DEFAULT public.uuid_generate_v4(),
    user_id BIGINT NOT NULL,
    salt_key VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    created_by BIGINT DEFAULT 0,
    updated_by BIGINT DEFAULT 0,
    deleted BOOLEAN DEFAULT FALSE
);

ALTER TABLE "salt" ADD CONSTRAINT pk_salt_id PRIMARY KEY (id);
ALTER TABLE "salt" ADD CONSTRAINT fk_salt_user_id_user_id FOREIGN KEY (user_id) REFERENCES "user"(id);

-- +migrate StatementEnd
