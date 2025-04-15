# coding: utf-8
try:
    import urllib.request # Python 3.x 的标准库，用于下载数据
except ImportError:
    raise ImportError('You should use Python 3.x')
import os.path # 用于处理文件路径
import gzip # 用于解压缩 MNIST 数据集（gzip 格式）
import pickle # 用于序列化和反序列化 Python 对象
import os
import numpy as np


url_base = 'https://ossci-datasets.s3.amazonaws.com/mnist/'  # mirror site
key_file = {
    'train_img':'train-images-idx3-ubyte.gz',  # 训练集图片
    'train_label':'train-labels-idx1-ubyte.gz', # 训练集标签
    'test_img':'t10k-images-idx3-ubyte.gz', # 测试集图片
    'test_label':'t10k-labels-idx1-ubyte.gz' # 测试集标签
}

dataset_dir = os.path.dirname(os.path.abspath(__file__))
save_file = dataset_dir + "/mnist.pkl"

train_num = 60000 # 训练数据数量
test_num = 10000 # 测试数据数量
img_dim = (1, 28, 28)  # 图片尺寸（通道数, 高, 宽）, ‌图片的通道数是指图片中每个像素点可以包含的颜色信息的数量
img_size = 784 # 28x28=784, 线性展开为一维向量, 784个像素点


def _download(file_name):
    file_path = dataset_dir + "/" + file_name
    
    if os.path.exists(file_path):
        return

    print("Downloading " + file_name + " ... ")
    urllib.request.urlretrieve(url_base + file_name, file_path)
    print("Done")
    
def download_mnist():
    """ 下载所有 MNIST 相关文件 """
    for v in key_file.values():
       _download(v)
        
def _load_label(file_name):
    """ 读取标签数据，并转换为 NumPy 数组 """
    file_path = dataset_dir + "/" + file_name
    
    print("Converting " + file_name + " to NumPy Array ...")
    with gzip.open(file_path, 'rb') as f:
            labels = np.frombuffer(f.read(), np.uint8, offset=8)
    print("Done")
    
    return labels

def _load_img(file_name):
    """ 读取图片数据，并转换为 NumPy 数组 """
    file_path = dataset_dir + "/" + file_name
    
    print("Converting " + file_name + " to NumPy Array ...")    
    with gzip.open(file_path, 'rb') as f:
            data = np.frombuffer(f.read(), np.uint8, offset=16)
    data = data.reshape(-1, img_size)
    print("Done")
    
    return data
    
def _convert_numpy():
    """ 将所有数据转换为 NumPy 数组，并存入字典 """
    dataset = {}
    dataset['train_img'] =  _load_img(key_file['train_img'])
    dataset['train_label'] = _load_label(key_file['train_label'])    
    dataset['test_img'] = _load_img(key_file['test_img'])
    dataset['test_label'] = _load_label(key_file['test_label'])
    
    return dataset

def init_mnist():
    download_mnist()
    dataset = _convert_numpy()
    print("Creating pickle file ...")
    with open(save_file, 'wb') as f:
        pickle.dump(dataset, f, -1)
    print("Done!")

def _change_one_hot_label(X):
    """ 将标签转换为 one-hot 形式 """
    T = np.zeros((X.size, 10))
    for idx, row in enumerate(T):
        row[X[idx]] = 1
        
    return T
    

def load_mnist(normalize=True, flatten=True, one_hot_label=False):
    """读入MNIST数据集
    
    Parameters
    ----------
    normalize : 将图像的像素值正规化为0.0~1.0
    one_hot_label : 
        one_hot_label为True的情况下，标签作为one-hot数组返回
        one-hot数组是指[0,0,1,0,0,0,0,0,0,0]这样的数组
    flatten : 是否将图像展开为一维数组
    
    Returns
    -------
    (训练图像, 训练标签), (测试图像, 测试标签)
    """
    if not os.path.exists(save_file):
        init_mnist()
        
    with open(save_file, 'rb') as f:
        dataset = pickle.load(f)

    # 避免数值过大 0 - 255, 转换成 0 - 1 之间，避免梯度爆炸或梯度消失
    if normalize:
        for key in ('train_img', 'test_img'):
            dataset[key] = dataset[key].astype(np.float32) # 转换为 float32
            # 127 → 127/255 ≈ 0.498
            # 255 → 255/255 = 1.0
            # 0 → 0/255 = 0.0
            # 超过 255 的整数（取模 256）
            # 负数（取模 256 变成正数）
            # 小数（截断小数部分后取模 256）
            dataset[key] /= 255.0 # 归一化到 [0,1]
            
    if one_hot_label:
        dataset['train_label'] = _change_one_hot_label(dataset['train_label'])
        dataset['test_label'] = _change_one_hot_label(dataset['test_label'])
    
    if not flatten:
         for key in ('train_img', 'test_img'):
            dataset[key] = dataset[key].reshape(-1, 1, 28, 28)

    return (dataset['train_img'], dataset['train_label']), (dataset['test_img'], dataset['test_label']) 


if __name__ == '__main__':
    init_mnist()
