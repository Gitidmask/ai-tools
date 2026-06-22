package main

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"ai_tools/internal/client/db"
	"ai_tools/internal/client/devbridge"
	"ai_tools/internal/client/sidecar"
)

//go:embed frontend/dist
var frontendDist embed.FS

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("AI 工具箱 starting...")

	// Initialize database
	database, err := db.Initialize()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.Close()

	// Find available port
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		log.Fatalf("Failed to find port: %v", err)
	}
	httpPort := listener.Addr().(*net.TCPAddr).Port
	listener.Close()

	// Start DevBridge
	devBridge := devbridge.New(httpPort + 1)
	if err := devBridge.Start(); err != nil {
		log.Printf("[startup] devbridge start failed: %v", err)
	} else {
		log.Printf("[startup] devbridge on 127.0.0.1:%d", devBridge.Port())
	}

	// Start sidecar
	sidecarSvc := sidecar.NewService()
	sidecarSvc.SetDB(database)
	if err := sidecarSvc.Start(devBridge.Port() + 1); err != nil {
		log.Printf("[startup] sidecar start failed: %v", err)
	} else {
		log.Printf("[startup] sidecar started (PID: %d)", sidecarSvc.Status().PID)
	}

	// Build HTTP mux
	mux := http.NewServeMux()

	// API routes
	mux.HandleFunc("/api/health", devbridge.HealthHandler)
	mux.HandleFunc("/api/sidecar-config", devbridge.ConfigHandler)

	// Frontend static files
	mux.Handle("/assets/", serveFrontendAssets())
	mux.HandleFunc("/", serveIndexHTML)

	// Start server
	server := &http.Server{
		Addr:    fmt.Sprintf("127.0.0.1:%d", httpPort),
		Handler: mux,
	}

	go func() {
		log.Printf("[startup] open http://127.0.0.1:%d in your browser", httpPort)
		if err := server.Serve(listener); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP server error: %v", err)
		}
	}()

	// Open browser
	openBrowser(fmt.Sprintf("http://127.0.0.1:%d", httpPort))

	// Graceful shutdown
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh
	log.Println("Shutting down...")
	sidecarSvc.Stop()
	devBridge.Stop()
	server.Close()
}

func serveFrontendAssets() http.Handler {
	// Try embedded assets first
	subFS, err := fs.Sub(frontendDist, "frontend/dist/assets")
	if err == nil {
		return http.FileServer(http.FS(subFS))
	}

	// Fallback to disk
	assetsDir := filepath.Join(".", "frontend", "dist", "assets")
	if _, err := os.Stat(assetsDir); err == nil {
		return http.StripPrefix("/assets/", http.FileServer(http.Dir(assetsDir)))
	}

	return http.NotFoundHandler()
}

func serveIndexHTML(w http.ResponseWriter, r *http.Request) {
	// Try embedded
	data, err := frontendDist.ReadFile("frontend/dist/index.html")
	if err == nil {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write(data)
		return
	}

	// Fallback to disk
	indexPath := filepath.Join(".", "frontend", "dist", "index.html")
	if _, err := os.Stat(indexPath); err == nil {
		http.ServeFile(w, r, indexPath)
		return
	}

	w.Write([]byte(`<html><body><h1>AI 工具箱</h1>
<p>Frontend not built. Run:</p>
<pre>cd frontend && npm install && npm run build</pre>
<p>Then restart this app.</p>
</body></html>`))
}

func openBrowser(url string) {
	// Try default browser on Windows
	proc, err := os.StartProcess("cmd", []string{"cmd", "/c", "start", url},
		&os.ProcAttr{Files: []*os.File{nil, nil, nil}})
	if err == nil {
		proc.Release()
		return
	}
	log.Printf("Open browser at: %s", url)
}
