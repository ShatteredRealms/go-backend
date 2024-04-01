package e2e_test

import (
	"bytes"
	"context"
	"encoding/gob"
	"os"
	"path/filepath"
	"testing"

	"github.com/bxcodec/faker/v4"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	clientset *kubernetes.Clientset
	testData  *setupData
)

func TestE2e(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "E2E Suite")
}

type setupData struct {
	// Namespace used for SRO
	Namespace string

	// Namespace used for Agones
	AgonesNamespace string
}

var _ = SynchronizedBeforeSuite(func(ctx context.Context) []byte {
	testData := setupData{
		Namespace:       "sro-testing-" + faker.Username(),
		AgonesNamespace: "agones-system-testing-" + faker.Username(),
	}
	setupClientSet()

	Expect(clientset.CoreV1().Namespaces().Create(
		ctx,
		&corev1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name: testData.Namespace,
			},
		},
		metav1.CreateOptions{},
	)).To(Succeed())
	Expect(clientset.CoreV1().Namespaces().Create(
		ctx,
		&corev1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name: testData.AgonesNamespace,
			},
		},
		metav1.CreateOptions{},
	)).To(Succeed())

	return testData.encode()
}, func(data []byte) {
	testData = decodeSetupData(data)
	setupClientSet()
})

var _ = SynchronizedAfterSuite(func() {}, func(ctx context.Context) {
	Expect(clientset.CoreV1().Namespaces().Delete(ctx, testData.Namespace, metav1.DeleteOptions{})).
		NotTo(HaveOccurred())
	Expect(clientset.CoreV1().Namespaces().Delete(ctx, testData.AgonesNamespace, metav1.DeleteOptions{})).
		NotTo(HaveOccurred())
})

func setupClientSet() {
	homeDir, err := os.UserHomeDir()
	Expect(err).NotTo(HaveOccurred(), "unable to get home directory")
	Expect(homeDir).NotTo(BeEmpty())

	config, err := clientcmd.BuildConfigFromFlags("", filepath.Join(homeDir, ".kube", "config"))
	Expect(err).NotTo(HaveOccurred(), "unable to get kubernetes config")

	clientset, err = kubernetes.NewForConfig(config)
	Expect(err).NotTo(HaveOccurred(), "kubernetes config invalid")
}

func (s *setupData) encode() []byte {
	var buf *bytes.Buffer
	enc := gob.NewEncoder(buf)
	Expect(enc.Encode(s)).NotTo(HaveOccurred())
	return buf.Bytes()
}

func decodeSetupData(in []byte) *setupData {
	buf := bytes.NewBuffer(in)
	out := &setupData{}
	dec := gob.NewDecoder(buf)
	Expect(dec.Decode(out)).NotTo(HaveOccurred())
	return out
}
