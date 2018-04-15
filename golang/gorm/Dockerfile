FROM mysql:5.6

ENV MYSQL_DATABASE=gorm
ENV MYSQL_USER=gorm
ENV MYSQL_PASSWORD=gorm
ENV MYSQL_RANDOM_ROOT_PASSWORD='yes'

EXPOSE 3306

CMD ["mysqld","--character-set-server=utf8", "--collation-server=utf8_unicode_ci"]