package controllers

// Context - контекст контроллера
type Context struct {
	// Handler http.Handler
	// ctx     *gin.Context
	// Version config.Version
	// App     *config.AppConfig
}

// service service.UsersService
func NewContext() *Context {
	return &Context{
		// usersService: service,
	}
}
