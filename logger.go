package bone

//BLogger 定义的日志接口 方便第三方实现
type BLogger interface {
	Trace(v ...interface{})
	Tracef(format string, v ...interface{})
	Debug(v ...interface{})
	Debugf(format string, v ...interface{})
	Info(v ...interface{})
	Infof(format string, v ...interface{})
	Warn(v ...interface{})
	Warnf(format string, v ...interface{})
	Error(v ...interface{})
	Errorf(format string, v ...interface{})
	Fatal(v ...interface{})
	Fatalf(format string, v ...interface{})
}

//emptyBLogger 默认的logger实现 忽略任何日志
type emptyBLogger struct {
}

func (eb *emptyBLogger) Trace(v ...interface{}) {

}
func (eb *emptyBLogger) Tracef(format string, v ...interface{}) {

}
func (eb *emptyBLogger) Debug(v ...interface{}) {

}
func (eb *emptyBLogger) Debugf(format string, v ...interface{}) {

}
func (eb *emptyBLogger) Info(v ...interface{}) {

}
func (eb *emptyBLogger) Infof(format string, v ...interface{}) {

}
func (eb *emptyBLogger) Warn(v ...interface{}) {

}
func (eb *emptyBLogger) Warnf(format string, v ...interface{}) {

}
func (eb *emptyBLogger) Error(v ...interface{}) {

}
func (eb *emptyBLogger) Errorf(format string, v ...interface{}) {

}
func (eb *emptyBLogger) Fatal(v ...interface{}) {

}
func (eb *emptyBLogger) Fatalf(format string, v ...interface{}) {

}
