CREATE TABLE IF NOT EXISTS users
(
    "id"         uuid         NOT NULL PRIMARY KEY,
    "user_id"    varchar(255) NOT NULL,
    "chat_id"    varchar(255) NOT NULL,
    "username"   varchar(255) NOT NULL,
    "first_name" varchar(255),
    "last_name"  varchar(255),
    "created_at" timestamp    NOT NULL,
    "updated_at" timestamp    NOT NULL
);

CREATE UNIQUE INDEX idx_user_id_chat_id ON users (user_id, chat_id);

-- Check docs for information about field lengths and types:
-- https://www.financnasprava.sk/_img/pfsedit/Dokumenty_PFS/Podnikatelia/eKasa/2019/2019.05.27_eKasa_rozhranie.pdf
CREATE TABLE IF NOT EXISTS receipts
(
    "id"                 uuid         NOT NULL PRIMARY KEY,
    "receipt_id"         varchar(44)  NOT NULL,
    "cash_register_code" varchar(17)  NOT NULL,
    "ico"                varchar(8)   NOT NULL,
    "ic_dph"             varchar(12),
    "dic"                varchar(10),
    "type"               varchar(2)   NOT NULL,
    "invoice_number"     varchar(50),
    "receipt_number"     varchar(10)  NOT NULL,
    "total_price"        float8       NOT NULL,
    "tax_base_basic"     float8       NOT NULL,
    "tax_base_reduced"   float8       NOT NULL,
    "vat_amount_basic"   float8       NOT NULL,
    "vat_amount_reduced" float8       NOT NULL,
    "vat_rate_basic"     float8       NOT NULL,
    "vat_rate_reduced"   float8       NOT NULL,
    "issue_date"         timestamp    NOT NULL,
    "create_date"        timestamp    NOT NULL,
    "organization"       json         NOT NULL,
    "unit"               json         NOT NULL,
    "items"              json         NOT NULL,
    "file_path"          varchar(255)
);

CREATE UNIQUE INDEX receipts_receipt_id_key ON receipts USING BTREE ("receipt_id");

COMMENT ON TABLE receipts
    IS 'Slovak republic cash register receipts';

COMMENT ON COLUMN receipts.id
    IS 'Internal ID of the receipt';

COMMENT ON COLUMN receipts.receipt_id
    IS 'Receipt ID from Financna sprava';

COMMENT ON COLUMN receipts.cash_register_code
    IS 'Code of the cash register issued to it by tax office';

COMMENT ON COLUMN receipts.ico
    IS 'Identification number of organization';

COMMENT ON COLUMN receipts.ic_dph
    IS 'Identification number for VAT, if the business is VAT payer';

COMMENT ON COLUMN receipts.dic
    IS 'Tax identification number';

COMMENT ON COLUMN receipts.type
    IS 'Type of receipt. Accepted values: PD (pokladničný doklad), UF (úhrada faktúry), ND (neplatný doklad), VY (výber), VK (vklad)';

COMMENT ON COLUMN receipts.invoice_number
    IS 'Number of invoice associated with this receipt';

COMMENT ON COLUMN receipts.file_path
IS 'Location of the receipt file in the filesystem, if it was submitted as a photo.';

CREATE TABLE IF NOT EXISTS user_receipts
(
    user_id       uuid NOT NULL,
    receipt_id_sk uuid NOT NULL,
    PRIMARY KEY ("user_id", "receipt_id_sk"),
    CONSTRAINT "fk_user_receipts_user_id" FOREIGN KEY ("user_id") REFERENCES users ("id") ON UPDATE CASCADE ON DELETE CASCADE,
    CONSTRAINT "fk_user_receipts_receipt_id" FOREIGN KEY ("receipt_id_sk") REFERENCES receipts ("id") ON UPDATE CASCADE ON DELETE CASCADE
);
