package gnome

// DBus service names for GNOME
const (
	// Shell is the GNOME Shell service
	Shell = "org.gnome.Shell"
	// Mutter is the Mutter window manager service
	Mutter = "org.gnome.Mutter"
	// SessionManager is the GNOME session manager service
	SessionManager = "org.gnome.SessionManager"
	// Notifications is the desktop notifications service
	Notifications = "org.freedesktop.Notifications"
	// DBus is the DBus service
	DBus = "org.freedesktop.DBus"
	// FileManager is the GNOME file manager service
	FileManager = "org.gnome.Nautilus"
	// Screenshot is the GNOME screenshot service
	Screenshot = "org.gnome.Screenshot"
	// Settings is the GNOME settings service
	Settings = "org.gnome.Settings"
	// MediaPlayer is the MPRIS media player service
	MediaPlayer = "org.mpris.MediaPlayer2"
	// ShellIntrospect is the GNOME Shell Introspect service
	ShellIntrospect = "org.gnome.Shell"
)

// DBus object paths for GNOME
const (
	// ShellPath is the GNOME Shell object path
	ShellPath = "/org/gnome/Shell"
	// MutterPath is the Mutter window manager object path
	MutterPath = "/org/gnome/Mutter"
	// SessionManagerPath is the GNOME session manager object path
	SessionManagerPath = "/org/gnome/SessionManager"
	// NotificationsPath is the desktop notifications object path
	NotificationsPath = "/org/freedesktop/Notifications"
	// DBusPath is the DBus object path
	DBusPath = "/org/freedesktop/DBus"
	// FileManagerPath is the GNOME file manager object path
	FileManagerPath = "/org/gnome/Nautilus"
	// ScreenshotPath is the GNOME screenshot object path
	ScreenshotPath = "/org/gnome/Screenshot"
	// SettingsPath is the GNOME settings object path
	SettingsPath = "/org/gnome/Settings"
	// MediaPlayerPath is the MPRIS media player object path
	MediaPlayerPath = "/org/mpris/MediaPlayer2"
	// ShellIntrospectPath is the GNOME Shell Introspect object path
	ShellIntrospectPath = "/org/gnome/Shell/Introspect"
)

// DBus interfaces for GNOME
const (
	// ShellInterface is the GNOME Shell interface
	ShellInterface = "org.gnome.Shell"
	// MutterInterface is the Mutter window manager interface
	MutterInterface = "org.gnome.Mutter"
	// SessionManagerInterface is the GNOME session manager interface
	SessionManagerInterface = "org.gnome.SessionManager"
	// NotificationsInterface is the desktop notifications interface
	NotificationsInterface = "org.freedesktop.Notifications"
	// DBusInterface is the DBus interface
	DBusInterface = "org.freedesktop.DBus"
	// FileManagerInterface is the GNOME file manager interface
	FileManagerInterface = "org.gnome.Nautilus"
	// ScreenshotInterface is the GNOME screenshot interface
	ScreenshotInterface = "org.gnome.Screenshot"
	// SettingsInterface is the GNOME settings interface
	SettingsInterface = "org.gnome.Settings"
	// MediaPlayerInterface is the MPRIS media player interface
	MediaPlayerInterface = "org.mpris.MediaPlayer2"
	// MediaPlayerPlayerInterface is the MPRIS media player player interface
	MediaPlayerPlayerInterface = "org.mpris.MediaPlayer2.Player"
	// ShellIntrospectInterface is the GNOME Shell Introspect interface
	ShellIntrospectInterface = "org.gnome.Shell.Introspect"
)

// Window manager DBus service names
const (
	// WindowManager is the window manager service
	WindowManager = "org.gnome.Shell"
	// WindowManagerPath is the window manager object path
	WindowManagerPath = "/org/gnome/Shell/WindowManager"
	// WindowManagerInterface is the window manager interface
	WindowManagerInterface = "org.gnome.Shell.WindowManager"
)

// Application launcher DBus service names
const (
	// AppLauncher is the application launcher service
	AppLauncher = "org.gnome.Shell"
	// AppLauncherPath is the application launcher object path
	AppLauncherPath = "/org/gnome/Shell"
	// AppLauncherInterface is the application launcher interface
	AppLauncherInterface = "org.gnome.Shell"
)

// Clipboard DBus service names
const (
	// Clipboard is the clipboard service
	Clipboard = "org.gnome.Shell"
	// ClipboardPath is the clipboard object path
	ClipboardPath = "/org/gnome/Shell"
	// ClipboardInterface is the clipboard interface
	ClipboardInterface = "org.gnome.Shell"
)
