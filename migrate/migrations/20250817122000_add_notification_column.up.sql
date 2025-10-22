SET statement_timeout = 0;

--bun:split

ALTER TABLE public.traffic ADD isnotified boolean DEFAULT false NULL;
CREATE INDEX ON public.traffic (isnotified);

