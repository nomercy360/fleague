PRAGMA foreign_keys = ON;

-- Таблица пользователей
CREATE TABLE users
(
    id                  TEXT PRIMARY KEY,
    first_name          TEXT,
    last_name           TEXT,
    username            TEXT,
    language_code       TEXT,
    chat_id             INTEGER UNIQUE, -- ID чата в телеграме
    referral_code       TEXT UNIQUE,
    referred_by         TEXT,
    created_at          DATETIME DEFAULT CURRENT_TIMESTAMP,
    total_points        INTEGER  DEFAULT 0,
    total_predictions   INTEGER  DEFAULT 0,
    correct_predictions INTEGER  DEFAULT 0,
    avatar_url          TEXT,
    FOREIGN KEY (referred_by) REFERENCES users (id) ON DELETE SET NULL
);

CREATE INDEX idx_users_chat_id ON users (chat_id);
CREATE INDEX idx_users_username ON users (username);

CREATE TABLE user_followers
(
    follower_id  TEXT NOT NULL,
    following_id TEXT NOT NULL,
    created_at   DATETIME DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (follower_id, following_id),
    FOREIGN KEY (follower_id) REFERENCES users (id) ON DELETE CASCADE,
    FOREIGN KEY (following_id) REFERENCES users (id) ON DELETE CASCADE
);

-- Таблица команд
CREATE TABLE teams
(
    id           TEXT PRIMARY KEY,
    name         TEXT NOT NULL UNIQUE, -- Название команды
    short_name   TEXT,                 -- Короткое название команды (опционально)
    abbreviation TEXT,                 -- Аббревиатура команды
    country      TEXT,                 -- Страна, к которой относится команда
    crest_url    TEXT                  -- URL на логотип команды
);

-- Таблица матчей
CREATE TABLE matches
(
    id           TEXT PRIMARY KEY,
    tournament   TEXT     NOT NULL,        -- Турнир, к которому относится матч
    home_team_id TEXT,        -- ID домашней команды
    away_team_id TEXT,        -- ID гостевой команды
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
    user_id              TEXT,
    match_id             TEXT,
    predicted_outcome    TEXT CHECK (predicted_outcome IN ('home', 'away', 'draw')),
    predicted_home_score INTEGER,
    predicted_away_score INTEGER,
    points_awarded       INTEGER  DEFAULT 0,
    created_at           DATETIME DEFAULT CURRENT_TIMESTAMP,
    completed_at         DATETIME, -- Дата когда матч завершился и прогноз был подсчитан
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
    FOREIGN KEY (match_id) REFERENCES matches (id) ON DELETE CASCADE,
    PRIMARY KEY (user_id, match_id)
);

CREATE TABLE seasons
(
    id         TEXT PRIMARY KEY,
    name       TEXT     NOT NULL, -- Название сезона
    start_date DATETIME NOT NULL,
    end_date   DATETIME NOT NULL,
    is_active  BOOLEAN DEFAULT 1  -- Активен ли сезон
);

-- Таблица лидербордов
CREATE TABLE leaderboards
(
    user_id   TEXT NOT NULL,
    points    INTEGER DEFAULT 0,
    season_id TEXT NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
    FOREIGN KEY (season_id) REFERENCES seasons (id) ON DELETE CASCADE,
    PRIMARY KEY (user_id, season_id)
);

INSERT INTO teams (id, name, short_name, abbreviation, country, crest_url)
VALUES ('1', 'Bayer 04 Leverkusen', 'Leverkusen', 'B04', 'DE', ''),
       ('2', 'Borussia Dortmund', 'Dortmund', 'BVB', 'DE', ''),
       ('3', 'FC Bayern München', 'Bayern', 'FCB', 'DE', ''),
       ('4', 'VfB Stuttgart', 'Stuttgart', 'VFB', 'DE', ''),
       ('5', 'Arsenal FC', 'Arsenal', 'ARS', 'EN', ''),
       ('6', 'Aston Villa FC', 'Aston Villa', 'AVL', 'EN', ''),
       ('7', 'Liverpool FC', 'Liverpool', 'LIV', 'EN', ''),
       ('8', 'Manchester City FC', 'Man City', 'MCI', 'EN', ''),
       ('9', 'Club Atlético de Madrid', 'Atleti', 'ATL', 'ES', ''),
       ('10', 'FC Barcelona', 'Barça', 'FCB', 'ES', ''),
       ('11', 'Real Madrid CF', 'Real Madrid', 'RMA', 'ES', ''),
       ('12', 'AC Milan', 'Milan', 'MIL', 'IT', ''),
       ('13', 'Atalanta BC', 'Atalanta', 'ATA', 'IT', ''),
       ('14', 'Bologna FC 1909', 'Bologna', 'BOL', 'IT', ''),
       ('15', 'FC Internazionale Milano', 'Inter', 'INT', 'IT', ''),
       ('16', 'Juventus FC', 'Juventus', 'JUV', 'IT', ''),
       ('17', 'Girona FC', 'Girona', 'GIR', 'ES', ''),
       ('18', 'Sporting Clube de Portugal', 'Sporting CP', 'SPO', 'PT', ''),
       ('19', 'Stade Brestois 29', 'Brest', 'BRE', 'FR', ''),
       ('20', 'Lille OSC', 'Lille', 'LIL', 'FR', ''),
       ('21', 'Paris Saint-Germain FC', 'PSG', 'PSG', 'FR', ''),
       ('22', 'AS Monaco FC', 'Monaco', 'ASM', 'MC', ''),
       ('23', 'PSV', 'PSV', 'PSV', 'NL', ''),
       ('24', 'Feyenoord Rotterdam', 'Feyenoord', 'FEY', 'NL', ''),
       ('25', 'RB Leipzig', 'RB Leipzig', 'RBL', 'DE', ''),
       ('26', 'Celtic FC', 'Celtic', 'CEL', 'SC', ''),
       ('27', 'GNK Dinamo Zagreb', 'Dinamo Zagreb', 'DIN', 'HR', ''),
       ('28', 'Club Brugge KV', 'Club Brugge', 'CLU', 'BE', ''),
       ('29', 'AC Sparta Praha', 'Sparta Praha', 'SPP', 'CZ', ''),
       ('30', 'BSC Young Boys', 'Young Boys', 'YOB', 'CH', ''),
       ('31', 'FC Red Bull Salzburg', 'RB Salzburg', 'RBS', 'AT', ''),
       ('32', 'FK Shakhtar Donetsk', 'Shaktar', 'SHD', 'UA', ''),
       ('33', 'Sport Lisboa e Benfica', 'SL Benfica', 'BEN', 'PT', ''),
       ('34', 'SK Sturm Graz', 'Sturm Graz', 'STU', 'AT', ''),
       ('35', 'FK Crvena Zvezda', 'Crvena Zvedza', 'CRV', 'RS', ''),
       ('36', 'ŠK Slovan Bratislava', 'Sl. Bratislava', 'SBA', 'SK', '');

INSERT INTO seasons (id, name, start_date, end_date)
VALUES ('1', '2024/25', '2024-08-01', '2025-05-31');
