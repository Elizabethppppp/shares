CREATE TABLE news(
                     id SERIAL PRIMARY KEY,
                     title VARCHAR(50) NOT NULL,
                     content TEXT NOT NULL,
                     is_active BOOLEAN NOT NULL DEFAULT false,
                     published_at DATE NOT NULL,
                     created_at  TIMESTAMP(0) DEFAULT CURRENT_TIMESTAMP,
                     updated_at TIMESTAMP(0),
                     preview_img VARCHAR(255),
                     img	VARCHAR(255)
);