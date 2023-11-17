package ssh_one

//type item interface {
//	Close()
//}
//
//type newFn func() (*Client, error)
//
//type pool struct {
//	mux   sync.RWMutex
//	items map[string]*Client
//	newFn func(server *Server) (*Client, error)
//}
//
//func (p *pool) Get(s Server) (*Client, error) {
//	p.mux.RLock()
//	if c, ok := p.items[p.key(s)]; ok {
//		p.mux.RUnlock()
//		return c, nil
//	}
//	p.mux.RUnlock()
//	p.mux.Lock()
//	defer p.mux.Unlock()
//	c, err := p.newFn(&s)
//	if err != nil {
//		return nil, err
//	}
//	p.items[p.key(s)] = c
//	return c, nil
//}
//
//func (p *pool) key(s Server) string {
//	return fmt.Sprintf("%s:%s@%s:%d", s.User, s.Password, s.Host, s.Port)
//}
