package odin_test

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/go-test/deep"

	"github.com/Devex/spaceflight/mcc/odin/odin"
)

type getRestoreDBInstanceFromDBSnapshotInputCase struct {
	name                         string
	identifier                   string
	restoreParams                odin.RestoreParams
	snapshot                     *rds.DBSnapshot
	expectedRestoreSnapshotInput *rds.RestoreDBInstanceFromDBSnapshotInput
	expectedError                string
}

var getRestoreDBInstanceFromDBSnapshotInputCases = []getRestoreDBInstanceFromDBSnapshotInputCase{
	// Params with Snapshot
	{
		name:       "Params with Snapshot",
		identifier: "production-rds",
		restoreParams: odin.RestoreParams{
			InstanceType:         "db.m1.medium",
			OriginalInstanceName: "production",
		},
		snapshot: exampleSnapshot1,
		expectedRestoreSnapshotInput: &rds.RestoreDBInstanceFromDBSnapshotInput{
			DBInstanceClass:      aws.String("db.m1.medium"),
			DBInstanceIdentifier: aws.String("production-rds"),
			DBSnapshotIdentifier: aws.String("rds:production-2015-06-11"),
			DBSubnetGroupName:    aws.String(""),
			Engine:               aws.String("postgres"),
		},
		expectedError: "",
	},
	// Params with Snapshot without OriginalInstanceName
	{
		name:       "Params with Snapshot without OriginalInstanceName",
		identifier: "production-rds",
		restoreParams: odin.RestoreParams{
			InstanceType: "db.m1.medium",
		},
		snapshot:                     exampleSnapshot1,
		expectedRestoreSnapshotInput: nil,
		expectedError:                "Original Instance Name was empty",
	},
	// Params with non existing Snapshot
	{
		name:       "Params with non existing Snapshot",
		identifier: "production-rds",
		restoreParams: odin.RestoreParams{
			InstanceType:         "db.m1.medium",
			OriginalInstanceName: "develop",
		},
		snapshot:                     exampleSnapshot1,
		expectedRestoreSnapshotInput: nil,
		expectedError:                "Couldn't find snapshot for develop instance",
	},
}

func TestGetRestoreDBInstanceFromDBSnapshotInput(t *testing.T) {
	svc := newMockRDSClient()
	for _, useCase := range getRestoreDBInstanceFromDBSnapshotInputCases {
		t.Run(
			useCase.name,
			func(t *testing.T) {
				if useCase.snapshot != nil {
					svc.dbSnapshots[*useCase.snapshot.DBInstanceIdentifier] = []*rds.DBSnapshot{useCase.snapshot}
				}
				restoreSnapshotInput, err := useCase.restoreParams.GetRestoreDBInstanceFromDBSnapshotInput(
					useCase.identifier,
					svc,
				)
				if err != nil {
					if useCase.expectedError == "" ||
						err.Error() != useCase.expectedError {
						t.Errorf(
							"Unexpected error happened: %v",
							err,
						)
					}
				} else {
					if diff := deep.Equal(
						restoreSnapshotInput,
						useCase.expectedRestoreSnapshotInput,
					); diff != nil {
						t.Errorf(
							"Unexpected output: %s should be %s",
							restoreSnapshotInput,
							useCase.expectedRestoreSnapshotInput,
						)
					}
				}
			},
		)
	}
}
