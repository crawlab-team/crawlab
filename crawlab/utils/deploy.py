import os, zipfile
from utils.log import other


def zip_file(source_dir, output_filename):
    """
    打包目录为zip文件（未压缩）
    :param source_dir: source directory
    :param output_filename: output file name
    """
    zipf = zipfile.ZipFile(output_filename, 'w')
    pre_len = len(os.path.dirname(source_dir))
    for parent, dirnames, filenames in os.walk(source_dir):
        for filename in filenames:
            pathfile = os.path.join(parent, filename)
            arcname = pathfile[pre_len:].strip(os.path.sep)  # 相对路径
            zipf.write(pathfile, arcname)
    zipf.close()


def unzip_file(zip_src, dst_dir):
    """
    Unzip file
    :param zip_src: source zip file
    :param dst_dir: destination directory
    """
    r = zipfile.is_zipfile(zip_src)
    if r:
        fz = zipfile.ZipFile(zip_src, 'r')
        for file in fz.namelist():
            fz.extract(file, dst_dir)
    else:
        other.info('This is not zip')
