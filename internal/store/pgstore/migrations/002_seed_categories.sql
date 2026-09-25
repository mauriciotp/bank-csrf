INSERT INTO categories (category_name) VALUES
('standard'),
('premium'),
('private'),
('micro_business'),
('small_business'),
('corporate');

---- create above / drop below ----

DELETE FROM categories
WHERE category_name IN (
  'standard', 'premium', 'private',
  'micro_business', 'small_business', 'corporate'
);
