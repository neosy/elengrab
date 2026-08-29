package deferred

func (m *migrations) initMigrations() {
	m.Add("02_fill_media_info", m.fillMediaInfo)
	m.Add("03_fill_thumbnails", m.fillThumbnails)
	m.Add("04_fill_title_for_instagram", m.fillTitleForInstagram)
	m.Add("05_fill_media_description", m.fillMediaDescription)
	m.Add("06_fill_media_duration", m.fillMediaDuration)
}
