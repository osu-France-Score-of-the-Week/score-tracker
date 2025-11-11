FROM postgres:17-alpine

# Copy the init script (shell script allows env var expansion)
COPY build/database/init.sh /docker-entrypoint-initdb.d/
RUN chmod +x /docker-entrypoint-initdb.d/init.sh