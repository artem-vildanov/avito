create table review (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    bid_id uuid references bid(id) on delete cascade,
    description text,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
)