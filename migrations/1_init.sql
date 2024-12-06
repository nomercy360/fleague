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

-- Таблица команд
CREATE TABLE teams
(
    id           INTEGER PRIMARY KEY AUTOINCREMENT,
    name         TEXT NOT NULL UNIQUE, -- Название команды
    short_name   TEXT,                 -- Короткое название команды (опционально)
    abbreviation TEXT,                 -- Аббревиатура команды
    country      TEXT,                 -- Страна, к которой относится команда
    crest_url    TEXT                  -- URL на логотип команды
);

-- Таблица матчей
CREATE TABLE matches
(
    id           INTEGER PRIMARY KEY AUTOINCREMENT,
    tournament   TEXT     NOT NULL,        -- Турнир, к которому относится матч
    home_team_id INTEGER  NOT NULL,        -- ID домашней команды
    away_team_id INTEGER  NOT NULL,        -- ID гостевой команды
    match_date   DATETIME NOT NULL,
    status       TEXT DEFAULT 'scheduled', -- Статусы: scheduled, ongoing, completed
    home_score   INTEGER,                  -- Заполняется после завершения
    away_score   INTEGER,                  -- Заполняется после завершения
    FOREIGN KEY (home_team_id) REFERENCES teams (id) ON DELETE CASCADE,
    FOREIGN KEY (away_team_id) REFERENCES teams (id) ON DELETE CASCADE
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

INSERT INTO teams (name, short_name, abbreviation, country, crest_url)
VALUES ('Bayer 04 Leverkusen', 'Leverkusen', 'B04', 'DE', ''),
       ('Borussia Dortmund', 'Dortmund', 'BVB', 'DE', ''),
       ('FC Bayern München', 'Bayern', 'FCB', 'DE', ''),
       ('VfB Stuttgart', 'Stuttgart', 'VFB', 'DE', ''),
       ('Arsenal FC', 'Arsenal', 'ARS', 'EN', ''),
       ('Aston Villa FC', 'Aston Villa', 'AVL', 'EN', ''),
       ('Liverpool FC', 'Liverpool', 'LIV', 'EN', ''),
       ('Manchester City FC', 'Man City', 'MCI', 'EN', ''),
       ('Club Atlético de Madrid', 'Atleti', 'ATL', 'ES', ''),
       ('FC Barcelona', 'Barça', 'FCB', 'ES', ''),
       ('Real Madrid CF', 'Real Madrid', 'RMA', 'ES', ''),
       ('AC Milan', 'Milan', 'MIL', 'IT', ''),
       ('Atalanta BC', 'Atalanta', 'ATA', 'IT', ''),
       ('Bologna FC 1909', 'Bologna', 'BOL', 'IT', ''),
       ('FC Internazionale Milano', 'Inter', 'INT', 'IT', ''),
       ('Juventus FC', 'Juventus', 'JUV', 'IT', ''),
       ('Girona FC', 'Girona', 'GIR', 'ES', ''),
       ('Sporting Clube de Portugal', 'Sporting CP', 'SPO', 'PT', ''),
       ('Stade Brestois 29', 'Brest', 'BRE', 'FR', ''),
       ('Lille OSC', 'Lille', 'LIL', 'FR', ''),
       ('Paris Saint-Germain FC', 'PSG', 'PSG', 'FR', ''),
       ('AS Monaco FC', 'Monaco', 'ASM', 'MC', ''),
       ('PSV', 'PSV', 'PSV', 'NL', ''),
       ('Feyenoord Rotterdam', 'Feyenoord', 'FEY', 'NL', ''),
       ('RB Leipzig', 'RB Leipzig', 'RBL', 'DE', ''),
       ('Celtic FC', 'Celtic', 'CEL', 'SC', ''),
       ('GNK Dinamo Zagreb', 'Dinamo Zagreb', 'DIN', 'HR', ''),
       ('Club Brugge KV', 'Club Brugge', 'CLU', 'BE', ''),
       ('AC Sparta Praha', 'Sparta Praha', 'SPP', 'CZ', ''),
       ('BSC Young Boys', 'Young Boys', 'YOB', 'CH', ''),
       ('FC Red Bull Salzburg', 'RB Salzburg', 'RBS', 'AT', ''),
       ('FK Shakhtar Donetsk', 'Shaktar', 'SHD', 'UA', ''),
       ('Sport Lisboa e Benfica', 'SL Benfica', 'BEN', 'PT', ''),
       ('SK Sturm Graz', 'Sturm Graz', 'STU', 'AT', ''),
       ('FK Crvena Zvezda', 'Crvena Zvedza', 'CRV', 'RS', ''),
       ('ŠK Slovan Bratislava', 'Sl. Bratislava', 'SBA', 'SK', '');