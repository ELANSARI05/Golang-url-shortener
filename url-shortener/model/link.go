package model

import (
	"url-shortener/db"
)

type ShortLink struct {
	ID          int
	UserID      int
	OriginalURL string
	ShortSlug   string
	ClickCount  int
	CreatedAt   string
}

// GetLinksByUserID returns all links for a given user
func GetLinksByUserID(userID int) ([]ShortLink, error) {
	rows, err := db.DB.Query(`
    SELECT id, user_id, original_url, short_slug, click_count, CAST(created_at AS DATETIME)
    FROM short_links
    WHERE user_id = ?
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var links []ShortLink
	for rows.Next() {
		var link ShortLink
		if err := rows.Scan(
			&link.ID,
			&link.UserID,
			&link.OriginalURL,
			&link.ShortSlug,
			&link.ClickCount,
			&link.CreatedAt,
		); err != nil {
			return nil, err
		}
		links = append(links, link)
	}

	return links, nil
}

// Create Link(insert new link into db)
func CreateLink(userID int, originalURL, shortSlug string) error {
	_, err := db.DB.Exec(`
        INSERT INTO short_links (user_id, original_url, short_slug)
        VALUES (?, ?, ?)
    `, userID, originalURL, shortSlug)
	return err
}

// DeleteLinkByID deletes a link by its ID
func DeleteLinkByID(userID, linkID int) error {
	_, err := db.DB.Exec(`
        DELETE FROM short_links
        WHERE id = ? AND user_id = ?
    `, linkID, userID)
	return err
}

// GetOriginalURLBySlug retrieves the original URL by its slug
func GetOriginalURLBySlug(slug string) (string, error) {
	var url string
	err := db.DB.QueryRow(`
        SELECT original_url FROM short_links WHERE short_slug = ?
    `, slug).Scan(&url)
	return url, err
}

// SlugExists checks if a slug already exists in the database(only one slug in the whole db)
func SlugExists(slug string) (bool, error) {
	var exists bool
	err := db.DB.QueryRow(`
		SELECT EXISTS(SELECT 1 FROM short_links WHERE short_slug = ?)
	`, slug).Scan(&exists)
	return exists, err
}
