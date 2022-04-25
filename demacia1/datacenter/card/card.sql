CREATE SCHEMA "card" AUTHORIZATION "postgres";

-- ----------------------------
-- Sequence structure for card_id_seq
-- ----------------------------
CREATE SEQUENCE "card"."card_id_seq"
    INCREMENT 1
    MINVALUE  1
    MAXVALUE 2147483647
    START 1
    CACHE 1;
ALTER SEQUENCE "card"."card_id_seq" OWNER TO "postgres";

-- ----------------------------
-- Table structure for card
-- ----------------------------
CREATE TABLE "card"."card" (
                               "id" int4 NOT NULL DEFAULT nextval('"card".card_id_seq'::regclass),
                               "card_num" varchar(255) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
                               "object_role" int2 NOT NULL,
                               "object_id" int4 NOT NULL,
                               "org_id" int4 NOT NULL DEFAULT 0
);
ALTER TABLE "card"."card" OWNER TO "postgres";
COMMENT ON COLUMN "card"."card"."object_role" IS '1：教师，2：学生';

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "card"."card_id_seq" OWNED BY "card"."card"."id";
SELECT setval('"card"."card_id_seq"', 52, true);

-- ----------------------------
-- Primary Key structure for table card
-- ----------------------------
ALTER TABLE "card"."card" ADD CONSTRAINT "card_pkey" PRIMARY KEY ("id");
