package cmd

import (
	"bufio"
	"fmt"
	"github.com/iotopo/topmq/config/secrets"
	"github.com/iotopo/topmq/internal/utils"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

var secretCmd = &cobra.Command{
	Use:   "secret [subcommand] [flags]",
	Short: "加密或解密文本",
}

// powershell，$ 会被当作特殊字符处理。 输入 secret encode "Tq7@nK2$wX9333" 时，$ 字符在命令行中被解释为变量引用。
// 需要使用 ^ 来转义 $ 字符：secret encode "Tq7@nK2^$wX9333"
// 推荐使用 secret encode --stdin 来输入文本

var encodeCmd = &cobra.Command{
	Use:   "encode [text]",
	Short: "加密给定文本",
	Long: `加密输入文本。

用法:
  secret encode <文本>              # 直接加密命令行参数中的文本
  secret encode --stdin             # 从标准输入读取文本进行加密

参数:
  text string    要加密的文本（可选，如果不提供则使用 --stdin 标志）

标志:
  --stdin        从标准输入读取文本

示例:
  secret encode "hello world"       # 加密 "hello world"
  secret encode --stdin             # 交互式输入文本进行加密

注意:
  - 在 PowerShell 中，$ 字符需要转义，建议使用 --stdin 标志`,
	Args: cobra.MaximumNArgs(1), // 允许0或1个参数
	Run: func(cmd *cobra.Command, args []string) {
		var text string

		// 检查是否使用 --stdin 标志
		if stdinFlag, _ := cmd.Flags().GetBool("stdin"); stdinFlag {
			// 从标准输入读取
			reader := bufio.NewReader(os.Stdin)
			fmt.Print("请输入要加密的文本: ")
			input, err := reader.ReadString('\n')
			if err != nil {
				fmt.Printf("读取输入时出错: %v\n", err)
				return
			}
			text = strings.TrimSpace(input)
		} else if len(args) > 0 {
			// 从命令行参数读取
			text = args[0]
		} else {
			fmt.Println("错误: 请提供要加密的文本或使用 --stdin 标志")
			cmd.Help()
			return
		}

		// fmt.Printf("encode %s\n", text)
		encoded, err := secrets.Encode(text)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("%s\n", encoded)
	},
}

var decodeCmd = &cobra.Command{
	Use:   "decode [text]",
	Short: "解密给定文本",
	Long: `解密输入文本。

用法:
  secret decode <文本>              # 直接解密命令行参数中的文本
  secret decode --stdin             # 从标准输入读取文本进行解密

参数:
  text string    要解密的文本（可选，如果不提供则使用 --stdin 标志）

标志:
  --stdin        从标准输入读取文本

示例:
  secret decode "encrypted_text"    # 解密 "encrypted_text"
  secret decode --stdin             # 交互式输入文本进行解密`,
	Args: cobra.MaximumNArgs(1), // 允许0或1个参数
	Run: func(cmd *cobra.Command, args []string) {
		var text string

		// 检查是否使用 --stdin 标志
		if stdinFlag, _ := cmd.Flags().GetBool("stdin"); stdinFlag {
			// 从标准输入读取
			reader := bufio.NewReader(os.Stdin)
			fmt.Print("请输入要解密的文本: ")
			input, err := reader.ReadString('\n')
			if err != nil {
				fmt.Printf("读取输入时出错: %v\n", err)
				return
			}
			text = strings.TrimSpace(input)
		} else if len(args) > 0 {
			// 从命令行参数读取
			text = args[0]
		} else {
			fmt.Println("错误: 请提供要解密的文本或使用 --stdin 标志")
			cmd.Help()
			return
		}

		// fmt.Printf("decode %s\n", text)
		decoded, err := secrets.Decode(text)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("%s\n", decoded)
	},
}

// getsecret -l <len>
var genSecretCmd = &cobra.Command{
	Use:   "gen",
	Short: "生成密钥",
	Long: `生成指定长度的密钥。

用法:
  secret gen -l <密钥长度>

参数:
  -l, --len number   密钥长度

示例:
  secret gen -l 8                    # 生成一个 8 位的密码`,
	Run: func(cmd *cobra.Command, args []string) {
		// 解析 username 和 projectID
		l, _ := cmd.Flags().GetInt("len")

		if l == 0 {
			l = 8
		}
		s := utils.GeneratePassword(l)
		fmt.Println(s)
	},
}

func init() {
	secretCmd.AddCommand(encodeCmd)
	secretCmd.AddCommand(decodeCmd)
	secretCmd.AddCommand(genSecretCmd)
	rootCmd.AddCommand(secretCmd)

	// 为 encodeCmd 和 decodeCmd 添加 --stdin 标志
	encodeCmd.Flags().BoolP("stdin", "s", false, "从标准输入读取文本")
	decodeCmd.Flags().BoolP("stdin", "s", false, "从标准输入读取文本")
	genSecretCmd.Flags().IntP("len", "l", 8, "密钥长度")
}
