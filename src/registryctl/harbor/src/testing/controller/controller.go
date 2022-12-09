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

package controller

//go:generate mockery --case snake --dir ../../controller/artifact --name Controller --output ./artifact --outpkg artifact
//go:generate mockery --case snake --dir ../../controller/blob --name Controller --output ./blob --outpkg blob
//go:generate mockery --case snake --dir ../../controller/project --name Controller --output ./project --outpkg project
//go:generate mockery --case snake --dir ../../controller/quota --name Controller --output ./quota --outpkg quota
//go:generate mockery --case snake --dir ../../controller/scan --name Controller --output ./scan --outpkg scan
//go:generate mockery --case snake --dir ../../controller/scan --name Checker --output ./scan --outpkg scan
//go:generate mockery --case snake --dir ../../controller/scanner --name Controller --output ./scanner --outpkg scanner
//go:generate mockery --case snake --dir ../../controller/replication --name Controller --output ./replication --outpkg replication
//go:generate mockery --case snake --dir ../../controller/robot --name Controller --output ./robot --outpkg robot
//go:generate mockery --case snake --dir ../../controller/proxy --name RemoteInterface --output ./proxy --outpkg proxy
//go:generate mockery --case snake --dir ../../controller/retention --name Controller --output ./retention --outpkg retention
//go:generate mockery --case snake --dir ../../controller/config --name Controller --output ./config --outpkg config
//go:generate mockery --case snake --dir ../../controller/user --name Controller --output ./user --outpkg user
//go:generate mockery --case snake --dir ../../controller/repository --name Controller --output ./repository --outpkg repository
//go:generate mockery --case snake --dir ../../controller/purge --name Controller --output ./purge --outpkg purge
//go:generate mockery --case snake --dir ../../controller/jobservice --name SchedulerController --output ./jobservice --outpkg jobservice
//go:generate mockery --case snake --dir ../../controller/systemartifact --name Controller --output ./systemartifact --outpkg systemartifact
//go:generate mockery --case snake --dir ../../controller/scandataexport --name Controller --output ./scandataexport --outpkg scandataexport
