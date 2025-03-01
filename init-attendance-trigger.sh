#!/bin/sh
set -e

echo "Ожидание PostgreSQL..."
# shellcheck disable=SC3020
until PGPASSWORD=password psql -U user -d biotime -c "SELECT 1" &>/dev/null; do
  sleep 1
done

PGPASSWORD=password psql -U user -d biotime <<EOSQL
    CREATE TABLE IF NOT EXISTS attendance_log (
        id SERIAL PRIMARY KEY,
        emp_id INT NOT NULL,
        punch_time TIMESTAMP DEFAULT now(),
        terminal_alias TEXT NOT NULL,
        processed BOOLEAN DEFAULT FALSE
    );

    CREATE OR REPLACE FUNCTION notify_attendance_event() RETURNS TRIGGER AS \$\$
    BEGIN
        PERFORM pg_notify('attendance_events', row_to_json(NEW)::text);
        RETURN NEW;
    END;
    \$\$ LANGUAGE plpgsql;

    DROP TRIGGER IF EXISTS attendance_trigger ON attendance_log;
    CREATE TRIGGER attendance_trigger
    AFTER INSERT ON attendance_log
    FOR EACH ROW EXECUTE FUNCTION notify_attendance_event();
EOSQL

echo "База данных готова!"