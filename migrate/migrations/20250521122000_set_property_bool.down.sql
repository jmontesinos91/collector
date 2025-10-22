SET statement_timeout = 0;

--bun:split

ALTER TABLE public.traffic ALTER COLUMN alarm TYPE bool USING alarm::varchar(256);
