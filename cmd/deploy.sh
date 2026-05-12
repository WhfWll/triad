src_dir=/opt/upgrade/deploy/
data_dir=/opt/laozhi/data
project_dir=/opt/laozhi

# 添加用户
docker exec laozhi_db mysql -uroot -proot.4dogs.cn mysql --default-character-set=utf8 -e "CREATE USER 'xiaozhi'@'localhost' IDENTIFIED BY 'xiaozhi.4dogs.cn';"

docker exec laozhi_db mysql -uroot -proot.4dogs.cn mysql --default-character-set=utf8 -e "GRANT ALL PRIVILEGES ON * . * TO 'xiaozhi'@'localhost';"


docker exec laozhi_db mysql -uroot -proot.4dogs.cn mysql --default-character-set=utf8 -e "CREATE USER 'xiaozhi'@'%' IDENTIFIED BY 'xiaozhi.4dogs.cn';"

docker exec laozhi_db mysql -uroot -proot.4dogs.cn mysql --default-character-set=utf8 -e "GRANT ALL PRIVILEGES ON * . * TO 'xiaozhi'@'%';"

docker exec laozhi_db mysql -uroot -proot.4dogs.cn mysql --default-character-set=utf8 -e "FLUSH PRIVILEGES;"

# 创建数据库
docker exec laozhi_db mysql -uxiaozhi -pxiaozhi.4dogs.cn mysql --default-character-set=utf8 -e "CREATE DATABASE smart CHARACTER SET 'utf8' COLLATE 'utf8_general_ci'"

docker exec laozhi_db mysql -uxiaozhi -pxiaozhi.4dogs.cn mysql --default-character-set=utf8 -e "CREATE DATABASE decision CHARACTER SET 'utf8' COLLATE 'utf8_general_ci'"


# 初始化数据结构
mkdir -p $project_dir/data/sql/

cp $src_dir/smart_struct.sql $project_dir/data/sql
docker exec laozhi_db mysql -uxiaozhi -pxiaozhi.4dogs.cn smart --default-character-set=utf8 -e "source /var/lib/mysql/sql/smart_struct.sql"

cp $src_dir/decision_struct.sql $project_dir/data/sql
docker exec laozhi_db mysql -uxiaozhi -pxiaozhi.4dogs.cn decision --default-character-set=utf8 -e "source /var/lib/mysql/sql/decision_struct.sql"

# 第六步  数据升级
cp $src_dir/smart_data.sql $project_dir/data/sql
docker exec laozhi_db mysql -uxiaozhi -pxiaozhi.4dogs.cn smart --default-character-set=utf8 -e "source /var/lib/mysql/sql/smart_data.sql"

cp $src_dir/decision_data.sql $project_dir/data/sql
docker exec laozhi_db mysql -uxiaozhi -pxiaozhi.4dogs.cn decision --default-character-set=utf8 -e "source /var/lib/mysql/sql/decision_data.sql"
