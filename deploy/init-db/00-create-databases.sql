SELECT 'CREATE DATABASE auth'
WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'auth')\gexec
SELECT 'CREATE DATABASE gateway'
WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'gateway')\gexec
SELECT 'CREATE DATABASE requests'
WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'requests')\gexec
