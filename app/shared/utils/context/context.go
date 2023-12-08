package utils_context

import (
	"context"
	"crypto/rand"
	"fmt"
	"cnores-skeleton-golang-app/app/shared/utils/log"
	"regexp"
	"strings"
	"time"
)

const (
	TraceVersionEmptyStructure    = "00"
	TraceIDEmptyStructure         = "00000000000000000000000000000000"
	ParentIDEmptyStructure        = "0000000000000000"
	TraceFlagsEmptyStructure      = "00"
	traceIDLarge                  = 32
	ParentIDLarge                 = 16
	traceVersionLarge             = 2
	traceFlagsLarge               = 2
	spanIDBytes                   = 8
	zeroPrefix                    = "0"
	validTraceParentRegex         = "^[0-9a-fA-F]{2}-[0-9a-fA-F]{32}-[0-9a-fA-F]{16}-[0-9a-fA-F]{2}$"
	TraceVersion                  = "version"
	TraceID                       = "traceId"
	ParentID                      = "parentId"
	TraceFlags                    = "traceFlags"
	SpanID                        = "spanId"
	CreatedAt                     = "createdAt"
	TraceParent                   = "traceparent"
	InputTraceParent              = "inputTraceParent"
	OutputTraceParent             = "outputTraceParent"
	traceParentComponentNumber    = 4
	TraceParentInvalidFormat      = "TraceParent header has an invalid format, new traceability components have been generated"
	TraceParentComponentSeparator = "-"
)

func CreateTraceContext(traceInput string) context.Context {
	traceContext := context.Background()
	traceGenerated, ok := BuildTraceParent(traceInput)
	if !ok {
		log.Error("%s %s Headers: %v", InputTraceParent, TraceParentInvalidFormat, traceInput)
	}
	traceContext = FillContextFromTraceComponent(traceGenerated, traceContext)
	return traceContext
}

func GetLogFromContext(ctx context.Context, layer string, module string) *log.Logger {
	return log.NewLogger(ctx, log.Fields{
		"Layer":  layer,
		"Module": module,
	})
}

func GetElapsedTimeFromContext(ctx context.Context) float64 {
	createdAt := ctx.Value("CreatedAt").(time.Time)
	return time.Since(createdAt).Seconds()
}

func BuildTraceParent(traceParent string) (map[string]interface{}, bool) {
	validFormat := false
	if validateRegex(traceParent, validTraceParentRegex) {
		validFormat = true
	}
	spanID := GenerateSpanID()

	traceParentComponent := []string{TraceVersionEmptyStructure, TraceIDEmptyStructure, ParentIDEmptyStructure, TraceFlagsEmptyStructure}
	if (len(strings.Split(traceParent, TraceParentComponentSeparator)) == traceParentComponentNumber) && validFormat {
		traceParentComponent = strings.Split(traceParent, TraceParentComponentSeparator)
	}
	traceComponent := map[string]interface{}{}
	traceComponent[TraceVersion] = fillOutTraceComponent(traceParentComponent[0], traceVersionLarge)
	traceComponent[TraceID] = fillOutTraceComponent(traceParentComponent[1], traceIDLarge)
	traceComponent[ParentID] = fillOutTraceComponent(traceParentComponent[2], ParentIDLarge)
	traceComponent[TraceFlags] = fillOutTraceComponent(traceParentComponent[3], traceFlagsLarge)
	traceComponent[SpanID] = spanID
	traceComponent[CreatedAt] = time.Now()
	return traceComponent, validFormat
}

func BuildOutputTraceParent(ctx context.Context) string {
	var outputParentId string
	version := ctx.Value("TraceVersion").(string)
	traceId := ctx.Value("TraceID").(string)
	parentId, ok := ctx.Value("ParentID").(string)
	outputParentId = parentId
	if !ok || parentId == ParentIDEmptyStructure {
		spanId := ctx.Value("SpanID").(string)
		outputParentId = spanId
	}
	traceFlags := ctx.Value("TraceFlags").(string)

	traceOutput := fmt.Sprintf("%s-%s-%s-%s", version, traceId, outputParentId, traceFlags)
	log.Info("%s traceparent generated %s", OutputTraceParent, traceOutput)
	return traceOutput
}

func FillContextFromTraceComponent(traceComponent map[string]interface{}, ctx context.Context) context.Context {
	ctx = context.WithValue(ctx, "TraceVersion", traceComponent[TraceVersion])
	ctx = context.WithValue(ctx, "TraceID", traceComponent[TraceID])
	ctx = context.WithValue(ctx, "ParentID", traceComponent[ParentID])
	ctx = context.WithValue(ctx, "TraceFlags", traceComponent[TraceFlags])
	ctx = context.WithValue(ctx, "SpanID", traceComponent[SpanID])
	ctx = context.WithValue(ctx, "CreatedAt", traceComponent[CreatedAt])
	return ctx
}

func FillContextFromMap(ctx context.Context, mapValue map[string]string) context.Context {
	ctx = context.WithValue(ctx, "fieldsToLog", mapValue)
	return ctx
}

func validateRegex(traceComponent, regex string) bool {
	match, _ := regexp.MatchString(regex, traceComponent)
	return match
}

func GenerateSpanID() string {
	spanMask := make([]byte, spanIDBytes)
	if _, err := rand.Read(spanMask); err != nil {
		panic(err)
	}
	spanID := fmt.Sprintf("%x", spanMask)
	return spanID
}

func fillOutTraceComponent(traceComponent string, totalLarge int) string {
	traceComponentLarge := len(traceComponent)
	for i := traceComponentLarge; i <= totalLarge; i++ {
		if len(traceComponent) < totalLarge {
			traceComponent = zeroPrefix + traceComponent
		}
	}
	return traceComponent
}
