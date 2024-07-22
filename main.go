package main

import (
	"gopkg.in/yaml.v2"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
)

type RouteConfig struct {
	Name    string `yaml:"name"`
	Pattern string `yaml:"pattern"`
	Target  string `yaml:"target"`
}

type ServerConfig struct {
	Address string `yaml:"address"`
}

type Config struct {
	Routes []RouteConfig `yaml:"routes"`
	Server ServerConfig  `yaml:"server"`
}

type Proxy struct {
	PortMapping map[string]string
}

func (s Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	urlParts := strings.SplitN(r.URL.Path, "/", 3)
	if len(urlParts) < 2 {
		log.Println("url path error", r.URL.Path)
		http.Error(w, "url path error", http.StatusInternalServerError)
		return
	}
	if targetHost, ok := s.PortMapping[urlParts[1]]; ok {
		// 创建反向代理
		urlPath := ""
		if len(urlParts) > 2 {
			urlPath += "/" + urlParts[2]
		}
		urlPath = targetHost + urlPath
		log.Println("[", urlParts[1], "] rq:", r.URL.Path, " to:", urlPath)
		targetURL, err := url.Parse(urlPath)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		proxy := &httputil.ReverseProxy{Director: func(req *http.Request) {
			req.URL.Scheme = targetURL.Scheme
			req.URL.Host = targetURL.Host
			req.Host = targetURL.Host
			req.URL.Path = targetURL.Path
			req.URL.RawQuery = r.URL.RawQuery
			req.Header = r.Header
			req.Body = r.Body
			req.ContentLength = r.ContentLength

		}}
		proxy.ServeHTTP(w, r)
	} else {
		log.Println("[", urlParts[1], "] http proxy no routes name")
		http.NotFound(w, r)
	}
	return
}

func loadConfig(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

func main() {
	// 加载配置
	config, err := loadConfig("config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
		return
	}

	// 初始化路由映射
	portMapping := make(map[string]string, len(config.Routes))
	for _, route := range config.Routes {
		portMapping[route.Pattern] = route.Target
	}

	// 在 6006 端口启动服务器
	log.Println("Starting server on ", config.Server.Address)
	if err := http.ListenAndServe(config.Server.Address, Proxy{PortMapping: portMapping}); err != nil {
		log.Fatal(err)
	}
}
