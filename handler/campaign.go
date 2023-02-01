package handler

import (
	"DevGo/campaign"
	"DevGo/helper"
	"DevGo/user"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type campaignHandler struct {
	service campaign.Service
}

func NewCampaignHandler(service campaign.Service) *campaignHandler {
	return &campaignHandler{service}
}

func (h *campaignHandler) GetCampaigns(c *gin.Context) {
	userID, _ := strconv.Atoi(c.Query("user_id"))

	campaigns, err := h.service.GetCampaigns(userID)

	if err != nil {
		response := helper.APIResponse("Error to get campaigns", http.StatusBadRequest, "error", nil)

		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIResponse("List of campaigns", http.StatusOK, "success", campaign.FormatCampaigns(campaigns))

	c.JSON(http.StatusOK, response)

}

func (h *campaignHandler) GetCampaign(c *gin.Context) {
	var input campaign.GetCampaignDetailInput

	err := c.ShouldBindUri(&input)

	if err != nil {
		response := helper.APIResponse("Failed to get detail of campaign", http.StatusBadRequest, "error", nil)

		c.JSON(http.StatusBadGateway, response)
		return
	}

	campaignDetail, err := h.service.GetCampaignById(input)

	if err != nil {
		response := helper.APIResponse("Failed to get detail of campaign", http.StatusBadRequest, "error", nil)

		c.JSON(http.StatusBadGateway, response)
		return
	}

	response := helper.APIResponse("Detail Campaign", http.StatusOK, "success", campaign.FormatCampaignDetail(campaignDetail))

	c.JSON(http.StatusOK, response)
}

func (h *campaignHandler) CreateCampaign(c *gin.Context) {
	var input campaign.CreateCampaignInput

	err := c.ShouldBindJSON(&input)
	if err != nil {

		errors := helper.FormatValidationError(err)

		errorMessage := gin.H{"errors": errors}

		response := helper.APIResponse("Create campaign failed", http.StatusUnprocessableEntity, "error", errorMessage)

		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	currentUser := c.MustGet("currentUser").(user.User)
	input.User = currentUser

	newCampaign, err := h.service.CreateCampaign(input)
	if err != nil {
		response := helper.APIResponse("Create campaign failed", http.StatusBadRequest, "error", nil)

		c.JSON(http.StatusBadRequest, response)
		return
	}

	formatter := campaign.FormatCampaign(newCampaign)

	response := helper.APIResponse("Campaign has been created", http.StatusOK, "Success", formatter)

	c.JSON(http.StatusOK, response)
}

func (h *campaignHandler) UpdateCampaign(c *gin.Context) {
	var inputID campaign.GetCampaignDetailInput

	err := c.ShouldBindUri(&inputID)

	if err != nil {
		response := helper.APIResponse("Failed to update campaign", http.StatusBadRequest, err.Error(), nil)

		c.JSON(http.StatusBadGateway, response)
		return
	}

	currentUser := c.MustGet("currentUser").(user.User)

	var input campaign.CreateCampaignInput
	input.User = currentUser

	err = c.ShouldBindJSON(&input)
	if err != nil {

		errors := helper.FormatValidationError(err)

		errorMessage := gin.H{"errors": errors}

		response := helper.APIResponse("Failed to update campaign", http.StatusUnprocessableEntity, err.Error(), errorMessage)

		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	updateCampaign, err := h.service.UpdateCampaign(inputID, input)
	if err != nil {
		response := helper.APIResponse("Failed to update campaign", http.StatusBadRequest, err.Error(), nil)

		c.JSON(http.StatusBadRequest, response)
		return
	}

	formatter := campaign.FormatCampaign(updateCampaign)

	response := helper.APIResponse("Campaign has been updated", http.StatusOK, "Success", formatter)

	c.JSON(http.StatusOK, response)
}

func (h *campaignHandler) UploadImage(c *gin.Context) {
	var input campaign.CreateCampaignImageInput
	err := c.ShouldBind(&input)
	if err != nil {
		response := helper.APIResponse("Failed to upload campaign image", http.StatusUnprocessableEntity, err.Error(), nil)

		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	currentUser := c.MustGet("currentUser").(user.User)
	input.User = currentUser

	file, err := c.FormFile("file")
	if err != nil {
		data := gin.H{"is_uploaded": false}

		response := helper.APIResponse("Failed Upload image !", http.StatusBadRequest, "error", data)

		c.JSON(http.StatusBadRequest, response)
		return

	}

	path := fmt.Sprintf("images/%s", file.Filename)

	err = c.SaveUploadedFile(file, path)
	if err != nil {
		data := gin.H{"is_uploaded": false}

		response := helper.APIResponse("Failed Upload imagae!", http.StatusBadRequest, "error", data)

		c.JSON(http.StatusBadRequest, response)
		return
	}
	_, err = h.service.SaveCampaignImage(input, path)

	if err != nil {
		data := gin.H{"is_uploaded": false}

		response := helper.APIResponse("Failed Upload iamage!", http.StatusBadRequest, err.Error(), data)

		c.JSON(http.StatusBadRequest, response)
		return
	}
	data := gin.H{"is_uploaded": true}

	response := helper.APIResponse("Upload image success !", http.StatusOK, "Success", data)

	c.JSON(http.StatusOK, response)
}