from __future__ import absolute_import
import unittest

from testutils import harbor_server, created_project, created_user, \
                      TEARDOWN, ADMIN_CLIENT,suppress_urllib3_warning
from library.repository import Repository
from library.repository import push_self_build_image_to_project, push_local_image_to_project
from library.system import System

class TestProjects(unittest.TestCase):
    @suppress_urllib3_warning
    def setUp(cls):
        cls.repo = Repository()
        cls.system = System()

    @unittest.skipIf(TEARDOWN == False, "Test data won't be erased.")
    def tearDown(self):
        print("Case completed")

    def testProjectQuota(self):
        """
        Test case:
            Project Quota
        Test step and expected result:
            1. Create a new user(UA);
            2. Create a new private project(PA) by user(UA);
            3. Add user(UA) as a member of project(PA) with project-admin role;
            4. Check quota is 0
            5. Push an image to project(PA) by user(UA), then check the project quota usage;
            6. Check quota change
            7. Push the image with another tag to project(PA) by user(UA)
            8. Check quota not changed
            9. Delete repository(RA) by user(UA);
            10. Delete image, the quota should be changed to 0.
        Tear down:
            1. Delete repository(RA) by user(UA);
            2. Delete project(PA);
            3. Delete user(UA);
        """
        user_001_password = "Aa123456"

        #1. Create a new user(UA);
        with created_user(user_001_password) as (user_id, user_name):
            #2. Create a new private project(PA) by user(UA);
            #3. Add user(UA) as a member of project(PA) with project-admin role;
            with created_project(metadata={"public": "false"}, user_id=user_id) as (project_id, project_name):
                #4. Check quota is 0
                quota = self.system.get_project_quota("project", project_id, **ADMIN_CLIENT)
                self.assertEqual(quota[0].used["storage"], 0)
                #5. Push an image to project(PA) by user(UA), then check the project quota usage;
                image, tag = "goharbor/alpine", "3.10"
                image_alias_name = "_alias"
                push_self_build_image_to_project(project_name, harbor_server, user_name, user_001_password, image, tag)

                #6. Get project quota
                quota = self.system.get_project_quota("project", project_id, **ADMIN_CLIENT)
                quota_size = quota[0].used["storage"]
                self.assertNotEqual(quota_size, 0)

                #7. Push the image with another tag to project(PA) by user(UA), the check the project quota usage;
                push_local_image_to_project(harbor_server, user_name, user_001_password, "{}/{}/{}".format(harbor_server, project_name, image), tag, "{}/{}/{}{}".format(harbor_server, project_name, image, image_alias_name), tag + image_alias_name)

                #8. Get project quota
                quota = self.system.get_project_quota("project", project_id, **ADMIN_CLIENT)
                self.assertEqual(quota[0].used["storage"], quota_size)

                #9. Delete repository(RA) by user(UA);
                self.repo.delete_repository(project_name, "goharbor%2Falpine", **ADMIN_CLIENT)
                self.repo.delete_repository(project_name, "goharbor%2Falpine"+image_alias_name, **ADMIN_CLIENT)

                #10. Quota should be 0
                quota = self.system.get_project_quota("project", project_id, **ADMIN_CLIENT)
                self.assertEqual(quota[0].used["storage"], 0)

if __name__ == '__main__':
    unittest.main()