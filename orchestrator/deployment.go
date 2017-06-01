package orchestrator

import (
	"fmt"

	"strings"

	"github.com/pkg/errors"
)

const ArtifactDirectory = "/var/vcap/store/bbr-backup"

//go:generate counterfeiter -o fakes/fake_deployment.go . Deployment
type Deployment interface {
	HasBackupScript() bool
	HasUniqueCustomBackupNames() bool
	CheckArtifactDir() error
	IsRestorable() bool
	PreBackupLock() error
	Backup() error
	PostBackupUnlock() error
	Restore() error
	CopyRemoteBackupToLocal(Artifact) error
	CopyLocalBackupToRemote(Artifact) error
	Cleanup() error
	Instances() []Instance
	CustomArtifactNamesMatch() error
}

type deployment struct {
	Logger

	instances instances
}

func NewDeployment(logger Logger, instancesArray []Instance) Deployment {
	return &deployment{Logger: logger, instances: instances(instancesArray)}
}

func (bd *deployment) HasBackupScript() bool {
	backupableInstances := bd.instances.AllBackupable()
	return !backupableInstances.IsEmpty()
}

func (bd *deployment) HasUniqueCustomBackupNames() bool {
	names := bd.instances.CustomBlobNames()

	uniqueNames := map[string]bool{}
	for _, name := range names {
		if _, found := uniqueNames[name]; found {
			return false
		}
		uniqueNames[name] = true
	}
	return true
}

func (bd *deployment) CheckArtifactDir() error {
	errs := []string{}

	for _, inst := range bd.instances {
		exists, err := inst.ArtifactDirExists()
		if err != nil {
			errs = append(errs, fmt.Sprintf("Error checking %s on instance %s/%s", ArtifactDirectory, inst.Name(), inst.ID()))
		} else if exists {
			errs = append(errs, fmt.Sprintf("Directory %s already exists on instance %s/%s", ArtifactDirectory, inst.Name(), inst.ID()))
		}
	}

	if len(errs) > 0 {
		return errors.New(strings.Join(errs, "\n"))
	}

	return nil
}

func (bd *deployment) PreBackupLock() error {
	bd.Logger.Info("bbr", "Running pre-backup scripts...")
	err := bd.instances.AllPreBackupLockable().PreBackupLock()
	bd.Logger.Info("bbr", "Done.")
	return err
}

func (bd *deployment) Backup() error {
	bd.Logger.Info("bbr", "Running backup scripts...")
	return bd.instances.AllBackupable().Backup()
}

func (bd *deployment) PostBackupUnlock() error {
	bd.Logger.Info("bbr", "Running post-backup scripts...")
	err := bd.instances.AllPostBackupUnlockable().PostBackupUnlock()
	bd.Logger.Info("bbr", "Done.")
	return err
}

func (bd *deployment) Restore() error {
	bd.Logger.Info("bbr", "Running restore scripts...")
	return bd.instances.AllRestoreable().Restore()
}

func (bd *deployment) Cleanup() error {
	return bd.instances.Cleanup()
}

func (bd *deployment) IsRestorable() bool {
	restoreableInstances := bd.instances.AllRestoreable()
	return !restoreableInstances.IsEmpty()
}

func (bd *deployment) CustomArtifactNamesMatch() error {
	for _, instance := range bd.Instances() {
		jobName := instance.Name()
		for _, restoreName := range instance.CustomRestoreBlobNames() {
			var found bool
			for _, backupName := range bd.instances.CustomBlobNames() {
				if restoreName == backupName {
					found = true
				}
			}
			if !found {
				return errors.New(
					fmt.Sprintf(
						"The %s restore script expects a backup script which produces %s artifact which is not present in the deployment.",
						jobName,
						restoreName,
					),
				)
			}
		}
	}
	return nil
}

func (bd *deployment) CopyRemoteBackupToLocal(artifact Artifact) error {
	instances := bd.instances.AllBackupable()
	for _, instance := range instances {
		for _, backupBlob := range instance.BlobsToBackup() {
			writer, err := artifact.CreateFile(backupBlob)

			if err != nil {
				return err
			}

			size, err := backupBlob.Size()
			if err != nil {
				return err
			}

			bd.Logger.Info("bbr", "Copying backup -- %s uncompressed -- from %s/%s...", size, instance.Name(), instance.ID())
			if err := backupBlob.StreamFromRemote(writer); err != nil {
				return err
			}

			if err := writer.Close(); err != nil {
				return err
			}
			bd.Logger.Info("bbr", "Finished copying backup -- from %s/%s...", instance.Name(), instance.ID())

			bd.Logger.Info("bbr", "Starting validity checks")
			localChecksum, err := artifact.CalculateChecksum(backupBlob)
			if err != nil {
				return err
			}

			remoteChecksum, err := backupBlob.Checksum()
			if err != nil {
				return err
			}
			bd.Logger.Debug("bbr", "Comparing shasums")
			if !localChecksum.Match(remoteChecksum) {
				return errors.Errorf("Backup artifact is corrupted, checksum failed for %s/%s %s,  remote file: %s, local file: %s", instance.Name(), instance.ID(), backupBlob.Name(), remoteChecksum, localChecksum)
			}

			artifact.AddChecksum(backupBlob, localChecksum)

			err = backupBlob.Delete()
			if err != nil {
				return err
			}
			bd.Logger.Info("bbr", "Finished validity checks")
		}
	}
	return nil
}

func (bd *deployment) CopyLocalBackupToRemote(artifact Artifact) error {
	instances := bd.instances.AllRestoreable()

	for _, instance := range instances {
		for _, blob := range instance.BlobsToRestore() {
			reader, err := artifact.ReadFile(blob)

			if err != nil {
				return err
			}

			bd.Logger.Info("bbr", "Copying backup to %s/%s...", instance.Name(), instance.Index())
			if err := blob.StreamToRemote(reader); err != nil {
				return err
			} else {
				instance.MarkArtifactDirCreated()
			}

			localChecksum, err := artifact.FetchChecksum(blob)
			if err != nil {
				return err
			}

			remoteChecksum, err := blob.Checksum()
			if err != nil {
				return err
			}
			if !localChecksum.Match(remoteChecksum) {
				return errors.Errorf("Backup couldn't be transfered, checksum failed for %s/%s %s,  remote file: %s, local file: %s", instance.Name(), instance.ID(), blob.Name(), remoteChecksum, localChecksum)
			}
			bd.Logger.Info("bbr", "Done.")
		}
	}
	return nil
}

func (bd *deployment) Instances() []Instance {
	return bd.instances
}
