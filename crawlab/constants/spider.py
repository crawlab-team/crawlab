class SpiderType:
    CONFIGURABLE = 'configurable'
    CUSTOMIZED = 'customized'


class LangType:
    PYTHON = 'python'
    JAVASCRIPT = 'javascript'
    JAVA = 'java'
    GO = 'go'
    OTHER = 'other'


class CronEnabled:
    ON = 1
    OFF = 0


class CrawlType:
    LIST = 'list'
    DETAIL = 'detail'
    LIST_DETAIL = 'list-detail'


class QueryType:
    CSS = 'css'
    XPATH = 'xpath'


class ExtractType:
    TEXT = 'text'
    ATTRIBUTE = 'attribute'


SUFFIX_IGNORE = [
    'pyc'
]

FILE_SUFFIX_LANG_MAPPING = {
    'py': LangType.PYTHON,
    'js': LangType.JAVASCRIPT,
    'java': LangType.JAVA,
    'go': LangType.GO,
}
