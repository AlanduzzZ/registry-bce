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
import { NgModule } from '@angular/core';
import { SharedModule } from '../../../shared/shared.module';
import { ProjectsComponent } from './projects.component';
import { ListProjectComponent } from './list-project/list-project.component';
import { CreateProjectComponent } from './create-project/create-project.component';
import { RouterModule, Routes } from '@angular/router';
import { StatisticsPanelComponent } from './statictics/statistics-panel.component';
import { ExportCveComponent } from './list-project/export-cve/export-cve.component';

const routes: Routes = [
    {
        path: '',
        component: ProjectsComponent,
    },
];

@NgModule({
    imports: [SharedModule, RouterModule.forChild(routes)],
    declarations: [
        ProjectsComponent,
        ListProjectComponent,
        CreateProjectComponent,
        StatisticsPanelComponent,
        ExportCveComponent,
    ],
    providers: [],
})
export class ProjectsModule {}
