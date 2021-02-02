package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"syscall"
	"unsafe"

	"github.com/creack/pty"
	"github.com/gliderlabs/ssh"
)

const (
	defaultPath  = "/debug-bridge"
	defaultShell = "/debug-bridge/ash"
)

func setWinsize(f *os.File, w, h int) {
	syscall.Syscall(syscall.SYS_IOCTL, f.Fd(), uintptr(syscall.TIOCSWINSZ),
		uintptr(unsafe.Pointer(&struct{ h, w, x, y uint16 }{uint16(h), uint16(w), 0, 0})))
}

func main() {
	listen := os.Getenv("SSH_LISTEN")
	if listen == "" {
		listen = ":2222"
	}
	sshKey := os.Getenv("SSH_AUTHORIZED_KEY")
	if sshKey == "" {
		log.Fatal("you need to specify SSH_AUTHORIZED_KEY")
	}
	path := os.Getenv("DEBUG_BRIDGE_PATH")
	if path == "" {
		path = defaultPath
	}
	shell := os.Getenv("DEBUG_BRIDGE_SHELL")
	if shell == "" {
		shell = defaultShell
	}

	myKey, _, _, _, err := ssh.ParseAuthorizedKey([]byte(sshKey))
	if err != nil {
		log.Fatalf("could not parse pubkey: %s", err.Error())
	}

	ssh.Handle(func(s ssh.Session) {
		log.Printf("accepted connection from %s\n", s.RemoteAddr().String())

		cmd := exec.Command(shell)
		ptyReq, winCh, isPty := s.Pty()
		if isPty {
			cmd.Env = append(cmd.Env, fmt.Sprintf("TERM=%s", ptyReq.Term))
			cmd.Env = append(cmd.Env, "PATH="+path)
			f, err := pty.Start(cmd)
			if err != nil {
				panic(err)
			}
			go func() {
				for win := range winCh {
					setWinsize(f, win.Width, win.Height)
				}
			}()
			go func() {
				io.Copy(f, s) // stdin
			}()
			io.Copy(s, f) // stdout
			cmd.Wait()
		} else {
			io.WriteString(s, "No PTY requested.\n")
			s.Exit(1)
		}
	})

	keyAuth := ssh.PublicKeyAuth(func(ctx ssh.Context, key ssh.PublicKey) bool {
		return ssh.KeysEqual(key, myKey)
	})

	log.Printf("starting ssh server on %s...\n", listen)
	log.Fatal(ssh.ListenAndServe(listen, nil, keyAuth))
}
