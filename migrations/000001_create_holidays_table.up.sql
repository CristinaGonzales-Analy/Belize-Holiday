CREATE TABLE IF NOT EXISTS holidays (
    id       BIGSERIAL PRIMARY KEY,
    day      TEXT NOT NULL,
    date     TEXT NOT NULL,
    month    INTEGER NOT NULL,
    year     INTEGER NOT NULL,
    occasion TEXT NOT NULL
);

INSERT INTO holidays (day, date, month, year, occasion) VALUES
    ('Thursday',  '1st January',    1,  2026, 'New Year''s Day'),
    ('Thursday',  '15th January',   1,  2026, 'George Price Day'),
    ('Monday',    '9th March',      3,  2026, 'National Heroes and Benefactor Day'),
    ('Friday',    '3rd April',      4,  2026, 'Good Friday'),
    ('Saturday',  '4th April',      4,  2026, 'Holy Saturday'),
    ('Monday',    '6th April',      4,  2026, 'Easter Monday'),
    ('Friday',    '1st May',        5,  2026, 'Labour Day'),
    ('Saturday',  '1st August',     8,  2026, 'Emancipation Day'),
    ('Thursday',  '10th September', 9,  2026, 'St. George''s Caye Day'),
    ('Monday',    '21st September', 9,  2026, 'Independence Day'),
    ('Monday',    '12th October',   10, 2026, 'Indigenous People''s Resistance Day'),
    ('Thursday',  '19th November',  11, 2026, 'Garifuna Settlement Day'),
    ('Friday',    '25th December',  12, 2026, 'Christmas Day'),
    ('Saturday',  '26th December',  12, 2026, 'Boxing Day');