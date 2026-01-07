package cmd

import (
	"bufio"
	"context"
	"fmt"
	auth2 "github.com/iotopo/topmq/auth"
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

// pwreset -u <username> [-p <project>]
var pwresetCmd = &cobra.Command{
	Use:   "pwreset",
	Short: "重置用户密码",
	Long: `重置用户密码命令。

用法:
  pwreset -u <用户名>

参数:
  -u, --username string   用户名（必需）

示例:
  pwreset -u admin                    # 重置 admin 用户的密码

说明:
  - 如果不输入密码，系统将自动生成一个12位随机密码
  - 如果输入密码，需要输入两次进行确认
  - 重置密码后，该用户的所有会话将被强制下线`,
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

		resetPassword(strings.TrimSpace(username), password)
	},
}

var (
	flagUsername string
)

func init() {
	rootCmd.AddCommand(pwresetCmd)

	// 确保标志在命令被使用之前就被注册
	pwresetCmd.Flags().StringVarP(&flagUsername, "username", "u", "", "用户名")
}

func resetPassword(username, password string) {
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
		password = utils.GeneratePassword(config.Conf.MinPwdLen)
	}
	passwordHash, passwordSalt := utils.GeneratePasswordHash(password)

	user := auth2.User{
		PasswordHash: passwordHash,
		PasswordSalt: passwordSalt,
	}

	err := db.DB.Transaction(func(tx *gorm.DB) error {
		result := tx.Model(auth2.Users).
			Where(clause.Eq{Column: "username", Value: username}).
			Select([]string{"id"}).
			Limit(1).
			Find(&user)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return fmt.Errorf("用户 %s 不存在", username)
		}

		return tx.Select([]string{"password_hash", "password_salt"}).Updates(&user).Error
	})
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// 强制已登录的用户下线还需要清空用户基础信息缓存
	auth2.ClearUserSession(context.Background(), user.ID)

	if genPassword {
		fmt.Printf("密码已重置为: %s\n", password)
	} else {
		fmt.Println("密码修改成功")
	}
}
