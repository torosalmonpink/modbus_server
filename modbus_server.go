package main

import (
	"context"
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/torosalmonpink/mbserver"
)

func main() {
	// Configure logging with custom timestamp format
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006/01/02 15:04:05.000000",
	})

	// Define command line options
	var address string
	var port int
	pflag.StringVarP(&address, "address", "a", "0.0.0.0", "Listening address")
	pflag.IntVarP(&port, "port", "p", 502, "Listening port")
	pflag.Parse()

	// Set up and start the server
	srv := mbserver.NewServer()
	srv.ConnectionAcceptedEvent = append(srv.ConnectionAcceptedEvent, onConnectionAccepted)
	srv.ConnectionClosedEvent = append(srv.ConnectionClosedEvent, onConnectionClosed)
	srv.RequestReceivedEvent = append(srv.RequestReceivedEvent, onRequestReceived)
	srv.ResponseSentEvent = append(srv.ResponseSentEvent, onResponseSent)
	srv.ServerStartedEvent = append(srv.ServerStartedEvent, onServerStarted)
	srv.ServerStoppedEvent = append(srv.ServerStoppedEvent, onServerStopped)
	err := srv.ListenTCP(fmt.Sprintf("%s:%d", address, port))
	if err != nil {
		logrus.Fatalf("Failed to start server: %v", err)
	}
	defer srv.Close()

	// Set up a context for cancellation
	ctx, cancel := context.WithCancel(context.Background())

	// Handle SIGINT and SIGTERM to cleanly stop the server
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sig
		cancel()
	}()

	<-ctx.Done()
}

// onConnectionAccepted logs when a client connects to the server.
func onConnectionAccepted(conn net.Conn) {
	if conn != nil {
		logrus.WithField("remote_addr", conn.RemoteAddr()).Info("Client connected")
	}
}

// onConnectionClosed logs when a client disconnects from the server.
func onConnectionClosed(conn net.Conn) {
	if conn != nil {
		logrus.WithField("remote_addr", conn.RemoteAddr()).Info("Disconnected")
	}
}

// onRequestReceived logs when a request is received from a client.
func onRequestReceived(conn io.ReadWriteCloser, framer mbserver.Framer) {
	entry := logrus.WithFields(logrus.Fields{
		"function": framer.GetFunction(),
		"data":     framer.GetData(),
	})

	if tcpConn, ok := conn.(net.Conn); ok {
		entry = entry.WithField("remote_addr", tcpConn.RemoteAddr())
	}

	entry.Info("Request received")
}

// onResponseSent logs when a response is sent to a client.
func onResponseSent(conn io.ReadWriteCloser, framer mbserver.Framer) {
	entry := logrus.WithFields(logrus.Fields{
		"function": framer.GetFunction(),
		"data":     framer.GetData(),
	})

	if tcpConn, ok := conn.(net.Conn); ok {
		entry = entry.WithField("remote_addr", tcpConn.RemoteAddr())
	}

	entry.Info("Response sent")
}

// onServerStarted logs when the server starts listening for connections.
func onServerStarted(listener net.Listener) {
	if listener != nil {
		logrus.WithField("address", listener.Addr()).Info("Server started")
	}
}

// onServerStopped logs when the server stops listening for connections.
func onServerStopped(listener net.Listener) {
	if listener != nil {
		logrus.WithField("address", listener.Addr()).Info("Server stopped")
	}
}
