CREATE TABLE "public"."bank_accounts" (
    "uuid" uuid DEFAULT gen_random_uuid() PRIMARY KEY,
    "legal_entity_uuid" uuid NOT NULL,
    "bic" varchar(9) NOT NULL,
    "bank_name" varchar(255) NOT NULL,
    "bank_address" varchar(255) NOT NULL,
    "corr_account" varchar(20) NOT NULL,
    "payment_account" varchar(20) NOT NULL,
    "currency" varchar(3) NOT NULL,
    "comment" text,
    "created_at" timestamptz NOT NULL DEFAULT now(),
    "updated_at" timestamptz NOT NULL DEFAULT now(),
    "deleted_at" timestamptz
);

-- migrations/xxxx_add_bank_accounts_relation.up.sql
ALTER TABLE "public"."bank_accounts"
ADD CONSTRAINT "fk_bank_accounts_legal_entity"
FOREIGN KEY ("legal_entity_uuid") 
REFERENCES "public"."legal_entities"("uuid")
ON DELETE CASCADE;
