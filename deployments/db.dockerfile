FROM mysql:5.7.22

COPY . .

EXPOSE 3306

ENV PORT 3306

#RUN mysql -u root -p password < deployments/init.sql

#RUN mysql -u root -e "CREATE DATABASE IF NOT EXISTS entryTask;"

#RUN /bin/bash -c "/usr/bin/mysqld_safe --skip-grant-tables &" && \
#      sleep 5 && \
#      mysql -u root -e "CREATE DATABASE IF NOT EXISTS entryTask"
      # && \
      #mysql -u root mydb < /tmp/dump.sql