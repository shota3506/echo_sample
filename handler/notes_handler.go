package handler

import (
	"../model"
	"github.com/labstack/echo"
	"net/http"
)

type noteParam struct {
	Title string
	Content string
	FolderId int `json:"folder_id"`
}

func (h *Handler) GetNote() echo.HandlerFunc {
	return func(c echo.Context) error {
		noteId := c.Param("id")
		note := model.Note{}
		result := h.DB.First(&note, "id=?", noteId)

		if result.Error != nil {
			return c.JSON(http.StatusNotFound, map[string]string{
				"status": "Not Found",
			})
		}
		if result.Error != nil { return h.return404(c, result.Error) }
		return c.JSON(http.StatusOK, struct {
			Note model.Note `json:"note"`
		} {
			Note: note,
		})
	}
}

func (h *Handler) CreateNote() echo.HandlerFunc {
	return func(c echo.Context) error {
		param := new(noteParam)
		if err := c.Bind(param); err != nil { return h.return400(c, err) }
		folder := model.Folder{}
		result := h.DB.First(&folder, "id=?", param.FolderId)
		if result.Error != nil { return h.return404(c, result.Error) }
		note := model.Note{
			Title: param.Title,
			Content: param.Content,
			Folder: folder,
		}

		if err := c.Validate(note); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error":  err.Error(),
			})
		}
		h.DB.Create(&note)

		result = h.DB.Create(&note)
		if result.Error != nil { return h.return400(c, result.Error) }

		return c.JSON(http.StatusOK, struct {
			Note model.Note `json:"note"`
		} {
			Note: note,
		})
	}
}

func (h *Handler) UpdateNote() echo.HandlerFunc {
	return func(c echo.Context) error {
		noteId := c.Param("id")
		note := model.Note{}
		result := h.DB.First(&note, "id=?", noteId)
		if result.Error != nil { return h.return404(c, result.Error) }

		param := new(noteParam)
		if err := c.Bind(param); err != nil { return h.return400(c, err) }

		note.Content = param.Content
		note.Title = param.Title
		if err := c.Validate(note); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error":  err.Error(),
			})
		}
		h.DB.Save(&note)

		return c.JSON(http.StatusOK, struct {
			Note model.Note `json:"note"`
		} {
			Note: note,
		})
	}
}
