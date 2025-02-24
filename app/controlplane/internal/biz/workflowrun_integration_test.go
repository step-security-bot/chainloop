//
// Copyright 2023 The Chainloop Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package biz_test

import (
	"context"
	"testing"

	"github.com/chainloop-dev/chainloop/app/controlplane/internal/biz"
	"github.com/chainloop-dev/chainloop/app/controlplane/internal/biz/testhelpers"
	"github.com/chainloop-dev/chainloop/internal/credentials"
	creds "github.com/chainloop-dev/chainloop/internal/credentials/mocks"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	"github.com/secure-systems-lab/go-securesystemslib/dsse"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

func (s *workflowRunIntegrationTestSuite) TestSaveAttestation() {
	assert := assert.New(s.T())
	ctx := context.Background()

	validEnvelope := &dsse.Envelope{}

	s.T().Run("non existing workflowRun", func(t *testing.T) {
		err := s.WorkflowRun.SaveAttestation(ctx, uuid.NewString(), validEnvelope, validDigest)
		assert.Error(err)
		assert.True(biz.IsNotFound(err))
	})

	s.T().Run("valid workflowRun", func(t *testing.T) {
		run, err := s.WorkflowRun.Create(ctx, &biz.WorkflowRunCreateOpts{
			WorkflowID: s.workflowOrg1.ID.String(), RobotaccountID: s.robotAccount.ID.String(), ContractRevisionUUID: s.contractVersion.Version.ID, CASBackendID: s.casBackend.ID,
		})
		assert.NoError(err)

		err = s.WorkflowRun.SaveAttestation(ctx, run.ID.String(), validEnvelope, validDigest)
		assert.NoError(err)

		// Retrieve attestation ref from storage and compare
		r, err := s.WorkflowRun.GetByIDInOrgOrPublic(ctx, s.org.ID, run.ID.String())
		assert.NoError(err)
		assert.Equal(r.Attestation, &biz.Attestation{Envelope: validEnvelope, Digest: validDigest})
	})

	s.T().Run("valid workflowRun attestation not stored in CAS", func(t *testing.T) {
		run, err := s.WorkflowRun.Create(ctx, &biz.WorkflowRunCreateOpts{
			WorkflowID: s.workflowOrg1.ID.String(), RobotaccountID: s.robotAccount.ID.String(), ContractRevisionUUID: s.contractVersion.Version.ID, CASBackendID: s.casBackend.ID,
		})
		assert.NoError(err)

		err = s.WorkflowRun.SaveAttestation(ctx, run.ID.String(), validEnvelope, "")
		assert.NoError(err)

		// Retrieve attestation ref from storage and compare
		r, err := s.WorkflowRun.GetByIDInOrgOrPublic(ctx, s.org.ID, run.ID.String())
		assert.NoError(err)
		assert.Equal(r.Attestation, &biz.Attestation{Envelope: validEnvelope, Digest: ""})
	})
}

func (s *workflowRunIntegrationTestSuite) TestGetByIDInOrgOrPublic() {
	assert := assert.New(s.T())
	ctx := context.Background()
	testCases := []struct {
		name    string
		orgID   string
		runID   string
		wantErr bool
	}{
		{
			name:    "non existing workflowRun",
			orgID:   s.org.ID,
			runID:   uuid.NewString(),
			wantErr: true,
		},
		{
			name:  "existing workflowRun in org1",
			orgID: s.org.ID,
			runID: s.runOrg1.ID.String(),
		},
		{
			name:    "can't access workflowRun from other org",
			orgID:   s.org.ID,
			runID:   s.runOrg2.ID.String(),
			wantErr: true,
		},
		{
			name:  "can access workflowRun from other org if public",
			orgID: s.org.ID,
			runID: s.runOrg2Public.ID.String(),
		},
	}

	for _, tc := range testCases {
		s.T().Run(tc.name, func(t *testing.T) {
			run, err := s.WorkflowRun.GetByIDInOrgOrPublic(ctx, tc.orgID, tc.runID)
			if tc.wantErr {
				assert.Error(err)
				assert.True(biz.IsNotFound(err))
			} else {
				assert.NoError(err)
				assert.Equal(tc.runID, run.ID.String())
			}
		})
	}
}

func (s *workflowRunIntegrationTestSuite) TestGetByDigestInOrgOrPublic() {
	assert := assert.New(s.T())
	ctx := context.Background()
	testCases := []struct {
		name           string
		orgID          string
		digest         string
		errTypeChecker func(err error) bool
	}{
		{
			name:           "non existing workflowRun",
			orgID:          s.org.ID,
			digest:         "sha256:b5bb9d8014a0f9b1d61e21e796d78dccdf1352f23cd32812f4850b878ae4944c",
			errTypeChecker: biz.IsNotFound,
		},
		{
			name:           "invalid digest",
			orgID:          s.org.ID,
			digest:         "b5bb9d8014a0f9b1d61e21e796d78dccdf1352f23cd32812f4850b878ae4944c",
			errTypeChecker: biz.IsErrValidation,
		},
		{
			name:   "existing workflowRun in org1",
			orgID:  s.org.ID,
			digest: validDigest,
		},
		{
			name:           "can't access workflowRun from other org",
			orgID:          s.org.ID,
			digest:         validDigest2,
			errTypeChecker: biz.IsNotFound,
		},
		{
			name:   "can access workflowRun from other org if public",
			orgID:  s.org.ID,
			digest: validDigest3,
		},
	}

	for _, tc := range testCases {
		s.T().Run(tc.name, func(t *testing.T) {
			run, err := s.WorkflowRun.GetByDigestInOrgOrPublic(ctx, tc.orgID, tc.digest)
			if tc.errTypeChecker != nil {
				assert.Error(err)
				assert.True(tc.errTypeChecker(err))
			} else {
				assert.NoError(err)
				assert.Equal(tc.digest, run.Attestation.Digest)
			}
		})
	}
}

func (s *workflowRunIntegrationTestSuite) TestCreate() {
	assert := assert.New(s.T())
	ctx := context.Background()

	s.T().Run("valid workflowRun", func(t *testing.T) {
		run, err := s.WorkflowRun.Create(ctx, &biz.WorkflowRunCreateOpts{
			WorkflowID: s.workflowOrg1.ID.String(), RobotaccountID: s.robotAccount.ID.String(), ContractRevisionUUID: s.contractVersion.Version.ID, CASBackendID: s.casBackend.ID,
			RunnerType: "runnerType", RunnerRunURL: "runURL",
		})
		assert.NoError(err)
		if diff := cmp.Diff(&biz.WorkflowRun{
			RunnerType: "runnerType", RunURL: "runURL", State: string(biz.WorkflowRunInitialized), ContractVersionID: s.contractVersion.Version.ID,
			Workflow:    s.workflowOrg1,
			CASBackends: []*biz.CASBackend{s.casBackend},
		}, run,
			cmpopts.IgnoreFields(biz.WorkflowRun{}, "CreatedAt", "ID", "Workflow"),
			cmpopts.IgnoreFields(biz.CASBackend{}, "CreatedAt", "ValidatedAt", "OrganizationID"),
		); diff != "" {
			assert.Failf("mismatch (-want +got):\n%s", diff)
		}
	})
}

// Run the tests
func TestWorkflowRunUseCase(t *testing.T) {
	suite.Run(t, new(workflowRunIntegrationTestSuite))
}

// Utility struct to hold the test suite
type workflowRunIntegrationTestSuite struct {
	testhelpers.UseCasesEachTestSuite
	org, org2                                      *biz.Organization
	casBackend                                     *biz.CASBackend
	workflowOrg1, workflowOrg2, workflowPublicOrg2 *biz.Workflow
	runOrg1, runOrg2, runOrg2Public                *biz.WorkflowRun
	robotAccount                                   *biz.RobotAccount
	contractVersion                                *biz.WorkflowContractWithVersion
}

func (s *workflowRunIntegrationTestSuite) SetupTest() {
	var err error
	assert := assert.New(s.T())
	ctx := context.Background()
	// OCI repository credentials
	credsWriter := creds.NewReaderWriter(s.T())
	credsWriter.On(
		"SaveCredentials", ctx, mock.Anything, &credentials.OCIKeypair{Repo: "repo", Username: "username", Password: "pass"},
	).Return("stored-OCI-secret", nil)

	s.TestingUseCases = testhelpers.NewTestingUseCases(s.T(), testhelpers.WithCredsReaderWriter(credsWriter))

	s.org, err = s.Organization.Create(ctx, "testing org")
	assert.NoError(err)
	s.org2, err = s.Organization.Create(ctx, "second org")
	assert.NoError(err)

	// Workflow
	s.workflowOrg1, err = s.Workflow.Create(ctx, &biz.CreateOpts{Name: "test workflow", OrgID: s.org.ID})
	assert.NoError(err)
	s.workflowOrg2, err = s.Workflow.Create(ctx, &biz.CreateOpts{Name: "test workflow", OrgID: s.org2.ID})
	assert.NoError(err)
	// Public workflow
	s.workflowPublicOrg2, err = s.Workflow.Create(ctx, &biz.CreateOpts{Name: "test public workflow", OrgID: s.org2.ID, Public: true})
	assert.NoError(err)

	// Robot account
	s.robotAccount, err = s.RobotAccount.Create(ctx, "name", s.org.ID, s.workflowOrg1.ID.String())
	assert.NoError(err)

	// Find contract revision
	s.contractVersion, err = s.WorkflowContract.Describe(ctx, s.org.ID, s.workflowOrg1.ContractID.String(), 0)
	assert.NoError(err)

	s.casBackend, err = s.CASBackend.CreateOrUpdate(ctx, s.org.ID, "repo", "username", "pass", biz.CASBackendOCI, true)
	assert.NoError(err)

	// Let's create 3 runs, one in org1 and 2 in org2 (one public)
	s.runOrg1, err = s.WorkflowRun.Create(ctx,
		&biz.WorkflowRunCreateOpts{
			WorkflowID: s.workflowOrg1.ID.String(), RobotaccountID: s.robotAccount.ID.String(), ContractRevisionUUID: s.contractVersion.Version.ID, CASBackendID: s.casBackend.ID,
		})
	assert.NoError(err)
	assert.NoError(s.WorkflowRun.SaveAttestation(ctx, s.runOrg1.ID.String(), &dsse.Envelope{}, validDigest))

	s.runOrg2, err = s.WorkflowRun.Create(ctx,
		&biz.WorkflowRunCreateOpts{
			WorkflowID: s.workflowOrg2.ID.String(), RobotaccountID: s.robotAccount.ID.String(), ContractRevisionUUID: s.contractVersion.Version.ID, CASBackendID: s.casBackend.ID,
		})
	assert.NoError(err)
	assert.NoError(s.WorkflowRun.SaveAttestation(ctx, s.runOrg2.ID.String(), &dsse.Envelope{}, validDigest2))

	s.runOrg2Public, err = s.WorkflowRun.Create(ctx,
		&biz.WorkflowRunCreateOpts{
			WorkflowID: s.workflowPublicOrg2.ID.String(), RobotaccountID: s.robotAccount.ID.String(), ContractRevisionUUID: s.contractVersion.Version.ID, CASBackendID: s.casBackend.ID,
		})
	assert.NoError(err)
	assert.NoError(s.WorkflowRun.SaveAttestation(ctx, s.runOrg2Public.ID.String(), &dsse.Envelope{}, validDigest3))
}
