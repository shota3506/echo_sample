package handler

import (
	"../model"
	"github.com/labstack/echo"
	"net/http"
)

type folderParam struct {
	Title string
	ParentId int
}

func (h *Handler) GetFolder() echo.HandlerFunc {
	return func(c echo.Context) error {
		folderId := c.Param("id")
		folder := model.Folder{}
		result := h.DB.Where("tree_paths.length = ?", 1).Preload("Folders").First(&folder, "id=?", folderId)
		if result.Error != nil {
			return c.JSON(http.StatusNotFound, map[string]string{
				"status": "Not Found",
			})
		}
		return c.JSON(http.StatusOK, folder)
	}
}

func (h *Handler) GetFolders() echo.HandlerFunc {
	return func(c echo.Context) error {
		folders := []model.Folder{}
		result := h.DB.Preload("Folders","tree_paths.length = ?", 1).Preload("Folders").Find(&folders)

		if result.Error != nil {
			return c.JSON(http.StatusNotFound, map[string]string{
				"status": "Not Found",
			})
		}
		return c.JSON(http.StatusOK, folders)
	}
}

func (h *Handler) UpdateFolder() echo.HandlerFunc {
	return func(c echo.Context) error {
		folderId := c.Param("id")
		folder := model.Folder{}
		result := h.DB.First(&folder, "id=?", folderId)

		if result.Error != nil {
			return c.JSON(http.StatusNotFound, map[string]string{
				"status": "Not Found",
			})
		}

		param := new(folderParam)
		if err := c.Bind(param); err != nil {
			return err
		}

		folder.Title = param.Title
		h.DB.Save(&folder)

		return  c.JSON(http.StatusOK, folder)
	}
}

func (h *Handler) CreateFolder() echo.HandlerFunc {
	return func(c echo.Context) error {
		param := new(folderParam)
		if err := c.Bind(param); err != nil {
			return err
		}
		folder := model.Folder{
			Title: param.Title,
		}
		h.DB.Create(&folder)

		parent_tree_paths := []model.TreePath{}
		h.DB.Find(&parent_tree_paths, "descendant_id = ?", param.ParentId)

		for _, parent_tree_path := range parent_tree_paths {
			tree_path := model.TreePath{
				AncestorId: parent_tree_path.AncestorId,
				DescendantId: folder.ID,
				Length: parent_tree_path.Length + 1,
			}
			h.DB.Create(&tree_path)
		}

		return c.JSON(http.StatusOK, folder)
	}
}