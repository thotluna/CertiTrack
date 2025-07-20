-- backend/db_init/init.sql
-- Este script se ejecuta cuando el contenedor de PostgreSQL se inicia por primera vez.
-- Habilita la extensión uuid-ossp para generar UUIDs.
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";