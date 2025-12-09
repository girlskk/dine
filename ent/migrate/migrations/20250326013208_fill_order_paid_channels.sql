-- Fill "paid_channels" column in "orders" table
UPDATE orders
SET paid_channels = JSON_MERGE_PRESERVE(
    '[]',
    CASE WHEN cash_paid > 0 THEN JSON_ARRAY('cash') ELSE JSON_ARRAY() END,
    CASE WHEN wechat_paid > 0 THEN JSON_ARRAY('wechat') ELSE JSON_ARRAY() END,
    CASE WHEN alipay_paid > 0 THEN JSON_ARRAY('alipay') ELSE JSON_ARRAY() END,
    CASE WHEN points_paid > 0 THEN JSON_ARRAY('point') ELSE JSON_ARRAY() END
)
WHERE cash_paid > 0 OR wechat_paid > 0 OR alipay_paid > 0 OR points_paid > 0;
