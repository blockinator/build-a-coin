package data

import (
    "io"
    "io/ioutil"
    "os"
)

func LoadAsset(conf *Conf, path string) (io.ReadCloser, error) {
    file, err := os.Open(conf.AssetsDir() + "/" + path)
    if err != nil {
        return nil, err
    }
    return file, nil
}

func ReadAsset(conf *Conf, path string) ([]byte, error) {
    file, err := LoadAsset(conf, path)
    if err != nil {
        return nil, err
    }
    output, err := ioutil.ReadAll(file)
    if err != nil {
        return nil, err
    }

    return output, nil
}

func StringAsset(conf *Conf, path string) (string, error) {
    bytes, err := ReadAsset(conf, path)
    if err != nil {
        return "", err
    }
    return string(bytes), nil
}
