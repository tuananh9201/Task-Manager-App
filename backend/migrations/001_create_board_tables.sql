CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    name VARCHAR(100),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS boards (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    owner_id INTEGER REFERENCES users(id),
    is_public BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS board_members (
    board_id INTEGER REFERENCES boards(id) ON DELETE CASCADE,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    role VARCHAR(50) DEFAULT 'member',  -- member/admin/owner
    PRIMARY KEY (board_id, user_id)
);

CREATE TABLE IF NOT EXISTS lists (
    id SERIAL PRIMARY KEY,
    board_id INTEGER REFERENCES boards(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    position INTEGER NOT NULL,  -- For ordering lists
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS cards (
    id SERIAL PRIMARY KEY,
    list_id INTEGER REFERENCES lists(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    position INTEGER NOT NULL,  -- For ordering cards
    due_date TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Card Assignments
CREATE TABLE IF NOT EXISTS card_members (
    card_id INTEGER REFERENCES cards(id) ON DELETE CASCADE,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    PRIMARY KEY (card_id, user_id)
);

CREATE TABLE IF NOT EXISTS labels (
    id SERIAL PRIMARY KEY,
    board_id INTEGER REFERENCES boards(id) ON DELETE CASCADE,
    name VARCHAR(100),
    color VARCHAR(20)
);

CREATE TABLE IF NOT EXISTS comments (
    id SERIAL PRIMARY KEY,
    card_id INTEGER REFERENCES cards(id) ON DELETE CASCADE,
    user_id INTEGER REFERENCES users(id),
    content TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS attachments (
    id SERIAL PRIMARY KEY,
    card_id INTEGER REFERENCES cards(id) ON DELETE CASCADE,
    filename VARCHAR(255),
    url TEXT NOT NULL,
    uploaded_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
