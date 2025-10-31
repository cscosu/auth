BEGIN;
INSERT OR REPLACE INTO users 
(buck_id, discord_id, display_name, name_num, last_seen_timestamp, last_attended_timestamp, added_to_mailinglist, student, alum, employee, faculty, is_admin)
VALUES
(500123456, NULL, 'Brutus Buckeye', 'buckeye.1', NULL, NULL, 1,1,0,0,0,1);

INSERT OR REPLACE INTO users
(buck_id, discord_id, display_name, name_num, last_seen_timestamp, last_attended_timestamp, added_to_mailinglist, student, alum, employee, faculty)
VALUES
(0, 1014305712312172595, 'Mark Bundschuh', 'bundschuh.15', 0, 0, 0, 1, 0, 0, 0),
(0, 393776786812436491, 'Benjamin Bundschuh', 'bundschuh.13', 0, 0, 0, 1, 0, 0, 0);
COMMIT;