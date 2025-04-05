CREATE TABLE "public"."legal_entities" (
    "uuid" uuid DEFAULT gen_random_uuid() PRIMARY KEY,
    "name" varchar(50) NOT NULL DEFAULT '',
    "created_at" timestamptz NOT NULL DEFAULT now(),
    "updated_at" timestamptz NOT NULL DEFAULT now(),
    "deleted_at" timestamptz
);