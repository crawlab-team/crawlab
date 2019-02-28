import os
import re
from collections import defaultdict

SUFFIX_PATTERN = r'\.(\w{,10})$'
suffix_regex = re.compile(SUFFIX_PATTERN, re.IGNORECASE)

SUFFIX_LANG_MAPPING = {
    'py': 'python',
    'js': 'javascript',
    'sh': 'shell',
    'java': 'java',
    'c': 'c',
}


def get_file_suffix(file_name: str):
    file_name = file_name.lower()
    m = suffix_regex.search(file_name)
    if m is not None:
        return m.groups()[0]
    else:
        return file_name


def get_file_list(path):
    for root, dirs, file_names in os.walk(path):
        # print(root)  # 当前目录路径
        # print(dirs)  # 当前路径下所有子目录
        # print(file_names)  # 当前路径下所有非目录子文件

        for file_name in file_names:
            file_path = os.path.join(root, file_name)
            yield file_path


def get_file_suffix_stats(path) -> dict:
    stats = defaultdict(int)
    for file_path in get_file_list(path):
        suffix = get_file_suffix(file_path)
        stats[suffix] += 1
    return stats


def get_file_content(path) -> dict:
    with open(path) as f:
        suffix = get_file_suffix(path)
        lang = SUFFIX_LANG_MAPPING.get(suffix)
        return {
            'lang': lang,
            'suffix': suffix,
            'content': f.read()
        }
