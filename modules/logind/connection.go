// SPDX-License-Identifier: GPL-3.0-or-later

//go:build linux
// +build linux

package logind

import (
	"context"
	"time"

	"github.com/coreos/go-systemd/v22/login1"
	"github.com/godbus/dbus/v5"
)

type logindConnection interface {
	Close()

	ListSessions() ([]login1.Session, error)
	GetSessionProperties(dbus.ObjectPath) (map[string]dbus.Variant, error)

	ListUsers() ([]login1.User, error)
	GetUserProperty(dbus.ObjectPath, string) (dbus.Variant, error)
}

func newLogindConnection(timeout time.Duration) (logindConnection, error) {
	conn, err := login1.New()
	if err != nil {
		return nil, err
	}
	dbusConn, err := dbus.SystemBus()
	if err != nil {
		return nil, err
	}
	return &logindDBusConnection{
		conn:     conn,
		dbusConn: dbusConn,
		timeout:  timeout,
	}, nil
}

type logindDBusConnection struct {
	conn     *login1.Conn
	dbusConn *dbus.Conn
	timeout  time.Duration
}

func (c *logindDBusConnection) Close() {
	if c.conn != nil {
		c.conn.Close()
		c.conn = nil
	}

	if c.dbusConn != nil {
		_ = c.dbusConn.Close()
		c.dbusConn = nil
	}
}

func (c *logindDBusConnection) ListSessions() ([]login1.Session, error) {
	return c.conn.ListSessions()
}

func (c *logindDBusConnection) GetSessionProperties(path dbus.ObjectPath) (map[string]dbus.Variant, error) {
	return c.getProperties(path, "org.freedesktop.login1.Session")
}

func (c *logindDBusConnection) GetSessionProperty(path dbus.ObjectPath, property string) (dbus.Variant, error) {
	return c.getProperty(path, "org.freedesktop.login1.Session", property)
}

func (c *logindDBusConnection) ListUsers() ([]login1.User, error) {
	return c.conn.ListUsers()
}

func (c *logindDBusConnection) GetUserProperties(path dbus.ObjectPath) (map[string]dbus.Variant, error) {
	return c.getProperties(path, "org.freedesktop.login1.User")
}

func (c *logindDBusConnection) GetUserProperty(path dbus.ObjectPath, property string) (dbus.Variant, error) {
	return c.getProperty(path, "org.freedesktop.login1.User", property)
}

func (c *logindDBusConnection) getProperties(path dbus.ObjectPath, dbusInterface string) (map[string]dbus.Variant, error) {
	obj := c.dbusConn.Object("org.freedesktop.login1", path)

	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	var props map[string]dbus.Variant

	err := obj.CallWithContext(ctx, "org.freedesktop.DBus.Properties.GetAll", 0, dbusInterface).Store(&props)
	if err != nil {
		return nil, err
	}

	return props, nil
}

func (c *logindDBusConnection) getProperty(path dbus.ObjectPath, dbusInterface, property string) (dbus.Variant, error) {
	obj := c.dbusConn.Object("org.freedesktop.login1", path)

	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	var prop dbus.Variant

	err := obj.CallWithContext(ctx, "org.freedesktop.DBus.Properties.Get", 0, dbusInterface, property).Store(&prop)
	if err != nil {
		return dbus.Variant{}, err
	}

	return prop, nil
}
