package integration

import (
	"fmt"
	"os"

	"github.com/pivotal-cf-experimental/cf-webmock/mockbosh"
	"github.com/pivotal-cf-experimental/cf-webmock/mockhttp"
	"github.com/pivotal-cf/pcf-backup-and-restore/testssh"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Backup", func() {
	var director *mockhttp.Server
	var testSSHServer *testssh.Server

	AfterEach(func() {
		director.VerifyMocks()
	})
	BeforeEach(func() {
		director = mockbosh.NewTLS()
		director.ExpectedBasicAuth("admin", "admin")
		testSSHServer = testssh.New(GinkgoWriter)
	})
	AfterEach(func() {
		Expect(os.RemoveAll("/var/vcap/jobs")).To(Succeed())
	})

	Context("with deployment present", func() {
		BeforeEach(func() {
			director.VerifyAndMock(
				mockbosh.VMsForDeployment("my-new-deployment").RedirectsToTask(14),
				mockbosh.Task(14).RespondsWithTaskContainingState(mockbosh.TaskDone),
				mockbosh.Task(14).RespondsWithTaskContainingState(mockbosh.TaskDone),
				mockbosh.TaskEvent(14).RespondsWithVMsOutput([]string{}),
				mockbosh.TaskOutput(14).RespondsWithVMsOutput([]mockbosh.VMsOutput{
					{
						IPs:     []string{"10.0.0.1"},
						JobName: "redis-dedicated-node",
					},
				}),
				mockbosh.StartSSHSession("my-new-deployment").RedirectsToTask(15),
				mockbosh.Task(15).RespondsWithTaskContainingState(mockbosh.TaskDone),
				mockbosh.Task(15).RespondsWithTaskContainingState(mockbosh.TaskDone),
				mockbosh.TaskEvent(15).RespondsWith("{}"),
				mockbosh.TaskOutput(15).RespondsWith(fmt.Sprintf(`[{"status":"success",
				"ip":"127.0.0.1:%d",
				"host_public_key":"not-relevant",
				"index":0}]`, testSSHServer.Port)),
				mockbosh.CleanupSSHSession("my-new-deployment").RedirectsToTask(16),
				mockbosh.Task(16).RespondsWithTaskContainingState(mockbosh.TaskDone),
			)
		})

		It("backs up deployment successfully", func() {
			filesExistOnVM(
				"/var/vcap/jobs/redis/bin/backup",
			)

			session := runBinary([]string{"BOSH_PASSWORD=admin"}, "--ca-cert", sslCertPath, "--username", "admin", "--target", director.URL, "--deployment", "my-new-deployment", "--debug", "backup")
			Expect(session.ExitCode()).To(BeZero())
		})

		It("errors if a deployment cant be backuped", func() {
			filesExistOnVM(
				"/var/vcap/jobs/redis/bin/ctl",
			)

			session := runBinary([]string{"BOSH_PASSWORD=admin"}, "--ca-cert", sslCertPath, "--username", "admin", "--target", director.URL, "--deployment", "my-new-deployment", "--debug", "backup")
			Expect(session.ExitCode()).NotTo(BeZero())
			Expect(string(session.Err.Contents())).To(ContainSubstring("Deployment 'my-new-deployment' has no backup scripts"))
		})
	})

	It("returns error if deployment not found", func() {
		director.VerifyAndMock(mockbosh.VMsForDeployment("my-new-deployment").NotFound())

		session := runBinary([]string{"BOSH_PASSWORD=admin"}, "--ca-cert", sslCertPath, "--username", "admin", "--target", director.URL, "--deployment", "my-new-deployment", "backup")

		Expect(session.ExitCode()).To(Equal(1))
		Expect(string(session.Err.Contents())).To(ContainSubstring("Director responded with non-successful status code"))
	})
})
