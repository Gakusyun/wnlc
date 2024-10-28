package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type Config struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type StatusResponse struct {
	Code    int        `json:"code"`
	Online  OnlineInfo `json:"online"`
	DialMsg string     `json:"dialMsg"`
}

type OnlineInfo struct {
	Name     string `json:"Name"`
	Username string `json:"Username"`
	UserIpv4 string `json:"UserIpv4"`
	UserMac  string `json:"UserMac"`
}

type LoginResponse struct {
	Code int `json:"code"`
}

func clearScreen() {
	cmd := exec.Command("clear")
	if os.Getenv("OS") == "Windows_NT" {
		cmd = exec.Command("cmd", "/c", "cls")
	}
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func getStatus() (*StatusResponse, error) {
	url := "http://192.168.96.110/api/account/status"
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result StatusResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	return &result, err
}

func login(username, password string) (*LoginResponse, error) {
	url := "http://192.168.96.110/api/account/login"
	data := map[string]interface{}{
		"username": username,
		"password": password,
		"nasld":    "2",
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result LoginResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	return &result, err
}

func logout() error {
	url := "http://192.168.96.110/api/account/logout"
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

func configDir() string {
	// 获取用户主目录
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Error getting user home directory:", err)
		return ""
	}

	// 构建 .wnlc 目录的路径
	wnlcDir := homeDir + "/.wnlc"
	return wnlcDir
}

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

		if _, err := strconv.Atoi(username); err == nil {
			if _, err := strconv.Atoi(password); err == nil {
				break
			}
		}
		fmt.Println("账号密码必须是数字，请重新输入")
	}

	config := Config{
		Username: username,
		Password: password,
	}

	jsonData, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}
	wnlcDir := configDir()
	err = os.MkdirAll(wnlcDir, 0755)
	if err != nil {
		return err
	}
	err = os.WriteFile(wnlcDir+"/config.json", jsonData, 0644)
	return err
}

func readConfig() (*Config, error) {
	wnlcDir := configDir()
	data, err := os.ReadFile(wnlcDir + "/config.json")
	if err != nil {
		return nil, err
	}

	var config Config
	err = json.Unmarshal(data, &config)
	return &config, err
}

func input(prompt string) string {
	fmt.Print(prompt)
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

func main() {
	statusResult, err := getStatus()
	if err != nil {
		fmt.Println("获取状态失败:", err)
		return
	}

	for {
		config, err := readConfig()
		if err != nil {
			fmt.Println("未找到配置文件，请先创建")
			if err := createConfig(); err != nil {
				fmt.Println("创建配置文件失败:", err)
				return
			}
			continue
		}

		username := config.Username
		password := config.Password

		if statusResult.Code == 1 {
			loginResult, err := login(username, password)
			if err != nil {
				fmt.Println("登录失败:", err)
				return
			}

			if loginResult.Code == 1 {
				fmt.Println("账号或密码错误")
				if input("输入0修改密码\n") == "0" {
					if err := createConfig(); err != nil {
						fmt.Println("创建配置文件失败:", err)
						return
					}
					continue
				}
			} else if loginResult.Code == 2 {
				fmt.Println("需要填写验证码，请前往 http://192.168.96.110/ 登录")
				break
			}

			statusResult, err = getStatus()
			if err != nil {
				fmt.Println("获取状态失败:", err)
				return
			}
			continue
		} else {
			clearScreen()
			online := statusResult.Online
			fmt.Printf("你好, %s\n", online.Name)
			fmt.Printf("当前账号: %s\n", online.Username)
			fmt.Printf("当前IP: %s\n", online.UserIpv4)
			fmt.Printf("当前MAC: %s\n", online.UserMac)
			fmt.Printf("%s\n", statusResult.DialMsg)
			fmt.Println("=============================================================")

			flag := input("0.退出登录\n1.修改账号\n3.关闭软件但不断开连接\n")
			switch flag {
			case "0":
				if err := logout(); err != nil {
					fmt.Println("注销失败:", err)
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
				continue
			case "3":
				return
			default:
				fmt.Println("无效的选项，请重试")
			}
		}
	}
}
