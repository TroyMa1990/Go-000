/*
 * Create By Xinwenjia 2020-05-28
 */

package hsbcodes_test

import (
    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"

    . "https://github.com/TroyMa1990/Go-000/tree/main/Week04/codes"
)

var _ = Describe("Codes", func() {
    Describe("Test Codes.", func() {
        Context("Test Codes String", func() {
            It("Should Return String Value", func() {
                Expect(OK.String()).To(Equal("OK"))
                Expect(MyServerLoadStageFailed.String()).To(Equal("MyServerLoadStageFailed"))
            })
        })

        Context("Test Codes Error", func() {
            It("Should Return String Value", func() {
                err := MyServerLoadStageFailed.Error()
                code := FromError(err)
                Expect(code).To(Equal(MyServerLoadStageFailed))
            })
        })
    })
})
