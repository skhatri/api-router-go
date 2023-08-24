package settings

type TlsConfig struct {
	Enabled    *bool  `json:"enabled,omitempty"`
	PublicKey  string `json:"public-key,omitempty"`
	PrivateKey string `json:"private-key,omitempty"`
}
type TransportSettings struct {
	Port int        `json:"port"`
	Tls  *TlsConfig `json:"tls,omitempty"`
}

type RouteSettings interface {
	Variable(name string) *string
	IsToggleOn(name string) bool
	IsToggleOff(name string) bool
	ResponseHeaders() map[string]string
	StaticMappings() map[string]string
	Transport() TransportSettings
	GetStringOption(name string) string
	GetBoolOption(name string, defaultValue bool) bool
}

type _RouteSettings struct {
	DefaultResponseHeaders  map[string]string      `json:"response-headers"`
	Toggles                 map[string]bool        `json:"toggles"`
	Variables               map[string]string      `json:"variables"`
	Static                  map[string]string      `json:"static"`
	TransportSettingsConfig TransportSettings      `json:"transport"`
	Options                 map[string]interface{} `json:"options"`
}

func (rs *_RouteSettings) Variable(name string) *string {
	value, ok := rs.Variables[name]
	if ok {
		return &value
	}
	return nil
}

func (rs *_RouteSettings) IsToggleOn(name string) bool {
	return rs.toggleState(name)
}

func (rs *_RouteSettings) IsToggleOff(name string) bool {
	return !rs.toggleState(name)
}

func (rs *_RouteSettings) toggleState(name string) bool {
	value, ok := rs.Toggles[name]
	if ok {
		return value
	}
	return false
}

func (rs *_RouteSettings) ResponseHeaders() map[string]string {
	return rs.DefaultResponseHeaders
}

func (rs *_RouteSettings) StaticMappings() map[string]string {
	return rs.Static
}

func (rs *_RouteSettings) Transport() TransportSettings {
	return rs.TransportSettingsConfig
}

func (rs *_RouteSettings) GetStringOption(name string) string {
	value, ok := rs.Options[name]
	if ok {
		return value.(string)
	}
	return ""
}

func (rs *_RouteSettings) GetBoolOption(name string, defaultValue bool) bool {
	value, ok := rs.Options[name]
	if ok {
		return value.(bool)
	}
	return defaultValue
}

type EmptySettings struct {
}

func (empty *EmptySettings) Variable(string) *string {
	return nil
}
func (empty *EmptySettings) IsToggleOn(string) bool {
	return false
}
func (empty *EmptySettings) IsToggleOff(string) bool {
	return true
}
func (empty *EmptySettings) ResponseHeaders() map[string]string {
	return nil
}
func (empty *EmptySettings) StaticMappings() map[string]string {
	return nil
}

func (empty *EmptySettings) Transport() TransportSettings {
	return TransportSettings{
		Port: 6100,
	}
}
func (empty *EmptySettings) GetStringOption(string) string {
	return ""
}

func (empty *EmptySettings) GetBoolOption(_ string, defaultValue bool) bool {
	return defaultValue
}
