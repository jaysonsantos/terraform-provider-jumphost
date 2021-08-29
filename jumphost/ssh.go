package jumphost

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"strings"
	"sync"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"golang.org/x/crypto/ssh"
)

type SshClient struct {
	sshConfig *ssh.ClientConfig
	sshClient *ssh.Client
	hostname  string
	port      int

	tunnels map[string]*SshTunnel
	mutex   sync.Mutex
}

type SshTunnel struct {
	LocalPort int
	ctx       context.Context
}

func NewSshClient(hostname, username, password string, port int) SshClient {
	config := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}
	return SshClient{
		sshConfig: config,
		hostname:  hostname,
		port:      port,
		tunnels:   make(map[string]*SshTunnel),
		mutex:     sync.Mutex{},
	}
}

func (s *SshClient) Connect() error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if s.sshClient != nil {
		return nil
	}
	s.tunnels = make(map[string]*SshTunnel)

	address := fmt.Sprintf("%s:%d", s.hostname, s.port)
	log.Printf("Connecting to jumphost on %s", address)
	client, err := ssh.Dial("tcp", address, s.sshConfig)
	if err != nil {
		return fmt.Errorf("failed to open a connection to the jumphost %s %w", address, err)
	}

	s.sshClient = client

	return nil
}

func (s *SshClient) GetTunnel(ctx context.Context, d *schema.ResourceData) (*SshTunnel, error) {
	err := s.Connect()
	if err != nil {
		return nil, err
	}
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if err != nil {
		return nil, fmt.Errorf("failed to open tunnel")
	}

	address := fmt.Sprintf("%s:%d", d.Get(hostNameAttr).(string), d.Get(portAttr).(int))
	remoteConn, err := s.sshClient.Dial("tcp", address)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to remote service %w", err)
	}
	listenerConfig := net.ListenConfig{}
	listener, err := listenerConfig.Listen(ctx, "tcp", "127.0.0.1:0")
	if err != nil {
		return nil, fmt.Errorf("failed to listen on local port %w", err)
	}
	localPort, err := strconv.Atoi(strings.Split(listener.Addr().String(), ":")[1])
	if err != nil {
		return nil, fmt.Errorf("failed to parse listened port %w", err)
	}

	go func() {
		for {
			localConn, err := listener.Accept()
			if err != nil {
				panic(err) // TODO: Deal with this
			}
			pipeConnections(localConn, remoteConn)

		}
	}()

	tunnel := &SshTunnel{
		LocalPort: localPort,
		ctx:       ctx,
	}

	return tunnel, nil
}

func getCacheKey(d *schema.ResourceData) string {
	return ""
}

func pipeConnections(localConn, remoteConn net.Conn) {
	go io.Copy(localConn, remoteConn)
	go io.Copy(remoteConn, localConn)
}