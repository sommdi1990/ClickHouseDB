USE
music_school;

INSERT INTO instruments (instrument_id, name, family)
VALUES (generateUUIDv4(), 'پیانو', 'کیبورد'),
       (generateUUIDv4(), 'گیتار کلاسیک', 'زهی'),
       (generateUUIDv4(), 'ویولن', 'زهی'),
       (generateUUIDv4(), 'فلوت', 'بادی چوبی');

INSERT INTO teachers (teacher_id, first_name, last_name, phone, email, hire_date, specialty)
VALUES (generateUUIDv4(), 'سارا', 'محمدی', '09120000001', 'sara@example.com', toDate('2023-01-10'), 'پیانو'),
       (generateUUIDv4(), 'علی', 'جعفری', '09120000002', 'ali@example.com', toDate('2022-09-05'), 'گیتار'),
       (generateUUIDv4(), 'نرگس', 'نعمتی', '09120000003', 'narges@example.com', toDate('2021-03-21'), 'ویولن');

-- Create few courses (instrument_id will be linked via subquery to first matches)
INSERT INTO courses (course_id, title, level, instrument_id, description)
SELECT generateUUIDv4(), title, level, instrument_id, description
FROM (SELECT 'پیانو مقدماتی' AS                                          title,
             'مقدماتی'       AS                                          level,
             (SELECT instrument_id FROM instruments WHERE name = 'پیانو' LIMIT 1) AS instrument_id, 'مبانی نوازندگی پیانو' AS description
UNION ALL
SELECT 'گیتار متوسط', 'متوسط', (SELECT instrument_id FROM instruments WHERE name = 'گیتار کلاسیک' LIMIT 1), 'آکوردها و ریتم‌های پیشرفته'
UNION ALL
SELECT 'ویولن پیشرفته',
       'پیشرفته',
       (SELECT instrument_id FROM instruments WHERE name = 'ویولن' LIMIT 1), 'تکنیک‌های پیشرفته'
);

-- Rooms
INSERT INTO rooms (room_id, name, capacity)
VALUES (generateUUIDv4(), 'اتاق 101', 4),
       (generateUUIDv4(), 'اتاق 102', 3);

-- Students
INSERT INTO students (student_id, first_name, last_name, phone, email, birth_date)
VALUES (generateUUIDv4(), 'محمد', 'رضایی', '09123000001', 'm.razai@example.com', toDate('2005-05-12')),
       (generateUUIDv4(), 'الهه', 'هاشمی', '09123000002', 'elahe@example.com', toDate('2003-11-23')),
       (generateUUIDv4(), 'رضا', 'مرادی', '09123000003', 'reza@example.com', toDate('1999-02-02'));


