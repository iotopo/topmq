package cmd

import (
	"bufio"
	"fmt"
	"github.com/iotopo/topmq/auth"
	"github.com/iotopo/topmq/cache"
	"github.com/iotopo/topmq/config"
	"github.com/iotopo/topmq/db"
	"github.com/iotopo/topmq/internal/utils"
	"github.com/spf13/cobra"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"os"
	"strings"
)

// adduser -u <username>
var adduserCmd = &cobra.Command{
	Use:   "adduser",
	Short: "添加管理员用户",
	Long: `添加管理员用户命令。

用法:
  adduser -u <用户名> [--admin]

参数:
  -u, --username string   用户名（必需）

示例:
  adduser -u admin                    # 创建 admin 用户

说明:
  - 如果不输入密码，系统将自动生成一个12位随机密码
  - 如果输入密码，需要输入两次进行确认`,
	Run: func(cmd *cobra.Command, args []string) {
		// 解析 username 和 projectID
		username, _ := cmd.Flags().GetString("username")
		if username == "" {
			fmt.Println("用户名是必需的")
			return
		}

		// 从标准输入读取
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("请输入密码（留空将自动生成）: ")
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("读取输入时出错: %v\n", err)
			return
		}
		password := strings.TrimSpace(input)
		// 空密码不再进行二次验证，直接自动生成
		if password != "" {
			fmt.Print("请再次输入密码: ")
			input, err = reader.ReadString('\n')
			if err != nil {
				fmt.Printf("读取输入时出错: %v\n", err)
				return
			}
			password2 := strings.TrimSpace(input)
			if password != password2 {
				fmt.Println("两次输入的密码不匹配")
				return
			}
		}

		adduser(strings.TrimSpace(username), password)
	},
}

var (
	adduserUsername string
)

func init() {
	rootCmd.AddCommand(adduserCmd)

	// 确保标志在命令被使用之前就被注册
	adduserCmd.Flags().StringVarP(&adduserUsername, "username", "u", "", "用户名")
}

func adduser(username, password string) {
	config.Conf.DB.ShowSQL = false
	config.Conf.DB.AutoCreateDatabase = false
	config.Conf.DB.AutoCreateAdmin = false

	cache.Init()
	defer cache.Close()

	db.Init()
	defer db.Close()

	genPassword := false
	if password == "" {
		genPassword = true
		password = utils.RandomString(12)
	}

	passwordHash, passwordSalt := utils.GeneratePasswordHash(password)
	userID := utils.XID()

	err := db.DB.Transaction(func(tx *gorm.DB) error {
		var count int64
		result := tx.Model(auth.Users).
			Where(clause.Eq{Column: "username", Value: username}).
			Count(&count)
		if result.Error != nil {
			return result.Error
		}
		if count > 0 {
			return fmt.Errorf("用户 %s 已存在", username)
		}
		// 创建一个用户，并设为项目管理员
		user := &auth.User{
			ID:           userID,
			Username:     username,
			PasswordHash: passwordHash,
			PasswordSalt: passwordSalt,
			Name:         "administrator",
			IsSuperuser:  true,
		}
		if username == "admin" {
			user.IsSuperuser = true
		}
		return tx.Create(user).Error
	})

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	if genPassword {
		fmt.Println("用户创建成功，密码为：" + password)
	} else {
		fmt.Println("管理员创建成功")
	}
}
