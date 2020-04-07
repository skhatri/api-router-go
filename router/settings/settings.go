package settings

type RouteSettings interface {
	Variable(name string) *string
	IsToggleOn(name string) bool
	IsToggleOff(name string) bool
	ResponseHeaders() map[string]string
}

type _RouteSettings struct {
	DefaultResponseHeaders map[string]string `json:"response-headers"`
	Toggles                map[string]bool   `json:"toggles"`
	Variables              map[string]string `json:"variables"`
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
