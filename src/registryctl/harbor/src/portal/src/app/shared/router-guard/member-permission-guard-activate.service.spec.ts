import { TestBed, inject } from '@angular/core/testing';
import { RouterTestingModule } from '@angular/router/testing';
import { MemberPermissionGuard } from './member-permission-guard-activate.service';
import { of } from 'rxjs';
import { ErrorHandler } from '../units/error-handler';
import { UserPermissionService } from '../services';

describe('MemberPermissionGuardActivateServiceGuard', () => {
    const fakeUserPermissionService = {
        getPermission() {
            return of(true);
        },
    };
    beforeEach(() => {
        TestBed.configureTestingModule({
            imports: [RouterTestingModule],
            providers: [
                ErrorHandler,
                MemberPermissionGuard,
                {
                    provide: UserPermissionService,
                    useValue: fakeUserPermissionService,
                },
            ],
        });
    });

    it('should ...', inject(
        [MemberPermissionGuard],
        (guard: MemberPermissionGuard) => {
            expect(guard).toBeTruthy();
        }
    ));
});
