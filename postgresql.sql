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
    value integer DEFAULT 0,
    unit character varying(32) NOT NULL
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
e10d7056e1964da897d97e5dd685d115	2023-05-07 20:32:58.749269+07	2023-05-07 20:32:58.790775+07	2023-05-07 20:32:58.792664+07	a102ed29f7114833befa8904b061e60e	0676194723ff46bda957ae4bdfab3c85	36bffb268d8a4a06bc92ee6ece927bce	12
f3d425dab36341a08b710ae206adbfa2	2023-05-07 20:32:58.754604+07	2023-05-07 20:32:58.790775+07	2023-05-07 20:32:58.792664+07	a102ed29f7114833befa8904b061e60e	0676194723ff46bda957ae4bdfab3c85	2a0e9c8b9a474e69ad9b5432579bfa7e	12
0a89602d6f6c47a49f8f13127353267e	2023-05-07 20:32:58.76112+07	2023-05-07 20:32:58.790775+07	2023-05-07 20:32:58.792664+07	a102ed29f7114833befa8904b061e60e	0676194723ff46bda957ae4bdfab3c85	9b433589f3954135a34633a56b3e2c12	3
\.


--
-- Data for Name: categories; Type: TABLE DATA; Schema: public; Owner: user
--

COPY public.categories (id, created_at, updated_at, deleted_at, product_id, category_id) FROM stdin;
78	2023-05-07 20:32:54.400765+07	2023-05-07 20:32:54.400765+07	\N	9b433589f3954135a34633a56b3e2c12	8abcfbfbc0334647934e54f6631493ce
79	2023-05-07 20:32:54.406868+07	2023-05-07 20:32:54.406868+07	\N	9b433589f3954135a34633a56b3e2c12	20f97f19c38041098a7707ddc318d338
80	2023-05-07 20:32:57.28081+07	2023-05-07 20:32:57.28081+07	\N	2a0e9c8b9a474e69ad9b5432579bfa7e	fef18112f36341d69f6175a1601c8264
81	2023-05-07 20:32:57.284141+07	2023-05-07 20:32:57.284141+07	\N	2a0e9c8b9a474e69ad9b5432579bfa7e	20f97f19c38041098a7707ddc318d338
82	2023-05-07 20:32:57.74855+07	2023-05-07 20:32:57.74855+07	\N	36bffb268d8a4a06bc92ee6ece927bce	ea32f6a8834948488c6561a0b4dd87d3
83	2023-05-07 20:32:57.750575+07	2023-05-07 20:32:57.750575+07	\N	36bffb268d8a4a06bc92ee6ece927bce	20f97f19c38041098a7707ddc318d338
84	2023-05-07 20:32:58.125175+07	2023-05-07 20:32:58.125175+07	\N	274b30be54c54319ac6c3e5c59f89f91	6e0e6259593547dab08724170314e757
85	2023-05-07 20:32:58.12733+07	2023-05-07 20:32:58.12733+07	\N	274b30be54c54319ac6c3e5c59f89f91	20f97f19c38041098a7707ddc318d338
\.


--
-- Data for Name: category; Type: TABLE DATA; Schema: public; Owner: user
--

COPY public.category (id, created_at, updated_at, deleted_at, name, description) FROM stdin;
8abcfbfbc0334647934e54f6631493ce	2023-05-07 20:32:54.398399+07	2023-05-07 20:32:54.398399+07	\N	pineapple	
20f97f19c38041098a7707ddc318d338	2023-05-07 20:32:54.404585+07	2023-05-07 20:32:54.404585+07	\N	raw	
fef18112f36341d69f6175a1601c8264	2023-05-07 20:32:57.277925+07	2023-05-07 20:32:57.277925+07	\N	apple	
ea32f6a8834948488c6561a0b4dd87d3	2023-05-07 20:32:57.746546+07	2023-05-07 20:32:57.746546+07	\N	papaya	
6e0e6259593547dab08724170314e757	2023-05-07 20:32:58.122802+07	2023-05-07 20:32:58.122802+07	\N	orange	
\.


--
-- Data for Name: nutrients; Type: TABLE DATA; Schema: public; Owner: user
--

COPY public.nutrients (id, created_at, updated_at, deleted_at, product_id, name, value, unit) FROM stdin;
170	2023-05-07 20:32:57.222895+07	2023-05-07 20:32:57.222895+07	\N	9b433589f3954135a34633a56b3e2c12	Carbohydrate, by difference	13	G
171	2023-05-07 20:32:57.23229+07	2023-05-07 20:32:57.23229+07	\N	9b433589f3954135a34633a56b3e2c12	Energy	50	KCAL
172	2023-05-07 20:32:57.234742+07	2023-05-07 20:32:57.234742+07	\N	9b433589f3954135a34633a56b3e2c12	Water	86	G
173	2023-05-07 20:32:57.237218+07	2023-05-07 20:32:57.237218+07	\N	9b433589f3954135a34633a56b3e2c12	Sugars, total including NLEA	9	G
174	2023-05-07 20:32:57.239927+07	2023-05-07 20:32:57.239927+07	\N	9b433589f3954135a34633a56b3e2c12	Fiber, total dietary	1	G
175	2023-05-07 20:32:57.243038+07	2023-05-07 20:32:57.243038+07	\N	9b433589f3954135a34633a56b3e2c12	Calcium, Ca	13	MG
176	2023-05-07 20:32:57.24565+07	2023-05-07 20:32:57.24565+07	\N	9b433589f3954135a34633a56b3e2c12	Magnesium, Mg	12	MG
177	2023-05-07 20:32:57.248455+07	2023-05-07 20:32:57.248455+07	\N	9b433589f3954135a34633a56b3e2c12	Phosphorus, P	8	MG
178	2023-05-07 20:32:57.251043+07	2023-05-07 20:32:57.251043+07	\N	9b433589f3954135a34633a56b3e2c12	Potassium, K	109	MG
179	2023-05-07 20:32:57.253414+07	2023-05-07 20:32:57.253414+07	\N	9b433589f3954135a34633a56b3e2c12	Sodium, Na	1	MG
180	2023-05-07 20:32:57.255728+07	2023-05-07 20:32:57.255728+07	\N	9b433589f3954135a34633a56b3e2c12	Vitamin A, RAE	3	UG
181	2023-05-07 20:32:57.257588+07	2023-05-07 20:32:57.257588+07	\N	9b433589f3954135a34633a56b3e2c12	Carotene, beta	35	UG
182	2023-05-07 20:32:57.259934+07	2023-05-07 20:32:57.259934+07	\N	9b433589f3954135a34633a56b3e2c12	Vitamin C, total ascorbic acid	47	MG
183	2023-05-07 20:32:57.261934+07	2023-05-07 20:32:57.261934+07	\N	9b433589f3954135a34633a56b3e2c12	Folate, total	18	UG
184	2023-05-07 20:32:57.26346+07	2023-05-07 20:32:57.26346+07	\N	9b433589f3954135a34633a56b3e2c12	Choline, total	5	MG
185	2023-05-07 20:32:57.265851+07	2023-05-07 20:32:57.265851+07	\N	9b433589f3954135a34633a56b3e2c12	Folate, food	18	UG
186	2023-05-07 20:32:57.267611+07	2023-05-07 20:32:57.267611+07	\N	9b433589f3954135a34633a56b3e2c12	Folate, DFE	18	UG
187	2023-05-07 20:32:57.699174+07	2023-05-07 20:32:57.699174+07	\N	2a0e9c8b9a474e69ad9b5432579bfa7e	Carbohydrate, by difference	14	G
188	2023-05-07 20:32:57.704693+07	2023-05-07 20:32:57.704693+07	\N	2a0e9c8b9a474e69ad9b5432579bfa7e	Energy	61	KCAL
189	2023-05-07 20:32:57.708161+07	2023-05-07 20:32:57.708161+07	\N	2a0e9c8b9a474e69ad9b5432579bfa7e	Water	84	G
190	2023-05-07 20:32:57.711385+07	2023-05-07 20:32:57.711385+07	\N	2a0e9c8b9a474e69ad9b5432579bfa7e	Sugars, total including NLEA	12	G
191	2023-05-07 20:32:57.715127+07	2023-05-07 20:32:57.715127+07	\N	2a0e9c8b9a474e69ad9b5432579bfa7e	Fiber, total dietary	2	G
192	2023-05-07 20:32:57.718217+07	2023-05-07 20:32:57.718217+07	\N	2a0e9c8b9a474e69ad9b5432579bfa7e	Calcium, Ca	5	MG
193	2023-05-07 20:32:57.720998+07	2023-05-07 20:32:57.720998+07	\N	2a0e9c8b9a474e69ad9b5432579bfa7e	Magnesium, Mg	5	MG
194	2023-05-07 20:32:57.723435+07	2023-05-07 20:32:57.723435+07	\N	2a0e9c8b9a474e69ad9b5432579bfa7e	Phosphorus, P	9	MG
195	2023-05-07 20:32:57.725608+07	2023-05-07 20:32:57.725608+07	\N	2a0e9c8b9a474e69ad9b5432579bfa7e	Potassium, K	104	MG
196	2023-05-07 20:32:57.727577+07	2023-05-07 20:32:57.727577+07	\N	2a0e9c8b9a474e69ad9b5432579bfa7e	Vitamin A, RAE	3	UG
197	2023-05-07 20:32:57.729573+07	2023-05-07 20:32:57.729573+07	\N	2a0e9c8b9a474e69ad9b5432579bfa7e	Carotene, beta	27	UG
198	2023-05-07 20:32:57.731382+07	2023-05-07 20:32:57.731382+07	\N	2a0e9c8b9a474e69ad9b5432579bfa7e	Cryptoxanthin, beta	11	UG
199	2023-05-07 20:32:57.733442+07	2023-05-07 20:32:57.733442+07	\N	2a0e9c8b9a474e69ad9b5432579bfa7e	Lutein + zeaxanthin	29	UG
200	2023-05-07 20:32:57.735319+07	2023-05-07 20:32:57.735319+07	\N	2a0e9c8b9a474e69ad9b5432579bfa7e	Vitamin C, total ascorbic acid	4	MG
201	2023-05-07 20:32:57.736882+07	2023-05-07 20:32:57.736882+07	\N	2a0e9c8b9a474e69ad9b5432579bfa7e	Choline, total	3	MG
202	2023-05-07 20:32:57.738645+07	2023-05-07 20:32:57.738645+07	\N	2a0e9c8b9a474e69ad9b5432579bfa7e	Vitamin K (phylloquinone)	2	UG
203	2023-05-07 20:32:58.066968+07	2023-05-07 20:32:58.066968+07	\N	36bffb268d8a4a06bc92ee6ece927bce	Carbohydrate, by difference	10	G
204	2023-05-07 20:32:58.071034+07	2023-05-07 20:32:58.071034+07	\N	36bffb268d8a4a06bc92ee6ece927bce	Energy	43	KCAL
205	2023-05-07 20:32:58.075045+07	2023-05-07 20:32:58.075045+07	\N	36bffb268d8a4a06bc92ee6ece927bce	Water	88	G
206	2023-05-07 20:32:58.077991+07	2023-05-07 20:32:58.077991+07	\N	36bffb268d8a4a06bc92ee6ece927bce	Sugars, total including NLEA	7	G
207	2023-05-07 20:32:58.081137+07	2023-05-07 20:32:58.081137+07	\N	36bffb268d8a4a06bc92ee6ece927bce	Fiber, total dietary	1	G
208	2023-05-07 20:32:58.084194+07	2023-05-07 20:32:58.084194+07	\N	36bffb268d8a4a06bc92ee6ece927bce	Calcium, Ca	20	MG
209	2023-05-07 20:32:58.086938+07	2023-05-07 20:32:58.086938+07	\N	36bffb268d8a4a06bc92ee6ece927bce	Magnesium, Mg	21	MG
210	2023-05-07 20:32:58.089324+07	2023-05-07 20:32:58.089324+07	\N	36bffb268d8a4a06bc92ee6ece927bce	Phosphorus, P	10	MG
211	2023-05-07 20:32:58.091864+07	2023-05-07 20:32:58.091864+07	\N	36bffb268d8a4a06bc92ee6ece927bce	Potassium, K	182	MG
212	2023-05-07 20:32:58.093415+07	2023-05-07 20:32:58.093415+07	\N	36bffb268d8a4a06bc92ee6ece927bce	Sodium, Na	8	MG
213	2023-05-07 20:32:58.095415+07	2023-05-07 20:32:58.095415+07	\N	36bffb268d8a4a06bc92ee6ece927bce	Vitamin A, RAE	47	UG
214	2023-05-07 20:32:58.097197+07	2023-05-07 20:32:58.097197+07	\N	36bffb268d8a4a06bc92ee6ece927bce	Carotene, beta	274	UG
215	2023-05-07 20:32:58.098839+07	2023-05-07 20:32:58.098839+07	\N	36bffb268d8a4a06bc92ee6ece927bce	Carotene, alpha	2	UG
216	2023-05-07 20:32:58.100365+07	2023-05-07 20:32:58.100365+07	\N	36bffb268d8a4a06bc92ee6ece927bce	Cryptoxanthin, beta	589	UG
217	2023-05-07 20:32:58.102598+07	2023-05-07 20:32:58.102598+07	\N	36bffb268d8a4a06bc92ee6ece927bce	Lycopene	1830	UG
218	2023-05-07 20:32:58.104077+07	2023-05-07 20:32:58.104077+07	\N	36bffb268d8a4a06bc92ee6ece927bce	Lutein + zeaxanthin	89	UG
219	2023-05-07 20:32:58.105815+07	2023-05-07 20:32:58.105815+07	\N	36bffb268d8a4a06bc92ee6ece927bce	Vitamin C, total ascorbic acid	60	MG
220	2023-05-07 20:32:58.107628+07	2023-05-07 20:32:58.107628+07	\N	36bffb268d8a4a06bc92ee6ece927bce	Folate, total	37	UG
221	2023-05-07 20:32:58.109844+07	2023-05-07 20:32:58.109844+07	\N	36bffb268d8a4a06bc92ee6ece927bce	Choline, total	6	MG
222	2023-05-07 20:32:58.111458+07	2023-05-07 20:32:58.111458+07	\N	36bffb268d8a4a06bc92ee6ece927bce	Vitamin K (phylloquinone)	2	UG
223	2023-05-07 20:32:58.113063+07	2023-05-07 20:32:58.113063+07	\N	36bffb268d8a4a06bc92ee6ece927bce	Folate, food	37	UG
224	2023-05-07 20:32:58.114624+07	2023-05-07 20:32:58.114624+07	\N	36bffb268d8a4a06bc92ee6ece927bce	Folate, DFE	37	UG
225	2023-05-07 20:32:58.693353+07	2023-05-07 20:32:58.693353+07	\N	274b30be54c54319ac6c3e5c59f89f91	Carbohydrate, by difference	11	G
226	2023-05-07 20:32:58.698811+07	2023-05-07 20:32:58.698811+07	\N	274b30be54c54319ac6c3e5c59f89f91	Energy	50	KCAL
227	2023-05-07 20:32:58.702037+07	2023-05-07 20:32:58.702037+07	\N	274b30be54c54319ac6c3e5c59f89f91	Water	86	G
228	2023-05-07 20:32:58.704341+07	2023-05-07 20:32:58.704341+07	\N	274b30be54c54319ac6c3e5c59f89f91	Sugars, total including NLEA	8	G
229	2023-05-07 20:32:58.706279+07	2023-05-07 20:32:58.706279+07	\N	274b30be54c54319ac6c3e5c59f89f91	Fiber, total dietary	2	G
230	2023-05-07 20:32:58.708132+07	2023-05-07 20:32:58.708132+07	\N	274b30be54c54319ac6c3e5c59f89f91	Calcium, Ca	42	MG
231	2023-05-07 20:32:58.710267+07	2023-05-07 20:32:58.710267+07	\N	274b30be54c54319ac6c3e5c59f89f91	Magnesium, Mg	10	MG
232	2023-05-07 20:32:58.712391+07	2023-05-07 20:32:58.712391+07	\N	274b30be54c54319ac6c3e5c59f89f91	Phosphorus, P	18	MG
233	2023-05-07 20:32:58.714505+07	2023-05-07 20:32:58.714505+07	\N	274b30be54c54319ac6c3e5c59f89f91	Potassium, K	174	MG
234	2023-05-07 20:32:58.717599+07	2023-05-07 20:32:58.717599+07	\N	274b30be54c54319ac6c3e5c59f89f91	Sodium, Na	4	MG
235	2023-05-07 20:32:58.719274+07	2023-05-07 20:32:58.719274+07	\N	274b30be54c54319ac6c3e5c59f89f91	Vitamin A, RAE	11	UG
236	2023-05-07 20:32:58.721088+07	2023-05-07 20:32:58.721088+07	\N	274b30be54c54319ac6c3e5c59f89f91	Carotene, beta	71	UG
237	2023-05-07 20:32:58.722816+07	2023-05-07 20:32:58.722816+07	\N	274b30be54c54319ac6c3e5c59f89f91	Carotene, alpha	11	UG
238	2023-05-07 20:32:58.724589+07	2023-05-07 20:32:58.724589+07	\N	274b30be54c54319ac6c3e5c59f89f91	Cryptoxanthin, beta	116	UG
239	2023-05-07 20:32:58.726676+07	2023-05-07 20:32:58.726676+07	\N	274b30be54c54319ac6c3e5c59f89f91	Lutein + zeaxanthin	129	UG
240	2023-05-07 20:32:58.72899+07	2023-05-07 20:32:58.72899+07	\N	274b30be54c54319ac6c3e5c59f89f91	Vitamin C, total ascorbic acid	56	MG
241	2023-05-07 20:32:58.731484+07	2023-05-07 20:32:58.731484+07	\N	274b30be54c54319ac6c3e5c59f89f91	Folate, total	28	UG
242	2023-05-07 20:32:58.733525+07	2023-05-07 20:32:58.733525+07	\N	274b30be54c54319ac6c3e5c59f89f91	Choline, total	8	MG
243	2023-05-07 20:32:58.735301+07	2023-05-07 20:32:58.735301+07	\N	274b30be54c54319ac6c3e5c59f89f91	Folate, food	28	UG
244	2023-05-07 20:32:58.736767+07	2023-05-07 20:32:58.736767+07	\N	274b30be54c54319ac6c3e5c59f89f91	Folate, DFE	28	UG
\.


--
-- Data for Name: products; Type: TABLE DATA; Schema: public; Owner: user
--

COPY public.products (id, created_at, updated_at, deleted_at, name, image, description, stocks, price) FROM stdin;
36bffb268d8a4a06bc92ee6ece927bce	2023-05-07 20:32:57.743913+07	2023-05-07 22:27:21.347082+07	\N	Papaya	papaya.png	The papaya, papaw, or pawpaw is the plant species Carica papaya, one of the 21 accepted species in the genus Carica of the family Caricaceae.	52	23
274b30be54c54319ac6c3e5c59f89f91	2023-05-07 20:32:58.11996+07	2023-05-07 20:32:58.11996+07	\N	Orange		An orange is a fruit of various citrus species in the family Rutaceae, it primarily refers to Citrus × sinensis, which is also called sweet orange.	120	6
2a0e9c8b9a474e69ad9b5432579bfa7e	2023-05-07 20:32:57.274974+07	2023-05-07 20:32:57.274974+07	\N	Apple		An edible fruit produced by an apple tree (Malus Domestica).	228	6
9b433589f3954135a34633a56b3e2c12	2023-05-07 20:32:54.395695+07	2023-05-07 20:32:54.395695+07	\N	Pineapple		A tropical plant with an edible fruit and the most economically significant plant in the family Bromeliaceae.	125	12
\.


--
-- Data for Name: sessions; Type: TABLE DATA; Schema: public; Owner: user
--

COPY public.sessions (id, created_at, updated_at, deleted_at, user_id, client_ip, user_agent, token, secret_key, expired, last_activated) FROM stdin;
41050c43f0c34f598101913a0904fd80	2023-05-07 22:18:32.282816+07	2023-05-07 22:27:21.33943+07	\N	ea6cef1735a34df4884c4e1dc4885919	127.0.0.1	PostmanRuntime/7.32.2	eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhZG1pbiI6dHJ1ZSwiZW1haWwiOiJhZG1pbkBtYWlsLmNvIiwiZXhwIjoxNjgzNDg3MTEyLCJpYXQiOjE2ODM0NzI3MTIsInVzZXJuYW1lIjoiYWRtaW4ifQ.wsnjAR16qDfbg_0qnoISq3s60_WwtNdGw-FTSTQnGu8	PKPvLtYRDuQZbeOb7UBSqOGTOApjbMF1Gk_xayvl3tU=	2023-05-07 19:18:32.229852	2023-05-07 15:27:21.339259
4754a990b4724040bebf91d23bcb8c03	2023-05-07 20:32:54.336331+07	2023-05-07 20:32:58.797731+07	\N	a102ed29f7114833befa8904b061e60e	127.0.0.1	Go-http-client/1.1	eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhZG1pbiI6dHJ1ZSwiZW1haWwiOiJ1c2VyQG1haWwuY28iLCJleHAiOjE2ODM0ODA3NzQsImlhdCI6MTY4MzQ2NjM3NCwidXNlcm5hbWUiOiJ1c2VyIn0.UB_L2AUB6bui0Q5tqYhfji1FDobQ3bAU8gSWHa-q7fU	2J7Z-ESkMWIhTBB_ey1gLzJ1CXb1n1WXh-TpIP6i0ns=	2023-05-07 17:32:54.288079	2023-05-07 13:32:58.797617
b7a5eea59e9d4544b8c97654f443cdcb	2023-05-07 20:33:23.640188+07	2023-05-07 20:34:08.463163+07	\N	a102ed29f7114833befa8904b061e60e	127.0.0.1	PostmanRuntime/7.32.2	eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhZG1pbiI6dHJ1ZSwiZW1haWwiOiJ1c2VyQG1haWwuY28iLCJleHAiOjE2ODM0ODA4MDMsImlhdCI6MTY4MzQ2NjQwMywidXNlcm5hbWUiOiJ1c2VyIn0.2Lrd0OnbhlLBLJblG0encYOg3iMNji3OsOtIknjSyEM	64S2jeUTnFKoilJcWTJHpaRFv-X6sTq_V8De4b20UeM=	2023-05-07 17:33:23.588414	2023-05-07 13:34:08.462686
5610926039c144bea643a2bb15d9b599	2023-05-07 20:32:54.387992+07	2023-05-07 20:32:58.117925+07	2023-05-07 21:06:53.406297+07	ea6cef1735a34df4884c4e1dc4885919	127.0.0.1	Go-http-client/1.1	eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhZG1pbiI6dHJ1ZSwiZW1haWwiOiJhZG1pbkBtYWlsLmNvIiwiZXhwIjoxNjgzNDgwNzc0LCJpYXQiOjE2ODM0NjYzNzQsInVzZXJuYW1lIjoiYWRtaW4ifQ.ACjAJYZ2te04WuQsyQb6GEZuwa8-2BGUYGRtxLUI0t0	Tk7D1ECU7Rlxgtj0NG2Kd-jSz7bS4RguV18ozxdHo9M=	2023-05-07 17:32:54.340112	2023-05-07 13:32:58.11785
35894d962f124ce0a7deac8536740342	2023-05-07 20:58:00.588851+07	2023-05-07 21:33:20.608155+07	2023-05-07 22:18:29.038533+07	ea6cef1735a34df4884c4e1dc4885919	127.0.0.1	PostmanRuntime/7.32.2	eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhZG1pbiI6dHJ1ZSwiZW1haWwiOiJhZG1pbkBtYWlsLmNvIiwiZXhwIjoxNjgzNDgyMjgwLCJpYXQiOjE2ODM0Njc4ODAsInVzZXJuYW1lIjoiYWRtaW4ifQ.cQdT6UUP2VAVnaeWSyCjlwdAXUB-zpLKQA9oIdiUeFk	UiQSzyJxtcqcbo1zij2hvFfBUROD81MvUOGmlq19cXk=	2023-05-07 17:58:00.534511	2023-05-07 14:33:20.607821
\.


--
-- Data for Name: transactions; Type: TABLE DATA; Schema: public; Owner: user
--

COPY public.transactions (id, created_at, updated_at, deleted_at, user_id, payment_method, verify) FROM stdin;
0676194723ff46bda957ae4bdfab3c85	2023-05-07 20:32:58.783535+07	2023-05-07 20:32:58.79451+07	\N	a102ed29f7114833befa8904b061e60e	visa-credit-card	t
\.


--
-- Data for Name: users; Type: TABLE DATA; Schema: public; Owner: user
--

COPY public.users (id, created_at, updated_at, deleted_at, name, username, email, password, gender, phone, dob, address, country_code, city, postal_code, admin, verify, balance) FROM stdin;
ea6cef1735a34df4884c4e1dc4885919	\N	\N	\N	\N	admin	admin@mail.co	$2a$10$u6sU3JFeWgd6ET7Ewwf5xuHzAmQhuBOpl.lBtSeI6GfO7bK5TTenW	\N	\N	\N	\N	\N	\N	\N	t	f	0
a102ed29f7114833befa8904b061e60e	\N	2023-05-07 20:33:36.669532+07	\N	Ujang, Slamet	user	user@mail.co	$2a$10$u6sU3JFeWgd6ET7Ewwf5xuHzAmQhuBOpl.lBtSeI6GfO7bK5TTenW	M	01234567890	2023-01-01 07:00:00	Heaven Road	ID	Semarang	00000	t	f	616
\.


--
-- Name: categories_id_seq; Type: SEQUENCE SET; Schema: public; Owner: user
--

SELECT pg_catalog.setval('public.categories_id_seq', 85, true);


--
-- Name: nutrients_id_seq; Type: SEQUENCE SET; Schema: public; Owner: user
--

SELECT pg_catalog.setval('public.nutrients_id_seq', 244, true);


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

