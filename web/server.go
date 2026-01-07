package web

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/iotopo/topmq/config"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/acme/autocert"
	"io"
	"io/fs"
	"net"
	"net/http"
	"strings"
	"time"
)

//go:embed all:static
var staticFS embed.FS

func StaticFS(router gin.IRouter) {
	fsys, _ := fs.Sub(staticFS, "static")

	httpFS := http.FS(fsys)

	entries, _ := fs.ReadDir(fsys, ".")
	for _, entry := range entries {
		if entry.IsDir() {
			dirFS, _ := fs.Sub(fsys, entry.Name())
			router.StaticFS("/"+entry.Name(), Dir(http.FS(dirFS)))
		} else {
			// index.html, favicon.ico
			router.StaticFileFS("/"+entry.Name(), entry.Name(), httpFS)
		}
	}
}

// CacheControl 缓存控制中间件
func CacheControl(maxAge time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Cache-Control", "public, max-age="+maxAge.String())
		c.Next()
	}
}

var server *http.Server

func runServer(ctx context.Context) error {
	logrus.Infof("start http server on port %d", config.Conf.Web.Port)
	defer func() {
		if e := recover(); e != nil {
			logrus.Fatalf("start web server panic: %v", e)
		}
	}()
	server = &http.Server{
		Addr:              fmt.Sprintf(":%d", config.Conf.Web.Port),
		Handler:           Router,
		ReadHeaderTimeout: 20 * time.Second,
		ReadTimeout:       200 * time.Second,
		WriteTimeout:      90 * time.Second,
		//BaseContext:       func(_ net.Listener) context.Context { return ctx },
	}

	// BaseContext 上的 cancel 可以通知到一些长连接的请求（sse、websocket）
	// 否则需要有读写操作才会发现连接关闭，会阻塞 server.Shutdown
	server.BaseContext = func(_ net.Listener) context.Context { return ctx }
	conf := config.Conf.Web
	if conf.TLS {
		if conf.AutoCert {
			m := autocert.Manager{
				Prompt: autocert.AcceptTOS,
				//HostPolicy: autocert.HostWhitelist("example1.com", "example2.com"),
				//Cache:      autocert.DirCache("/var/www/.cache"),
				Cache: autocert.DirCache("./tls/.cache"),
			}

			if len(conf.HostPolicy) > 0 {
				m.HostPolicy = autocert.HostWhitelist(conf.HostPolicy...)
			}
			//m.HTTPHandler(nil)
			server.TLSConfig = m.TLSConfig()
			return server.ListenAndServeTLS("", "")
		} else {
			// 可以通过 https 访问
			// 生成 cert.key openssl genrsa -out cert.key 2048
			// 生成 cert.cer openssl req -new -x509 -key cert.key -out cert.cer -days 3650
			if err := initCert(); err != nil {
				return fmt.Errorf("web init cert err: %v", err)
			}
			return server.ListenAndServeTLS(certFile, keyFile)
		}
	} else {
		return server.ListenAndServe()
	}
}

func Run(ctx context.Context) {
	go func() {
		if err := runServer(ctx); err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				logrus.WithError(err).Fatal("failed starting http server")
			}
		}
	}()
	defer func() {
		if server != nil {
			_ = server.Shutdown(context.Background())
		}
	}()
	<-ctx.Done()
}

//var FS stuffbin.FileSystem

var Router *gin.Engine

func createRouter() {
	conf := config.Conf.Web
	if conf.Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	//gin.DefaultWriter = logrus.StandardLogger().WriterLevel(logrus.InfoLevel)
	gin.DefaultErrorWriter = logrus.StandardLogger().WriterLevel(logrus.ErrorLevel)

	Router = gin.Default()

	if conf.PProf {
		// go tool pprof http://localhost:8000/debug/pprof/profile?seconds=30
		// 命令执行完成后会生成CPU性能分析文件保存到本地，并自动进入分析操作的终端界面。
		// top：显示消耗CPU时间最多的函数。
		// list function_name：显示特定函数的详细信息。
		// web：生成火焰图并在浏览器中查看。
		// exit：退出分析操作的终端界面。
		pprof.Register(Router, "/debug/pprof")
	}

	if len(conf.TrustedProxies) > 0 {
		_ = Router.SetTrustedProxies(conf.TrustedProxies)
	}
	Router.Use(gzip.Gzip(gzip.DefaultCompression))

	Router.Use(func(c *gin.Context) {
		sizeLimit := config.Conf.Web.ClientMaxBodySize << 20
		if c.Request.ContentLength > sizeLimit {
			c.AbortWithStatusJSON(http.StatusRequestEntityTooLarge, gin.H{
				"code": "err_request_entity_too_large",
				"args": gin.H{
					"limit": fmt.Sprintf("%dMB", sizeLimit),
				},
			})
		}
	})

	if !conf.DisableCache {
		Router.Use(func(c *gin.Context) {
			if strings.HasSuffix(c.Request.URL.Path, ".css") ||
				strings.HasSuffix(c.Request.URL.Path, ".js") ||
				strings.HasSuffix(c.Request.URL.Path, ".json") ||
				strings.HasSuffix(c.Request.URL.Path, ".png") {
				c.Header("Cache-Control", "public, max-age=31536000") // 1 year
				c.Next()
			}
		})
	}

	if len(conf.AllowOrigins) > 0 {
		Router.Use(cors.New(cors.Config{
			AllowCredentials: true,
			AllowOrigins:     conf.AllowOrigins,
			AllowHeaders: []string{
				"Origin", "Content-Length", "Content-Type", "Authentication",
			},
			AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD"},
			MaxAge:       12 * time.Hour,
		}))
	}

	StaticFS(Router)

	Router.NoRoute(func(c *gin.Context) {
		if c.Request.Method != http.MethodHead && c.Request.Method != http.MethodGet {
			return
		}

		if strings.Contains(c.Request.URL.Path, "/api/v1/") || strings.HasPrefix(c.Request.URL.Path, "/assets") {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}

		c.AbortWithStatus(http.StatusOK)
		reader, _ := staticFS.Open("static/index.html")
		io.Copy(c.Writer, reader)
	})

	Router.GET("api/v1/config", func(c *gin.Context) {
		SuccessResponse(c, gin.H{
			"metrics":   config.Conf.Metrics.Enabled,
			"appName":   config.AppName,
			"version":   config.Version,
			"minPwdLen": config.Conf.MinPwdLen,
		})
	})
}
