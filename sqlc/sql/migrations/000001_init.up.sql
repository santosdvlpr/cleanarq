CREATE TABLE orders (
  id   serial  NOT NULL PRIMARY KEY,
  descricao varchar(60) not null,
  preco decimal(10,2) NOT NULL,
  taxa decimal(10,2) NOT NULL,
  preco_total decimal(10,2) NULL
);
