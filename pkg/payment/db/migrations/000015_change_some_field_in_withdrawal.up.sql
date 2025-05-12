ALTER TABLE withdrawal RENAME COLUMN amount to gross_amount;
ALTER TABLE withdrawal ADD COLUMN net_amount NUMERIC(20, 2) NULL DEFAULT 0;
ALTER TABLE withdrawal ADD COLUMN charge_amount NUMERIC(20, 2) NULL DEFAULT 0;