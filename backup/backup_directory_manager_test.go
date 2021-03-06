package backup_test

import (
	"fmt"
	"os"

	"time"

	. "github.com/cloudfoundry-incubator/bosh-backup-and-restore/backup"
	. "github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/config"
	. "github.com/onsi/gomega"
	"io/ioutil"
)

var _ = Context("BackupManager", func() {
	var artifactPath string
	var deploymentName string
	var backupName string
	var backupManager = BackupDirectoryManager{}
	var err error
	fakeClock := func() time.Time {
		return time.Date(2015, 10, 21, 02, 2, 3, 0, time.FixedZone("UTC+1", 3600))
	}

	BeforeEach(func() {
		artifactPath, err = ioutil.TempDir("", "test-backup-artifact-dir")
		Expect(err).NotTo(HaveOccurred())

		deploymentName = fmt.Sprintf("my-cool-redis-%d", config.GinkgoConfig.ParallelNode)
		backupName = deploymentName + "_20151021T010203Z"
	})

	AfterEach(func() {
		Expect(os.RemoveAll(backupName)).To(Succeed())
	})

	Describe("Create", func() {
		JustBeforeEach(func() {
			_, err = backupManager.Create(artifactPath, deploymentName, nil, fakeClock)
		})

		It("creates a directory with the given name", func() {
			Expect(err).NotTo(HaveOccurred())
			Expect(fmt.Sprintf("%s/%s", artifactPath, backupName)).To(BeADirectory())
		})

		Context("when the artifact directory cannot be created", func() {
			BeforeEach(func() {
				Expect(os.MkdirAll(fmt.Sprintf("%s/%s", artifactPath, backupName), 0777)).To(Succeed())
			})

			It("returns an error", func() {
				Expect(err).To(MatchError(ContainSubstring("failed creating artifact directory")))
			})
		})

		Context("when the artifact path does not exist", func() {
			BeforeEach(func() {
				artifactPath = "/myawesomedir"
			})

			It("returns an error", func() {
				Expect(err).To(HaveOccurred())
			})
		})

		Context("when the artifact path is not a directory", func() {
			BeforeEach(func() {
				file, err := ioutil.TempFile(os.TempDir(), "test-backup-artifact-not-a-dir")
				Expect(err).NotTo(HaveOccurred())
				artifactPath = file.Name()
			})

			It("returns an error", func() {
				Expect(err).To(MatchError(fmt.Sprintf("artifact path %s is not a directory", artifactPath)))
			})
		})
	})

	Describe("Open", func() {
		Context("when the directory exists", func() {
			BeforeEach(func() {
				err := os.MkdirAll(backupName, 0700)
				Expect(err).NotTo(HaveOccurred())
			})

			It("does not create a directory", func() {
				_, err := backupManager.Open(backupName, nil)
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("when the directory does not exist", func() {
			It("fails", func() {
				_, err := backupManager.Open(backupName, nil)
				Expect(err).To(MatchError(ContainSubstring("failed opening the directory")))
				Expect(backupName).NotTo(BeADirectory())
			})
		})
	})
})
