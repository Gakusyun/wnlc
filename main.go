package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

// 常量定义
const (
	StatusLoggedIn    = 0
	StatusNotLoggedIn = 1

	LoginSuccess = 0
	LoginInvalid = 1
	LoginCaptcha = 2
)

// 配置结构体
type Config struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// 状态响应结构体
type StatusResponse struct {
	Code    int        `json:"code"`
	Online  OnlineInfo `json:"online"`
	DialMsg string     `json:"dialMsg"`
}

// 在线信息结构体
type OnlineInfo struct {
	Name     string `json:"Name"`
	Username string `json:"Username"`
	UserIpv4 string `json:"UserIpv4"`
	UserMac  string `json:"UserMac"`
}

// 登录响应结构体
type LoginResponse struct {
	Code int `json:"code"`
}

// 清屏函数
func clearScreen() {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "cls")
	} else {
		cmd = exec.Command("clear")
	}
	cmd.Stdout = os.Stdout
	_ = cmd.Run()
}

// 获取状态
func getStatus() (*StatusResponse, error) {
	url := "http://192.168.96.110/api/account/status"
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("无法连接到服务器: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("获取状态失败，服务器返回: %s", resp.Status)
	}

	var result StatusResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, fmt.Errorf("解析状态数据失败: %v", err)
	}
	return &result, nil
}

// 登录
func login(username, password string) (*LoginResponse, error) {
	url := "http://192.168.96.110/api/account/login"
	data := map[string]interface{}{
		"username": username,
		"password": password,
		"nasld":    "2",
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("序列化登录数据失败: %v", err)
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("无法连接到服务器: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("登录失败，服务器返回: %s", resp.Status)
	}

	var result LoginResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, fmt.Errorf("解析登录响应失败: %v", err)
	}
	return &result, nil
}

// 注销
func logout() error {
	url := "http://192.168.96.110/api/account/logout"
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("无法连接到服务器: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("注销失败，服务器返回: %s", resp.Status)
	}
	return nil
}

// 获取配置目录
func configDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("获取用户主目录失败: %v", err)
	}
	wnlcDir := filepath.Join(homeDir, ".wnlc")
	return wnlcDir, nil
}

// 创建配置文件
func createConfig() error {
	reader := bufio.NewReader(os.Stdin)
	var username, password string

	for {
		fmt.Print("请输入你的学号: ")
		username, _ = reader.ReadString('\n')
		username = strings.TrimSpace(username)

		fmt.Print("请输入身份证后8位: ")
		password, _ = reader.ReadString('\n')
		password = strings.TrimSpace(password)

		if _, err := strconv.Atoi(username); err == nil && len(username) > 0 {
			if _, err := strconv.Atoi(password); err == nil && len(password) == 8 {
				break
			}
		}
		fmt.Println("账号或密码格式错误，请重新输入")
	}

	config := Config{
		Username: username,
		Password: password,
	}

	jsonData, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("序列化配置数据失败: %v", err)
	}

	wnlcDir, err := configDir()
	if err != nil {
		return err
	}

	err = os.MkdirAll(wnlcDir, 0755)
	if err != nil {
		return fmt.Errorf("创建配置目录失败: %v", err)
	}

	err = os.WriteFile(filepath.Join(wnlcDir, "config.json"), jsonData, 0644)
	if err != nil {
		return fmt.Errorf("写入配置文件失败: %v", err)
	}
	return nil
}

// 读取配置文件
func readConfig() (*Config, error) {
	wnlcDir, err := configDir()
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(filepath.Join(wnlcDir, "config.json"))
	if err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %v", err)
	}

	var config Config
	err = json.Unmarshal(data, &config)
	if err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %v", err)
	}
	return &config, nil
}

// 获取用户输入
func input(prompt string) string {
	fmt.Print(prompt)
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

// 显示帮助信息
func showHelp() {
	fmt.Println("Usage: wnlc [OPTION]")
	fmt.Println("WUIT Network Login client written in Go.")
	fmt.Println("")
	fmt.Println("Options:")
	fmt.Println("-l\tLogin to the server")
	fmt.Println("-s\tGet status of the server")
	fmt.Println("-c\tCreate config file")
	fmt.Println("-h\tShow help")
}

// 处理登录
func handleLogin() {
	config, err := readConfig()
	if err != nil {
		fmt.Println("未找到配置文件，请先创建")
		if err := createConfig(); err != nil {
			fmt.Println("创建配置文件失败:", err)
			return
		}
		config, err = readConfig()
		if err != nil {
			fmt.Println("读取配置文件失败:", err)
			return
		}
	}

	username := config.Username
	password := config.Password

	loginResult, err := login(username, password)
	if err != nil {
		fmt.Println("登录失败:", err)
		return
	}

	switch loginResult.Code {
	case LoginSuccess:
		fmt.Println("登录成功")
	case LoginInvalid:
		fmt.Println("账号或密码错误")
		if input("输入0修改密码: ") == "0" {
			if err := createConfig(); err != nil {
				fmt.Println("创建配置文件失败:", err)
				return
			}
		}
	case LoginCaptcha:
		fmt.Println("需要填写验证码，请前往 http://192.168.96.110/ 登录")
	default:
		fmt.Println("未知的登录错误")
	}
}

// 显示状态
func handleStatus() {
	statusResult, err := getStatus()
	if err != nil {
		fmt.Println("获取状态失败:", err)
		return
	}

	if statusResult.Code == StatusLoggedIn {
		online := statusResult.Online
		fmt.Printf("你好, %s\n", online.Name)
		fmt.Printf("当前账号: %s\n", online.Username)
		fmt.Printf("当前IP: %s\n", online.UserIpv4)
		fmt.Printf("当前MAC: %s\n", online.UserMac)
		fmt.Printf("%s\n", statusResult.DialMsg)
	} else {
		fmt.Println("当前未登录")
	}
}

// 交互模式
func handleInteractiveMode() {
	for {
		statusResult, err := getStatus()
		if err != nil {
			fmt.Println("获取状态失败:", err)
			return
		}

		if statusResult.Code == StatusNotLoggedIn {
			config, err := readConfig()
			if err != nil {
				fmt.Println("未找到配置文件，请先创建")
				if err := createConfig(); err != nil {
					fmt.Println("创建配置文件失败:", err)
					return
				}
				config, err = readConfig()
				if err != nil {
					fmt.Println("读取配置文件失败:", err)
					return
				}
			}

			username := config.Username
			password := config.Password

			loginResult, err := login(username, password)
			if err != nil {
				fmt.Println("登录失败:", err)
				return
			}

			switch loginResult.Code {
			case LoginSuccess:
				fmt.Println("登录成功")
			case LoginInvalid:
				fmt.Println("账号或密码错误")
				if input("输入0修改密码: ") == "0" {
					if err := createConfig(); err != nil {
						fmt.Println("创建配置文件失败:", err)
						return
					}
				}
			case LoginCaptcha:
				fmt.Println("需要填写验证码，请前往 http://192.168.96.110/ 登录")
			default:
				fmt.Println("未知的登录错误")
			}
		} else {
			clearScreen()
			online := statusResult.Online
			fmt.Printf("你好, %s\n", online.Name)
			fmt.Printf("当前账号: %s\n", online.Username)
			fmt.Printf("当前IP: %s\n", online.UserIpv4)
			fmt.Printf("当前MAC: %s\n", online.UserMac)
			fmt.Printf("%s\n", statusResult.DialMsg)
			fmt.Println("=============================================================")

			choice := input("0.退出登录\n1.修改账号\n3.关闭软件但不断开连接\n请选择: ")
			switch choice {
			case "0":
				if err := logout(); err != nil {
					fmt.Println("注销失败:", err)
				} else {
					fmt.Println("注销成功")
				}
				return
			case "1":
				if err := createConfig(); err != nil {
					fmt.Println("创建配置文件失败:", err)
					return
				}
				if err := logout(); err != nil {
					fmt.Println("注销失败:", err)
				}
			case "3":
				return
			default:
				fmt.Println("无效的选项，请重试")
			}
		}
	}
}

// 主函数
func main() {
	if len(os.Args) < 2 {
		handleInteractiveMode()
	} else {
		switch os.Args[1] {
		case "-l":
			handleLogin()
		case "-s":
			handleStatus()
		case "-c":
			if err := createConfig(); err != nil {
				fmt.Println("创建配置文件失败:", err)
			} else {
				fmt.Println("配置文件创建成功")
			}
		case "-h":
			showHelp()
		default:
			fmt.Println("无效的选项，请使用 -h 查看帮助")
		}
	}
}
