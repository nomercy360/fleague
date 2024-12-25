DELETE
FROM users WHERE 1=1;
DELETE
FROM matches WHERE 1=1;
DELETE
FROM predictions WHERE 1=1;
DELETE
FROM user_followers WHERE 1=1;
DELETE
FROM leaderboards WHERE 1=1;

CREATE TRIGGER update_user_total_predictions
    AFTER INSERT
    ON predictions
BEGIN
    UPDATE users
    SET total_predictions = total_predictions + 1
    WHERE id = NEW.user_id;
END;

INSERT INTO users (id, first_name, last_name, username, language_code, chat_id, referred_by,
                   total_points, total_predictions, correct_predictions, avatar_url)
VALUES ('1', 'Maksim', NULL, 'mkkksim', 'en', 927635965, 'ref1', 0, 0, 0,
        'https://assets.peatch.io/fb/users/uY8YwvCn.jpg'),
       ('2', 'Gor', NULL, 'cronaldo', 'en', 428630919, 'ref2', 0, 0, 0,
        'https://assets.peatch.io/avatars/gor.jpeg'),
       ('3', 'Neymar', 'Junior', 'njunior', 'en', 345678912, 1, 0, 0, 0,
        'https://assets.peatch.io/avatars/neymar.webp'),
       ('4', 'Kevin', 'De Bruyne', 'kdebruyne', 'en', 456789123, 2, 0, 0, 0,
        'https://assets.peatch.io/avatars/kevin.webp'),
       ('5', 'Kylian', 'Mbappe', 'kmbappe', 'en', 567891234, 3, 0, 0, 0,
        'https://assets.peatch.io/avatars/mbappe.jpg'),
       ('6', 'Erling', 'Haaland', 'ehaaland', 'en', 678912345, NULL, 0, 0, 0,
        'https://assets.peatch.io/avatars/haaland.jpeg'),
       ('7', 'Robert', 'Lewandowski', 'rlewandowski', 'en', 789123456, 5, 0, 0, 0,
        'https://assets.peatch.io/avatars/lewandowski.jpg'),
       ('8', 'Sadio', 'Mane', 'smane', 'en', 891234567, 4, 0, 0, 0,
        'https://assets.peatch.io/avatars/sadio-mane.gif'),
       ('9', 'Virgil', 'van Dijk', 'vvdijk', 'en', 912345678, 6, 0, 0, 0,
        'https://assets.peatch.io/avatars/van_dijk.jpg'),
       ('10', 'Mohamed', 'Salah', 'msalah', 'en', 123456789, 8, 0, 0, 0,
        'https://assets.peatch.io/avatars/salah.jpg');

-- schedule matches
INSERT INTO matches (id, tournament, home_team_id, away_team_id, match_date, status, home_score, away_score)
VALUES ('10', 'Champions League', '3', '12', datetime('now', '+10 days'), 'scheduled', NULL, NULL),
       ('11', 'Premier League', '7', '15', datetime('now', '+2 days'), 'scheduled', NULL, NULL),
       ('12', 'Champions League', '9', '20', datetime('now', '+4 days'), 'scheduled', NULL, NULL),
       ('13', 'La Liga', '5', '8', datetime('now', '+3 days'), 'scheduled', NULL, NULL),
       ('14', 'La Liga', '10', '14', datetime('now', '+5 days'), 'scheduled', NULL, NULL),
       ('15', 'La Liga', '18', '21', datetime('now', '+5 days'), 'scheduled', NULL, NULL),
       ('16', 'Champions League', '6', '17', datetime('now', '+5 days'), 'scheduled', NULL, NULL),
       ('17', 'La Liga', '12', '15', datetime('now', '+20 days'), 'scheduled', NULL, NULL),
       ('18', 'Premier League', '3', '7', datetime('now', '+15 days'), 'scheduled', NULL, NULL),
       ('19', 'Champions League', '2', '10', datetime('now', '+6 days'), 'scheduled', NULL, NULL),
       ('20', 'Premier League', '4', '13', datetime('now', '+7 days'), 'scheduled', NULL, NULL);

-- completed matches
INSERT INTO matches (id, tournament, home_team_id, away_team_id, match_date, status, home_score, away_score)
VALUES ('1', 'Champions League', '4', '13', datetime('now', '-20 days'), 'completed', 0, 0),
       ('2', 'Premier League', '11', '16', datetime('now', '+15 days'), 'completed', 2, 3),
       ('3', 'Premier League', '2', '6', datetime('now', '-10 days'), 'completed', 1, 2),
       ('4', 'Champions League', '8', '19', datetime('now', '-15 days'), 'completed', 2, 2),
       ('5', 'La Liga', '5', '11', datetime('now', '-25 days'), 'completed', 0, 1),
       ('6', 'La Liga', '6', '20', datetime('now', '-30 days'), 'completed', 3, 1),
       ('7', 'Champions League', '9', '16', datetime('now', '-35 days'), 'completed', 3, 3),
       ('8', 'Premier League', '1', '11', datetime('now', '-40 days'), 'completed', 2, 0),
       ('9', 'Premier League', '9', '14', datetime('now', '-5 days'), 'completed', 3, 1);


-- INSERT INTO predictions (user_id, match_id, predicted_outcome, predicted_home_score, predicted_away_score,
--                          points_awarded, created_at, completed_at)

-- completed predictions, wrong. Format: (note date before match date and match should be completed). Can be either score or outcome
--- (1, 3, 'home', NULL, NULL, 0, datetime('now', '-25 days'), NULL)
--- or
--- (1, 3, NULL, 0, 5, 0, datetime('now', '-25 days'), NULL)

INSERT INTO predictions (user_id, match_id, predicted_outcome, predicted_home_score, predicted_away_score,
                         points_awarded, created_at, completed_at)
VALUES ('1', '3', 'home', NULL, NULL, 0, datetime('now', '-25 days'), NULL),
       ('1', '4', NULL, 1, 2, 0, datetime('now', '-15 days'), NULL),
       ('1', '5', 'home', NULL, NULL, 0, datetime('now', '-30 days'), NULL),
       ('2', '6', 'away', NULL, NULL, 0, datetime('now', '-35 days'), NULL),
       ('2', '7', NULL, 2, 2, 0, datetime('now', '-40 days'), NULL),
       ('2', '8', 'draw', NULL, NULL, 0, datetime('now', '-5 days'), NULL),
       ('3', '9', 'away', NULL, NULL, 0, datetime('now', '-20 days'), NULL),
       ('3', '6', NULL, 1, 3, 0, datetime('now', '-10 days'), NULL),
       ('3', '4', 'home', NULL, NULL, 0, datetime('now', '-25 days'), NULL),
       ('4', '5', NULL, 2, 1, 0, datetime('now', '-15 days'), NULL),
       ('4', '6', 'away', NULL, NULL, 0, datetime('now', '-25 days'), NULL),
       ('4', '7', NULL, 0, 2, 0, datetime('now', '-30 days'), NULL),
       ('5', '8', 'draw', NULL, NULL, 0, datetime('now', '-35 days'), NULL),
       ('5', '9', NULL, 3, 3, 0, datetime('now', '-40 days'), NULL),
       ('5', '1', 'away', NULL, NULL, 0, datetime('now', '-5 days'), NULL),
       ('6', '2', 'home', NULL, NULL, 0, datetime('now', '-20 days'), NULL),
       ('6', '3', NULL, 2, 4, 0, datetime('now', '-10 days'), NULL),
       ('6', '4', 'away', NULL, NULL, 0, datetime('now', '-25 days'), NULL),
       ('7', '1', NULL, 3, 5, 0, datetime('now', '-15 days'), NULL),
       ('7', '2', 'home', NULL, NULL, 0, datetime('now', '-25 days'), NULL),
       ('7', '3', NULL, 2, 2, 0, datetime('now', '-30 days'), NULL),
       ('8', '4', 'home', NULL, NULL, 0, datetime('now', '-35 days'), NULL),
       ('8', '5', NULL, 3, 1, 0, datetime('now', '-40 days'), NULL),
       ('8', '6', 'draw', NULL, NULL, 0, datetime('now', '-5 days'), NULL),
       ('9', '7', 'home', NULL, NULL, 0, datetime('now', '-20 days'), NULL),
       ('9', '8', NULL, 1, 4, 0, datetime('now', '-10 days'), NULL),
       ('9', '9', 'draw', NULL, NULL, 0, datetime('now', '-25 days'), NULL),
       ('10', '8', NULL, 3, 1, 0, datetime('now', '-15 days'), NULL),
       ('10', '7', 'home', NULL, NULL, 0, datetime('now', '-25 days'), NULL),
       ('10', '6', NULL, 0, 1, 0, datetime('now', '-30 days'), NULL);

-- completed predictions, correct, by outcome. Format: (note date before match date and match should be completed)
-- Each prediction gives 3 points
--- (3, 6, 'home', NULL, NULL, 0, datetime('now', '-35 days'), NULL)

INSERT INTO predictions (user_id, match_id, predicted_outcome, predicted_home_score, predicted_away_score,
                         points_awarded, created_at, completed_at)
VALUES ('2', '1', 'draw', NULL, NULL, 0, datetime('now', '-25 days'), NULL),
       ('1', '2', 'away', NULL, NULL, 0, datetime('now', '-25 days'), NULL),
       ('3', '2', 'away', NULL, NULL, 0, datetime('now', '-25 days'), NULL),
       ('4', '2', 'away', NULL, NULL, 0, datetime('now', '-25 days'), NULL),
       ('4', '3', 'away', NULL, NULL, 0, datetime('now', '-25 days'), NULL),
       ('5', '3', 'away', NULL, NULL, 0, datetime('now', '-25 days'), NULL),
       ('9', '3', 'away', NULL, NULL, 0, datetime('now', '-25 days'), NULL),
       ('8', '3', 'away', NULL, NULL, 0, datetime('now', '-25 days'), NULL),
       ('7', '4', 'draw', NULL, NULL, 0, datetime('now', '-25 days'), NULL),
       ('9', '4', 'draw', NULL, NULL, 0, datetime('now', '-25 days'), NULL),
       ('6', '5', 'away', NULL, NULL, 0, datetime('now', '-25 days'), NULL),
       ('9', '5', 'away', NULL, NULL, 0, datetime('now', '-25 days'), NULL),
       ('3', '5', 'away', NULL, NULL, 0, datetime('now', '-25 days'), NULL),
       ('2', '5', 'away', NULL, NULL, 0, datetime('now', '-25 days'), NULL),
       ('5', '6', 'home', NULL, NULL, 0, datetime('now', '-25 days'), NULL),
       ('9', '6', 'home', NULL, NULL, 0, datetime('now', '-25 days'), NULL),
       ('1', '6', 'home', NULL, NULL, 0, datetime('now', '-25 days'), NULL),
       ('6', '7', 'draw', NULL, NULL, 0, datetime('now', '-25 days'), NULL),
       ('5', '7', 'draw', NULL, NULL, 0, datetime('now', '-25 days'), NULL),
       ('3', '7', 'draw', NULL, NULL, 0, datetime('now', '-25 days'), NULL),
       ('4', '8', 'home', NULL, NULL, 0, datetime('now', '-25 days'), NULL),
       ('6', '8', 'home', NULL, NULL, 0, datetime('now', '-25 days'), NULL),
       ('6', '9', 'home', NULL, NULL, 0, datetime('now', '-25 days'), NULL),
       ('7', '9', 'home', NULL, NULL, 0, datetime('now', '-25 days'), NULL),
       ('1', '9', 'home', NULL, NULL, 0, datetime('now', '-25 days'), NULL);

-- 1: 9, 2: 6, 3: 9, 4: 9, 5: 9, 6: 12, 7: 6, 8: 3, 9: 12, 10: 0

-- completed predictions, correct, by score. Format: (note date before match date and match should be completed)
-- Each prediction gives 5 points
--- (5, 6, NULL, 3, 1, 0, datetime('now', '-60 days'), NULL)
INSERT INTO predictions (user_id, match_id, predicted_outcome, predicted_home_score, predicted_away_score,
                         points_awarded, created_at, completed_at)
VALUES ('10', '1', NULL, 0, 0, 0, datetime('now', '-25 days'), NULL),
       ('1', '1', NULL, 0, 0, 0, datetime('now', '-25 days'), NULL),
       ('10', '2', NULL, 2, 3, 0, datetime('now', '-25 days'), NULL),
       ('2', '2', NULL, 2, 3, 0, datetime('now', '-25 days'), NULL),
       ('2', '3', NULL, 1, 2, 0, datetime('now', '-25 days'), NULL),
       ('10', '3', NULL, 1, 2, 0, datetime('now', '-25 days'), NULL),
       ('3', '3', NULL, 1, 2, 0, datetime('now', '-25 days'), NULL),
       ('4', '4', NULL, 2, 2, 0, datetime('now', '-25 days'), NULL),
       ('10', '4', NULL, 2, 2, 0, datetime('now', '-25 days'), NULL),
       ('5', '4', NULL, 2, 2, 0, datetime('now', '-25 days'), NULL),
       ('10', '5', NULL, 0, 1, 0, datetime('now', '-25 days'), NULL),
       ('5', '5', NULL, 0, 1, 0, datetime('now', '-25 days'), NULL),
       ('6', '6', NULL, 3, 1, 0, datetime('now', '-25 days'), NULL),
       ('7', '6', NULL, 3, 1, 0, datetime('now', '-25 days'), NULL),
       ('7', '7', NULL, 3, 3, 0, datetime('now', '-25 days'), NULL),
       ('1', '7', NULL, 3, 3, 0, datetime('now', '-25 days'), NULL),
       ('7', '8', NULL, 2, 0, 0, datetime('now', '-25 days'), NULL),
       ('8', '8', NULL, 2, 0, 0, datetime('now', '-25 days'), NULL),
       ('1', '8', NULL, 2, 0, 0, datetime('now', '-25 days'), NULL),
       ('10', '9', NULL, 3, 1, 0, datetime('now', '-25 days'), NULL),
       ('4', '9', NULL, 3, 1, 0, datetime('now', '-25 days'), NULL),
       ('8', '9', NULL, 3, 1, 0, datetime('now', '-25 days'), NULL);

-- 1: 15, 2: 10, 3:5, 4: 10, 5: 10, 6: 5, 7: 15, 8:10, 9: 0, 10: 30

-- ongoing predictions. Format: (note date before match date and match should be scheduled). Can be either score or outcome
INSERT INTO predictions (user_id, match_id, predicted_outcome, predicted_home_score, predicted_away_score,
                         points_awarded, created_at, completed_at)
VALUES ('1', '11', 'draw', NULL, NULL, 0, datetime('now', '-25 days'), NULL),
       ('2', '11', NULL, 0, 0, 0, datetime('now', '-25 days'), NULL),
       ('3', '12', 'away', NULL, NULL, 0, datetime('now', '-25 days'), NULL),
       ('4', '13', NULL, 0, 0, 0, datetime('now', '-25 days'), NULL),
       ('5', '13', 'home', NULL, NULL, 0, datetime('now', '-25 days'), NULL),
       ('7', '14', NULL, 0, 0, 0, datetime('now', '-25 days'), NULL),
       ('9', '14', NULL, 0, 0, 0, datetime('now', '-25 days'), NULL),
       ('6', '15', NULL, 0, 0, 0, datetime('now', '-25 days'), NULL),
       ('7', '15', 'away', NULL, NULL, 0, datetime('now', '-25 days'), NULL),
       ('8', '15', NULL, 0, 0, 0, datetime('now', '-25 days'), NULL),
       ('10', '16', NULL, 0, 0, 0, datetime('now', '-25 days'), NULL),
       ('8', '17', 'draw', NULL, NULL, 0, datetime('now', '-25 days'), NULL),
       ('6', '17', NULL, 0, 0, 0, datetime('now', '-25 days'), NULL),
       ('5', '17', NULL, 0, 0, 0, datetime('now', '-25 days'), NULL),
       ('3', '17', 'draw', NULL, NULL, 0, datetime('now', '-25 days'), NULL),
       ('4', '18', NULL, 0, 0, 0, datetime('now', '-25 days'), NULL),
       ('2', '18', 'home', NULL, NULL, 0, datetime('now', '-25 days'), NULL),
       ('1', '19', NULL, 0, 0, 0, datetime('now', '-25 days'), NULL);

-- at the end print report about calculations how many user have: total_points, total_predictions, correct_predictions
-- or give SQL query to get this information

-- Query for report
SELECT id                             AS user_id,
       first_name || ' ' || last_name AS name,
       total_points,
       total_predictions,
       correct_predictions
FROM users
ORDER BY total_points DESC, correct_predictions DESC;

INSERT INTO user_followers (follower_id, following_id)
VALUES ('1', '2'),
       ('1', '3'),
       ('1', '4'),
       ('1', '5'),
       ('1', '6'),
       ('1', '7'),
       ('1', '8'),
       ('1', '9'),
       ('1', '10');


DROP TRIGGER update_user_total_predictions;
