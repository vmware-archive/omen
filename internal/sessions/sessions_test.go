package sessions_test

import (
	"github.com/pivotal-cloudops/omen/internal/sessions"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal-cloudops/omen/internal/sessions/sessionsfakes"
	"errors"
)

var _ = Describe("Sessions", func() {

	Describe("clear sessions", func() {

		It("deletes all active sessions via the opsman api", func() {
			fakeClient := &sessionsfakes.FakeOpsmanClient{}
			fakeClient.DeleteReturns(nil)

			manager := sessions.NewSessionManager(fakeClient)
			Expect(manager.ClearAll()).To(Succeed())

			Expect(fakeClient.DeleteCallCount()).To(Equal(1))
			req, _ := fakeClient.DeleteArgsForCall(0)
			Expect(req).To(ContainSubstring("/api/v0/sessions"))
		})

		It("fails if the status code is not 200", func() {
			fakeClient := &sessionsfakes.FakeOpsmanClient{}
			fakeClient.DeleteReturns(errors.New("Something is wrong"))

			manager := sessions.NewSessionManager(fakeClient)
			Expect(manager.ClearAll()).ToNot(Succeed())
		})

	})

})
