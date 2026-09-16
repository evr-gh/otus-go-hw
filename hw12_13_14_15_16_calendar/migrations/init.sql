CREATE SCHEMA IF NOT EXISTS calendar AUTHORIZATION "user";

CREATE TABLE IF NOT EXISTS calendar.events (
    id BIGSERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    description TEXT NOT NULL,
    "time" TIMESTAMPTZ NOT NULL,
    duration BIGINT NOT NULL,
    owner TEXT NOT NULL,
    notifyleadtime BIGINT NOT NULL,
    scheduled BOOLEAN NOT NULL DEFAULT FALSE
);
