echo "enter root password"
echo "CREATE DATABASE pajbot2 CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;GRANT ALL PRIVILEGES ON pajbot2.* TO 'pajbot2'@'localhost' IDENTIFIED BY 'password';USE pajbot2;CREATE TABLE pb_command (triggers VARCHAR(512));" | mysql -uroot -p
