--
-- PostgreSQL database dump
--

-- Dumped from database version 14.5 (Ubuntu 14.5-1ubuntu1)
-- Dumped by pg_dump version 14.5 (Ubuntu 14.5-1ubuntu1)

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
-- Name: uuid-ossp; Type: EXTENSION; Schema: -; Owner: -
--

CREATE EXTENSION IF NOT EXISTS "uuid-ossp" WITH SCHEMA public;


--
-- Name: EXTENSION "uuid-ossp"; Type: COMMENT; Schema: -; Owner: 
--

COMMENT ON EXTENSION "uuid-ossp" IS 'generate universally unique identifiers (UUIDs)';


SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: carts; Type: TABLE; Schema: public; Owner: user
--

CREATE TABLE public.carts (
    id character varying(32) NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    user_id character varying(32) NOT NULL,
    transaction_id character varying(32),
    product_id character varying(32) NOT NULL,
    qty bigint
);


ALTER TABLE public.carts OWNER TO "user";

--
-- Name: categories; Type: TABLE; Schema: public; Owner: user
--

CREATE TABLE public.categories (
    id bigint NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    product_id character varying(32) NOT NULL,
    category_id character varying(32) NOT NULL
);


ALTER TABLE public.categories OWNER TO "user";

--
-- Name: categories_id_seq; Type: SEQUENCE; Schema: public; Owner: user
--

CREATE SEQUENCE public.categories_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.categories_id_seq OWNER TO "user";

--
-- Name: categories_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: user
--

ALTER SEQUENCE public.categories_id_seq OWNED BY public.categories.id;


--
-- Name: category; Type: TABLE; Schema: public; Owner: user
--

CREATE TABLE public.category (
    id character varying(32) NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    name character varying(32) NOT NULL,
    description text
);


ALTER TABLE public.category OWNER TO "user";

--
-- Name: nutrients; Type: TABLE; Schema: public; Owner: user
--

CREATE TABLE public.nutrients (
    id bigint NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    product_id character varying(32) NOT NULL,
    name character varying(32) NOT NULL,
    value integer DEFAULT 0
);


ALTER TABLE public.nutrients OWNER TO "user";

--
-- Name: nutrients_id_seq; Type: SEQUENCE; Schema: public; Owner: user
--

CREATE SEQUENCE public.nutrients_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.nutrients_id_seq OWNER TO "user";

--
-- Name: nutrients_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: user
--

ALTER SEQUENCE public.nutrients_id_seq OWNED BY public.nutrients.id;


--
-- Name: products; Type: TABLE; Schema: public; Owner: user
--

CREATE TABLE public.products (
    id character varying(32) NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    name text NOT NULL,
    image text,
    description text,
    stocks bigint,
    price text
);


ALTER TABLE public.products OWNER TO "user";

--
-- Name: sessions; Type: TABLE; Schema: public; Owner: user
--

CREATE TABLE public.sessions (
    id character varying(32) NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    user_id character varying(32) NOT NULL,
    client_ip character varying(40),
    user_agent text,
    token text NOT NULL,
    secret_key text NOT NULL,
    expired timestamp without time zone NOT NULL,
    last_activated timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);


ALTER TABLE public.sessions OWNER TO "user";

--
-- Name: transactions; Type: TABLE; Schema: public; Owner: user
--

CREATE TABLE public.transactions (
    id character varying(32) NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    user_id character varying(32) NOT NULL,
    payment_method text,
    verify boolean DEFAULT false
);


ALTER TABLE public.transactions OWNER TO "user";

--
-- Name: users; Type: TABLE; Schema: public; Owner: user
--

CREATE TABLE public.users (
    id character varying(32) NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    name character varying(52),
    username character varying(16) NOT NULL,
    email character varying(254) NOT NULL,
    password character varying(128) NOT NULL,
    gender character varying(1),
    phone character varying(24),
    dob timestamp without time zone,
    address character varying(128),
    country_code character varying(4),
    city character varying(64),
    postal_code character varying(10),
    admin boolean DEFAULT false,
    verify boolean DEFAULT false,
    balance text DEFAULT '0'::text
);


ALTER TABLE public.users OWNER TO "user";

--
-- Name: categories id; Type: DEFAULT; Schema: public; Owner: user
--

ALTER TABLE ONLY public.categories ALTER COLUMN id SET DEFAULT nextval('public.categories_id_seq'::regclass);


--
-- Name: nutrients id; Type: DEFAULT; Schema: public; Owner: user
--

ALTER TABLE ONLY public.nutrients ALTER COLUMN id SET DEFAULT nextval('public.nutrients_id_seq'::regclass);


--
-- Name: carts carts_pkey; Type: CONSTRAINT; Schema: public; Owner: user
--

ALTER TABLE ONLY public.carts
    ADD CONSTRAINT carts_pkey PRIMARY KEY (id);


--
-- Name: categories categories_pkey; Type: CONSTRAINT; Schema: public; Owner: user
--

ALTER TABLE ONLY public.categories
    ADD CONSTRAINT categories_pkey PRIMARY KEY (id);


--
-- Name: category category_pkey; Type: CONSTRAINT; Schema: public; Owner: user
--

ALTER TABLE ONLY public.category
    ADD CONSTRAINT category_pkey PRIMARY KEY (id);


--
-- Name: nutrients nutrients_pkey; Type: CONSTRAINT; Schema: public; Owner: user
--

ALTER TABLE ONLY public.nutrients
    ADD CONSTRAINT nutrients_pkey PRIMARY KEY (id);


--
-- Name: products products_pkey; Type: CONSTRAINT; Schema: public; Owner: user
--

ALTER TABLE ONLY public.products
    ADD CONSTRAINT products_pkey PRIMARY KEY (id);


--
-- Name: sessions sessions_pkey; Type: CONSTRAINT; Schema: public; Owner: user
--

ALTER TABLE ONLY public.sessions
    ADD CONSTRAINT sessions_pkey PRIMARY KEY (id);


--
-- Name: sessions sessions_secret_key_key; Type: CONSTRAINT; Schema: public; Owner: user
--

ALTER TABLE ONLY public.sessions
    ADD CONSTRAINT sessions_secret_key_key UNIQUE (secret_key);


--
-- Name: sessions sessions_token_key; Type: CONSTRAINT; Schema: public; Owner: user
--

ALTER TABLE ONLY public.sessions
    ADD CONSTRAINT sessions_token_key UNIQUE (token);


--
-- Name: transactions transactions_pkey; Type: CONSTRAINT; Schema: public; Owner: user
--

ALTER TABLE ONLY public.transactions
    ADD CONSTRAINT transactions_pkey PRIMARY KEY (id);


--
-- Name: users users_email_key; Type: CONSTRAINT; Schema: public; Owner: user
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_email_key UNIQUE (email);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: user
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- Name: users users_username_key; Type: CONSTRAINT; Schema: public; Owner: user
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_username_key UNIQUE (username);


--
-- Name: idx_carts_deleted_at; Type: INDEX; Schema: public; Owner: user
--

CREATE INDEX idx_carts_deleted_at ON public.carts USING btree (deleted_at);


--
-- Name: idx_categories_deleted_at; Type: INDEX; Schema: public; Owner: user
--

CREATE INDEX idx_categories_deleted_at ON public.categories USING btree (deleted_at);


--
-- Name: idx_category_deleted_at; Type: INDEX; Schema: public; Owner: user
--

CREATE INDEX idx_category_deleted_at ON public.category USING btree (deleted_at);


--
-- Name: idx_nutrients_deleted_at; Type: INDEX; Schema: public; Owner: user
--

CREATE INDEX idx_nutrients_deleted_at ON public.nutrients USING btree (deleted_at);


--
-- Name: idx_products_deleted_at; Type: INDEX; Schema: public; Owner: user
--

CREATE INDEX idx_products_deleted_at ON public.products USING btree (deleted_at);


--
-- Name: idx_sessions_deleted_at; Type: INDEX; Schema: public; Owner: user
--

CREATE INDEX idx_sessions_deleted_at ON public.sessions USING btree (deleted_at);


--
-- Name: idx_transactions_deleted_at; Type: INDEX; Schema: public; Owner: user
--

CREATE INDEX idx_transactions_deleted_at ON public.transactions USING btree (deleted_at);


--
-- Name: idx_users_deleted_at; Type: INDEX; Schema: public; Owner: user
--

CREATE INDEX idx_users_deleted_at ON public.users USING btree (deleted_at);


--
-- Name: carts fk_transactions_carts; Type: FK CONSTRAINT; Schema: public; Owner: user
--

ALTER TABLE ONLY public.carts
    ADD CONSTRAINT fk_transactions_carts FOREIGN KEY (transaction_id) REFERENCES public.transactions(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: carts fk_users_carts; Type: FK CONSTRAINT; Schema: public; Owner: user
--

ALTER TABLE ONLY public.carts
    ADD CONSTRAINT fk_users_carts FOREIGN KEY (user_id) REFERENCES public.users(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: sessions fk_users_sessions; Type: FK CONSTRAINT; Schema: public; Owner: user
--

ALTER TABLE ONLY public.sessions
    ADD CONSTRAINT fk_users_sessions FOREIGN KEY (user_id) REFERENCES public.users(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: transactions fk_users_transactions; Type: FK CONSTRAINT; Schema: public; Owner: user
--

ALTER TABLE ONLY public.transactions
    ADD CONSTRAINT fk_users_transactions FOREIGN KEY (user_id) REFERENCES public.users(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- PostgreSQL database dump complete
--

