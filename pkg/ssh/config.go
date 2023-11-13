package ssh

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"os"
	"time"
)

// Config ssh配置
type Config struct {
	IdentityFile     string        `help:"免密登陆密钥地址" default:"$HOME/.ssh/id_rsa"`
	IdentityPassword string        `help:"免密登陆密钥密码" default:""`
	Timeout          time.Duration `help:"连接超时" default:"30s"`
}

type ServerConfig struct {
	Host           string        `json:"host"`
	User           string        `json:"user"`
	Password       string        `json:"password"` //如果密码为空，则认为是免密登陆
	Port           int           `json:"port"`
	Timeout        time.Duration `json:"timeout"`
	IdentitySigner ssh.Signer    `json:"-"`
}

func (s *ServerConfig) String() (key string) {
	return fmt.Sprintf("%s:%s@%s:%d", s.User, s.Password, s.Host, s.Port)
}

func (conf *Config) IdentitySigner() (signer ssh.Signer, err error) {
	if conf == nil {
		return nil, ErrSSH.New("config can't nil ")
	}
	_, err = os.Stat(conf.IdentityFile)
	if err != nil {
		return nil, ErrSSH.New("ssh config IdentityFile: %s not exists", conf.IdentityFile)
	}
	bytes, err := os.ReadFile(conf.IdentityFile)
	if err != nil {
		return nil, ErrSSH.Wrap(err)
	}
	signer, err = ssh.ParsePrivateKey(bytes)
	if _, ok := err.(*ssh.PassphraseMissingError); ok {
		signer, err = ssh.ParsePrivateKeyWithPassphrase(bytes, []byte(conf.IdentityPassword))
	}
	if err != nil {
		err = ErrSSH.Wrap(err)
	}
	return
}
