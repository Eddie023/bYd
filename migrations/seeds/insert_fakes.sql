-- Inserting data into the state table
INSERT INTO state (id, name, code) VALUES
(1, 'New South Wales', 'NSW'),
(2, 'Victoria', 'VIC'),
(3, 'Queensland', 'QLD'),
(4, 'South Australia', 'SA'),
(5, 'Western Australia', 'WA'),
(6, 'Tasmania', 'TAS'),
(7, 'Northern Territory', 'NT'),
(8, 'Australian Capital Territory', 'ACT');

-- Inserting data into the post_type table
INSERT INTO post_type (id, name) VALUES
(1, 'Blog Post'),
(2, 'News Article'),
(3, 'Discussion'),
(4, 'Question'),
(5, 'Review');

-- Inserting data into the users table
INSERT INTO users (id, email, first_name, last_name, state_id, created_at, updated_at, is_verified) VALUES
('user_001', 'user_01@gmail.com','John', 'Doe', 1, NOW(), NOW(), TRUE),
('user_002', 'user_02@gmail.com','Jane', 'Smith', 2, NOW(), NOW(), FALSE),
('user_003', 'user_03@gmail.com','Alice', 'Johnson', 3, NOW(), NOW(), TRUE),
('user_004', 'user_04@gmail.com','Bob', 'Brown', 4, NOW(), NOW(), FALSE),
('user_005', 'user_05@gmail.com','Charlie', 'Davis', 5, NOW(), NOW(), TRUE);

-- Inserting data into the post table
INSERT INTO post (user_id, title, is_anon, description, type_id, created_at, updated_at) VALUES
('user_001', 'My First Blog', FALSE, 'This is the description of my first blog post.', 1, NOW(), NOW()),
('user_002', 'Breaking News', FALSE, 'This is a breaking news article.', 2, NOW(), NOW()),
('user_003', 'Golang Discussion', TRUE, 'Discussion about Golang best practices.', 3, NOW(), NOW()),
('user_004', 'Question on SQL', TRUE, 'How to optimize SQL queries?', 4, NOW(), NOW()),
('user_005', 'Product Review', FALSE, 'This is a review of a new product.', 5, NOW(), NOW());

-- Inserting data into the comment table
INSERT INTO comment (post_id, user_id, description, created_at, updated_at) VALUES
(1, 'user_002', 'Great post!', NOW(), NOW()),
(2, 'user_001', 'Interesting news.', NOW(), NOW()),
(3, 'user_004', 'I agree with your points.', NOW(), NOW()),
(4, 'user_003', 'Thanks for the insights.', NOW(), NOW()),
(5, 'user_005', 'Very detailed review.', NOW(), NOW());
