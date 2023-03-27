package proxy

import (
	"crypto/tls"
	"fmt"
	"golang.org/x/net/proxy"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

/**
*简洁版本用GET方法访问指定的URL，并将页面结果以字符串格式返回
 */
func HttpGet(url string) (string, error) {
	result := ""
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("访问url:"+url+"异常：", err.Error())
		return result, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("访问url:"+url+"异常：", err.Error())
		return result, err
	}
	result = string(body)
	return result, nil
}

func HttpProxyGet(url string, proxy string) (string, error) {
	if proxy == "" {
		return HttpGet(url)
	} else {
		if strings.Contains(proxy, "Socks5://") {
			proxy = strings.Replace(proxy, "Socks5://", "", -1)
			return HttpProxySocks5GetTimeout(url, &proxy, 30)
		}
	}
	return HttpProxyGetTimeout(url, &proxy, 30)
}

func HttpProxySocks5GetTimeout(url_str string, proxy_addr *string, timeOut int64) (string, error) {
	timeout := time.Duration(timeOut) * time.Second

	dialer, err := proxy.SOCKS5("tcp", *proxy_addr, nil, proxy.Direct)
	if err != nil {
		return "", err
	}
	//tr := &http.Transport{}
	//tr.Dial =dialer.Dial

	tr := &http.Transport{
		TLSClientConfig:   &tls.Config{InsecureSkipVerify: true},
		DisableKeepAlives: true,
		Dial:              dialer.Dial,
	}

	httpClient := &http.Client{}
	httpClient.Transport = tr
	httpClient.Timeout = timeout

	req, _ := http.NewRequest("GET", url_str, nil)
	res, err := httpClient.Do(req)
	if res != nil {
		defer res.Body.Close()
	} else {
		fmt.Println("访问url:"+url_str+"异常：", err.Error())
		return "", err
	}
	body, _ := ioutil.ReadAll(res.Body)
	retult := string(body)
	return retult, nil
}

func HttpProxyGetTimeout(url_str string, proxy_addr *string, timeOut int64) (string, error) {

	httpClient := &http.Client{}
	url_i := url.URL{}
	url_proxy, _ := url_i.Parse(*proxy_addr)

	timeout := time.Duration(timeOut) * time.Second

	tr := &http.Transport{
		Proxy: http.ProxyURL(url_proxy),
	}
	if strings.Contains(url_str, "https") {
		tr = &http.Transport{
			Proxy:             http.ProxyURL(url_proxy),
			TLSClientConfig:   &tls.Config{InsecureSkipVerify: true},
			DisableKeepAlives: true,
			Dial: func(netw, addr string) (net.Conn, error) {
				c, err := net.DialTimeout(netw, addr, timeout) //设置建立连接超时
				if err != nil {
					return nil, err
				}
				c.SetDeadline(time.Now().Add(timeout)) //设置发送接收数据超时
				return c, nil
			},
		}
	}

	httpClient.Transport = tr
	httpClient.Timeout = timeout

	req, _ := http.NewRequest("GET", url_str, nil)
	res, err := httpClient.Do(req)
	if res != nil {
		defer res.Body.Close()
	} else {
		fmt.Println("访问url:"+url_str+"异常：", err.Error())
		return "", err
	}
	body, _ := ioutil.ReadAll(res.Body)
	retult := string(body)
	return retult, nil
}

func HttpProxyPost(url string, proxy string, data io.Reader, contentType, authorization string) (string, error) {
	if proxy == "" {
		return HttpPost_(url, data, contentType, authorization)
	} else {
		if strings.Contains(proxy, "Socks5://") {
			proxy = strings.Replace(proxy, "Socks5://", "", -1)
			return HttpProxySocks5PostTimeout(url, &proxy, 30, data, contentType, authorization)
		}
	}
	return HttpProxyPostTimeout(url, &proxy, 30, data, contentType, authorization)
}
func HttpPost_(url string, data io.Reader, contentType, authorization string) (string, error) {
	httpClient := &http.Client{}
	tr := &http.Transport{
		DisableKeepAlives: true,
		TLSClientConfig:   &tls.Config{InsecureSkipVerify: false},
		Dial: func(netw, addr string) (net.Conn, error) {
			c, err := net.DialTimeout(netw, addr, time.Second*5) //设置建立连接超时
			if err != nil {
				return nil, err
			}
			c.SetDeadline(time.Now().Add(5 * time.Second)) //设置发送接收数据超时
			return c, nil
		},
	}
	httpClient.Transport = tr
	req, _ := http.NewRequest("POST", url, data)
	req.Header.Set("Content-Type", contentType)
	if authorization != "" {
		req.Header.Set("Authorization", "Bearer "+authorization)
	}
	res, err := httpClient.Do(req)
	if err == nil {
		defer res.Body.Close()
	} else {
		return "", err
	}
	body, _ := ioutil.ReadAll(res.Body)
	retult := string(body)
	return retult, nil
}

func HttpProxySocks5PostTimeout(url_str string, proxy_addr *string, timeOut int64, data io.Reader, contentType, authorization string) (string, error) {
	timeout := time.Duration(timeOut) * time.Second

	dialer, err := proxy.SOCKS5("tcp", *proxy_addr, nil, proxy.Direct)
	if err != nil {
		return "", err
	}
	//tr := &http.Transport{}
	//tr.Dial =dialer.Dial

	tr := &http.Transport{
		TLSClientConfig:   &tls.Config{InsecureSkipVerify: true},
		DisableKeepAlives: true,
		Dial:              dialer.Dial,
	}

	httpClient := &http.Client{}
	httpClient.Transport = tr
	httpClient.Timeout = timeout

	req, _ := http.NewRequest("POST", url_str, data)
	req.Header.Set("Content-Type", contentType)
	if authorization != "" {
		req.Header.Set("Authorization", "Bearer "+authorization)
	}
	res, err := httpClient.Do(req)
	if res != nil {
		defer res.Body.Close()
	} else {
		fmt.Println("访问url:"+url_str+"异常：", err.Error())
		return "", err
	}
	body, _ := ioutil.ReadAll(res.Body)
	retult := string(body)
	return retult, nil
}

func HttpProxyPostTimeout(url_str string, proxy_addr *string, timeOut int64, data io.Reader, contentType, authorization string) (string, error) {

	httpClient := &http.Client{}
	url_i := url.URL{}
	url_proxy, _ := url_i.Parse(*proxy_addr)

	timeout := time.Duration(timeOut) * time.Second

	tr := &http.Transport{
		Proxy: http.ProxyURL(url_proxy),
	}
	if strings.Contains(url_str, "https") {
		tr = &http.Transport{
			Proxy:             http.ProxyURL(url_proxy),
			TLSClientConfig:   &tls.Config{InsecureSkipVerify: true},
			DisableKeepAlives: true,
			Dial: func(netw, addr string) (net.Conn, error) {
				c, err := net.DialTimeout(netw, addr, timeout) //设置建立连接超时
				if err != nil {
					return nil, err
				}
				c.SetDeadline(time.Now().Add(timeout)) //设置发送接收数据超时
				return c, nil
			},
		}
	}

	httpClient.Transport = tr
	httpClient.Timeout = timeout

	req, _ := http.NewRequest("POST", url_str, data)
	req.Header.Set("Content-Type", contentType)
	if authorization != "" {
		req.Header.Set("Authorization", "Bearer "+authorization)
	}
	res, err := httpClient.Do(req)
	if res != nil {
		defer res.Body.Close()
	} else {
		fmt.Println("访问url:"+url_str+"异常：", err.Error())
		return "", err
	}
	body, _ := ioutil.ReadAll(res.Body)
	retult := string(body)
	return retult, nil
}
