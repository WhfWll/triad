//go:build ignore

package main

import (
	"encoding/json"
	"fmt"
	"github.com/shirou/gopsutil/disk"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"
)

// 中广核 辅助测试脚本

/*
备份命令
cd /opt/laozhi
sudo docker-compose stop db
sudo cp -r data data240816
sudo docker-compose start db
sudo supervisorctl restart all
恢复命令
cd /opt/laozhi
sudo docker-compose stop db
sudo mv data /tmp/
sudo cp -r data240816 data
sudo docker-compose start db
sudo supervisorctl restart all
*/

func main() {
	fmt.Println("start check ")
	filePath := "record.txt"
	// 尝试打开文件
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		fmt.Printf("The file %s does not exist.\n", filePath)
	} else if err != nil {
		fmt.Printf("An error occurred while checking the file: %v \n", err)
	} else {
		fmt.Printf("The file %s exists.\n", filePath)
	}
	// 打开或创建一个名为 "example.txt" 的文件
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	// 检查磁盘空间
	res0 := checkDisk()
	_, err = file.WriteString("disk: ")
	_, err = file.WriteString(res0)
	_, err = file.WriteString("\n")

	// 执行date命令
	res1 := execCommandDate()
	_, err = file.WriteString("date: ")
	_, err = file.Write(res1)
	_, err = file.WriteString("\n")

	// 执行timedatectl命令
	res2 := execCommandTimedatectl()
	_, err = file.WriteString("timedatectl: ")
	_, err = file.Write(res2)
	_, err = file.WriteString("\n")

	// 执行备份命令
	res3 := execCommandBak()
	_, err = file.WriteString("bak: ")
	_, err = file.Write(res3)
	_, err = file.WriteString("\n")

	// 连接数据库
	db, err := dbConnect()
	if err != nil {
		_, err = file.WriteString("get db error: " + err.Error())
		_, err = file.WriteString("\n")
		os.Exit(1)
	}
	// 查询正在运行的任务
	res4, err := dbSelectTaskRunning(db)
	if err != nil {
		_, err = file.WriteString("get task error: " + err.Error())
		_, err = file.WriteString("\n")
		os.Exit(1)
	}
	_, err = file.WriteString("running task list: ")
	_, err = file.Write(res4)
	_, err = file.WriteString("\n")
	// 结束正在运行的任务
	err = dbFinishTaskRunning(db)
	if err != nil {
		_, err = file.WriteString("finish task error: " + err.Error())
		_, err = file.WriteString("\n")
	}
	// 查询正在运行的任务
	res5, err := dbSelectTaskRunning(db)
	if err != nil {
		_, err = file.WriteString("get task error: " + err.Error())
		_, err = file.WriteString("\n")
		os.Exit(1)
	}
	_, err = file.WriteString("running task list: ")
	_, err = file.Write(res5)
	_, err = file.WriteString("\n")
	// 记下当前时间
	fmt.Println(time.Now())
	_, err = file.WriteString(time.Now().String())
	_, err = file.WriteString("\n")

	_, err = file.WriteString("check over")
	_, err = file.WriteString("\n")
	fmt.Println("check over")
}

func checkDisk() string {
	// 指定要查询的磁盘路径
	path := "/"
	// 获取磁盘信息
	usage, err := disk.Usage(path)
	if err != nil {
		fmt.Println("Failed to get disk usage: " + err.Error())
	}
	if usage.Free < 5 {
		fmt.Println("Free disk space is less than 5GB")
		os.Exit(1)
	}
	tempString := "Total disk space: " + strconv.Itoa(int(usage.Total))
	tempString += " Used disk space: " + strconv.Itoa(int(usage.Used))
	tempString += " Free disk space: " + strconv.Itoa(int(usage.Free))
	return tempString
}

func execCommandDate() []byte {
	cmd := exec.Command("date")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return []byte("Failed to execute command: " + err.Error())
	}
	return output
}

func execCommandTimedatectl() []byte {
	cmd := exec.Command("timedatectl")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return []byte("Failed to execute command: " + err.Error())
	}
	return output
}

func execCommandBak() []byte {
	srcDir := "/opt/laozhi/data"
	dstDir := "./data20240816"
	if err := copyDir(srcDir, dstDir); err != nil {
		return []byte("Failed to copy directory: " + err.Error())
	}
	return []byte("file backup success")
}

func zghCopyFile(src string, dst string) error {
	s, err := os.Open(src)
	if err != nil {
		return err
	}
	defer s.Close()
	d, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer d.Close()
	// 复制文件内容
	_, err = io.Copy(d, s)
	if err != nil {
		return err
	}
	// 设置目标文件的权限
	info, err := os.Stat(src)
	if err != nil {
		return err
	}
	err = os.Chmod(dst, info.Mode())
	if err != nil {
		return err
	}
	return nil
}

func copyDir(src string, dst string) error {
	// 获取源目录的信息
	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}
	// 检查源目录是否确实是一个目录
	if !srcInfo.IsDir() {
		return fmt.Errorf("%s is not a directory", src)
	}
	// 创建目标目录
	err = os.MkdirAll(dst, srcInfo.Mode())
	if err != nil {
		return err
	}
	// 遍历源目录中的所有文件和子目录
	err = filepath.WalkDir(src, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		// 构建目标文件或目录的路径
		targetPath := filepath.Join(dst, path[len(src):])
		// 如果是文件，则复制文件
		if !d.IsDir() {
			return zghCopyFile(path, targetPath)
		}
		// 如果是目录，则创建目录
		return os.MkdirAll(targetPath, d.Type().Perm())
	})
	return err
}

func dbConnect() (*gorm.DB, error) {
	dsn := "xiaozhi:xiaozhi.4dogs.cn@tcp(127.0.0.1:33306)/qiming_db2?charset=utf8mb4&parseTime=True&loc=Local"
	// 连接数据库
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	// 从 GORM 的 DB 对象中获取底层的 sql.DB 对象
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	// 测试连接
	err = sqlDB.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}

type TaskCheckTask struct {
	TaskName     string
	TaskType     int
	TargetNumber int
	RiskLevel    int
	TaskStatus   int
	UpdateTime   time.Time
	CreateTime   time.Time
}

func (t *TaskCheckTask) TableName() string {
	return "task_checktask"
}

func dbSelectTaskRunning(db *gorm.DB) ([]byte, error) {
	var taskList []TaskCheckTask
	err := db.Debug().Where("task_status = ?", 3).Find(&taskList).Error
	if err != nil {
		return nil, err
	}
	taskListByte, err := json.Marshal(taskList)
	if err != nil {
		return taskListByte, err
	}
	return taskListByte, nil
}

func dbFinishTaskRunning(db *gorm.DB) error {
	var task TaskCheckTask
	err := db.Model(&task).Debug().Where("task_status = ?", 3).Update("task_status", 4).Error
	if err != nil {
		return err
	}
	return nil
}
