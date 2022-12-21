import os

dir_path = '/app/dist/assets'

for file_name in os.listdir(dir_path):
    if not file_name.endswith('.js'):
        continue

    file_path = os.path.join(dir_path, file_name)

    api_url = 'http://localhost:8000'

    with open(file_path, 'r') as f:
        content = f.read()

    if api_url not in content:
        continue

    content = content.replace(api_url, '/api')

    with open(file_path, 'w') as f:
        f.write(content)

    print(f'replaced api url in file: {file_name}')
