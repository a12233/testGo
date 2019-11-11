--
-- PostgreSQL database dump
--

-- Dumped from database version 11.5
-- Dumped by pg_dump version 11.5

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
-- Name: ci_job_status; Type: TYPE; Schema: public; Owner: postgres
--

CREATE TYPE public.ci_job_status AS ENUM (
    'new',
    'initializing',
    'initialized',
    'running',
    'success',
    'error'
);


ALTER TYPE public.ci_job_status OWNER TO postgres;

--
-- Name: ci_jobs_status_notify(); Type: FUNCTION; Schema: public; Owner: postgres
--

CREATE FUNCTION public.ci_jobs_status_notify() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
DECLARE
    payload varchar;
    mid uuid; 
BEGIN
 select * into mid from ci_jobs where status='initializing';
    payload = CAST(NEW.id AS text) ||
    ',' || CAST(NEW.jobname AS text) ||  ',' || CAST(NEW.status AS text) ||
     ',' || CAST(NEW.status_change_time AS text);
PERFORM pg_notify('ci_jobs_status_channel', payload);
RETURN NEW;
END;
$$;


ALTER FUNCTION public.ci_jobs_status_notify() OWNER TO postgres;

SET default_tablespace = '';

SET default_with_oids = false;

--
-- Name: ci_jobs; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.ci_jobs (
    id integer NOT NULL,
    jobname character varying(256),
    status public.ci_job_status,
    status_change_time timestamp without time zone
);


ALTER TABLE public.ci_jobs OWNER TO postgres;

--
-- Name: ci_jobs_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.ci_jobs_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.ci_jobs_id_seq OWNER TO postgres;

--
-- Name: ci_jobs_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.ci_jobs_id_seq OWNED BY public.ci_jobs.id;


--
-- Name: mid; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.mid (
    id integer,
    jobname character varying(256),
    status public.ci_job_status,
    status_change_time timestamp without time zone
);


ALTER TABLE public.mid OWNER TO postgres;

--
-- Name: ci_jobs id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.ci_jobs ALTER COLUMN id SET DEFAULT nextval('public.ci_jobs_id_seq'::regclass);


--
-- Data for Name: ci_jobs; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.ci_jobs (id, jobname, status, status_change_time) FROM stdin;
1	job1	new	2019-11-08 14:23:59.423696
\.


--
-- Data for Name: mid; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.mid (id, jobname, status, status_change_time) FROM stdin;
1	job1	initializing	2019-11-08 14:23:59.423696
\.


--
-- Name: ci_jobs_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.ci_jobs_id_seq', 1, true);


--
-- Name: ci_jobs ci_jobs_status; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER ci_jobs_status AFTER INSERT OR UPDATE OF status ON public.ci_jobs FOR EACH ROW EXECUTE PROCEDURE public.ci_jobs_status_notify();


--
-- PostgreSQL database dump complete
--

