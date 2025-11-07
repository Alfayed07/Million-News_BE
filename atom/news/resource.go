package atom_news

import (
	"errors"
	"strings"
)

func ListUseCase(category string, page, limit int) (ListResponse, error) {
    items, total, err := listNews(category, page, limit)
    if err != nil { return ListResponse{}, err }
    return ListResponse{Items: items, Total: total, Page: page, Limit: limit}, nil
}

func TopUseCase(limit int) ([]NewsItem, error) { return topStories(limit) }
func TrendingUseCase(limit int) ([]NewsItem, error) { return trending(limit) }
func SearchUseCase(q string, page, limit int) (ListResponse, error) {
    if strings.TrimSpace(q) == "" { return ListResponse{Items: []NewsItem{}, Total: 0, Page: page, Limit: limit}, nil }
    items, total, err := searchNews(q, page, limit)
    if err != nil { return ListResponse{}, err }
    return ListResponse{Items: items, Total: total, Page: page, Limit: limit}, nil
}
func DetailUseCase(id int64) (NewsItem, error) {
    if id <= 0 { return NewsItem{}, errors.New("invalid id") }
    return getNewsByID(id)
}

func ListCommentsUseCase(newsID int64, limit int) ([]Comment, error) {
    if newsID <= 0 { return []Comment{}, errors.New("invalid id") }
    return listCommentsForNews(newsID, limit)
}

func AddCommentUseCase(newsID int64, userID *int64, content string) (Comment, error) {
    if newsID <= 0 || strings.TrimSpace(content) == "" {
        return Comment{}, errors.New("invalid payload")
    }
    return insertCommentForNews(newsID, userID, content)
}

func ListCategoriesUseCase() ([]Category, error) {
    return listCategories()
}

func RecordViewUseCase(newsID int64) error {
    if newsID <= 0 { return errors.New("invalid id") }
    return incrementNewsView(newsID)
}
