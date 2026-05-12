#!/bin/bash
# 执行升级的脚本
# 2023.09.15    wuhaifeng

# todo  config.json 待替换

src_dir=/tmp/94eb5a1cbed058a9
data_dir=/opt/laozhi/data
project_dir=/opt/laozhi

# 第一步 检查运行环境
uid=`id -u`
echo $uid
if [ "$uid" -ne 0 ]; then
    echo 请使用root用户权限执行升级
    rm -rf $src_dir
    exit
fi

# 第二步 获取版本
version=`docker exec laozhi_db mysql -uxiaozhi -pxiaozhi.4dogs.cn smart --default-character-set=utf8 -e "select obj_value from map_set where obj_key='systemVersion'"|grep currentVersion|jq '.currentVersion'`
version=${version#"\""}
version=${version%"\""}
echo '当前系统版本是：'$version
backup_dir=/opt/laozhi/backup/$version

# 第三步 进行备份
rm -rf $backup_dir
mkdir -p $backup_dir
cp $project_dir/decision/decision $backup_dir                 # 备份decision
cp $project_dir/smart/smart $backup_dir                       # 备份smart
cp -r $project_dir/nginx/smart_vue $backup_dir                # 备份smart_vue
cd $project_dir                                               # 备份mysql数据
zip -r data.zip data
mv data.zip $backup_dir

# 第四步 文件替换
cd $project_dir
mv $src_dir/decision $project_dir/decision
mv $src_dir/smart $project_dir/smart
rm -rf $project_dir/nginx/smart_vue
mv $src_dir/smart_vue $project_dir/nginx

# 第五步 数据结构升级
mkdir -p $project_dir/data/sql
cp $src_dir/smart_struct.sql $project_dir/data/sql
docker exec laozhi_db mysql -uxiaozhi -pxiaozhi.4dogs.cn smart --default-character-set=utf8 -e "source /var/lib/mysql/sql/smart_struct.sql"

cp $src_dir/decision_struct.sql $project_dir/data/sql
docker exec laozhi_db mysql -uxiaozhi -pxiaozhi.4dogs.cn decision --default-character-set=utf8 -e "source /var/lib/mysql/sql/decision_struct.sql"

# 第六步  数据升级
cp $src_dir/smart_data.sql $project_dir/data/sql
docker exec laozhi_db mysql -uxiaozhi -pxiaozhi.4dogs.cn smart --default-character-set=utf8 -e "source /var/lib/mysql/sql/smart_data.sql"

cp $src_dir/decision_data.sql $project_dir/data/sql
docker exec laozhi_db mysql -uxiaozhi -pxiaozhi.4dogs.cn decision --default-character-set=utf8 -e "source /var/lib/mysql/sql/decision_data.sql"

# 第七步  替换nginx和docker-compose配置文件
cp $src_dir/nginx.conf $project_dir/nginx/

cp $src_dir/docker-compose.yml $project_dir/

cd $project_dir
docker-compose stop nginx
docker-compose up -d

# 第八步  服务重启
supervisord ctl restart smart decision

# 第九步  版本更新
newVersion=`cat $src_dir/version.txt`
updateSql=`docker exec laozhi_db mysql -uxiaozhi -pxiaozhi.4dogs.cn smart --default-character-set=utf8 -e "select obj_value from map_set where obj_key='systemVersion'"|grep currentVersion|jq --arg new_value "$newVersion" '.currentVersion = $new_value'`

# 第十步  删除掉临时解压文件
rm -rf $src_dir