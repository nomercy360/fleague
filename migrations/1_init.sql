PRAGMA foreign_keys = ON;

-- Таблица пользователей
CREATE TABLE users
(
    id            INTEGER PRIMARY KEY AUTOINCREMENT,
    first_name    TEXT,
    last_name     TEXT,
    username      TEXT,
    language_code TEXT,
    chat_id       INTEGER UNIQUE, -- ID чата в телеграме
    created_at    DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Таблица матчей
CREATE TABLE matches
(
    id         INTEGER PRIMARY KEY AUTOINCREMENT,
    tournament TEXT     NOT NULL,        -- Турнир, к которому относится матч
    home_team  TEXT     NOT NULL,
    away_team  TEXT     NOT NULL,
    match_date DATETIME NOT NULL,
    status     TEXT DEFAULT 'scheduled', -- Статусы: scheduled, ongoing, completed
    home_score INTEGER,                  -- Заполняется после завершения
    away_score INTEGER                   -- Заполняется после завершения
);

-- Таблица прогнозов
CREATE TABLE predictions
(
    id               INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id          TEXT    NOT NULL,
    match_id         INTEGER NOT NULL,
    predicted_winner TEXT    NOT NULL,   -- Прогноз на победителя (home/away/draw)
    points_awarded   INTEGER  DEFAULT 0, -- Очки за прогноз
    created_at       DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
    FOREIGN KEY (match_id) REFERENCES matches (id) ON DELETE CASCADE
);

-- Таблица лиг
CREATE TABLE leagues
(
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    name        TEXT NOT NULL,      -- Название лиги
    description TEXT,               -- Описание (опционально)
    owner_id    TEXT NOT NULL,      -- Владелец лиги
    is_active   BOOLEAN  DEFAULT 1, -- Активна ли лига
    created_at  DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (owner_id) REFERENCES users (id) ON DELETE CASCADE
);

-- Участники лиг
CREATE TABLE league_members
(
    id        INTEGER PRIMARY KEY AUTOINCREMENT,
    league_id INTEGER NOT NULL,
    user_id   TEXT    NOT NULL,
    joined_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (league_id) REFERENCES leagues (id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);

-- Связь матчей с лигами
CREATE TABLE league_matches
(
    id        INTEGER PRIMARY KEY AUTOINCREMENT,
    league_id INTEGER NOT NULL,
    match_id  INTEGER NOT NULL,
    FOREIGN KEY (league_id) REFERENCES leagues (id) ON DELETE CASCADE,
    FOREIGN KEY (match_id) REFERENCES matches (id) ON DELETE CASCADE
);

-- Таблица лидербордов
CREATE TABLE leaderboards
(
    id        INTEGER PRIMARY KEY AUTOINCREMENT,
    league_id INTEGER NOT NULL,
    user_id   TEXT    NOT NULL,
    points    INTEGER DEFAULT 0,
    FOREIGN KEY (league_id) REFERENCES leagues (id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);

-- Таблица сезонов
CREATE TABLE seasons
(
    id         INTEGER PRIMARY KEY AUTOINCREMENT,
    name       TEXT     NOT NULL, -- Название сезона
    start_date DATETIME NOT NULL,
    end_date   DATETIME NOT NULL,
    is_active  BOOLEAN DEFAULT 1  -- Активен ли сезон
);

-- Привязка матчей к сезонам
CREATE TABLE season_matches
(
    id        INTEGER PRIMARY KEY AUTOINCREMENT,
    season_id INTEGER NOT NULL,
    match_id  INTEGER NOT NULL,
    FOREIGN KEY (season_id) REFERENCES seasons (id) ON DELETE CASCADE,
    FOREIGN KEY (match_id) REFERENCES matches (id) ON DELETE CASCADE
);