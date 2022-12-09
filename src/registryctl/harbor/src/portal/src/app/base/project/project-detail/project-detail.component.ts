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
import {
    Component,
    OnInit,
    HostListener,
    AfterViewInit,
    OnDestroy,
    ChangeDetectorRef,
} from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';
import { Project } from '../project';
import { SessionService } from '../../../shared/services/session.service';
import { AppConfigService } from '../../../services/app-config.service';
import { forkJoin, Observable, Subject, Subscription } from 'rxjs';
import {
    UserPermissionService,
    USERSTATICPERMISSION,
} from '../../../shared/services';
import { ErrorHandler } from '../../../shared/units/error-handler';
import { debounceTime } from 'rxjs/operators';
import { DOWN, SHOW_ELLIPSIS_WIDTH, UP } from './project-detail.const';
import { ProjectService } from '../../../../../ng-swagger-gen/services/project.service';
import { ProjectSummaryQuota } from '../../../../../ng-swagger-gen/models/project-summary-quota';
import {
    QUOTA_DANGER_COEFFICIENT,
    QUOTA_WARNING_COEFFICIENT,
    QuotaUnits,
} from '../../../shared/entities/shared.const';
import {
    clone,
    GetIntegerAndUnit,
    getSizeNumber,
    getSizeUnit,
} from '../../../shared/units/utils';
import {
    EventService,
    HarborEvent,
} from '../../../services/event-service/event.service';

@Component({
    selector: 'project-detail',
    templateUrl: 'project-detail.component.html',
    styleUrls: ['project-detail.component.scss'],
})
export class ProjectDetailComponent
    implements OnInit, AfterViewInit, OnDestroy
{
    hasSignedIn: boolean;
    currentProject: Project;

    isMember: boolean;
    roleName: string;
    projectId: number;
    hasProjectReadPermission: boolean;
    hasHelmChartsListPermission: boolean;
    hasRepositoryListPermission: boolean;
    hasMemberListPermission: boolean;
    hasLabelListPermission: boolean;
    hasLabelCreatePermission: boolean;
    hasLogListPermission: boolean;
    hasConfigurationListPermission: boolean;
    hasRobotListPermission: boolean;
    hasTagRetentionPermission: boolean;
    hasTagImmutablePermission: boolean;
    hasWebhookListPermission: boolean;
    hasScannerReadPermission: boolean;
    hasP2pProviderReadPermission: boolean;
    hasQuotaReadPermission: boolean = false;
    tabLinkNavList = [
        {
            linkName: 'summary',
            tabLinkInOverflow: false,
            showTabName: 'PROJECT_DETAIL.SUMMARY',
            permissions: () => this.hasProjectReadPermission,
        },
        {
            linkName: 'repositories',
            tabLinkInOverflow: false,
            showTabName: 'PROJECT_DETAIL.REPOSITORIES',
            permissions: () => this.hasRepositoryListPermission,
        },
        {
            linkName: 'helm-charts',
            tabLinkInOverflow: false,
            showTabName: 'PROJECT_DETAIL.HELMCHART',
            permissions: () =>
                this.withHelmChart && this.hasHelmChartsListPermission,
        },
        {
            linkName: 'members',
            tabLinkInOverflow: false,
            showTabName: 'PROJECT_DETAIL.USERS',
            permissions: () => this.hasMemberListPermission,
        },
        {
            linkName: 'labels',
            tabLinkInOverflow: false,
            showTabName: 'PROJECT_DETAIL.LABELS',
            permissions: () =>
                this.hasLabelListPermission &&
                this.hasLabelCreatePermission &&
                !this.withAdmiral,
        },
        {
            linkName: 'scanner',
            tabLinkInOverflow: false,
            showTabName: 'SCANNER.SCANNER',
            permissions: () => this.hasScannerReadPermission,
        },
        {
            linkName: 'p2p-provider',
            tabLinkInOverflow: false,
            showTabName: 'P2P_PROVIDER.P2P_PROVIDER',
            permissions: () => this.hasP2pProviderReadPermission,
        },
        {
            linkName: 'tag-strategy',
            tabLinkInOverflow: false,
            showTabName: 'PROJECT_DETAIL.POLICY',
            permissions: () =>
                this.hasTagRetentionPermission ||
                this.hasTagImmutablePermission,
        },
        {
            linkName: 'robot-account',
            tabLinkInOverflow: false,
            showTabName: 'PROJECT_DETAIL.ROBOT_ACCOUNTS',
            permissions: () => this.hasRobotListPermission,
        },
        {
            linkName: 'webhook',
            tabLinkInOverflow: false,
            showTabName: 'PROJECT_DETAIL.WEBHOOKS',
            permissions: () => this.hasWebhookListPermission,
        },
        {
            linkName: 'logs',
            tabLinkInOverflow: false,
            showTabName: 'PROJECT_DETAIL.LOGS',
            permissions: () => this.hasLogListPermission,
        },
        {
            linkName: 'configs',
            tabLinkInOverflow: false,
            showTabName: 'PROJECT_DETAIL.CONFIG',
            permissions: () =>
                this.isSessionValid && this.hasConfigurationListPermission,
        },
    ];
    previousWindowWidth: number;
    private _subject = new Subject<string>();
    private _subscription: Subscription;
    isProxyCacheProject: boolean = false;
    projectQuota: ProjectSummaryQuota;
    quotaDangerCoefficient: number = QUOTA_DANGER_COEFFICIENT;
    quotaWarningCoefficient: number = QUOTA_WARNING_COEFFICIENT;
    eventSub: Subscription;
    constructor(
        private projectService: ProjectService,
        private route: ActivatedRoute,
        private router: Router,
        private sessionService: SessionService,
        private appConfigService: AppConfigService,
        private userPermissionService: UserPermissionService,
        private errorHandler: ErrorHandler,
        private cdf: ChangeDetectorRef,
        private event: EventService
    ) {
        this.hasSignedIn = this.sessionService.getCurrentUser() !== null;
        this.route.data.subscribe(data => {
            this.currentProject = <Project>data['projectResolver'];
            if (this.currentProject.registry_id) {
                this.isProxyCacheProject = true;
            }
            this.isMember = this.currentProject.is_member;
            this.roleName = this.currentProject.role_name;
        });
    }
    ngOnInit() {
        this.projectId = this.route.snapshot.params['id'];
        this.getPermissionsList(this.projectId);
        if (!this._subscription) {
            this._subscription = this._subject
                .pipe(debounceTime(100))
                .subscribe(type => {
                    if (type === DOWN) {
                        this.resetTabsForDownSize();
                    } else {
                        this.resetTabsForUpSize();
                    }
                });
        }
        if (!this.eventSub) {
            this.eventSub = this.event.subscribe(
                HarborEvent.REFRESH_PROJECT_INFO,
                () => {
                    this.refreshProjectInfo();
                }
            );
        }
    }

    ngAfterViewInit() {
        this.previousWindowWidth = window.innerWidth;
        setTimeout(() => {
            this.resetTabsForDownSize();
        }, 0);
    }
    ngOnDestroy() {
        if (this._subscription) {
            this._subscription.unsubscribe();
            this._subscription = null;
        }
        if (this.eventSub) {
            this.eventSub.unsubscribe();
            this.eventSub = null;
        }
    }
    getPermissionsList(projectId: number): void {
        let permissionsList: Array<Observable<boolean>> = [];
        permissionsList.push(
            this.userPermissionService.getPermission(
                projectId,
                USERSTATICPERMISSION.PROJECT.KEY,
                USERSTATICPERMISSION.PROJECT.VALUE.READ
            )
        );
        permissionsList.push(
            this.userPermissionService.getPermission(
                projectId,
                USERSTATICPERMISSION.LOG.KEY,
                USERSTATICPERMISSION.LOG.VALUE.LIST
            )
        );
        permissionsList.push(
            this.userPermissionService.getPermission(
                projectId,
                USERSTATICPERMISSION.CONFIGURATION.KEY,
                USERSTATICPERMISSION.CONFIGURATION.VALUE.READ
            )
        );
        permissionsList.push(
            this.userPermissionService.getPermission(
                projectId,
                USERSTATICPERMISSION.MEMBER.KEY,
                USERSTATICPERMISSION.MEMBER.VALUE.LIST
            )
        );
        permissionsList.push(
            this.userPermissionService.getPermission(
                projectId,
                USERSTATICPERMISSION.LABEL.KEY,
                USERSTATICPERMISSION.LABEL.VALUE.LIST
            )
        );
        permissionsList.push(
            this.userPermissionService.getPermission(
                projectId,
                USERSTATICPERMISSION.REPOSITORY.KEY,
                USERSTATICPERMISSION.REPOSITORY.VALUE.LIST
            )
        );
        permissionsList.push(
            this.userPermissionService.getPermission(
                projectId,
                USERSTATICPERMISSION.HELM_CHART.KEY,
                USERSTATICPERMISSION.HELM_CHART.VALUE.LIST
            )
        );
        permissionsList.push(
            this.userPermissionService.getPermission(
                projectId,
                USERSTATICPERMISSION.ROBOT.KEY,
                USERSTATICPERMISSION.ROBOT.VALUE.LIST
            )
        );
        permissionsList.push(
            this.userPermissionService.getPermission(
                projectId,
                USERSTATICPERMISSION.LABEL.KEY,
                USERSTATICPERMISSION.LABEL.VALUE.CREATE
            )
        );
        permissionsList.push(
            this.userPermissionService.getPermission(
                projectId,
                USERSTATICPERMISSION.TAG_RETENTION.KEY,
                USERSTATICPERMISSION.TAG_RETENTION.VALUE.READ
            )
        );
        permissionsList.push(
            this.userPermissionService.getPermission(
                projectId,
                USERSTATICPERMISSION.IMMUTABLE_TAG.KEY,
                USERSTATICPERMISSION.IMMUTABLE_TAG.VALUE.LIST
            )
        );
        permissionsList.push(
            this.userPermissionService.getPermission(
                projectId,
                USERSTATICPERMISSION.WEBHOOK.KEY,
                USERSTATICPERMISSION.WEBHOOK.VALUE.LIST
            )
        );
        permissionsList.push(
            this.userPermissionService.getPermission(
                projectId,
                USERSTATICPERMISSION.SCANNER.KEY,
                USERSTATICPERMISSION.SCANNER.VALUE.READ
            )
        );
        permissionsList.push(
            this.userPermissionService.getPermission(
                projectId,
                USERSTATICPERMISSION.P2P_PROVIDER.KEY,
                USERSTATICPERMISSION.P2P_PROVIDER.VALUE.READ
            )
        );
        permissionsList.push(
            this.userPermissionService.getPermission(
                projectId,
                USERSTATICPERMISSION.QUOTA.KEY,
                USERSTATICPERMISSION.QUOTA.VALUE.READ
            )
        );

        forkJoin(permissionsList).subscribe(
            Rules => {
                [
                    this.hasProjectReadPermission,
                    this.hasLogListPermission,
                    this.hasConfigurationListPermission,
                    this.hasMemberListPermission,
                    this.hasLabelListPermission,
                    this.hasRepositoryListPermission,
                    this.hasHelmChartsListPermission,
                    this.hasRobotListPermission,
                    this.hasLabelCreatePermission,
                    this.hasTagRetentionPermission,
                    this.hasTagImmutablePermission,
                    this.hasWebhookListPermission,
                    this.hasScannerReadPermission,
                    this.hasP2pProviderReadPermission,
                    this.hasQuotaReadPermission,
                ] = Rules;
                if (this.hasQuotaReadPermission) {
                    this.getQuotaInfo();
                }
            },
            error => this.errorHandler.error(error)
        );
    }
    getQuotaInfo() {
        this.projectService
            .getProjectSummary({
                projectNameOrId: this.projectId.toString(),
            })
            .subscribe(
                res => {
                    if (res && res.quota) {
                        this.projectQuota = res.quota;
                    }
                },
                error => {
                    this.errorHandler.error(error);
                }
            );
    }

    public get isSessionValid(): boolean {
        return this.sessionService.getCurrentUser() != null;
    }

    public get withAdmiral(): boolean {
        return this.appConfigService.getConfig().with_admiral;
    }

    public get withHelmChart(): boolean {
        return this.appConfigService.getConfig().with_chartmuseum;
    }

    backToProject(): void {
        if (window.sessionStorage) {
            window.sessionStorage.setItem('fromDetails', 'true');
        }
        this.router.navigate(['/harbor', 'projects']);
    }
    isDefaultTab(tab, index) {
        return (
            this.route.snapshot.children[0].routeConfig.path !== tab.linkName &&
            index === 0
        );
    }

    isTabLinkInOverFlow() {
        return this.tabLinkNavList.some(tab => {
            return (
                tab.tabLinkInOverflow &&
                this.route.snapshot.children[0].routeConfig.path ===
                    tab.linkName
            );
        });
    }

    @HostListener('window:resize', ['$event'])
    onResize(event) {
        if (this.previousWindowWidth) {
            // down size
            if (this.previousWindowWidth > event.target.innerWidth) {
                this._subject.next(DOWN);
            } else {
                // up size
                this._subject.next(UP);
            }
        }
        this.previousWindowWidth = event.target.innerWidth;
    }

    resetTabsForDownSize(): void {
        this.tabLinkNavList
            .filter(item => !item.tabLinkInOverflow)
            .forEach((item, index) => {
                const tabEle: HTMLElement = document.getElementById(
                    'project-' + item.linkName
                );
                // strengthen code
                if (tabEle && tabEle.getBoundingClientRect) {
                    const right: number =
                        window.innerWidth -
                        document
                            .getElementById('project-' + item.linkName)
                            .getBoundingClientRect().right;
                    if (right < SHOW_ELLIPSIS_WIDTH) {
                        this.tabLinkNavList[index].tabLinkInOverflow = true;
                    }
                }
            });
    }
    resetTabsForUpSize() {
        // 1.Set tabLinkInOverflow to false for all tabs(show all tabs)
        for (let i = 0; i < this.tabLinkNavList.length; i++) {
            this.tabLinkNavList[i].tabLinkInOverflow = false;
        }
        // 2.Manually  detect changes to rerender dom
        this.cdf.detectChanges();
        // 3. Hide overflowed tabs
        this.resetTabsForDownSize();
    }
    getIntegerAndUnit(hardValue, usedValue) {
        return GetIntegerAndUnit(
            hardValue,
            clone(QuotaUnits),
            usedValue,
            clone(QuotaUnits)
        );
    }

    getSizeNumber(): number | string {
        if (
            this.projectQuota &&
            this.projectQuota.used &&
            this.projectQuota.used.storage
        ) {
            return getSizeNumber(this.projectQuota.used.storage);
        }
        return 0;
    }
    getSizeUnit(): number | string {
        if (this.projectQuota) {
            return getSizeUnit(this.projectQuota.used.storage);
        }
        return null;
    }
    refreshProjectInfo() {
        this.getQuotaInfo();
        this.getProject();
    }
    getProject() {
        this.projectService
            .getProject({
                projectNameOrId: this.projectId.toString(),
            })
            .subscribe(res => {
                if (res) {
                    this.currentProject = (<any>res) as Project;
                }
            });
    }
}
