SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET transaction_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

--
-- Name: timescaledb; Type: EXTENSION; Schema: -; Owner: -
--

CREATE EXTENSION IF NOT EXISTS timescaledb WITH SCHEMA public;


--
-- Name: EXTENSION timescaledb; Type: COMMENT; Schema: -; Owner: -
--

COMMENT ON EXTENSION timescaledb IS 'Enables scalable inserts and complex queries for time-series data (Community Edition)';


SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: _compressed_hypertable_4; Type: TABLE; Schema: _timescaledb_internal; Owner: -
--

CREATE TABLE _timescaledb_internal._compressed_hypertable_4 (
);


--
-- Name: metrics; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.metrics (
    "time" timestamp with time zone NOT NULL,
    device_name text,
    device_id integer,
    register_name text,
    register_address integer,
    value integer
);


--
-- Name: schema_migrations; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.schema_migrations (
    version character varying(255) NOT NULL
);


--
-- Name: schema_migrations schema_migrations_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.schema_migrations
    ADD CONSTRAINT schema_migrations_pkey PRIMARY KEY (version);


--
-- Name: metrics_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX metrics_time_idx ON public.metrics USING btree ("time" DESC);


--
-- Name: _compressed_hypertable_4 ts_insert_blocker; Type: TRIGGER; Schema: _timescaledb_internal; Owner: -
--

CREATE TRIGGER ts_insert_blocker BEFORE INSERT ON _timescaledb_internal._compressed_hypertable_4 FOR EACH ROW EXECUTE FUNCTION _timescaledb_functions.insert_blocker();


--
-- Name: metrics ts_insert_blocker; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER ts_insert_blocker BEFORE INSERT ON public.metrics FOR EACH ROW EXECUTE FUNCTION _timescaledb_functions.insert_blocker();


--
-- PostgreSQL database dump complete
--


--
-- Dbmate schema migrations
--

INSERT INTO public.schema_migrations (version) VALUES
    ('20250618141808'),
    ('20250618141940');
