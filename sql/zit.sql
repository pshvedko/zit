--
-- PostgreSQL database dump
--

-- Dumped from database version 13.9
-- Dumped by pg_dump version 13.7

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

--
-- Name: zit; Type: DATABASE; Schema: -; Owner: postgres
--

CREATE DATABASE zit WITH TEMPLATE = template0 ENCODING = 'UTF8' LC_COLLATE = 'C' LC_CTYPE = 'ru_RU.UTF-8';


ALTER DATABASE zit OWNER TO postgres;

\connect zit

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

--
-- Name: log_notify(); Type: FUNCTION; Schema: public; Owner: postgres
--

CREATE FUNCTION public.log_notify() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
	PERFORM pg_notify('log', NEW.id || ' ' || NEW.ip);
	RETURN NULL;
END;
$$;


ALTER FUNCTION public.log_notify() OWNER TO postgres;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: log; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.log (
    id bigint,
    ip inet check ( family(ip) = 4 ),
    date timestamp without time zone DEFAULT now()
);


ALTER TABLE public.log OWNER TO postgres;

--
-- Data for Name: log; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.log (id, ip, date) FROM stdin;
1	127.0.0.1	2023-11-01 18:59:47.421634
2	127.0.0.1	2023-11-01 19:00:53.485579
1	127.0.0.1	2023-11-01 19:05:53.053217
1	127.0.0.2	2023-11-01 19:06:07.708868
2	127.0.0.2	2023-11-01 19:06:18.58874
2	127.0.0.3	2023-11-01 19:06:33.740626
3	127.0.0.3	2023-11-01 19:06:50.42884
3	127.0.0.1	2023-11-01 19:06:55.372614
4	127.0.0.1	2023-11-01 19:07:05.404575
\.


--
-- Name: log log_notify; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER log_notify AFTER INSERT ON public.log FOR EACH ROW EXECUTE FUNCTION public.log_notify();


--
-- PostgreSQL database dump complete
--
