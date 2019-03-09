class SpiderType:
    SCRAPY = 'scrapy'
    PYSPIDER = 'pyspider'
    WEBMAGIC = 'webmagic'


class LangType:
    PYTHON = 'python'
    JAVASCRIPT = 'javascript'
    JAVA = 'java'
    GO = 'go'
    OTHER = 'other'


class CronEnabled:
    ON = 1
    OFF = 0


SUFFIX_IGNORE = [
    'pyc'
]

FILE_SUFFIX_LANG_MAPPING = {
    'py': LangType.PYTHON,
    'js': LangType.JAVASCRIPT,
    'java': LangType.JAVA,
    'go': LangType.GO,
}
