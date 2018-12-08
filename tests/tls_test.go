package tests

import (
	"fmt"
	"strconv"

	"github.com/deiscc/workflow-e2e/tests/cmd"
	"github.com/deiscc/workflow-e2e/tests/cmd/apps"
	"github.com/deiscc/workflow-e2e/tests/cmd/auth"
	"github.com/deiscc/workflow-e2e/tests/cmd/builds"
	"github.com/deiscc/workflow-e2e/tests/model"
	"github.com/deiscc/workflow-e2e/tests/settings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	. "github.com/onsi/gomega/gexec"
)

var _ = Describe("deis tls", func() {

	Context("with an existing user", func() {

		var user model.User

		BeforeEach(func() {
			user = auth.RegisterAndLogin()
		})

		AfterEach(func() {
			auth.Cancel(user)
		})

		Context("who owns an existing app that has already been deployed", func() {

			var app model.App

			BeforeEach(func() {
				app = apps.Create(user, "--no-remote")
				builds.Create(user, app)
			})

			AfterEach(func() {
				apps.Destroy(user, app)
			})

			Specify("can enable/disable tls", func() {
				sess, err := cmd.Start("deis tls:enable --app=%s", &user, app.Name)
				Eventually(sess, settings.MaxEventuallyTimeout).Should(Say("done"))
				Expect(err).NotTo(HaveOccurred())
				Eventually(sess).Should(Exit(0))

				// curl the app's root URL and ensure we get a 301 redirect
				cmdRetryTimeout := 60
				curlCmd := model.Cmd{CommandLineString: fmt.Sprintf(`curl -sL -w "%%{http_code}\\n" "%s" -o /dev/null`, app.URL)}
				Eventually(cmd.Retry(curlCmd, strconv.Itoa(301), cmdRetryTimeout)).Should(BeTrue())

				sess, err = cmd.Start("deis tls:disable --app=%s", &user, app.Name)
				Eventually(sess, settings.MaxEventuallyTimeout).Should(Say("done"))
				Expect(err).NotTo(HaveOccurred())
				Eventually(sess).Should(Exit(0))

				cmdRetryTimeout = 60
				curlCmd = model.Cmd{CommandLineString: fmt.Sprintf(`curl -sL -w "%%{http_code}\\n" "%s" -o /dev/null`, app.URL)}
				Eventually(cmd.Retry(curlCmd, strconv.Itoa(200), cmdRetryTimeout)).Should(BeTrue())
			})
		})
	})
})
