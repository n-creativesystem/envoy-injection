package cert

import (
	"crypto/tls"
	"fmt"
	"path/filepath"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/n-creativesystem/kubernetes-extensions/pkg/logger"
)

type CertNotify struct {
	mu       sync.RWMutex
	certFile string
	keyFile  string
	keyPair  *tls.Certificate
	watcher  *fsnotify.Watcher
	watching chan bool
}

func New(certFile, keyFile string) (*CertNotify, error) {
	var err error
	certFile, err = filepath.Abs(certFile)
	if err != nil {
		return nil, err
	}
	keyFile, err = filepath.Abs(keyFile)
	if err != nil {
		return nil, err
	}
	cm := &CertNotify{
		mu:       sync.RWMutex{},
		certFile: certFile,
		keyFile:  keyFile,
	}
	return cm, nil
}

func (cm *CertNotify) Watch() error {
	var err error
	if cm.watcher, err = fsnotify.NewWatcher(); err != nil {
		return fmt.Errorf("cert: can't create watcher: %s", err)
	}
	if err = cm.watcher.Add(cm.certFile); err != nil {
		return fmt.Errorf("cert: can't watch cert file: %s", err)
	}
	if err = cm.watcher.Add(cm.keyFile); err != nil {
		return fmt.Errorf("cert: can't watch key file: %s", err)
	}
	if err := cm.load(); err != nil {
		return fmt.Errorf("cert: can't load cert or key file: %v", err)
	}
	logger.Infof("cert: watching for cert and key change")
	cm.watching = make(chan bool)
	go cm.run()
	return nil
}

func (cm *CertNotify) load() error {
	keyPair, err := tls.LoadX509KeyPair(cm.certFile, cm.keyFile)
	if err == nil {
		cm.mu.Lock()
		cm.keyPair = &keyPair
		cm.mu.Unlock()
		logger.Infof("cert: certificate and key loaded")
	}
	return err
}

func (cm *CertNotify) run() {
loop:
	for {
		select {
		case <-cm.watching:
			break loop
		case event := <-cm.watcher.Events:
			logger.Infof("cert: watch event: %v", event)
			if err := cm.load(); err != nil {
				logger.Errorf(fmt.Sprintf("cert: can't load cert or key file: %v", err))
			}
		case err := <-cm.watcher.Errors:
			logger.Errorf(fmt.Sprintf("cert: error watching files: %v", err))
		}
	}
	logger.Infof("cert: stopped watching")
	cm.watcher.Close()
}

func (cm *CertNotify) GetCertificate(hello *tls.ClientHelloInfo) (*tls.Certificate, error) {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	return cm.keyPair, nil
}

func (cm *CertNotify) Stop() {
	cm.watching <- false
}
