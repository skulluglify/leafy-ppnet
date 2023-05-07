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
-- Data for Name: carts; Type: TABLE DATA; Schema: public; Owner: user
--

COPY public.carts (id, created_at, updated_at, deleted_at, user_id, transaction_id, product_id, qty) FROM stdin;
a19a65c0ffe2442a950601189f9e0a6b	2023-05-07 12:30:52.602404+07	2023-05-07 12:30:52.640008+07	2023-05-07 12:30:52.641782+07	6d02e33a79a04a9ead9b2faf707cef2a	f964a8819563479fb4a9ea1a38589d1c	b119a103f10c408cb13ce35ca4141459	12
fb6ef2c522384b5e8b152dff3134fb21	2023-05-07 12:30:52.607551+07	2023-05-07 12:30:52.640008+07	2023-05-07 12:30:52.641782+07	6d02e33a79a04a9ead9b2faf707cef2a	f964a8819563479fb4a9ea1a38589d1c	543d71b22cff46b68d0a5515d7909fc8	12
d9ea4b2a2a9447b7bba98d79f630696d	2023-05-07 12:30:52.612512+07	2023-05-07 12:30:52.640008+07	2023-05-07 12:30:52.641782+07	6d02e33a79a04a9ead9b2faf707cef2a	f964a8819563479fb4a9ea1a38589d1c	6f064ad3a99f48e88a58805718f2fc91	3
\.


--
-- Data for Name: categories; Type: TABLE DATA; Schema: public; Owner: user
--

COPY public.categories (id, created_at, updated_at, deleted_at, product_id, category_id) FROM stdin;
36	2023-05-07 12:30:52.563714+07	2023-05-07 12:30:52.563714+07	\N	6f064ad3a99f48e88a58805718f2fc91	e67a76cf339744d694fac64bdc421b8a
37	2023-05-07 12:30:52.567062+07	2023-05-07 12:30:52.567062+07	\N	6f064ad3a99f48e88a58805718f2fc91	b58c9170db5647a18daae18724fba7a6
38	2023-05-07 12:30:52.573914+07	2023-05-07 12:30:52.573914+07	\N	543d71b22cff46b68d0a5515d7909fc8	c162241ff2f048a0b919d6ae6c3ffb98
39	2023-05-07 12:30:52.575773+07	2023-05-07 12:30:52.575773+07	\N	543d71b22cff46b68d0a5515d7909fc8	b58c9170db5647a18daae18724fba7a6
40	2023-05-07 12:30:52.582873+07	2023-05-07 12:30:52.582873+07	\N	b119a103f10c408cb13ce35ca4141459	1390cf02d00a49fe9f51cdaac3464e06
41	2023-05-07 12:30:52.584683+07	2023-05-07 12:30:52.584683+07	\N	b119a103f10c408cb13ce35ca4141459	b58c9170db5647a18daae18724fba7a6
42	2023-05-07 12:30:52.592611+07	2023-05-07 12:30:52.592611+07	\N	3d6ed7e8019b44b0ae50135691b7ba96	a143baf91bd849aabbbda0fac3790ebc
43	2023-05-07 12:30:52.594497+07	2023-05-07 12:30:52.594497+07	\N	3d6ed7e8019b44b0ae50135691b7ba96	b58c9170db5647a18daae18724fba7a6
\.


--
-- Data for Name: category; Type: TABLE DATA; Schema: public; Owner: user
--

COPY public.category (id, created_at, updated_at, deleted_at, name, description) FROM stdin;
e67a76cf339744d694fac64bdc421b8a	2023-05-07 12:30:52.561957+07	2023-05-07 12:30:52.561957+07	\N	pineapple	
b58c9170db5647a18daae18724fba7a6	2023-05-07 12:30:52.565571+07	2023-05-07 12:30:52.565571+07	\N	raw	
c162241ff2f048a0b919d6ae6c3ffb98	2023-05-07 12:30:52.572623+07	2023-05-07 12:30:52.572623+07	\N	apple	
1390cf02d00a49fe9f51cdaac3464e06	2023-05-07 12:30:52.581354+07	2023-05-07 12:30:52.581354+07	\N	papaya	
a143baf91bd849aabbbda0fac3790ebc	2023-05-07 12:30:52.591074+07	2023-05-07 12:30:52.591074+07	\N	orange	
\.


--
-- Data for Name: nutrients; Type: TABLE DATA; Schema: public; Owner: user
--

COPY public.nutrients (id, created_at, updated_at, deleted_at, product_id, name, value) FROM stdin;
\.


--
-- Data for Name: products; Type: TABLE DATA; Schema: public; Owner: user
--

COPY public.products (id, created_at, updated_at, deleted_at, name, image, description, stocks, price) FROM stdin;
3d6ed7e8019b44b0ae50135691b7ba96	2023-05-07 12:30:52.589254+07	2023-05-07 12:30:52.589254+07	\N	Orange		An orange is a fruit of various citrus species in the family Rutaceae (see list of plants known as orange); it primarily refers to Citrus × sinensis,[1] which is also called sweet orange, to distinguish it from the related Citrus × aurantium, referred to as bitter orange.	120	6
b119a103f10c408cb13ce35ca4141459	2023-05-07 12:30:52.579606+07	2023-05-07 12:30:52.579606+07	\N	Papaya		The papaya, papaw, or pawpaw is the plant species Carica papaya, one of the 21 accepted species in the genus Carica of the family Caricaceae.	52	23
543d71b22cff46b68d0a5515d7909fc8	2023-05-07 12:30:52.571084+07	2023-05-07 12:30:52.571084+07	\N	Apple		An edible fruit produced by an apple tree (Malus Domestica).	228	6
6f064ad3a99f48e88a58805718f2fc91	2023-05-07 12:30:52.559721+07	2023-05-07 12:30:52.559721+07	\N	Pineapple		A tropical plant with an edible fruit and the most economically significant plant in the family Bromeliaceae.	125	12
\.


--
-- Data for Name: sessions; Type: TABLE DATA; Schema: public; Owner: user
--

COPY public.sessions (id, created_at, updated_at, deleted_at, user_id, client_ip, user_agent, token, secret_key, expired, last_activated) FROM stdin;
ae2cc72c915743c1a67753203ff37f03	2023-05-07 12:30:52.499312+07	2023-05-07 12:30:52.646038+07	\N	6d02e33a79a04a9ead9b2faf707cef2a	127.0.0.1	Go-http-client/1.1	eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhZG1pbiI6dHJ1ZSwiZW1haWwiOiJ1c2VyQG1haWwuY28iLCJleHAiOjE2ODM0NTE4NTIsImlhdCI6MTY4MzQzNzQ1MiwidXNlcm5hbWUiOiJ1c2VyIn0.AUT0BZbRrUPFDAHxlrzDJE0g_eMzvY3-VcToZ6x2fwY	GkHsb1t-7vXI6m0wwptxgZry4nkC9EgoCJ4cm81PiUM=	2023-05-07 09:30:52.446855	2023-05-07 05:30:52.645938
d69e962635584db2af30f634d65aad96	2023-05-07 12:30:52.553863+07	2023-05-07 12:30:52.58731+07	\N	a4dd3bab9fd84cf188e6c4093698690f	127.0.0.1	Go-http-client/1.1	eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhZG1pbiI6dHJ1ZSwiZW1haWwiOiJhZG1pbkBtYWlsLmNvIiwiZXhwIjoxNjgzNDUxODUyLCJpYXQiOjE2ODM0Mzc0NTIsInVzZXJuYW1lIjoiYWRtaW4ifQ.iD-yghUhM0jfXcofSQuxq1fGE2P8aW9zDIaiQ9EDYvo	I6kTs8knQQWQVRSbxb1ubo-o6pvdOPVqWhMbp-kbsvg=	2023-05-07 09:30:52.502988	2023-05-07 05:30:52.587211
\.


--
-- Data for Name: transactions; Type: TABLE DATA; Schema: public; Owner: user
--

COPY public.transactions (id, created_at, updated_at, deleted_at, user_id, payment_method, verify) FROM stdin;
f964a8819563479fb4a9ea1a38589d1c	2023-05-07 12:30:52.63254+07	2023-05-07 12:30:52.643456+07	\N	6d02e33a79a04a9ead9b2faf707cef2a	visa-credit-card	t
\.


--
-- Data for Name: users; Type: TABLE DATA; Schema: public; Owner: user
--

COPY public.users (id, created_at, updated_at, deleted_at, name, username, email, password, gender, phone, dob, address, country_code, city, postal_code, admin, verify, balance) FROM stdin;
a4dd3bab9fd84cf188e6c4093698690f	\N	\N	\N	\N	admin	admin@mail.co	$2a$10$dMF8Q6vCvgxynyttdS.lY./5rSW9tBs4DHJmqG8uXBP8i4miH.RRG	\N	\N	\N	\N	\N	\N	\N	t	f	0
6d02e33a79a04a9ead9b2faf707cef2a	\N	2023-05-07 12:30:52.630164+07	\N	\N	user	user@mail.co	$2a$10$dMF8Q6vCvgxynyttdS.lY./5rSW9tBs4DHJmqG8uXBP8i4miH.RRG	\N	\N	\N	\N	\N	\N	\N	t	f	616
\.


--
-- Name: categories_id_seq; Type: SEQUENCE SET; Schema: public; Owner: user
--

SELECT pg_catalog.setval('public.categories_id_seq', 43, true);


--
-- Name: nutrients_id_seq; Type: SEQUENCE SET; Schema: public; Owner: user
--

SELECT pg_catalog.setval('public.nutrients_id_seq', 1, false);


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

