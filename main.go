package main

import (
	"flag"
	"fmt"
	"golang.org/x/crypto/ssh"
	"log"
	"net"
	"os"
	"time"
)

// connect 创建连接并获取会话
func connect(user, password, host string, port int) (*ssh.Session, error) {
	var (
		addr         string
		clientConfig *ssh.ClientConfig
		client       *ssh.Client
		session      *ssh.Session
		err          error
	)
	clientConfig = &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		Timeout: 30 * time.Second,
		//HostKeyCallback: ssh.InsecureIgnoreHostKey(), // 验证服务端
		//HostKeyCallback: ssh.FixedHostKey(nil), // 验证服务端
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error { // 验证服务端
			log.Printf("主机名称：%s 地址：%s 公钥：%s", hostname, remote, string(ssh.MarshalAuthorizedKey(key)))
			return nil
		},
	}
	// connect to ssh
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

func init() {
	loadConfig()
}

func main() {
	log.Println("about：https://github.com/bajins/bajins-shell-go")
	var (
		initConf bool
	)
	flag.StringVar(&config.Host, "h", config.Host, "地址")
	flag.IntVar(&config.Port, "p", config.Port, "端口")
	flag.StringVar(&config.Username, "u", config.Username, "用户")
	flag.StringVar(&config.Password, "P", config.Password, "密码")
	flag.StringVar(&config.Cmd, "c", config.Cmd, "命令，多条命令以“;”分割")
	flag.BoolVar(&initConf, "i", false, "初始化配置文件")
	flag.PrintDefaults()
	//flag.Usage() // 打印帮助 -h/-help
	flag.Parse()

	if initConf {
		initConfig()
	}
	// 当参数不存在时，尝试使用配置文件
	if config.Host == "" {
		log.Fatalln("远程地址不能为空")
	}
	if config.Username == "" {
		log.Fatalln("用户名不能为空")
	}
	if config.Cmd == "" {
		log.Fatalln("需要执行的命令不能为空")
	}
	session, err := connect(config.Username, config.Password, config.Host, config.Port)
	if err != nil {
		log.Fatal("远程连接错误：", err)
	}
	defer func(session *ssh.Session) { // 正常退出或异常时，释放资源
		err := session.Close()
		if err != nil {
			return
		}
	}(session)
	session.Stdout = os.Stdout
	session.Stderr = os.Stderr
	session.Stdin = os.Stdin

	err = session.Run(config.Cmd) // 多行命令以分号分隔
	if err != nil {
		log.Println("执行命令错误：", err)
	}

	// 执行交互式命令
	/*fd := int(os.Stdin.Fd())
	  oldState, err := terminal.MakeRaw(fd)
	  if err != nil {
	      panic(err)
	  }
	  defer func(fd int, oldState *terminal.State) { // 正常退出或异常时，释放资源
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
