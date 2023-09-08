package runtime

import (
	"bufio"
	"fmt"
	"io"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	logging "github.com/openshift/cluster-logging-operator/apis/logging/v1"
	"github.com/openshift/cluster-logging-operator/internal/runtime"
	"github.com/openshift/cluster-logging-operator/test"
	"github.com/openshift/cluster-logging-operator/test/client"
	. "github.com/openshift/cluster-logging-operator/test/matchers"
	corev1 "k8s.io/api/core/v1"
)

var _ = Describe("Object", func() {
	var (
		clf = NewClusterLogForwarder()
	)

	DescribeTable("Decode",
		func(manifest string, o runtime.Object) {
			got := runtime.Decode(manifest)
			Expect(got).To(EqualDiff(o), "%#v", manifest)
		},
		Entry("YAML string clf", test.YAMLString(clf), clf),
	)
})

var _ = Describe("NewClusterLogging", func() {
	It("should stub a ClusterLogging", func() {
		Expect(NewClusterLogging().Spec).To(Equal(logging.ClusterLoggingSpec{
			Collection: &logging.CollectionSpec{
				Type: logging.LogCollectionTypeFluentd,
			},
			ManagementState: logging.ManagementStateManaged,
		}))
	})
})

var _ = Describe("pod reader and writer", func() {
	var (
		t   *client.Test
		pod *corev1.Pod
	)
	BeforeEach(func() {
		t = client.NewTest()
		pod = runtime.NewPod(t.NS.Name, "testpod", corev1.Container{
			Name: "testpod", Image: "quay.io/quay/busybox", Args: []string{"sleep", "1h"},
		})
		Expect(t.Create(pod)).To(Succeed())
		Expect(t.WaitFor(pod, client.PodRunning)).To(Succeed(), test.YAMLString(pod))
	})

	AfterEach(func() {
		_ = t.Delete(pod)
		t.Close()
	})

	Context("NewPodWriter", func() {
		It("write files to pod", func() {
			w, err := PodWriter(pod, "", "/tmp/testme")
			Expect(err).To(Succeed())
			_, err = fmt.Fprintln(w, "hello world")
			Expect(err).To(Succeed())
			Expect(w.Close()).To(Succeed())

			Eventually(func() string {
				out, err := Exec(pod, "", "cat", "/tmp/testme").CombinedOutput()
				Expect(err).To(Succeed())
				return string(out)
			}).WithTimeout(time.Second).Should(Equal("hello world\n"))
		})
	})

	Context("NewPodReader", func() {
		It("reads existing file", func() {
			Expect(Exec(pod, "", "sh", "-c", "echo hello world > /tmp/testme").Run()).To(Succeed())
			r, err := PodReader(pod, "", "/tmp/testme")
			Expect(err).To(Succeed())
			b, err := io.ReadAll(r)
			Expect(err).To(Succeed())
			Expect(string(b)).To(Equal("hello world\n"))
		})
	})

	Context("NewPodTailReader", func() {
		It("tails non-existent file", func() {
			// Create tail reader before file exists
			r, err := PodTailReader(pod, "", "/tmp/testme")
			Expect(err).To(Succeed())

			// Write some data
			w, err := PodWriter(pod, "", "/tmp/testme")
			Expect(err).To(Succeed())
			fmt.Fprintln(w, "hello world")

			// Read from the tail reader.
			scan := bufio.NewScanner(r)
			Expect(scan.Scan()).To(BeTrue(), scan.Err())
			Expect(scan.Text()).To(Equal("hello world"))

			// Write and read some more data
			fmt.Fprintln(w, "goodbye")
			Expect(scan.Scan()).To(BeTrue(), scan.Err())
			Expect(scan.Text()).To(Equal("goodbye"))
		})
	})
})
