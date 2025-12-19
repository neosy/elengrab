package sqliterep

func (r *Repositories) Backup(path string) error {
	_, err := r.db.Exec(`
        VACUUM INTO ?
    `, path)
	return err
}
