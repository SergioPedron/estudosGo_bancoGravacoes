/*
Criação do banco de dados mySql gravacoes para testes de acesso com golang 
 - create database gravacoes;
 - use gravacoes;
 - source fontes/golang/banco_gravacoes/inicializa_banco.sql
*/

DROP TABLE IF EXISTS album;
CREATE TABLE album (
  id         INT AUTO_INCREMENT NOT NULL,
  titulo     VARCHAR(128) NOT NULL,
  artista    VARCHAR(255) NOT NULL,
  preco      DECIMAL(5,2) NOT NULL,
  PRIMARY KEY (`id`)
);

INSERT INTO album
  (titulo, artista, preco)
VALUES
  ('Fique Onde Eu Possa Te Ver', 'Pato Fu', 56.99),
  ('Efeito Dominó', 'Ira', 63.99),
  ('Amor e Fé', 'Pixote', 17.99),
  ('O Passageiro', 'Capital Inicial', 34.98);