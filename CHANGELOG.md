# 0.6.0 (TBC)

(TBC)

# 0.5.1 (2020-07-31)

### Features / Enhancement
- **Added error message details**.
- **Added Golang programming language support**.
- **Added web driver installation scripts for Chrome Driver and Firefox**.
- **Support system tasks**. A "system task" is similar to normal spider task, it allows users to view logs of general tasks such as installing languages.
- **Changed methods of installing languages from RPC to system tasks**.

### Bug Fixes
- **Fixed first download repo 500 error in Spider Market page**. [#808](https://github.com/crawlab-team/crawlab/issues/808)
- **Fixed some translation issues**.
- **Fixed 500 error in task detail page**. [#810](https://github.com/crawlab-team/crawlab/issues/810)
- **Fixed password reset issue**. [#811](https://github.com/crawlab-team/crawlab/issues/811)
- **Fixed unable to download CSV issue**. [#812](https://github.com/crawlab-team/crawlab/issues/812)
- **Fixed unable to install node.js issue**. [#813](https://github.com/crawlab-team/crawlab/issues/813)
- **Fixed disabled status for batch adding schedules**. [#814](https://github.com/crawlab-team/crawlab/issues/814)

# 0.5.0 (2020-07-19)
### Features / Enhancement
- **Spider Market**. Allow users to download open-source spiders into Crawlab.
- **Batch actions**. Allow users to interact with Crawlab in batch fashions, e.g. batch run tasks, batch delete spiders, ect.
- **Migrate MongoDB driver to `MongoDriver`**.
- **Refactor and optmize node-related logics**.
- **Change default `task.workers` to 16**.
- **Change default nginx `client_max_body_size` to 200m**.
- **Support writing logs to ElasticSearch**.
- **Display error details in Scrapy page**.
- **Removed Challenge page**.
- **Moved Feedback and Dislaimer pages to navbar**.

### Bug Fixes
- **Fixed log not expiring issue because of failure to create TTL index**.
- **Set default log expire duration to 1 day**.
- **`task_id` index not created**.
- **`docker-compose.yml` fix**.
- **Fixed 404 page**.
- **Fixed unable to create worker node before master node issue**.

# 0.4.10 (2020-04-21)
### Features / Enhancement
- **Enhanced Log Management**. Centralizing log storage in MongoDB, reduced the dependency of PubSub, allowing log error detection.
- **API Token**. Allow users to generate API tokens and use them to integrate into their own systems.
- **Web Hook**. Trigger a Web Hook http request to pre-defined URL when a task starts or finishes.
- **Auto Install Dependencies**. Allow installing dependencies automatically from `requirements.txt` or `package.json`.
- **Auto Results Collection**. Set results collection to `results_<spider_name>` if it is not set.
- **Optimized Project List**. Not display "No Project" item in the project list.
- **Upgrade Node.js**. Upgrade Node.js version from v8.12 to v10.19.
- **Add Run Button in Schedule Page**. Allow users to manually run task in Schedule Page.

### Bug Fixes
- **Cannot register**. [#670](https://github.com/crawlab-team/crawlab/issues/670)
- **Spider schedule tab cron expression shows second**. [#678](https://github.com/crawlab-team/crawlab/issues/678)
- **Missing daily stats in spider**. [#684](https://github.com/crawlab-team/crawlab/issues/684)
- **Results count not update in time**. [#689](https://github.com/crawlab-team/crawlab/issues/689)

# 0.4.9 (2020-03-31)
### Features / Enhancement
- **Challenges**. Users can achieve different challenges based on their actions.
- **More Advanced Access Control**. More granular access control, e.g. normal users can only view/manage their own spiders/projects and admin users can view/manage all spiders/projects.
- **Feedback**. Allow users to send feedbacks and ratings to Crawlab team.
- **Better Home Page Metrics**. Optimized metrics display on home page.
- **Configurable Spiders Converted to Customized Spiders**. Allow users to convert their configurable spiders into customized spiders which are also Scrapy spiders.
- **View Tasks Triggered by Schedule**. Allow users to view tasks triggered by a schedule. [#648](https://github.com/crawlab-team/crawlab/issues/648)
- **Support Results De-Duplication**. Allow users to configure de-duplication of results. [#579](https://github.com/crawlab-team/crawlab/issues/579)
- **Support Task Restart**. Allow users to re-run historical tasks.

### Bug Fixes
- **CLI unable to use on Windows**. [#580](https://github.com/crawlab-team/crawlab/issues/580)
- **Re-upload error**. [#643](https://github.com/crawlab-team/crawlab/issues/643) [#640](https://github.com/crawlab-team/crawlab/issues/640)
- **Upload missing folders**. [#646](https://github.com/crawlab-team/crawlab/issues/646)
- **Unable to add schedules in Spider Page**.

# 0.4.8 (2020-03-11)
### Features / Enhancement
- **Support Installations of More Programming Languages**. Now users can install or pre-install more programming languages including Java, .Net Core and PHP.
- **Installation UI Optimization**. Users can better view and manage installations on Node List page.
- **More Git Support**. Allow users to view Git Commits record, and allow checkout to corresponding commit.
- **Support Hostname Node Registration Type**. Users can set hostname as the node key as the unique identifier.
- **RPC Support**. Added RPC support to better manage node communication.
- **Run On Master Switch**. Users can determine whether to run tasks on master. If not, all tasks will be run only on worker nodes.
- **Disabled Tutorial by Default**.
- **Added Related Documentation Sidebar**.
- **Loading Page Optimization**.

### Bug Fixes
- **Duplicated Nodes**. [#391](https://github.com/crawlab-team/crawlab/issues/391)
- **Duplicated Spider Upload**. [#603](https://github.com/crawlab-team/crawlab/issues/603)
- **Failure in dependencies installation results in unusable dependency installation functionalities.**. [#609](https://github.com/crawlab-team/crawlab/issues/609)
- **Create Tasks for Offline Nodes**. [#622](https://github.com/crawlab-team/crawlab/issues/622)

# 0.4.7 (2020-02-24)
### Features / Enhancement
- **Better Support for Scrapy**. Spiders identification, `settings.py` configuration, log level selection, spider selection. [#435](https://github.com/crawlab-team/crawlab/issues/435)
- **Git Sync**. Allow users to sync git projects to Crawlab.
- **Long Task Support**. Users can add long-task spiders which is supposed to run without finishing. [#425](https://github.com/crawlab-team/crawlab/issues/425)
- **Spider List Optimization**. Tasks count by status, tasks detail popup, legend. [#425](https://github.com/crawlab-team/crawlab/issues/425)
- **Upgrade Check**. Check latest version and notifiy users to upgrade.
- **Spiders Batch Operation**. Allow users to run/stop spider tasks and delete spiders in batches.
- **Copy Spiders**. Allow users to copy an existing spider to create a new one.
- **Wechat Group QR Code**.

### Bug Fixes
- **Schedule Spider Selection Issue**. Fields not responding to spider change.
- **Cron Jobs Conflict**. Possible bug when two spiders set to the same time of their cron jobs. [#515](https://github.com/crawlab-team/crawlab/issues/515) [#565](https://github.com/crawlab-team/crawlab/issues/565)
- **Task Log Issue**. Different tasks write to the same log file if triggered at the same time. [#577](https://github.com/crawlab-team/crawlab/issues/577)
- **Task List Filter Options Incomplete**.

# 0.4.6 (2020-02-13)
### Features / Enhancement
- **SDK for Node.js**. Users can apply SDK in their Node.js spiders.
- **Log Management Optimization**. Log search, error highlight, auto-scrolling.
- **Task Execution Process Optimization**. Allow users to be redirected to task detail page after triggering a task.
- **Task Display Optimization**. Added "Param" in the Latest Tasks table in the spider detail page. [#295](https://github.com/crawlab-team/crawlab/issues/295)
- **Spider List Optimization**. Added "Update Time" and "Create Time" in spider list page.
- **Page Loading Placeholder**. 

### Bug Fixes
- **Lost Focus in Schedule Configuration**. [#519](https://github.com/crawlab-team/crawlab/issues/519)
- **Unable to Upload Spider using CLI**. [#524](https://github.com/crawlab-team/crawlab/issues/524)

# 0.4.5 (2020-02-03)
### Features / Enhancement
- **Interactive Tutorial**. Guide users through the main functionalities of Crawlab.
- **Global Environment Variables**. Allow users to set global environment variables, which will be passed into all spider programs. [#177](https://github.com/crawlab-team/crawlab/issues/177)
- **Project**. Allow users to link spiders to projects. [#316](https://github.com/crawlab-team/crawlab/issues/316)
- **Demo Spiders**. Added demo spiders when Crawlab is initialized. [#379](https://github.com/crawlab-team/crawlab/issues/379)
- **User Admin Optimization**. Restrict privilleges of admin users. [#456](https://github.com/crawlab-team/crawlab/issues/456)
- **Setting Page Optimization**.
- **Task Results Optimization**.

### Bug Fixes
- **Unable to find spider file error**. [#485](https://github.com/crawlab-team/crawlab/issues/485)
- **Click delete button results in redirect**. [#480](https://github.com/crawlab-team/crawlab/issues/480)
- **Unable to create files in an empty spider**. [#479](https://github.com/crawlab-team/crawlab/issues/479)
- **Download results error**. [#465](https://github.com/crawlab-team/crawlab/issues/465)
- **crawlab-sdk CLI error**. [#458](https://github.com/crawlab-team/crawlab/issues/458)
- **Page refresh issue**. [#441](https://github.com/crawlab-team/crawlab/issues/441)
- **Results not support JSON**. [#202](https://github.com/crawlab-team/crawlab/issues/202)
- **Getting all spider after deleting a spider**.
- **i18n warning**.

# 0.4.4 (2020-01-17)
### Features / Enhancement
- **Email Notification**. Allow users to send email notifications.
- **DingTalk Robot Notification**. Allow users to send DingTalk Robot notifications.
- **Wechat Robot Notification**. Allow users to send Wechat Robot notifications.
- **API Address Optimization**. Added relative URL path in frontend so that users don't have to specify `CRAWLAB_API_ADDRESS` explicitly.
- **SDK Compatiblity**. Allow users to integrate Scrapy or general spiders with Crawlab SDK.
- **Enhanced File Management**. Added tree-like file sidebar to allow users to edit files much more easier.
- **Advanced Schedule Cron**. Allow users to edit schedule cron with visualized cron editor.

### Bug Fixes
- **`nil retuened` error**.
- **Error when using HTTPS**.
- **Unable to run Configurable Spiders on Spider List**.
- **Missing form validation before uploading spider files**.

# 0.4.3 (2020-01-07)

### Features / Enhancement
- **Dependency Installation**. Allow users to install/uninstall dependencies and add programming languages (Node.js only for now) on the platform web interface.
- **Pre-install Programming Languages in Docker**. Allow Docker users to set `CRAWLAB_SERVER_LANG_NODE` as `Y` to pre-install `Node.js` environments.
- **Add Schedule List in Spider Detail Page**. Allow users to view / add / edit schedule cron jobs in the spider detail page. [#360](https://github.com/crawlab-team/crawlab/issues/360)
- **Align Cron Expression with Linux**. Change the expression of 6 elements to 5 elements as aligned in Linux.
- **Enable/Disable Schedule Cron**. Allow users to enable/disable the schedule jobs. [#297](https://github.com/crawlab-team/crawlab/issues/297)
- **Better Task Management**. Allow users to batch delete tasks. [#341](https://github.com/crawlab-team/crawlab/issues/341)
- **Better Spider Management**. Allow users to sort and filter spiders in the spider list page.
- **Added Chinese `CHANGELOG`**.
- **Added Github Star Button at Nav Bar**.

### Bug Fixes
- **Schedule Cron Task Issue**. [#423](https://github.com/crawlab-team/crawlab/issues/423)
- **Upload Spider Zip File Issue**. [#403](https://github.com/crawlab-team/crawlab/issues/403) [#407](https://github.com/crawlab-team/crawlab/issues/407)
- **Exit due to Network Failure**. [#340](https://github.com/crawlab-team/crawlab/issues/340)
- **Cron Jobs not Running Correctly**
- **Schedule List Columns Mis-positioned**
- **Clicking Refresh Button Redirected to 404 Page**

# 0.4.2 (2019-12-26)
### Features / Enhancement
- **Disclaimer**. Added page for Disclaimer.
- **Call API to fetch version**. [#371](https://github.com/crawlab-team/crawlab/issues/371)
- **Configure to allow user registration**. [#346](https://github.com/crawlab-team/crawlab/issues/346)
- **Allow adding new users**.
- **More Advanced File Management**. Allow users to add / edit / rename / delete files. [#286](https://github.com/crawlab-team/crawlab/issues/286)
- **Optimized Spider Creation Process**. Allow users to create an empty customized spider before uploading the zip file.
- **Better Task Management**. Allow users to filter tasks by selecting through certian criterions. [#341](https://github.com/crawlab-team/crawlab/issues/341)

### Bug Fixes
- **Duplicated nodes**. [#391](https://github.com/crawlab-team/crawlab/issues/391)
- **"mongodb no reachable" error**. [#373](https://github.com/crawlab-team/crawlab/issues/373)

# 0.4.1 (2019-12-13)
### Features / Enhancement
- **Spiderfile Optimization**. Stages changed from dictionary to array. [#358](https://github.com/crawlab-team/crawlab/issues/358)
- **Baidu Tongji Update**.

### Bug Fixes
- **Unable to display schedule tasks**. [#353](https://github.com/crawlab-team/crawlab/issues/353)
- **Duplicate node registration**. [#334](https://github.com/crawlab-team/crawlab/issues/334)

# 0.4.0 (2019-12-06)
### Features / Enhancement
- **Configurable Spider**. Allow users to add spiders using *Spiderfile* to configure crawling rules.
- **Execution Mode**. Allow users to select 3 modes for task execution: *All Nodes*, *Selected Nodes* and *Random*.

### Bug Fixes
- **Task accidentally killed**. [#306](https://github.com/crawlab-team/crawlab/issues/306)
- **Documentation fix**. [#301](https://github.com/crawlab-team/crawlab/issues/258) [#301](https://github.com/crawlab-team/crawlab/issues/258)
- **Direct deploy incompatible with Windows**. [#288](https://github.com/crawlab-team/crawlab/issues/288)
- **Log files lost**. [#269](https://github.com/crawlab-team/crawlab/issues/269)

# 0.3.5 (2019-10-28)
### Features / Enhancement
- **Graceful Showdown**. [detail](https://github.com/crawlab-team/crawlab/commit/63fab3917b5a29fd9770f9f51f1572b9f0420385)
- **Node Info Optimization**. [detail](https://github.com/crawlab-team/crawlab/commit/973251a0fbe7a2184ac0da09e0404a17c736aee7)
- **Append System Environment Variables to Tasks**. [detail](https://github.com/crawlab-team/crawlab/commit/4ab4892471965d6342d30385578ca60dc51f8ad3)
- **Auto Refresh Task Log**. [detail](https://github.com/crawlab-team/crawlab/commit/4ab4892471965d6342d30385578ca60dc51f8ad3)
- **Enable HTTPS Deployment**. [detail](https://github.com/crawlab-team/crawlab/commit/5d8f6f0c56768a6e58f5e46cbf5adff8c7819228)

### Bug Fixes
- **Unable to fetch spider list info in schedule jobs**. [detail](https://github.com/crawlab-team/crawlab/commit/311f72da19094e3fa05ab4af49812f58843d8d93)
- **Unable to fetch node info from worker nodes**. [detail](https://github.com/crawlab-team/crawlab/commit/6af06efc17685a9e232e8c2b5fd819ec7d2d1674)
- **Unable to select node when trying to run spider tasks**. [detail](https://github.com/crawlab-team/crawlab/commit/31f8e03234426e97aed9b0bce6a50562f957edad)
- **Unable to fetch result count when result volume is large**. [#260](https://github.com/crawlab-team/crawlab/issues/260)
- **Node issue in schedule tasks**. [#244](https://github.com/crawlab-team/crawlab/issues/244)


# 0.3.1 (2019-08-25)
### Features / Enhancement
- **Docker Image Optimization**. Split docker further into master, worker, frontend with alpine image.
- **Unit Tests**. Covered part of the backend code with unit tests.
- **Frontend Optimization**. Login page, button size, hints of upload UI optimization. 
- **More Flexible Node Registration**. Allow users to pass a variable as key for node registration instead of MAC by default.

### Bug Fixes
- **Uploading Large Spider Files Error**. Memory crash issue when uploading large spider files. [#150](https://github.com/crawlab-team/crawlab/issues/150)
- **Unable to Sync Spiders**. Fixes through increasing level of write permission when synchronizing spider files. [#114](https://github.com/crawlab-team/crawlab/issues/114)
- **Spider Page Issue**. Fixes through removing the field "Site". [#112](https://github.com/crawlab-team/crawlab/issues/112)
- **Node Display Issue**. Nodes do not display correctly when running docker containers on multiple machines. [#99](https://github.com/crawlab-team/crawlab/issues/99)

# 0.3.0 (2019-07-31)
### Features / Enhancement
- **Golang Backend**: Refactored code from Python backend to Golang, much more stability and performance.
- **Node Network Graph**: Visualization of node typology.
- **Node System Info**: Available to see system info including OS, CPUs and executables.
- **Node Monitoring Enhancement**: Nodes are monitored and registered through Redis.
- **File Management**: Available to edit spider files online, including code highlight.
- **Login/Regiser/User Management**: Require users to login to use Crawlab, allow user registration and user management, some role-based authorization.
- **Automatic Spider Deployment**: Spiders are deployed/synchronized to all online nodes automatically.
- **Smaller Docker Image**: Slimmed Docker image and reduced Docker image size from 1.3G to \~700M by applying Multi-Stage Build.

### Bug Fixes
- **Node Status**. Node status does not change even though it goes offline actually. [#87](https://github.com/tikazyq/crawlab/issues/87) 
- **Spider Deployment Error**. Fixed through Automatic Spider Deployment [#83](https://github.com/tikazyq/crawlab/issues/83) 
- **Node not showing**. Node not able to show online [#81](https://github.com/tikazyq/crawlab/issues/81) 
- **Cron Job not working**. Fixed through new Golang backend [#64](https://github.com/tikazyq/crawlab/issues/64) 
- **Flower Error**. Fixed through new Golang backend [#57](https://github.com/tikazyq/crawlab/issues/57) 

# 0.2.4 (2019-07-07)
### Features / Enhancement
- **Documentation**: Better and much more detailed documentation.
- **Better Crontab**: Make crontab expression through crontab UI.
- **Better Performance**: Switched from native flask engine to `gunicorn`. [#78](https://github.com/tikazyq/crawlab/issues/78)

### Bugs Fixes
- **Deleting Spider**. Deleting a spider does not only remove record in db but also removing related folder, tasks and schedules. [#69](https://github.com/tikazyq/crawlab/issues/69)
- **MongoDB Auth**. Allow user to specify `authenticationDatabase` to connect to `mongodb`. [#68](https://github.com/tikazyq/crawlab/issues/68)
- **Windows Compatibility**. Added `eventlet` to `requirements.txt`. [#59](https://github.com/tikazyq/crawlab/issues/59)


# 0.2.3 (2019-06-12)
### Features / Enhancement
- **Docker**: User can run docker image to speed up deployment.
- **CLI**: Allow user to use command-line interface to execute Crawlab programs.
- **Upload Spider**: Allow user to upload Customized Spider to Crawlab.
- **Edit Fields on Preview**: Allow user to edit fields when previewing data in Configurable Spider.

### Bugs Fixes
- **Spiders Pagination**. Fixed pagination problem in spider page.

# 0.2.2 (2019-05-30)
### Features / Enhancement
- **Automatic Extract Fields**: Automatically extracting data fields in list pages for configurable spider.
- **Download Results**: Allow downloading results as csv file.
- **Baidu Tongji**: Allow users to choose to report usage info to Baidu Tongji.

### Bug Fixes
- **Results Page Pagination**: Fixes so the pagination of results page is working correctly. [#45](https://github.com/tikazyq/crawlab/issues/45)
- **Schedule Tasks Duplicated Triggers**: Set Flask DEBUG as False so that schedule tasks won't trigger twice. [#32](https://github.com/tikazyq/crawlab/issues/32)
- **Frontend Environment**: Added `VUE_APP_BASE_URL` as production mode environment variable so the API call won't be always `localhost` in deployed env [#30](https://github.com/tikazyq/crawlab/issues/30)

# 0.2.1 (2019-05-27)
- **Configurable Spider**: Allow users to create a spider to crawl data without coding.

# 0.2 (2019-05-10)

- **Advanced Stats**: Advanced analytics in spider detail view.
- **Sites Data**: Added sites list (China) for users to check info such as robots.txt and home page response time/code.

# 0.1.1 (2019-04-23)

- **Basic Stats**: User can view basic stats such as number of failed tasks and number of results in spiders and tasks pages.
- **Near Realtime Task Info**: Periodically (5 sec) polling data from server to allow view task info in a near-realtime fashion.
- **Scheduled Tasks**: Allow users to set up cron-like scheduled/periodical tasks using apscheduler.

# 0.1 (2019-04-17)

- **Initial Release**
