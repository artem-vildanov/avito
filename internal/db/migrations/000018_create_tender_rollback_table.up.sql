CREATE TABLE tender_rollback (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tender_id UUID NOT NULL,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    status tender_status NOT NULL,
    service_type service_type NOT NULL,
    creator_username VARCHAR(50) REFERENCES employee(username) ON DELETE CASCADE,
    organization_id UUID REFERENCES organization(id) ON DELETE CASCADE,
    version INT NOT NULL,
    created_at TIMESTAMP not null
);
