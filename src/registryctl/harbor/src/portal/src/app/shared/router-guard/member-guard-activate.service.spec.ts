import { TestBed, inject } from '@angular/core/testing';
import { RouterTestingModule } from '@angular/router/testing';
import { SessionService } from '../services/session.service';
import { MemberGuard } from './member-guard-activate.service';
import { ProjectService } from '../services';

describe('MemberGuard', () => {
    const fakeSessionService = null;
    const fakeProjectService = null;

    beforeEach(() => {
        TestBed.configureTestingModule({
            imports: [RouterTestingModule],
            providers: [
                MemberGuard,
                { provide: SessionService, useValue: fakeSessionService },
                { provide: ProjectService, useValue: fakeProjectService },
            ],
        });
    });

    it('should ...', inject([MemberGuard], (guard: MemberGuard) => {
        expect(guard).toBeTruthy();
    }));
});
