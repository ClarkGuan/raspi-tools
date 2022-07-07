package rtools

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/bramvdbogaerde/go-scp"
	"golang.org/x/crypto/ssh"
)

var (
	IP               = "192.168.1.8"
	RootUser         = "pi"
	RootUserPassword = "raspberry"
)

const (
	dirDefault = "/home/%s/"

	cvIPEnvName       = "RASPI_IP"
	cvUserEnvName     = "RASPI_USER"
	cvPasswordEnvName = "RASPI_PASSWORD"
)

func init() {
	prepareConfig()
}

func prepareConfig() {
	if newIP, found := os.LookupEnv(cvIPEnvName); found {
		IP = newIP
	}
	if newUser, found := os.LookupEnv(cvUserEnvName); found {
		RootUser = newUser
	}
	if newPassword, found := os.LookupEnv(cvPasswordEnvName); found {
		RootUserPassword = newPassword
	}
}

func DirDefault() string {
	return fmt.Sprintf(dirDefault, RootUser)
}

func RunCmd(cmd string, args ...string) error {
	client, err := dialClient(IP, RootUser, RootUserPassword)
	if err != nil {
		return err
	}
	defer silentClose(client)

	// scp copy
	dst, err := scpCopy(cmd, client)
	if err != nil {
		return err
	}

	session, err := client.NewSession()
	if err != nil {
		return err
	}
	defer silentClose(session)

	session.Stdout = os.Stdout
	session.Stderr = os.Stderr
	dir := filepath.Dir(dst)
	base := filepath.Base(dst)
	finalCmd := fmt.Sprintf("cd %[1]s && ./%[2]s %[3]s || printf \"==================================\\ncode:[%%d]\\n\" $?; rm -f %[2]s",
		dir, base, strings.Join(args, " "))

	fmt.Printf("%s@%s: \n==================================\n", RootUser, IP)
	_ = os.Stdout.Sync()
	return session.Run(finalCmd)
}

func dialClient(ip, user, password string) (*ssh.Client, error) {
	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	return ssh.Dial("tcp", ip+":22", config)
}

func scpCopy(src string, client *ssh.Client) (string, error) {
	dst := filepath.Join(DirDefault(), filepath.Base(src))
	srcFile, err := os.Open(src)
	if err != nil {
		return "", err
	}
	defer silentClose(srcFile)

	scpClient, err := scp.NewClientBySSH(client)
	if err != nil {
		return "", err
	}
	defer scpClient.Close()

	if err = scpClient.CopyFile(context.Background(), srcFile, dst, "0775"); err != nil {
		return "", err
	}

	return dst, nil
}

func silentClose(c io.Closer) {
	_ = c.Close()
}
