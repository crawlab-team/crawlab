# 0.3.0 (2019-07-31)
### Features / Enhancement
- **Golang**: Refactored code from Python backend to Golang, much more stability and performance.
- **Node Network Graph**: Visualization of node typology.
- **Node System Info**: Available to see system info including OS, CPUs and executables.

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
