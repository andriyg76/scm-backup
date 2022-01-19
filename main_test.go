package main

import (
	"github.com/andriyg76/glog"
	"testing"
)

func init() {
	glog.ToFileAndConsole("file.log", glog.INFO, glog.TRACE)
}

func TestFileLog(t *testing.T) {
	glog.Trace("trace %s %s", "'a'", "'b'")
	glog.Debug("debug %s %s", "'a'", "'b'")
	glog.Info("info %s %s", "'a'", "'b'")
	glog.Warn("warn %s %s", "'a'", "'b'")
	glog.Error("trace %s %s", "'a'", "'b'")
	func() {
		defer func() {
			recover()
		}()
		glog.Panic("panic %s %s", "'a'", "'b'")
	}()
	glog.Fatal("fatal %s %s", "'a'", "'b'")
}
