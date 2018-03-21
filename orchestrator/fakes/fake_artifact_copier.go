// Code generated by counterfeiter. DO NOT EDIT.
package fakes

import (
	"sync"

	"github.com/cloudfoundry-incubator/bosh-backup-and-restore/orchestrator"
)

type FakeArtifactCopier struct {
	DownloadBackupFromDeploymentStub        func(orchestrator.Backup, orchestrator.Deployment) error
	downloadBackupFromDeploymentMutex       sync.RWMutex
	downloadBackupFromDeploymentArgsForCall []struct {
		arg1 orchestrator.Backup
		arg2 orchestrator.Deployment
	}
	downloadBackupFromDeploymentReturns struct {
		result1 error
	}
	downloadBackupFromDeploymentReturnsOnCall map[int]struct {
		result1 error
	}
	UploadBackupToDeploymentStub        func(orchestrator.Backup, orchestrator.Deployment) error
	uploadBackupToDeploymentMutex       sync.RWMutex
	uploadBackupToDeploymentArgsForCall []struct {
		arg1 orchestrator.Backup
		arg2 orchestrator.Deployment
	}
	uploadBackupToDeploymentReturns struct {
		result1 error
	}
	uploadBackupToDeploymentReturnsOnCall map[int]struct {
		result1 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeArtifactCopier) DownloadBackupFromDeployment(arg1 orchestrator.Backup, arg2 orchestrator.Deployment) error {
	fake.downloadBackupFromDeploymentMutex.Lock()
	ret, specificReturn := fake.downloadBackupFromDeploymentReturnsOnCall[len(fake.downloadBackupFromDeploymentArgsForCall)]
	fake.downloadBackupFromDeploymentArgsForCall = append(fake.downloadBackupFromDeploymentArgsForCall, struct {
		arg1 orchestrator.Backup
		arg2 orchestrator.Deployment
	}{arg1, arg2})
	fake.recordInvocation("DownloadBackupFromDeployment", []interface{}{arg1, arg2})
	fake.downloadBackupFromDeploymentMutex.Unlock()
	if fake.DownloadBackupFromDeploymentStub != nil {
		return fake.DownloadBackupFromDeploymentStub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1
	}
	return fake.downloadBackupFromDeploymentReturns.result1
}

func (fake *FakeArtifactCopier) DownloadBackupFromDeploymentCallCount() int {
	fake.downloadBackupFromDeploymentMutex.RLock()
	defer fake.downloadBackupFromDeploymentMutex.RUnlock()
	return len(fake.downloadBackupFromDeploymentArgsForCall)
}

func (fake *FakeArtifactCopier) DownloadBackupFromDeploymentArgsForCall(i int) (orchestrator.Backup, orchestrator.Deployment) {
	fake.downloadBackupFromDeploymentMutex.RLock()
	defer fake.downloadBackupFromDeploymentMutex.RUnlock()
	return fake.downloadBackupFromDeploymentArgsForCall[i].arg1, fake.downloadBackupFromDeploymentArgsForCall[i].arg2
}

func (fake *FakeArtifactCopier) DownloadBackupFromDeploymentReturns(result1 error) {
	fake.DownloadBackupFromDeploymentStub = nil
	fake.downloadBackupFromDeploymentReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeArtifactCopier) DownloadBackupFromDeploymentReturnsOnCall(i int, result1 error) {
	fake.DownloadBackupFromDeploymentStub = nil
	if fake.downloadBackupFromDeploymentReturnsOnCall == nil {
		fake.downloadBackupFromDeploymentReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.downloadBackupFromDeploymentReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeArtifactCopier) UploadBackupToDeployment(arg1 orchestrator.Backup, arg2 orchestrator.Deployment) error {
	fake.uploadBackupToDeploymentMutex.Lock()
	ret, specificReturn := fake.uploadBackupToDeploymentReturnsOnCall[len(fake.uploadBackupToDeploymentArgsForCall)]
	fake.uploadBackupToDeploymentArgsForCall = append(fake.uploadBackupToDeploymentArgsForCall, struct {
		arg1 orchestrator.Backup
		arg2 orchestrator.Deployment
	}{arg1, arg2})
	fake.recordInvocation("UploadBackupToDeployment", []interface{}{arg1, arg2})
	fake.uploadBackupToDeploymentMutex.Unlock()
	if fake.UploadBackupToDeploymentStub != nil {
		return fake.UploadBackupToDeploymentStub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1
	}
	return fake.uploadBackupToDeploymentReturns.result1
}

func (fake *FakeArtifactCopier) UploadBackupToDeploymentCallCount() int {
	fake.uploadBackupToDeploymentMutex.RLock()
	defer fake.uploadBackupToDeploymentMutex.RUnlock()
	return len(fake.uploadBackupToDeploymentArgsForCall)
}

func (fake *FakeArtifactCopier) UploadBackupToDeploymentArgsForCall(i int) (orchestrator.Backup, orchestrator.Deployment) {
	fake.uploadBackupToDeploymentMutex.RLock()
	defer fake.uploadBackupToDeploymentMutex.RUnlock()
	return fake.uploadBackupToDeploymentArgsForCall[i].arg1, fake.uploadBackupToDeploymentArgsForCall[i].arg2
}

func (fake *FakeArtifactCopier) UploadBackupToDeploymentReturns(result1 error) {
	fake.UploadBackupToDeploymentStub = nil
	fake.uploadBackupToDeploymentReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeArtifactCopier) UploadBackupToDeploymentReturnsOnCall(i int, result1 error) {
	fake.UploadBackupToDeploymentStub = nil
	if fake.uploadBackupToDeploymentReturnsOnCall == nil {
		fake.uploadBackupToDeploymentReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.uploadBackupToDeploymentReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeArtifactCopier) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.downloadBackupFromDeploymentMutex.RLock()
	defer fake.downloadBackupFromDeploymentMutex.RUnlock()
	fake.uploadBackupToDeploymentMutex.RLock()
	defer fake.uploadBackupToDeploymentMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeArtifactCopier) recordInvocation(key string, args []interface{}) {
	fake.invocationsMutex.Lock()
	defer fake.invocationsMutex.Unlock()
	if fake.invocations == nil {
		fake.invocations = map[string][][]interface{}{}
	}
	if fake.invocations[key] == nil {
		fake.invocations[key] = [][]interface{}{}
	}
	fake.invocations[key] = append(fake.invocations[key], args)
}

var _ orchestrator.ArtifactCopier = new(FakeArtifactCopier)
