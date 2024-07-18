

CREATE TABLE IF NOT EXISTS transactions (
    id UUID DEFAULT generateUUIDv4(),
    sku String,
    amount Int32,
    qty Int32,
    created_by String,
    created_at DateTime DEFAULT now(),
    updated_at DateTime DEFAULT now()
) ENGINE = MergeTree()
ORDER BY created_at;
