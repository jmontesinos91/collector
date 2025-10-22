SET statement_timeout = 0;

--bun:split

ALTER TABLE public.traffic ADD counter int8 DEFAULT 0 NULL;
CREATE INDEX ON public.traffic (counter);

