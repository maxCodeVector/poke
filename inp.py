import json


def loadFont(path):
    f = open(path, encoding='utf-8')
    setting = json.load(f)
    return setting
