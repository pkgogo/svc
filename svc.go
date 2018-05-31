package svc

import (
        "os"
        "os/signal"
        "syscall"
)

var DefaultSignals = []os.Signal{syscall.SIGINT, syscall.SIGTERM}

type Service interface {
        Start() error
        Stop() error
}

type SignalHandler interface {
        HandleSignal(os.Signal) bool
        CatchSignals() []os.Signal
}

func Run(svc Service) error {
        if err := svc.Start(); err != nil {
                return err
        }

        signalChan := make(chan os.Signal, 1)

        ss := DefaultSignals

        sh, ok := svc.(SignalHandler)
        if ok {
                ss = sh.CatchSignals()
        }

        signal.Notify(signalChan, ss...)

        for {
                s := <-signalChan
                if ok && sh.HandleSignal(s) {
                        break
                }
        }

        return svc.Stop()
}
