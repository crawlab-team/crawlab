# 0.3.5 (2019-10-28)
### Features / Enhancement
- **Graceful Showdown**. [detail](https://github.com/crawlab-team/crawlab/commit/63fab3917b5a29fd9770f9f51f1572b9f0420385)**
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
