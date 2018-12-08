package tests

import (
	"fmt"

	"github.com/deiscc/workflow-e2e/tests/cmd"
	"github.com/deiscc/workflow-e2e/tests/cmd/auth"
	"github.com/deiscc/workflow-e2e/tests/cmd/keys"
	"github.com/deiscc/workflow-e2e/tests/model"
	"github.com/deiscc/workflow-e2e/tests/settings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	. "github.com/onsi/gomega/gexec"
)

var _ = Describe("deis keys", func() {

	Context("with an existing user", func() {

		var user model.User

		BeforeEach(func() {
			user = auth.RegisterAndLogin()
		})

		AfterEach(func() {
			auth.Cancel(user)
		})

		Context("who has at least one key", func() {

			var keyName string

			BeforeEach(func() {
				keyName, _ = keys.Add(user)
			})

			Specify("that user can list their own keys", func() {
				sess, err := cmd.Start("deis keys:list", &user)
				Eventually(sess, settings.MaxEventuallyTimeout).Should(Say(fmt.Sprintf("%s ssh-rsa", keyName)))
				Expect(err).NotTo(HaveOccurred())
				Eventually(sess).Should(Exit(0))
			})

		})

		Specify("that user can add and remove keys", func() {
			keyName, _ := keys.Add(user)
			keys.Remove(user, keyName)
		})

	})

})
