-- FORK-001: persistent support ticket center.
-- Conversations are normalized into tickets/messages/read positions. Replies
-- are intentionally append-only; no API or foreign-key cascade deletes them.

CREATE TABLE IF NOT EXISTS tickets (
    id                  BIGSERIAL PRIMARY KEY,
    user_id             BIGINT NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    subject             VARCHAR(200) NOT NULL,
    category            VARCHAR(32) NOT NULL,
    status              VARCHAR(32) NOT NULL DEFAULT 'pending_admin',
    priority            VARCHAR(16) NOT NULL DEFAULT 'normal',
    related_request_id  VARCHAR(128) NOT NULL DEFAULT '',
    assignee_id         BIGINT REFERENCES users(id) ON DELETE SET NULL,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    last_message_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    closed_at           TIMESTAMPTZ,
    CONSTRAINT tickets_category_check CHECK (category IN ('account', 'billing', 'api_model', 'incident', 'feature', 'other')),
    CONSTRAINT tickets_status_check CHECK (status IN ('pending_admin', 'pending_user', 'resolved', 'closed')),
    CONSTRAINT tickets_priority_check CHECK (priority IN ('low', 'normal', 'high', 'urgent'))
);

CREATE INDEX IF NOT EXISTS idx_tickets_user_last_message ON tickets(user_id, last_message_at DESC, id DESC);
CREATE INDEX IF NOT EXISTS idx_tickets_status_last_message ON tickets(status, last_message_at DESC, id DESC);
CREATE INDEX IF NOT EXISTS idx_tickets_category ON tickets(category);
CREATE INDEX IF NOT EXISTS idx_tickets_priority ON tickets(priority);
CREATE INDEX IF NOT EXISTS idx_tickets_assignee ON tickets(assignee_id);
CREATE INDEX IF NOT EXISTS idx_tickets_related_request ON tickets(related_request_id) WHERE related_request_id <> '';

CREATE TABLE IF NOT EXISTS ticket_messages (
    id          BIGSERIAL PRIMARY KEY,
    ticket_id   BIGINT NOT NULL REFERENCES tickets(id) ON DELETE RESTRICT,
    sender_id   BIGINT NOT NULL,
    sender_type VARCHAR(16) NOT NULL,
    content     TEXT NOT NULL,
    request_id  VARCHAR(128) NOT NULL DEFAULT '',
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT ticket_messages_sender_type_check CHECK (sender_type IN ('user', 'admin'))
);

CREATE INDEX IF NOT EXISTS idx_ticket_messages_ticket_id ON ticket_messages(ticket_id, id);
CREATE INDEX IF NOT EXISTS idx_ticket_messages_sender ON ticket_messages(sender_type, id);
CREATE INDEX IF NOT EXISTS idx_ticket_messages_request_id ON ticket_messages(request_id) WHERE request_id <> '';

CREATE TABLE IF NOT EXISTS ticket_reads (
    id                   BIGSERIAL PRIMARY KEY,
    ticket_id            BIGINT NOT NULL REFERENCES tickets(id) ON DELETE RESTRICT,
    reader_id            BIGINT NOT NULL,
    last_read_message_id BIGINT NOT NULL DEFAULT 0,
    updated_at           TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT ticket_reads_ticket_reader_unique UNIQUE(ticket_id, reader_id)
);

CREATE INDEX IF NOT EXISTS idx_ticket_reads_reader ON ticket_reads(reader_id, ticket_id);

