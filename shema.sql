
--CREATE DATABASE "News" WITH TEMPLATE = template0 ;



CREATE TABLE public."NewsTable" (
    "Id" integer NOT NULL,
    "Title" character varying NOT NULL,
    "Description" character varying NOT NULL,
    "Time" bigint NOT NULL,
    "Url" character varying,
    "GUID" character varying
);


ALTER TABLE public."NewsTable" OWNER TO postgres;


CREATE SEQUENCE public."NewsTable_Id_seq"
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public."NewsTable_Id_seq" OWNER TO postgres;

ALTER SEQUENCE public."NewsTable_Id_seq" OWNED BY public."NewsTable"."Id";



CREATE TABLE public."Settings" (
    "ID" integer NOT NULL,
    "Description" character varying,
    "ValueStr" character varying,
    "TypeSetting" integer NOT NULL,
    "ValueInt" integer,
    "Active" boolean DEFAULT false NOT NULL
);


ALTER TABLE public."Settings" OWNER TO postgres;


CREATE SEQUENCE public."Settings_ID_seq"
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public."Settings_ID_seq" OWNER TO postgres;


ALTER SEQUENCE public."Settings_ID_seq" OWNED BY public."Settings"."ID";


ALTER TABLE ONLY public."NewsTable" ALTER COLUMN "Id" SET DEFAULT nextval('public."NewsTable_Id_seq"'::regclass);



ALTER TABLE ONLY public."Settings" ALTER COLUMN "ID" SET DEFAULT nextval('public."Settings_ID_seq"'::regclass);


SELECT pg_catalog.setval('public."NewsTable_Id_seq"', 1, false);

SELECT pg_catalog.setval('public."Settings_ID_seq"', 4, true);


ALTER TABLE ONLY public."NewsTable"
    ADD CONSTRAINT "NewsTable_pk" PRIMARY KEY ("Id");

CREATE UNIQUE INDEX "NewsTable_GUID_IDX" ON public."NewsTable" USING btree ("GUID");

INSERT INTO public."Settings" VALUES (1, 'URL', 'https://habr.com/ru/rss/hub/go/all/?fl=ru', 1, NULL, true);
INSERT INTO public."Settings" VALUES (2, 'URL', 'https://habr.com/ru/rss/hub/go/all/?fl=ru', 1, NULL, true);
INSERT INTO public."Settings" VALUES (3, 'URL', 'http://www.bashorg.org/rss.xml', 1, NULL, true);
INSERT INTO public."Settings" VALUES (4, 'Time', NULL, 2, 300, true);
