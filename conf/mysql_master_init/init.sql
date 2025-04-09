create user 'slave'@'%' identified with mysql_native_password by 'miyike3716';
grant replication slave, replication client on *.* to 'slave'@'%';
grant select on *.* to 'slave'@'%';
flush privileges;