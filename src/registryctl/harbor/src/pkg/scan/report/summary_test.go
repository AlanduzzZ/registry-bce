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

package report

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/goharbor/harbor/src/pkg/scan/dao/scan"
	v1 "github.com/goharbor/harbor/src/pkg/scan/rest/v1"
	"github.com/goharbor/harbor/src/pkg/scan/vuln"
)

// SummaryTestSuite is test suite for testing report summary.
type SummaryTestSuite struct {
	suite.Suite

	r *scan.Report
}

// TestSummary is the entry point of SummaryTestSuite.
func TestSummary(t *testing.T) {
	suite.Run(t, &SummaryTestSuite{})
}

// SetupSuite prepares testing env for the testing suite.
func (suite *SummaryTestSuite) SetupSuite() {
	rp := vuln.Report{
		GeneratedAt: time.Now().UTC().String(),
		Scanner: &v1.Scanner{
			Name:    "Trivy",
			Vendor:  "Harbor",
			Version: "0.1.0",
		},
		Severity: vuln.High,
		Vulnerabilities: []*vuln.VulnerabilityItem{
			{
				ID:          "2019-0980-0909",
				Package:     "dpkg",
				Version:     "0.9.1",
				FixVersion:  "0.9.2",
				Severity:    vuln.High,
				Description: "mock one",
				Links:       []string{"https://vuln1.com"},
			},
			{
				ID:          "2019-0980-1010",
				Package:     "dpkg",
				Version:     "5.0.1",
				FixVersion:  "5.0.2",
				Severity:    vuln.Medium,
				Description: "mock two",
				Links:       []string{"https://vuln2.com"},
			},
		},
	}

	jsonData, err := json.Marshal(rp)
	require.NoError(suite.T(), err)

	suite.r = &scan.Report{
		ID:               1,
		UUID:             "r-uuid-001",
		Digest:           "digest-code",
		RegistrationUUID: "reg-uuid-001",
		MimeType:         v1.MimeTypeNativeReport,
		Status:           "Success",
		Report:           string(jsonData),
	}
}

// TestSummaryGenerateSummaryNoOptions ...
func (suite *SummaryTestSuite) TestSummaryGenerateSummaryNoOptions() {
	summaries, err := GenerateSummary(suite.r)
	require.NoError(suite.T(), err)
	require.NotNil(suite.T(), summaries)

	nativeSummary, ok := summaries.(*vuln.NativeReportSummary)
	require.Equal(suite.T(), true, ok)

	suite.Equal(vuln.High, nativeSummary.Severity)
	suite.Nil(nativeSummary.CVEBypassed)
	suite.Equal(2, nativeSummary.Summary.Total)

	suite.Equal("Trivy", nativeSummary.Scanner.Name)
	suite.Equal("Harbor", nativeSummary.Scanner.Vendor)
	suite.Equal("0.1.0", nativeSummary.Scanner.Version)
}

// TestSummaryGenerateSummaryWrongMime ...
func (suite *SummaryTestSuite) TestSummaryGenerateSummaryWrongMime() {
	suite.r.MimeType = "wrong-mime"
	defer func() {
		suite.r.MimeType = v1.MimeTypeNativeReport
	}()

	_, err := GenerateSummary(suite.r)
	require.Error(suite.T(), err)
}
