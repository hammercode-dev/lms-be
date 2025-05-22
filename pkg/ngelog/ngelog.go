package ngelog

import (
	"context"

	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/trace"
)

type AddFields = map[string]interface{}

type MasterLog struct {
	service   string
	namespace string
	log       *logrus.Logger
}

var defaultFormatJSON = &logrus.JSONFormatter{
	FieldMap: logrus.FieldMap{
		logrus.FieldKeyMsg:  "message",
		logrus.FieldKeyTime: "@timestamp",
	},
	DisableHTMLEscape: true,
}

func NewMasterLog(serviceName, nameSpace string) *MasterLog {
	log := logrus.New()
	log.SetFormatter(defaultFormatJSON)
	return &MasterLog{
		log:       log,
		service:   serviceName,
		namespace: nameSpace,
	}
}

var mlog = NewMasterLog("LMS-BE", "staging")

func SetServiceName(name string) {
	mlog.service = name
}

// staging or production
func SetNameSpace(name string) {
	mlog.namespace = name
}

func (mLog MasterLog) combineFields(fields []map[string]any) logrus.Fields {
	merged := logrus.Fields{}
	merged["service"] = mLog.service
	if mLog.namespace != "" {
		merged["namespace"] = mLog.namespace
	}
	for _, f := range fields {
		for k, v := range f {
			merged[k] = v
		}
	}
	return merged
}

// get all log fields and return log entry
func (mLog MasterLog) getLogEntry(ctx context.Context, fields ...map[string]any) *logrus.Entry {
	span := trace.SpanFromContext(ctx)
	sc := span.SpanContext()

	allFields := mLog.combineFields(fields)

	// Add tracing info if available
	if sc.IsValid() {
		allFields["trace_id"] = sc.TraceID().String()
		allFields["span_id"] = sc.SpanID().String()
	}

	allFields["service"] = mLog.service
	allFields["namespace"] = mLog.namespace

	return mLog.log.WithContext(ctx).WithFields(allFields)
}

func Info(ctx context.Context, message string, fields ...map[string]any) {
	mlog.getLogEntry(ctx, fields...).Info(message)
}

func Error(ctx context.Context, message string, err error, fields ...map[string]any) {
	mlog.getLogEntry(ctx, fields...).WithError(err).Error(message)
}

func Fatal(ctx context.Context, message string, err error, fields ...map[string]any) {
	mlog.getLogEntry(ctx, fields...).WithError(err).Fatal(message)
}

func FatalPanic(ctx context.Context, message string, err error, fields ...map[string]any) {
	mlog.getLogEntry(ctx, fields...).WithError(err).Fatalf(message)
}
