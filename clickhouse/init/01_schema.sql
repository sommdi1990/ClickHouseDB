-- Create database if not exists
CREATE
DATABASE IF NOT EXISTS music_school;
USE
music_school;

-- Dimension tables
CREATE TABLE IF NOT EXISTS students
(
    student_id
    UUID
    DEFAULT
    generateUUIDv4
(
),
    first_name String,
    last_name String,
    phone String,
    email String,
    birth_date Date,
    created_at DateTime DEFAULT now
(
)
    )
    ENGINE = MergeTree
    PRIMARY KEY
(
    student_id
)
    ORDER BY
(
    student_id
);

CREATE TABLE IF NOT EXISTS teachers
(
    teacher_id
    UUID
    DEFAULT
    generateUUIDv4
(
),
    first_name String,
    last_name String,
    phone String,
    email String,
    hire_date Date,
    specialty String,
    created_at DateTime DEFAULT now
(
)
    )
    ENGINE = MergeTree
    PRIMARY KEY
(
    teacher_id
)
    ORDER BY
(
    teacher_id
);

CREATE TABLE IF NOT EXISTS instruments
(
    instrument_id
    UUID
    DEFAULT
    generateUUIDv4
(
),
    name String,
    family String,
    created_at DateTime DEFAULT now
(
)
    )
    ENGINE = MergeTree
    PRIMARY KEY
(
    instrument_id
)
    ORDER BY
(
    instrument_id
);

CREATE TABLE IF NOT EXISTS rooms
(
    room_id
    UUID
    DEFAULT
    generateUUIDv4
(
),
    name String,
    capacity UInt16,
    created_at DateTime DEFAULT now
(
)
    )
    ENGINE = MergeTree
    PRIMARY KEY
(
    room_id
)
    ORDER BY
(
    room_id
);

CREATE TABLE IF NOT EXISTS courses
(
    course_id
    UUID
    DEFAULT
    generateUUIDv4
(
),
    title String,
    level String,
    instrument_id UUID,
    description String,
    created_at DateTime DEFAULT now
(
)
    )
    ENGINE = MergeTree
    PRIMARY KEY
(
    course_id
)
    ORDER BY
(
    course_id
);

-- Fact-like event tables
CREATE TABLE IF NOT EXISTS enrollments
(
    enrollment_id
    UUID
    DEFAULT
    generateUUIDv4
(
),
    student_id UUID,
    course_id UUID,
    enrolled_at DateTime DEFAULT now
(
),
    status Enum8
(
    'active' =
    1,
    'completed' =
    2,
    'cancelled' =
    3
)
    )
    ENGINE = MergeTree
    PRIMARY KEY
(
    enrollment_id
)
    ORDER BY
(
    enrollment_id
);

CREATE TABLE IF NOT EXISTS lessons
(
    lesson_id
    UUID
    DEFAULT
    generateUUIDv4
(
),
    course_id UUID,
    teacher_id UUID,
    room_id UUID,
    scheduled_at DateTime,
    duration_min UInt16,
    topic String
    )
    ENGINE = MergeTree
    PRIMARY KEY
(
    lesson_id
)
    ORDER BY
(
    lesson_id
);

CREATE TABLE IF NOT EXISTS attendance
(
    attendance_id
    UUID
    DEFAULT
    generateUUIDv4
(
),
    lesson_id UUID,
    student_id UUID,
    attended UInt8,
    notes String,
    marked_at DateTime DEFAULT now
(
)
    )
    ENGINE = MergeTree
    PRIMARY KEY
(
    attendance_id
)
    ORDER BY
(
    attendance_id
);

CREATE TABLE IF NOT EXISTS payments
(
    payment_id
    UUID
    DEFAULT
    generateUUIDv4
(
),
    student_id UUID,
    amount Decimal
(
    12,
    2
),
    method Enum8
(
    'cash' =
    1,
    'card' =
    2,
    'online' =
    3
),
    paid_at DateTime DEFAULT now
(
),
    description String
    )
    ENGINE = MergeTree
    PRIMARY KEY
(
    payment_id
)
    ORDER BY
(
    payment_id
);

CREATE TABLE IF NOT EXISTS practice_logs
(
    log_id
    UUID
    DEFAULT
    generateUUIDv4
(
),
    student_id UUID,
    instrument_id UUID,
    practiced_at Date,
    minutes UInt16,
    notes String
    )
    ENGINE = MergeTree
    PRIMARY KEY
(
    log_id
)
    ORDER BY
(
    log_id
);

-- Materialized view for daily practice summary
CREATE TABLE IF NOT EXISTS practice_daily_agg
(
    student_id
    UUID,
    practiced_at
    Date,
    total_minutes
    UInt32
)
    ENGINE = SummingMergeTree
    PRIMARY KEY
(
    student_id,
    practiced_at
)
    ORDER BY
(
    student_id,
    practiced_at
);

CREATE
MATERIALIZED VIEW IF NOT EXISTS mv_practice_daily_agg
TO practice_daily_agg AS
SELECT student_id,
       practiced_at,
       sum(minutes) AS total_minutes
FROM practice_logs
GROUP BY student_id,
         practiced_at;


