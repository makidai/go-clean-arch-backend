DROP TABLE IF EXISTS article;
CREATE TABLE article
(
    id INT NOT NULL,
    title varchar(255) NOT NULL,
    content text NOT NULL,
    author_id INT NOT NULL default '0',
    updated_at TIMESTAMP default CURRENT_TIMESTAMP NOT NULL,
    created_at TIMESTAMP default CURRENT_TIMESTAMP NOT NULL,
    PRIMARY KEY (id)
);

DROP TABLE IF EXISTS article_category;
CREATE TABLE article_category 
(
    id INT NOT NULL,
    article_id int NOT NULL,
    category_id int NOT NULL,
    PRIMARY KEY (id),
    CONSTRAINT composite UNIQUE (article_id,category_id)
);

DROP TABLE IF EXISTS author;
CREATE TABLE author
(
    id INT NOT NULL,
    name varchar(255) DEFAULT '""',
    created_at TIMESTAMP default CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP default CURRENT_TIMESTAMP NOT NULL,
    PRIMARY KEY (id)
);

DROP TABLE IF EXISTS category;
CREATE TABLE category (
    id INT NOT NULL,
    name varchar(255) NOT NULL,
    tag varchar(255) NOT NULL,
    created_at TIMESTAMP default CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP default CURRENT_TIMESTAMP NOT NULL,
    PRIMARY KEY (id)
);