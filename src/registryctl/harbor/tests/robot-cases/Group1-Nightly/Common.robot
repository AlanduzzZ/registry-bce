# Copyright Project Harbor Authors
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#	http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License

*** Settings ***
Documentation  Harbor BATs
Library  ../../apitests/python/testutils.py
Library  ../../apitests/python/library/oras.py
Library  ../../apitests/python/library/singularity.py
Resource  ../../resources/Util.robot
Default Tags  Nightly

*** Variables ***
${HARBOR_URL}  https://${ip}
${SSH_USER}  root
${HARBOR_ADMIN}  admin

*** Test Cases ***
Test Case - Sign With Admin
    [tags]  admin
    Init Chrome Driver
    Sign In Harbor  ${HARBOR_URL}  ${HARBOR_ADMIN}  ${HARBOR_PASSWORD}
    Close Browser

Test Case - Push ORAS and Display
    [Tags]  push_oras
    Init Chrome Driver
    ${d}=    Get Current Date    result_format=%m%s

    Sign In Harbor  ${HARBOR_URL}  user010  Test1@34
    Create An New Project And Go Into Project  test${d}

    ${repo_name}=  Set Variable  hello-oras-artifact
    ${tag}=  Set Variable  1.0.0
    Retry Keyword N Times When Error  5  Oras Push  ${ip}  user010  Test1@34  test${d}  ${repo_name}  ${tag}

    Go Into Project  test${d}
    Wait Until Page Contains  test${d}/${repo_name}

    Go Into Repo  test${d}/${repo_name}
    Wait Until Page Contains  ${tag}
    Close Browser

## TODO: uncomment it once #14470 fixed
# Test Case - Push SIF and Display
#     [Tags]  push_sif
#     Init Chrome Driver
#     ${d}=    Get Current Date    result_format=%m%s
#     ${user}=  Set Variable  user010
#     ${pwd}=  Set Variable  Test1@34

#     Sign In Harbor  ${HARBOR_URL}  ${user}  ${pwd}
#     Create An New Project And Go Into Project  test${d}

#     ${repo_name}=  Set Variable  busybox
#     ${tag}=  Set Variable  1.28
#     Retry Keyword N Times When Error  5  Push Singularity To Harbor  library:  library/default/  ${ip}  ${user}  ${pwd}  test${d}  ${repo_name}  ${tag}

#     Go Into Project  test${d}
#     Wait Until Page Contains  test${d}/${repo_name}

#     Go Into Repo  test${d}/${repo_name}
#     Wait Until Page Contains  ${tag}
#     Close Browser

Test Case - Push CNAB Bundle and Display
    [Tags]  push_cnab
    Init Chrome Driver
    ${d}=    Get Current Date    result_format=%m%s
    ${index1_image1}=  Set Variable  busybox
    ${index1_image2}=  Set Variable  alpine
    ${index2_image1}=  Set Variable  hello-world
    ${index2_image2}=  Set Variable  redis
    ${user}=  Set Variable  user010
    ${pwd}=  Set Variable  Test1@34
    Sign In Harbor  ${HARBOR_URL}  ${user}  ${pwd}
    Create An New Project And Go Into Project  test${d}
    ${index1}  ${index2}=  Prepare Cnab Push Test Data  ${ip}  ${user}  ${pwd}  test${d}  ${index1_image1}  ${index1_image2}  ${index2_image1}  ${index2_image2}
    ${target}=  Set Variable  ${ip}/test${d}/cnab${d}:cnab_tag${d}
    Retry Keyword N Times When Error  5  CNAB Push Bundle  ${ip}  ${user}  ${pwd}  ${target}  ./tests/robot-cases/Group0-Util/bundle.json  ${ip}  test${d}  ${index1}  ${index2}

    Go Into Project  test${d}
    Wait Until Page Contains  test${d}/cnab${d}

    Go Into Repo  test${d}/cnab${d}
    Wait Until Page Contains  cnab_tag${d}
    Go Into Project  test${d}
    Wait Until Page Contains  test${d}/cnab${d}
    Go Into Repo  test${d}/cnab${d}
    Go Into Index And Contain Artifacts  cnab_tag${d}  total_artifact_count=3  archive_count=2
    Close Browser

Test Case - Create An New Project
    Init Chrome Driver
    ${d}=    Get Current Date    result_format=%m%s
    Sign In Harbor  ${HARBOR_URL}  ${HARBOR_ADMIN}  ${HARBOR_PASSWORD}
    Create An New Project And Go Into Project  test${d}
    Close Browser

Test Case - Delete A Project
    Init Chrome Driver
    Sign In Harbor  ${HARBOR_URL}  ${HARBOR_ADMIN}  ${HARBOR_PASSWORD}
    Delete A Project Without Sign In Harbor
    Close Browser

Test Case - Repo Size
    Init Chrome Driver
    ${d}=  Get Current Date    result_format=%m%s
    Sign In Harbor  ${HARBOR_URL}  ${HARBOR_ADMIN}  ${HARBOR_PASSWORD}
    Push Image With Tag  ${ip}  ${HARBOR_ADMIN}  ${HARBOR_PASSWORD}  library  alpine  2.6  2.6
    Go Into Project  library
    Go Into Repo  alpine
    Wait Until Page Contains  1.92MiB
    Close Browser

Test Case - Staticsinfo
    Init Chrome Driver
    ${d}=  Get Current Date    result_format=%m%s
    Sign In Harbor  ${HARBOR_URL}  ${HARBOR_ADMIN}  ${HARBOR_PASSWORD}
    ${element}=  Set Variable  ${project_statistics_private_repository_icon}
    Wait Until Element Is Visible  ${element}
    ${privaterepocount1}=  Get Statics Private Repo
    ${privateprojcount1}=  Get Statics Private Project
    ${publicrepocount1}=  Get Statics Public Repo
    ${publicprojcount1}=  Get Statics Public Project
    ${totalrepocount1}=  Get Statics Total Repo
    ${totalprojcount1}=  Get Statics Total Project
    Create An New Project And Go Into Project  private${d}
    Create An New Project And Go Into Project  public${d}  true
    Push Image  ${ip}  ${HARBOR_ADMIN}  ${HARBOR_PASSWORD}  private${d}  hello-world
    Push Image  ${ip}  ${HARBOR_ADMIN}  ${HARBOR_PASSWORD}  public${d}  hello-world
    Reload Page
    ${privateprojcount}=  evaluate  ${privateprojcount1}+1
    ${privaterepocount}=  evaluate  ${privaterepocount1}+1
    ${publicprojcount}=  evaluate  ${publicprojcount1}+1
    ${publicrepocount}=  evaluate  ${publicrepocount1}+1
    ${totalrepocount}=  evaluate  ${totalrepocount1}+2
    ${totalprojcount}=  evaluate  ${totalprojcount1}+2
    Navigate To Projects
    Wait Until Element Is Visible  ${element}
    ${privaterepocountStr}=  Convert To String  ${privaterepocount}
    Wait Until Element Contains  ${element}  ${privaterepocountStr}
    ${privaterepocount2}=  Get Statics Private Repo
    ${privateprojcount2}=  get statics private project
    ${publicrepocount2}=  get statics public repo
    ${publicprojcount2}=  get statics public project
    ${totalrepocount2}=  get statics total repo
    ${totalprojcount2}=  get statics total project
    Should Be Equal As Integers  ${privateprojcount2}  ${privateprojcount}
    Should be equal as integers  ${privaterepocount2}  ${privaterepocount}
    Should be equal as integers  ${publicprojcount2}  ${publicprojcount}
    Should be equal as integers  ${publicrepocount2}  ${publicrepocount}
    Should be equal as integers  ${totalprojcount2}  ${totalprojcount}
    Should be equal as integers  ${totalrepocount2}  ${totalrepocount}
    Close Browser

Test Case - Push Image
    Init Chrome Driver
    ${d}=    Get Current Date    result_format=%m%s
    Sign In Harbor  ${HARBOR_URL}  ${HARBOR_ADMIN}  ${HARBOR_PASSWORD}
    Create An New Project And Go Into Project  test${d}

    Push image  ${ip}  ${HARBOR_ADMIN}  ${HARBOR_PASSWORD}  test${d}  hello-world:latest
    Go Into Project  test${d}
    Wait Until Page Contains  test${d}/hello-world
    Close Browser

Test Case - Project Level Policy Public
    Init Chrome Driver
    ${d}=  Get Current Date    result_format=%m%s
    Sign In Harbor  ${HARBOR_URL}  ${HARBOR_ADMIN}  ${HARBOR_PASSWORD}
    Create An New Project And Go Into Project  project${d}
    Goto Project Config
    Click Project Public
    Save Project Config
    # Verify
    Public Should Be Selected
    # Project${d}  default should be private
    # Here logout and login to try avoid a bug only in autotest
    Logout Harbor
    Sign In Harbor  ${HARBOR_URL}  ${HARBOR_ADMIN}  ${HARBOR_PASSWORD}
    Retry Double Keywords When Error  Filter Project  project${d}  Project Should Be Public  project${d}
    Close Browser

Test Case - Verify Download Ca Link
    Init Chrome Driver
    Sign In Harbor  ${HARBOR_URL}  ${HARBOR_ADMIN}  ${HARBOR_PASSWORD}
    Switch To System Settings
    Page Should Contain  Registry Root Certificate
    Close Browser

Test Case - Edit Email Settings
    Init Chrome Driver
    Sign In Harbor  ${HARBOR_URL}  ${HARBOR_ADMIN}  ${HARBOR_PASSWORD}

    Switch To Email
    Config Email

    Logout Harbor
    Sign In Harbor  ${HARBOR_URL}  ${HARBOR_ADMIN}  ${HARBOR_PASSWORD}

    Switch To Email
    Verify Email

    Close Browser

Test Case - Edit Token Expire
    Init Chrome Driver
    Sign In Harbor  ${HARBOR_URL}  ${HARBOR_ADMIN}  ${HARBOR_PASSWORD}
    Switch To System Settings
    Modify Token Expiration  20
    Logout Harbor

    Sign In Harbor  ${HARBOR_URL}  ${HARBOR_ADMIN}  ${HARBOR_PASSWORD}
    Switch To System Settings
    Token Must Be Match  20

    #reset to default
    Modify Token Expiration  30
    Close Browser

Test Case - Create A New Labels
    Init Chrome Driver
    ${d}=    Get Current Date
    Sign In Harbor  ${HARBOR_URL}  ${HARBOR_ADMIN}  ${HARBOR_PASSWORD}
    Switch To System Labels
    Create New Labels  label_${d}
    Close Browser

Test Case - Update Label
   Init Chrome Driver
   ${d}=    Get Current Date

   Sign In Harbor  ${HARBOR_URL}  ${HARBOR_ADMIN}  ${HARBOR_PASSWORD}
   Switch To System Labels
   Create New Labels  label_${d}
   Sleep  3
   ${d1}=    Get Current Date
   Update A Label  label_${d}
   Close Browser

Test Case - Delete Label
    Init Chrome Driver
    ${d}=    Get Current Date
    Sign In Harbor  ${HARBOR_URL}  ${HARBOR_ADMIN}  ${HARBOR_PASSWORD}
    Switch To System Labels
    Create New Labels  label_${d}
    Sleep  3
    Delete A Label  label_${d}
    Close Browser

Test Case - User View Projects
    Init Chrome Driver
    ${d}=    Get Current Date    result_format=%m%s
    Sign In Harbor  ${HARBOR_URL}  user001  Test1@34
    Create An New Project And Go Into Project  test${d}1
    Create An New Project And Go Into Project  test${d}2
    Create An New Project And Go Into Project  test${d}3
    Switch To Logs
    Wait Until Page Contains  test${d}1
    Wait Until Page Contains  test${d}2
    Wait Until Page Contains  test${d}3
    Close Browser

Test Case - User View Logs
    [tags]  user_view_logs
    Init Chrome Driver
    ${d}=   Get Current Date    result_format=%m%s
    ${img}=    Set Variable    kong
    ${tag}=    Set Variable    latest
    ${user}=    Set Variable    user002
    ${pwd}=    Set Variable    Test1@34
    &{image_with_tag}=	 Create Dictionary  image=for_log_view  tag=base
    ${replication_image}=  Get From Dictionary  ${image_with_tag}  image
    ${replication_tag}=  Get From Dictionary  ${image_with_tag}  tag
    @{target_images}=  Create List  '&{image_with_tag}'

    Sign In Harbor  ${HARBOR_URL}  ${user}  ${pwd}
    Create An New Project And Go Into Project  project${d}
    Logout Harbor

    Body Of Replication Of Pull Images from Registry To Self   harbor  https://cicd.harbor.vmwarecna.net  ${null}  ${null}  nightly/${replication_image}  project${d}  N  Flatten 1 Level  @{target_images}

    Push image  ${ip}  ${user}  ${pwd}  project${d}  ${img}:${tag}
    Pull image  ${ip}  ${user}  ${pwd}  project${d}  ${replication_image}:${replication_tag}
    Close Browser

    Init Chrome Driver
    Sign In Harbor  ${HARBOR_URL}  ${user}  ${pwd}
    Go Into Project  project${d}
    Delete Repo  project${d}  ${replication_image}
    Delete Repo  project${d}  ${img}

    Sleep  3

    Go To Project Log
    Advanced Search Should Display

    Do Log Advanced Search
    Close Browser

Test Case - Manage Project Member
    Init Chrome Driver
    ${user}=    Set Variable    user004
    ${pwd}=    Set Variable    Test1@34
    Sign In Harbor  ${HARBOR_URL}  ${user}  ${pwd}
    Manage Project Member Without Sign In Harbor  ${user}  ${pwd}
    Close Browser

Test Case - Manage project publicity
    Body Of Manage project publicity

Test Case - Assign Sys Admin
    Init Chrome Driver
    ${d}=    Get Current Date    result_format=%m%s
    Sign In Harbor  ${HARBOR_URL}  user009  Test1@34
    Logout Harbor
    Sign In Harbor  ${HARBOR_URL}  ${HARBOR_ADMIN}  ${HARBOR_PASSWORD}
    Switch to User Tag
    Assign User Admin  user009
    Logout Harbor
    Sign In Harbor  ${HARBOR_URL}  user009  Test1@34
    Administration Tag Should Display
    Close Browser

Test Case - Edit Project Creation
    # Create normal user and login
    Init Chrome Driver
    ${d}=    Get Current Date    result_format=%m%s
    Sign In Harbor  ${HARBOR_URL}  user010  Test1@34
    Project Creation Should Display
    Logout Harbor

    Sleep  3
    Sign In Harbor  ${HARBOR_URL}  ${HARBOR_ADMIN}  ${HARBOR_PASSWORD}
    Set Pro Create Admin Only
    Logout Harbor

    Sign In Harbor  ${HARBOR_URL}  user010  Test1@34
    Project Creation Should Not Display
    Logout Harbor

    Sign In Harbor  ${HARBOR_URL}  ${HARBOR_ADMIN}  ${HARBOR_PASSWORD}
    Set Pro Create Every One
    Close browser

Test Case - Edit Repo Info
    Init Chrome Driver
    ${d}=  Get Current Date  result_format=%m%s

    Sign In Harbor  ${HARBOR_URL}  user011  Test1@34
    Create An New Project And Go Into Project  project${d}
    Push Image  ${ip}  user011  Test1@34  project${d}  hello-world
    Go Into Project  project${d}
    Go Into Repo  project${d}/hello-world
    Edit Repo Info
    Close Browser

Test Case - Delete Multi Project
    Init Chrome Driver
    ${d}=    Get Current Date    result_format=%m%s

    Sign In Harbor  ${HARBOR_URL}  user012  Test1@34
    Create An New Project And Go Into Project  projecta${d}
    Create An New Project And Go Into Project  projectb${d}
    Push Image  ${ip}  user012  Test1@34  projecta${d}  hello-world
    Navigate To Projects
    Filter Object  project
    Retry Wait Element Not Visible  //clr-datagrid/div/div[2]
    @{project_list}  Create List  projecta  projectb
    FOR  ${project}  IN  @{project_list}
        Retry Element Click  //clr-dg-row[contains(.,'${project}')]//label
    END
    Retry Element Click  ${project_action_xpath}
    Retry Element Click  ${project_delete_btn}
    Retry Element Click  ${repo_delete_on_card_view_btn}
    # Verify delete project with image should not be deleted directly
    Delete Fail  projecta${d}
    Delete Success  projectb${d}
    Close Browser

Test Case - Delete Multi Repo
    Init Chrome Driver
    ${d}=   Get Current Date    result_format=%m%s

    Sign In Harbor  ${HARBOR_URL}  user013  Test1@34
    Create An New Project And Go Into Project  project${d}
    Push Image  ${ip}  user013  Test1@34  project${d}  hello-world
    Push Image  ${ip}  user013  Test1@34  project${d}  busybox
    Sleep  2
    Go Into Project  project${d}
    @{repo_list}  Create List  hello-world  busybox
    Multi-delete Object  ${repo_delete_btn}  @{repo_list}
    # Verify
    Delete Success  hello-world  busybox
    Close Browser

Test Case - Delete Multi Artifacts
    Init Chrome Driver
    ${d}=   Get Current Date    result_format=%m%s
    Sign In Harbor  ${HARBOR_URL}  user014  Test1@34
    Create An New Project And Go Into Project  project${d}
    Push Image With Tag  ${ip}  user014  Test1@34  project${d}  redis  3.2.10-alpine  3.2.10-alpine
    Push Image With Tag  ${ip}  user014  Test1@34  project${d}  redis  4.0.7-alpine  4.0.7-alpine
    Go Into Project  project${d}
    Go Into Repo  redis
    @{tag_list}  Create List  3.2.10-alpine  4.0.7-alpine
    Multi-delete Artifact  @{tag_list}
    # Verify
    Delete Success  sha256:dd179737  sha256:28a85227
    Close Browser

Test Case - Delete Repo on CardView
    Init Chrome Driver
    ${d}=   Get Current Date  result_format=%m%s
    Sign In Harbor  ${HARBOR_URL}  user015  Test1@34
    Create An New Project And Go Into Project  project${d}
    Push Image  ${ip}  user015  Test1@34  project${d}  hello-world
    Push Image  ${ip}  user015  Test1@34  project${d}  busybox
    Go Into Project  project${d}
    Switch To CardView
    Delete Repo on CardView  busybox
    # Verify
    Delete Success  busybox
    Close Browser

Test Case - Delete Multi Member
    Init Chrome Driver
    ${d}=   Get Current Date    result_format=%m%s
    Sign In Harbor  ${HARBOR_URL}  user016  Test1@34
    Create An New Project And Go Into Project  project${d}
    Switch To Member
    Add Guest Member To Project  user017
    Add Guest Member To Project  user018
    Multi-delete Member  user017  user018
    Delete Success  user017  user018
    Close Browser

Test Case - Project Admin Operate Labels
    Init Chrome Driver
    ${d}=   Get Current Date    result_format=%m%s
    Sign In Harbor  ${HARBOR_URL}  user019  Test1@34
    Create An New Project And Go Into Project  project${d}
    Sleep  2
    # Add labels
    Switch To Project Label
    Create New Labels  label_${d}
    Sleep  2
    Update A Label  label_${d}
    Sleep  2
    Delete A Label  label_${d}
    Close Browser

Test Case - Project Admin Add Labels To Repo
    Init Chrome Driver
    ${d}=   Get Current Date    result_format=%m%s
    Sign In Harbor  ${HARBOR_URL}  user020  Test1@34
    Create An New Project And Go Into Project  project${d}
    Push Image With Tag  ${ip}  user020  Test1@34  project${d}  redis  3.2.10-alpine  3.2.10-alpine
    Push Image With Tag  ${ip}  user020  Test1@34  project${d}  redis  4.0.7-alpine  4.0.7-alpine
    Go Into Project  project${d}
    Sleep  2
    # Add labels
    Switch To Project Label
    Create New Labels  label111
    Create New Labels  label22
    Sleep  2
    Switch To Project Repo
    Go Into Repo  project${d}/redis
    Add Labels To Tag  3.2.10-alpine  label111
    Add Labels To Tag  4.0.7-alpine  label22
    Filter Labels In Tags  label111  label22
    Close Browser

Test Case - Developer Operate Labels
    Init Chrome Driver
    ${d}=   Get Current Date    result_format=%m%s

    Sign In Harbor  ${HARBOR_URL}  user021  Test1@34
    Create An New Project And Go Into Project  project${d}
    Logout Harbor

    Manage Project Member  user021  Test1@34  project${d}  user022  Add  ${false}
    Change User Role In Project  user021  Test1@34  project${d}  user022  Developer

    Sign In Harbor  ${HARBOR_URL}  user022  Test1@34
    Go Into Project  project${d}  has_image=${false}
    Sleep  3
    Retry Wait Until Page Not Contains Element  xpath=//a[contains(.,'Labels')]
    Close Browser

Test Case - Copy A Image
    Init Chrome Driver
    ${random_num1}=   Get Current Date    result_format=%m%s
    ${random_num2}=   Evaluate  str(random.randint(1000,9999))  modules=random

    Sign In Harbor  ${HARBOR_URL}  user028  Test1@34
    Create An New Project And Go Into Project  project${random_num1}${random_num2}
    Create An New Project And Go Into Project  project${random_num1}

    Sleep  1
    Push Image With Tag  ${ip}  user028  Test1@34  project${random_num1}  redis  ${image_tag}
    Sleep  1
    Go Into Repo  project${random_num1}/redis
    Copy Image  ${image_tag}  project${random_num1}${random_num2}  ${target_image_name}
    Retry Wait Element Not Visible  ${repo_retag_confirm_dlg}
    Navigate To Projects
    Go Into Project  project${random_num1}${random_num2}
    Sleep  1
    Page Should Contain  ${target_image_name}
    Go Into Repo  project${random_num1}${random_num2}/${target_image_name}
    Sleep  1
    Retry Wait Until Page Contains Element  xpath=${tag_value_xpath}
    Close Browser

Test Case - Copy A Image And Accessory
    [Tags]  copy_image_and_accessory
    Init Chrome Driver
    ${d}=  Get Current Date  result_format=%m%s
    ${source_project}=  Set Variable  source_project${d}
    ${target_project}=  Set Variable  target_project${d}
    ${user}=  Set Variable  user006
    ${pwd}=  Set Variable  Test1@34
    ${image}=  Set Variable  redis
    ${tag}=  Set Variable  latest

    Sign In Harbor  ${HARBOR_URL}  ${user}  ${pwd}
    Create An New Project And Go Into Project  ${target_project}
    Create An New Project And Go Into Project  ${source_project}

    Push Image With Tag  ${ip}  ${user}  ${pwd}  ${source_project}  ${image}  ${tag}
    Cosign Generate Key Pair
    Docker Login  ${ip}  ${user}  ${pwd}
    Cosign Sign  ${ip}/${source_project}/${image}:${tag}
    Docker Logout  ${ip}
    Retry Double Keywords When Error  Go Into Repo  ${source_project}/${image}  Should Be Signed By Cosign  ${tag}
    
    Copy Image  ${tag}  ${target_project}  ${image}
    Retry Wait Until Page Contains  Copy artifact successfully

    Retry Double Keywords When Error  Go Into Project  ${target_project}  Retry Wait Until Page Contains  ${image}
    Retry Double Keywords When Error  Go Into Repo  ${target_project}/${image}  Retry Wait Until Page Contains Element  //clr-dg-row[contains(.,${tag})]
    Should Be Signed By Cosign  ${tag}
    Close Browser

Test Case - Create An New Project With Quotas Set
    Init Chrome Driver
    ${d}=  Get Current Date  result_format=%m%s
    ${storage_quota}=  Set Variable  600
    ${storage_quota_unit}=  Set Variable  GiB
    Sign In Harbor  ${HARBOR_URL}  ${HARBOR_ADMIN}  ${HARBOR_PASSWORD}
    Create An New Project And Go Into Project    project${d}  storage_quota=${storage_quota}  storage_quota_unit=${storage_quota_unit}
    ${storage_quota_ret}=  Get Project Storage Quota Text From Project Quotas List  project${d}
    Should Be Equal As Strings  ${storage_quota_ret}  0Byte of ${storage_quota}${storage_quota_unit}
    Close Browser

Test Case - Project Storage Quotas Dispaly And Control
    Init Chrome Driver
    ${d}=  Get Current Date  result_format=%m%s
    ${storage_quota}=  Set Variable  350
    ${storage_quota_unit}=  Set Variable  MiB
    ${image_a}=  Set Variable  one_layer
    ${image_b}=  Set Variable  redis
    ${image_a_size}=    Set Variable   330.83MiB
    ${image_b_size}=    Set Variable   34.1\\dMiB
    ${image_a_ver}=  Set Variable  1.0
    ${image_b_ver}=  Set Variable  donotremove5.0
    Sign In Harbor  ${HARBOR_URL}  ${HARBOR_ADMIN}  ${HARBOR_PASSWORD}
    Create An New Project And Go Into Project  project${d}  storage_quota=${storage_quota}  storage_quota_unit=${storage_quota_unit}
    Push Image With Tag  ${ip}  ${HARBOR_ADMIN}  ${HARBOR_PASSWORD}  project${d}  ${image_b}  tag=${image_b_ver}  tag1=${image_b_ver}
    ${storage_quota_ret}=  Get Project Storage Quota Text From Project Quotas List  project${d}
    Should Match Regexp  ${storage_quota_ret}  ${image_b_size} of ${storage_quota}${storage_quota_unit}
    Cannot Push image  ${ip}  ${HARBOR_ADMIN}  ${HARBOR_PASSWORD}  project${d}  ${image_a}:${image_a_ver}  err_msg=adding 330.1 MiB of storage resource, which when updated to current usage of   err_msg_2=MiB will exceed the configured upper limit of ${storage_quota}.0 MiB
    Go Into Project  project${d}
    Delete Repo  project${d}  ${image_b}
    Push Image With Tag  ${ip}  ${HARBOR_ADMIN}  ${HARBOR_PASSWORD}  project${d}  ${image_a}  tag=${image_a_ver}  tag1=${image_a_ver}
    ${storage_quota_ret}=  Get Project Storage Quota Text From Project Quotas List  project${d}
    ${storage_quota_ret_str_left}  Fetch From Left  ${storage_quota_ret}  25.
    Log  ${storage_quota_ret_str_left}
    ${storage_quota_ret_str_right}  Fetch From Left  ${storage_quota_ret}  25.
    Log  ${storage_quota_ret_str_right}
    Log  ${storage_quota_ret_str_left}${storage_quota_ret_str_right}
    Should Be Equal As Strings  ${storage_quota_ret}  ${image_a_size} of ${storage_quota}${storage_quota_unit}
    Close Browser

Test Case - Project Quotas Control Under Copy
    Init Chrome Driver
    ${d}=  Get Current Date  result_format=%m%s
    ${image_a}=  Set Variable  redis
    ${image_b}=  Set Variable  logstash
    ${image_a_ver}=  Set Variable  donotremove5.0
    ${image_b_ver}=  Set Variable  do_not_remove_6.8.3
    ${storage_quota}=  Set Variable  330
    ${storage_quota_unit}=  Set Variable  MiB
    Sign In Harbor  ${HARBOR_URL}  ${HARBOR_ADMIN}  ${HARBOR_PASSWORD}
    Create An New Project And Go Into Project  project_a_${d}
    Create An New Project And Go Into Project  project_b_${d}  storage_quota=${storage_quota}  storage_quota_unit=${storage_quota_unit}
    Push Image With Tag  ${ip}  ${HARBOR_ADMIN}  ${HARBOR_PASSWORD}  project_a_${d}  ${image_a}  tag=${image_a_ver}  tag1=${image_a_ver}
    Push Image With Tag  ${ip}  ${HARBOR_ADMIN}  ${HARBOR_PASSWORD}  project_a_${d}  ${image_b}  tag=${image_b_ver}  tag1=${image_b_ver}
    Go Into Project  project_a_${d}
    Go Into Repo  project_a_${d}/${image_a}
    Copy Image  ${image_a_ver}  project_b_${d}  ${image_a}
    Retry Wait Element Not Visible  ${repo_retag_confirm_dlg}
    Go Into Project  project_a_${d}
    Go Into Repo  project_a_${d}/${image_b}
    Copy Image  ${image_b_ver}  project_b_${d}  ${image_b}
    Retry Wait Element Not Visible  ${repo_retag_confirm_dlg}
    Sleep  2
    Go Into Project  project_b_${d}
    Sleep  2
    Retry Wait Until Page Contains Element  xpath=//clr-dg-cell[contains(.,'${image_a}')]/a
    Retry Wait Until Page Not Contains Element  xpath=//clr-dg-cell[contains(.,'${image_b}')]/a
    Close Browser

Test Case - Tag CRUD
    Init Chrome Driver
    Sign In Harbor  ${HARBOR_URL}  ${HARBOR_ADMIN}  ${HARBOR_PASSWORD}
    ${d}=    Get Current Date    result_format=%m%s
    Create An New Project And Go Into Project  project${d}
    Push Image With Tag  ${ip}  ${HARBOR_ADMIN}  ${HARBOR_PASSWORD}  project${d}  hello-world  latest
    Switch To Project Repo
    Go Into Repo   hello-world
    Go Into Artifact   latest
    Should Contain Tag   latest
    # add more than one tag
    Add A New Tag   123
    Add A New Tag   456
    Delete A Tag  latest
    Close Browser

Test Case - Tag Retention
    Init Chrome Driver
    Sign In Harbor  ${HARBOR_URL}  ${HARBOR_ADMIN}  ${HARBOR_PASSWORD}
    ${d}=    Get Current Date    result_format=%m%s
    ${image_sample_1}=    Set Variable  hello-world
    ${image_sample_2}=    Set Variable  memcached
    Create An New Project And Go Into Project  project${d}
    Switch To Tag Retention
    Add A Tag Retention Rule
    Delete A Tag Retention Rule
    Add A Tag Retention Rule
    Edit A Tag Retention Rule    **   latest
    Push Image With Tag  ${ip}  ${HARBOR_ADMIN}  ${HARBOR_PASSWORD}  project${d}  ${image_sample_1}  latest
    Push Image With Tag  ${ip}  ${HARBOR_ADMIN}  ${HARBOR_PASSWORD}  project${d}  ${image_sample_2}   123
    Set Daily Schedule
    Execute Dry Run  ${image_sample_2}  0/1
    Execute Run  ${image_sample_2}  0/1
    Execute Dry Run  ${image_sample_1}  1/1
    Execute Run  ${image_sample_1}  1/1
    Close Browser

Test Case - Tag Immutability
    Init Chrome Driver
    Sign In Harbor  ${HARBOR_URL}  ${HARBOR_ADMIN}  ${HARBOR_PASSWORD}
    ${d}=    Get Current Date    result_format=%m%s
    Create An New Project And Go Into Project  project${d}
    Switch To Tag Immutability
    @{param}  Create List  1212  3434
    Retry Add A Tag Immutability Rule  @{param}
    Delete A Tag Immutability Rule
    @{param}  Create List  5566  7788
    Retry Add A Tag Immutability Rule  @{param}
    Edit A Tag Immutability Rule  hello-world  latest
    Push Image With Tag  ${ip}  ${HARBOR_ADMIN}  ${HARBOR_PASSWORD}  project${d}  hello-world  latest
    Push Image With Tag  ${ip}  ${HARBOR_ADMIN}  ${HARBOR_PASSWORD}  project${d}  busybox  latest
    Go Into Project  project${d}
    @{repo_list}  Create List  hello-world  busybox
    Multi-delete Object  ${repo_delete_btn}  @{repo_list}
    # Verify
    Delete Fail  hello-world
    Delete Success  busybox
    Close Browser

Test Case - Robot Account
    [tags]  robot_account
    Init Chrome Driver
    ${d}=    Get Current Date    result_format=%m%s
    Sign In Harbor    ${HARBOR_URL}    ${HARBOR_ADMIN}    ${HARBOR_PASSWORD}
    Create An New Project And Go Into Project    project${d}
    ${token}=    Create A Robot Account And Return Token    project${d}    robot${d}
    Log To Console    ${token}
    Log    ${token}
    Push image  ${ip}  robot${d}  ${token}  project${d}  hello-world:latest  is_robot=${true}
    Pull image  ${ip}  robot${d}  ${token}  project${d}  hello-world:latest  is_robot=${true}
    Close Browser

Test Case - Push Docker Manifest Index and Display
    Init Chrome Driver
    ${d}=    Get Current Date    result_format=%m%s
    ${image_a}=  Set Variable  hello-world
    ${image_b}=  Set Variable  busybox
    ${image_a_ver}=  Set Variable  latest
    ${image_b_ver}=  Set Variable  latest

    Sign In Harbor  ${HARBOR_URL}  user010  Test1@34
    Create An New Project And Go Into Project  test${d}

    Push image  ${ip}  ${HARBOR_ADMIN}  ${HARBOR_PASSWORD}  test${d}  ${image_a}:${image_a_ver}
    Go Into Project  test${d}
    Wait Until Page Contains  test${d}/${image_a}

    Push image  ${ip}  ${HARBOR_ADMIN}  ${HARBOR_PASSWORD}  test${d}  ${image_b}:${image_b_ver}
    Go Into Project  test${d}
    Wait Until Page Contains  test${d}/${image_b}

    Docker Push Index  ${ip}  user010  Test1@34  ${ip}/test${d}/index${d}:index_tag${d}  ${ip}/test${d}/${image_a}:${image_a_ver}  ${ip}/test${d}/${image_b}:${image_b_ver}

    Go Into Project  test${d}
    Wait Until Page Contains  test${d}/index${d}

    Go Into Repo  test${d}/index${d}
    Wait Until Page Contains  index_tag${d}

    Go Into Project  test${d}
    Wait Until Page Contains  test${d}/index${d}
    Go Into Repo  test${d}/index${d}
    Go Into Index And Contain Artifacts  index_tag${d}  total_artifact_count=2
    Close Browser

Test Case - Push Helm Chart and Display
    Init Chrome Driver
    ${d}=    Get Current Date    result_format=%m%s
    ${chart_file}=  Set Variable  https://storage.googleapis.com/harbor-builds/helm-chart-test-files/harbor-0.2.0.tgz
    ${archive}=  Set Variable  harbor/
    ${verion}=  Set Variable  0.2.0
    ${repo_name}=  Set Variable  harbor_chart_test

    Sign In Harbor  ${HARBOR_URL}  user010  Test1@34
    Create An New Project And Go Into Project  test${d}

    Retry Action Keyword  Helm Chart Push  ${ip}  user010  Test1@34  ${chart_file}  ${archive}  test${d}  ${repo_name}  ${verion}

    Go Into Project  test${d}
    Wait Until Page Contains  test${d}/${repo_name}

    Go Into Repo  test${d}/${repo_name}
    Wait Until Page Contains  ${repo_name}
    Go Into Project  test${d}
    Wait Until Page Contains  test${d}/${repo_name}
    Retry Double Keywords When Error  Go Into Repo  test${d}/${repo_name}  Page Should Contain Element  ${tag_table_column_vulnerabilities}
    Close Browser

Test Case - Can Not Copy Image In ReadOnly Mode
    Init Chrome Driver
    ${random_num1}=   Get Current Date    result_format=%m%s
    ${random_num2}=   Evaluate  str(random.randint(1000,9999))  modules=random

    Sign In Harbor  ${HARBOR_URL}  ${HARBOR_ADMIN}  ${HARBOR_PASSWORD}
    Create An New Project And Go Into Project  project${random_num1}${random_num2}
    Create An New Project And Go Into Project  project${random_num1}

    Sleep  1
    Push Image With Tag  ${ip}  ${HARBOR_ADMIN}  ${HARBOR_PASSWORD}  project${random_num1}  redis  ${image_tag}
    Sleep  1
    Enable Read Only
    Go Into Repo  project${random_num1}/redis
    Copy Image  ${image_tag}  project${random_num1}${random_num2}  ${target_image_name}
    Retry Wait Element Not Visible  ${repo_retag_confirm_dlg}
    Navigate To Projects
    Go Into Project  project${random_num1}${random_num2}  has_image=${false}
    Sleep  10
    Go Into Project  project${random_num1}${random_num2}  has_image=${false}
    Disable Read Only
    Sleep  10
    Close Browser

Test Case - Read Only Mode
    Init Chrome Driver
    ${d}=   Get Current Date    result_format=%m%s
    Sign In Harbor  ${HARBOR_URL}  ${HARBOR_ADMIN}  ${HARBOR_PASSWORD}
    Create An New Project And Go Into Project  project${d}

    Enable Read Only
    Cannot Push image  ${ip}  ${HARBOR_ADMIN}  ${HARBOR_PASSWORD}  project${d}  busybox:latest

    Disable Read Only
    Sleep  5
    Push image  ${ip}  ${HARBOR_ADMIN}  ${HARBOR_PASSWORD}  project${d}  busybox:latest
    Close Browser

Test Case - System Robot Account Cover All Projects
    [Tags]  sys_robot_account_cover
    ${d}=  Get Current Date    result_format=%m%s
    ${pro_name}=  Set Variable  project_${d}
    Init Chrome Driver
    Sign In Harbor  ${HARBOR_URL}  ${HARBOR_ADMIN}  ${HARBOR_PASSWORD}
    Create An New Project And Go Into Project  ${pro_name}
    ${name}  ${secret}=  Create A New System Robot Account  is_cover_all=${true}
    Navigate To Projects
    Switch To Robot Account
    System Robot Account Exist  ${name}  all
    Close Browser

Test Case - System Robot Account
    [Tags]  sys_robot_account
    ${d}=  Get Current Date    result_format=%m%s
    ${project_count}=  Evaluate  random.randint(3, 5)
    ${pro_name}=  Set Variable  project_${d}
    Init Chrome Driver
    Sign In Harbor  ${HARBOR_URL}  ${HARBOR_ADMIN}  ${HARBOR_PASSWORD}
    ${project_permission_list}=  Create A Random Project Permission List  ${project_count}
    ${name}  ${secret}=  Create A New System Robot Account  project_permission_list=${project_permission_list}
    System Robot Account Exist  ${name}  ${project_count}
    Close Browser

Test Case - Go To Harbor Api Page
    [Tags]  go_to_harbor_api_page
    Init Chrome Driver
    Sign In Harbor  ${HARBOR_URL}  ${HARBOR_ADMIN}  ${HARBOR_PASSWORD}
    Retry Keyword N Times When Error  4  Check Harbor Api Page
    Close Browser

Test Case - WASM Push And Pull To Harbor
    [Tags]  wasm_push_and_pull_to_harbor
    Init Chrome Driver
    ${user}=    Set Variable    user004
    ${pwd}=    Set Variable    Test1@34
    Sign In Harbor  ${HARBOR_URL}  ${user}  ${pwd}
    ${d}=   Get Current Date    result_format=%m%s
    Create An New Project And Go Into Project  project${d}
    Run  wget https://github.com/engineerd/wasm-to-oci/blob/v0.1.2/testdata/hello.wasm
    Wait Unitl Command Success  docker login -u ${user} -p ${pwd} ${ip}
    Wait Unitl Command Success  wasm-to-oci push hello.wasm ${ip}/project${d}/wasm-to-oci:v1
    Wait Unitl Command Success  wasm-to-oci pull ${ip}/project${d}/wasm-to-oci:v1 --out test.wasm
    Wait Unitl Command Success  docker logout ${ip}
    Retry file should exist  test.wasm
    Close Browser

Test Case - Carvel Imgpkg Push And Pull To Harbor
    [Tags]  imgpkg_push_and_pull
    Init Chrome Driver
    ${user}=  Set Variable  user004
    ${pwd}=  Set Variable  Test1@34
    ${out_path}=  Set Variable  /tmp/my-bundle
    ${repository}=  Set Variable  my-bundle
    ${tag}=  Set Variable  v1.0.0
    Sign In Harbor  ${HARBOR_URL}  ${user}  ${pwd}
    ${d}=  Get Current Date  result_format=%m%s
    Create An New Project And Go Into Project  project${d}
    Prepare Image Package Test Files  ${EXECDIR}/config
    Docker Login  ${ip}  ${user}  ${pwd}
    Imgpkg Push  ${ip}  project${d}  ${repository}  ${tag}  ${EXECDIR}/config
    Imgpkg Pull  ${ip}  project${d}  ${repository}  ${tag}  ${out_path}
    Docker Logout  ${ip}
    Retry File Should Exist  ${out_path}/.imgpkg/bundle.yml
    Retry File Should Exist  ${out_path}/.imgpkg/images.yml
    Close Browser

Test Case - Cosign And Cosign Deployment Security Policy
    [Tags]  cosign
    Init Chrome Driver
    ${user}=  Set Variable  user006
    ${pwd}=  Set Variable  Test1@34
    ${d}=  Get Current Date  result_format=%m%s
    ${image}=  Set Variable  hello-world
    ${tag}=  Set Variable  latest
    Sign In Harbor  ${HARBOR_URL}  ${user}  ${pwd}
    Create An New Project And Go Into Project  project${d}
    Goto Project Config
    Click Cosign Deployment Security
    Save Project Config
    Content Cosign Deployment security Be Selected

    Push Image With Tag  ${ip}  ${user}  ${pwd}  project${d}  ${image}  ${tag}
    Go Into Project  project${d}
    Retry Double Keywords When Error  Go Into Repo  project${d}/${image}  Should Not Be Signed By Cosign  ${tag}
    Cannot Pull Image  ${ip}  ${user}  ${pwd}  project${d}  ${image}:${tag}  err_msg=The image is not signed in Cosign.
    Cosign Generate Key Pair
    Cosign Verify  ${ip}/project${d}/${image}:${tag}  ${false}

    Cosign Sign  ${ip}/project${d}/${image}:${tag}
    Cosign Verify  ${ip}/project${d}/${image}:${tag}  ${true}
    Retry Double Keywords When Error  Retry Element Click  ${artifact_list_refresh_btn}  Should Be Signed By Cosign  ${tag}
    Pull image  ${ip}  ${user}  ${pwd}  project${d}  ${image}:${tag}

    Retry Double Keywords When Error  Delete Accessory  ${tag}  Should be Accessory deleted  ${tag}
    Close Browser

Test Case - Audit Log And Purge
    [Tags]  audit_log_and_purge
    Init Chrome Driver
    ${user}=  Set Variable  user003
    ${pwd}=  Set Variable  Test1@34
    ${d}=  Get Current Date  result_format=%m%s
    ${image}=  Set Variable  alpine
    ${tag1}=  Set Variable  3.10
    ${tag2}=  Set Variable  test
    ${sha256}=  Set Variable  sha256:de78803598bc4c940fc4591d412bffe488205d5d953f94751c6308deeaaa7eb8
    Sign In Harbor  ${HARBOR_URL}  ${user}  ${pwd}
    # create project
    Create An New Project And Go Into Project  project${d}
    Switch To Logs
    Verify Log  ${user}  project${d}  project  create
    # create artifact
    Push Image With Tag  ${ip}  ${user}  ${pwd}  project${d}  ${image}  ${tag1}  ${tag1}
    Clean All Local Images
    Verify Log  ${user}  project${d}/${image}:${tag1}  artifact  create
    Go Into Project  project${d}
    Go Into Repo  ${image}
    Go Into Artifact  ${tag1}
    # create tag
    Add A New Tag   ${tag2}
    # delete tag
    Delete A Tag  ${tag2}
    Switch To Logs
    Verify Log  ${user}  project${d}/${image}:${tag2}  tag  delete
    Verify Log  ${user}  project${d}/${image}:${tag2}  tag  create  2
    Docker Login  ${ip}  ${user}  ${pwd}
    # pull artifact
    Docker Pull  ${ip}/project${d}/${image}:${tag1}
    Docker Logout  ${ip}
    Verify Log  ${user}  project${d}/${image}:${sha256}  artifact  pull
    Go Into Project  project${d}
    Go Into Repo  project${d}/${image}
    # delete artifact
    @{tag_list}  Create List  ${tag1}
    Multi-delete Artifact  @{tag_list}
    Switch To Logs
    Verify Log  ${user}  project${d}/${image}:${sha256}  artifact  delete
    Go Into Project  project${d}
    # delete repository
    Delete Repo  project${d}  ${image}
    Switch To Logs
    Verify Log  ${user}  project${d}/${image}  repository  delete
    # delete project
    Delete Project  project${d}
    Logout Harbor
    Sign In Harbor  ${HARBOR_URL}  ${HARBOR_ADMIN}  ${HARBOR_PASSWORD}
    Switch To Logs
    Verify Log  ${user}  project${d}  project  delete
    Switch To Log Rotation
    Purge Now  1  Hours
    Close Browser

Test Case - Audit Log Forward
    [Tags]  audit_log_forward
    Init Chrome Driver
    ${d}=  Get Current Date  result_format=%m%s
    ${audit_log_path}=  Set Variable  ${log_path}/audit.log
    ${syslog_endpoint}=  Set Variable  harbor-log:10514
    ${test_endpoint}=  Set Variable  test.endpoint
    ${image}=  Set Variable  alpine
    ${tag1}=  Set Variable  3.10
    ${tag2}=  Set Variable  test
    Sign In Harbor  ${HARBOR_URL}  ${HARBOR_ADMIN}  ${HARBOR_PASSWORD}
    Create An New Project And Go Into Project  project${d}
    Switch To Logs
    Verify Log  ${HARBOR_ADMIN}  project${d}  project  create
    Switch To System Settings
    Retry Wait Element Should Be Disabled  ${skip_audit_log_database_checkbox}
    Set Audit Log Forward  ${test_endpoint}  bad request: could not connect to the audit endpoint: ${test_endpoint}
    # Set Audit Log Forward
    Set Audit Log Forward  ${syslog_endpoint}  Configuration has been successfully saved.
    Wait Until Element Is Enabled  ${skip_audit_log_database_checkbox}
    # create artifact
    Push Image With Tag  ${ip}  ${HARBOR_ADMIN}  ${HARBOR_PASSWORD}  project${d}  ${image}  ${tag1}  ${tag1}
    Switch To Logs
    Verify Log  ${HARBOR_ADMIN}  project${d}/${image}:${tag1}  artifact  create
    Verify Log In File  ${HARBOR_ADMIN}  project${d}/${image}:${tag1}  artifact  create
    # Enable Skip Audit Log Database
    Enable Skip Audit Log Database
    Go Into Project  project${d}
    Go Into Repo  ${image}
    Go Into Artifact  ${tag1}
    # create tag
    Add A New Tag   ${tag2}
    Switch To Logs
    Verify Log  ${HARBOR_ADMIN}  project${d}/${image}:${tag1}  artifact  create
    Verify Log In File  ${HARBOR_ADMIN}  project${d}/${image}:${tag2}  artifact  create
    Set Audit Log Forward  ${null}  Configuration has been successfully saved.
    Retry Wait Element Should Be Disabled  ${skip_audit_log_database_checkbox}
    Checkbox Should Not Be Selected  ${skip_audit_log_database_checkbox}
    Go Into Project  project${d}
    Go Into Repo  ${image}
    Go Into Artifact  ${tag1}
    # delete tag
    Delete A Tag  ${tag2}
    Switch To Logs
    Verify Log  ${HARBOR_ADMIN}  project${d}/${image}:${tag2}  tag  delete
    Verify Log In File  ${HARBOR_ADMIN}  project${d}/${image}:${tag2}  artifact  create
    Close Browser

Test Case - Export CVE
    [Tags]  export_cve
    Init Chrome Driver
    ${d}=  Get Current Date  result_format=%m%s
    ${user}=  Set Variable  user023
    ${pwd}=  Set Variable  Test1@34
    &{images}=  Create Dictionary  nginx=1.14.0  redis=5.0  alpine=3.9.4  photon=4.0-20210226  postgres=9.6
    ${labels}=  Create List  sys_level_export${d}  proj_level_export${d}
    ${cve_ids}=  Create List  CVE-2019-18224  CVE-2021-3997  CVE-2022-25315
    ${nginx_sha256}=  Set Variable  sha256:d43aa3719937f9df0502f8258f3034a21b720b5b9bbf01bbfdbd09871aac8930
    ${redis_sha256}=  Set Variable  sha256:e4b315ad03a1d1d9ff0c111e648a1a91066c09ead8352d3d6a48fa971a82922c
    ${expected_cve_data}=  Create List  project${d}/nginx,${nginx_sha256},${cve_ids}[1],libudev1,232-25+deb9u4,,Medium,,"{""CVSS"": {""redhat"": {""V3Score"": 5.5, ""V3Vector"": ""CVSS:3.1/AV:L/AC:L/PR:L/UI:N/S:U/C:N/I:N/A:H""}}}",Trivy
    ...                                 project${d}/nginx,${nginx_sha256},${cve_ids}[1],libsystemd0,232-25+deb9u4,,Medium,,"{""CVSS"": {""redhat"": {""V3Score"": 5.5, ""V3Vector"": ""CVSS:3.1/AV:L/AC:L/PR:L/UI:N/S:U/C:N/I:N/A:H""}}}",Trivy
    ...                                 project${d}/redis,${redis_sha256},${cve_ids}[0],libidn2-0,2.0.5-1,2.0.5-1+deb10u1,Critical,CWE-787,"{""CVSS"": {""nvd"": {""V2Score"": 7.5, ""V3Score"": 9.8, ""V2Vector"": ""AV:N/AC:L/Au:N/C:P/I:P/A:P"", ""V3Vector"": ""CVSS:3.1/AV:N/AC:L/PR:N/UI:N/S:U/C:H/I:H/A:H""}, ""redhat"": {""V3Score"": 5.6, ""V3Vector"": ""CVSS:3.0/AV:N/AC:H/PR:N/UI:N/S:U/C:L/I:L/A:L""}}}",Trivy
    ...                                 project${d}/redis,${redis_sha256},${cve_ids}[1],libudev1,241-7~deb10u2,,Medium,,"{""CVSS"": {""redhat"": {""V3Score"": 5.5, ""V3Vector"": ""CVSS:3.1/AV:L/AC:L/PR:L/UI:N/S:U/C:N/I:N/A:H""}}}",Trivy
    ...                                 project${d}/nginx,${nginx_sha256},${cve_ids}[2],libexpat1,2.2.0-2+deb9u1,2.2.0-2+deb9u5,Critical,CWE-190,"{""CVSS"": {""nvd"": {""V2Score"": 7.5, ""V3Score"": 9.8, ""V2Vector"": ""AV:N/AC:L/Au:N/C:P/I:P/A:P"", ""V3Vector"": ""CVSS:3.1/AV:N/AC:L/PR:N/UI:N/S:U/C:H/I:H/A:H""}, ""redhat"": {""V3Score"": 9.8, ""V3Vector"": ""CVSS:3.1/AV:N/AC:L/PR:N/UI:N/S:U/C:H/I:H/A:H""}}}",Trivy
    ...                                 project${d}/redis,${redis_sha256},${cve_ids}[1],libsystemd0,241-7~deb10u2,,Medium,,"{""CVSS"": {""redhat"": {""V3Score"": 5.5, ""V3Vector"": ""CVSS:3.1/AV:L/AC:L/PR:L/UI:N/S:U/C:N/I:N/A:H""}}}",Trivy
    Sign In Harbor  ${HARBOR_URL}  ${HARBOR_ADMIN}  ${HARBOR_PASSWORD}
    Switch To System Labels
    Create New Labels  ${labels}[0]
    Logout Harbor
    Sign In Harbor  ${HARBOR_URL}  ${user}  ${pwd}
    Create An New Project And Go Into Project  project${d}
    # push images
    FOR  ${image}  IN  @{images.keys()}
        Push Image With Tag  ${ip}  ${user}  ${pwd}  project${d}  ${image}  ${images['${image}']}  ${images['${image}']}
    END
    # scan images
    Refresh Repositories
    FOR  ${image}  IN  @{images.keys()}
        Go Into Repo  ${image}
        Scan Repo  ${images['${image}']}  Succeed
        Back Project Home  project${d}
    END
    Switch To Project Label
    Create New Labels  ${labels}[1]
    Switch To Project Repo
    Go Into Repo  nginx
    Add Labels To Tag  ${images['nginx']}  ${labels}[0]
    Back Project Home  project${d}
    Go Into Repo  redis
    Add Labels To Tag  ${images['redis']}  ${labels}[1]
    Navigate To Projects
    Should Not Be Export CVEs
    Retry Element Click  //clr-dg-row[1]//label
    Retry Element Click  //clr-dg-row[2]//label
    Should Not Be Export CVEs
    Export CVEs  project${d}  photon,postgres,nginx,redis  ${images['photon']},${images['nginx']},${images['redis']}  ${labels}  ${cve_ids}[0],${cve_ids}[1],${cve_ids}[2]
    ${csv_file_path}=  Download Latest CVE CSV File
    ${csv_file}=  OperatingSystem.Get File  ${csv_file_path}
    ${csv_file_content}=  Create List  ${csv_file}
    ${actual_cve_data}=  Split To Lines  @{csv_file_content}  1
    Lists Should Be Equal  ${expected_cve_data}  ${actual_cve_data}  ignore_order=True
    Close Browser