package jackdbus

import (
	"github.com/godbus/dbus"
	"github.com/godbus/dbus/introspect"
)

const (
	DBusPatchbayInterface       = "org.jackaudio.JackPatchbay"
	DBusControlInterface        = "org.jackaudio.JackControl"
	DBusSessionManagerInterface = "org.jackaudio.SessionManager"
	DBusConfigureInterface      = "org.jackaudio.Configure"
	DBusServiceName             = "org.jackaudio.service"
	DBusRootPath                = dbus.ObjectPath("/org/jackaudio/Controller")
)

// Client connects to JACK server.
type Client struct {
	conn        *dbus.Conn
	eventRouter *eventRouter
}

func New() (*Client, error) {
	eventRouter := newEventRouter()
	conn, err := dbus.SessionBusPrivateHandler(dbus.NewDefaultHandler(), eventRouter)
	eventRouter.conn = conn
	if err != nil {
		return nil, err
	}
	if err = conn.Auth(nil); err != nil {
		conn.Close()
		return nil, err
	}
	if err = conn.Hello(); err != nil {
		conn.Close()
		return nil, err
	}

	return &Client{conn, eventRouter}, nil
}

// Introspect returns introspection info about JACK server object
func (client *Client) Introspect() (*introspect.Node, error) {
	return introspect.Call(client.conn.Object(DBusServiceName, DBusRootPath))
}

// Control is a proxy for org.jackaudio.JackControl interface
func (client *Client) Control() *Control {
	obj := object{client.conn.Object(DBusServiceName, DBusRootPath)}
	return &Control{
		object:      obj,
		iface:       DBusControlInterface,
		eventRouter: client.eventRouter,
	}
}

// Patchbay is a proxy for org.jackaudio.JackPatchbay interface
func (client *Client) Patchbay() *Patchbay {
	obj := object{client.conn.Object(DBusServiceName, DBusRootPath)}
	return &Patchbay{
		object:      obj,
		iface:       DBusPatchbayInterface,
		eventRouter: client.eventRouter,
	}
}

// SessionManager is a proxy for org.jackaudio.SessionManager interface
func (client *Client) SessionManager() *SessionManager {
	obj := object{client.conn.Object(DBusServiceName, DBusRootPath)}
	return &SessionManager{
		object:      obj,
		iface:       DBusSessionManagerInterface,
		eventRouter: client.eventRouter,
	}
}

// Configure is a proxy for org.jackaudio.Configure interface
func (client *Client) Configure() *Configure {
	obj := object{client.conn.Object(DBusServiceName, DBusRootPath)}
	return &Configure{
		object: obj,
		iface:  DBusConfigureInterface,
	}
}

// Close closes client's connection to dbus
func (client *Client) Close() error {
	return client.conn.Close()
}
