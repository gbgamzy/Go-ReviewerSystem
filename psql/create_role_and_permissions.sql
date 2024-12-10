-- SQL script to create the client role and grant permissions
CREATE ROLE client WITH LOGIN PASSWORD 'password';

-- Grant client role access to the ReviewerDB database
GRANT CONNECT ON DATABASE "ReviewerDB" TO client;

-- Grant schema and table-level privileges
\c ReviewerDB
GRANT USAGE ON SCHEMA public TO client;
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO client;
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO client;

-- Set default privileges for future tables
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON TABLES TO client;
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON SEQUENCES TO client;
