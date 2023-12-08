package utils_context_test

import (
	"context"
	"fmt"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"cnores-skeleton-golang-app/app/shared/utils/context"
)

var _ = Describe("Context", func() {

	When("An empty traceparent is received", func() {
		It("Should fill the trace parent components with zero and generate span id successfully", func() {

			emptyTraceParent := ""
			traceParentComponent, ok := utils_context.BuildTraceParent(emptyTraceParent)
			Expect(ok).To(BeFalse())
			Expect(traceParentComponent[utils_context.TraceVersion]).To(Equal(utils_context.TraceFlagsEmptyStructure))
			Expect(traceParentComponent[utils_context.ParentID]).To(Equal(utils_context.ParentIDEmptyStructure))
			Expect(traceParentComponent[utils_context.TraceID]).To(Equal(utils_context.TraceIDEmptyStructure))
			Expect(traceParentComponent[utils_context.TraceFlags]).To(Equal(utils_context.TraceFlagsEmptyStructure))
			Expect(len(traceParentComponent[utils_context.SpanID].(string))).To(Equal(utils_context.ParentIDLarge))
		})
	})

	When("A right traceparent is received", func() {
		It("Should fill all trace parent components successfully", func() {

			version := "07"
			traceId := "0af7651916cd43dd8448eb211c80319c"
			parentId := "b7ad6b7169203331"
			traceFlags := "00"

			emptyTraceParent := fmt.Sprintf("%s-%s-%s-%s", version, traceId, parentId, traceFlags)

			traceParentComponent, ok := utils_context.BuildTraceParent(emptyTraceParent)
			Expect(ok).To(BeTrue())
			Expect(traceParentComponent[utils_context.TraceVersion]).To(Equal(version))
			Expect(traceParentComponent[utils_context.ParentID]).To(Equal(parentId))
			Expect(traceParentComponent[utils_context.TraceID]).To(Equal(traceId))
			Expect(traceParentComponent[utils_context.TraceFlags]).To(Equal(traceFlags))
			Expect(len(traceParentComponent[utils_context.SpanID].(string))).To(Equal(utils_context.ParentIDLarge))
		})
	})

	When("A right parentId  is received by output ", func() {
		It("Should generate traceparent with parentId", func() {

			version := "07"
			traceId := "0af7651916cd43dd8448eb211c80319c"
			parentId := "b7ad6b7169203331"
			traceFlags := "00"
			spanId := "0000007169203331"

			expectedTraceOutput := fmt.Sprintf("%s-%s-%s-%s", version, traceId, parentId, traceFlags)

			traceComponent := map[string]interface{}{}
			traceComponent[utils_context.TraceVersion] = version
			traceComponent[utils_context.TraceID] = traceId
			traceComponent[utils_context.ParentID] = parentId
			traceComponent[utils_context.TraceFlags] = traceFlags
			traceComponent[utils_context.SpanID] = spanId
			ctx := context.Background()
			ctx = utils_context.FillContextFromTraceComponent(traceComponent, ctx)

			outputTraceParent := utils_context.BuildOutputTraceParent(ctx)

			Expect(outputTraceParent).To(Equal(expectedTraceOutput))
		})
	})

	When("A wrong parentId  is received by output ", func() {
		It("Should generate traceparent with spanId", func() {

			version := "07"
			traceId := "0af7651916cd43dd8448eb211c80319c"
			parentId := utils_context.ParentIDEmptyStructure
			traceFlags := "00"
			spanId := "0000007169203331"

			expectedTraceOutput := fmt.Sprintf("%s-%s-%s-%s", version, traceId, spanId, traceFlags)

			traceComponent := map[string]interface{}{}
			traceComponent[utils_context.TraceVersion] = version
			traceComponent[utils_context.TraceID] = traceId
			traceComponent[utils_context.ParentID] = parentId
			traceComponent[utils_context.TraceFlags] = traceFlags
			traceComponent[utils_context.SpanID] = spanId
			ctx := context.Background()
			ctx = utils_context.FillContextFromTraceComponent(traceComponent, ctx)

			outputTraceParent := utils_context.BuildOutputTraceParent(ctx)

			Expect(outputTraceParent).To(Equal(expectedTraceOutput))
		})
	})

	When("A right fieldsToLog is received", func() {
		It("Should fill all fields attributes as map successfully", func() {

			orderInformationMap := map[string]string{"orderId": "3201706872", "shippingGroupId": "3201706872"}
			ctx := context.Background()

			contextWithFieldsToLog := utils_context.FillContextFromMap(ctx, orderInformationMap)
			Expect(contextWithFieldsToLog.Value("fieldsToLog")).To(Equal(orderInformationMap))
		})
	})

})
