#!/bin/bash
set -e

#Installing Postgres
sudo yum install -y postgresql-server postgresql-contrib

#Initialize the Postgres
sudo postgresql-setup --initdb
sudo systemctl enable postgresql
sudo systemctl start postgresql

# Create PostgreSQL user and database
sudo -u postgres psql -c "CREATE USER ${POSTGRES_USER} WITH ENCRYPTED PASSWORD '${POSTGRES_PASSWORD}';"
sudo -u postgres psql -c "GRANT ALL PRIVILEGES ON DATABASE ${POSTGRES_DB} TO ${POSTGRES_USER};"
