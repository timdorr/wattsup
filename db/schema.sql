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
-- Name: metrics; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.metrics (
    "time" timestamp with time zone NOT NULL,
    device_id integer,
    register_address integer,
    value integer
);


--
-- Name: _direct_view_9; Type: VIEW; Schema: _timescaledb_internal; Owner: -
--

CREATE VIEW _timescaledb_internal._direct_view_9 AS
 SELECT register_address,
    public.time_bucket('01:00:00'::interval, "time") AS bucket,
    avg(value) AS value
   FROM public.metrics
  WHERE ("time" > date_trunc('day'::text, now()))
  GROUP BY register_address, (public.time_bucket('01:00:00'::interval, "time"));


--
-- Name: _hyper_5_2_chunk; Type: TABLE; Schema: _timescaledb_internal; Owner: -
--

CREATE TABLE _timescaledb_internal._hyper_5_2_chunk (
    CONSTRAINT constraint_2 CHECK ((("time" >= '2025-06-11 20:00:00-04'::timestamp with time zone) AND ("time" < '2025-06-18 20:00:00-04'::timestamp with time zone)))
)
INHERITS (public.metrics);


--
-- Name: _hyper_5_3_chunk; Type: TABLE; Schema: _timescaledb_internal; Owner: -
--

CREATE TABLE _timescaledb_internal._hyper_5_3_chunk (
    CONSTRAINT constraint_3 CHECK ((("time" >= '2025-06-18 20:00:00-04'::timestamp with time zone) AND ("time" < '2025-06-25 20:00:00-04'::timestamp with time zone)))
)
INHERITS (public.metrics);


--
-- Name: _materialized_hypertable_9; Type: TABLE; Schema: _timescaledb_internal; Owner: -
--

CREATE TABLE _timescaledb_internal._materialized_hypertable_9 (
    register_address integer,
    bucket timestamp with time zone NOT NULL,
    value numeric
);


--
-- Name: _hyper_9_6_chunk; Type: TABLE; Schema: _timescaledb_internal; Owner: -
--

CREATE TABLE _timescaledb_internal._hyper_9_6_chunk (
    CONSTRAINT constraint_6 CHECK (((bucket >= '2025-06-20 20:00:00-04'::timestamp with time zone) AND (bucket < '2025-06-30 20:00:00-04'::timestamp with time zone)))
)
INHERITS (_timescaledb_internal._materialized_hypertable_9);


--
-- Name: _partial_view_9; Type: VIEW; Schema: _timescaledb_internal; Owner: -
--

CREATE VIEW _timescaledb_internal._partial_view_9 AS
 SELECT register_address,
    public.time_bucket('01:00:00'::interval, "time") AS bucket,
    avg(value) AS value
   FROM public.metrics
  WHERE ("time" > date_trunc('day'::text, now()))
  GROUP BY register_address, (public.time_bucket('01:00:00'::interval, "time"));


--
-- Name: metrics_daily; Type: VIEW; Schema: public; Owner: -
--

CREATE VIEW public.metrics_daily AS
 SELECT register_address,
    bucket,
    value
   FROM _timescaledb_internal._materialized_hypertable_9;


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
-- Name: _hyper_5_2_chunk_metrics_register_address_time_idx; Type: INDEX; Schema: _timescaledb_internal; Owner: -
--

CREATE INDEX _hyper_5_2_chunk_metrics_register_address_time_idx ON _timescaledb_internal._hyper_5_2_chunk USING btree (register_address, "time" DESC);


--
-- Name: _hyper_5_2_chunk_metrics_time_idx; Type: INDEX; Schema: _timescaledb_internal; Owner: -
--

CREATE INDEX _hyper_5_2_chunk_metrics_time_idx ON _timescaledb_internal._hyper_5_2_chunk USING btree ("time" DESC);


--
-- Name: _hyper_5_3_chunk_metrics_register_address_time_idx; Type: INDEX; Schema: _timescaledb_internal; Owner: -
--

CREATE INDEX _hyper_5_3_chunk_metrics_register_address_time_idx ON _timescaledb_internal._hyper_5_3_chunk USING btree (register_address, "time" DESC);


--
-- Name: _hyper_5_3_chunk_metrics_time_idx; Type: INDEX; Schema: _timescaledb_internal; Owner: -
--

CREATE INDEX _hyper_5_3_chunk_metrics_time_idx ON _timescaledb_internal._hyper_5_3_chunk USING btree ("time" DESC);


--
-- Name: _hyper_9_6_chunk__materialized_hypertable_9_bucket_idx; Type: INDEX; Schema: _timescaledb_internal; Owner: -
--

CREATE INDEX _hyper_9_6_chunk__materialized_hypertable_9_bucket_idx ON _timescaledb_internal._hyper_9_6_chunk USING btree (bucket DESC);


--
-- Name: _hyper_9_6_chunk__materialized_hypertable_9_register_address_bu; Type: INDEX; Schema: _timescaledb_internal; Owner: -
--

CREATE INDEX _hyper_9_6_chunk__materialized_hypertable_9_register_address_bu ON _timescaledb_internal._hyper_9_6_chunk USING btree (register_address, bucket DESC);


--
-- Name: _materialized_hypertable_9_bucket_idx; Type: INDEX; Schema: _timescaledb_internal; Owner: -
--

CREATE INDEX _materialized_hypertable_9_bucket_idx ON _timescaledb_internal._materialized_hypertable_9 USING btree (bucket DESC);


--
-- Name: _materialized_hypertable_9_register_address_bucket_idx; Type: INDEX; Schema: _timescaledb_internal; Owner: -
--

CREATE INDEX _materialized_hypertable_9_register_address_bucket_idx ON _timescaledb_internal._materialized_hypertable_9 USING btree (register_address, bucket DESC);


--
-- Name: metrics_register_address_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX metrics_register_address_time_idx ON public.metrics USING btree (register_address, "time" DESC);


--
-- Name: metrics_time_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX metrics_time_idx ON public.metrics USING btree ("time" DESC);


--
-- Name: _hyper_5_2_chunk ts_cagg_invalidation_trigger; Type: TRIGGER; Schema: _timescaledb_internal; Owner: -
--

CREATE TRIGGER ts_cagg_invalidation_trigger AFTER INSERT OR DELETE OR UPDATE ON _timescaledb_internal._hyper_5_2_chunk FOR EACH ROW EXECUTE FUNCTION _timescaledb_functions.continuous_agg_invalidation_trigger('5');


--
-- Name: _hyper_5_3_chunk ts_cagg_invalidation_trigger; Type: TRIGGER; Schema: _timescaledb_internal; Owner: -
--

CREATE TRIGGER ts_cagg_invalidation_trigger AFTER INSERT OR DELETE OR UPDATE ON _timescaledb_internal._hyper_5_3_chunk FOR EACH ROW EXECUTE FUNCTION _timescaledb_functions.continuous_agg_invalidation_trigger('5');


--
-- Name: _materialized_hypertable_9 ts_insert_blocker; Type: TRIGGER; Schema: _timescaledb_internal; Owner: -
--

CREATE TRIGGER ts_insert_blocker BEFORE INSERT ON _timescaledb_internal._materialized_hypertable_9 FOR EACH ROW EXECUTE FUNCTION _timescaledb_functions.insert_blocker();


--
-- Name: metrics ts_cagg_invalidation_trigger; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER ts_cagg_invalidation_trigger AFTER INSERT OR DELETE OR UPDATE ON public.metrics FOR EACH ROW EXECUTE FUNCTION _timescaledb_functions.continuous_agg_invalidation_trigger('5');


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
    ('20250618141940'),
    ('20250622203229'),
    ('20250622205135'),
    ('20250622210447');
