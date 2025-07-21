CREATE TABLE subscriptions(
    id BIGSERIAL PRIMARY KEY,
    user_id UUID NOT NULL,
    service_name VARCHAR NOT NULL,
    price INTEGER CHECK (price > 0) NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE
);
