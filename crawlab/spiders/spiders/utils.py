import itertools
import re


def generate_urls(base_url: str) -> str:
    url = base_url

    # number range list
    list_arr = []
    for i, res in enumerate(re.findall(r'{(\d+),(\d+)}', base_url)):
        try:
            _min = int(res[0])
            _max = int(res[1])
        except ValueError as err:
            raise ValueError(f'{base_url} is not a valid URL pattern')

        # list
        _list = range(_min, _max + 1)

        # key
        _key = f'n{i}'

        # append list and key
        list_arr.append((_list, _key))

        # replace url placeholder with key
        url = url.replace('{' + res[0] + ',' + res[1] + '}', '{' + _key + '}', 1)

    # string list
    for i, res in enumerate(re.findall(r'\[(.+)\]', base_url)):
        # list
        _list = res.split(',')

        # key
        _key = f's{i}'

        # append list and key
        list_arr.append((_list, _key))

        # replace url placeholder with key
        url = url.replace('[' + ','.join(_list) + ']', '{' + _key + '}', 1)

    # combine together
    _list_arr = []
    for res in itertools.product(*map(lambda x: x[0], list_arr)):
        _url = url
        for _arr, _rep in zip(list_arr, res):
            _list, _key = _arr
            _url = _url.replace('{' + _key + '}', str(_rep), 1)
        yield _url
