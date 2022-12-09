// Copyright Project Harbor Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package contenttrust

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/goharbor/harbor/src/common/security"
	"github.com/goharbor/harbor/src/controller/artifact"
	"github.com/goharbor/harbor/src/controller/artifact/processor/image"
	"github.com/goharbor/harbor/src/controller/project"
	"github.com/goharbor/harbor/src/lib"
	"github.com/goharbor/harbor/src/pkg/accessory"
	accessorymodel "github.com/goharbor/harbor/src/pkg/accessory/model"
	basemodel "github.com/goharbor/harbor/src/pkg/accessory/model/base"
	proModels "github.com/goharbor/harbor/src/pkg/project/models"
	securitytesting "github.com/goharbor/harbor/src/testing/common/security"
	artifacttesting "github.com/goharbor/harbor/src/testing/controller/artifact"
	projecttesting "github.com/goharbor/harbor/src/testing/controller/project"
	"github.com/goharbor/harbor/src/testing/mock"
	accessorytesting "github.com/goharbor/harbor/src/testing/pkg/accessory"
)

type MiddlewareTestSuite struct {
	suite.Suite

	originalArtifactController artifact.Controller
	artifactController         *artifacttesting.Controller

	originalProjectController project.Controller
	projectController         *projecttesting.Controller

	artifact *artifact.Artifact
	project  *proModels.Project

	originalAccessMgr accessory.Manager
	accessMgr         *accessorytesting.Manager

	isArtifactSigned func(req *http.Request, art lib.ArtifactInfo) (bool, error)
	next             http.Handler
}

func (suite *MiddlewareTestSuite) SetupTest() {
	suite.originalArtifactController = artifact.Ctl
	suite.artifactController = &artifacttesting.Controller{}
	artifact.Ctl = suite.artifactController

	suite.originalProjectController = project.Ctl
	suite.projectController = &projecttesting.Controller{}
	project.Ctl = suite.projectController

	suite.originalAccessMgr = accessory.Mgr
	suite.accessMgr = &accessorytesting.Manager{}
	accessory.Mgr = suite.accessMgr

	suite.isArtifactSigned = isArtifactSigned
	suite.artifact = &artifact.Artifact{}
	suite.artifact.Type = image.ArtifactTypeImage
	suite.artifact.ProjectID = 1
	suite.artifact.RepositoryName = "library/photon"
	suite.artifact.Digest = "digest"

	suite.project = &proModels.Project{
		ProjectID: suite.artifact.ProjectID,
		Name:      "library",
		Metadata: map[string]string{
			proModels.ProMetaEnableContentTrust: "true",
		},
	}

	suite.next = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	isArtifactSigned = func(req *http.Request, art lib.ArtifactInfo) (bool, error) {
		return false, nil
	}
}

func (suite *MiddlewareTestSuite) TearDownTest() {
	artifact.Ctl = suite.originalArtifactController
	project.Ctl = suite.originalProjectController
	accessory.Mgr = suite.originalAccessMgr
}

func (suite *MiddlewareTestSuite) makeRequest() *http.Request {
	req := httptest.NewRequest("GET", "/v1/library/photon/manifests/2.0", nil)
	info := lib.ArtifactInfo{
		Repository: "library/photon",
		Reference:  "2.0",
		Tag:        "2.0",
		Digest:     "",
	}
	return req.WithContext(lib.WithArtifactInfo(req.Context(), info))
}

func (suite *MiddlewareTestSuite) TestGetArtifactFailed() {
	mock.OnAnything(suite.artifactController, "GetByReference").Return(nil, fmt.Errorf("error"))
	mock.OnAnything(suite.projectController, "GetByName").Return(suite.project, nil)

	req := suite.makeRequest()
	rr := httptest.NewRecorder()

	Notary()(suite.next).ServeHTTP(rr, req)
	suite.Equal(rr.Code, http.StatusInternalServerError)
}

func (suite *MiddlewareTestSuite) TestGetProjectFailed() {
	mock.OnAnything(suite.artifactController, "GetByReference").Return(suite.artifact, nil)
	mock.OnAnything(suite.projectController, "GetByName").Return(nil, fmt.Errorf("err"))

	req := suite.makeRequest()
	rr := httptest.NewRecorder()

	Notary()(suite.next).ServeHTTP(rr, req)
	suite.Equal(rr.Code, http.StatusInternalServerError)
}

func (suite *MiddlewareTestSuite) TestContentTrustDisabled() {
	mock.OnAnything(suite.artifactController, "GetByReference").Return(suite.artifact, nil)
	suite.project.Metadata[proModels.ProMetaEnableContentTrust] = "false"
	mock.OnAnything(suite.projectController, "GetByName").Return(suite.project, nil)

	req := suite.makeRequest()
	rr := httptest.NewRecorder()

	Notary()(suite.next).ServeHTTP(rr, req)
	suite.Equal(rr.Code, http.StatusOK)
}

func (suite *MiddlewareTestSuite) TestNoneArtifact() {
	req := httptest.NewRequest("GET", "/v1/library/photon/manifests/nonexist", nil)
	rr := httptest.NewRecorder()

	Notary()(suite.next).ServeHTTP(rr, req)
	suite.Equal(rr.Code, http.StatusNotFound)
}

func (suite *MiddlewareTestSuite) TestAuthenticatedUserPulling() {
	mock.OnAnything(suite.artifactController, "GetByReference").Return(suite.artifact, nil)
	mock.OnAnything(suite.projectController, "GetByName").Return(suite.project, nil)
	mock.OnAnything(suite.accessMgr, "List").Return([]accessorymodel.Accessory{}, nil)
	securityCtx := &securitytesting.Context{}
	mock.OnAnything(securityCtx, "Name").Return("local")
	mock.OnAnything(securityCtx, "Can").Return(true, nil)
	mock.OnAnything(securityCtx, "IsAuthenticated").Return(true)

	req := suite.makeRequest()
	req = req.WithContext(security.NewContext(req.Context(), securityCtx))
	rr := httptest.NewRecorder()

	Notary()(suite.next).ServeHTTP(rr, req)
	suite.Equal(rr.Code, http.StatusPreconditionFailed)
}

func (suite *MiddlewareTestSuite) TestScannerPulling() {
	mock.OnAnything(suite.artifactController, "GetByReference").Return(suite.artifact, nil)
	mock.OnAnything(suite.projectController, "GetByName").Return(suite.project, nil)
	mock.OnAnything(suite.accessMgr, "List").Return([]accessorymodel.Accessory{}, nil)
	securityCtx := &securitytesting.Context{}
	mock.OnAnything(securityCtx, "Name").Return("v2token")
	mock.OnAnything(securityCtx, "Can").Return(true, nil)
	mock.OnAnything(securityCtx, "IsAuthenticated").Return(true)

	req := suite.makeRequest()
	req = req.WithContext(security.NewContext(req.Context(), securityCtx))
	rr := httptest.NewRecorder()

	Notary()(suite.next).ServeHTTP(rr, req)
	suite.Equal(rr.Code, http.StatusOK)
}

// pull a public project a un-signed image when policy checker is enabled.
func (suite *MiddlewareTestSuite) TestUnAuthenticatedUserPulling() {
	mock.OnAnything(suite.artifactController, "GetByReference").Return(suite.artifact, nil)
	mock.OnAnything(suite.projectController, "GetByName").Return(suite.project, nil)
	mock.OnAnything(suite.accessMgr, "List").Return([]accessorymodel.Accessory{}, nil)
	securityCtx := &securitytesting.Context{}
	mock.OnAnything(securityCtx, "Name").Return("local")
	mock.OnAnything(securityCtx, "Can").Return(true, nil)
	mock.OnAnything(securityCtx, "IsAuthenticated").Return(false)

	req := suite.makeRequest()
	rr := httptest.NewRecorder()

	Notary()(suite.next).ServeHTTP(rr, req)
	suite.Equal(rr.Code, http.StatusPreconditionFailed)
}

// pull cosign signature when policy checker is enabled.
func (suite *MiddlewareTestSuite) TestSignaturePulling() {
	mock.OnAnything(suite.artifactController, "GetByReference").Return(suite.artifact, nil)
	mock.OnAnything(suite.projectController, "GetByName").Return(suite.project, nil)
	acc := &basemodel.Default{
		Data: accessorymodel.AccessoryData{
			ID:            1,
			ArtifactID:    2,
			SubArtifactID: 1,
			Type:          accessorymodel.TypeCosignSignature,
		},
	}
	mock.OnAnything(suite.accessMgr, "List").Return([]accessorymodel.Accessory{
		acc,
	}, nil)

	req := suite.makeRequest()
	rr := httptest.NewRecorder()

	Notary()(suite.next).ServeHTTP(rr, req)
	suite.Equal(rr.Code, http.StatusOK)
}

func TestMiddlewareTestSuite(t *testing.T) {
	suite.Run(t, &MiddlewareTestSuite{})
}
