from app import api
from routes.base import BaseApi
from tasks.spider import get_baidu_html


class TestApi(BaseApi):
    col_name = 'test'

    def __init__(self):
        super(TestApi).__init__()
        self.parser.add_argument('keyword', type=str)

    def get(self, id=None):
        args = self.parser.parse_args()
        for i in range(100):
            get_baidu_html.delay(args.keyword)
        return {
            'status': 'ok'
        }


# add api to resources
api.add_resource(TestApi,
                 '/api/test',
                 )
