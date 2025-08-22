func (s *Storage) Close() error {
	if s.db != nil {
		return s.db.Close()
	}
	return nil
}