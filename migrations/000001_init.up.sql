CREATE TABLE bookings(
    id SERIAL PRIMARY KEY,
    place_id INTEGER NOT NULL,
    user_name VARCHAR(100) NOT NULL,
    user_phone VARCHAR(100) NOT NULL UNIQUE,
    start_time VARCHAR(100) NOT NULL,
    end_time VARCHAR(100)
);
