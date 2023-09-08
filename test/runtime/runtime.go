// package runtime provides functions to run processes in containers and to create logging resources.
package runtime

import (
	"fmt"
	"io"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/openshift/cluster-logging-operator/internal/utils"
	corev1 "k8s.io/api/core/v1"
)

// Exec returns an `oc exec` exec.Cmd to run 'cmd args...' on a container.
// If container=="", use the default container.
// The caller is in control of running the returned Cmd and redirecting IO.
func Exec(pod *corev1.Pod, container, cmd string, args ...string) *exec.Cmd {
	ocArgs := []string{
		"exec",
		"-i",
		"-n", pod.GetNamespace(),
	}
	if container != "" {
		ocArgs = append(ocArgs, "-c", strings.ToLower(container))
	}
	ocArgs = append(ocArgs, "pod/"+pod.GetName(), "--", cmd)
	ocArgs = append(ocArgs, args...)
	return exec.Command("oc", ocArgs...)
}

type reader struct {
	io.ReadCloser
	cmd *exec.Cmd
}

func (r *reader) Close() error {
	err := r.ReadCloser.Close()
	_ = r.cmd.Wait()
	return err
}

// PodReader returns an io.ReadCloser that reads from an existing file in a container on a pod.
func PodReader(pod *corev1.Pod, container, filepath string) (io.ReadCloser, error) {
	cmd := Exec(pod, container, "cat", filepath)
	out, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	return &reader{cmd: cmd, ReadCloser: out}, utils.WrapError(cmd.Start())
}

// PodTailReader returns an io.ReadCloser that tails a file in a container on a pod.
//
// Uses 'tail -F' so file does not need to exist immediately.
func PodTailReader(pod *corev1.Pod, container, filepath string) (io.ReadCloser, error) {
	cmd := Exec(pod, container, "tail", "-F", filepath)
	out, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	err = cmd.Start()
	return &reader{cmd: cmd, ReadCloser: out}, err
}

type writer struct {
	io.WriteCloser // Pipe to stdin
	cmd            *exec.Cmd
}

func (w writer) Close() error {
	err := w.WriteCloser.Close()
	_ = w.cmd.Wait()
	return err
}

// PodWriter returns an io.WriteCloser that appends to a file in a pod container.
//
// If container == "", use default container.
// If the file has a directory part, attempts to create it with mkdir -p.
func PodWriter(pod *corev1.Pod, container, filename string) (io.WriteCloser, error) {
	dir := filepath.Dir(filename)
	cmd := Exec(pod, container, "sh", "-c", fmt.Sprintf("mkdir -p %v && cat > %v", dir, filename))
	w, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}
	return writer{WriteCloser: w, cmd: cmd}, utils.WrapError(cmd.Start())
}

// PodWrite is a shortcut for NewPodWriter(), Write(), Close()
func PodWrite(pod *corev1.Pod, container, filename string, data []byte) error {
	w, err := PodWriter(pod, container, filename)
	if err == nil {
		defer w.Close()
		_, err = w.Write(data)
	}
	return err
}
