SET statement_timeout = 0;

--bun:split

ALTER TABLE public.traffic DROP COLUMN isnotified;
DROP INDEX IF EXISTS traffic_counter;
