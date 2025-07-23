FROM postgres:15

COPY backup.sh /usr/local/bin/backup.sh
COPY ../.env /usr/local/bin/.env
COPY forbidden_tables.txt /usr/local/bin/forbidden_tables.txt
RUN chmod +x /usr/local/bin/backup.sh
