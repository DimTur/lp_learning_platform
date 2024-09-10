CREATE TABLE IF NOT EXISTS lp_channels
(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    created_by INTEGER NOT NULL,
    last_modified_by INTEGER NOT NULL,
    public INTEGER NOT NULL CHECK(public IN (0, 1))
);

CREATE TABLE IF NOT EXISTS lp_plans
(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    created_by INTEGER NOT NULL,
    last_modified_by INTEGER NOT NULL,
    is_published INTEGER NOT NULL CHECK(is_published IN (0, 1)),
    published_at TIMESTAMP NOT NULL,
    public INTEGER NOT NULL CHECK(public IN (0, 1))
);

CREATE TABLE IF NOT EXISTS lp_lessons
(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    pass_percentage INTEGER NOT NULL,
    created_by INTEGER NOT NULL,
    last_modified_by INTEGER NOT NULL
);

CREATE TABLE IF NOT EXISTS lp_channel_plans
(
    channel_id INTEGER NOT NULL,
    plan_id INTEGER NOT NULL,
    PRIMARY KEY (channel_id, plan_id),
    CONSTRAINT fk_channel
        FOREIGN KEY (channel_id)
        REFERENCES lp_channels(id)
        ON DELETE CASCADE,
    CONSTRAINT fk_plan
        FOREIGN KEY (plan_id)
        REFERENCES lp_plans(id)
        ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS lp_plan_lessons
(
    plan_id INTEGER NOT NULL,
    lesson_id INTEGER NOT NULL,
    PRIMARY KEY (plan_id, lesson_id),
    CONSTRAINT fk_plan
        FOREIGN KEY (plan_id)
        REFERENCES lp_plans(id)
        ON DELETE CASCADE,
    CONSTRAINT fk_lesson
        FOREIGN KEY (lesson_id)
        REFERENCES lp_lessons(id)
        ON DELETE CASCADE
);