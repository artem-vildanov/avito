CREATE TABLE bid_rollback (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    bid_id UUID NOT NULL,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    status bid_status NOT NULL,
    tender_id UUID REFERENCES tender(id) ON DELETE CASCADE,
    author_id UUID NOT NULL,
    author_type author_type not null,
    version INT NOT NULL,
    created_at TIMESTAMP NOT NULL
);
