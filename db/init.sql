-- Создаем новую таблицу для хранения информации о целевых хостах и их IP-адресах
CREATE TABLE ip_table (
    id SERIAL PRIMARY KEY,
    listContainer VARCHAR(255) NOT NULL,
    ip_address INET NOT NULL,
    time TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);