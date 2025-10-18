package tests

import (
	"os"
	"os/exec"
	"syscall"
	"testing"
	"time"
)

func TestSyscalls(t *testing.T) {
	cmd := exec.Command(os.Args[0], "-test.run=^helperSignalHandler$")
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true} // OWN IT'S OWN PID SO WE CAN SEND SIGTERM'S TO IT INSTEAD OF THE TEST FUNC.
	if err := cmd.Start(); err != nil {
		t.Fatalf("Error starting sub-process: %v", err)
	}
	time.Sleep(50 * time.Millisecond)
	if err := syscall.Kill(-cmd.Process.Pid, syscall.SIGTERM); err != nil {
		t.Fatalf("Error at sending SIGTERM to sub-process: %v", err)
	}
	select {
	case <-waitChan(cmd):
		return
	case <-time.After(50 * time.Millisecond):
		t.Fatalf("Graceful shutdown isn't triggered in 50 ms")
	} // either get the signal in 50 ms,or fail.

}
func waitChan(cmd *exec.Cmd) chan error {
	ch := make(chan error)
	go func() {
		ch <- cmd.Wait()
	}()
	return ch
}
