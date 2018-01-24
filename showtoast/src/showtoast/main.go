package main

import (
	"errors"
	"runtime"
	"unsafe"

	"golang.org/x/mobile/app"
	"golang.org/x/mobile/event/lifecycle"
	"golang.org/x/mobile/event/paint"
)

/*
#cgo LDFLAGS: -landroid -llog
#include <jni.h>
#include <stdlib.h>
#include <android/log.h>

JavaVM* current_vm;
jobject current_ctx;

#define LOG_INFO(...) __android_log_print(ANDROID_LOG_INFO, "Go", __VA_ARGS__)
#define LOG_FATAL(...) __android_log_print(ANDROID_LOG_FATAL, "Go", __VA_ARGS__)

char* _lockJNI(uintptr_t* envp, int* attachedp) {
	JNIEnv* env;
	if (current_vm == NULL) {
		return "no current JVM";
	}
	*attachedp = 0;
	switch ((*current_vm)->GetEnv(current_vm, (void**)&env, JNI_VERSION_1_6)) {
	case JNI_OK:
		break;
	case JNI_EDETACHED:
		if ((*current_vm)->AttachCurrentThread(current_vm, &env, 0) != 0) {
			return "cannot attach to JVM";
		}
		*attachedp = 1;
		break;
	case JNI_EVERSION:
		return "bad JNI version";
	default:
		return "unknown JNI error from GetEnv";
	}
	*envp = (uintptr_t)env;
	return NULL;
}

char* _checkException(uintptr_t jnienv) {
	jthrowable exc;
	JNIEnv* env = (JNIEnv*)jnienv;
	if (!(*env)->ExceptionCheck(env)) {
		return NULL;
	}
	exc = (*env)->ExceptionOccurred(env);
	(*env)->ExceptionClear(env);
	jclass clazz = (*env)->FindClass(env, "java/lang/Throwable");
	jmethodID toString = (*env)->GetMethodID(env, clazz, "toString", "()Ljava/lang/String;");
	jobject msgStr = (*env)->CallObjectMethod(env, exc, toString);
	return (char*)(*env)->GetStringUTFChars(env, msgStr, 0);
}

void _unlockJNI() {
	(*current_vm)->DetachCurrentThread(current_vm);
}

void showToast(JNIEnv* env) {
	jclass clazz = (*env)->FindClass(env, "com/example/tenntenn/showtoast/ToastViewer");
    jmethodID methodShow = (*env)->GetStaticMethodID(env, clazz, "showToast", "(Landroid/app/Activity)V");
    if (methodShow == NULL) {
		LOG_INFO("cannot find method");
        return;
    }
    (*env)->CallVoidMethod(env, clazz, methodShow, current_ctx);
}
*/
import "C"

func RunOnJVM(fn func(vm, env, ctx uintptr) error) error {
	errch := make(chan error)
	go func() {
		runtime.LockOSThread()
		defer runtime.UnlockOSThread()

		env := C.uintptr_t(0)
		attached := C.int(0)
		if errStr := C._lockJNI(&env, &attached); errStr != nil {
			errch <- errors.New(C.GoString(errStr))
			return
		}
		if attached != 0 {
			defer C._unlockJNI()
		}

		vm := uintptr(unsafe.Pointer(C.current_vm))
		if err := fn(vm, uintptr(env), uintptr(C.current_ctx)); err != nil {
			errch <- err
			return
		}

		if exc := C._checkException(env); exc != nil {
			errch <- errors.New(C.GoString(exc))
			C.free(unsafe.Pointer(exc))
			return
		}
		errch <- nil
	}()
	return <-errch
}

func main() {
	app.Main(func(a app.App) {
		for e := range a.Events() {
			switch a.Filter(e).(type) {
			case lifecycle.Event:
			case paint.Event:
				RunOnJVM(func(vm, jniEnv, ctx uintptr) error {
					env := (*C.JNIEnv)(unsafe.Pointer(jniEnv))
					C.showToast(env)
					return nil
				})
				a.Publish()
			}
		}
	})
}
