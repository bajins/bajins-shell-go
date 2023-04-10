package main

import (
	"flag"
	"fmt"
	"golang.org/x/crypto/ssh"
	"log"
	"os"
	"time"
)

// connect 创建连接并获取会话
func connect(user, password, host string, port int) (*ssh.Session, error) {
	var (
		auth         []ssh.AuthMethod
		addr         string
		clientConfig *ssh.ClientConfig
		client       *ssh.Client
		session      *ssh.Session
		err          error
	)
	// get auth method
	auth = make([]ssh.AuthMethod, 0)
	auth = append(auth, ssh.Password(password))

	clientConfig = &ssh.ClientConfig{
		User:    user,
		Auth:    auth,
		Timeout: 30 * time.Second,
	}
	// connet to ssh
	addr = fmt.Sprintf("%s:%d", host, port)

	if client, err = ssh.Dial("tcp", addr, clientConfig); err != nil {
		return nil, err
	}
	// create session
	if session, err = client.NewSession(); err != nil {
		return nil, err
	}
	return session, nil
}

func main() {
	fmt.Println(`命令执行参数说明：
    -h 地址，默认127.0.0.1
    -p 端口，默认22
    -u 用户名，默认root
    -P 密码
    -c 命令
    -i 初始化配置文件`)
	var (
		host     string
		port     int
		username string
		password string
		cmd      string
		initConf bool
	)
	flag.StringVar(&host, "h", "", "地址")
	flag.IntVar(&port, "p", 0, "端口")
	flag.StringVar(&username, "u", "", "用户")
	flag.StringVar(&password, "P", "", "密码")
	flag.StringVar(&cmd, "c", "", "命令")
	flag.BoolVar(&initConf, "i", false, "初始化")
	flag.Parse()

	if initConf {
		initConfig()
	}
	// 当参数不存在时，尝试使用配置文件
	if host == "" && port == 0 && username == "" && password == "" {
		isSuccess := loadConfig()
		if !isSuccess {
			return
		}
		host = config.Host
		port = config.Port
		username = config.Username
		password = config.Password
		cmd = config.Cmd
	}
	if cmd == "" {
		log.Fatalln("需要执行的命令不能为空")
	}
	session, err := connect(username, password, host, port)
	if err != nil {
		log.Fatal(err)
	}
	defer func(session *ssh.Session) { // 处理异常
		err := session.Close()
		if err != nil {
			log.Println(err)
		}
	}(session)
	session.Stdout = os.Stdout
	session.Stderr = os.Stderr
	session.Stdin = os.Stdin

	err = session.Run(cmd) // 多行命令以分号分隔
	if err != nil {
		log.Fatal(err)
	}

	// 执行交互式命令
	/*fd := int(os.Stdin.Fd())
	  oldState, err := terminal.MakeRaw(fd)
	  if err != nil {
	      panic(err)
	  }
	  defer func(fd int, oldState *terminal.State) { // 处理异常
	      err := terminal.Restore(fd, oldState)
	      if err != nil {
	          log.Println(err)
	      }
	  }(fd, oldState)
	  session.Stdout = os.Stdout
	  session.Stderr = os.Stderr
	  session.Stdin = os.Stdin

	  termWidth, termHeight, err := terminal.GetSize(fd)
	  if err != nil {
	      panic(err)
	  }
	  // Set up terminal modes
	  modes := ssh.TerminalModes{
	      ssh.ECHO:          1,     // enable echoing
	      ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
	      ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	  }
	  // Request pseudo terminal
	  if err := session.RequestPty("xterm-256color", termHeight, termWidth, modes); err != nil {
	      log.Fatal(err)
	  }
	  err = session.Run("top")
	  if err != nil {
	      log.Fatal(err)
	  }*/
}
