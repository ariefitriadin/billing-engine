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

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: billing_schedule; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.billing_schedule (
    id integer NOT NULL,
    loan_id integer NOT NULL,
    week integer NOT NULL,
    amount numeric NOT NULL,
    due_date date NOT NULL,
    paid boolean DEFAULT false
);


--
-- Name: billing_schedule_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.billing_schedule_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: billing_schedule_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.billing_schedule_id_seq OWNED BY public.billing_schedule.id;


--
-- Name: template_table; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.template_table (
    id integer NOT NULL,
    createdat timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updatedat timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    deletedat timestamp without time zone
);


--
-- Name: template_table_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.template_table_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: template_table_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.template_table_id_seq OWNED BY public.template_table.id;


--
-- Name: borrowers; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.borrowers (
    id integer DEFAULT nextval('public.template_table_id_seq'::regclass) NOT NULL,
    createdat timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updatedat timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    deletedat timestamp without time zone,
    name character varying(255) NOT NULL,
    email character varying(255) NOT NULL,
    phone character varying(20) NOT NULL
);


--
-- Name: loans; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.loans (
    id integer DEFAULT nextval('public.template_table_id_seq'::regclass) NOT NULL,
    createdat timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updatedat timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    deletedat timestamp without time zone,
    borrower_id integer NOT NULL,
    amount numeric(15,2) NOT NULL,
    interest_rate numeric(5,2) NOT NULL,
    duration_weeks integer NOT NULL,
    outstanding numeric(15,2) NOT NULL,
    delinquent_weeks integer NOT NULL,
    installment_amount numeric(15,2)
);


--
-- Name: schema_migrations; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.schema_migrations (
    version character varying(128) NOT NULL
);


--
-- Name: billing_schedule id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.billing_schedule ALTER COLUMN id SET DEFAULT nextval('public.billing_schedule_id_seq'::regclass);


--
-- Name: template_table id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.template_table ALTER COLUMN id SET DEFAULT nextval('public.template_table_id_seq'::regclass);


--
-- Name: billing_schedule billing_schedule_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.billing_schedule
    ADD CONSTRAINT billing_schedule_pkey PRIMARY KEY (id);


--
-- Name: borrowers borrowers_email_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.borrowers
    ADD CONSTRAINT borrowers_email_key UNIQUE (email);


--
-- Name: borrowers borrowers_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.borrowers
    ADD CONSTRAINT borrowers_pkey PRIMARY KEY (id);


--
-- Name: loans loans_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.loans
    ADD CONSTRAINT loans_pkey PRIMARY KEY (id);


--
-- Name: schema_migrations schema_migrations_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.schema_migrations
    ADD CONSTRAINT schema_migrations_pkey PRIMARY KEY (version);


--
-- Name: template_table template_table_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.template_table
    ADD CONSTRAINT template_table_pkey PRIMARY KEY (id);


--
-- Name: billing_schedule billing_schedule_loan_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.billing_schedule
    ADD CONSTRAINT billing_schedule_loan_id_fkey FOREIGN KEY (loan_id) REFERENCES public.loans(id);


--
-- Name: loans fk_borrower; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.loans
    ADD CONSTRAINT fk_borrower FOREIGN KEY (borrower_id) REFERENCES public.borrowers(id);


--
-- PostgreSQL database dump complete
--


--
-- Dbmate schema migrations
--

INSERT INTO public.schema_migrations (version) VALUES
    ('20240518034949'),
    ('20240518045226'),
    ('20240518045252'),
    ('20240518090238'),
    ('20240519035433'),
    ('20240519070436');
