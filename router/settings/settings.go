package settings

type RouteSettings interface {
	Variable(name string) *string
	IsToggleOn(name string) bool
	IsToggleOff(name string) bool
	ResponseHeaders() map[string]string
	StaticMappings() map[string]string
}

type _RouteSettings struct {
	DefaultResponseHeaders map[string]string `json:"response-headers"`
	Toggles                map[string]bool   `json:"toggles"`
	Variables              map[string]string `json:"variables"`
	Static                 map[string]string `json:"static"`
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
