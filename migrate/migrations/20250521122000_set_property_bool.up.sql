SET statement_timeout = 0;

--bun:split

CREATE INDEX ON public.traffic(created_at);
ALTER TABLE public.traffic RENAME COLUMN alarm TO "isAlarm";
ALTER TABLE public.traffic ALTER COLUMN "isAlarm" TYPE bool USING "isAlarm"::bool;
CREATE INDEX ON public.traffic (imei) WHERE "isAlarm" = false;
CREATE INDEX ON public.traffic (imei, "isAlarm");

