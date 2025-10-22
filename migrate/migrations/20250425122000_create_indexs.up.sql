SET statement_timeout = 0;

--bun:split

CREATE INDEX IF NOT EXISTS traffic_imei_idx ON traffic (imei);
CREATE INDEX IF NOT EXISTS traffic_imei_alarm_idx ON traffic (alarm);
CREATE INDEX IF NOT EXISTS traffic_pkey ON traffic (imei) WHERE alarm = '0';
