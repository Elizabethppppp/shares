CREATE TABLE news(
                     id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
                     title VARCHAR(50) NOT NULL,
                     content TEXT NOT NULL,
                     is_active BOOLEAN NOT NULL DEFAULT false,
                     published_at DATE NOT NULL,
                     created_at  TIMESTAMPTZ(0) NOT NULL DEFAULT CURRENT_TIMESTAMP,
                     updated_at TIMESTAMPTZ(0) DEFAULT NULL,
                     preview_img VARCHAR(255),
                     img	VARCHAR(255)
);

CREATE OR REPLACE FUNCTION update_timestamp_column()
    RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER set_timestamp
    BEFORE UPDATE ON news
    FOR EACH ROW
EXECUTE FUNCTION update_timestamp_column();

